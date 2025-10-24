package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	
	"github.com/knetz-io/knetz/pkg/cluster"
	"github.com/knetz-io/knetz/pkg/config"
)

var (
	clusterName string
	allClusters bool
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage cluster connections",
	Long:  `Test connectivity, list clusters, and manage cluster information.`,
}

// clusterTestCmd tests cluster connectivity
var clusterTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test cluster connectivity",
	Long:  `Test connectivity to one or more configured clusters.`,
	RunE:  runClusterTest,
}

// clusterListCmd lists configured clusters
var clusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured clusters",
	Long:  `Display all configured clusters and their details.`,
	RunE:  runClusterList,
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterTestCmd)
	clusterCmd.AddCommand(clusterListCmd)

	// Flags for cluster test
	clusterTestCmd.Flags().StringVar(&clusterName, "cluster", "", "cluster name to test")
	clusterTestCmd.Flags().BoolVar(&allClusters, "all", false, "test all clusters")

	// Flags for cluster list
	clusterListCmd.Flags().StringVar(&cfgFile, "tenant", "", "filter by tenant")
}

func runClusterTest(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var clustersToTest []config.ClusterConfig

	if allClusters {
		clustersToTest = cfg.Clusters
	} else if clusterName != "" {
		clusterCfg := cfg.GetCluster(clusterName)
		if clusterCfg == nil {
			return fmt.Errorf("cluster %s not found in config", clusterName)
		}
		clustersToTest = []config.ClusterConfig{*clusterCfg}
	} else {
		return fmt.Errorf("specify --cluster or --all flag")
	}

	fmt.Println("Testing cluster connectivity...")
	fmt.Println()

	successCount := 0
	for _, clusterCfg := range clustersToTest {
		fmt.Printf("Testing %s...\n", clusterCfg.Name)

		client, err := cluster.NewClient(&clusterCfg)
		if err != nil {
			fmt.Printf("  ❌ Failed to create client: %v\n", err)
			continue
		}

		if err := client.Test(ctx); err != nil {
			fmt.Printf("  ❌ Connection failed: %v\n", err)
			continue
		}

		version, err := client.GetServerVersion(ctx)
		if err != nil {
			fmt.Printf("  ⚠️  Connected but could not get version: %v\n", err)
		} else {
			fmt.Printf("  ✅ Connected successfully\n")
			fmt.Printf("     Kubernetes Version: %s\n", version)
			fmt.Printf("     Context: %s\n", clusterCfg.Context)
			if clusterCfg.Tenant != "" {
				fmt.Printf("     Tenant: %s\n", clusterCfg.Tenant)
			}
		}
		fmt.Println()
		successCount++
	}

	fmt.Printf("Results: %d/%d clusters reachable\n", successCount, len(clustersToTest))

	if successCount < len(clustersToTest) {
		return fmt.Errorf("some clusters unreachable")
	}

	return nil
}

func runClusterList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	fmt.Println("Configured Clusters:")
	fmt.Println()

	for _, cluster := range cfg.Clusters {
		fmt.Printf("• %s\n", cluster.Name)
		fmt.Printf("  Context: %s\n", cluster.Context)
		if cluster.Tenant != "" {
			fmt.Printf("  Tenant: %s\n", cluster.Tenant)
		}
		if len(cluster.Namespaces) > 0 {
			fmt.Printf("  Namespaces: %v\n", cluster.Namespaces)
		}
		if cluster.Platform != "" {
			fmt.Printf("  Platform: %s\n", cluster.Platform)
		}
		fmt.Println()
	}

	return nil
}

