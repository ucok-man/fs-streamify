package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ucok-man/streamify/cmd/cli/db"
)

func init() {
	rootCmd.AddCommand(db.DBCmd)
}

var rootCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "streamify-cli",
	Short:   "streamify-cli - Tools for manage streamify api",
	Example: "- streamify-cli db seed\n- streamify-cli db restart",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
