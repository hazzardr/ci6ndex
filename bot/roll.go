package bot

import (
	"ci6ndex/ci6ndex"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func HandleConfirmRollDraft(b *Bot) handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		b.Logger.Info("HandleConfirmRollDraft")
		err := e.DeferCreateMessage(true)
		if err != nil {
			b.Logger.Error("Failed to defer message", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				b.Logger.Error(desc)
			}
			return err
		}

		guild, err := parseGuildId(e.GuildID().String())
		players, err := b.Ci6ndex.GetPlayersFromActiveDraft(guild)
		var playerIds = make([]int64, len(players))
		for i, player := range players {
			playerIds[i] = player.ID
		}
		var rules = make([]ci6ndex.Rule, 0)
		rules = append(rules, &ci6ndex.MinTierRule{MinTier: 3})
		rules = append(rules, &ci6ndex.NoOpRule{})
		rules = append(rules, &ci6ndex.NoOpRule{})
		offers, err := b.Ci6ndex.RollForPlayers(
			guild,
			playerIds,
			rules,
		)
		if err != nil {
			b.Logger.Error("Failed to roll for players", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				b.Logger.Error(desc)
			}
			return err
		}
		b.Logger.Info("HandleConfirmRollDraft", "offers", offers)
		rows := make([]discord.ContainerSubComponent, len(offers))
		for i, offer := range offers {
			leaderStr := ""
			for _, leader := range offer.Leaders {
				leaderStr += fmt.Sprintf("%s %s,", leader.DiscordEmojiString.String, leader.LeaderName)
			}
			// Strip final ,
			leaderStr = leaderStr[:len(leaderStr)-1]
			rows[i] = discord.NewTextDisplayf(
				"<@%d>: %s",
				offer.Player.ID, leaderStr,
			)
		}

		layout := []discord.LayoutComponent{
			discord.NewContainer().AddComponents(rows...).WithAccentColor(colorSuccess),
		}
		_, err = e.CreateFollowupMessage(
			discord.MessageCreate{
				Flags:      discord.MessageFlagIsComponentsV2,
				Components: layout,
			},
		)
		if err != nil {
			b.Logger.Error("Failed to create test message", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				b.Logger.Error(desc)
			}
			return err
		}
		return nil
	}
}
