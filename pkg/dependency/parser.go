package dependency

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
	
	"github.com/knetz-io/knetz/internal/models"
)

// Spec represents a dependency specification
type Spec struct {
	APIVersion string       `yaml:"apiVersion" json:"apiVersion"`
	Kind       string       `yaml:"kind" json:"kind"`
	Metadata   SpecMetadata `yaml:"metadata" json:"metadata"`
	Spec       SpecContent  `yaml:"spec" json:"spec"`
}

// SpecMetadata contains metadata about the service
type SpecMetadata struct {
	Name      string `yaml:"name" json:"name"`
	Version   string `yaml:"version" json:"version"`
	Namespace string `yaml:"namespace" json:"namespace"`
	Cluster   string `yaml:"cluster" json:"cluster"`
	Tenant    string `yaml:"tenant" json:"tenant"`
}

// SpecContent contains the dependency specifications
type SpecContent struct {
	Dependencies            []DependencySpec `yaml:"dependencies" json:"dependencies"`
	CrossTenantDependencies []DependencySpec `yaml:"cross_tenant_dependencies,omitempty" json:"cross_tenant_dependencies,omitempty"`
	APIContracts            []APIContract    `yaml:"api_contracts,omitempty" json:"api_contracts,omitempty"`
}

// DependencySpec represents a single dependency
type DependencySpec struct {
	Name        string `yaml:"name" json:"name"`
	Version     string `yaml:"version" json:"version"`
	Namespace   string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Cluster     string `yaml:"cluster,omitempty" json:"cluster,omitempty"`
	Tenant      string `yaml:"tenant,omitempty" json:"tenant,omitempty"`
	Required    bool   `yaml:"required" json:"required"`
	Type        string `yaml:"type,omitempty" json:"type,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

// APIContract represents an API contract
type APIContract struct {
	Path          string `yaml:"path" json:"path"`
	Method        string `yaml:"method" json:"method"`
	SchemaVersion string `yaml:"schema_version" json:"schema_version"`
}

// Parser handles dependency specification parsing
type Parser struct{}

// NewParser creates a new dependency parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseYAML parses a YAML dependency specification
func (p *Parser) ParseYAML(data []byte) (*Spec, error) {
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("could not parse YAML: %w", err)
	}

	if err := p.validate(&spec); err != nil {
		return nil, fmt.Errorf("invalid spec: %w", err)
	}

	return &spec, nil
}

// ParseJSON parses a JSON dependency specification
func (p *Parser) ParseJSON(data []byte) (*Spec, error) {
	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("could not parse JSON: %w", err)
	}

	if err := p.validate(&spec); err != nil {
		return nil, fmt.Errorf("invalid spec: %w", err)
	}

	return &spec, nil
}

// ParseAnnotation parses dependency information from Kubernetes annotations
func (p *Parser) ParseAnnotation(annotations map[string]string) ([]models.Dependency, error) {
	depJSON, ok := annotations["knetz.io/dependencies"]
	if !ok {
		return nil, nil
	}

	var deps []struct {
		Name      string  `json:"name"`
		Version   string  `json:"version"`
		Namespace string  `json:"namespace,omitempty"`
		Cluster   string  `json:"cluster,omitempty"`
		Required  bool    `json:"required,omitempty"`
		Type      string  `json:"type,omitempty"`
	}

	if err := json.Unmarshal([]byte(depJSON), &deps); err != nil {
		return nil, fmt.Errorf("could not parse annotation: %w", err)
	}

	var dependencies []models.Dependency
	for _, d := range deps {
		dependencies = append(dependencies, models.Dependency{
			ServiceName: d.Name,
			Version:     d.Version,
			Namespace:   d.Namespace,
			Cluster:     d.Cluster,
			Required:    d.Required,
			Type:        d.Type,
			Source:      "annotation",
			Confidence:  1.0,
		})
	}

	return dependencies, nil
}

// validate validates a dependency specification
func (p *Parser) validate(spec *Spec) error {
	if spec.APIVersion == "" {
		return fmt.Errorf("apiVersion is required")
	}

	if spec.Kind != "ServiceDependency" {
		return fmt.Errorf("kind must be ServiceDependency, got %s", spec.Kind)
	}

	if spec.Metadata.Name == "" {
		return fmt.Errorf("metadata.name is required")
	}

	if spec.Metadata.Version == "" {
		return fmt.Errorf("metadata.version is required")
	}

	// Validate dependencies
	for i, dep := range spec.Spec.Dependencies {
		if dep.Name == "" {
			return fmt.Errorf("dependency[%d].name is required", i)
		}
		if dep.Version == "" {
			return fmt.Errorf("dependency[%d].version is required", i)
		}
	}

	return nil
}

// ToModels converts a Spec to internal dependency models
func (p *Parser) ToModels(spec *Spec) []models.Dependency {
	var dependencies []models.Dependency

	for _, dep := range spec.Spec.Dependencies {
		dependencies = append(dependencies, models.Dependency{
			ServiceName: dep.Name,
			Version:     dep.Version,
			Namespace:   dep.Namespace,
			Cluster:     dep.Cluster,
			Tenant:      dep.Tenant,
			Required:    dep.Required,
			Type:        dep.Type,
			Source:      "manual",
			Confidence:  1.0,
		})
	}

	for _, dep := range spec.Spec.CrossTenantDependencies {
		dependencies = append(dependencies, models.Dependency{
			ServiceName: dep.Name,
			Version:     dep.Version,
			Namespace:   dep.Namespace,
			Cluster:     dep.Cluster,
			Tenant:      dep.Tenant,
			Required:    dep.Required,
			Type:        "cross-tenant",
			Source:      "manual",
			Confidence:  1.0,
		})
	}

	return dependencies
}

