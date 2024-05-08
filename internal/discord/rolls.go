package discord

import (
	"ci6ndex/domain"
	"ci6ndex/internal"
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

var (
	RollCivs = &discordgo.ApplicationCommand{
		Name:        "roll-civs",
		Description: "Rolls civs for a draft, if there is one.",
	}
	PickCiv = &discordgo.ApplicationCommand{
		Name:        "pick-civ",
		Description: "Shows a user what civs they can pick from.",
	}
	PickCivSelectorId = "pick-civ-selector"
)

func getRollCivCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, 0)
	cmds = append(cmds, RollCivs)
	cmds = append(cmds, PickCiv)
	return cmds
}

func getRollCivsHandlers(db *internal.DatabaseOperations, mb *MessageBuilder) map[string]CommandHandler {
	handlers := make(map[string]CommandHandler)

	handlers[RollCivs.Name] = getRollCivsHandler(db, mb)
	handlers[PickCiv.Name] = pickCivHandler(db)
	handlers[PickCivSelectorId] = pickCivSelectHandler(db)
	return handlers
}

func getRollCivsHandler(db *internal.DatabaseOperations, mb *MessageBuilder) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name, "interactionId", i.ID)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
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
		slog.Info("offers", offers)
		for _, offer := range offers {
			leaderIDs := make([]int64, len(offer.Leaders))
			for j, l := range offer.Leaders {
				leaderIDs[j] = l.ID
			}
			u, err := db.Queries.GetUserByDiscordName(ctx, offer.User)
			if err != nil {
				reportError("failed to fetch user information required for roll", err, s, i, true)
				return
			}
			param := domain.WriteOfferedParams{
				DraftID: activeDraft.ID,
				UserID:  u.ID,
				Offered: leaderIDs,
			}

			_, err = db.Queries.WriteOffered(ctx, param)
			if err != nil {
				reportError("failed to persist offer to database", err, s, i, true)
				return
			}
		}

		content, err := mb.WriteDraftOfferings(RollCivs.Name, offers)
		if err != nil {
			reportError("Error writing discord message", err, s, i, true)
			return
		}
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: content,
		})

		if err != nil {
			slog.Error("error responding to user", "error", err)
			return
		}
		slog.Info("successfully rolled civs", "interactionId", i.ID)
	}
}

func pickCivHandler(db *internal.DatabaseOperations) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name, "interactionId", i.ID)

		var minValues int
		minValues = 0

		ctx := context.Background()
		drafts, err := db.Queries.GetActiveDrafts(ctx)
		if err != nil {
			reportError("error checking active drafts", err, s, i, true)
			return
		}
		var activeDraft domain.Ci6ndexDraft

		if len(drafts) == 0 {
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
			activeDraft = drafts[0]
		}
		u, err := db.Queries.GetUserByDiscordName(ctx, i.Interaction.Member.User.GlobalName)
		if err != nil {
			reportError("error fetching user", err, s, i, true)
			return
		}
		offers, err := db.Queries.ReadOfferedForUser(ctx, domain.ReadOfferedForUserParams{
			UserID:  u.ID,
			DraftID: activeDraft.ID,
		})
		if err != nil {
			reportError("error fetching offered civs", err, s, i, true)
			return
		}

		leaders := make([]domain.Ci6ndexLeader, 0)
		for _, o := range offers.Offered {
			l, err := db.Queries.GetLeader(ctx, o)
			if err != nil {
				reportError("error fetching leader", err, s, i, true)
				return
			}
			leaders = append(leaders, l)
		}
		options := make([]discordgo.SelectMenuOption, len(leaders))
		for i, l := range leaders {
			trimmed := strings.Trim(l.DiscordEmojiString.String, "<:>")

			parts := strings.Split(trimmed, ":")
			options[i] = discordgo.SelectMenuOption{
				Label: fmt.Sprintf("%s (%s)", l.LeaderName, l.CivName),
				Value: l.LeaderName,
				Emoji: &discordgo.ComponentEmoji{
					Name:     parts[0],
					ID:       parts[1],
					Animated: false,
				},
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				CustomID: "submit-pick",
				Flags:    discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    PickCivSelectorId,
								Placeholder: "Please select your civ",
								MinValues:   &minValues,
								MaxValues:   1,
								MenuType:    discordgo.StringSelectMenu,
								Options:     options,
							},
						},
					},
				},
			},
		})
		if err != nil {
			reportError("error responding to user", err, s, i, true)
			return
		}
	}
}

func pickCivSelectHandler(db *internal.DatabaseOperations) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "command", i.Interaction.MessageComponentData().CustomID,
			"interactionId", i.ID)
		slog.Info("msg", i.Interaction.MessageComponentData())
		//ctx := context.Background()
		//drafts, err := db.Queries.GetActiveDrafts(ctx)
		//if err != nil {
		//	reportError("error checking active drafts", err, s, i, true)
		//	return
		//}
		//var activeDraft domain.Ci6ndexDraft
		//
		//if len(drafts) == 0 {
		//	reportError("No active draft found", errors.New("no active draft found"), s, i, true)
		//}
		//
		//if len(drafts) > 1 {
		//	msg := "There are multiple active drafts. This should not be possible."
		//	reportError(msg, errors.New(msg), s, i, true)
		//	return
		//}
		//
		//if len(drafts) == 1 {
		//	activeDraft = drafts[0]
		//}
		//u, err := db.Queries.GetUserByDiscordName(ctx, i.Interaction.Member.User.Username)
		//if err != nil {
		//	reportError("error fetching user", err, s, i, true)
		//	return
		//}

	}
}
