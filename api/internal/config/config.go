package config

import (
	"errors"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ucok-man/streamify-api/internal/validator"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

type Config struct {
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
	Log  struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
	DB struct {
		MongoURI      string        `mapstructure:"mongo_uri"`
		DatabaseName  string        `mapstructure:"database_name"`
		MaxConnecting uint64        `mapstructure:"max_connecting"`
		MaxPoolSize   uint64        `mapstructure:"max_pool_size"`
		MaxIdleTime   time.Duration `mapstructure:"max_idle_time"`
	} `mapstructure:"db"`

	Cors struct {
		Origins []string `mapstructure:"origins"` // viper/mapstructure automatically split by comma and convert it
	} `mapstructure:"cors"`
	StreamIO struct {
		ApiKey    string `mapstructure:"api_key"`
		ApiSecret string `mapstructure:"api_secret"`
	} `mapstructure:"streamio"`
}

func New() Config {
	viper.SetConfigName("config") // Config file name without extension
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(".")      // Look for the config file in the current directory

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("Error reading config file")
	}

	var config Config
	err := viper.Unmarshal(&config, func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnset = true
		dc.ErrorUnused = true
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
