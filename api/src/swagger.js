import swaggerJsdoc from 'swagger-jsdoc';

const options = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'KnetZ API',
      version: '1.0.0',
      description: `
# KnetZ - Dependency Management Platform API

A comprehensive API for tracking service versions, dependencies, and automating version management across multiple sources.

## Features

- 🌍 **Public Repository Support** - Import and track public GitHub repositories without authentication
- 🔐 **Private Repository Support** - Full OAuth support for private repositories
- 📦 **Multi-Language Support** - Parse manifests for Node.js, Go, Python, Rust, and more
- 🔄 **Automated Syncing** - Cron-based and webhook-triggered synchronization
- 🔗 **Dependency Tracking** - Automatic dependency extraction and linking
- 📊 **Version Management** - Track releases, tags, and version history

## Authentication

### Public Repositories
For public GitHub repositories, **no authentication is required**. You can import and sync public repositories without any setup.

### Private Repositories
For private repositories, authentication is required:
- **API Key**: Include \`x-api-key\` header
- **GitHub OAuth**: Authenticate via OAuth flow

## Rate Limits

- **Unauthenticated**: 60 requests/hour per IP (GitHub public API limit)
- **Authenticated**: 5,000 requests/hour (GitHub authenticated limit)
      `,
      contact: {
        name: 'KnetZ Team',
        url: 'https://github.com/yourusername/knetz'
      },
      license: {
        name: 'MIT',
        url: 'https://opensource.org/licenses/MIT'
      }
    },
    servers: [
      {
        url: 'http://localhost:3000/api/v1',
        description: 'Development server'
      },
      {
        url: 'https://api.knetz.io/v1',
        description: 'Production server'
      }
    ],
    tags: [
      {
        name: 'Health',
        description: 'Health check and platform statistics'
      },
      {
        name: 'Services',
        description: 'Service management operations'
      },
      {
        name: 'GitHub',
        description: 'GitHub integration endpoints'
      },
      {
        name: 'Webhooks',
        description: 'Webhook handlers'
      },
      {
        name: 'Events',
        description: 'Event management and notifications'
      },
      {
        name: 'Logs',
        description: 'Sync logs and operation history'
      }
    ],
    components: {
      securitySchemes: {
        ApiKeyAuth: {
          type: 'apiKey',
          in: 'header',
          name: 'x-api-key',
          description: 'API key for authentication (required for private repositories)'
        }
      },
      schemas: {
        Service: {
          type: 'object',
          properties: {
            _id: {
              type: 'string',
              example: '507f1f77bcf86cd799439011'
            },
            name: {
              type: 'string',
              example: 'express'
            },
            alias: {
              type: 'string',
              example: 'express-framework'
            },
            repository: {
              type: 'string',
              enum: ['github', 'gitlab', 'docker', 'quay'],
              example: 'github'
            },
            repoURL: {
              type: 'string',
              example: 'https://github.com/expressjs/express'
            },
            repoOwner: {
              type: 'string',
              example: 'expressjs'
            },
            repoName: {
              type: 'string',
              example: 'express'
            },
            visibility: {
              type: 'string',
              enum: ['public', 'private'],
              example: 'public'
            },
            trackingMethod: {
              type: 'string',
              enum: ['manual', 'cron', 'webhook'],
              example: 'manual'
            },
            lastSyncedAt: {
              type: 'string',
              format: 'date-time'
            },
            createdAt: {
              type: 'string',
              format: 'date-time'
            },
            updatedAt: {
              type: 'string',
              format: 'date-time'
            },
            versions: {
              type: 'array',
              items: {
                $ref: '#/components/schemas/Version'
              }
            },
            dependencies: {
              type: 'array',
              items: {
                $ref: '#/components/schemas/Dependency'
              }
            },
            metadata: {
              type: 'object',
              properties: {
                language: {
                  type: 'string',
                  example: 'javascript'
                },
                packageManager: {
                  type: 'string',
                  example: 'npm'
                },
                manifestFile: {
                  type: 'string',
                  example: 'package.json'
                }
              }
            }
          }
        },
        Version: {
          type: 'object',
          properties: {
            version: {
              type: 'string',
              example: '4.18.2'
            },
            type: {
              type: 'string',
              enum: ['release', 'tag', 'commit', 'image'],
              example: 'release'
            },
            releasedAt: {
              type: 'string',
              format: 'date-time'
            },
            changelog: {
              type: 'string',
              example: 'Bug fixes and improvements'
            },
            artifacts: {
              type: 'array',
              items: {
                type: 'string'
              }
            },
            isMajor: {
              type: 'boolean',
              example: false
            },
            isPrerelease: {
              type: 'boolean',
              example: false
            }
          }
        },
        Dependency: {
          type: 'object',
          properties: {
            name: {
              type: 'string',
              example: 'body-parser'
            },
            currentVersion: {
              type: 'string',
              example: '1.20.1'
            },
            requiredVersion: {
              type: 'string',
              example: '1.20.1'
            },
            serviceId: {
              type: 'string',
              example: '507f1f77bcf86cd799439011'
            },
            status: {
              type: 'string',
              enum: ['up-to-date', 'outdated', 'external', 'unknown', 'ahead'],
              example: 'up-to-date'
            }
          }
        },
        Event: {
          type: 'object',
          properties: {
            _id: {
              type: 'string',
              example: '507f1f77bcf86cd799439011'
            },
            type: {
              type: 'string',
              example: 'version.detected'
            },
            severity: {
              type: 'string',
              enum: ['info', 'warning', 'error', 'critical'],
              example: 'info'
            },
            serviceId: {
              type: 'string'
            },
            userId: {
              type: 'string'
            },
            message: {
              type: 'string',
              example: 'New version 1.2.3 detected'
            },
            details: {
              type: 'object'
            },
            createdAt: {
              type: 'string',
              format: 'date-time'
            },
            isArchived: {
              type: 'boolean',
              example: false
            }
          }
        },
        SyncLog: {
          type: 'object',
          properties: {
            _id: {
              type: 'string'
            },
            serviceId: {
              type: 'string'
            },
            serviceName: {
              type: 'string'
            },
            syncType: {
              type: 'string',
              enum: ['manual', 'cron', 'webhook', 'cli']
            },
            status: {
              type: 'string',
              enum: ['success', 'failed', 'partial']
            },
            triggeredBy: {
              type: 'string',
              enum: ['system', 'user', 'webhook']
            },
            versionsAdded: {
              type: 'number'
            },
            versionsUpdated: {
              type: 'number'
            },
            errors: {
              type: 'array',
              items: {
                type: 'object'
              }
            },
            duration: {
              type: 'number',
              description: 'Duration in milliseconds'
            },
            startedAt: {
              type: 'string',
              format: 'date-time'
            },
            completedAt: {
              type: 'string',
              format: 'date-time'
            }
          }
        },
        Error: {
          type: 'object',
          properties: {
            error: {
              type: 'object',
              properties: {
                code: {
                  type: 'string',
                  example: 'SERVICE_NOT_FOUND'
                },
                message: {
                  type: 'string',
                  example: 'Service not found'
                },
                details: {
                  type: 'object'
                },
                timestamp: {
                  type: 'string',
                  format: 'date-time'
                }
              }
            }
          }
        }
      },
      responses: {
        UnauthorizedError: {
          description: 'Authentication required',
          content: {
            'application/json': {
              schema: {
                $ref: '#/components/schemas/Error'
              }
            }
          }
        },
        NotFoundError: {
          description: 'Resource not found',
          content: {
            'application/json': {
              schema: {
                $ref: '#/components/schemas/Error'
              }
            }
          }
        },
        ValidationError: {
          description: 'Validation error',
          content: {
            'application/json': {
              schema: {
                $ref: '#/components/schemas/Error'
              }
            }
          }
        }
      }
    }
  },
  apis: ['./src/routes/*.js', './src/controllers/*.js']
};

const swaggerSpec = swaggerJsdoc(options);

export default swaggerSpec;

