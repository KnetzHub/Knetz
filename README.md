# Knetz - Cross-Cluster Dependency & Version Manager

![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8.svg)

A standalone CLI tool that provides global visibility, dependency tracking, and version management across multiple Kubernetes clusters.

## 🎯 Overview

Knetz addresses the complexity of managing microservices dependencies and version drift in multi-cluster, multi-namespace, and multi-tenant environments. It helps you:

- **Discover** services across multiple Kubernetes/OpenShift clusters
- **Track** version drift and incompatibilities  
- **Validate** dependencies before deployment
- **Visualize** dependency graphs
- **Plan** safe rollouts with dependency awareness

## ✨ Key Features

### Multi-Cluster Service Inventory
- Connect to multiple Kubernetes clusters via kubeconfig contexts
- Support for AWS EKS, GCP GKE, Azure AKS, Red Hat OpenShift, and vanilla Kubernetes
- Automatic service discovery across clusters and namespaces

### Intelligent Dependency Discovery
- Manual declaration via YAML specs or annotations
- Automatic discovery from Helm charts, environment variables, ConfigMaps, and service mesh configs
- Confidence scoring for auto-discovered dependencies

### Version Drift Detection
- Compare versions across clusters within the same tenant
- Detect version skew across namespaces
- Identify incompatible dependency versions
- Highlight services violating version constraints

### Dependency Graph & Analysis
- Build comprehensive dependency graphs with scope awareness
- Detect circular dependencies
- Calculate dependency chains and impact radius
- Multi-level visualization (tenant → cluster → namespace → service)

### Safe Rollout Planning
- Suggest optimal upgrade order based on dependency graph
- Show blast radius for planned updates
- Validate compatibility before deployment

## 🚀 Quick Start

### Installation

```bash
# Using Go
go install github.com/knetz-io/knetz/cmd/knetz@latest

# Or build from source
git clone https://github.com/knetz-io/knetz.git
cd knetz
make build
sudo mv bin/knetz /usr/local/bin/
```

### Initialize Configuration

```bash
# Generate starter config
knetz init

# Edit the config file
vim ~/.knetz/config.yaml
```

### Basic Usage

```bash
# Test cluster connectivity
knetz cluster test

# Scan a specific cluster
knetz scan --cluster prod-us-east

# Scan all configured clusters
knetz scan --all

# Compare versions across clusters
knetz diff --cluster prod-us-east --cluster prod-eu-west

# Check for version drift and violations
knetz check --tenant acme-corp

# Generate dependency graph
knetz graph --tenant acme-corp --output graph.svg
```

## 📋 Requirements

- Go 1.23 or higher
- Access to one or more Kubernetes clusters (via kubeconfig)
- Appropriate RBAC permissions to list resources in target clusters

## 🏗️ Architecture

Knetz uses a hierarchical scope model:

```
Tenant (Organization)
├── Cluster (Kubernetes/OpenShift)
│   ├── Namespace (Logical grouping)
│   │   ├── Service A (v1.3.0)
│   │   └── Service B (v2.0.1)
│   └── Namespace (Another grouping)
│       └── Service C (v3.1.0)
└── Cluster (Another environment)
    └── ...
```

## 📖 Documentation

- [Configuration Guide](docs/configuration.md)
- [CLI Reference](docs/cli-reference.md)
- [Dependency Specification](docs/dependency-spec.md)
- [Architecture](docs/architecture.md)

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🌟 Project Status

**Phase 1 (MVP)**: ✅ In Active Development

- [x] Project Setup & Foundation
- [ ] Configuration System
- [ ] Multi-Cluster Connection
- [ ] Service Discovery
- [ ] Dependency Management
- [ ] Version Drift Detection
- [ ] CLI Tools
- [ ] Reporting

## 💬 Community

- GitHub Issues: [Report bugs or request features](https://github.com/knetz-io/knetz/issues)
- Discussions: [Ask questions and share ideas](https://github.com/knetz-io/knetz/discussions)

---

**Note**: Knetz is under active development. APIs and features may change as we work towards a stable v1.0 release.

