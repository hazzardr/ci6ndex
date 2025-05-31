package bot

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	startDraft,
}

var startDraft = discord.SlashCommandCreate{
	Name:        "draft",
	Description: "Start a civ draft",
}

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Replies with pong",
}
