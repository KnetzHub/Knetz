package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	
	"github.com/knetz-io/knetz/pkg/config"
	"github.com/knetz-io/knetz/pkg/storage"
)

var (
	diffClusters []string
	diffService  string
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare versions across clusters",
	Long: `Compare service versions across different clusters, namespaces, or tenants.
	
Shows version differences and highlights mismatches.`,
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringSliceVar(&diffClusters, "cluster", []string{}, "clusters to compare (specify multiple times)")
	diffCmd.Flags().StringVar(&diffService, "service", "", "specific service to compare")
}

func runDiff(cmd *cobra.Command, args []string) error {
	if len(diffClusters) < 2 {
		return fmt.Errorf("specify at least 2 clusters to compare using --cluster flag")
	}

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	store, err := storage.New(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("could not open storage: %w", err)
	}
	defer store.Close()

	fmt.Printf("🔍 Comparing versions across clusters: %v\n\n", diffClusters)

	// Get services from each cluster
	for _, clusterName := range diffClusters {
		services, err := store.GetServicesByCluster(clusterName)
		if err != nil {
			fmt.Printf("❌ Error loading services from %s: %v\n", clusterName, err)
			continue
		}

		fmt.Printf("Cluster: %s\n", clusterName)
		fmt.Printf("  Services found: %d\n", len(services))
		
		if diffService != "" {
			// Filter for specific service
			found := false
			for _, svc := range services {
				if svc.Name == diffService {
					fmt.Printf("  %s: %s\n", svc.Name, svc.Version)
					found = true
				}
			}
			if !found {
				fmt.Printf("  Service '%s' not found\n", diffService)
			}
		}
		fmt.Println()
	}

	return nil
}

