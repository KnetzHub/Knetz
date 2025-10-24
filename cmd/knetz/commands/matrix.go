package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	
	"github.com/knetz-io/knetz/pkg/config"
	"github.com/knetz-io/knetz/pkg/storage"
)

var (
	matrixTenant string
	matrixClusters []string
)

// matrixCmd represents the matrix command
var matrixCmd = &cobra.Command{
	Use:   "matrix",
	Short: "Display version matrix across clusters",
	Long: `Show a multi-dimensional matrix view of service versions across
clusters and namespaces.`,
	RunE: runMatrix,
}

func init() {
	rootCmd.AddCommand(matrixCmd)
	matrixCmd.Flags().StringVar(&matrixTenant, "tenant", "", "tenant to display")
	matrixCmd.Flags().StringSliceVar(&matrixClusters, "cluster", []string{}, "clusters to include")
}

func runMatrix(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	store, err := storage.New(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("could not open storage: %w", err)
	}
	defer store.Close()

	// Determine clusters to display
	clustersToShow := matrixClusters
	if matrixTenant != "" {
		tenantClusters := cfg.GetClustersByTenant(matrixTenant)
		clustersToShow = make([]string, len(tenantClusters))
		for i, c := range tenantClusters {
			clustersToShow[i] = c.Name
		}
	}

	if len(clustersToShow) == 0 {
		clustersToShow = make([]string, len(cfg.Clusters))
		for i, c := range cfg.Clusters {
			clustersToShow[i] = c.Name
		}
	}

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	if matrixTenant != "" {
		fmt.Printf("║   Service Version Matrix - Tenant: %-25s ║\n", matrixTenant)
	} else {
		fmt.Printf("║   Service Version Matrix - All Clusters                      ║\n")
	}
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Collect all services
	serviceMap := make(map[string]map[string]string) // service -> cluster -> version

	for _, clusterName := range clustersToShow {
		services, err := store.GetServicesByCluster(clusterName)
		if err != nil {
			fmt.Printf("⚠️  Could not load services from %s: %v\n", clusterName, err)
			continue
		}

		for _, svc := range services {
			if serviceMap[svc.Name] == nil {
				serviceMap[svc.Name] = make(map[string]string)
			}
			serviceMap[svc.Name][clusterName] = svc.Version
		}
	}

	// Print header
	fmt.Printf("%-30s", "Service")
	for _, cluster := range clustersToShow {
		fmt.Printf(" │ %-15s", cluster)
	}
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")

	// Print services
	for serviceName, clusterVersions := range serviceMap {
		fmt.Printf("%-30s", serviceName)
		for _, cluster := range clustersToShow {
			version := clusterVersions[cluster]
			if version == "" {
				version = "-"
			}
			fmt.Printf(" │ %-15s", version)
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Printf("Total unique services: %d\n", len(serviceMap))
	fmt.Printf("Clusters compared: %d\n", len(clustersToShow))

	return nil
}

