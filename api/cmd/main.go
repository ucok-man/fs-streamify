package main

import (
	"github.com/rs/zerolog/log"
	"github.com/ucok-man/streamify-api/internal/config"
)

type application struct {
	config config.Config
}

func main() {
	cfg := config.New()
	// app := &application{
	// 	config: cfg,
	// }
	// if err := app.serve(); err != nil {
	// 	slog.Error("Error running server", err)
	// 	os.Exit(1)
	// }
	log.Info().CallerSkipFrame(2).Any("config", cfg).Msg("Hello From Main!")
}
