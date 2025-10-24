package dependency

import (
	"testing"
)

func TestParseYAML(t *testing.T) {
	yamlData := `
apiVersion: knetz.io/v1
kind: ServiceDependency
metadata:
  name: test-service
  version: "1.0.0"
  namespace: default
  cluster: test-cluster
  tenant: test-tenant
spec:
  dependencies:
    - name: dep-service
      version: ">=1.0.0"
      namespace: default
      required: true
`

	parser := NewParser()
	spec, err := parser.ParseYAML([]byte(yamlData))
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if spec.Metadata.Name != "test-service" {
		t.Errorf("expected name 'test-service', got '%s'", spec.Metadata.Name)
	}

	if spec.Metadata.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", spec.Metadata.Version)
	}

	if len(spec.Spec.Dependencies) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(spec.Spec.Dependencies))
	}

	if spec.Spec.Dependencies[0].Name != "dep-service" {
		t.Errorf("expected dependency 'dep-service', got '%s'", spec.Spec.Dependencies[0].Name)
	}
}

func TestValidate(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name        string
		spec        *Spec
		expectError bool
	}{
		{
			name: "valid spec",
			spec: &Spec{
				APIVersion: "knetz.io/v1",
				Kind:       "ServiceDependency",
				Metadata: SpecMetadata{
					Name:    "test",
					Version: "1.0.0",
				},
				Spec: SpecContent{
					Dependencies: []DependencySpec{
						{Name: "dep1", Version: "1.0.0"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "missing apiVersion",
			spec: &Spec{
				Kind: "ServiceDependency",
				Metadata: SpecMetadata{
					Name:    "test",
					Version: "1.0.0",
				},
			},
			expectError: true,
		},
		{
			name: "wrong kind",
			spec: &Spec{
				APIVersion: "knetz.io/v1",
				Kind:       "WrongKind",
				Metadata: SpecMetadata{
					Name:    "test",
					Version: "1.0.0",
				},
			},
			expectError: true,
		},
		{
			name: "missing metadata.name",
			spec: &Spec{
				APIVersion: "knetz.io/v1",
				Kind:       "ServiceDependency",
				Metadata: SpecMetadata{
					Version: "1.0.0",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.validate(tt.spec)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseAnnotation(t *testing.T) {
	parser := NewParser()

	annotations := map[string]string{
		"knetz.io/dependencies": `[
			{"name": "service-b", "version": ">=2.0.0", "namespace": "production", "required": true},
			{"name": "service-c", "version": "^1.0.0", "required": false}
		]`,
	}

	deps, err := parser.ParseAnnotation(annotations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(deps))
	}

	if deps[0].ServiceName != "service-b" {
		t.Errorf("expected 'service-b', got '%s'", deps[0].ServiceName)
	}

	if deps[0].Version != ">=2.0.0" {
		t.Errorf("expected '>=2.0.0', got '%s'", deps[0].Version)
	}

	if deps[0].Confidence != 1.0 {
		t.Errorf("expected confidence 1.0, got %f", deps[0].Confidence)
	}
}

