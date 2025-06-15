package bot

import (
	"bytes"
	"errors"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	md "github.com/nao1215/markdown"
	"io"
	"log/slog"
)

func (b *Bot) draftScreen() ([]discord.LayoutComponent, error) {
	me, _ := b.Client.Caches.SelfUser()

	var draftHeader, recentGames bytes.Buffer
	err := renderDraftMainScreen(&draftHeader, &recentGames)
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
			discord.NewPrimaryButton("Leaders", "/leaders").WithEmoji(discord.ComponentEmoji{
				Name: notebook,
			}),
		),
	).WithAccentColor(0x5c5fea),
	}, nil
}

func (b *Bot) handleManageDraft() handler.SlashCommandHandler {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		if e.GuildID() == nil {
			slog.Error("missing guild ID ", slog.Any("vars", e))
			return errors.New("missing guild id on event")
		}
		// Create "form" for rolling settings
		flags := discord.MessageFlagIsComponentsV2
		flags = flags.Add(discord.MessageFlagEphemeral)
		draft, err := b.draftScreen()
		if err != nil {
			return err
		}

		if err := e.CreateMessage(discord.MessageCreate{
			Flags:      flags,
			Components: draft,
		}); err != nil {
			slog.Error("Failed to create test message", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
		}

		return nil
	}
}

func (b *Bot) handleManageDraftButton() handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		slog.Info("handleCreateDraft")

		if e.GuildID() == nil {
			slog.Error("missing guild ID ", slog.Any("vars", e))
			return errors.New("missing guild id on event")
		}
		// Create "form" for rolling settings
		flags := discord.MessageFlagIsComponentsV2
		flags = flags.Add(discord.MessageFlagEphemeral)
		draft, err := b.draftScreen()
		if err != nil {
			return err
		}

		if err := e.UpdateMessage(discord.MessageUpdate{
			Components: &draft,
		}); err != nil {
			slog.Error("Failed to create test message", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
		}
		return nil
	}
}

func (b *Bot) handleCreateDraft() handler.ButtonComponentHandler {
	return func(bid discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		slog.Info("handleCreateDraft")
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
			slog.Error("Failed to create test message", "error", err)
			desc, ok := errorDescription(err)
			if ok {
				slog.Error(desc)
			}
		}
		return nil
	}
}

func renderDraftMainScreen(header, previousGame io.Writer) error {
	err := renderDraftHeader(header)
	if err != nil {
		return errors.Join(err, errors.New("failed to render draft header"))
	}
	err = renderPreviousGameSummary(previousGame)
	if err != nil {
		return errors.Join(err, errors.New("failed to render recent games"))
	}
	return nil
}

func renderDraftHeader(header io.Writer) error {
	return md.NewMarkdown(header).H1("Ci6ndex Draft Manager").
		PlainText("Civ (VI) Index helps manage drafts and stores match history.").
		Build()
}

func renderPreviousGameSummary(output io.Writer) error {
	return md.NewMarkdown(output).H2("Previous Game").
		H3f("**%s Winner:** <:Nzinga_Mbande_Civ6:1229393600790663220> <@135218870494429184>",
			partyEmoji).
		Build()
}
