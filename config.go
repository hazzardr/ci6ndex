package main

import "github.com/caarlos0/env/v11"

type Config struct {
	DiscordToken     string `env:"DISCORD_API_TOKEN"`
	BotApplicationID string `env:"DISCORD_BOT_APPLICATION_ID"`
	GuildIDs         string `env:"GUILD_IDS"`
	ListenToGuildID  string `env:"LISTEN_TO_GUILD_ID"`
}

func loadConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
