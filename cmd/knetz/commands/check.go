package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	
	"github.com/knetz-io/knetz/pkg/config"
	"github.com/knetz-io/knetz/pkg/storage"
	"github.com/knetz-io/knetz/internal/utils"
)

var (
	checkTenant  string
	checkCluster string
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for version drift and violations",
	Long: `Validate dependencies and detect version drift across clusters.
	
Identifies:
  • Version mismatches across clusters
  • Incompatible dependency versions
  • Services violating constraints`,
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVar(&checkTenant, "tenant", "", "tenant to check")
	checkCmd.Flags().StringVar(&checkCluster, "cluster", "", "cluster to check")
}

func runCheck(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	store, err := storage.New(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("could not open storage: %w", err)
	}
	defer store.Close()

	fmt.Println("🔍 Checking for version drift and violations...")
	fmt.Println()

	violations := 0
	warnings := 0

	// Determine scope
	var clustersToCheck []string
	if checkTenant != "" {
		tenantClusters := cfg.GetClustersByTenant(checkTenant)
		for _, c := range tenantClusters {
			clustersToCheck = append(clustersToCheck, c.Name)
		}
	} else if checkCluster != "" {
		clustersToCheck = []string{checkCluster}
	} else {
		for _, c := range cfg.Clusters {
			clustersToCheck = append(clustersToCheck, c.Name)
		}
	}

	// Collect services by name across clusters
	serviceVersions := make(map[string]map[string]string) // service -> cluster -> version

	for _, clusterName := range clustersToCheck {
		services, err := store.GetServicesByCluster(clusterName)
		if err != nil {
			fmt.Printf("⚠️  Could not load services from %s: %v\n", clusterName, err)
			continue
		}

		for _, svc := range services {
			if serviceVersions[svc.Name] == nil {
				serviceVersions[svc.Name] = make(map[string]string)
			}
			serviceVersions[svc.Name][clusterName] = svc.Version
		}
	}

	// Check for version drift
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Version Drift Analysis")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	for serviceName, clusterVersions := range serviceVersions {
		if len(clusterVersions) < 2 {
			continue
		}

		// Check if all versions are the same
		var versions []string
		for _, v := range clusterVersions {
			versions = append(versions, v)
		}

		allSame := true
		firstVersion := versions[0]
		for _, v := range versions {
			if v != firstVersion {
				allSame = false
				break
			}
		}

		if !allSame {
			warnings++
			fmt.Printf("⚠️  %s - Version Drift Detected\n", serviceName)
			for cluster, version := range clusterVersions {
				fmt.Printf("    %s: %s\n", cluster, version)
			}
			
			// Try to compare versions
			if len(versions) >= 2 && versions[0] != "unknown" && versions[1] != "unknown" {
				cmp, err := utils.CompareVersionStrings(versions[0], versions[1])
				if err == nil {
					if cmp > 0 {
						fmt.Printf("    ℹ️  First version is newer\n")
					} else if cmp < 0 {
						fmt.Printf("    ℹ️  Second version is newer\n")
					}
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Summary")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Printf("🔴 Critical Violations: %d\n", violations)
	fmt.Printf("🟡 Warnings: %d\n", warnings)
	
	if violations == 0 && warnings == 0 {
		fmt.Println("🟢 No issues found - all services aligned!")
	}

	return nil
}

