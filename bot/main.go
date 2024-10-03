/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/spf13/viper"
	"log/slog"
)

type AppConfig struct {
	DiscordToken     string `mapstructure:"DISCORD_API_TOKEN"`
	BotApplicationID string `mapstructure:"DISCORD_BOT_APPLICATION_ID"`
	GuildId          string `mapstructure:"CS_GUILD_ID"`
	DatabaseUrl      string `mapstructure:"SQLITE_DB_URL"`
}

func loadConfig() (*AppConfig, error) {
	var config AppConfig
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}
	_, err = InitNewApp(config)
	if err != nil {
		panic(err)
	}
	slog.Info("initialization successful")
}
