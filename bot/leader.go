package bot

import (
	"bytes"
	"ci6ndex/ci6ndex/generated"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	md "github.com/nao1215/markdown"
)

func (b *Bot) handleManageLeaders() handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		slog.Info("handleManageLeaders", "path", bid.CustomID())

		if e.GuildID() == nil {
			slog.Error("missing guild ID")
			return errors.New("missing guild id on event")
		}
		gid, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Join(err, errors.New("failed to parse guild ID"))
		}
		// Create form for managing leaders
		flags := discord.MessageFlagIsComponentsV2
		flags = flags.Add(discord.MessageFlagEphemeral)
		r, err := b.leadersScreen(gid)
		if err != nil {
			return err
		}

		if err := e.UpdateMessage(discord.MessageUpdate{
			Components: &r,
		}); err != nil {
			slog.Error("Failed to create leaders screen", "error", slog.Any("err", err))
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
		}
		return nil
	}
}

func (b *Bot) handleLeaderDetails() handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		lid, err := strconv.ParseInt(e.Vars["leaderId"], 10, 64)
		if err != nil {
			return errors.Join(err, errors.New("failed to parse leaderId from event"))
		}
		slog.Info("handleLeaderDetails", slog.Uint64("lid", uint64(lid)), slog.String("path", bid.CustomID()))

		if e.GuildID() == nil {
			slog.Error("missing guild ID ", slog.Any("err", e))
			return errors.New("missing guild id on event")
		}
		guildId, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Join(err, errors.New("failed to parse guild ID"))
		}

		// Get leader details from database
		leader, err := b.Ci6ndex.GetLeaderById(guildId, uint64(lid))
		if err != nil {
			return errors.Join(err, errors.New(fmt.Sprintf("failed to fetch leader with ID %d", lid)))
		}

		// Create components to display leader details
		components, err := b.leaderDetailsScreen(leader, guildId)
		if err != nil {
			return err
		}

		if err := e.UpdateMessage(discord.MessageUpdate{
			Components: &components,
		}); err != nil {
			slog.Error("Failed to create leader details screen", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
		}
		return nil
	}
}

func (b *Bot) leaderDetailsScreen(leader generated.Leader, guildId uint64) ([]discord.LayoutComponent, error) {
	me, _ := b.Client.Caches.SelfUser()

	var detailsBuffer bytes.Buffer
	err := renderLeaderDetails(&detailsBuffer, leader)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to render leader details"))
	}

	layout := []discord.LayoutComponent{
		discord.NewContainer().AddComponents(
			discord.NewSection(
				discord.NewTextDisplay(detailsBuffer.String()),
			).WithAccessory(discord.NewThumbnail(me.EffectiveAvatarURL())),
			discord.NewLargeSeparator(),
			discord.NewActionRow().WithComponents(
				discord.NewPrimaryButton("Back", "/leaders").WithEmoji(discord.ComponentEmoji{
					Name: backArrow,
				}),
				discord.NewSecondaryButton("Edit Tier", fmt.Sprintf("/leaders/%d/edit", leader.ID)),
			)).
			WithAccentColor(colorSuccess),
	}

	return layout, nil
}

func renderLeaderDetails(buffer io.Writer, leader generated.Leader) error {
	md := md.NewMarkdown(buffer)

	emoji := leader.DiscordEmojiString.String
	if emoji == "" {
		emoji = "👑"
	}

	// Create a more detailed leader profile
	err := md.H1(emoji + " " + leader.LeaderName + " of " + leader.CivName).
		PlainText("**Tier**: " + fmt.Sprintf("%.1f", leader.Tier)).
		LineBreak().
		LineBreak().
		PlainText("**Status**: " + getBannedStatus(leader.Banned)).
		LineBreak().
		LineBreak().
		PlainText("This leader represents the " + leader.CivName + " civilization in Civilization VI.").
		LineBreak().
		LineBreak().
		PlainText("Each leader in the game has unique abilities and bonuses that affect gameplay.").
		LineBreak().
		LineBreak().
		PlainText("**ID**: " + fmt.Sprintf("%d", leader.ID)).
		Build()

	if err != nil {
		return errors.Join(err, errors.New("failed to build leader details markdown"))
	}

	return nil
}

func getBannedStatus(banned bool) string {
	if banned {
		return "❌ Banned from drafts"
	}
	return "✅ Available for drafts"
}

func (b *Bot) leadersScreen(guildId uint64) ([]discord.LayoutComponent, error) {
	me, _ := b.Client.Caches.SelfUser()

	var leadersHeader bytes.Buffer
	err := renderLeadersMainScreen(&leadersHeader)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to draft card"))
	}
	var start uint64 = 1
	var end uint64 = 10
	leads, err := b.Ci6ndex.GetLeadersInRange(guildId, start, end)
	if err != nil {
		return nil, errors.Join(err, errors.New(fmt.Sprintf("failed to fetch leaders in range %d - %d", start, end)))
	}

	leaderRows := make([]discord.ContainerSubComponent, len(leads))
	for i, leader := range leads {
		var leaderRow bytes.Buffer
		err := md.NewMarkdown(&leaderRow).
			H1(leader.DiscordEmojiString.String + fmt.Sprintf(" %s ", leader.LeaderName)).
			Build()
		if err != nil {
			return nil, errors.Join(err, errors.New(fmt.Sprintf("failed to build leader: %v", leader)))
		}

		leaderRoute := fmt.Sprintf("/leaders/%d", leader.ID)

		leaderRows[i] = discord.NewSection(
			discord.NewTextDisplay(
				leaderRow.String(),
			),
		).WithAccessory(discord.NewSecondaryButton(
			"Details",
			leaderRoute,
		))
	}

	layout := []discord.LayoutComponent{
		discord.NewContainer().AddComponents(
			discord.NewSection(
				discord.NewTextDisplay(leadersHeader.String()),
			).WithAccessory(discord.NewThumbnail(me.EffectiveAvatarURL())),
			discord.NewLargeSeparator()).
			AddComponents(leaderRows...).
			AddComponents(
				discord.NewLargeSeparator(),
				discord.NewActionRow().WithComponents(
					discord.NewPrimaryButton("Back", "/draft").WithEmoji(discord.ComponentEmoji{
						Name: backArrow,
					}),
				)).
			WithAccentColor(colorSuccess),
	}

	return layout, nil
}

func renderLeadersMainScreen(header io.Writer) error {
	err := renderLeadersHeader(header)
	if err != nil {
		return errors.Join(err, errors.New("failed to render leaders header"))
	}
	return nil
}

func renderLeadersHeader(header io.Writer) error {
	return md.NewMarkdown(header).H1("Browse Civs").
		PlainText("View statistics about each leader, and optionally update your personal rank of them.").
		Build()
}
