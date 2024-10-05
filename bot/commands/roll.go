package commands

import (
	"ci6ndex-bot/ci6ndex"
	"github.com/disgoorg/disgo/discord"
)

var rollCivsCommand = discord.SlashCommandCreate{
	Name:        "roll-civs",
	Description: "Rolls a random set of civilizations",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name: "	",
		},
	},
}

func HandleRollCivs(c *ci6ndex.Ci6ndex) error {

}
