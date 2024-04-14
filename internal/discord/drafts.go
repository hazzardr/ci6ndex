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
)

func getDraftCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, 0)
	cmds = append(cmds, GetActiveDraftCommand)
	return cmds
}

func getDraftHandlers(db *internal.DatabaseOperations) map[string]CommandHandler {
	handlers := make(map[string]CommandHandler)
	handlers[GetActiveDraftCommand.Name] = getActiveDraftHandler(db)
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