package bot

import (
	"ci6ndex/ci6ndex"
	"ci6ndex/ci6ndex/generated"
	"context"
	"database/sql"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/pkg/errors"
)

const (
	DefaultPoolSize = 3
)

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

func HandlePlayerSelect(b *Bot) handler.SelectMenuComponentHandler {
	return func(data discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
		b.Logger.Info("event received", "guild", e.GuildID(), "eventId", e.ID())
		users := data.(discord.UserSelectMenuInteractionData)

		guild, err := parseGuildId(e.GuildID().String())
		if err != nil {
			return errors.Wrap(err, "failed to parse guild id from event")
		}
		d, err := b.Ci6ndex.GetOrCreateActiveDraft(guild)
		aurl := e.User().EffectiveAvatarURL()
		if err != nil {
			return errors.Wrap(err, "failed to get active draft")
		}
		var players []generated.AddPlayerParams
		for id, user := range users.Resolved.Users {
			gn := ci6ndex.ResolveOptionalString(user.GlobalName)
			av := ci6ndex.ResolveOptionalString(&aurl)
			params := generated.AddPlayerParams{
				ID:            int64(id),
				Username:      user.Username,
				GlobalName:    gn,
				DiscordAvatar: av,
			}
			_, err := b.Ci6ndex.GetPlayer(context.TODO(), guild, int64(id))
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err := b.Ci6ndex.AddPlayer(context.TODO(), guild, params)
					if err != nil {
						return errors.Wrap(err, "failed to add player")
					}
				}
			}
			players = append(players, params)
		}

		errs := b.Ci6ndex.SetPlayersForDraft(guild, d.ID, players)
		if len(errs) > 0 {
			b.Client.Logger.Error("failed to add players to draft", "errors", errs)
			return errors.New("failed to add players to draft")
		}
		return e.DeferUpdateMessage()
	}
}

func HandleConfirmRoll(c *Bot) handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		c.Logger.Info("event received", "guild", e.GuildID())
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
		var pids = make([]int64, 0)
		for _, player := range players {
			pids = append(pids, player.ID)
		}

		var rules = make([]ci6ndex.Rule, DefaultPoolSize)
		rules[0] = &ci6ndex.MinTierRule{MinTier: 3}
		for i := 1; i < DefaultPoolSize; i++ {
			rules[i] = &ci6ndex.NoOpRule{} // No additional rules for now
		}
		rolls, err := c.Ci6ndex.RollForPlayers(gid, pids, rules)
		if err != nil {
			return errors.Wrap(err, "failed to roll for players")
		}

		pf := make([]discord.EmbedField, len(players))
		for i, roll := range rolls {
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
