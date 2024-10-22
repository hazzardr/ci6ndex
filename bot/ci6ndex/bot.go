package ci6ndex

import (
	"ci6ndex-bot/domain"
	"context"
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

type Ci6ndex struct {
	Client bot.Client
	Logger log.Logger
	DB     *domain.DatabaseOperations
	Config AppConfig
}

func New(config AppConfig, db *domain.DatabaseOperations, log log.Logger) *Ci6ndex {
	return &Ci6ndex{
		Config: config,
		Logger: log,
		DB:     db,
	}
}

func (c *Ci6ndex) Configure(r handler.Router) error {
	c.Logger.Info("Configuring Ci6ndex...")
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

func (c *Ci6ndex) Start() error {
	c.Logger.Info("Starting Ci6ndex...")

	if err := c.Client.OpenGateway(context.Background()); err != nil {
		c.Logger.Errorf("Failed to connect to discord gateway: %s", err)
		return err
	}
	return nil
}

func (c *Ci6ndex) GracefulShutdown() {
	c.Logger.Info("Shutting down Ci6ndex...")
	c.Client.Close(context.Background())
	c.DB.Close()
}

func (c *Ci6ndex) OnReady(_ *events.Ready) {
	c.Logger.Info("Ci6ndex ready! Listening for new events...")
	err := c.Client.SetPresence(context.Background(),
		gateway.WithListeningActivity("Ian and Alex arguing"),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	)
	if err != nil {
		c.Logger.Errorf("Failed to set presence: %s", err)
	}
}

func (c *Ci6ndex) SyncCommands() {
	c.Logger.Info("Syncing commands...")
	ids := strings.Split(c.Config.GuildIds, ",")

	for _, id := range ids {
		_, err := c.Client.Rest().SetGuildCommands(
			snowflake.MustParse(c.Config.BotApplicationID),
			snowflake.MustParse(id),
			Commands,
		)
		if err != nil {
			c.Logger.Errorf("Failed to sync commands: %v", err)
		}
	}
}

func parseGuildId(guild string) (uint64, error) {
	id, err := strconv.ParseUint(guild, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse guild id")
	}
	return id, nil
}
