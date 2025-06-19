package bot

import (
	"bytes"
	"ci6ndex/ci6ndex/generated"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
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

		var currentOffset uint64 = 1 // Default offset (1-indexed)
		var currentLimit uint64 = 10 // Default items per page

		parsedURL, pErr := url.Parse(bid.CustomID())
		if pErr == nil { // Path could be /leaders or /leaders?offset=...
			queryParams := parsedURL.Query()
			if offsetStr := queryParams.Get("offset"); offsetStr != "" {
				if parsedOffset, errConv := strconv.ParseUint(offsetStr, 10, 64); errConv == nil {
					if parsedOffset < 1 {
						currentOffset = 1 // Ensure offset is at least 1
					} else {
						currentOffset = parsedOffset
					}
				} else {
					slog.Warn("Failed to parse offset from CustomID", "offsetStr", offsetStr, "error", errConv)
				}
			}
			if limitStr := queryParams.Get("limit"); limitStr != "" {
				if parsedLimit, errConv := strconv.ParseUint(limitStr, 10, 64); errConv == nil && parsedLimit > 0 {
					currentLimit = parsedLimit
				} else if errConv != nil {
					slog.Warn("Failed to parse limit from CustomID", "limitStr", limitStr, "error", errConv)
				}
			}
		} else if bid.CustomID() != "/leaders" { // Allow initial call with /leaders
			// Log error only if CustomID is not the base path and parsing failed
			slog.Warn("Failed to parse CustomID for pagination", "customID", bid.CustomID(), "error", pErr)
		}

		r, err := b.leadersScreen(gid, currentOffset, currentLimit)
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

	// Create a more detailed leader profile
	err := md.H1(emoji+" "+leader.LeaderName+" of "+leader.CivName).
		H2f("**Tier**: %.1f", leader.Tier).
		PlainText("\n**Status**: " + getBannedStatus(leader.Banned)).
		Build()

	if err != nil {
		return errors.Join(err, errors.New("failed to build leader details markdown"))
	}

	return nil
}

func getBannedStatus(banned bool) string {
	if banned {
		return "🧑‍⚖️ Banned from drafts"
	}
	return "✅ Available for draft"
}

func (b *Bot) leadersScreen(guildId uint64, currentOffset uint64, limit uint64) ([]discord.LayoutComponent, error) {
	me, _ := b.Client.Caches.SelfUser()

	var leadersHeader bytes.Buffer
	err := renderLeadersMainScreen(&leadersHeader)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to render leaders main screen"))
	}

	// Assuming GetLeadersInRange expects a 1-based offset and a limit (number of items)
	leads, err := b.Ci6ndex.GetLeadersInRange(guildId, currentOffset, limit)
	if err != nil {
		return nil, errors.Join(err, errors.New(fmt.Sprintf("failed to fetch leaders with offset %d, limit %d", currentOffset, limit)))
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
				b.createLeaderPaginationActionRow(currentOffset, limit, len(leads)),
			).
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

func (b *Bot) createLeaderPaginationActionRow(currentOffset, limit uint64, numLeadsFetched int) discord.ActionRowComponent {
	// Previous Button
	prevButtonOffset := currentOffset - limit
	if prevButtonOffset < 1 {
		prevButtonOffset = 1 // Ensure offset doesn't go below 1
	}
	prevButton := discord.NewSecondaryButton("Previous", fmt.Sprintf("/leaders?offset=%d&limit=%d", prevButtonOffset, limit)).
		WithEmoji(discord.ComponentEmoji{Name: "⬅️"})
	if currentOffset <= 1 {
		prevButton = prevButton.WithDisabled(true)
	}

	// Next Button
	nextButtonOffset := currentOffset + limit
	nextButton := discord.NewSecondaryButton("Next", fmt.Sprintf("/leaders?offset=%d&limit=%d", nextButtonOffset, limit)).
		WithEmoji(discord.ComponentEmoji{Name: "➡️"})
	if numLeadsFetched < int(limit) {
		nextButton = nextButton.WithDisabled(true)
	}

	// Back Button (to draft menu)
	backToDraftButton := discord.NewPrimaryButton("Back", "/draft").WithEmoji(discord.ComponentEmoji{
		Name: backArrow,
	})

	return discord.NewActionRow().WithComponents(prevButton, nextButton, backToDraftButton)
}
