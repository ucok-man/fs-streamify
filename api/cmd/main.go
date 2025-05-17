package main

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ucok-man/streamify-api/internal/config"
	"github.com/ucok-man/streamify-api/internal/logger"
)

type application struct {
	config config.Config
	logger *zerolog.Logger
	wg     sync.WaitGroup
}

func main() {
	cfg := config.New()
	applog, err := logger.New(cfg.Log.Level, cfg.Env)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed initialize logger")
	}
	app := &application{
		config: cfg,
		logger: applog,
	}

	if err := app.serve(); err != nil {
		log.Fatal().Err(err).Msg("Failed running server")
	}
}
