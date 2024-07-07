# The draft is now closed 
# The game will start on {{.Start}}

{{range .Picks}}
- **{{.DiscordName}} picked: {{if eq .DiscordEmojiString.Valid true}}{{.DiscordEmojiString.String}} {{end}}{{.LeaderName}}** {{end}}
