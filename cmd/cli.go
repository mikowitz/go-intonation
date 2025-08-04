package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-intonation",
	Short: "A CLI tool for working with just intonation",
	Long:  `A command line tool for working with just intonation, ratios, diamonds, and lattices.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(ratioCmd)
	rootCmd.AddCommand(diamondCmd)
	rootCmd.AddCommand(latticeCmd)
}

func Run() {
	Execute()
}
