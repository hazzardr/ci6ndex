package ci6ndex

import (
	"ci6ndex-bot/domain"
	"ci6ndex-bot/domain/generated"
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

func HandleRollCivs(c *Ci6ndex) handler.CommandHandler {
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

func HandlePlayerSelect(c *Ci6ndex) handler.SelectMenuComponentHandler {
	return func(data discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
		c.Logger.Info("event received", "guild", e.GuildID(), "eventId", e.ID())
		users := data.(discord.UserSelectMenuInteractionData)

		guild, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Wrap(err, "failed to parse guild id from event")
		}
		d, err := c.DB.GetOrCreateActiveDraft(guild)
		aurl := e.User().EffectiveAvatarURL()
		if err != nil {
			return errors.Wrap(err, "failed to get active draft")
		}
		var players []generated.AddPlayerParams
		for id, user := range users.Resolved.Users {
			gn := domain.ResolveOptionalString(user.GlobalName)
			av := domain.ResolveOptionalString(&aurl)
			players = append(players, generated.AddPlayerParams{
				ID:            int64(id),
				Username:      user.Username,
				GlobalName:    gn,
				DiscordAvatar: av,
			})
		}

		errs := c.DB.SetPlayersForDraft(guild, d.ID, players)
		if len(errs) > 0 {
			c.Client.Logger().Error("failed to add players to draft", "errors", errs)
			return errors.New("failed to add players to draft")
		}
		return e.DeferUpdateMessage()
	}
}

func HandleConfirmRoll(c *Ci6ndex) handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		c.Logger.Info("event received", "guild", e.GuildID(), "")
		gid, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Wrap(err, "failed to parse guild id from event")
		}
		players, err := c.DB.GetPlayersFromActiveDraft(gid)
		if err != nil {
			return errors.Wrap(err, "failed to get players from active draft")
		}

		rolls, err := c.DB.RollForPlayers(gid, 3)
		if err != nil {
			return errors.Wrap(err, "failed to roll for players")
		}

		pf := make([]discord.EmbedField, len(players))
		for i, roll := range rolls {
			inline := true
			pf[i] = getRollEmbedField(roll, &inline)
		}

		me, _ := c.Client.Caches().SelfUser()
		return e.UpdateMessage(discord.NewMessageUpdateBuilder().
			SetEmbeds(discord.NewEmbedBuilder().
				SetTitle("Rolls:").
				SetThumbnail(me.EffectiveAvatarURL()).
				AddFields(pf...).
				Build()).
			ClearContainerComponents().
			Build(),
		)
	}
}

func getRollEmbedField(offer domain.Offering, inline *bool) discord.EmbedField {
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
