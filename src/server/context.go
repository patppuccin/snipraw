package server

import (
	"github.com/patppuccin/snipraw/src/config"
	"github.com/rs/zerolog"
)

type ctxKey int

const appCtxKey ctxKey = iota

type appCtx struct {
	dir    string
	config *config.Config
	logger *zerolog.Logger
}
