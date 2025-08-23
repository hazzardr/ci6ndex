package bot

import (
	"bytes"
	"ci6ndex/ci6ndex"
	"ci6ndex/ci6ndex/generated"
	"context"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	json "github.com/disgoorg/json/v2"
	snowflake "github.com/disgoorg/snowflake/v2"
	"github.com/pkg/errors"
)

type Bot struct {
	Client          *bot.Client
	Ci6ndex         *ci6ndex.Ci6ndex
	discordToken    string
	guildIDs        string
	listenToGuildID snowflake.ID
	leadersCache    map[uint64][]generated.Leader
	wg              sync.WaitGroup
}

func New(c *ci6ndex.Ci6ndex, discordToken, guildIDs string, listenToGuildID string) *Bot {
	return &Bot{
		Ci6ndex:         c,
		discordToken:    discordToken,
		guildIDs:        guildIDs,
		listenToGuildID: snowflake.MustParse(listenToGuildID),
		leadersCache:    make(map[uint64][]generated.Leader),
		wg:              sync.WaitGroup{},
	}
}

func (b *Bot) Configure() error {
	slog.Info("configuring Discord Bot...")
	r := handler.New()
	r.SlashCommand("/ping", HandlePing)
	// r.SlashCommand("/leader", b.handleGetLeaderSlashCommand())

	r.Use(FilterGuildMiddleware(b.listenToGuildID))

	r.Group(func(r handler.Router) {
		r.SlashCommand("/draft", b.handleManageDraft())
		r.ButtonComponent("/draft", b.handleManageDraftButton())
		r.ButtonComponent("/create-draft", b.handleCreateDraft())
	})
	r.Group(func(r handler.Router) {
		r.ButtonComponent("/confirm-roll", b.handleConfirmRoll())
		r.ButtonComponent("/confirm-roll-draft", b.handleConfirmRollDraft())
	})
	r.SlashCommand("/leader", b.handleSearchLeaderSlashCommand())
	r.Route("/leaders", func(r handler.Router) {
		// r.Use(middleware.Logger)
		r.SlashCommand("/", b.handleManageLeadersSlashCommand())
		r.ButtonComponent("/", b.handleManageLeadersButtonCommand())
		r.ButtonComponent("/page/{page}", b.handleManageLeadersButtonCommand())
		r.ButtonComponent("/{leaderId}", b.handleLeaderDetailsButtonCommand())
		r.SelectMenuComponent("/{leaderId}/rating", b.handleRateLeaderMenuSelectCommand())
	})

	r.SelectMenuComponent("/select-player", b.handlePlayerSelect())
	//r.ButtonComponent("/game/latest", HandleViewLatestCompletedGame(b))

	var err error
	b.Client, err = disgo.New(b.discordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuilds,
				gateway.IntentDirectMessages,
			),
			gateway.WithCompress(true),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("loading..."),
			),
		),
		bot.WithEventListenerFunc(b.onReady),
		bot.WithEventListeners(r),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create discord client")
	}
	return nil
}

func Start(b *Bot) error {
	slog.Info("Starting Bot...",
		slog.String("guildID", b.listenToGuildID.String()),
		slog.Bool("tokenProvided", b.discordToken != ""),
	)

	if err := b.Client.OpenGateway(context.Background()); err != nil {
		slog.Error("failed to connect to discord gateway", slog.Any("err", err))
		return err
	}
	slog.Info("Successfully connected to Discord gateway")
	return nil
}

func GracefulShutdown(b *Bot) {
	slog.Info("Shutting down Bot...")
	b.Client.Close(context.Background())
	b.Ci6ndex.Close()
}

func (b *Bot) onReady(_ *events.Ready) {
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

func FilterGuildMiddleware(guildID snowflake.ID) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(event *handler.InteractionEvent) error {
			if event.GuildID() == nil {
				slog.Debug("DROP event", "reason", "only serve guild messages", "event", event)
				return nil
			}
			if !(*event.GuildID() == guildID) {
				slog.Info("DROP event", "reason", "guild id does not match this deployment", "allowedGuildID", guildID, "event", event)
				return nil
			}
			return next(event)
		}
	}
}

// background is a convenience function to use a singular waitgroup for bot operations.
// This is helpful when we want to gracefully shut down, as we can wg.wait() with a timeout and ensure
// background tasks are attempted to be finished up
func (b *Bot) background(fn func()) {
	b.wg.Go(fn)
}
