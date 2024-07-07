package discord

import (
	"bytes"
	"ci6ndex-bot/domain"
	"ci6ndex-bot/pkg"
	"path/filepath"
	"text/template"
)

var (
	templates = map[string]string{
		RollCivs.Name:              "./templates/discord/rolls.md",
		CreateDraftConfirmId:       "./templates/discord/start-draft.md",
		CheckActiveDraftCommand.ID: "./templates/discord/get-active-draft.md",
		RevealCivs.Name:            "./templates/discord/reveal-civs.md",
	}
)

type MessageBuilder struct {
	tmpl *template.Template
}

type DraftStartProps struct {
	Players []string
	Date    string
}

type GetActiveDraftProps struct {
	Players   []string
	NoPickYet []string
	Date      string
}

// NewDiscTemplate Parse all disc message templates to be usable
func NewDiscTemplate() *MessageBuilder {
	tmpl := template.New("discord-messages")
	tmplFiles := make([]string, 0, len(templates))
	for _, tmplFile := range templates {
		tmplFiles = append(tmplFiles, tmplFile)
	}
	template.Must(tmpl.ParseFiles(tmplFiles...))
	return &MessageBuilder{
		tmpl: tmpl,
	}
}

// WriteDraftOfferings writes info in a discord friendly format given the draft offerings.
func (mb *MessageBuilder) WriteDraftOfferings(cmdName string, offers []pkg.DraftOffering) (string, error) {
	tmplName := filepath.Base(templates[cmdName])
	var buf bytes.Buffer
	err := mb.tmpl.ExecuteTemplate(&buf, tmplName, offers)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (mb *MessageBuilder) WriteConfirmDraft(cmdName string, players []string, date string) (string, error) {
	tmplName := filepath.Base(templates[cmdName])
	var buf bytes.Buffer
	err := mb.tmpl.ExecuteTemplate(&buf, tmplName, &DraftStartProps{
		Players: players,
		Date:    date,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (mb *MessageBuilder) WriteActiveDraft(cmdName string, players []string, noPickYet []string, date string) (string, error) {
	tmplName := filepath.Base(templates[cmdName])
	var buf bytes.Buffer
	draft := &GetActiveDraftProps{
		Players:   players,
		NoPickYet: noPickYet,
		Date:      date,
	}
	err := mb.tmpl.ExecuteTemplate(&buf, tmplName, draft)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type RevealCivsProps struct {
	Picks []domain.GetDenormalizedDraftPicksForDraftRow
	Start string
}

func (mb *MessageBuilder) WriteFinalizedPicks(cmdName string,
	picks []domain.GetDenormalizedDraftPicksForDraftRow) (string, error) {
	tmplName := filepath.Base(templates[cmdName])

	date := picks[0].StartDate.Time.Format("2006-01-02")
	props := &RevealCivsProps{
		Picks: picks,
		Start: date,
	}

	var buf bytes.Buffer
	err := mb.tmpl.ExecuteTemplate(&buf, tmplName, props)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
