package discovery

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	
	"github.com/knetz-io/knetz/internal/models"
	"github.com/knetz-io/knetz/pkg/cluster"
	"github.com/knetz-io/knetz/pkg/config"
	"github.com/knetz-io/knetz/pkg/storage"
)

// Scanner discovers services in a Kubernetes cluster
type Scanner struct {
	clusterConfig *config.ClusterConfig
	config        *config.Config
	storage       *storage.Storage
	client        *cluster.Client
}

// NewScanner creates a new scanner
func NewScanner(clusterCfg *config.ClusterConfig, cfg *config.Config, store *storage.Storage) *Scanner {
	return &Scanner{
		clusterConfig: clusterCfg,
		config:        cfg,
		storage:       store,
	}
}

// Scan scans the cluster for services
func (s *Scanner) Scan(ctx context.Context) ([]*models.Service, error) {
	// Create cluster client
	client, err := cluster.NewClient(s.clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create cluster client: %w", err)
	}
	s.client = client

	var allServices []*models.Service

	// Determine namespaces to scan
	namespacesToScan := s.clusterConfig.Namespaces
	if len(namespacesToScan) == 0 {
		// Scan all namespaces
		nsList, err := s.client.GetNamespaces(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not list namespaces: %w", err)
		}
		namespacesToScan = nsList
	}

	// Apply filters
	namespacesToScan = s.filterNamespaces(namespacesToScan)

	// Scan each namespace
	for _, ns := range namespacesToScan {
		services, err := s.scanNamespace(ctx, ns)
		if err != nil {
			return nil, fmt.Errorf("error scanning namespace %s: %w", ns, err)
		}
		allServices = append(allServices, services...)
	}

	return allServices, nil
}

// filterNamespaces applies namespace filters
func (s *Scanner) filterNamespaces(namespaces []string) []string {
	var filtered []string

	for _, ns := range namespaces {
		// Check exclude list
		excluded := false
		for _, excludeNS := range s.config.Filters.ExcludeNamespaces {
			if ns == excludeNS {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		// If include list is specified, check it
		if len(s.config.Filters.IncludeNamespaces) > 0 {
			included := false
			for _, includeNS := range s.config.Filters.IncludeNamespaces {
				if ns == includeNS {
					included = true
					break
				}
			}
			if !included {
				continue
			}
		}

		filtered = append(filtered, ns)
	}

	return filtered
}

// scanNamespace scans a single namespace
func (s *Scanner) scanNamespace(ctx context.Context, namespace string) ([]*models.Service, error) {
	var services []*models.Service

	// Scan deployments
	if s.config.Discovery.ScanDeployments {
		depServices, err := s.scanDeployments(ctx, namespace)
		if err != nil {
			return nil, err
		}
		services = append(services, depServices...)
	}

	// Scan StatefulSets
	if s.config.Discovery.ScanStatefulSets {
		stsServices, err := s.scanStatefulSets(ctx, namespace)
		if err != nil {
			return nil, err
		}
		services = append(services, stsServices...)
	}

	// Scan DaemonSets
	if s.config.Discovery.ScanDaemonSets {
		dsServices, err := s.scanDaemonSets(ctx, namespace)
		if err != nil {
			return nil, err
		}
		services = append(services, dsServices...)
	}

	return services, nil
}

// scanDeployments scans deployments in a namespace
func (s *Scanner) scanDeployments(ctx context.Context, namespace string) ([]*models.Service, error) {
	clientset := s.client.GetClientset()
	deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var services []*models.Service
	for _, dep := range deployments.Items {
		service := s.deploymentToService(&dep, namespace)
		services = append(services, service)
	}

	return services, nil
}

// scanStatefulSets scans statefulsets in a namespace
func (s *Scanner) scanStatefulSets(ctx context.Context, namespace string) ([]*models.Service, error) {
	clientset := s.client.GetClientset()
	statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var services []*models.Service
	for _, sts := range statefulSets.Items {
		service := s.statefulSetToService(&sts, namespace)
		services = append(services, service)
	}

	return services, nil
}

// scanDaemonSets scans daemonsets in a namespace
func (s *Scanner) scanDaemonSets(ctx context.Context, namespace string) ([]*models.Service, error) {
	clientset := s.client.GetClientset()
	daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var services []*models.Service
	for _, ds := range daemonSets.Items {
		service := s.daemonSetToService(&ds, namespace)
		services = append(services, service)
	}

	return services, nil
}

// deploymentToService converts a Deployment to a Service model
func (s *Scanner) deploymentToService(dep interface{}, namespace string) *models.Service {
	// Type assertion would go here in real implementation
	// For now, creating a basic service structure
	
	name := "example-service" // Would extract from dep
	version := s.extractVersion(nil, nil)
	
	id := s.generateServiceID(name, namespace, s.clusterConfig.Name)
	
	return &models.Service{
		ID:          id,
		Name:        name,
		Version:     version,
		ClusterName: s.clusterConfig.Name,
		Namespace:   namespace,
		TenantName:  s.clusterConfig.Tenant,
		Type:        "deployment",
		Platform:    "kubernetes",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// statefulSetToService converts a StatefulSet to a Service model
func (s *Scanner) statefulSetToService(sts interface{}, namespace string) *models.Service {
	name := "example-statefulset"
	version := s.extractVersion(nil, nil)
	id := s.generateServiceID(name, namespace, s.clusterConfig.Name)
	
	return &models.Service{
		ID:          id,
		Name:        name,
		Version:     version,
		ClusterName: s.clusterConfig.Name,
		Namespace:   namespace,
		TenantName:  s.clusterConfig.Tenant,
		Type:        "statefulset",
		Platform:    "kubernetes",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// daemonSetToService converts a DaemonSet to a Service model
func (s *Scanner) daemonSetToService(ds interface{}, namespace string) *models.Service {
	name := "example-daemonset"
	version := s.extractVersion(nil, nil)
	id := s.generateServiceID(name, namespace, s.clusterConfig.Name)
	
	return &models.Service{
		ID:          id,
		Name:        name,
		Version:     version,
		ClusterName: s.clusterConfig.Name,
		Namespace:   namespace,
		TenantName:  s.clusterConfig.Tenant,
		Type:        "daemonset",
		Platform:    "kubernetes",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// extractVersion extracts version from labels, annotations, or image tag
func (s *Scanner) extractVersion(labels, annotations map[string]string) string {
	// Priority:
	// 1. app.kubernetes.io/version label
	// 2. version label
	// 3. version annotation
	// 4. image tag (would be extracted from container spec)
	
	if labels != nil {
		if v, ok := labels["app.kubernetes.io/version"]; ok {
			return v
		}
		if v, ok := labels["version"]; ok {
			return v
		}
	}
	
	if annotations != nil {
		if v, ok := annotations["version"]; ok {
			return v
		}
	}
	
	return "unknown"
}

// generateServiceID generates a unique ID for a service
func (s *Scanner) generateServiceID(name, namespace, cluster string) string {
	data := fmt.Sprintf("%s/%s/%s", cluster, namespace, name)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

