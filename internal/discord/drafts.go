package discord

import (
	"ci6ndex/internal"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

var (
	GetActiveDraftCommand = &discordgo.ApplicationCommand{
		Name:        "get-active-draft",
		Description: "Gets details of the active draft, if there is one",
	}

	PromptCreateDraftCommand = &discordgo.ApplicationCommand{
		Name:        "create-draft",
		Description: "Creates a draft, if no active draft already exists.",
	}

	CreateDraftSelectPlayersId  = "create-draft-select-player"
	CreateDraftSelectStrategyId = "create-draft-select-strategy"
	CreateDraftConfirmButtonId  = "create-draft-confirm-button"
)

func getDraftCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, 0)
	cmds = append(cmds, GetActiveDraftCommand)
	cmds = append(cmds, PromptCreateDraftCommand)
	return cmds
}

func getDraftHandlers(db *internal.DatabaseOperations, mb *MessageBuilder) map[string]CommandHandler {
	handlers := make(map[string]CommandHandler)
	handlers[GetActiveDraftCommand.Name] = getActiveDraftHandler(db)
	handlers[PromptCreateDraftCommand.Name] = createDraftHandler(db)
	handlers[CreateDraftSelectPlayersId] = handlePlayerPickerInteraction()
	handlers[CreateDraftSelectStrategyId] = handleDraftStrategyPickerInteraction()
	handlers[CreateDraftConfirmButtonId] = handleCreateDraftButtonInteraction(mb)
	return handlers
}

func getActiveDraftHandler(db *internal.DatabaseOperations) CommandHandler {
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

		strategies, err := db.Queries.GetDraftStrategies(context.Background())
		if err != nil {
			reportError("could not fetch draft strategies for discord interaciton", err, s, i, false)
		}
		stratMenuOptions := make([]discordgo.SelectMenuOption, 0, len(strategies))
		for _, s := range strategies {
			if s.Name != "RandomPickNoRepeats" {
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
								CustomID: CreateDraftConfirmButtonId,
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

func handlePlayerPickerInteraction() CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "messageComponentId", i.Interaction.MessageComponentData().CustomID, "interactionId", i.ID)

		players := make([]string, 0)
		for _, p := range i.Interaction.MessageComponentData().Resolved.Users {
			players = append(players, p.GlobalName)
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

// TODO: update original message with template from mb
func handleCreateDraftButtonInteraction(mb *MessageBuilder) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("event received", "messageComponentId", i.Interaction.MessageComponentData().CustomID, "interactionId", i.ID)
		strategy := i.Interaction.MessageComponentData().Values[0]

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Selected Strategy: %v", strategy),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction.ID, "error", err)
		}
	}
}
