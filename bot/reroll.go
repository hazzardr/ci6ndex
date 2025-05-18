package bot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/pkg/errors"
)

var rerollCivsCommand = discord.SlashCommandCreate{
	Name:        "reroll",
	Description: "Re-rolls a random set of civilizations for the given players",
}

func HandleReRollCivs(c *Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		var minPlayers int
		minPlayers = 1
		maxPlayers := 14
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("SelectUsersToReRoll").
			AddActionRow(
				discord.UserSelectMenuComponent{
					CustomID:  "select-reroll-player",
					MinValues: &minPlayers,
					MaxValues: maxPlayers,
				}).
			AddActionRow(
				discord.ButtonComponent{
					Style:    discord.ButtonStylePrimary,
					Label:    "Confirm Reroll! \U0001F3B2",
					CustomID: "confirm-reroll",
				}).
			SetEphemeral(true).
			Build(),
		)
	}
}

func HandleConfirmReRoll(c *Bot) handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		c.Logger.Info("event received", "guild", e.GuildID(), "")
		err := e.DeferCreateMessage(false)
		if err != nil {
			return err
		}
		gid, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Wrap(err, "failed to parse guild id from event")
		}

		rolls, err := c.Ci6ndex.ReRollForPlayers(gid, 3)
		if err != nil {
			return errors.Wrap(err, "failed to roll for players")
		}

		pf := make([]discord.EmbedField, len(rolls))
		for i, roll := range rolls {
			inline := true
			pf[i] = getRollEmbedField(roll, &inline)
		}

		me, _ := c.Client.Caches.SelfUser()
		_, err = e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitle("ReRolls:").
				SetThumbnail(me.EffectiveAvatarURL()).
				AddFields(pf...).
				Build()).
			Build(),
		)
		if err != nil {
			return err
		}
		return err
	}
}
