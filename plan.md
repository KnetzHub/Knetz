# Cross-Cluster Dependency & Version Manager

## Product Plan & Development Roadmap

---

## 📋 Executive Summary

A standalone CLI tool that provides global visibility, dependency tracking, and version management across multiple Kubernetes clusters. This tool addresses the complexity of managing microservices dependencies and version drift in multi-cluster, multi-namespace, and multi-tenant environments.

---

## 🎯 Problem Statement

### Current Challenges in Modern Cloud-Native Organizations

**Microservices Sprawl**
- Dozens to hundreds of services, each with independent version lifecycles
- No unified view of what versions are running where
- Tribal knowledge and spreadsheets as primary tracking mechanisms

**Multi-Cluster Reality**
- Multiple environments: staging, production, regional clusters, hybrid cloud
- Each cluster running different versions of the same service
- No standardized way to track version consistency across clusters

**API Version Drift**
- Service A in Cluster-1 runs v1.3, but Cluster-2 runs v1.1
- Breaking changes propagate unexpectedly across clusters
- Consumers fail due to incompatible API versions

**Dependency Management Complexity**
- Service A depends on Service B >=v2.0, but some clusters still run v1.9
- Dependencies span across namespaces and clusters
- No visibility into cross-cluster dependency chains
- Circular dependencies go undetected

**Risky Rollouts**
- Coordinating upgrades across multiple clusters is complex and error-prone
- No clear understanding of blast radius for updates
- Breaking changes discovered after deployment

### Existing Tool Limitations

- **Helm/ArgoCD/Flux**: Solve deployment automation but don't provide global dependency & version awareness
- **Service Meshes**: Handle runtime traffic but don't validate version compatibility
- **Monitoring Tools**: Show what's running but don't validate dependencies or suggest safe upgrade paths

---

## 💡 Solution Overview

### What This Tool Provides

**1. Multi-Cluster Service Inventory**
- Connect to multiple Kubernetes clusters via kubeconfig contexts
- Automatically discover and catalog all services (Deployments, StatefulSets, Helm releases, OpenShift DeploymentConfigs)
- Track versions across clusters, namespaces, and tenants
- Support for AWS EKS, GCP GKE, Azure AKS, Red Hat OpenShift, and vanilla Kubernetes

**2. Multi-Level Scope Management**
- **Tenant Level**: Logical grouping of clusters by organization/business unit
- **Cluster Level**: Individual Kubernetes cluster tracking
- **Namespace Level**: Granular service isolation within clusters
- Cross-scope comparison and aggregation

**3. Intelligent Dependency Discovery**
- Manual declaration via YAML specs or annotations
- Automatic discovery from:
  - Helm charts and Chart.yaml
  - Environment variables in manifests
  - ConfigMap/Secret references
  - Service mesh configurations (Istio, Linkerd)
  - Network policies and ingress rules
  - Prometheus metrics (traffic patterns)
- Confidence scoring for auto-discovered dependencies

**4. Version Drift Detection**
- Compare versions across clusters within the same tenant
- Detect version skew within a cluster across namespaces
- Identify incompatible dependency versions
- Highlight services violating version constraints

**5. Dependency Graph & Analysis**
- Build comprehensive dependency graphs with scope awareness
- Detect circular dependencies (within and across scopes)
- Calculate dependency chains and impact radius
- Visualize relationships at tenant, cluster, and namespace levels

**6. Safe Rollout Planning**
- Suggest optimal upgrade order based on dependency graph
- Show blast radius for planned updates
- Validate compatibility before deployment
- Generate rollout recommendations

**7. Multi-Format Reporting**
- CLI matrix view showing versions across all scopes
- Dependency tree visualization
- Mismatch and violation reports
- Export to JSON, YAML, CSV, Markdown, and SVG graphs

---

## 🏗️ Phase 1: Core Standalone Tool (MVP)

### 1. Project Setup & Foundation

**Tasks:**
- [x] Initialize Go project with Cobra CLI framework
- [x] Setup project structure and module organization
- [x] Configure build system and Makefile
- [x] Setup semantic versioning for tool releases
- [x] Create comprehensive README with quickstart
- [x] Add LICENSE file (Apache 2.0 recommended for CNCF)
- [x] Setup CI/CD pipeline (GitHub Actions)

**Deliverables:**
- ✅ Working Go project structure
- ✅ Build and release automation
- ✅ Documentation foundation

---

### 2. Configuration System

**Tasks:**
- [x] Design configuration file format (YAML-based)
- [x] Implement config file parser and loader
- [x] Build multi-kubeconfig support with context mapping
- [x] Create cluster connection settings
- [x] Add namespace-level configuration with filters
- [x] Add tenant/organization grouping support
- [x] Add output format preferences
- [x] Implement configuration validation with helpful error messages
- [x] Build `init` command to generate starter config

**Configuration Structure:**

```yaml
# ~/.knetz/config.yaml

# Cluster definitions with multi-kubeconfig support
clusters:
  - name: prod-us-east
    kubeconfig: ~/.kube/config
    context: prod-us-east-context
    namespaces: [default, production, shared-services]
    tenant: acme-corp
    
  - name: prod-eu-west
    kubeconfig: ~/.kube/config-eu
    context: prod-eu-context
    namespaces: [default, production]
    tenant: acme-corp
    
  - name: staging
    kubeconfig: ~/.kube/config
    context: staging-context
    namespaces: [staging, dev, qa]
    tenant: acme-corp
    
  - name: prod-asia
    kubeconfig: ~/.kube/config-asia
    context: prod-asia-context
    namespaces: [production, shared-services]
    tenant: beta-company
    
  - name: openshift-prod
    kubeconfig: ~/.kube/config-openshift
    context: openshift-prod-context
    namespaces: [production, staging]  # OpenShift projects
    tenant: acme-corp
    platform: openshift
    api_url: https://api.openshift.example.com:6443

# Tenant/Organization grouping
tenants:
  - name: acme-corp
    description: "Main production tenant"
    clusters: [prod-us-east, prod-eu-west, staging, openshift-prod]
    
  - name: beta-company
    description: "Beta customer tenant"
    clusters: [prod-asia, staging-asia]

# Scan scope configuration
scan_scope:
  cluster_level: true
  namespace_level: true
  tenant_level: true
  
# Namespace filters
filters:
  include_namespaces: [production, staging, shared-services, dev]
  exclude_namespaces: [kube-system, kube-public, kube-node-lease]
  include_labels:
    managed-by: knetz
  exclude_labels:
    ignore: "true"

# Service discovery settings
discovery:
  scan_deployments: true
  scan_statefulsets: true
  scan_daemonsets: true
  scan_helm_releases: true
  scan_openshift_deploymentconfigs: true  # OpenShift-specific
  scan_openshift_routes: true             # OpenShift-specific
  auto_discover_dependencies: true
  dependency_sources:
    - helm_charts
    - env_variables
    - configmaps
    - service_mesh
    - network_policies
    - openshift_routes  # OpenShift-specific

# Output preferences
output:
  default_format: table
  color_enabled: true
  show_scope_indicators: true
  timezone: UTC

# Storage
storage:
  type: sqlite
  path: ~/.knetz/data.db
  retention_days: 90
```

**Deliverables:**
- ✅ Configuration parser and validator
- ✅ Multi-kubeconfig support
- ✅ Tenant and cluster grouping
- ✅ CLI command: `knetz init`

---

### 3. Multi-Cluster Connection

**Tasks:**
- [x] Implement kubeconfig parser (multiple files support)
- [x] Build cluster authentication manager
- [x] Support multiple authentication providers:
  - AWS EKS (IAM authenticator)
  - GCP GKE (gcloud credentials)
  - Azure AKS (Azure AD)
  - Red Hat OpenShift (OAuth token, service account)
  - Vanilla Kubernetes (certificate/token auth)
- [x] Create cluster connection pool with reuse
- [x] Implement cluster health check and connectivity test
- [x] Add support for multiple kubeconfig contexts
- [x] Build namespace discovery per cluster (including OpenShift projects)
- [x] Implement tenant-based cluster grouping
- [x] Add connection retry logic with exponential backoff
- [x] Add timeout configuration
- [x] Implement connection caching
- [x] Support OpenShift-specific resources (DeploymentConfig, BuildConfig, Routes)

**OpenShift-Specific Considerations:**
- Projects (OpenShift namespaces) support
- DeploymentConfig detection (in addition to Deployments)
- OpenShift Routes (in addition to Ingress)
- ImageStreams and BuildConfigs for version tracking
- OAuth authentication flow
- OpenShift-specific labels and annotations

**Deliverables:**
- ✅ Robust multi-cluster connection manager
- ✅ Support for all major cloud providers and OpenShift
- ✅ Connection health monitoring
- ✅ OpenShift-aware resource detection
- ✅ CLI command: `knetz cluster test`

---

### 4. Scope Management System

**Tasks:**
- [x] Design multi-level scope hierarchy (Tenant → Cluster → Namespace)
- [x] Implement cluster-level inventory tracking
- [x] Implement namespace-level inventory tracking
- [x] Implement tenant-level aggregation
- [x] Build scope selector and filter engine
- [x] Create scope validation logic
- [x] Add cross-scope comparison support
- [x] Build scope isolation boundaries
- [x] Implement scope-aware caching
- [x] Create scope resolution logic for queries

**Scope Hierarchy:**

```
Tenant (acme-corp)
├── Cluster (prod-us-east)
│   ├── Namespace (production)
│   │   ├── service-a v1.3.0
│   │   └── service-b v2.0.1
│   ├── Namespace (shared-services)
│   │   └── auth-service v3.1.0
│   └── Namespace (staging)
│       └── service-a v1.2.0
├── Cluster (prod-eu-west)
│   └── Namespace (production)
│       ├── service-a v1.3.0
│       └── service-b v1.9.5
└── Cluster (staging)
    └── Namespace (dev)
        └── service-a v1.4.0-beta
```

**Deliverables:**
- ✅ Hierarchical scope system
- ✅ Scope-aware filtering and querying
- ✅ Cross-scope comparison engine

---

### 5. Service Discovery & Inventory

**Tasks:**
- [x] Build Kubernetes Deployment collector
- [x] Build OpenShift DeploymentConfig collector
- [x] Implement Helm release detector (using Helm storage secrets/configmaps)
- [x] Create StatefulSet collector
- [x] Create DaemonSet collector
- [x] Detect OpenShift Routes (in addition to Ingress)
- [x] Extract version information from:
  - Labels (app.kubernetes.io/version, version)
  - Annotations (version, app.version)
  - Container image tags
  - Helm chart metadata
- [x] Build service metadata extractor
- [x] Implement manifest parser
- [x] Add namespace tagging to all services
- [x] Add cluster tagging to all services
- [x] Add tenant tagging to all services
- [x] Build multi-namespace scanner with parallel execution
- [x] Implement namespace isolation tracking
- [x] Create cross-namespace service discovery
- [x] Add service type classification

**Version Extraction Priority:**
1. `app.kubernetes.io/version` label
2. `version` label
3. `version` annotation
4. Container image tag (if semantic version)
5. Helm chart version

**Deliverables:**
- ✅ Comprehensive service discovery engine
- ✅ Multi-source version extraction
- ✅ Scope-tagged service inventory
- ✅ CLI command: `knetz scan --cluster <name>`

---

### 6. Data Model & Storage

**Tasks:**
- [x] Design core data structures (Service, Cluster, Namespace, Tenant, Version)
- [x] Implement SQLite storage layer with scope indexing
- [x] Create database schema with cluster/namespace/tenant tables
- [x] Build CRUD operations for services with scope context
- [x] Implement cluster registry
- [x] Implement namespace registry per cluster
- [x] Implement tenant registry
- [x] Add version history tracking (per scope level)
- [x] Build scope-based query interface
- [x] Create indexes for efficient scope-based queries
- [x] Implement data retention policies
- [x] Add database migration system

**Data Model (Internal Representation):**

```go
// Pseudo-Go structures for internal data model

type Service struct {
    ID          string
    Name        string
    Version     string
    ClusterName string
    Namespace   string
    TenantName  string
    Type        string // deployment, statefulset, daemonset, helm, deploymentconfig (OpenShift)
    Labels      map[string]string
    Annotations map[string]string
    ImageTag    string
    Metadata    ServiceMetadata
    Dependencies []Dependency
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Platform    string // kubernetes, openshift
}

type ServiceMetadata struct {
    Replicas    int
    Status      string
    Helm        *HelmInfo
    OpenShift   *OpenShiftInfo
}

type HelmInfo struct {
    ChartName    string
    ChartVersion string
    ReleaseName  string
}

type OpenShiftInfo struct {
    DeploymentConfigName string
    BuildConfigName      string
    ImageStreamName      string
    RouteName            string
    RouteHost            string
}

type Cluster struct {
    Name       string
    Context    string
    Kubeconfig string
    TenantName string
    Namespaces []string
    Status     string // healthy, unhealthy, unreachable
    LastScan   time.Time
    Metadata   ClusterMetadata
}

type ClusterMetadata struct {
    Provider      string // eks, gke, aks, openshift, vanilla
    Version       string // Kubernetes/OpenShift version
    Region        string
    ServicesCount int
    IsOpenShift   bool
}

type Namespace struct {
    Name        string
    ClusterName string
    Services    []string
    Labels      map[string]string
    Status      string
}

type Tenant struct {
    Name         string
    Description  string
    Clusters     []string
    TotalServices int
    CreatedAt    time.Time
}

type Dependency struct {
    ServiceName string
    Version     string
    Namespace   string
    Cluster     string
    Tenant      string
    Required    bool
    Type        string // internal, external
    Source      string // manual, helm, envvar, service-mesh
    Confidence  float64 // 0.0 to 1.0 for auto-discovered
}
```

**Database Schema (SQLite):**

```sql
-- Tables structure (pseudo-SQL)

CREATE TABLE tenants (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP
);

CREATE TABLE clusters (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    context TEXT,
    kubeconfig_path TEXT,
    tenant_id TEXT REFERENCES tenants(id),
    status TEXT,
    provider TEXT,
    k8s_version TEXT,
    last_scan TIMESTAMP,
    created_at TIMESTAMP,
    UNIQUE(name, tenant_id)
);

CREATE TABLE namespaces (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cluster_id TEXT REFERENCES clusters(id),
    labels JSON,
    status TEXT,
    created_at TIMESTAMP,
    UNIQUE(name, cluster_id)
);

CREATE TABLE services (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT,
    type TEXT,
    namespace_id TEXT REFERENCES namespaces(id),
    cluster_id TEXT REFERENCES clusters(id),
    tenant_id TEXT REFERENCES tenants(id),
    labels JSON,
    annotations JSON,
    image_tag TEXT,
    metadata JSON,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE dependencies (
    id TEXT PRIMARY KEY,
    service_id TEXT REFERENCES services(id),
    target_service_name TEXT,
    target_version TEXT,
    target_namespace TEXT,
    target_cluster TEXT,
    target_tenant TEXT,
    required BOOLEAN,
    type TEXT,
    source TEXT,
    confidence REAL,
    created_at TIMESTAMP
);

CREATE TABLE version_history (
    id TEXT PRIMARY KEY,
    service_id TEXT REFERENCES services(id),
    old_version TEXT,
    new_version TEXT,
    changed_at TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_services_cluster ON services(cluster_id);
CREATE INDEX idx_services_namespace ON services(namespace_id);
CREATE INDEX idx_services_tenant ON services(tenant_id);
CREATE INDEX idx_services_name_version ON services(name, version);
CREATE INDEX idx_dependencies_service ON dependencies(service_id);
```

**Deliverables:**
- ✅ Complete data model
- ✅ SQLite storage with migrations
- ✅ Efficient query interface
- ✅ Version history tracking

---

### 7. Dependency Declaration System

**Tasks:**
- [x] Design dependency.yaml specification with scope awareness
- [x] Implement annotation-based dependency parser
- [x] Create dependency validation logic (cross-namespace, cross-cluster)
- [x] Build dependency resolver with scope context
- [x] Implement semantic versioning comparator
- [x] Add ConfigMap-based dependency support

**Auto-Discovery Implementation:**
- [x] Extract from Helm charts (Chart.yaml, values.yaml, requirements.yaml)
- [x] Extract from environment variables in Deployment/DeploymentConfig specs
- [x] Extract from ConfigMap/Secret references
- [x] Extract from Service mesh configs (Istio VirtualService, DestinationRule, Linkerd ServiceProfile)
- [x] Extract from Kubernetes Network Policies
- [x] Extract from Ingress/Gateway rules
- [x] Extract from OpenShift Routes
- [x] Extract from Service object selectors and endpoints
- [ ] Optional: Extract from Prometheus metrics (actual service-to-service traffic) - Phase 2
- [ ] Optional: Scan application code repositories (go.mod, package.json, pom.xml, requirements.txt) - Phase 2

**Auto-Discovery Tasks:**
- [x] Build multi-source dependency collector with plugins
- [x] Implement confidence scoring algorithm for discovered dependencies
- [x] Create dependency merge and deduplication logic
- [x] Build cross-validation between sources
- [x] Implement fallback to manual declaration
- [x] Add dependency source tracking

**Dependency Specification Format:**

```yaml
# Example: service-a-dependency.yaml
# Can be stored as ConfigMap in cluster or in Git repo

apiVersion: knetz.io/v1
kind: ServiceDependency
metadata:
  name: service-a
  version: "1.3.0"
  namespace: production
  cluster: prod-us-east
  tenant: acme-corp
  
spec:
  # Same namespace, same cluster dependency
  dependencies:
    - name: service-b
      version: ">=2.0.0"
      namespace: production
      cluster: prod-us-east
      required: true
      description: "Service B provides data processing"
      
    # Different namespace, same cluster
    - name: auth-service
      version: "^3.0.0"
      namespace: shared-services
      cluster: prod-us-east
      required: true
      description: "Authentication service for user validation"
      
    # Cross-cluster dependency
    - name: payment-gateway
      version: ">=1.5.0"
      namespace: production
      cluster: prod-eu-west
      tenant: acme-corp
      required: false
      description: "EU payment processing service"
      
    # External dependency (outside cluster)
    - name: postgres
      version: ">=13.0.0"
      type: external
      required: true
      description: "PostgreSQL database"
      
    - name: redis
      version: ">=6.0.0"
      type: external
      required: true
      
  # Cross-tenant dependencies (for shared services)
  cross_tenant_dependencies:
    - name: shared-api
      tenant: platform-team
      cluster: shared-prod
      namespace: apis
      version: ">=3.0.0"
      required: true
      description: "Platform API gateway"
      
  # API contracts (for future compatibility checking)
  api_contracts:
    - path: /api/v1/users
      method: GET
      schema_version: "1.0"
      
    - path: /api/v1/orders
      method: POST
      schema_version: "2.0"
```

**Annotation-Based Declaration:**

```yaml
# In Kubernetes Deployment manifest
apiVersion: apps/v1
kind: Deployment
metadata:
  name: service-a
  namespace: production
  annotations:
    knetz.io/dependencies: |
      [
        {"name": "service-b", "version": ">=2.0.0", "namespace": "production"},
        {"name": "auth-service", "version": "^3.0.0", "namespace": "shared-services"}
      ]
    knetz.io/version: "1.3.0"
spec:
  # ... deployment spec
```

**Deliverables:**
- ✅ Dependency declaration format (YAML + annotations)
- ✅ Multi-source auto-discovery engine
- ✅ Dependency resolver with scope awareness
- ✅ Semantic versioning support
- ✅ CLI commands: `knetz deps show/validate/export`

---

### 8. Dependency Graph Builder

**Tasks:**
- [x] Create graph data structure with scope nodes (tenant, cluster, namespace, service)
- [x] Implement dependency relationship mapping:
  - Intra-namespace dependencies
  - Cross-namespace dependencies (same cluster)
  - Cross-cluster dependencies (same tenant)
  - Cross-tenant dependencies
- [x] Build graph traversal algorithms (DFS, BFS)
- [x] Detect circular dependencies at all scope levels
- [x] Calculate dependency chains and depth
- [x] Build cluster-level aggregated dependency graph
- [x] Build namespace-level dependency graph
- [x] Build tenant-level aggregated dependency graph
- [x] Implement cross-scope dependency tracking
- [x] Add graph pruning and filtering
- [x] Calculate service impact scores (how many services depend on this)

**Graph Types:**

```
1. Service-Level Graph (most granular)
   service-a@prod-us-east:production → service-b@prod-us-east:production
   service-a@prod-us-east:production → auth-service@prod-us-east:shared-services
   
2. Namespace-Level Graph
   prod-us-east:production → prod-us-east:shared-services
   prod-us-east:production → prod-eu-west:production
   
3. Cluster-Level Graph
   prod-us-east → prod-eu-west
   
4. Tenant-Level Graph
   acme-corp → platform-team (cross-tenant dependency)
```

**Deliverables:**
- ✅ Multi-level dependency graph builder
- ✅ Circular dependency detection
- ✅ Impact analysis engine
- ✅ CLI command: `knetz graph --tenant <name> --output graph.svg`
- ✅ Transitive dependency calculation
- ✅ Depth and impact scoring

---

### 9. Version Drift Detection

**Tasks:**
- [ ] Implement cross-cluster version comparator
- [ ] Build cross-namespace version comparator (within cluster)
- [ ] Build cross-tenant version comparator
- [ ] Create version mismatch detection algorithm at each scope level
- [ ] Build compatibility checker with scope awareness
- [ ] Implement violation scoring system (severity: critical, warning, info)
- [ ] Build dependency constraint validator (cross-scope)
- [ ] Detect version skew across clusters in same tenant
- [ ] Detect version inconsistencies within cluster across namespaces
- [ ] Implement version drift thresholds (e.g., >2 minor versions = warning)
- [ ] Add custom drift rules engine

**Drift Detection Rules:**

```yaml
# Example drift rules
drift_rules:
  - name: production-consistency
    description: "All production namespaces should run same version"
    scope: tenant
    rule: |
      All services in namespaces matching 'production' 
      across all clusters must be within 1 patch version
    severity: critical
    
  - name: staging-ahead-of-prod
    description: "Staging should not be behind production"
    scope: cluster
    rule: |
      Service version in 'staging' namespace must be >= 
      service version in 'production' namespace
    severity: warning
    
  - name: dependency-satisfaction
    description: "All dependencies must satisfy version constraints"
    scope: service
    rule: |
      All declared dependencies must have compatible 
      versions available in target scope
    severity: critical
```

**Deliverables:**
- Multi-level drift detection engine
- Compatibility checker
- Violation reporting
- Custom rules engine
- CLI command: `knetz check --tenant <name>`

---

### 10. CLI Tool Development

**Tasks:**
- [ ] Design CLI command structure with scope flags
- [ ] Implement `init` command (generate starter config)
- [ ] Build `scan` command with scope options
- [ ] Create `diff` command for scope comparison
- [ ] Implement `matrix` command for multi-dimensional view
- [ ] Build `check` command for validation
- [ ] Add `status` command for health summary
- [ ] Create `graph` command for visualization
- [ ] Add `deps` command group for dependency management
- [ ] Add `cluster` command group for cluster operations
- [ ] Implement global flags: `--cluster`, `--namespace`, `--tenant`, `--all`
- [ ] Add output format flags: `--output json|yaml|table|csv`
- [ ] Implement verbose and debug flags
- [ ] Add dry-run mode
- [ ] Build interactive prompts for missing parameters

**CLI Command Structure:**

```bash
# ============================================
# INITIALIZATION
# ============================================

# Generate initial configuration
knetz init
knetz init --config-path ~/.knetz/config.yaml

# ============================================
# CLUSTER MANAGEMENT
# ============================================

# Test cluster connectivity
knetz cluster test
knetz cluster test --cluster prod-us-east

# List all configured clusters
knetz cluster list
knetz cluster list --tenant acme-corp

# Show cluster details
knetz cluster info --cluster prod-us-east

# ============================================
# SCANNING & DISCOVERY
# ============================================

# Scan specific cluster
knetz scan --cluster prod-us-east

# Scan specific namespace in cluster
knetz scan --cluster prod-us-east --namespace production

# Scan all clusters in a tenant
knetz scan --tenant acme-corp

# Scan everything configured
knetz scan --all

# Scan with filters
knetz scan --cluster prod-us-east --include-namespaces production,staging

# ============================================
# DEPENDENCY DISCOVERY
# ============================================

# Discover dependencies automatically
knetz deps discover --cluster prod-us-east

# Show dependencies for a service
knetz deps show --service service-a --cluster prod-us-east --namespace production

# Validate dependency declarations
knetz deps validate --cluster prod-us-east

# Export dependencies
knetz deps export --cluster prod-us-east --output deps.yaml

# ============================================
# VERSION COMPARISON
# ============================================

# Compare versions across clusters
knetz diff --cluster prod-us-east --cluster prod-eu-west

# Compare namespaces within a cluster
knetz diff --cluster prod-us-east --namespace production --namespace staging

# Compare tenants
knetz diff --tenant acme-corp --tenant beta-company

# Compare specific service across all clusters
knetz diff --service service-a --tenant acme-corp

# ============================================
# MATRIX VIEW (Multi-dimensional)
# ============================================

# Matrix view for entire tenant
knetz matrix --tenant acme-corp

# Matrix view for specific clusters
knetz matrix --cluster prod-us-east --cluster prod-eu-west

# Matrix view with filters
knetz matrix --tenant acme-corp --service service-a,service-b

# Export matrix to CSV
knetz matrix --tenant acme-corp --output matrix.csv

# ============================================
# VALIDATION & CHECKING
# ============================================

# Check for version drift and violations
knetz check --tenant acme-corp

# Check specific cluster
knetz check --cluster prod-us-east --namespace production

# Check with custom rules
knetz check --rules ./custom-rules.yaml --tenant acme-corp

# Check specific service dependencies
knetz check --service service-a --cluster prod-us-east

# ============================================
# STATUS & HEALTH
# ============================================

# Overall system status
knetz status

# Status for specific tenant
knetz status --tenant acme-corp

# Status with detailed view
knetz status --tenant acme-corp --verbose

# ============================================
# GRAPH VISUALIZATION
# ============================================

# Generate dependency graph for tenant
knetz graph --tenant acme-corp --output graph.svg

# Graph for specific cluster
knetz graph --cluster prod-us-east --output graph.png

# Graph for specific service and its dependencies
knetz graph --service service-a --cluster prod-us-east --depth 2

# Graph with filtering
knetz graph --tenant acme-corp --exclude-external --output graph.svg

# ============================================
# REPORTING
# ============================================

# Generate drift report
knetz report drift --tenant acme-corp --output report.md

# Generate dependency report
knetz report dependencies --cluster prod-us-east --output deps.html

# Generate summary report
knetz report summary --tenant acme-corp --output summary.json

# ============================================
# OUTPUT FORMATS
# ============================================

# JSON output
knetz scan --cluster prod-us-east --output json

# YAML output
knetz matrix --tenant acme-corp --output yaml

# Table output (default)
knetz scan --cluster prod-us-east --output table

# CSV output
knetz matrix --tenant acme-corp --output csv

# ============================================
# GLOBAL FLAGS
# ============================================

# Specify custom config
knetz scan --config ~/.knetz/custom-config.yaml --cluster prod-us-east

# Verbose output
knetz scan --cluster prod-us-east --verbose

# Debug mode
knetz scan --cluster prod-us-east --debug

# No color output
knetz matrix --tenant acme-corp --no-color

# Dry run (show what would be done)
knetz scan --cluster prod-us-east --dry-run
```

**Deliverables:**
- ✅ Complete CLI tool with all commands
- ✅ Scope-aware filtering
- ✅ Multiple output formats
- ✅ Interactive mode
- ✅ Shell completion scripts

---

### 11. Output & Reporting

**Tasks:**
- [ ] Create ASCII table renderer for matrix view (multi-dimensional)
- [ ] Implement service dependency tree printer with scope indicators
- [ ] Build mismatch summary generator (by cluster, namespace, tenant)
- [ ] Create violation report formatter with scope context
- [ ] Add JSON export with scope metadata
- [ ] Add YAML export with scope metadata
- [ ] Add CSV export with scope columns
- [ ] Implement colored terminal output with scope highlighting
- [ ] Build Markdown report generator with scope sections
- [ ] Create HTML report generator (optional)
- [ ] Create cluster comparison reports
- [ ] Create namespace comparison reports
- [ ] Create tenant-level summary reports
- [ ] Add report templates
- [ ] Implement report scheduling (for future web version)

**Matrix View Example:**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Service Version Matrix - Tenant: acme-corp                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ Service        │ Cluster: prod-us-east       │ Cluster: prod-eu-west        │
│                │ NS: prod    │ NS: staging   │ NS: prod    │ NS: staging   │
│────────────────┼─────────────┼───────────────┼─────────────┼───────────────│
│ service-a      │ v1.3.0 ✓   │ v1.2.0 ⚠️     │ v1.3.0 ✓   │ v1.2.0 ⚠️     │
│ service-b      │ v2.0.1 ✓   │ v2.0.1 ✓     │ v1.9.5 ❌   │ v2.0.0 ✓     │
│ auth-service   │ v3.1.0 ✓   │ v3.1.0 ✓     │ v3.1.0 ✓   │ v3.0.5 ⚠️     │
│ payment-svc    │ v2.5.3 ✓   │ v2.5.3 ✓     │ v2.5.3 ✓   │ v2.5.2 ✓     │
│                                                                              │
│────────────────┼─────────────┼───────────────┼─────────────┼───────────────│
│                                                                              │
│ Cluster: staging                                                            │
│ NS: dev         │ NS: qa                                                    │
│────────────────┼───────────────────────────────────────────────────────────│
│ v1.4.0-beta 🚀 │ v1.3.0 ✓                                                  │
│ v2.1.0-rc 🚀   │ v2.0.1 ✓                                                  │
│ v3.2.0-dev 🚀  │ v3.1.0 ✓                                                  │
│ v2.6.0-beta 🚀 │ v2.5.3 ✓                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

Legend:
✓  Version aligned
⚠️  Minor drift detected
❌  Critical version mismatch
🚀 Pre-release version
```

**Dependency Tree Example:**

```
service-a@prod-us-east:production (v1.3.0)
├─── service-b@prod-us-east:production (v2.0.1) ✓
│    ├─── postgres@external (v13.2) ✓
│    └─── redis@external (v6.2.1) ✓
│
├─── auth-service@prod-us-east:shared-services (v3.1.0) ✓
│    ├─── user-db@external (v14.0) ✓
│    └─── session-cache@external (v6.0) ✓
│
└─── payment-gateway@prod-eu-west:production (v1.9.5) ❌
     └─── VIOLATION: Requires >=1.5.0, but dependency constraint violated
          in cross-cluster context
```

**Violation Report Example:**

```markdown
# Version Drift & Violation Report
**Tenant:** acme-corp
**Generated:** 2025-10-24 10:30:00 UTC
**Scanned Clusters:** 3
**Total Services:** 47

## Summary
- 🔴 **3 Critical violations**
- 🟡 **7 Warnings**
- 🟢 **37 Services healthy**

## Critical Violations

### 1. service-b Version Mismatch (CRITICAL)
**Service:** service-b  
**Issue:** Incompatible versions across production clusters

| Cluster       | Namespace  | Version | Status |
|---------------|------------|---------|--------|
| prod-us-east  | production | v2.0.1  | ✓      |
| prod-eu-west  | production | v1.9.5  | ❌     |

**Impact:** 
- service-a in prod-us-east depends on service-b >=v2.0.0
- Cross-cluster calls may fail due to API incompatibility

**Recommendation:** Upgrade service-b in prod-eu-west to v2.0.1

---

### 2. Circular Dependency Detected (CRITICAL)
**Cycle:** service-c → service-d → service-e → service-c  
**Scope:** prod-us-east:production  

**Services Involved:**
- service-c@prod-us-east:production
- service-d@prod-us-east:production  
- service-e@prod-us-east:production

**Recommendation:** Refactor to break circular dependency

---

## Warnings

### 1. Version Drift in Staging (WARNING)
**Service:** service-a  
**Issue:** Staging running older version than production

| Cluster       | Namespace  | Version | Expected |
|---------------|------------|---------|----------|
| prod-us-east  | production | v1.3.0  | -        |
| prod-us-east  | staging    | v1.2.0  | >=v1.3.0 |

**Recommendation:** Update staging to v1.3.0 before next production release

---
```

**Deliverables:**
- ✅ Rich terminal output with colors and tables
- ✅ Multi-format export (JSON, YAML, CSV, Markdown, HTML)
- ✅ Matrix view renderer
- ✅ Dependency tree visualizer
- ✅ Violation report generator

---

### 12. Testing & Documentation

**Tasks:**
- [ ] Write unit tests for core components (target: 80%+ coverage)
  - Configuration parser
  - Cluster connection manager
  - Service discovery
  - Dependency resolver
  - Version comparator
  - Graph builder
- [ ] Create integration tests with mock Kubernetes clusters (kind/k3s)
- [ ] Build end-to-end test scenarios:
  - Multi-cluster scanning
  - Cross-namespace dependency validation
  - Version drift detection
  - Graph generation
- [ ] Create test fixtures and mock data
- [ ] Write CLI usage documentation with examples
- [ ] Create configuration guide with tenant/cluster/namespace setup
- [ ] Create dependency declaration guide with best practices
- [ ] Build quickstart guide (5-minute setup)
- [ ] Add troubleshooting documentation
- [ ] Create scope management best practices guide
- [ ] Write architecture documentation
- [ ] Create contributing guide (for open source)
- [ ] Add examples directory with sample configs

**Documentation Structure:**

```
docs/
├── README.md                 # Overview and quickstart
├── installation.md           # Installation guide
├── configuration.md          # Configuration reference
├── cli-reference.md          # Complete CLI command reference
├── dependency-spec.md        # Dependency declaration format
├── scope-management.md       # Multi-level scope guide
├── auto-discovery.md         # Dependency auto-discovery guide
├── architecture.md           # System architecture
├── troubleshooting.md        # Common issues and solutions
├── best-practices.md         # Usage best practices
├── contributing.md           # Contribution guidelines
└── examples/
    ├── basic-setup/
    ├── multi-tenant/
    ├── cross-cluster/
    └── gitops-integration/
```

**Deliverables:**
- ✅ Comprehensive test suite
- ✅ Complete documentation
- ✅ Examples and tutorials
- ✅ Troubleshooting guide

---

### 13. Build & Distribution

**Tasks:**
- [ ] Setup cross-platform build (Linux/Mac/Windows)
- [ ] Configure build for multiple architectures (amd64, arm64)
- [ ] Create release automation (GitHub Actions)
- [ ] Build installation scripts (curl/wget one-liner)
- [ ] Add Homebrew formula (for macOS)
- [ ] Create Debian/RPM packages
- [ ] Add Docker container image
- [ ] Setup container registry (Docker Hub, GHCR)
- [ ] Add shell completion scripts (bash/zsh/fish) with scope completion
- [ ] Create checksum and signature generation
- [ ] Setup automated release notes generation
- [ ] Add version update notifier in CLI

**Distribution Methods:**

```bash
# Homebrew (macOS/Linux)
brew install knetz

# Curl install script
curl -sSL https://get.knetz.io | bash

# Go install
go install github.com/your-org/knetz@latest

# Docker
docker run --rm -v ~/.kube:/root/.kube knetz/cli scan --all

# Download binary
wget https://github.com/your-org/knetz/releases/latest/download/knetz-linux-amd64
chmod +x knetz-linux-amd64
sudo mv knetz-linux-amd64 /usr/local/bin/knetz

# Debian/Ubuntu
wget https://github.com/your-org/knetz/releases/latest/download/knetz_amd64.deb
sudo dpkg -i knetz_amd64.deb

# RHEL/CentOS
wget https://github.com/your-org/knetz/releases/latest/download/knetz_amd64.rpm
sudo rpm -i knetz_amd64.rpm
```

**Deliverables:**
- ✅ Cross-platform binaries
- ✅ Multiple installation methods
- ✅ Automated releases
- ✅ Shell completions

---

### 14. OpenShift-Specific Features

**Overview:**
Red Hat OpenShift support is built into the tool as a first-class citizen, allowing seamless management of hybrid Kubernetes and OpenShift environments.

**Supported OpenShift Resources:**
- [ ] **DeploymentConfig:** OpenShift's deployment mechanism (in addition to standard Deployments)
- [ ] **Routes:** OpenShift's external access mechanism (in addition to Ingress)
- [ ] **BuildConfig:** CI/CD pipeline definitions for version tracking
- [ ] **ImageStream:** Container image management and versioning
- [ ] **Projects:** OpenShift namespaces with additional RBAC and quotas

**OpenShift Authentication:**
- [ ] OAuth token authentication
- [ ] Service account tokens
- [ ] Certificate-based authentication
- [ ] Integration with OpenShift login sessions

**OpenShift Version Detection:**
- [ ] Extract versions from DeploymentConfig image tags
- [ ] Track versions through ImageStream tags
- [ ] Monitor BuildConfig outputs for version information
- [ ] Support OpenShift-specific labels: `app.openshift.io/version`

**OpenShift Dependency Discovery:**
- [ ] Parse Route specifications for service dependencies
- [ ] Extract dependencies from BuildConfig source references
- [ ] Analyze ImageStream triggers for inter-service dependencies
- [ ] Support OpenShift Templates for dependency declarations

**Cross-Platform Features:**
- [ ] Unified view of Kubernetes Deployments and OpenShift DeploymentConfigs
- [ ] Combine Ingress and Route information in dependency graphs
- [ ] Cross-platform dependency validation (K8s → OpenShift, OpenShift → K8s)
- [ ] Platform-aware reporting with visual indicators

**OpenShift Security:**
- [ ] Respect OpenShift Security Context Constraints (SCC)
- [ ] Support OpenShift multi-tenancy model
- [ ] Integrate with OpenShift RBAC policies
- [ ] Handle OpenShift project isolation

**CLI Examples:**

```bash
# Detect OpenShift cluster automatically
knetz cluster test --cluster openshift-prod
# Output: Platform: OpenShift 4.12, Kubernetes: 1.25

# Scan OpenShift-specific resources
knetz scan --cluster openshift-prod --scan-deploymentconfigs --scan-routes

# Show OpenShift Route information
knetz deps show --service frontend --cluster openshift-prod
# Output includes:
#   Platform: OpenShift
#   DeploymentConfig: frontend-dc
#   Route: frontend-route.apps.example.com
#   ImageStream: frontend:v1.2.0

# Compare Kubernetes and OpenShift clusters
knetz matrix --cluster k8s-prod --cluster openshift-prod --show-platform

# Export with platform info
knetz report summary --tenant hybrid-org --include-platform-details
```

**Configuration Example:**

```yaml
# OpenShift-specific configuration
clusters:
  - name: openshift-prod
    platform: openshift
    kubeconfig: ~/.kube/config-openshift
    context: openshift-prod/admin
    api_url: https://api.openshift.example.com:6443
    oauth_url: https://oauth-openshift.apps.example.com
    namespaces: [production, staging, shared-services]
    tenant: hybrid-org
    
    # OpenShift-specific settings
    openshift:
      scan_deploymentconfigs: true
      scan_routes: true
      scan_buildconfigs: true
      scan_imagestreams: true
      use_projects: true  # Use OpenShift projects instead of namespaces
      
# Discovery settings for OpenShift
discovery:
  openshift_sources:
    - routes
    - buildconfigs
    - imagestreams
    - templates
```

**Deliverables:**
- Full OpenShift support across all features
- Automatic platform detection
- Unified Kubernetes + OpenShift management
- Platform-aware visualizations and reports

---

## 📊 Technical Stack

### Core Technologies

**Backend/CLI**
- **Language:** Go 1.21+ (for performance, concurrency, and CNCF ecosystem fit)
- **CLI Framework:** Cobra (standard in Kubernetes ecosystem)
- **Configuration:** Viper (YAML/JSON config management)

**Kubernetes Integration**
- **Client Library:** client-go (official Kubernetes Go client)
- **OpenShift Integration:** OpenShift client-go / openshift/client-go
- **Helm Integration:** Helm SDK / helm.sh/helm/v3
- **Multi-cluster:** Multiple kubeconfig support with context switching
- **OpenShift Resources:** Support for DeploymentConfig, Route, BuildConfig, ImageStream APIs

**Storage**
- **Database:** SQLite (embedded, no external dependencies for MVP)
- **Migrations:** golang-migrate or custom migration system
- **Optional:** PostgreSQL for enterprise/hosted version

**Data Structures**
- **Graph:** Custom graph implementation or gonum/graph
- **Versioning:** Masterminds/semver (semantic version parsing and comparison)

**Output & Reporting**
- **Tables:** olekukonko/tablewriter or pterm
- **Colors:** fatih/color or pterm
- **JSON/YAML:** encoding/json, gopkg.in/yaml.v3
- **Graph Visualization:** Graphviz DOT format → SVG/PNG

**Testing**
- **Unit Tests:** Go standard testing package
- **Mocking:** golang/mock or testify/mock
- **Kubernetes Testing:** kind (Kubernetes in Docker) or k3s

**Build & CI/CD**
- **Build:** Make, GoReleaser
- **CI/CD:** GitHub Actions
- **Linting:** golangci-lint

---

## 🎯 Success Criteria for MVP

### Functional Requirements
- [ ] Successfully connect to at least 3 Kubernetes clusters simultaneously
- [ ] Discover and inventory services across multiple namespaces
- [ ] Parse and validate dependency declarations (manual + auto-discovered)
- [ ] Generate dependency graph with cycle detection
- [ ] Detect version drift across clusters with severity scoring
- [ ] Produce matrix view showing versions across all scopes
- [ ] Export reports in JSON, YAML, and Markdown formats
- [ ] Support tenant/cluster/namespace hierarchy

### Non-Functional Requirements
- [ ] Scan 100+ services across 5 clusters in < 30 seconds
- [ ] Handle clusters with 1000+ pods without memory issues
- [ ] Provide helpful error messages for connection failures
- [ ] Work offline after initial scan (local database)
- [ ] Cross-platform support (Linux, macOS, Windows)

### User Experience
- [ ] Single binary with no external dependencies
- [ ] Configuration in < 5 minutes for basic setup
- [ ] Clear, actionable output
- [ ] Comprehensive CLI help and examples

---

## 🚀 Phase 2: Enhanced Features (Post-MVP)

### 1. Web UI Dashboard
- Interactive dependency graph visualization (d3.js/cytoscape.js)
- Real-time cluster status monitoring
- Historical version tracking charts
- Drill-down views (tenant → cluster → namespace → service)
- Alerts and notifications dashboard

### 2. GitOps Integration
- Auto-generate PRs to fix version drift
- Integration with ArgoCD/Flux for deployment validation
- Pre-deployment compatibility checks
- Automated rollout suggestions

### 3. Advanced Dependency Features
- API schema validation (OpenAPI/Protobuf)
- Runtime traffic analysis integration (from service mesh)
- ML-based dependency confidence improvement
- Dependency health scoring

### 4. Policy Engine
- Custom policy definitions (OPA/Rego)
- Policy-as-code for version constraints
- Automated policy enforcement in CI/CD
- Compliance reporting

### 5. Alerting & Notifications
- Slack/Teams/PagerDuty integration
- Webhook support for custom integrations
- Email notifications for drift detection
- Scheduled reports

### 6. Service Mesh Integration
- Deep integration with Istio/Linkerd/Consul
- Traffic-based dependency discovery
- Canary deployment compatibility checking
- Service mesh configuration validation

### 7. Observability Integration
- Prometheus metrics export
- Grafana dashboard templates
- Jaeger/Tempo trace analysis for dependencies
- Log aggregation for dependency errors

### 8. Enterprise Features
- Multi-tenancy with RBAC
- PostgreSQL backend with replication
- Audit logging
- SSO/SAML integration
- API server for programmatic access

---

## 📈 CNCF Sandbox Proposal Path

### Why This Tool Fits CNCF

**1. Addresses Ecosystem Gap**
- No existing CNCF project focuses on cross-cluster dependency management
- Complements existing tools (Helm, ArgoCD, Flux) rather than competing
- Fills the "global visibility" gap in Kubernetes ecosystem

**2. Cloud-Native Principles**
- Multi-cluster by design (cloud-native reality)
- Kubernetes-native (uses standard APIs and patterns)
- Portable and platform-agnostic
- Declarative dependency specification

**3. Community Value**
- Solves real pain points for organizations running Kubernetes at scale
- Useful for any organization with >1 cluster (broad applicability)
- Reduces operational complexity and risk
- Improves system reliability

**4. Incremental Growth Path**
- Start small: CLI tool with core features
- Grow organically: add UI, integrations, advanced features
- Community-driven: open source from day one

### Sandbox Application Checklist

- [ ] **Open Source License:** Apache 2.0 (CNCF preferred)
- [ ] **Public Repository:** GitHub with clear structure
- [ ] **Documentation:** README, architecture, user guide
- [ ] **Governance:** GOVERNANCE.md with decision-making process
- [ ] **Code of Conduct:** Adopt CNCF Code of Conduct
- [ ] **Contributing Guide:** CONTRIBUTING.md
- [ ] **Maintainers:** List of initial maintainers (MAINTAINERS.md)
- [ ] **Roadmap:** PUBLIC_ROADMAP.md
- [ ] **Logo & Brand:** Design project logo and branding
- [ ] **Demo:** Working demo video or live demo
- [ ] **Presentation:** Slide deck for CNCF TAG presentation
- [ ] **Community:** Start building community (Slack, discussions)

### Sandbox Proposal Structure

```markdown
# CNCF Sandbox Proposal: Cross-Cluster Dependency Manager

## Problem Statement
[Describe the multi-cluster dependency tracking problem]

## Solution
[Describe what the tool does]

## Current State
- MVP completed with core features
- X GitHub stars, Y contributors
- Used by Z organizations
- Active community on Slack

## Alignment with CNCF
- Complements: Helm, ArgoCD, Flux, Kubernetes
- Cloud-native architecture
- Multi-cloud and hybrid cloud support

## Roadmap
- Phase 1: [Completed] CLI tool with core features
- Phase 2: [Q1 2026] Web UI and advanced features
- Phase 3: [Q3 2026] Enterprise features and integrations

## Governance
- Open governance model
- Maintainers from multiple organizations
- Transparent decision-making

## Community
- Monthly community meetings
- Active Slack channel
- Regular release cycle

## Ask
We request CNCF Sandbox status to:
- Gain visibility in cloud-native community
- Access CNCF resources and infrastructure
- Collaborate with other CNCF projects
```

---

## 📅 Implementation Timeline

### Month 1-2: Foundation
- Week 1-2: Project setup, configuration system
- Week 3-4: Multi-cluster connection, scope management
- Week 5-6: Service discovery and inventory
- Week 7-8: Data model and storage

### Month 3-4: Core Features
- Week 9-10: Dependency declaration system
- Week 11-12: Auto-discovery implementation
- Week 13-14: Dependency graph builder
- Week 15-16: Version drift detection

### Month 5-6: CLI & Polish
- Week 17-18: CLI command implementation
- Week 19-20: Output and reporting
- Week 21-22: Testing and documentation
- Week 23-24: Build, distribution, MVP release

### Month 7-9: Community & Enhancement
- Month 7: Community building, feedback gathering
- Month 8: Bug fixes, performance optimization
- Month 9: Phase 2 planning, CNCF Sandbox preparation

---

## 🎯 Key Differentiators

### vs. Existing Tools

**vs. Helm**
- Helm manages individual releases, not global dependencies
- Our tool: Cross-cluster visibility and dependency validation

**vs. ArgoCD/Flux**
- GitOps tools deploy, but don't validate cross-cluster dependencies
- Our tool: Pre-deployment compatibility checking

**vs. Service Mesh (Istio/Linkerd)**
- Service mesh handles runtime traffic, not version compatibility
- Our tool: Design-time dependency validation and version management

**vs. Monitoring (Prometheus/Grafana)**
- Monitoring shows current state, not dependency relationships
- Our tool: Proactive dependency management and drift detection

**vs. Spreadsheets/Scripts**
- Manual tracking is error-prone and doesn't scale
- Our tool: Automated, real-time, comprehensive

### Unique Value Propositions

1. **Multi-Level Scope Awareness:** Tenant → Cluster → Namespace hierarchy
2. **Intelligent Auto-Discovery:** Multiple sources with confidence scoring
3. **Cross-Scope Dependencies:** Handle dependencies across namespaces, clusters, and tenants
4. **Version Drift Detection:** Proactive identification of version mismatches
5. **Dependency Graph Analysis:** Visual and programmatic dependency understanding
6. **Safe Rollout Planning:** Suggest upgrade order based on dependency graph
7. **Zero Infrastructure:** Standalone CLI with no servers or agents required
8. **Multi-Platform Support:** Works with Kubernetes, OpenShift, and all major cloud providers (EKS, GKE, AKS)

---

## 📝 Example Use Cases

### Use Case 1: Multi-Region Deployment
**Scenario:** E-commerce company with US, EU, and Asia clusters

**Challenge:**
- Payment service v2.0 deployed in US cluster
- EU cluster still running v1.8, missing critical features
- Order service in EU fails due to API incompatibility

**Solution:**
```bash
# Scan all production clusters
knetz scan --tenant production --cluster us-prod --cluster eu-prod --cluster asia-prod

# Compare versions across regions
knetz diff --service payment-service --tenant production

# Check for violations
knetz check --tenant production

# Output shows:
# ❌ CRITICAL: payment-service version mismatch
#    us-prod: v2.0.1
#    eu-prod: v1.8.3 (order-service requires >=2.0.0)
#    asia-prod: v2.0.0
```

### Use Case 2: Safe Microservices Upgrade
**Scenario:** Upgrading auth-service from v2.x to v3.x with breaking changes

**Challenge:**
- 15 services depend on auth-service
- Need to identify which services are compatible with v3.x
- Need safe rollout order

**Solution:**
```bash
# Check current dependencies
knetz deps show --service auth-service --cluster prod

# Shows:
# auth-service@prod:shared-services (v2.5.0)
# ├─ Depended by:
# │  ├─ user-service (requires ^2.0.0) ⚠️
# │  ├─ order-service (requires >=2.3.0) ⚠️
# │  ├─ api-gateway (requires >=2.0.0) ⚠️
# │  └─ ... 12 more services

# Test compatibility before upgrade
knetz check --what-if auth-service=v3.0.0 --cluster prod

# Output suggests:
# 1. Update user-service constraint to ^3.0.0
# 2. Update order-service constraint to >=3.0.0
# 3. Then upgrade auth-service to v3.0.0
```

### Use Case 3: Staging-to-Production Validation
**Scenario:** Before promoting staging to production, ensure version compatibility

**Challenge:**
- Staging has newer versions than production
- Need to ensure dependencies are satisfied

**Solution:**
```bash
# Compare staging vs production
knetz diff --cluster prod --cluster staging

# Check if staging versions can be safely promoted
knetz check --cluster staging --validate-against prod

# Matrix view
knetz matrix --cluster prod --cluster staging
```

### Use Case 4: Cross-Namespace Service Discovery
**Scenario:** New developer needs to understand service architecture

**Challenge:**
- Services spread across multiple namespaces
- Dependencies not documented

**Solution:**
```bash
# Auto-discover all dependencies in cluster
knetz deps discover --cluster prod --all-namespaces

# Generate visual dependency graph
knetz graph --cluster prod --output architecture.svg

# Export documentation
knetz report dependencies --cluster prod --output README.md
```

### Use Case 5: Hybrid Kubernetes & OpenShift Environment
**Scenario:** Organization running both vanilla Kubernetes and OpenShift clusters

**Challenge:**
- Different deployment mechanisms (Deployments vs DeploymentConfigs)
- Different networking (Ingress vs Routes)
- Need unified view across both platforms
- Services in Kubernetes cluster depend on services in OpenShift cluster

**Solution:**
```bash
# Scan both Kubernetes and OpenShift clusters
knetz scan --cluster k8s-prod --cluster openshift-prod

# Compare versions across platforms
knetz diff --cluster k8s-prod --cluster openshift-prod

# Show cross-platform dependencies
knetz deps show --service frontend-app --cluster k8s-prod

# Output shows:
# frontend-app@k8s-prod:production (v2.1.0)
# ├─ Depends on:
# │  ├─ backend-api@k8s-prod:production (requires >=3.0.0) ✓
# │  ├─ auth-service@openshift-prod:shared-services (requires ^2.5.0) ✓
# │  └─ payment-api@openshift-prod:production (requires >=1.8.0) ⚠️
#
# ⚠️  Cross-platform dependency detected
#     OpenShift Route: payment-api-route.apps.openshift.example.com

# Matrix view showing both platforms
knetz matrix --tenant hybrid-org --output hybrid-services.csv

# Generate unified dependency graph
knetz graph --tenant hybrid-org --include-platform-info --output unified.svg
```

**Benefits:**
- Single tool for both Kubernetes and OpenShift
- Automatic detection of DeploymentConfigs and Routes
- Cross-platform dependency validation
- Unified reporting and visualization

---

## 🔐 Security Considerations

### Data Security
- [ ] Kubeconfig credentials stored securely (use OS keychain)
- [ ] No sensitive data in logs
- [ ] Database encryption at rest (optional)
- [ ] TLS for all cluster connections

### Access Control
- [ ] Respect Kubernetes RBAC (only access allowed namespaces)
- [ ] Respect OpenShift RBAC and Security Context Constraints (SCC)
- [ ] Support OpenShift Projects and multi-tenancy model
- [ ] Read-only operations by default
- [ ] No cluster modification without explicit flags
- [ ] Audit logging for compliance

### Best Practices
- [ ] Follow OWASP secure coding guidelines
- [ ] Regular security audits
- [ ] Dependency vulnerability scanning (Dependabot)
- [ ] Supply chain security (sign releases)

---

## 📚 References & Inspiration

### Existing Tools (Partial Solutions)
- **Helm:** Package management and versioning (single cluster)
- **ArgoCD/Flux:** GitOps deployment (per cluster)
- **Kyverno/OPA:** Policy enforcement (runtime)
- **Crossplane:** Infrastructure as code
- **Backstage:** Service catalog (manual tracking)

### CNCF Projects to Integrate With
- Kubernetes (core API)
- Helm (release detection)
- Prometheus (metrics-based discovery)
- OpenTelemetry (trace-based discovery)
- Istio/Linkerd (service mesh integration)

### Related Concepts
- Semantic Versioning (semver.org)
- Dependency management in software (npm, maven, go modules)
- Configuration management (Ansible, Terraform)

---

## 🤝 Contributing & Community

### Community Building Strategy
1. **Open Source from Day One:** Public GitHub repository
2. **Clear Communication:** Regular blog posts and demos
3. **Community Meetings:** Monthly video calls
4. **Responsive Maintainers:** Quick issue triage and PR reviews
5. **Good First Issues:** Tag beginner-friendly tasks
6. **Documentation:** Comprehensive and beginner-friendly

### Contribution Areas
- Core feature development
- Auto-discovery plugins for new sources
- Output format renderers
- Documentation and examples
- Testing and bug reports
- Integrations with other tools

---

## 📞 Support & Contact

### Getting Help
- **Documentation:** docs.knetz.io (future)
- **GitHub Issues:** Bug reports and feature requests
- **Slack:** #knetz on CNCF Slack (future)
- **Stack Overflow:** Tag `knetz`

### Roadmap & Feedback
- Public roadmap in GitHub Projects
- Community feedback via GitHub Discussions
- Feature voting system

---

## ✅ Next Steps

### Immediate Actions (Before Development)
1. [ ] Review and approve this plan
2. [ ] Setup GitHub repository structure
3. [ ] Create initial project skeleton
4. [ ] Setup development environment
5. [ ] Create project tracking board (GitHub Projects)

### Week 1 Goals
1. [ ] Initialize Go project with Cobra
2. [ ] Implement basic config system
3. [ ] Create `init` command
4. [ ] Write first test cases

### Month 1 Goals
1. [ ] Complete foundation (setup, config, connection, scope)
2. [ ] Successfully connect to multiple test clusters
3. [ ] Basic service discovery working
4. [ ] Initial database schema implemented

---

## 📄 Appendix

### Glossary
- **Tenant:** Logical grouping of clusters (usually by organization or business unit)
- **Scope:** Level of granularity (tenant, cluster, namespace, service)
- **Drift:** Version mismatch across scopes
- **Dependency Chain:** Transitive dependencies (A → B → C)
- **Blast Radius:** Impact scope of a change

### Sample Config Files
See inline examples throughout this document for:
- Configuration file (`config.yaml`)
- Dependency specification (`dependency.yaml`)
- Annotation-based dependencies

### Architecture Diagram
```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Interface (Cobra)                    │
└───────────────┬─────────────────────────────────────────────┘
                │
    ┌───────────┴───────────┬──────────────┬─────────────┐
    │                       │              │             │
┌───▼────────┐   ┌─────────▼─────┐   ┌───▼─────┐  ┌───▼────────┐
│ Config     │   │ Multi-Cluster │   │ Scope   │  │ Dependency │
│ Manager    │   │ Connector     │   │ Manager │  │ Resolver   │
└────────────┘   └───────────────┘   └─────────┘  └────────────┘
                         │
            ┌────────────┴────────────┐
            │                         │
    ┌───────▼────────┐      ┌────────▼────────┐
    │ Service        │      │ Dependency      │
    │ Discovery      │      │ Discovery       │
    └────────────────┘      └─────────────────┘
            │                         │
            └────────────┬────────────┘
                         │
                ┌────────▼─────────┐
                │  Storage Layer   │
                │    (SQLite)      │
                └──────────────────┘
                         │
            ┌────────────┴────────────┐
            │                         │
    ┌───────▼────────┐      ┌────────▼────────┐
    │ Graph Builder  │      │ Drift Detector  │
    └────────────────┘      └─────────────────┘
            │                         │
            └────────────┬────────────┘
                         │
                ┌────────▼─────────┐
                │ Report Generator │
                │ & Output Formatter
                └──────────────────┘
```

---

**Document Version:** 1.0  
**Last Updated:** October 24, 2025  
**Status:** Draft - Awaiting Review

---

*This plan is a living document and will be updated as the project evolves.*

