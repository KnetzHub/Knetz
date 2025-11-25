# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial Knetz MVP implementation
- Cross-cluster dependency tracking
- Version drift detection
- Semantic versioning support
- Dependency graph analysis
- SQLite storage backend
- CLI commands: version, init, cluster, scan, status, diff, matrix, check, deps, graph
- GitHub Actions CI/CD pipelines
- GoReleaser configuration
- Comprehensive documentation

### Features Implemented
- **Multi-Cluster Management:** Connect and manage multiple Kubernetes clusters
- **Service Discovery:** Automatic service discovery from Deployments, StatefulSets, DaemonSets
- **Dependency Management:** YAML-based dependency specification (knetz.io/v1 API)
- **Version Tracking:** Semantic version parsing and comparison
- **Drift Detection:** Identify version mismatches across clusters
- **Graph Analysis:** Dependency graph with cycle detection and impact analysis
- **Tenant Support:** Multi-tenant architecture (Tenant → Cluster → Namespace)
- **OpenShift Configuration:** OpenShift-specific configuration support

### Infrastructure
- Go 1.24 with Cobra CLI framework
- Viper for configuration management
- SQLite for local storage
- Kubernetes client-go integration
- GitHub Actions for CI/CD
- GoReleaser for cross-platform releases

## [0.1.0] - 2025-10-24

### Added
- Initial project structure
- Core MVP functionality
- Basic CLI commands
- Documentation framework

---

**Note:** This changelog is automatically generated from git commits following [Conventional Commits](https://www.conventionalcommits.org/).

### Commit Message Format

We follow the Conventional Commits specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Test additions or updates
- `build`: Build system changes
- `ci`: CI/CD changes
- `chore`: Other changes (maintenance, etc.)

**Breaking Changes:**
Add `!` after type or `BREAKING CHANGE:` in footer for major version bumps.

**Examples:**
```
feat(cluster): add AWS EKS auto-discovery
fix(deps): resolve version comparison edge case
docs(readme): update installation instructions
feat!: change dependency specification format
```


