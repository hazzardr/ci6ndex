package discord

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

var (
	RollCivs = &discordgo.ApplicationCommand{
		Name:        "roll-civs",
		Description: "Rolls civs for a draft, if there is one.",
	}
)

func GetRollCivsCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, 0)
	cmds = append(cmds, RollCivs)
	return cmds
}

func GetRollCivsHandlers() map[string]CommandHandler {
	handlers := make(map[string]CommandHandler)
	handlers[RollCivs.Name] = rollCivsHandler
	return handlers
}

func rollCivsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name)
	
}
