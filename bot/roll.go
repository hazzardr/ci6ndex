package bot

import (
	"ci6ndex/ci6ndex"
	"ci6ndex/ci6ndex/generated"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/pkg/errors"
)

var rollCivsCommand = discord.SlashCommandCreate{
	Name:        "roll",
	Description: "Rolls a random set of civilizations",
	//Options: []discord.ApplicationCommandOption{
	//	discord.ApplicationCommandOptionInt{
	//		Name:        "pool-size",
	//		Description: "The number of civilizations a player will have to choose from",
	//		Required:    true,
	//	},
	//	discord.ApplicationCommandOptionBool{
	//		Name:        "randomize",
	//		Description: "Whether or not to randomize picked civs",
	//		Required:    true,
	//	},
	//},
}

func HandleRollCivs(c *Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
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
			AddActionRow(
				discord.ButtonComponent{
					Style:    discord.ButtonStylePrimary,
					Label:    "Confirm Players! \U0001F3B2",
					CustomID: "confirm-roll",
				}).
			SetEphemeral(true).
			Build(),
		)
	}
}

func HandlePlayerSelect(c *Bot) handler.SelectMenuComponentHandler {
	return func(data discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
		c.Logger.Info("event received", "guild", e.GuildID(), "eventId", e.ID())
		users := data.(discord.UserSelectMenuInteractionData)

		guild, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Wrap(err, "failed to parse guild id from event")
		}
		d, err := c.Ci6ndex.GetOrCreateActiveDraft(guild)
		aurl := e.User().EffectiveAvatarURL()
		if err != nil {
			return errors.Wrap(err, "failed to get active draft")
		}
		var players []generated.AddPlayerParams
		for id, user := range users.Resolved.Users {
			gn := ci6ndex.ResolveOptionalString(user.GlobalName)
			av := ci6ndex.ResolveOptionalString(&aurl)
			players = append(players, generated.AddPlayerParams{
				ID:            int64(id),
				Username:      user.Username,
				GlobalName:    gn,
				DiscordAvatar: av,
			})
		}

		errs := c.Ci6ndex.SetPlayersForDraft(guild, d.ID, players)
		if len(errs) > 0 {
			c.Client.Logger.Error("failed to add players to draft", "errors", errs)
			return errors.New("failed to add players to draft")
		}
		return e.DeferUpdateMessage()
	}
}

func HandleConfirmRoll(c *Bot) handler.ButtonComponentHandler {
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
		players, err := c.Ci6ndex.GetPlayersFromActiveDraft(gid)
		if err != nil {
			return errors.Wrap(err, "failed to get players from active draft")
		}

		reRolls, err := c.Ci6ndex.RollForPlayers(gid, 3)
		if err != nil {
			return errors.Wrap(err, "failed to roll for players")
		}

		pf := make([]discord.EmbedField, len(players))
		for i, roll := range reRolls {
			inline := true
			pf[i] = getRollEmbedField(roll, &inline)
		}

		me, _ := c.Client.Caches.SelfUser()
		_, err = e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitle("Rolls:").
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

func getRollEmbedField(offer ci6ndex.Offering, inline *bool) discord.EmbedField {
	var civs string
	for _, leader := range offer.Leaders {
		if leader.DiscordEmojiString.Valid {
			civs += leader.DiscordEmojiString.String + " "
		}
		civs += leader.LeaderName + ": " + leader.CivName + "\n"
	}
	return discord.EmbedField{
		Name:   offer.Player.Username,
		Value:  civs,
		Inline: inline,
	}
}
