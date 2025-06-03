package bot

import (
	"errors"
	md "github.com/nao1215/markdown"
	"io"
)

const (
	colorSuccess    = 0x5c5fea // Using the Unicode escape sequence
	partyEmoji      = "\U0001F389"
	magnifyingGlass = "\U0001F50D"
	crossedSwords   = "\u2694\uFE0F"
	backArrow       = "\u2B05\uFE0F"
)

func renderMainScreen(header, previousGame io.Writer) error {
	err := renderDraftHeader(header)
	if err != nil {
		return errors.Join(err, errors.New("failed to render draft header"))
	}
	err = renderPreviousGameSummary(previousGame)
	if err != nil {
		return errors.Join(err, errors.New("failed to render recent games"))
	}
	return nil
}

func renderDraftHeader(header io.Writer) error {
	return md.NewMarkdown(header).H1("Ci6ndex Draft Manager").
		PlainText("Civ (VI) Index helps manage drafts and stores match history.").
		Build()
}

func renderPreviousGameSummary(output io.Writer) error {
	return md.NewMarkdown(output).H2("Previous Game").
		H3f("**%s Winner:** <:Nzinga_Mbande_Civ6:1229393600790663220> <@135218870494429184>",
			partyEmoji).
		Build()
}
