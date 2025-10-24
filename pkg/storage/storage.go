package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/knetz-io/knetz/internal/models"
	_ "modernc.org/sqlite"
)

// Storage handles data persistence
type Storage struct {
	db *sql.DB
}

// New creates a new Storage instance
func New(dbPath string) (*Storage, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	storage := &Storage{db: db}

	// Initialize schema
	if err := storage.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("could not initialize schema: %w", err)
	}

	return storage, nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// initSchema creates the database schema
func (s *Storage) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tenants (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS clusters (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		context TEXT,
		kubeconfig_path TEXT,
		tenant_id TEXT,
		status TEXT,
		provider TEXT,
		k8s_version TEXT,
		last_scan TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id),
		UNIQUE(name, tenant_id)
	);

	CREATE TABLE IF NOT EXISTS namespaces (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		cluster_id TEXT,
		labels TEXT,
		status TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (cluster_id) REFERENCES clusters(id),
		UNIQUE(name, cluster_id)
	);

	CREATE TABLE IF NOT EXISTS services (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		version TEXT,
		type TEXT,
		namespace_id TEXT,
		cluster_id TEXT,
		tenant_id TEXT,
		labels TEXT,
		annotations TEXT,
		image_tag TEXT,
		metadata TEXT,
		platform TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (namespace_id) REFERENCES namespaces(id),
		FOREIGN KEY (cluster_id) REFERENCES clusters(id),
		FOREIGN KEY (tenant_id) REFERENCES tenants(id)
	);

	CREATE TABLE IF NOT EXISTS dependencies (
		id TEXT PRIMARY KEY,
		service_id TEXT,
		target_service_name TEXT,
		target_version TEXT,
		target_namespace TEXT,
		target_cluster TEXT,
		target_tenant TEXT,
		required INTEGER,
		type TEXT,
		source TEXT,
		confidence REAL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (service_id) REFERENCES services(id)
	);

	CREATE TABLE IF NOT EXISTS version_history (
		id TEXT PRIMARY KEY,
		service_id TEXT,
		old_version TEXT,
		new_version TEXT,
		changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (service_id) REFERENCES services(id)
	);

	CREATE INDEX IF NOT EXISTS idx_services_cluster ON services(cluster_id);
	CREATE INDEX IF NOT EXISTS idx_services_namespace ON services(namespace_id);
	CREATE INDEX IF NOT EXISTS idx_services_tenant ON services(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_services_name_version ON services(name, version);
	CREATE INDEX IF NOT EXISTS idx_dependencies_service ON dependencies(service_id);
	`

	_, err := s.db.Exec(schema)
	return err
}

// SaveService saves a service to the database
func (s *Storage) SaveService(service *models.Service) error {
	query := `
	INSERT OR REPLACE INTO services (
		id, name, version, type, cluster_id, namespace_id, tenant_id,
		labels, annotations, image_tag, metadata, platform, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		service.ID,
		service.Name,
		service.Version,
		service.Type,
		service.ClusterName,
		service.Namespace,
		service.TenantName,
		"", // labels (JSON)
		"", // annotations (JSON)
		service.ImageTag,
		"", // metadata (JSON)
		service.Platform,
		time.Now(),
	)

	return err
}

// GetServicesByCluster returns all services for a cluster
func (s *Storage) GetServicesByCluster(clusterName string) ([]*models.Service, error) {
	query := `
	SELECT id, name, version, type, cluster_id, namespace_id, tenant_id,
		   image_tag, platform, created_at, updated_at
	FROM services
	WHERE cluster_id = ?
	ORDER BY name, namespace_id`

	rows, err := s.db.Query(query, clusterName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*models.Service
	for rows.Next() {
		service := &models.Service{}
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Version,
			&service.Type,
			&service.ClusterName,
			&service.Namespace,
			&service.TenantName,
			&service.ImageTag,
			&service.Platform,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		service.CreatedAt = createdAt
		service.UpdatedAt = updatedAt
		services = append(services, service)
	}

	return services, rows.Err()
}

// GetServicesByTenant returns all services for a tenant
func (s *Storage) GetServicesByTenant(tenantName string) ([]*models.Service, error) {
	query := `
	SELECT id, name, version, type, cluster_id, namespace_id, tenant_id,
		   image_tag, platform, created_at, updated_at
	FROM services
	WHERE tenant_id = ?
	ORDER BY cluster_id, namespace_id, name`

	rows, err := s.db.Query(query, tenantName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*models.Service
	for rows.Next() {
		service := &models.Service{}
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Version,
			&service.Type,
			&service.ClusterName,
			&service.Namespace,
			&service.TenantName,
			&service.ImageTag,
			&service.Platform,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		service.CreatedAt = createdAt
		service.UpdatedAt = updatedAt
		services = append(services, service)
	}

	return services, rows.Err()
}

// SaveCluster saves cluster information
func (s *Storage) SaveCluster(cluster *models.Cluster) error {
	query := `
	INSERT OR REPLACE INTO clusters (
		id, name, context, kubeconfig_path, tenant_id, status,
		provider, k8s_version, last_scan
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query,
		cluster.Name,
		cluster.Name,
		cluster.Context,
		cluster.Kubeconfig,
		cluster.TenantName,
		cluster.Status,
		cluster.Metadata.Provider,
		cluster.Metadata.Version,
		cluster.LastScan,
	)

	return err
}

// GetCluster retrieves cluster information
func (s *Storage) GetCluster(name string) (*models.Cluster, error) {
	query := `
	SELECT id, name, context, kubeconfig_path, tenant_id, status,
		   provider, k8s_version, last_scan, created_at
	FROM clusters
	WHERE name = ?`

	cluster := &models.Cluster{
		Metadata: models.ClusterMetadata{},
	}
	var createdAt time.Time

	err := s.db.QueryRow(query, name).Scan(
		&cluster.Name,
		&cluster.Name,
		&cluster.Context,
		&cluster.Kubeconfig,
		&cluster.TenantName,
		&cluster.Status,
		&cluster.Metadata.Provider,
		&cluster.Metadata.Version,
		&cluster.LastScan,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// DeleteServicesByCluster deletes all services for a cluster (for re-scanning)
func (s *Storage) DeleteServicesByCluster(clusterName string) error {
	query := `DELETE FROM services WHERE cluster_id = ?`
	_, err := s.db.Exec(query, clusterName)
	return err
}

