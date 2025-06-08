package bot

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	md "github.com/nao1215/markdown"
	"io"
	"log/slog"
	"strconv"
)

func (b *Bot) HandleManageLeaders() handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		slog.Info("HandleManageLeaders", "path", bid.CustomID(), slog.Any("vars", e.Vars))

		lid, err := strconv.ParseInt(e.Vars["leaderId"], 10, 64)
		if err == nil {
			return b.HandleLeaderDetails(uint64(lid))(bid, e)
		}
		if e.GuildID() == nil {
			slog.Error("missing guild ID ", slog.Any("err", err))
			return errors.New("missing guild id on event")
		}
		gid, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Join(err, errors.New("failed to parse guild ID"))
		}
		// Create form for ranking leaders
		flags := discord.MessageFlagIsComponentsV2
		flags = flags.Add(discord.MessageFlagEphemeral)
		r, err := b.rankScreen(gid)
		if err != nil {
			return err
		}

		if err := e.UpdateMessage(discord.MessageUpdate{
			Components: &r,
		}); err != nil {
			slog.Error("Failed to create rankings screen", "error", slog.Any("err", err))
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
		}
		return nil
	}
}

func (b *Bot) HandleLeaderDetails(lid uint64) handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		slog.Info("HandleLeaderDetails", slog.Uint64("lid", lid), slog.String("path", bid.CustomID()))

		if e.GuildID() == nil {
			slog.Error("missing guild ID ", slog.Any("err", e))
			return errors.New("missing guild id on event")
		}
		_, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Join(err, errors.New("failed to parse guild ID"))
		}

		flags := discord.MessageFlagIsComponentsV2
		flags = flags.Add(discord.MessageFlagEphemeral)
		if err != nil {
			return err
		}

		err = e.DeferUpdateMessage()
		if err != nil {
			return err
		}

		//if err := e.UpdateMessage(discord.MessageUpdate{
		//	Components: &r,
		//}); err != nil {
		//	slog.Error("Failed to create rankings screen", "error", err)
		//	desc, ok := errorDescription(err)
		//	if ok {
		//		slog.Error(desc)
		//	}
		//}
		return nil
	}
}

func (b *Bot) rankScreen(guildId uint64) ([]discord.LayoutComponent, error) {
	me, _ := b.Client.Caches.SelfUser()

	var ranksHeader bytes.Buffer
	err := renderRankMainScreen(&ranksHeader)
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
				discord.NewTextDisplay(ranksHeader.String()),
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

func renderRankMainScreen(header io.Writer) error {
	err := renderRankHeader(header)
	if err != nil {
		return errors.Join(err, errors.New("failed to render ranks header"))
	}
	return nil
}

func renderRankHeader(header io.Writer) error {
	return md.NewMarkdown(header).H1("Browse Civs").
		PlainText("View statistics about each leader, and optionally update your personal rank of them.").
		Build()
}
