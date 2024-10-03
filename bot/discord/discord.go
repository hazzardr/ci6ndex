package discord

import (
	"ci6ndex-bot/domain"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	s  *discordgo.Session
	db *domain.DatabaseOperations
}
