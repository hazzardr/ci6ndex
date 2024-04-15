package discord

import (
	"bytes"
	"ci6ndex/internal"
	"path/filepath"
	"text/template"
)

var (
	templates = map[string]string{
		RollCivs.Name: "./templates/discord/rolls.md",
	}
)

type MessageBuilder struct {
	tmpl *template.Template
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

// WriteDraftInfo writes info in a discord friendly format given the draft offerings.
func (mb *MessageBuilder) WriteDraftInfo(cmdName string, offers []internal.DraftOffering) (string, error) {
	tmplName := filepath.Base(templates[cmdName])
	var buf bytes.Buffer
	err := mb.tmpl.ExecuteTemplate(&buf, tmplName, offers)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
