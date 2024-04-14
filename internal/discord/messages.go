package discord

import (
	"bytes"
	"ci6ndex/internal"
	"text/template"
)

type MessageBuilder struct {
	tmpl *template.Template
}

func NewDiscTemplate() *MessageBuilder {
	tmpl := template.New("discord-messages")
	template.Must(tmpl.ParseFiles("./templates/discord/rolls.md"))
	return &MessageBuilder{
		tmpl: tmpl,
	}
}

func (mb *MessageBuilder) WriteRollsMessage(offers []internal.DraftOffering) (string, error) {
	var buf bytes.Buffer
	err := mb.tmpl.ExecuteTemplate(&buf, "rolls.md", offers)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
