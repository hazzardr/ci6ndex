package discord

import (
	"ci6ndex/domain"
	"ci6ndex/internal"
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

var (
	RollCivs = &discordgo.ApplicationCommand{
		Name:        "roll-civs",
		Description: "Rolls civs for a draft, if there is one.",
	}
)

func getRollCivCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, 0)
	cmds = append(cmds, RollCivs)
	return cmds
}

func getRollCivsHandlers(db *internal.DatabaseOperations) map[string]CommandHandler {
	handlers := make(map[string]CommandHandler)

	handlers[RollCivs.Name] = getRollCivsHandler(db)
	return handlers
}

func getRollCivsHandler(db *internal.DatabaseOperations) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name, "interactionId", i.ID)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Rolling Civs. This may take a few seconds!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction, "error", err)
		}
		ctx := context.Background()
		drafts, err := db.Queries.GetActiveDrafts(ctx)
		if err != nil {
			reportError("error checking active drafts", err, s, i, true)
			return
		}

		var activeDraft domain.Ci6ndexDraft

		if len(drafts) == 0 {
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "There is no active draft. Roll will not be attached to a game...",
				Flags:   discordgo.MessageFlagsEphemeral,
			})
			if err != nil {
				slog.Error("error responding to user", "error", err)
			}

			// dummy draft as a default
			activeDraft, err = db.Queries.GetDraft(ctx, -1)
			if err != nil {
				reportError("Failed to roll civs for mock draft.", err, s, i, true)
				return
			}

		}

		if len(drafts) > 1 {
			msg := "There are multiple active drafts. This should not be possible."
			reportError(msg, errors.New(msg), s, i, true)
			return
		}

		if len(drafts) == 1 {
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Rolled civs will be attached to the active draft...",
				Flags:   discordgo.MessageFlagsEphemeral,
			})

			if err != nil {
				slog.Error("error responding to user", "error", err)
			}

			activeDraft = drafts[0]
		}
		leaders, err := db.Queries.GetLeaders(ctx)
		if err != nil {
			reportError("Failed to fetch leaders for draft.", err, s, i, true)
			return
		}

		strat, err := db.Queries.GetDraftStrategy(ctx, activeDraft.DraftStrategy)
		shuffler := internal.NewCivShuffler(leaders, activeDraft.Players, strat, db)
		offers, err := shuffler.Shuffle()
		if err != nil {
			reportError("Error when shuffling civs.", err, s, i, true)
			return
		}

		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("The following picks were rolled: %v", offers),
			Flags:   discordgo.MessageFlagsEphemeral,
		})

		if err != nil {
			slog.Error("error responding to user", "error", err)
			return
		}
		slog.Info("successfully rolled civs", "interactionId", i.ID)
	}
}
