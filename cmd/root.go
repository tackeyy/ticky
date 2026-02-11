package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version     = "dev"
	commit      = "none"
	date        = "unknown"
	outputJSON  bool
	outputPlain bool
)

var rootCmd = &cobra.Command{
	Use:   "ticky",
	Short: "TickTick CLI tool",
	Long:  "ticky — A CLI tool for TickTick task management. Designed for both human use and AI agent integration.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&outputPlain, "plain", false, "Output in TSV format")
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("ticky version %s (commit: %s, built: %s)\n", version, commit, date))
}
