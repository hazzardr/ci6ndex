package bot

import (
	"ci6ndex/ci6ndex"
	"context"
	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/pkg/errors"
	"log/slog"
	"strconv"
	"strings"
)

type Bot struct {
	Client  *bot.Client
	Logger  log.Logger
	Ci6ndex *ci6ndex.Ci6ndex
	Config  Config
}

func New(config Config, c *ci6ndex.Ci6ndex) *Bot {
	return &Bot{
		Config:  config,
		Logger:  *c.Logger,
		Ci6ndex: c,
	}
}

func Configure(b *Bot, r handler.Router) error {
	b.Logger.Info("Configuring Discord Bot...")

	r.Command("/ping", HandlePing)
	r.Command("/roll", HandleRollCivs(b))
	r.SelectMenuComponent("/select-player", HandlePlayerSelect(b))
	r.ButtonComponent("/confirm-roll", HandleConfirmRoll(b))

	var err error
	b.Client, err = disgo.New(b.Config.DiscordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages),
			gateway.WithCompress(true),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("loading..."),
			),
		),
		bot.WithEventListenerFunc(b.OnReady),
		bot.WithEventListeners(r),
		bot.WithLogger(slog.New(&b.Logger)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create discord client")
	}
	return nil
}

func Start(b *Bot) error {
	b.Logger.Info("Starting Bot...")

	if err := b.Client.OpenGateway(context.Background()); err != nil {
		b.Logger.Errorf("Failed to connect to discord gateway: %s", err)
		return err
	}
	return nil
}

func GracefulShutdown(b *Bot) {
	b.Logger.Info("Shutting down Bot...")
	b.Client.Close(context.Background())
	b.Ci6ndex.Close()
}

func (b *Bot) OnReady(_ *events.Ready) {
	b.Logger.Info("Bot is ready! Listening for new events...")
	err := b.Client.SetPresence(context.Background(),
		gateway.WithListeningActivity("Counting Tiles Between Cities..."),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		b.Logger.Errorf("Failed to set presence: %s", err)
	}
}

func (b *Bot) SyncCommands() {
	b.Logger.Info("Syncing commands...")
	ids := strings.Split(b.Config.GuildIds, ",")

	for _, id := range ids {
		_, err := b.Client.Rest.SetGuildCommands(
			snowflake.MustParse(b.Config.BotApplicationID),
			snowflake.MustParse(id),
			Commands,
		)
		if err != nil {
			b.Logger.Errorf("Failed to sync commands: %v", err)
		}
	}
}

type Config struct {
	DiscordToken     string `env:"DISCORD_API_TOKEN"`
	BotApplicationID string `env:"DISCORD_BOT_APPLICATION_ID"`
	GuildIds         string `env:"GUILD_IDS"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func parseGuildId(guild string) (uint64, error) {
	id, err := strconv.ParseUint(guild, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse guild id")
	}
	return id, nil
}
