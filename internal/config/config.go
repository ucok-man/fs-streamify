package config

import (
	"context"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ucok-man/streamify/internal/validator"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

type Config struct {
	Port int    `mapstructure:"API_PORT"`
	Env  string `mapstructure:"API_ENV"`
	Log  struct {
		Level string `mapstructure:"API_LOG_LEVEL"`
	} `mapstructure:",squash"`
	DB struct {
		MongoURI      string        `mapstructure:"API_DB_MONGO_URI"`
		DatabaseName  string        `mapstructure:"API_DB_DATABASE_NAME"`
		MaxConnecting uint64        `mapstructure:"API_DB_MAX_CONNECTING"`
		MaxPoolSize   uint64        `mapstructure:"API_DB_MAX_POOL_SIZE"`
		MaxIdleTime   time.Duration `mapstructure:"API_DB_MAX_IDLE_TIME"`
	} `mapstructure:",squash"`
	Cors struct {
		Origins []string `mapstructure:"API_CORS_ORIGINS"`
	} `mapstructure:",squash"`
	GetStreamIO struct {
		ApiKey    string `mapstructure:"API_GETSTREAMIO_API_KEY"`
		ApiSecret string `mapstructure:"API_GETSTREAMIO_API_SECRET"`
	} `mapstructure:",squash"`
	JWT struct {
		AuthSecret string `mapstructure:"API_JWT_AUTH_SECRET"`
	} `mapstructure:",squash"`
}

func New() Config {
	viper.SetConfigFile(".env") // Config file name without extension
	viper.SetConfigType("env")  // Config file type
	viper.AddConfigPath(".")    // Look for the config file in the current directory
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Error reading config file")
	}

	var config Config
	err := viper.Unmarshal(&config, func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnset = true
		// dc.ErrorUnused = true
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to decode config file")
	}

	errmap := validator.Schema().Config.Validate(&config)
	if errmap != nil {
		log.Fatal().
			Err(errors.New("invalid or missing config")).
			Any("err_detail", validator.Sanitize(errmap)).
			Msg("Invalid or missing config")
	}

	return config
}

func (cfg Config) OpenDB() (*mongo.Client, error) {
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
