package bot

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	checkLeaders,
	startDraft,
}

var startDraft = discord.SlashCommandCreate{
	Name:        "draft",
	Description: "Start a civ draft",
}

var checkLeaders = discord.SlashCommandCreate{
	Name:        "leaders",
	Description: "Browse available leaders",
}

var getLeader = discord.SlashCommandCreate{
	Name:        "leader",
	Description: "Get details about a specific leader",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "search",
			Description: "Term to search by: Can be partial name of civ or name",
			Required:    true,
		},
	},
}

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Replies with pong",
}
