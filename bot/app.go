package main

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"os"
	"time"
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

type App struct {
	db     *DatabaseOperations
	config *AppConfig
	logger *log.Logger
}

func InitNewApp() (*App, error) {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}
	db, err := newDBConnection(config.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	return &App{
		db:     db,
		config: config,
		logger: logger,
	}, nil
}
