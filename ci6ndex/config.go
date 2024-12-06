package ci6ndex

import "github.com/caarlos0/env/v11"

type AppConfig struct {
	DiscordToken     string `env:"DISCORD_API_TOKEN"`
	BotApplicationID string `env:"DISCORD_BOT_APPLICATION_ID"`
	GuildIds         string `env:"GUILD_IDS"`
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
