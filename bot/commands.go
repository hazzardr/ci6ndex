package bot

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	checkLeaders,
	startDraft,
	getLeader,
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
			Name:        "leader-name",
			Description: "name of leader to search for",
			Required:    false,
		},
		discord.ApplicationCommandOptionString{
			Name:        "civ-name",
			Description: "name of civ to search for",
			Required:    false,
		},
	},
}

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Replies with pong",
}
