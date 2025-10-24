package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version
type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Metadata   string
}

var versionRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// ParseVersion parses a version string into a Version struct
func ParseVersion(v string) (*Version, error) {
	matches := versionRegex.FindStringSubmatch(v)
	if matches == nil {
		return nil, fmt.Errorf("invalid version format: %s", v)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: matches[4],
		Metadata:   matches[5],
	}, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		s += "-" + v.PreRelease
	}
	if v.Metadata != "" {
		s += "+" + v.Metadata
	}
	return s
}

// Compare compares two versions
// Returns: -1 if v < other, 0 if v == other, 1 if v > other
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}

	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}

	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}

	// If both have no pre-release, they're equal
	if v.PreRelease == "" && other.PreRelease == "" {
		return 0
	}

	// Version with pre-release is less than version without
	if v.PreRelease == "" {
		return 1
	}
	if other.PreRelease == "" {
		return -1
	}

	// Compare pre-release versions
	return strings.Compare(v.PreRelease, other.PreRelease)
}

// IsGreaterThan checks if v > other
func (v *Version) IsGreaterThan(other *Version) bool {
	return v.Compare(other) > 0
}

// IsLessThan checks if v < other
func (v *Version) IsLessThan(other *Version) bool {
	return v.Compare(other) < 0
}

// IsEqual checks if v == other
func (v *Version) IsEqual(other *Version) bool {
	return v.Compare(other) == 0
}

// CompareVersionStrings compares two version strings
func CompareVersionStrings(v1, v2 string) (int, error) {
	ver1, err := ParseVersion(v1)
	if err != nil {
		return 0, err
	}

	ver2, err := ParseVersion(v2)
	if err != nil {
		return 0, err
	}

	return ver1.Compare(ver2), nil
}

// SatisfiesConstraint checks if a version satisfies a constraint
// Supports: >=, <=, >, <, =, ^, ~
func SatisfiesConstraint(version, constraint string) (bool, error) {
	constraint = strings.TrimSpace(constraint)

	// Handle ^ (caret) - allows changes that do not modify left-most non-zero digit
	if strings.HasPrefix(constraint, "^") {
		return satisfiesCaretConstraint(version, strings.TrimPrefix(constraint, "^"))
	}

	// Handle ~ (tilde) - allows patch-level changes
	if strings.HasPrefix(constraint, "~") {
		return satisfiesTildeConstraint(version, strings.TrimPrefix(constraint, "~"))
	}

	// Handle >=
	if strings.HasPrefix(constraint, ">=") {
		targetVer := strings.TrimPrefix(constraint, ">=")
		cmp, err := CompareVersionStrings(version, targetVer)
		if err != nil {
			return false, err
		}
		return cmp >= 0, nil
	}

	// Handle <=
	if strings.HasPrefix(constraint, "<=") {
		targetVer := strings.TrimPrefix(constraint, "<=")
		cmp, err := CompareVersionStrings(version, targetVer)
		if err != nil {
			return false, err
		}
		return cmp <= 0, nil
	}

	// Handle >
	if strings.HasPrefix(constraint, ">") {
		targetVer := strings.TrimPrefix(constraint, ">")
		cmp, err := CompareVersionStrings(version, targetVer)
		if err != nil {
			return false, err
		}
		return cmp > 0, nil
	}

	// Handle <
	if strings.HasPrefix(constraint, "<") {
		targetVer := strings.TrimPrefix(constraint, "<")
		cmp, err := CompareVersionStrings(version, targetVer)
		if err != nil {
			return false, err
		}
		return cmp < 0, nil
	}

	// Handle = or exact match
	if strings.HasPrefix(constraint, "=") {
		constraint = strings.TrimPrefix(constraint, "=")
	}
	
	cmp, err := CompareVersionStrings(version, constraint)
	if err != nil {
		return false, err
	}
	return cmp == 0, nil
}

// satisfiesCaretConstraint checks if version satisfies ^constraint
func satisfiesCaretConstraint(version, constraint string) (bool, error) {
	ver, err := ParseVersion(version)
	if err != nil {
		return false, err
	}

	target, err := ParseVersion(constraint)
	if err != nil {
		return false, err
	}

	// Must be >= target
	if ver.Compare(target) < 0 {
		return false, nil
	}

	// Major version must match if > 0
	if target.Major > 0 {
		return ver.Major == target.Major, nil
	}

	// If major is 0, minor must match if > 0
	if target.Minor > 0 {
		return ver.Major == 0 && ver.Minor == target.Minor, nil
	}

	// If both major and minor are 0, patch must match
	return ver.Major == 0 && ver.Minor == 0 && ver.Patch == target.Patch, nil
}

// satisfiesTildeConstraint checks if version satisfies ~constraint
func satisfiesTildeConstraint(version, constraint string) (bool, error) {
	ver, err := ParseVersion(version)
	if err != nil {
		return false, err
	}

	target, err := ParseVersion(constraint)
	if err != nil {
		return false, err
	}

	// Must be >= target
	if ver.Compare(target) < 0 {
		return false, nil
	}

	// Major and minor must match
	return ver.Major == target.Major && ver.Minor == target.Minor, nil
}

