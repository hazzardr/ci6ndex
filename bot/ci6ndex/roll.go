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
		discord.ApplicationCommandOptionUser{
			Name:        "player",
			Description: "The player to roll civilizations for",
			Required:    true,
		},
	},
}

func HandleRollCivs(c *Ci6ndex) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		data := e.SlashCommandInteractionData()

		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContentf("Rolling %d civilizations...", data.Int("pool-size")).
			SetEphemeral(true).
			Build(),
		)
	}
}
