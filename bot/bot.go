package bot

import (
	"bytes"
	"ci6ndex/ci6ndex"
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/json/v2"
	"github.com/disgoorg/snowflake/v2"
	"github.com/pkg/errors"
	"log/slog"
	"strconv"
	"strings"
)

type Bot struct {
	Client       *bot.Client
	Ci6ndex      *ci6ndex.Ci6ndex
	discordToken string
	guildIDs     string
}

func New(c *ci6ndex.Ci6ndex, discordToken, guildIDs string) *Bot {
	return &Bot{
		Ci6ndex:      c,
		discordToken: discordToken,
		guildIDs:     guildIDs,
	}
}

func (b *Bot) Configure() error {
	slog.Info("configuring Discord Bot...")
	r := handler.New()
	r.SlashCommand("/ping", HandlePing)
	r.SlashCommand("/draft", b.HandleManageDraft())
	r.ButtonComponent("/draft", b.HandleManageDraftButton())
	r.ButtonComponent("/create-draft", b.HandleCreateDraft())
	//r.ButtonComponent("/game/latest", HandleViewLatestCompletedGame(b))
	r.ButtonComponent("/confirm-roll", b.HandleConfirmRoll())
	r.ButtonComponent("/confirm-roll-draft", b.HandleConfirmRollDraft())
	r.ButtonComponent("/leaders", b.HandleManageLeaders())
	r.ButtonComponent("/leaders/{leaderId}", b.HandleManageLeaders())
	r.SelectMenuComponent("/select-player", b.HandlePlayerSelect())

	var err error
	b.Client, err = disgo.New(b.discordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages),
			gateway.WithCompress(true),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("loading..."),
			),
		),
		bot.WithEventListenerFunc(b.OnReady),
		bot.WithEventListeners(r),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create discord client")
	}
	return nil
}

func Start(b *Bot) error {
	slog.Info("Starting Bot...")

	if err := b.Client.OpenGateway(context.Background()); err != nil {
		slog.Error("failed to connect to discord gateway", slog.Any("err", err))
		return err
	}
	return nil
}

func GracefulShutdown(b *Bot) {
	slog.Info("Shutting down Bot...")
	b.Client.Close(context.Background())
	b.Ci6ndex.Close()
}

func (b *Bot) OnReady(_ *events.Ready) {
	slog.Info("Bot is ready! Listening for new events...")
	err := b.Client.SetPresence(context.Background(),
		gateway.WithListeningActivity("Counting Tiles Between Cities..."),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		slog.Error("Failed to set presence: ", slog.Any("err", err))
	}
}

func (b *Bot) SyncCommands() error {
	slog.Info("Syncing commands...")
	ids := strings.Split(b.guildIDs, ",")
	var guildIds = make([]snowflake.ID, len(ids))
	for i, id := range ids {
		guildIds[i] = snowflake.MustParse(id)
	}

	err := handler.SyncCommands(b.Client, Commands, guildIds)
	if err != nil {
		var restErr rest.Error
		if errors.As(err, &restErr) {
			if err != nil {
				slog.Error("failed to sync commands", "error", err)
				return err
			}
		} else {
			slog.Error("failed to sync commands", "err", slog.Any("err", err))
			return err

		}
		slog.Error("failed to sync commands: ", slog.Any("err", err))
	}
	slog.Info("Done!")
	return nil
}

func parseGuildId(guild string) (uint64, error) {
	id, err := strconv.ParseUint(guild, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse guild id")
	}
	return id, nil
}

func errorDescription(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	var restErr rest.Error
	if errors.As(err, &restErr) {
		if len(restErr.Errors) == 0 {
			return string(restErr.RsBody), true
		}

		// Pretty format the JSON
		var prettyJSON bytes.Buffer
		if json.Indent(&prettyJSON, restErr.Errors, "", "  ") == nil {
			return prettyJSON.String(), true
		}

		// Fallback to the original string if formatting fails
		return string(restErr.Errors), true
	}
	return err.Error(), true
}
