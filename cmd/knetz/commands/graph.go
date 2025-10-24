package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	
	"github.com/knetz-io/knetz/pkg/config"
	"github.com/knetz-io/knetz/pkg/storage"
	"github.com/knetz-io/knetz/pkg/graph"
)

var (
	graphTenant  string
	graphCluster string
	graphService string
	graphOutput  string
	graphDepth   int
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Generate dependency graphs",
	Long: `Generate and visualize dependency graphs for services across clusters.
	
Supports multiple output formats including SVG, PNG, and DOT.`,
	RunE: runGraph,
}

func init() {
	rootCmd.AddCommand(graphCmd)

	graphCmd.Flags().StringVar(&graphTenant, "tenant", "", "tenant to graph")
	graphCmd.Flags().StringVar(&graphCluster, "cluster", "", "cluster to graph")
	graphCmd.Flags().StringVar(&graphService, "service", "", "specific service to graph")
	graphCmd.Flags().StringVar(&graphOutput, "output", "", "output file (e.g., graph.svg)")
	graphCmd.Flags().IntVar(&graphDepth, "depth", 0, "max depth for dependency traversal (0 = unlimited)")
}

func runGraph(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	store, err := storage.New(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("could not open storage: %w", err)
	}
	defer store.Close()

	fmt.Println("🔗 Generating dependency graph...")
	fmt.Println()

	// Build the graph
	builder := graph.NewBuilder()

	var clustersToGraph []string
	if graphTenant != "" {
		tenantClusters := cfg.GetClustersByTenant(graphTenant)
		for _, c := range tenantClusters {
			clustersToGraph = append(clustersToGraph, c.Name)
		}
	} else if graphCluster != "" {
		clustersToGraph = []string{graphCluster}
	} else {
		for _, c := range cfg.Clusters {
			clustersToGraph = append(clustersToGraph, c.Name)
		}
	}

	totalServices := 0
	for _, clusterName := range clustersToGraph {
		services, err := store.GetServicesByCluster(clusterName)
		if err != nil {
			fmt.Printf("⚠️  Could not load services from %s: %v\n", clusterName, err)
			continue
		}

		for _, svc := range services {
			if graphService == "" || svc.Name == graphService {
				builder.AddService(svc)
				totalServices++
			}
		}
	}

	g := builder.Build()
	
	fmt.Printf("Nodes: %d\n", len(g.GetNodes()))
	fmt.Printf("Edges: %d\n", len(g.GetEdges()))
	fmt.Println()

	// Detect cycles
	cycles := g.DetectCycles()
	if len(cycles) > 0 {
		fmt.Printf("⚠️  Detected %d circular dependencies:\n", len(cycles))
		for i, cycle := range cycles {
			fmt.Printf("  %d. %v\n", i+1, cycle)
		}
		fmt.Println()
	} else {
		fmt.Println("✅ No circular dependencies detected")
		fmt.Println()
	}

	// Calculate statistics
	if graphService != "" {
		for nodeID, node := range g.GetNodes() {
			if node.ServiceName == graphService {
				depth := g.CalculateDepth(nodeID)
				impact := g.CalculateImpactScore(nodeID)
				
				fmt.Printf("Service: %s\n", graphService)
				fmt.Printf("  Dependency depth: %d\n", depth)
				fmt.Printf("  Impact score: %d services depend on this\n", impact)
				fmt.Printf("  Direct dependencies: %d\n", len(g.GetDependencies(nodeID)))
				fmt.Printf("  Direct dependents: %d\n", len(g.GetDependents(nodeID)))
				break
			}
		}
	}

	if graphOutput != "" {
		fmt.Printf("\nℹ️  Graph visualization to %s will be available in Phase 2\n", graphOutput)
		fmt.Println("   (DOT/SVG/PNG export coming soon)")
	}

	return nil
}

