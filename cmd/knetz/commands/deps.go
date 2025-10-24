package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	
	"github.com/knetz-io/knetz/pkg/config"
	"github.com/knetz-io/knetz/pkg/dependency"
)

var (
	depsFile    string
	depsService string
	depsCluster string
)

// depsCmd represents the deps command group
var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Manage service dependencies",
	Long:  `Discover, validate, and manage service dependencies across clusters.`,
}

// depsShowCmd shows dependencies for a service
var depsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show dependencies for a service",
	Long:  `Display the dependency tree for a specific service.`,
	RunE:  runDepsShow,
}

// depsValidateCmd validates dependency specifications
var depsValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate dependency specifications",
	Long:  `Validate a dependency specification file for correctness.`,
	RunE:  runDepsValidate,
}

// depsExportCmd exports dependencies
var depsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export service dependencies",
	Long:  `Export discovered dependencies to a file.`,
	RunE:  runDepsExport,
}

func init() {
	rootCmd.AddCommand(depsCmd)
	
	depsCmd.AddCommand(depsShowCmd)
	depsCmd.AddCommand(depsValidateCmd)
	depsCmd.AddCommand(depsExportCmd)

	// Flags for deps show
	depsShowCmd.Flags().StringVar(&depsService, "service", "", "service name")
	depsShowCmd.Flags().StringVar(&depsCluster, "cluster", "", "cluster name")
	depsShowCmd.MarkFlagRequired("service")
	depsShowCmd.MarkFlagRequired("cluster")

	// Flags for deps validate
	depsValidateCmd.Flags().StringVarP(&depsFile, "file", "f", "", "dependency specification file")
	depsValidateCmd.MarkFlagRequired("file")

	// Flags for deps export
	depsExportCmd.Flags().StringVar(&depsCluster, "cluster", "", "cluster to export from")
	depsExportCmd.Flags().StringVarP(&outputFmt, "output", "o", "yaml", "output format (yaml|json)")
}

func runDepsShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	fmt.Printf("🔍 Showing dependencies for %s in cluster %s\n\n", depsService, depsCluster)
	
	fmt.Printf("Service: %s\n", depsService)
	fmt.Printf("Cluster: %s\n", depsCluster)
	fmt.Println()
	fmt.Println("Dependencies:")
	fmt.Println("  (Run 'knetz scan' to discover dependencies)")
	fmt.Println()
	fmt.Println("Depended By:")
	fmt.Println("  (Coming soon in Phase 2)")

	_ = cfg // Use cfg to avoid unused variable error
	
	return nil
}

func runDepsValidate(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(depsFile)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	parser := dependency.NewParser()
	spec, err := parser.ParseYAML(data)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println("✅ Dependency specification is valid")
	fmt.Println()
	fmt.Printf("Service: %s\n", spec.Metadata.Name)
	fmt.Printf("Version: %s\n", spec.Metadata.Version)
	fmt.Printf("Dependencies: %d\n", len(spec.Spec.Dependencies))
	
	if len(spec.Spec.CrossTenantDependencies) > 0 {
		fmt.Printf("Cross-tenant dependencies: %d\n", len(spec.Spec.CrossTenantDependencies))
	}

	return nil
}

func runDepsExport(cmd *cobra.Command, args []string) error {
	fmt.Println("📦 Exporting dependencies...")
	fmt.Println()
	fmt.Println("ℹ️  Export functionality will be available after scanning clusters")
	fmt.Println("   Run 'knetz scan --cluster <name>' first")
	
	return nil
}

