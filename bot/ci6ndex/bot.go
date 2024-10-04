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
	"os"
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

func (c *Ci6ndex) Configure() error {
	c.Logger.Info("Configuring Bot...", "guildId", c.Config.GuildId)
	var err error
	if c.Client, err = disgo.New(c.Config.DiscordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages),
			gateway.WithCompress(true),
			gateway.WithPresenceOpts(
				gateway.WithPlayingActivity("loading..."),
			),
		),
		bot.WithEventListenerFunc(c.OnReady),
	); err != nil {
		c.Logger.Errorf("Failed to configure bot: %s", err)
	}

	return nil
}

func (c *Ci6ndex) Start() {
	c.Logger.Info("Starting Bot...")

	if err := c.Client.OpenGateway(context.Background()); err != nil {
		c.Logger.Errorf("Failed to connect to discord gateway: %s", err)
	}
	defer func() {
		c.Logger.Info("Shutting down bot...")
		c.Client.Close(context.Background())
		c.DB.Close()
	}()

	s := make(chan os.Signal, 1)
	<-s
}

func (c *Ci6ndex) OnReady(_ *events.Ready) {
	c.Logger.Info("Bot started! Listening for new events.")
	if err := c.Client.SetPresence(context.Background(),
		gateway.WithListeningActivity("Ian and Alex arguing		"),
		gateway.WithOnlineStatus(discord.OnlineStatusOnline),
	); err != nil {
		c.Logger.Errorf("Failed to set presence: %s", err)
	}
}
