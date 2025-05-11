package ci6ndex

import "github.com/charmbracelet/log"

type Ci6ndex struct {
	Logger      *log.Logger
	Connections map[uint64]*DB
	Path        string
}
