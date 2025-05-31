package bot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

func HandleTest(b *Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		switch data := e.Data.(type) {
		case discord.SlashCommandInteractionData:
			flags := discord.MessageFlagIsComponentsV2
			if ephemeral, ok := data.OptBool("ephemeral"); !ok || ephemeral {
				flags = flags.Add(discord.MessageFlagEphemeral)
			}
			if err := e.CreateMessage(discord.MessageCreate{
				Flags: flags,
				Components: []discord.LayoutComponent{
					discord.NewContainer(
						discord.NewSection(
							discord.NewTextDisplay("test test!"),
						).WithAccessory(discord.NewPrimaryButton("Test Button", "test-button").
							WithEmoji(discord.ComponentEmoji{
								ID: snowflake.MustParse("1229388680745975868"),
							}),
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
				b.Logger.Error(errorDescription(err))
			}
		}
		return nil
	}
}
