package ci6ndex

import (
	"errors"
	"github.com/spf13/viper"
)

type AppConfig struct {
	DiscordToken     string `mapstructure:"DISCORD_API_TOKEN"`
	BotApplicationID string `mapstructure:"DISCORD_BOT_APPLICATION_ID"`
	GuildIds         string `mapstructure:"GUILD_IDS"`
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig

	viper.AutomaticEnv()

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	// Read the .env file if it exists
	err := viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
