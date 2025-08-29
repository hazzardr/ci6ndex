package cmd

import (
	"ci6ndex/bot"
	"errors"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Bot    Bot           `cmd:"" help:"Perform bot actions."`
	Import ImportCommand `cmd:"" help:"Import rankings from google sheets"`
}

type ImportCommand struct{}

func (i *ImportCommand) Run(b *bot.Bot) error {
	return b.SyncRankings()
}

func Exec(b *bot.Bot) error {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("ci6ndex"),
		kong.Description("Ci6ndex Management CLI."),
		kong.UsageOnError(),
	)

	err := ctx.Run(b)
	if err != nil {
		ctx.FatalIfErrorf(errors.Join(err, errors.New("failed to run bot")))
	}
	return nil
}
