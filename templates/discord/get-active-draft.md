The following players are in an active draft, with game starting {{.Date}}:
{{range .Players}}
- {{.}}{{end}}

{{ if eq (len .NoPickYet) 0 }}
All players have picked a Civ.
{{ else }} Players who have not picked a Civ yet:
{{range .NoPickYet}}
- {{.}}{{end}}
{{ end }}
