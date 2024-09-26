package pkg

type AppConfig struct {
	DiscordToken                   string `mapstructure:"DISCORD_API_TOKEN"`
	DatabaseUrl                    string `mapstructure:"POSTGRES_URL"`
	GoogleCloudCredentialsLocation string `mapstructure:"GCLOUD_CREDS_LOC"`
	CivRankingSheetId              string `mapstructure:"RANKING_SHEET_ID"`
	BotApplicationID               string `mapstructure:"DISCORD_BOT_APPLICATION_ID"`
	GuildId                        string `mapstructure:"CS_GUILD_ID"`
}
