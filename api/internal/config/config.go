package config

import (
	"errors"

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
	}
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
		log.Fatal().Err(errors.New("invalid or missing config")).Any("err_detail", validator.Sanitize(errmap)).Msg("Invalid or missing config")
	}

	return config
}
