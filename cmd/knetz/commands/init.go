package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	initConfigPath string
	initForce      bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Knetz configuration",
	Long: `Generate a starter configuration file for Knetz.

This command creates a sample configuration file with example cluster definitions,
tenant groupings, and discovery settings. You can then customize this file to
match your environment.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&initConfigPath, "config-path", "", "path for config file (default is $HOME/.knetz/config.yaml)")
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "overwrite existing config file")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine config path
	configPath := initConfigPath
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not determine home directory: %w", err)
		}
		configPath = filepath.Join(home, ".knetz", "config.yaml")
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil && !initForce {
		return fmt.Errorf("config file already exists at %s (use --force to overwrite)", configPath)
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	// Generate sample config
	sampleConfig := generateSampleConfig()

	// Marshal to YAML
	yamlData, err := yaml.Marshal(sampleConfig)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	fmt.Printf("✅ Configuration file created at: %s\n\n", configPath)
	fmt.Println("Next steps:")
	fmt.Println("1. Edit the config file to add your cluster details")
	fmt.Println("2. Test cluster connectivity: knetz cluster test")
	fmt.Println("3. Scan your clusters: knetz scan --all")

	return nil
}

func generateSampleConfig() map[string]interface{} {
	return map[string]interface{}{
		"clusters": []map[string]interface{}{
			{
				"name":       "example-cluster",
				"kubeconfig": "~/.kube/config",
				"context":    "example-context",
				"namespaces": []string{"default", "production"},
				"tenant":     "example-tenant",
			},
		},
		"tenants": []map[string]interface{}{
			{
				"name":        "example-tenant",
				"description": "Example tenant for demonstration",
				"clusters":    []string{"example-cluster"},
			},
		},
		"scan_scope": map[string]bool{
			"cluster_level":   true,
			"namespace_level": true,
			"tenant_level":    true,
		},
		"filters": map[string]interface{}{
			"include_namespaces": []string{"production", "staging"},
			"exclude_namespaces": []string{"kube-system", "kube-public", "kube-node-lease"},
			"include_labels": map[string]string{
				"managed-by": "knetz",
			},
			"exclude_labels": map[string]string{
				"ignore": "true",
			},
		},
		"discovery": map[string]interface{}{
			"scan_deployments":      true,
			"scan_statefulsets":     true,
			"scan_daemonsets":       true,
			"scan_helm_releases":    true,
			"auto_discover_dependencies": true,
			"dependency_sources": []string{
				"helm_charts",
				"env_variables",
				"configmaps",
				"service_mesh",
			},
		},
		"output": map[string]interface{}{
			"default_format":        "table",
			"color_enabled":         true,
			"show_scope_indicators": true,
			"timezone":              "UTC",
		},
		"storage": map[string]interface{}{
			"type":           "sqlite",
			"path":           "~/.knetz/data.db",
			"retention_days": 90,
		},
	}
}

