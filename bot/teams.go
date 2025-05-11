package bot

import (
	"github.com/disgoorg/disgo/discord"
)

var randomTeamsCommand = discord.SlashCommandCreate{
	Name:        "teams",
	Description: "Create and assign random teams",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "teamSize",
			Description: "The number of players per team",
			Required:    true,
		},
	},
}

//
//func HandleRandomizeTeams(c *Bot) handler.CommandHandler {
//	return func(e *handler.CommandEvent) error {
//		teamSize := e.M("teamSize")
//		if teamSize < 1 {
//			return e.Respond(discord.InteractionResponseTypeChannelMessageWithSource, discord.NewMessageCreateBuilder().SetContent("Team size must be greater than 0").Build())
//		}
