import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import dotenv from 'dotenv';
import swaggerUi from 'swagger-ui-express';
import { connectDB } from './config/database.js';
import cronService from './services/cronService.js';
import logger from './utils/logger.js';
import { errorHandler } from './middleware/errorHandler.js';
import { rateLimiter } from './middleware/rateLimiter.js';
import routes from './routes/index.js';
import swaggerSpec from './swagger.js';

// Load environment variables
dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;
const API_VERSION = process.env.API_VERSION || 'v1';

// Middleware
app.use(helmet({
  contentSecurityPolicy: false // Disable CSP for Swagger UI
})); // Security headers
app.use(cors({
  origin: process.env.CORS_ORIGIN || '*',
  credentials: true
}));
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(rateLimiter);

// Request logging
app.use((req, res, next) => {
  logger.info(`${req.method} ${req.path}`, {
    ip: req.ip,
    userAgent: req.get('user-agent')
  });
  next();
});

// Swagger API Documentation
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerSpec, {
  customCss: '.swagger-ui .topbar { display: none }',
  customSiteTitle: 'KnetZ API Documentation',
  customfavIcon: '/favicon.ico'
}));

// Swagger JSON
app.get('/api-docs.json', (req, res) => {
  res.setHeader('Content-Type', 'application/json');
  res.send(swaggerSpec);
});

// API Routes
app.use(`/api/${API_VERSION}`, routes);

// Root endpoint
app.get('/', (req, res) => {
  res.json({
    name: 'KnetZ API',
    version: '1.0.0',
    status: 'running',
    documentation: {
      swagger: `http://localhost:${PORT}/api-docs`,
      openapi: `http://localhost:${PORT}/api-docs.json`,
      postman: 'Import KnetZ_API.postman_collection.json'
    },
    endpoints: {
      health: `http://localhost:${PORT}/api/${API_VERSION}/health`,
      services: `http://localhost:${PORT}/api/${API_VERSION}/services`,
      github: `http://localhost:${PORT}/api/${API_VERSION}/integrations/github`
    },
    timestamp: new Date().toISOString()
  });
});

// 404 Handler
app.use((req, res) => {
  res.status(404).json({
    error: {
      code: 'NOT_FOUND',
      message: `Route ${req.method} ${req.path} not found`,
      timestamp: new Date().toISOString()
    }
  });
});

// Error Handler (must be last)
app.use(errorHandler);

// Start Server
const startServer = async () => {
  try {
    // Connect to MongoDB
    await connectDB();
    
    // Start cron jobs (only in production or if explicitly enabled)
    if (process.env.NODE_ENV === 'production' || process.env.ENABLE_CRON === 'true') {
      cronService.start();
      logger.info('⏰ Cron jobs started');
    } else {
      logger.info('⏰ Cron jobs disabled (set ENABLE_CRON=true to enable)');
    }
    
    // Start listening
    app.listen(PORT, () => {
      logger.info(`🚀 KnetZ API server running on port ${PORT}`);
      logger.info(`📝 Environment: ${process.env.NODE_ENV || 'development'}`);
      logger.info(`🔗 API Base URL: http://localhost:${PORT}/api/${API_VERSION}`);
      logger.info(`📚 API Documentation: http://localhost:${PORT}/api-docs`);
      logger.info(`📄 OpenAPI Spec: http://localhost:${PORT}/api-docs.json`);
    });
  } catch (error) {
    logger.error('Failed to start server:', error);
    process.exit(1);
  }
};

// Handle unhandled rejections
process.on('unhandledRejection', (err) => {
  logger.error('Unhandled Rejection:', err);
  process.exit(1);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.info('SIGTERM received, shutting down gracefully');
  cronService.stop();
  process.exit(0);
});

process.on('SIGINT', () => {
  logger.info('SIGINT received, shutting down gracefully');
  cronService.stop();
  process.exit(0);
});

startServer();

export default app;

