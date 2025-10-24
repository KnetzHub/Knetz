package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display version, commit hash, and build date for Knetz.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Knetz Version: %s\n", Version)
		fmt.Printf("Git Commit: %s\n", Commit)
		fmt.Printf("Build Date: %s\n", Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

