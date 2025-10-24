package utils

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		expectError bool
		expected    *Version
	}{
		{
			name:        "simple version",
			version:     "1.2.3",
			expectError: false,
			expected:    &Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:        "version with v prefix",
			version:     "v1.2.3",
			expectError: false,
			expected:    &Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:        "version with prerelease",
			version:     "1.2.3-beta.1",
			expectError: false,
			expected:    &Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "beta.1"},
		},
		{
			name:        "version with metadata",
			version:     "1.2.3+build.123",
			expectError: false,
			expected:    &Version{Major: 1, Minor: 2, Patch: 3, Metadata: "build.123"},
		},
		{
			name:        "invalid version",
			version:     "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseVersion(tt.version)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Major != tt.expected.Major {
				t.Errorf("Major = %d, want %d", result.Major, tt.expected.Major)
			}
			if result.Minor != tt.expected.Minor {
				t.Errorf("Minor = %d, want %d", result.Minor, tt.expected.Minor)
			}
			if result.Patch != tt.expected.Patch {
				t.Errorf("Patch = %d, want %d", result.Patch, tt.expected.Patch)
			}
			if result.PreRelease != tt.expected.PreRelease {
				t.Errorf("PreRelease = %s, want %s", result.PreRelease, tt.expected.PreRelease)
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"equal versions", "1.0.0", "1.0.0", 0},
		{"v1 greater major", "2.0.0", "1.0.0", 1},
		{"v1 greater minor", "1.1.0", "1.0.0", 1},
		{"v1 greater patch", "1.0.1", "1.0.0", 1},
		{"v2 greater major", "1.0.0", "2.0.0", -1},
		{"v2 greater minor", "1.0.0", "1.1.0", -1},
		{"v2 greater patch", "1.0.0", "1.0.1", -1},
		{"prerelease less than release", "1.0.0-beta", "1.0.0", -1},
		{"release greater than prerelease", "1.0.0", "1.0.0-beta", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareVersionStrings(tt.v1, tt.v2)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("CompareVersionStrings(%s, %s) = %d, want %d", 
					tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestSatisfiesConstraint(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		constraint string
		expected   bool
	}{
		{"equal exact", "1.2.3", "1.2.3", true},
		{"not equal", "1.2.3", "1.2.4", false},
		{"greater than", "1.2.3", ">1.2.0", true},
		{"not greater than", "1.2.0", ">1.2.0", false},
		{"greater or equal", "1.2.3", ">=1.2.3", true},
		{"less than", "1.2.0", "<1.2.3", true},
		{"less or equal", "1.2.3", "<=1.2.3", true},
		{"caret major", "1.5.0", "^1.2.0", true},
		{"caret fail", "2.0.0", "^1.2.0", false},
		{"tilde patch", "1.2.5", "~1.2.0", true},
		{"tilde fail minor", "1.3.0", "~1.2.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SatisfiesConstraint(tt.version, tt.constraint)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("SatisfiesConstraint(%s, %s) = %v, want %v", 
					tt.version, tt.constraint, result, tt.expected)
			}
		})
	}
}

