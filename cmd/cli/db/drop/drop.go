package drop

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/ucok-man/streamify/internal/config"
	"github.com/ucok-man/streamify/internal/logger"
)

var DropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop in database",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		logger, err := logger.New(cfg.Log.Level, cfg.Env)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed initialize logger")
		}

		conn, err := cfg.OpenDB()
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed initialize db connection")
		}
		defer conn.Disconnect(context.Background())

		err = conn.Database(cfg.DB.DatabaseName).Drop(context.Background())
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to drop database")
		}

		logger.Info().Msgf("Success dropping %v database", cfg.DB.DatabaseName)
	},
}
