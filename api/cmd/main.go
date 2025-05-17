package main

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ucok-man/streamify-api/internal/config"
	"github.com/ucok-man/streamify-api/internal/logger"
	"github.com/ucok-man/streamify-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type application struct {
	config config.Config
	logger *zerolog.Logger
	models models.Models
	wg     sync.WaitGroup
}

func main() {
	cfg := config.New()
	applog, err := logger.New(cfg.Log.Level, cfg.Env)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed initialize logger")
	}

	dbclient, err := openDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed initialize db connection")
	}
	defer dbclient.Disconnect(context.Background())

	app := &application{
		config: cfg,
		logger: applog,
		models: models.NewModels(dbclient.Database(cfg.DB.DatabaseName)),
	}

	if err := app.serve(); err != nil {
		log.Fatal().Err(err).Msg("Failed running server")
	}
}

func openDB(cfg config.Config) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(cfg.DB.MongoURI).
		SetServerAPIOptions(serverAPI)

	// Nilai ini tidak membatasi total koneksi dalam pool, hanya membatasi
	// berapa banyak koneksi baru yang bisa dibuat secara bersamaan.
	opts.SetMaxConnecting(cfg.DB.MaxConnecting)

	// Jumlah maksimum koneksi dalam connection pool
	// Jika mencapai nilai maksimum, permintaan baru ke server akan diblokir
	// (menunggu) sampai ada koneksi yang tersedia.
	opts.SetMaxPoolSize(cfg.DB.MaxPoolSize)

	// batas waktu maksimum sebuah koneksi boleh menganggur (idle) di dalam
	// connection pool sebelum koneksi tersebut ditutup dan dibuang.
	opts.SetMaxConnIdleTime(cfg.DB.MaxIdleTime)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, err
}
