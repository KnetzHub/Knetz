package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	
	"github.com/knetz-io/knetz/pkg/config"
	"github.com/knetz-io/knetz/pkg/discovery"
	"github.com/knetz-io/knetz/pkg/storage"
)

var (
	scanCluster   string
	scanNamespace string
	scanTenant    string
	scanAll       bool
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan clusters for services",
	Long: `Discover and inventory services across Kubernetes clusters.
	
Scans configured clusters to discover services and their versions.`,
	RunE: runScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVar(&scanCluster, "cluster", "", "cluster name to scan")
	scanCmd.Flags().StringVar(&scanNamespace, "namespace", "", "namespace to scan (requires --cluster)")
	scanCmd.Flags().StringVar(&scanTenant, "tenant", "", "tenant to scan")
	scanCmd.Flags().BoolVar(&scanAll, "all", false, "scan all configured clusters")
}

func runScan(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Open storage
	store, err := storage.New(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("could not open storage: %w", err)
	}
	defer store.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Determine which clusters to scan
	var clustersToScan []config.ClusterConfig

	if scanAll {
		clustersToScan = cfg.Clusters
	} else if scanTenant != "" {
		clustersToScan = cfg.GetClustersByTenant(scanTenant)
		if len(clustersToScan) == 0 {
			return fmt.Errorf("no clusters found for tenant %s", scanTenant)
		}
	} else if scanCluster != "" {
		clusterCfg := cfg.GetCluster(scanCluster)
		if clusterCfg == nil {
			return fmt.Errorf("cluster %s not found", scanCluster)
		}
		clustersToScan = []config.ClusterConfig{*clusterCfg}
	} else {
		return fmt.Errorf("specify --cluster, --tenant, or --all")
	}

	fmt.Printf("🔍 Scanning %d cluster(s)...\n\n", len(clustersToScan))

	totalServices := 0
	for _, clusterCfg := range clustersToScan {
		fmt.Printf("Scanning cluster: %s\n", clusterCfg.Name)

		scanner := discovery.NewScanner(&clusterCfg, cfg, store)
		services, err := scanner.Scan(ctx)
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("  ✅ Found %d services\n", len(services))
		totalServices += len(services)

		// Save services to storage
		for _, service := range services {
			if err := store.SaveService(service); err != nil {
				fmt.Printf("  ⚠️  Could not save service %s: %v\n", service.Name, err)
			}
		}
		fmt.Println()
	}

	fmt.Printf("✅ Scan complete! Total services discovered: %d\n", totalServices)

	return nil
}

