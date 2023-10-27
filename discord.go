package main

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func messageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	slog.Info("event received", "channel", e.ChannelID, "content", e.Content)
	if e.Author.ID == s.State.User.ID {
		return
	}
}

func ready(s *discordgo.Session, e *discordgo.Ready) {
	err := s.UpdateGameStatus(0, "!ci6ndex")
	if err != nil {
		slog.Warn("could not update discord status on startup")
	}
	slog.Info("bot initialized and ready to receive events")
}
