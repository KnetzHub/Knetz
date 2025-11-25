# KnetZ Backend Plan

## Overview
KnetZ is a dependency management platform that tracks service versions across multiple sources (VCS, container registries, Kubernetes platforms). The backend will be built with a MongoDB database, RESTful API server, and CLI integration.

---

## Architecture Components

### 1. API Server
- **Framework**: Node.js (Express/Fastify) or Go (Gin/Echo)
- **Database**: MongoDB
- **Authentication**: Phase 1 (Open/Basic), Phase 2 (OAuth2, SSO)
- **API Style**: RESTful with potential GraphQL for complex queries

### 2. Dashboard Backend
- **Purpose**: Serve frontend, provide real-time logs and events
- **Features**: WebSocket support for real-time updates
- **Auth Integration**: Gmail, GitHub OAuth (Phase 1), SSO (Phase 2)

### 3. CLI
- **Purpose**: Pipeline integration, API client
- **Commands**: `init`, `check`, `update`, `sync`
- **Configuration**: `.knetz` config file support

---

## MongoDB Database Schema

### Collection: `dependencies_hub`

**Purpose**: Store service metadata and version history

```json
{
  "_id": ObjectId,
  "name": String,              // Service name
  "alias": String,             // User-defined alias (e.g., "demo")
  "repository": String,        // Source type: "github", "gitlab", "docker", "quay"
  "repoURL": String,           // Full repository URL
  "repoOwner": String,         // Repository owner/organization
  "repoName": String,          // Repository name
  "visibility": String,        // "public" or "private"
  "trackingMethod": String,    // "cron", "webhook", "manual"
  "lastSyncedAt": ISODate,     // Last sync timestamp
  "createdAt": ISODate,
  "updatedAt": ISODate,
  "userId": ObjectId,          // Reference to user who registered (Phase 1: nullable)
  "versions": [
    {
      "version": String,       // Semantic version (1.2.3)
      "type": String,          // "release", "tag", "commit", "image"
      "releasedAt": ISODate,
      "changelog": String,     // Release notes
      "artifacts": Array,      // Associated artifacts/assets
      "isMajor": Boolean,      // Breaking change indicator
      "isPrerelease": Boolean
    }
  ],
  "dependencies": [
    {
      "name": String,          // Dependency name
      "currentVersion": String,
      "requiredVersion": String,
      "serviceId": ObjectId,   // Reference to dependencies_hub
      "status": String         // "up-to-date", "outdated", "deprecated"
    }
  ],
  "metadata": {
    "language": String,        // Programming language
    "framework": String,
    "packageManager": String,  // npm, pip, maven, go.mod, etc.
    "manifestFile": String,    // package.json, requirements.txt, etc.
    "manifestPath": String     // Path to manifest in repo
  }
}
```

**Indexes**:
```javascript
{
  "repoURL": 1,        // Unique index
  "name": 1,
  "alias": 1,
  "userId": 1,
  "repository": 1,
  "versions.version": 1,
  "lastSyncedAt": 1
}
```

---

### Collection: `users`

**Purpose**: User authentication and authorization (Phase 1: Basic, Phase 2: Full)

```json
{
  "_id": ObjectId,
  "email": String,             // Unique
  "username": String,
  "authProvider": String,      // "gmail", "github", "sso"
  "authProviderId": String,    // External provider ID
  "apiKeys": [
    {
      "key": String,           // Hashed API key
      "name": String,          // Key identifier
      "createdAt": ISODate,
      "expiresAt": ISODate,
      "lastUsed": ISODate
    }
  ],
  "tokens": {
    "github": {
      "accessToken": String,   // Encrypted
      "refreshToken": String,
      "expiresAt": ISODate
    },
    "gitlab": {},
    "docker": {}
  },
  "settings": {
    "defaultTrackingMethod": String,
    "notifications": Boolean,
    "webhookURL": String
  },
  "createdAt": ISODate,
  "updatedAt": ISODate,
  "lastLoginAt": ISODate
}
```

**Indexes**:
```javascript
{
  "email": 1,          // Unique
  "username": 1,       // Unique
  "authProviderId": 1
}
```

---

### Collection: `sync_logs`

**Purpose**: Track sync operations, cron jobs, and API calls

```json
{
  "_id": ObjectId,
  "serviceId": ObjectId,       // Reference to dependencies_hub
  "serviceName": String,
  "syncType": String,          // "cron", "webhook", "manual", "cli"
  "status": String,            // "success", "failed", "partial"
  "triggeredBy": String,       // "system", "user", "webhook"
  "userId": ObjectId,
  "versionsAdded": Number,
  "versionsUpdated": Number,
  "errors": [
    {
      "message": String,
      "code": String,
      "timestamp": ISODate
    }
  ],
  "duration": Number,          // Milliseconds
  "startedAt": ISODate,
  "completedAt": ISODate,
  "metadata": {
    "source": String,
    "trigger": String,
    "rateLimitRemaining": Number
  }
}
```

**Indexes**:
```javascript
{
  "serviceId": 1,
  "startedAt": -1,     // Descending for recent logs
  "status": 1,
  "triggeredBy": 1
}
```

---

### Collection: `events`

**Purpose**: Application events for dashboard and audit trail

```json
{
  "_id": ObjectId,
  "type": String,              // "service.created", "version.detected", "dependency.outdated"
  "severity": String,          // "info", "warning", "error", "critical"
  "serviceId": ObjectId,
  "userId": ObjectId,
  "message": String,
  "details": Object,           // Event-specific data
  "createdAt": ISODate,
  "readBy": [ObjectId],        // Users who marked as read
  "isArchived": Boolean
}
```

**Indexes**:
```javascript
{
  "createdAt": -1,
  "type": 1,
  "severity": 1,
  "userId": 1,
  "isArchived": 1
}
```

---

### Collection: `webhooks`

**Purpose**: Webhook management for GitHub, GitLab, etc. (Phase 1+)

```json
{
  "_id": ObjectId,
  "serviceId": ObjectId,
  "userId": ObjectId,
  "provider": String,          // "github", "gitlab"
  "webhookId": String,         // External webhook ID
  "events": [String],          // ["push", "release", "tag"]
  "secret": String,            // Webhook secret (encrypted)
  "active": Boolean,
  "lastTriggered": ISODate,
  "createdAt": ISODate
}
```

---

## API Endpoints

### Phase 1: Core Endpoints

#### Services Management
```
POST   /api/v1/services                 - Register a new service
GET    /api/v1/services                 - List all services
GET    /api/v1/services/:id             - Get service details
GET    /api/v1/services/:id/versions    - Get version history
PUT    /api/v1/services/:id             - Update service metadata
DELETE /api/v1/services/:id             - Remove service
POST   /api/v1/services/:id/sync        - Trigger manual sync
```

#### Version Management
```
GET    /api/v1/services/:id/versions/:version    - Get specific version
POST   /api/v1/services/:id/versions             - Add version manually
GET    /api/v1/services/:id/latest               - Get latest version
GET    /api/v1/services/:id/compare/:v1/:v2      - Compare versions
```

#### Dependency Checking
```
GET    /api/v1/services/:id/dependencies         - Get dependencies
POST   /api/v1/services/:id/dependencies/check   - Check dependency status
GET    /api/v1/dependencies/outdated             - List all outdated deps
```

#### GitHub Integration
```
POST   /api/v1/integrations/github/connect       - Connect GitHub account
GET    /api/v1/integrations/github/repos         - List accessible repos
POST   /api/v1/integrations/github/import        - Import repo as service
POST   /api/v1/integrations/github/webhook       - Setup webhook
```

#### Logs & Events
```
GET    /api/v1/logs                              - Get sync logs
GET    /api/v1/logs/:id                          - Get specific log
GET    /api/v1/events                            - Get events feed
POST   /api/v1/events/:id/read                   - Mark event as read
```

#### Health & Status
```
GET    /api/v1/health                            - API health check
GET    /api/v1/stats                             - Platform statistics
```

### Phase 2: Additional Endpoints
```
POST   /api/v1/auth/sso                          - SSO authentication
POST   /api/v1/integrations/docker               - Docker registry integration
POST   /api/v1/integrations/gitlab               - GitLab integration
POST   /api/v1/integrations/quay                 - Quay.io integration
GET    /api/v1/services/:id/helm                 - Helm chart versions
POST   /api/v1/services/:id/branches             - Manage branch tracking
```

---

## Implementation Details

### 1. Service Registration Flow

#### For GitHub (Phase 1)

**Public Repository - Cron Method**:
```
1. User calls: POST /api/v1/services
   Body: { "repoURL": "github.com/owner/repo", "alias": "demo" }
2. API validates URL and fetches metadata from GitHub API
3. Create document in dependencies_hub collection
4. Parse package manifest (package.json, go.mod, etc.)
5. Extract current version and dependencies
6. Schedule cron job for periodic updates
7. Return service ID and status
```

**Public Repository - CLI Method**:
```
1. User runs: knetz init --url=github.com/owner/repo --alias=demo
2. CLI calls POST /api/v1/services with credentials
3. API performs registration (same as above)
4. CLI stores service config in local .knetz file
```

**Manual Update Flow**:
```
1. User runs: knetz update <repoURL>
2. CLI calls: POST /api/v1/services/:id/sync
3. API fetches latest releases/tags from GitHub
4. Compare with existing versions, add new ones
5. Update dependencies if manifest changed
6. Log sync operation in sync_logs
7. Emit events for new versions
8. Return sync summary
```

**Check Dependencies Flow**:
```
1. User runs: knetz check <repoURL>
2. CLI calls: GET /api/v1/services/:id/dependencies
3. API retrieves service and dependencies
4. For each dependency, check if newer version exists
5. Calculate status: up-to-date, outdated, deprecated
6. Return dependency tree with status
7. CLI displays formatted output
```

---

### 2. GitHub Integration Architecture

#### Authentication
- OAuth2 flow for user authentication
- Personal Access Tokens (PAT) for CLI
- GitHub App for webhook management (recommended)

#### API Interactions
```javascript
// Fetch releases
GET /repos/:owner/:repo/releases

// Fetch tags
GET /repos/:owner/:repo/tags

// Fetch specific file (manifest)
GET /repos/:owner/:repo/contents/:path

// Create webhook
POST /repos/:owner/:repo/hooks
```

#### Rate Limiting Strategy
- Authenticated requests: 5,000 req/hour
- Cache responses with TTL
- Implement exponential backoff
- Priority queue for user-triggered syncs

---

### 3. Version Detection Logic

```javascript
async function detectVersions(service) {
  const strategy = {
    github: async () => {
      // 1. Fetch all releases
      const releases = await githubAPI.getReleases(service.repoURL);
      
      // 2. Fetch all tags (if no releases)
      const tags = await githubAPI.getTags(service.repoURL);
      
      // 3. Parse semantic versions
      const versions = [...releases, ...tags]
        .map(parseVersion)
        .filter(isValidSemver)
        .sort(semverCompare);
      
      return versions;
    },
    docker: async () => {
      // Fetch image tags from Docker Hub/Registry API
    },
    gitlab: async () => {
      // Similar to GitHub
    }
  };
  
  return await strategy[service.repository]();
}
```

---

### 4. Dependency Parsing

**Supported Manifest Files**:
- `package.json` (Node.js)
- `go.mod` (Go)
- `requirements.txt`, `Pipfile` (Python)
- `pom.xml`, `build.gradle` (Java)
- `Cargo.toml` (Rust)
- `composer.json` (PHP)

**Parsing Flow**:
```javascript
async function parseDependencies(service) {
  // 1. Detect manifest file
  const manifest = await detectManifest(service);
  
  // 2. Fetch manifest content from repo
  const content = await fetchManifestFile(service, manifest.path);
  
  // 3. Parse based on type
  const parser = manifestParsers[manifest.type];
  const deps = parser.parse(content);
  
  // 4. Resolve each dependency
  for (const dep of deps) {
    // Check if dependency exists in our DB
    const depService = await findServiceByName(dep.name);
    
    if (depService) {
      // Link to existing service
      dep.serviceId = depService._id;
      dep.status = compareVersions(dep.version, depService.versions[0].version);
    } else {
      // External dependency (not tracked yet)
      dep.status = "external";
    }
  }
  
  return deps;
}
```

---

### 5. Cron Job Architecture

**Scheduler**: Node-cron or Bull Queue (with Redis)

```javascript
// Sync schedule priorities
const syncSchedules = {
  high: "*/30 * * * *",    // Every 30 minutes (critical services)
  medium: "0 */4 * * *",   // Every 4 hours (standard)
  low: "0 0 * * *"         // Daily (archived/stable)
};

// Cron job implementation
cron.schedule(syncSchedules.medium, async () => {
  const services = await db.dependencies_hub.find({
    trackingMethod: "cron",
    visibility: "public"
  });
  
  for (const service of services) {
    await queueSync(service._id);
  }
});

async function queueSync(serviceId) {
  const log = {
    serviceId,
    syncType: "cron",
    status: "pending",
    triggeredBy: "system",
    startedAt: new Date()
  };
  
  try {
    const versions = await detectVersions(service);
    const newVersions = await addNewVersions(serviceId, versions);
    
    log.status = "success";
    log.versionsAdded = newVersions.length;
    
    // Emit events for new versions
    for (const version of newVersions) {
      await emitEvent({
        type: "version.detected",
        severity: version.isMajor ? "warning" : "info",
        serviceId,
        message: `New version ${version.version} detected`
      });
    }
  } catch (error) {
    log.status = "failed";
    log.errors = [{ message: error.message, timestamp: new Date() }];
  } finally {
    log.completedAt = new Date();
    log.duration = log.completedAt - log.startedAt;
    await db.sync_logs.insertOne(log);
  }
}
```

---

### 6. Webhook Implementation

**GitHub Webhook Handler**:
```javascript
POST /api/v1/webhooks/github/:serviceId

async function handleGitHubWebhook(req, res) {
  // 1. Verify webhook signature
  const signature = req.headers["x-hub-signature-256"];
  if (!verifySignature(req.body, signature)) {
    return res.status(401).send("Invalid signature");
  }
  
  // 2. Parse event type
  const event = req.headers["x-github-event"];
  
  // 3. Handle specific events
  switch (event) {
    case "release":
      await handleReleaseEvent(req.body);
      break;
    case "push":
      await handlePushEvent(req.body);
      break;
    case "create":
      if (req.body.ref_type === "tag") {
        await handleTagEvent(req.body);
      }
      break;
  }
  
  res.status(200).send("OK");
}
```

---

## CLI Implementation

### Configuration File: `.knetz`

```yaml
version: 1.0
api_endpoint: https://api.knetz.io
api_key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
services:
  - name: service-a
    alias: demo
    url: github.com/owner/service-a
    id: 507f1f77bcf86cd799439011
  - name: service-b
    alias: prod
    url: github.com/owner/service-b
    id: 507f1f77bcf86cd799439012
```

### CLI Commands

#### `knetz init`
```bash
knetz init --url=github.com/owner/repo --alias=demo

# Workflow:
# 1. Validate URL format
# 2. Call POST /api/v1/services
# 3. Store service in .knetz config
# 4. Output: "Service registered successfully with ID: xxx"
```

#### `knetz check`
```bash
knetz check --url=github.com/owner/repo
knetz check --alias=demo

# Workflow:
# 1. Resolve service from .knetz or URL
# 2. Call GET /api/v1/services/:id/dependencies
# 3. Display formatted dependency tree:
#
# Service: service-a (v1.2.3)
# Dependencies:
#   ✓ express: 4.18.0 (up-to-date)
#   ⚠ lodash: 4.17.20 (outdated, latest: 4.17.21)
#   ✗ moment: 2.29.1 (deprecated, use date-fns)
```

#### `knetz update`
```bash
knetz update --url=github.com/owner/repo
knetz update --alias=demo
knetz update --all

# Workflow:
# 1. Resolve service(s)
# 2. Call POST /api/v1/services/:id/sync
# 3. Output sync results
```

#### `knetz list`
```bash
knetz list

# Workflow:
# 1. Call GET /api/v1/services
# 2. Display table:
#
# | Name       | Alias | Latest Version | Last Synced    |
# |------------|-------|----------------|----------------|
# | service-a  | demo  | 1.2.3          | 2 hours ago    |
# | service-b  | prod  | 2.0.1          | 1 day ago      |
```

#### `knetz sync` (Phase 2)
```bash
knetz sync

# Workflow:
# 1. Read .knetz file
# 2. Sync all configured services
# 3. Update local cache
```

---

## Technology Stack Recommendations

### Backend API
**Option 1: Node.js**
- Framework: Express.js or Fastify
- ORM: Mongoose (MongoDB)
- Auth: Passport.js, JWT
- Validation: Joi or Zod
- Testing: Jest

**Option 2: Go**
- Framework: Gin or Echo
- MongoDB: Official Go driver
- Auth: golang-jwt
- Validation: validator
- Testing: testify

### Background Jobs
- Bull Queue (Redis-based) for Node.js
- Cron: node-cron or agenda
- Alternative: Temporal.io for complex workflows

### Real-time (Dashboard)
- WebSocket: Socket.io or native WS
- Server-Sent Events (SSE) for log streaming

### CLI
- Node.js: Commander.js or oclif
- Go: Cobra

### Monitoring
- Logging: Winston (Node) or Zap (Go)
- Metrics: Prometheus
- Tracing: OpenTelemetry
- APM: DataDog or New Relic

---

## Development Phases

### Phase 1: Core Features (MVP)
**Timeline**: 8-12 weeks

**Week 1-2: Database & API Foundation**
- [ ] MongoDB schema implementation
- [ ] Basic CRUD API for services
- [ ] API authentication (API keys)
- [ ] Health check endpoints

**Week 3-4: GitHub Integration**
- [ ] GitHub OAuth implementation
- [ ] Fetch releases/tags from GitHub API
- [ ] Parse common manifest files (package.json, go.mod)
- [ ] Version detection and storage

**Week 5-6: Sync Mechanism**
- [ ] Cron job scheduler
- [ ] Manual sync endpoints
- [ ] Sync logging
- [ ] Error handling and retry logic

**Week 7-8: CLI Development**
- [ ] CLI framework setup
- [ ] init, check, update, list commands
- [ ] .knetz config management
- [ ] API client integration

**Week 9-10: Dashboard Backend**
- [ ] User authentication (Gmail, GitHub OAuth)
- [ ] Service management endpoints
- [ ] Events and logs API
- [ ] Real-time WebSocket setup

**Week 11-12: Testing & Documentation**
- [ ] Unit tests (80% coverage)
- [ ] Integration tests
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Deployment setup

### Phase 2: Advanced Features
**Timeline**: 6-8 weeks

**Features**:
- [ ] GitLab integration
- [ ] Docker registry support (Docker Hub, Quay)
- [ ] Private repository support
- [ ] SSO authentication
- [ ] Webhook management UI
- [ ] Advanced dependency analysis
- [ ] Helm chart version tracking
- [ ] Branch-specific tracking
- [ ] Kubernetes platform integrations
- [ ] Email/Slack notifications
- [ ] Dependency vulnerability scanning
- [ ] CLI sync command with .knetz

---

## Security Considerations

### Phase 1
1. **API Keys**: Generate secure random API keys (32+ chars)
2. **Rate Limiting**: Implement per-IP and per-user limits
3. **Input Validation**: Sanitize all inputs (URL, service names)
4. **Secrets Storage**: Encrypt OAuth tokens at rest
5. **HTTPS Only**: Enforce TLS for all API calls

### Phase 2
1. **SSO Integration**: SAML 2.0 or OAuth2/OIDC
2. **RBAC**: Role-based access control for teams
3. **Audit Logs**: Track all sensitive operations
4. **Webhook Signature Verification**: Validate webhook sources
5. **Private Token Encryption**: Use KMS for token storage

---

## Scalability Considerations

1. **Database Indexing**: Optimize queries with proper indexes
2. **Caching Layer**: Redis for frequently accessed data
3. **API Rate Limiting**: Prevent abuse and manage resources
4. **Background Job Queue**: Separate sync operations from API
5. **Horizontal Scaling**: Stateless API servers behind load balancer
6. **Database Sharding**: Partition by userId or serviceId if needed
7. **CDN**: Cache static dashboard assets

---

## Error Handling

### API Error Response Format
```json
{
  "error": {
    "code": "SERVICE_NOT_FOUND",
    "message": "Service with ID 123 not found",
    "details": {},
    "timestamp": "2025-01-15T10:30:00Z",
    "requestId": "req_abc123"
  }
}
```

### Common Error Codes
- `AUTH_FAILED`: Authentication failed
- `INVALID_INPUT`: Validation error
- `SERVICE_NOT_FOUND`: Service doesn't exist
- `RATE_LIMIT_EXCEEDED`: Too many requests
- `SYNC_FAILED`: Sync operation failed
- `GITHUB_API_ERROR`: GitHub API error
- `VERSION_PARSE_ERROR`: Cannot parse version

---

## Monitoring & Observability

### Key Metrics
1. **API Metrics**:
   - Request rate (req/sec)
   - Response time (p50, p95, p99)
   - Error rate (%)
   - Active connections

2. **Sync Metrics**:
   - Sync success rate (%)
   - Sync duration (avg, max)
   - Services synced per hour
   - Sync queue length

3. **Business Metrics**:
   - Total services tracked
   - Active users
   - Version updates detected per day
   - CLI usage stats

### Logging Strategy
- Structured JSON logs
- Log levels: DEBUG, INFO, WARN, ERROR, CRITICAL
- Correlation IDs for request tracing
- Separate logs: access.log, error.log, sync.log

---

## Testing Strategy

### Unit Tests
- Service registration logic
- Version parsing and comparison
- Dependency resolution
- Authentication middleware

### Integration Tests
- API endpoint flows
- GitHub API integration
- Database operations
- Cron job execution

### E2E Tests
- CLI commands
- Full user workflows
- Dashboard interactions

### Performance Tests
- API load testing (Apache Bench, k6)
- Database query optimization
- Sync job throughput

---

## Deployment

### Infrastructure
- **API Server**: Docker containers on Kubernetes/ECS
- **Database**: MongoDB Atlas (managed) or self-hosted replica set
- **Redis**: ElastiCache or managed Redis
- **Load Balancer**: AWS ALB/NLB or Nginx

### CI/CD Pipeline
1. Code pushed to GitHub
2. Run tests (unit, integration)
3. Build Docker image
4. Push to container registry
5. Deploy to staging
6. Run E2E tests
7. Deploy to production (blue-green)

### Environment Variables
```bash
NODE_ENV=production
MONGO_URI=mongodb://...
REDIS_URL=redis://...
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
JWT_SECRET=xxx
API_PORT=3000
LOG_LEVEL=info
```

---

## Next Steps

1. **Finalize Tech Stack**: Choose between Node.js or Go
2. **Setup Development Environment**: MongoDB, Redis, IDE
3. **Initialize Git Repository**: Create monorepo structure
4. **Design API Contracts**: OpenAPI specification
5. **Create Database Migrations**: Schema versioning
6. **Build MVP**: Focus on core GitHub integration
7. **Deploy Alpha**: Internal testing
8. **Gather Feedback**: Iterate on features

---

## Questions to Resolve

1. Should we support monorepos with multiple services?
2. How to handle private dependencies (npm private packages)?
3. Do we need to track historical dependency changes over time?
4. Should the CLI work offline with cached data?
5. How to handle breaking changes detection automatically?
6. Should we support custom version schemes (non-semver)?

---

## References

- [GitHub REST API Documentation](https://docs.github.com/en/rest)
- [Semantic Versioning](https://semver.org/)
- [MongoDB Best Practices](https://www.mongodb.com/docs/manual/administration/production-notes/)
- [OAuth 2.0 Specification](https://oauth.net/2/)

