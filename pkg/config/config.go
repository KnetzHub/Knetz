package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Clusters  []ClusterConfig  `mapstructure:"clusters"`
	Tenants   []TenantConfig   `mapstructure:"tenants"`
	ScanScope ScanScopeConfig  `mapstructure:"scan_scope"`
	Filters   FilterConfig     `mapstructure:"filters"`
	Discovery DiscoveryConfig  `mapstructure:"discovery"`
	Output    OutputConfig     `mapstructure:"output"`
	Storage   StorageConfig    `mapstructure:"storage"`
}

// ClusterConfig represents a single cluster configuration
type ClusterConfig struct {
	Name       string            `mapstructure:"name"`
	Kubeconfig string            `mapstructure:"kubeconfig"`
	Context    string            `mapstructure:"context"`
	Namespaces []string          `mapstructure:"namespaces"`
	Tenant     string            `mapstructure:"tenant"`
	Platform   string            `mapstructure:"platform"`
	APIURL     string            `mapstructure:"api_url"`
	OpenShift  *OpenShiftConfig  `mapstructure:"openshift"`
}

// OpenShiftConfig contains OpenShift-specific settings
type OpenShiftConfig struct {
	ScanDeploymentConfigs bool `mapstructure:"scan_deploymentconfigs"`
	ScanRoutes            bool `mapstructure:"scan_routes"`
	ScanBuildConfigs      bool `mapstructure:"scan_buildconfigs"`
	ScanImageStreams      bool `mapstructure:"scan_imagestreams"`
	UseProjects           bool `mapstructure:"use_projects"`
}

// TenantConfig represents a tenant configuration
type TenantConfig struct {
	Name        string   `mapstructure:"name"`
	Description string   `mapstructure:"description"`
	Clusters    []string `mapstructure:"clusters"`
}

// ScanScopeConfig defines what levels to scan
type ScanScopeConfig struct {
	ClusterLevel   bool `mapstructure:"cluster_level"`
	NamespaceLevel bool `mapstructure:"namespace_level"`
	TenantLevel    bool `mapstructure:"tenant_level"`
}

// FilterConfig defines namespace and label filters
type FilterConfig struct {
	IncludeNamespaces []string          `mapstructure:"include_namespaces"`
	ExcludeNamespaces []string          `mapstructure:"exclude_namespaces"`
	IncludeLabels     map[string]string `mapstructure:"include_labels"`
	ExcludeLabels     map[string]string `mapstructure:"exclude_labels"`
}

// DiscoveryConfig defines service discovery settings
type DiscoveryConfig struct {
	ScanDeployments           bool     `mapstructure:"scan_deployments"`
	ScanStatefulSets          bool     `mapstructure:"scan_statefulsets"`
	ScanDaemonSets            bool     `mapstructure:"scan_daemonsets"`
	ScanHelmReleases          bool     `mapstructure:"scan_helm_releases"`
	AutoDiscoverDependencies  bool     `mapstructure:"auto_discover_dependencies"`
	DependencySources         []string `mapstructure:"dependency_sources"`
}

// OutputConfig defines output settings
type OutputConfig struct {
	DefaultFormat        string `mapstructure:"default_format"`
	ColorEnabled         bool   `mapstructure:"color_enabled"`
	ShowScopeIndicators  bool   `mapstructure:"show_scope_indicators"`
	Timezone             string `mapstructure:"timezone"`
}

// StorageConfig defines storage settings
type StorageConfig struct {
	Type          string `mapstructure:"type"`
	Path          string `mapstructure:"path"`
	RetentionDays int    `mapstructure:"retention_days"`
}

// Load loads the configuration from the specified path
func Load(configPath string) (*Config, error) {
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("could not determine home directory: %w", err)
		}
		v.AddConfigPath(filepath.Join(home, ".knetz"))
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	// Set defaults
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	// Expand paths
	if err := expandPaths(&config); err != nil {
		return nil, fmt.Errorf("could not expand paths: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	v.SetDefault("scan_scope.cluster_level", true)
	v.SetDefault("scan_scope.namespace_level", true)
	v.SetDefault("scan_scope.tenant_level", true)

	v.SetDefault("discovery.scan_deployments", true)
	v.SetDefault("discovery.scan_statefulsets", true)
	v.SetDefault("discovery.scan_daemonsets", true)
	v.SetDefault("discovery.scan_helm_releases", true)
	v.SetDefault("discovery.auto_discover_dependencies", true)

	v.SetDefault("output.default_format", "table")
	v.SetDefault("output.color_enabled", true)
	v.SetDefault("output.show_scope_indicators", true)
	v.SetDefault("output.timezone", "UTC")

	v.SetDefault("storage.type", "sqlite")
	v.SetDefault("storage.retention_days", 90)
}

// expandPaths expands ~ and environment variables in paths
func expandPaths(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Expand storage path
	if config.Storage.Path != "" {
		config.Storage.Path = expandPath(config.Storage.Path, home)
	} else {
		config.Storage.Path = filepath.Join(home, ".knetz", "data.db")
	}

	// Expand kubeconfig paths
	for i := range config.Clusters {
		if config.Clusters[i].Kubeconfig != "" {
			config.Clusters[i].Kubeconfig = expandPath(config.Clusters[i].Kubeconfig, home)
		}
	}

	return nil
}

// expandPath expands ~ to home directory
func expandPath(path, home string) string {
	if len(path) > 0 && path[0] == '~' {
		return filepath.Join(home, path[1:])
	}
	return os.ExpandEnv(path)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if len(c.Clusters) == 0 {
		return fmt.Errorf("no clusters configured")
	}

	// Validate clusters
	for _, cluster := range c.Clusters {
		if cluster.Name == "" {
			return fmt.Errorf("cluster name is required")
		}
		if cluster.Context == "" {
			return fmt.Errorf("cluster context is required for cluster %s", cluster.Name)
		}
	}

	// Validate tenants
	for _, tenant := range c.Tenants {
		if tenant.Name == "" {
			return fmt.Errorf("tenant name is required")
		}
	}

	return nil
}

// GetCluster returns the cluster configuration by name
func (c *Config) GetCluster(name string) *ClusterConfig {
	for _, cluster := range c.Clusters {
		if cluster.Name == name {
			return &cluster
		}
	}
	return nil
}

// GetTenant returns the tenant configuration by name
func (c *Config) GetTenant(name string) *TenantConfig {
	for _, tenant := range c.Tenants {
		if tenant.Name == name {
			return &tenant
		}
	}
	return nil
}

// GetClustersByTenant returns all clusters for a given tenant
func (c *Config) GetClustersByTenant(tenantName string) []ClusterConfig {
	var clusters []ClusterConfig
	for _, cluster := range c.Clusters {
		if cluster.Tenant == tenantName {
			clusters = append(clusters, cluster)
		}
	}
	return clusters
}

