package db

import (
	"github.com/spf13/cobra"
	"github.com/ucok-man/streamify/cmd/cli/db/drop"
	"github.com/ucok-man/streamify/cmd/cli/db/seed"
)

func init() {
	DBCmd.AddCommand(seed.SeedCmd, drop.DropCmd)
}

var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage database migration",
}
