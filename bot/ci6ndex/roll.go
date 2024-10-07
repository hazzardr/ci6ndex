package ci6ndex

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/pkg/errors"
	"strconv"
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
		guild := e.GuildID()
		if guild == nil {
			c.Logger.Error("unable to parse guild from interaction", "event", e)
			return fmt.Errorf("unable to parse guild from interaction")
		}
		gid, err := strconv.ParseUint(guild.String(), 10, 64)
		if err != nil {
			return errors.Wrapf(err, "failed to parse guild id %s", guild.String())
		}

		d, err := c.DB.GetOrCreateActiveDraft(gid)
		if err != nil {
			return errors.Wrap(err, "failed to get active draft")
		}
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(fmt.Sprintf("draft toime %v", d)).
			SetEphemeral(true).
			Build(),
		)
	}
}
