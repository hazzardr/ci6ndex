package bot

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	rollCivsCommand,
	randomTeamsCommand,
	rerollCivsCommand,
}
