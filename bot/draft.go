package bot

import (
	"bytes"
	"errors"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func draftScreen(b *Bot) ([]discord.LayoutComponent, error) {
	me, _ := b.Client.Caches.SelfUser()

	var draftHeader, recentGames bytes.Buffer
	err := renderMainScreen(&draftHeader, &recentGames)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to draft card"))
	}

	return []discord.LayoutComponent{discord.NewContainer(
		discord.NewSection(
			discord.NewTextDisplay(draftHeader.String()),
		).WithAccessory(discord.NewThumbnail(me.EffectiveAvatarURL())),
		discord.NewLargeSeparator(),
		discord.NewSection().WithComponents(
			discord.NewTextDisplay(recentGames.String()),
		).WithAccessory(
			discord.NewPrimaryButton("Details", "/game/latest").AsDisabled().WithEmoji(discord.ComponentEmoji{
				Name: magnifyingGlass,
			})),
		discord.NewLargeSeparator(),
		discord.NewActionRow(
			discord.NewPrimaryButton("New Draft", "/create-draft").WithEmoji(discord.ComponentEmoji{
				Name: crossedSwords,
			}),
		),
	).WithAccentColor(0x5c5fea),
	}, nil
}

func HandleManageDraft(b *Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		if e.GuildID() == nil {
			b.Logger.Errorf("missing guild ID %v", e)
			return errors.New("missing guild id on event")
		}
		// Create "form" for rolling settings
		flags := discord.MessageFlagIsComponentsV2
		flags = flags.Add(discord.MessageFlagEphemeral)
		draft, err := draftScreen(b)
		if err != nil {
			return err
		}

		if err := e.CreateMessage(discord.MessageCreate{
			Flags:      flags,
			Components: draft,
		}); err != nil {
			b.Logger.Error("Failed to create test message", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				b.Logger.Error(desc)
			}
		}

		return nil
	}
}

func HandleManageDraftButton(b *Bot) handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		b.Logger.Info("HandleCreateDraft")

		if e.GuildID() == nil {
			b.Logger.Errorf("missing guild ID %v", e)
			return errors.New("missing guild id on event")
		}
		// Create "form" for rolling settings
		flags := discord.MessageFlagIsComponentsV2
		flags = flags.Add(discord.MessageFlagEphemeral)
		draft, err := draftScreen(b)
		if err != nil {
			return err
		}

		if err := e.UpdateMessage(discord.MessageUpdate{
			Components: &draft,
		}); err != nil {
			b.Logger.Error("Failed to create test message", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				b.Logger.Error(desc)
			}
		}
		return nil
	}
}

func HandleCreateDraft(b *Bot) handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		b.Logger.Info("HandleCreateDraft")
		err := e.UpdateMessage(
			discord.MessageUpdate{
				Components: &[]discord.LayoutComponent{
					discord.NewContainer(
						discord.NewTextDisplay("## Create a Draft"),
						discord.NewSmallSeparator(),
						discord.NewActionRow().WithComponents(
							discord.NewUserSelectMenu("/select-player", "Select users").
								WithMinValues(2).
								WithMaxValues(12),
						),
						discord.NewActionRow().WithComponents(
							discord.NewPrimaryButton("Back", "/draft").WithEmoji(discord.ComponentEmoji{
								Name: backArrow,
							}),
							discord.NewPrimaryButton("Roll!", "/confirm-roll-draft").WithEmoji(discord.ComponentEmoji{
								Name: crossedSwords,
							}),
						),
					).WithAccentColor(0x5c5fea),
				},
			},
		)
		if err != nil {
			b.Logger.Error("Failed to create test message", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				b.Logger.Error(desc)
			}
		}
		return nil
	}
}
