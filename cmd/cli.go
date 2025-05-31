package cmd

import (
	"ci6ndex/bot"
	"ci6ndex/ci6ndex"
	"embed"
	"errors"
	"github.com/alecthomas/kong"
	"github.com/disgoorg/disgo/handler"
)

type CLI struct {
	Bot Bot `cmd:"" help:"Perform bot actions."`
}

func Exec(fs embed.FS) {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("ci6ndex"),
		kong.Description("Ci6ndex Management CLI."),
		kong.UsageOnError(),
	)

	config, err := bot.LoadConfig()
	if err != nil {
		ctx.FatalIfErrorf(errors.Join(err, errors.New("failed to load bot config")))
	}
	c, err := ci6ndex.New(fs)
	if err != nil {
		ctx.FatalIfErrorf(errors.Join(err, errors.New("failed to initialize ci6ndex")))
	}
	b, err := initBot(c, config)
	if err != nil {
		ctx.FatalIfErrorf(errors.Join(err, errors.New("failed to initialize bot")))
	}

	err = ctx.Run(b)
	if err != nil {
		ctx.FatalIfErrorf(errors.Join(err, errors.New("failed to run bot")))
	}
}

func initBot(c *ci6ndex.Ci6ndex, config *bot.Config) (*bot.Bot, error) {
	b := bot.New(*config, c)
	err := bot.Configure(b, handler.New())
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to configure bot"))
	}
	return b, nil
}
