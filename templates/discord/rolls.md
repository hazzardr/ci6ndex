# Rolls

{{range .}}
### {{.User}}
{{range .Leaders}}
- {{if eq .DiscordEmojiString.Valid true}}{{.DiscordEmojiString.String}} {{end}}{{.LeaderName}} {{end}} {{end}}
