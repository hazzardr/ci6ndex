package discord

import (
	"ci6ndex/domain"
	"ci6ndex/internal"
	"context"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgtype"
	"log/slog"
	"time"
)

var (
	CheckActiveDraftCommand = &discordgo.ApplicationCommand{
		Name:        "get-active-draft",
		Description: "Gets details of the active draft, if there is one",
	}

	PromptCreateDraftCommand = &discordgo.ApplicationCommand{
		Name:        "create-draft",
		Description: "Creates a draft, if no active draft already exists.",
	}

	CancelAllActiveDraftsCommand = &discordgo.ApplicationCommand{
		Name:        "cancel-all-active-drafts",
		Description: "Sets all drafts to inactive in the database",
	}

	CreateDraftSelectPlayersId  = "create-draft-select-player"
	CreateDraftSelectStrategyId = "create-draft-select-strategy"
	CreateDraftConfirmId        = "create-draft-confirm-button"
)

func getDraftCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, 0)
	cmds = append(cmds, CheckActiveDraftCommand)
	cmds = append(cmds, PromptCreateDraftCommand)
	cmds = append(cmds, CancelAllActiveDraftsCommand)
	return cmds
}

func getDraftHandlers(db *internal.DatabaseOperations, mb *MessageBuilder) map[string]CommandHandler {
	handlers := make(map[string]CommandHandler)
	handlers[CheckActiveDraftCommand.Name] = getActiveDraftHandler(db, mb)
	handlers[PromptCreateDraftCommand.Name] = createDraftHandler(db)
	handlers[CancelAllActiveDraftsCommand.Name] = handleCancelDraftsCommand(db)
	handlers[CreateDraftSelectPlayersId] = handlePlayerPickerInteraction(db)
	handlers[CreateDraftSelectStrategyId] = handleDraftStrategyPickerInteraction()
	handlers[CreateDraftConfirmId] = handleCreateDraftConfirmInteraction(db, mb)
	return handlers
}

func getActiveDraftHandler(db *internal.DatabaseOperations, mb *MessageBuilder) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name)
		drafts, err := db.Queries.GetActiveDrafts(context.Background())
		if err != nil {
			reportError("An error occured while trying to get the active draft.", err, s, i, false)
			return
		}
		if drafts == nil || len(drafts) == 0 {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "There is no active draft.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				slog.Error("failed to respond to interaction", "i", i.Interaction, "error", err)
			}
			return
		}
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The active draft id is %d", drafts[0].ID),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction, "error", err)
		}
	}
}

func createDraftHandler(db *internal.DatabaseOperations) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name, "interactionId", i.ID)
		minPlayersAllowed := 1
		maxPlayersAllowed := 12
		ctx := context.Background()

		drafts, err := db.Queries.GetActiveDrafts(ctx)
		if err != nil {
			reportError("error checking active drafts", err, s, i, false)
		}
		if len(drafts) > 0 {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "There is already an active draft! Cannot create a new one.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				slog.Error("failed to respond to interaction", "i", i.Interaction, "error", err)
			}
			return
		}
		strategies, err := db.Queries.GetDraftStrategies(ctx)
		if err != nil {
			reportError("could not fetch draft strategies for discord interaciton", err, s, i, false)
		}
		stratMenuOptions := make([]discordgo.SelectMenuOption, 0, len(strategies))
		defaultStrat := "RandomPickNoRepeats"
		for _, s := range strategies {
			if s.Name != defaultStrat {
				stratMenuOptions = append(stratMenuOptions, discordgo.SelectMenuOption{
					Label: s.Name,
					Value: s.Name,
				})
			} else {
				stratMenuOptions = append(stratMenuOptions, discordgo.SelectMenuOption{
					Label:   s.Name,
					Value:   s.Name,
					Default: true,
				})
			}
		}

		d, err := db.Queries.CreateActiveDraft(ctx, defaultStrat)
		if err != nil {
			reportError("error creating draft", err, s, i, false)
		}

		est, err := time.LoadLocation("America/New_York")
		if err != nil {
			reportError("error loading timezone", err, s, i, false)
		}

		// Convert current time to EST timezone
		nowInEST := time.Now().In(est)
		defaultDate := pgtype.Date{
			Valid: true,
			Time:  nowInEST,
		}
		_, err = db.Queries.CreateGameFromDraft(ctx, domain.CreateGameFromDraftParams{
			DraftID:   d.ID,
			StartDate: defaultDate,
		})
		if err != nil {
			reportError("error scaffolding game", err, s, i, false)
		}
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				CustomID: "new-draft",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    CreateDraftSelectPlayersId,
								Placeholder: "Select Members for game",
								MinValues:   &minPlayersAllowed,
								MaxValues:   maxPlayersAllowed,
								MenuType:    discordgo.UserSelectMenu,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								CustomID:    CreateDraftSelectStrategyId,
								Placeholder: "Select Draft Type",
								MenuType:    discordgo.StringSelectMenu,
								Options:     stratMenuOptions,
							},
						},
					},
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								CustomID: CreateDraftConfirmId,
								Style:    discordgo.PrimaryButton,
								Label:    "Start Draft",
							},
						},
					},
				},
			},
		})
		if err != nil {
			reportError("unable to create draft!", err, s, i, false)
		}
	}
}

func handlePlayerPickerInteraction(db *internal.DatabaseOperations) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "messageComponentId", i.Interaction.MessageComponentData().CustomID, "interactionId", i.ID)

		players := make([]string, 0)
		for _, p := range i.Interaction.MessageComponentData().Resolved.Users {
			players = append(players, p.GlobalName)
		}
		_, err := db.Queries.AddPlayersToActiveDraft(context.Background(), players)
		if err != nil {
			reportError("error adding players to draft", err, s, i, false)
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction.ID, "error", err)
		}
	}
}

func handleDraftStrategyPickerInteraction() CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "messageComponentId", i.Interaction.MessageComponentData().CustomID, "interactionId", i.ID)
		_ = i.Interaction.MessageComponentData().Values[0]

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction.ID, "error", err)
		}
	}
}

func handleCreateDraftConfirmInteraction(db *internal.DatabaseOperations, mb *MessageBuilder) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "messageComponentId", i.Interaction.MessageComponentData().CustomID, "interactionId", i.ID)

		drafts, err := db.Queries.GetActiveDrafts(context.Background())
		if err != nil {
			reportError("error getting active draft", err, s, i, false)
		}

		var activeDraft domain.Ci6ndexDraft

		if len(drafts) > 1 || len(drafts) == 0 {
			msg := "error: there should be exactly one active draft"
			reportError(msg, errors.New(msg), s, i, false)
			return
		}

		if len(drafts) == 1 {
			activeDraft = drafts[0]
		}
		if err != nil {
			reportError("error fetching game from active draft", err, s, i, false)
		}

		if err != nil {
			reportError("error finalizing draft", err, s, i, false)
		}
		//displayMessage, err := mb.WriteConfirmDraft(CreateDraftConfirmId, activeDraft.Players)
		_ = activeDraft.Players
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "displayMessage",
			},
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction.ID, "error", err)
		}
	}
}

func handleCancelDraftsCommand(db *internal.DatabaseOperations) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "messageComponentId", i.Interaction.ApplicationCommandData().Name, "interactionId", i.ID)
		_, err := db.Queries.CancelActiveDrafts(context.Background())
		if err != nil {
			reportError("unable to cancel drafts", err, s, i, false)
			return
		}
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Successfully cancelled drafts.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction.ID, "error", err)
		}
	}
}
