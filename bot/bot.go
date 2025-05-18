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

func New(config Config, c *ci6ndex.Ci6ndex, log log.Logger) *Bot {
	return &Bot{
		Config:  config,
		Logger:  log,
		Ci6ndex: c,
	}
}

func Configure(c *Bot, r handler.Router) error {
	c.Logger.Info("Configuring Discord Bot...")
	var err error
	c.Client, err = disgo.New(c.Config.DiscordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages),
			gateway.WithCompress(true),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("loading..."),
			),
		),
		bot.WithEventListenerFunc(c.OnReady),
		bot.WithEventListeners(r),
		bot.WithLogger(slog.New(&c.Logger)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create discord client")
	}
	return nil
}

func Start(c *Bot) error {
	c.Logger.Info("Starting Bot...")

	if err := c.Client.OpenGateway(context.Background()); err != nil {
		c.Logger.Errorf("Failed to connect to discord gateway: %s", err)
		return err
	}
	return nil
}

func GracefulShutdown(c *Bot) {
	c.Logger.Info("Shutting down Bot...")
	c.Client.Close(context.Background())
	c.Ci6ndex.Close()
}

func (c *Bot) OnReady(_ *events.Ready) {
	c.Logger.Info("Bot is ready! Listening for new events...")
	err := c.Client.SetPresence(context.Background(),
		gateway.WithListeningActivity("Ian and Alex arguing"),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		c.Logger.Errorf("Failed to set presence: %s", err)
	}
}

func (c *Bot) SyncCommands() {
	c.Logger.Info("Syncing commands...")
	ids := strings.Split(c.Config.GuildIds, ",")

	for _, id := range ids {
		_, err := c.Client.Rest.SetGuildCommands(
			snowflake.MustParse(c.Config.BotApplicationID),
			snowflake.MustParse(id),
			Commands,
		)
		if err != nil {
			c.Logger.Errorf("Failed to sync commands: %v", err)
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
