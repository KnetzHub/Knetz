import express from 'express';
import serviceRoutes from './serviceRoutes.js';
import healthRoutes from './healthRoutes.js';
import eventRoutes from './eventRoutes.js';
import logRoutes from './logRoutes.js';
import githubRoutes from './githubRoutes.js';
import webhookRoutes from './webhookRoutes.js';

const router = express.Router();

// Mount routes
router.use('/health', healthRoutes);
router.use('/services', serviceRoutes);
router.use('/events', eventRoutes);
router.use('/logs', logRoutes);
router.use('/integrations/github', githubRoutes);
router.use('/webhooks', webhookRoutes);

// API info endpoint
router.get('/', (req, res) => {
  res.json({
    name: 'KnetZ API',
    version: '1.0.0',
    documentation: '/api/v1/docs',
    endpoints: {
      health: '/api/v1/health',
      services: '/api/v1/services',
      events: '/api/v1/events',
      logs: '/api/v1/logs',
      github: '/api/v1/integrations/github',
      webhooks: '/api/v1/webhooks'
    }
  });
});

export default router;

