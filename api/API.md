# KnetZ API Documentation

Version: 1.0.0  
Base URL: `http://localhost:3000/api/v1`

---

## Authentication

### 🌍 Public Repository Access (No Auth Required)

For **public GitHub repositories**, authentication is **optional**. You can import and sync public repositories without any authentication:

```bash
# Import public repository - no auth needed
curl -X POST http://localhost:3000/api/v1/integrations/github/import \
  -H "Content-Type: application/json" \
  -d '{
    "repoURL": "https://github.com/expressjs/express",
    "alias": "express"
  }'

# Sync public repository - no auth needed
curl -X POST http://localhost:3000/api/v1/services/{serviceId}/sync
```

### 🔐 Private Repository Access (Auth Required)

For **private GitHub repositories**, you need to authenticate:

#### API Key Authentication
Include your API key in the request header:
```
X-API-Key: your-api-key-here
```

Or as a query parameter:
```
?apiKey=your-api-key-here
```

#### GitHub OAuth
1. Redirect user to GitHub OAuth:
   ```
   https://github.com/login/oauth/authorize?client_id={GITHUB_CLIENT_ID}&scope=repo
   ```

2. Handle callback:
   ```
   GET /api/v1/integrations/github/callback?code={code}
   ```

### Rate Limits

- **Unauthenticated requests**: 60 requests/hour per IP (GitHub public API limit)
- **Authenticated requests**: 5,000 requests/hour (GitHub authenticated limit)

---

## Endpoints

### Health & Status

#### Health Check
```http
GET /api/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-11-25T10:00:00Z",
  "uptime": 3600,
  "environment": "development",
  "version": "1.0.0",
  "services": {
    "database": {
      "status": "connected",
      "name": "knetz"
    },
    "api": {
      "status": "running"
    }
  }
}
```

#### Platform Statistics
```http
GET /api/v1/health/stats
```

**Response:**
```json
{
  "success": true,
  "data": {
    "overview": {
      "totalServices": 50,
      "totalUsers": 10,
      "totalSyncs": 500,
      "recentEvents": 25
    },
    "sync": {
      "last24Hours": {
        "success": 45,
        "failed": 5,
        "total": 50,
        "successRate": "90.00%"
      }
    }
  }
}
```

---

### Services

#### Create Service
```http
POST /api/v1/services
Content-Type: application/json
X-API-Key: {apiKey} (optional)
```

**Request Body:**
```json
{
  "name": "express",
  "alias": "express-js",
  "repository": "github",
  "repoURL": "https://github.com/expressjs/express",
  "repoOwner": "expressjs",
  "repoName": "express",
  "visibility": "public",
  "trackingMethod": "manual"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "name": "express",
    "alias": "express-js",
    "repository": "github",
    "repoURL": "https://github.com/expressjs/express",
    "repoOwner": "expressjs",
    "repoName": "express",
    "visibility": "public",
    "trackingMethod": "manual",
    "versions": [],
    "dependencies": [],
    "createdAt": "2025-11-25T10:00:00Z",
    "updatedAt": "2025-11-25T10:00:00Z"
  }
}
```

#### List Services
```http
GET /api/v1/services?page=1&limit=20&repository=github&search=express
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)
- `repository` (optional): Filter by repository type
- `trackingMethod` (optional): Filter by tracking method
- `search` (optional): Search by name, alias, or URL

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "name": "express",
      "alias": "express-js",
      "repository": "github",
      "repoURL": "https://github.com/expressjs/express",
      "lastSyncedAt": "2025-11-25T09:00:00Z",
      "createdAt": "2025-11-25T08:00:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "page": 1,
    "limit": 20,
    "pages": 1
  }
}
```

#### Get Service by ID
```http
GET /api/v1/services/:id
```

#### Update Service
```http
PUT /api/v1/services/:id
Content-Type: application/json
```

**Request Body:**
```json
{
  "alias": "new-alias",
  "trackingMethod": "cron"
}
```

#### Delete Service
```http
DELETE /api/v1/services/:id
```

#### Trigger Manual Sync
```http
POST /api/v1/services/:id/sync
X-API-Key: {apiKey} (optional)
```

**Response:**
```json
{
  "success": true,
  "message": "Sync completed successfully",
  "data": {
    "success": true,
    "versionsAdded": 5,
    "versionsUpdated": 2,
    "service": { ... }
  }
}
```

---

### Versions

#### Get Service Versions
```http
GET /api/v1/services/:id/versions?limit=50&type=release
```

**Query Parameters:**
- `limit` (optional): Number of versions to return (default: 50)
- `type` (optional): Filter by type (release, tag, commit, image)

**Response:**
```json
{
  "success": true,
  "data": {
    "serviceId": "507f1f77bcf86cd799439011",
    "serviceName": "express",
    "totalVersions": 100,
    "versions": [
      {
        "version": "4.18.2",
        "type": "release",
        "releasedAt": "2023-10-15T12:00:00Z",
        "changelog": "Bug fixes and improvements",
        "artifacts": [],
        "isMajor": false,
        "isPrerelease": false
      }
    ]
  }
}
```

#### Get Latest Version
```http
GET /api/v1/services/:id/latest
```

---

### Dependencies

#### Get Service Dependencies
```http
GET /api/v1/services/:id/dependencies
```

**Response:**
```json
{
  "success": true,
  "data": {
    "serviceId": "507f1f77bcf86cd799439011",
    "serviceName": "express",
    "currentVersion": "4.18.2",
    "dependencies": [
      {
        "name": "body-parser",
        "currentVersion": "1.20.1",
        "requiredVersion": "1.20.1",
        "status": "up-to-date"
      }
    ],
    "metadata": {
      "language": "javascript",
      "packageManager": "npm",
      "manifestFile": "package.json"
    }
  }
}
```

---

### GitHub Integration

#### OAuth Callback
```http
GET /api/v1/integrations/github/callback?code={code}
```

#### Get User Repositories
```http
GET /api/v1/integrations/github/repos?page=1&perPage=30
X-API-Key: {apiKey}
```

#### Import Repository
```http
POST /api/v1/integrations/github/import
Content-Type: application/json
X-API-Key: {apiKey} (optional for public repos, required for private)
```

**Request Body:**
```json
{
  "repoURL": "https://github.com/expressjs/express",
  "alias": "express-js",
  "trackingMethod": "manual"
}
```

**Note:** 
- ✅ **Public repositories**: No authentication required
- 🔐 **Private repositories**: Requires `X-API-Key` header with valid API key
- The API automatically detects repository visibility

#### Setup Webhook
```http
POST /api/v1/integrations/github/webhook
Content-Type: application/json
X-API-Key: {apiKey}
```

**Request Body:**
```json
{
  "serviceId": "507f1f77bcf86cd799439011",
  "events": ["push", "release", "create"]
}
```

---

### Events

#### Get Events Feed
```http
GET /api/v1/events?page=1&limit=20&severity=info&type=version.detected
```

**Query Parameters:**
- `page`, `limit`: Pagination
- `severity`: Filter by severity (info, warning, error, critical)
- `type`: Filter by event type
- `includeArchived`: Include archived events (default: false)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "type": "version.detected",
      "severity": "info",
      "serviceId": { ... },
      "message": "5 new version(s) detected for express",
      "details": { "versionsAdded": 5 },
      "readBy": [],
      "isArchived": false,
      "createdAt": "2025-11-25T10:00:00Z"
    }
  ],
  "unreadCount": 10,
  "pagination": { ... }
}
```

#### Mark Event as Read
```http
POST /api/v1/events/:id/read
X-API-Key: {apiKey}
```

#### Archive Event
```http
POST /api/v1/events/:id/archive
```

---

### Logs

#### Get Sync Logs
```http
GET /api/v1/logs?serviceId={id}&status=success&page=1&limit=50
```

**Query Parameters:**
- `serviceId`: Filter by service
- `status`: Filter by status (success, failed, partial)
- `syncType`: Filter by sync type (cron, webhook, manual, cli)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "_id": "507f1f77bcf86cd799439011",
      "serviceId": { ... },
      "serviceName": "express",
      "syncType": "manual",
      "status": "success",
      "triggeredBy": "user",
      "versionsAdded": 5,
      "versionsUpdated": 2,
      "errors": [],
      "duration": 2500,
      "startedAt": "2025-11-25T10:00:00Z",
      "completedAt": "2025-11-25T10:00:02Z"
    }
  ],
  "pagination": { ... }
}
```

#### Get Log by ID
```http
GET /api/v1/logs/:id
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "SERVICE_NOT_FOUND",
    "message": "Service with ID 123 not found",
    "timestamp": "2025-11-25T10:00:00Z",
    "requestId": "req_abc123"
  }
}
```

### Common Error Codes
- `AUTH_FAILED` (401): Authentication failed
- `INVALID_INPUT` (400): Validation error
- `SERVICE_NOT_FOUND` (404): Service doesn't exist
- `RATE_LIMIT_EXCEEDED` (429): Too many requests
- `SYNC_FAILED` (500): Sync operation failed
- `GITHUB_API_ERROR` (500): GitHub API error

---

## Rate Limiting

- Default: 100 requests per 15 minutes per IP
- Authenticated: May have higher limits
- Headers included in response:
  - `RateLimit-Limit`: Request limit
  - `RateLimit-Remaining`: Remaining requests
  - `RateLimit-Reset`: Time when limit resets

---

## Webhooks

### GitHub Webhook Endpoint
```
POST /api/v1/webhooks/github/:serviceId
```

**Headers:**
- `X-GitHub-Event`: Event type
- `X-Hub-Signature-256`: HMAC signature

**Supported Events:**
- `release`: New release published
- `push`: Push to main/master branch
- `create`: New tag created

The webhook automatically triggers a sync when relevant events are received.

---

## Examples

### Complete Flow: Import and Sync a Repository

```bash
# 1. Import repository
curl -X POST http://localhost:3000/api/v1/integrations/github/import \
  -H "Content-Type: application/json" \
  -d '{
    "repoURL": "https://github.com/expressjs/express",
    "alias": "express-js"
  }'

# Response: { "success": true, "data": { "_id": "SERVICE_ID", ... } }

# 2. Get service details
curl http://localhost:3000/api/v1/services/SERVICE_ID

# 3. Trigger manual sync
curl -X POST http://localhost:3000/api/v1/services/SERVICE_ID/sync

# 4. Get versions
curl http://localhost:3000/api/v1/services/SERVICE_ID/versions

# 5. Get dependencies
curl http://localhost:3000/api/v1/services/SERVICE_ID/dependencies

# 6. Check sync logs
curl http://localhost:3000/api/v1/logs?serviceId=SERVICE_ID
```

