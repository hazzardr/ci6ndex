package ci6ndex

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var rollCivsCommand = discord.SlashCommandCreate{
	Name:        "roll",
	Description: "Rolls a random set of civilizations",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "pool-size",
			Description: "The number of civilizations a player will have to choose from",
			Required:    true,
		},
		discord.ApplicationCommandOptionBool{
			Name:        "randomize",
			Description: "Whether or not to randomize picked civs",
			Required:    true,
		},
	},
}

func HandleRollCivs(c *Ci6ndex) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		_ = e.SlashCommandInteractionData()

		var minPlayers int
		minPlayers = 1
		maxPlayers := 14
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Add Users to draft").
			AddActionRow(
				discord.UserSelectMenuComponent{
					CustomID:  "select-player",
					MinValues: &minPlayers,
					MaxValues: maxPlayers,
				}).
			SetEphemeral(true).
			Build(),
		)
	}
}

func HandlePlayerSelect(c *Ci6ndex) handler.SelectMenuComponentHandler {
	return func(data discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
		c.Logger.Info("player selected")
		_, err := e.UpdateInteractionResponse(
			discord.NewMessageUpdateBuilder().
				SetContent("player selected!").
				Build(),
		)
		if err != nil {
			return err
		}
		return nil
	}
}
