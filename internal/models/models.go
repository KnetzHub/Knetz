package models

import "time"

// Service represents a service instance in a Kubernetes cluster
type Service struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	ClusterName string            `json:"cluster_name"`
	Namespace   string            `json:"namespace"`
	TenantName  string            `json:"tenant_name"`
	Type        string            `json:"type"` // deployment, statefulset, daemonset, helm, deploymentconfig
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	ImageTag    string            `json:"image_tag"`
	Metadata    ServiceMetadata   `json:"metadata"`
	Dependencies []Dependency     `json:"dependencies"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Platform    string            `json:"platform"` // kubernetes, openshift
}

// ServiceMetadata contains additional service information
type ServiceMetadata struct {
	Replicas  int            `json:"replicas"`
	Status    string         `json:"status"`
	Helm      *HelmInfo      `json:"helm,omitempty"`
	OpenShift *OpenShiftInfo `json:"openshift,omitempty"`
}

// HelmInfo contains Helm-specific information
type HelmInfo struct {
	ChartName    string `json:"chart_name"`
	ChartVersion string `json:"chart_version"`
	ReleaseName  string `json:"release_name"`
}

// OpenShiftInfo contains OpenShift-specific information
type OpenShiftInfo struct {
	DeploymentConfigName string `json:"deployment_config_name,omitempty"`
	BuildConfigName      string `json:"build_config_name,omitempty"`
	ImageStreamName      string `json:"image_stream_name,omitempty"`
	RouteName            string `json:"route_name,omitempty"`
	RouteHost            string `json:"route_host,omitempty"`
}

// Cluster represents a Kubernetes/OpenShift cluster
type Cluster struct {
	Name       string          `json:"name"`
	Context    string          `json:"context"`
	Kubeconfig string          `json:"kubeconfig"`
	TenantName string          `json:"tenant_name"`
	Namespaces []string        `json:"namespaces"`
	Status     string          `json:"status"` // healthy, unhealthy, unreachable
	LastScan   time.Time       `json:"last_scan"`
	Metadata   ClusterMetadata `json:"metadata"`
}

// ClusterMetadata contains additional cluster information
type ClusterMetadata struct {
	Provider      string `json:"provider"` // eks, gke, aks, openshift, vanilla
	Version       string `json:"version"`  // Kubernetes/OpenShift version
	Region        string `json:"region"`
	ServicesCount int    `json:"services_count"`
	IsOpenShift   bool   `json:"is_openshift"`
}

// Namespace represents a Kubernetes namespace or OpenShift project
type Namespace struct {
	Name        string            `json:"name"`
	ClusterName string            `json:"cluster_name"`
	Services    []string          `json:"services"`
	Labels      map[string]string `json:"labels"`
	Status      string            `json:"status"`
}

// Tenant represents a logical grouping of clusters
type Tenant struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Clusters      []string  `json:"clusters"`
	TotalServices int       `json:"total_services"`
	CreatedAt     time.Time `json:"created_at"`
}

// Dependency represents a service dependency
type Dependency struct {
	ServiceName string  `json:"service_name"`
	Version     string  `json:"version"`
	Namespace   string  `json:"namespace"`
	Cluster     string  `json:"cluster"`
	Tenant      string  `json:"tenant"`
	Required    bool    `json:"required"`
	Type        string  `json:"type"`   // internal, external
	Source      string  `json:"source"` // manual, helm, envvar, service-mesh
	Confidence  float64 `json:"confidence"` // 0.0 to 1.0 for auto-discovered
}

// VersionHistory tracks version changes
type VersionHistory struct {
	ID         string    `json:"id"`
	ServiceID  string    `json:"service_id"`
	OldVersion string    `json:"old_version"`
	NewVersion string    `json:"new_version"`
	ChangedAt  time.Time `json:"changed_at"`
}

// DriftReport represents version drift findings
type DriftReport struct {
	TenantName   string           `json:"tenant_name"`
	GeneratedAt  time.Time        `json:"generated_at"`
	Violations   []DriftViolation `json:"violations"`
	Warnings     []DriftWarning   `json:"warnings"`
	HealthyCount int              `json:"healthy_count"`
}

// DriftViolation represents a critical version mismatch
type DriftViolation struct {
	ServiceName string            `json:"service_name"`
	Severity    string            `json:"severity"` // critical, warning, info
	Issue       string            `json:"issue"`
	Instances   []ServiceInstance `json:"instances"`
	Impact      string            `json:"impact"`
	Recommendation string         `json:"recommendation"`
}

// DriftWarning represents a non-critical version drift
type DriftWarning struct {
	ServiceName string            `json:"service_name"`
	Issue       string            `json:"issue"`
	Instances   []ServiceInstance `json:"instances"`
	Recommendation string         `json:"recommendation"`
}

// ServiceInstance represents a specific instance of a service
type ServiceInstance struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Version     string `json:"version"`
	Status      string `json:"status"`
}

// DependencyGraph represents the complete dependency graph
type DependencyGraph struct {
	Nodes map[string]*GraphNode `json:"nodes"`
	Edges []GraphEdge           `json:"edges"`
}

// GraphNode represents a node in the dependency graph
type GraphNode struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
	Cluster     string `json:"cluster"`
	Namespace   string `json:"namespace"`
	Tenant      string `json:"tenant"`
}

// GraphEdge represents an edge (dependency) in the graph
type GraphEdge struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	Required   bool    `json:"required"`
	Confidence float64 `json:"confidence"`
}

