package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	
	"github.com/knetz-io/knetz/pkg/config"
	"github.com/knetz-io/knetz/pkg/storage"
)

var statusTenant string

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show system status",
	Long:  `Display overall status of clusters, services, and recent scans.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVar(&statusTenant, "tenant", "", "filter by tenant")
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	// Open storage
	store, err := storage.New(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("could not open storage: %w", err)
	}
	defer store.Close()

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Knetz Status                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Show configuration summary
	fmt.Printf("Configured Clusters: %d\n", len(cfg.Clusters))
	fmt.Printf("Configured Tenants: %d\n", len(cfg.Tenants))
	fmt.Println()

	// Show storage info
	fmt.Printf("Storage Type: %s\n", cfg.Storage.Type)
	fmt.Printf("Storage Path: %s\n", cfg.Storage.Path)
	fmt.Println()

	fmt.Println("✅ Knetz is operational")
	fmt.Println()
	fmt.Println("Run 'knetz scan --all' to discover services")

	return nil
}

