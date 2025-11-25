import express from 'express';
import {
  githubCallback,
  getUserRepos,
  importRepository,
  triggerSync,
  setupWebhook
} from '../controllers/githubController.js';
import { optionalAuth, authenticateAPIKey } from '../middleware/auth.js';

const router = express.Router();

// OAuth
router.get('/callback', githubCallback);

// Repository management
router.get('/repos', authenticateAPIKey, getUserRepos);
router.post('/import', optionalAuth, importRepository);

// Webhook management
router.post('/webhook', authenticateAPIKey, setupWebhook);

export default router;

