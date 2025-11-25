# KnetZ API

Backend API for KnetZ - A dependency management platform that tracks service versions across multiple sources.

## Features

✅ **GitHub Integration**
- 🌍 Public repository support (no authentication required)
- 🔐 Private repository support (OAuth authentication)
- 📦 Automated version tracking (releases & tags)
- 🎣 Webhook support for real-time updates

✅ **Manifest Parsing**
- Node.js (package.json)
- Go (go.mod)
- Python (requirements.txt, Pipfile)
- Rust (Cargo.toml)
- More coming soon...

✅ **Automated Syncing**
- ⏰ Cron-based scheduled syncing
- 🔄 Manual sync triggers
- 🎯 Webhook-based real-time syncing

✅ **Dependency Tracking**
- 🔗 Automatic dependency extraction from manifests
- 📊 Link dependencies to tracked services
- ⚠️ Version comparison and outdated detection

✅ **API Features**
- RESTful API design
- Optional authentication for public repos
- Rate limiting and security headers
- Comprehensive error handling
- Structured logging

## Getting Started

### Prerequisites
- Node.js >= 18.x
- MongoDB >= 6.x
- npm or yarn

### Installation

1. Install dependencies:
```bash
npm install
```

2. Setup environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

Required environment variables:
```bash
# Server
NODE_ENV=development
PORT=3000
API_BASE_URL=http://localhost:3000

# MongoDB
MONGO_URI=mongodb://localhost:27017/knetz

# Authentication
API_KEY_SECRET=your-super-secret-key
JWT_SECRET=your-jwt-secret

# GitHub OAuth (optional for public repos only)
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# Cron Jobs (set to true to enable automated syncing)
ENABLE_CRON=false
```

3. Start MongoDB (if running locally):
```bash
mongod
```

4. Run the development server:
```bash
npm run dev
```

The API will be available at `http://localhost:3000`

## API Documentation

### Base URL
```
http://localhost:3000/api/v1
```

### Quick Usage Examples

#### 1. Import a Public GitHub Repository (No Auth Required)

```bash
curl -X POST http://localhost:3000/api/v1/integrations/github/import \
  -H "Content-Type: application/json" \
  -d '{
    "repoURL": "https://github.com/expressjs/express",
    "alias": "express",
    "trackingMethod": "manual"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "_id": "...",
    "name": "express",
    "alias": "express",
    "repository": "github",
    "visibility": "public",
    "versions": [...],
    "dependencies": [...]
  }
}
```

#### 2. Get Service Information

```bash
# List all services
curl http://localhost:3000/api/v1/services

# Get specific service
curl http://localhost:3000/api/v1/services/{serviceId}

# Get version history
curl http://localhost:3000/api/v1/services/{serviceId}/versions

# Get latest version
curl http://localhost:3000/api/v1/services/{serviceId}/latest

# Get dependencies
curl http://localhost:3000/api/v1/services/{serviceId}/dependencies
```

#### 3. Trigger Manual Sync

```bash
curl -X POST http://localhost:3000/api/v1/services/{serviceId}/sync
```

#### 4. Health Check

```bash
curl http://localhost:3000/api/v1/health
```

### Private Repository Access

For private repositories, you need to authenticate:

1. **Authenticate with GitHub OAuth**
   - Visit: `http://localhost:3000/api/v1/integrations/github/callback?code={github_oauth_code}`
   
2. **Use API Key for requests**

```bash
curl -X POST http://localhost:3000/api/v1/integrations/github/import \
  -H "Content-Type: application/json" \
  -H "x-api-key: knetz_your_api_key" \
  -d '{
    "repoURL": "https://github.com/yourorg/private-repo",
    "alias": "my-private-service"
  }'
```

### Available Endpoints

For complete API documentation, see [API.md](./API.md)

## Project Structure

```
api/
├── src/
│   ├── config/        # Configuration files
│   ├── models/        # MongoDB models
│   ├── routes/        # API routes
│   ├── controllers/   # Route controllers
│   ├── middleware/    # Express middleware
│   ├── services/      # Business logic
│   ├── utils/         # Utility functions
│   └── server.js      # Entry point
├── tests/             # Test files
├── .env.example       # Environment template
└── package.json
```

## Configuration

### Enabling Automated Syncing

To enable cron-based automated syncing, set `ENABLE_CRON=true` in your `.env` file:

```bash
ENABLE_CRON=true
```

The cron service will automatically sync services with `trackingMethod: "cron"` at these intervals:
- **Medium Priority**: Every 4 hours
- **High Priority**: Every 30 minutes (coming soon)
- **Low Priority**: Daily at midnight (coming soon)

### Webhook Setup

To enable real-time updates for private repositories:

1. Setup webhook via API (requires authentication):
```bash
curl -X POST http://localhost:3000/api/v1/integrations/github/webhook \
  -H "x-api-key: your_api_key" \
  -H "Content-Type: application/json" \
  -d '{
    "serviceId": "your_service_id",
    "events": ["push", "release", "create"]
  }'
```

The webhook will automatically trigger syncs when:
- New releases are published
- New tags are created
- Commits are pushed to main/master branch

## Scripts

- `npm start` - Start production server
- `npm run dev` - Start development server with hot reload
- `npm test` - Run tests
- `npm run lint` - Lint code

## Supported Manifest Files

| Language | Manifest File | Status |
|----------|--------------|--------|
| Node.js | package.json | ✅ Supported |
| Go | go.mod | ✅ Supported |
| Python | requirements.txt | ✅ Supported |
| Python | Pipfile | ✅ Supported |
| Rust | Cargo.toml | ✅ Supported |
| Java | pom.xml | 🔄 Coming Soon |
| Java/Kotlin | build.gradle | 🔄 Coming Soon |
| PHP | composer.json | 🔄 Coming Soon |

## License

MIT

