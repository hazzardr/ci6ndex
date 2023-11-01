package main

import (
	"ci6ndex/internal"
	"os"
)

func main() {
	if os.Args[1] == "mode" {
		mode := os.Args[2]
		internal.Start(mode)
		return
	}
	internal.Start("bot")
	return
}
