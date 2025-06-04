package main

import (
	"context"
	"sync"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ucok-man/streamify/internal/config"
	"github.com/ucok-man/streamify/internal/logger"
	"github.com/ucok-man/streamify/internal/models"
)

type application struct {
	config config.Config
	logger *zerolog.Logger
	models models.Models
	stream *stream.Client
	wg     sync.WaitGroup
}

func main() {
	cfg := config.New()
	applog, err := logger.New(cfg.Log.Level, cfg.Env)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed initialize logger")
	}

	dbclient, err := cfg.OpenDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed initialize db connection")
	}
	defer dbclient.Disconnect(context.Background())

	// Initialize stream chat client
	streamChatClient, err := stream.NewClient(cfg.GetStreamIO.ApiKey, cfg.GetStreamIO.ApiSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed initialize stream chat client")
	}

	app := &application{
		config: cfg,
		logger: applog,
		stream: streamChatClient,
		models: models.NewModels(dbclient.Database(cfg.DB.DatabaseName), applog),
	}

	if err := app.serve(); err != nil {
		log.Fatal().Err(err).Msg("Failed running server")
	}
}
