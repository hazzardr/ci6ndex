package cmd

import (
	"ci6ndex/bot"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

type ServeCommand struct{}
type SyncCommand struct{}
type Bot struct {
	Serve ServeCommand `cmd:"" help:"Start the Discord Bot"`
	Sync  SyncCommand  `cmd:"" help:"Sync the Discord Bot Application Commands"`
}

func (s *SyncCommand) Run(b *bot.Bot) error {
	return b.SyncCommands()
}

func (s *ServeCommand) Run(b *bot.Bot) error {
	defer bot.GracefulShutdown(b)
	err := bot.Start(b)
	if err != nil {
		return errors.Join(errors.New("failed to start the discord bot"), err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sig
	return nil
}
