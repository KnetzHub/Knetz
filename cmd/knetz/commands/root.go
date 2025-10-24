package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Version information
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"

	// Global flags
	cfgFile   string
	verbose   bool
	debug     bool
	noColor   bool
	outputFmt string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "knetz",
	Short: "Cross-Cluster Dependency & Version Manager",
	Long: `Knetz is a CLI tool for managing dependencies and versions across multiple 
Kubernetes clusters, providing global visibility and intelligent drift detection.

Knetz helps you:
  • Discover services across multiple clusters and namespaces
  • Track version drift and incompatibilities
  • Validate dependencies before deployment
  • Visualize dependency graphs
  • Plan safe rollouts`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.knetz/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode with detailed logs")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "table", "output format (table|json|yaml|csv)")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".knetz" (without extension)
		configDir := filepath.Join(home, ".knetz")
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Read in environment variables that match
	viper.SetEnvPrefix("KNETZ")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	}
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return viper.GetBool("verbose")
}

// IsDebug returns whether debug mode is enabled
func IsDebug() bool {
	return viper.GetBool("debug")
}

// NoColor returns whether color output is disabled
func NoColor() bool {
	return viper.GetBool("no-color")
}

// GetOutputFormat returns the configured output format
func GetOutputFormat() string {
	return viper.GetString("output")
}

