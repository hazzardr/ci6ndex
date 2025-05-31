package bot

import (
	"context"
	"errors"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func HandleManageDraft(b *Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		var guildId uint64
		if e.GuildID() == nil {
			b.Logger.Errorf("missing guild ID %v", e)
			return errors.New("missing guild id on event")
		}
		guildId = uint64(*e.GuildID())

		switch data := e.Data.(type) {
		case discord.SlashCommandInteractionData:
			ctx := context.TODO()
			// Create "form" for rolling settings
			flags := discord.MessageFlagIsComponentsV2
			if ephemeral, ok := data.OptBool("ephemeral"); !ok || ephemeral {
				flags = flags.Add(discord.MessageFlagEphemeral)
			}
			me, _ := b.Client.Caches.SelfUser()
			members, err := b.Ci6ndex.GetPlayers(ctx, guildId)
			if err != nil {
				b.Logger.Error("failed to fetch member", err)
			}
			b.Logger.Info("fetched members: ", members)
			if err := e.CreateMessage(discord.MessageCreate{
				Flags: flags,
				Components: []discord.LayoutComponent{
					discord.NewContainer(
						discord.NewSection(
							discord.NewTextDisplay(`# Ci6ndex Draft Manager
Civ (VI) Index helps manage drafts and stores match history.
`,
							),
						).WithAccessory(discord.NewThumbnail(me.EffectiveAvatarURL())),
						discord.NewLargeSeparator(),
						discord.NewSection(
							discord.NewTextDisplay(`## Recent Games
abc123
`),
						),
						discord.NewActionRow(
							discord.NewUserSelectMenu("user-select", "Select a user").
								WithMinValues(1).
								WithMaxValues(14),
						),
					).WithAccentColor(0x5c5fea),
				},
			}); err != nil {
				b.Logger.Error("Failed to create test message", "error", err)
				desc, ok := errorDescription(err)
				if ok {
					b.Logger.Error(desc)
				}
			}
		}
		return nil
	}
}
