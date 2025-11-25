import express from 'express';
import { handleGitHubWebhook, handleGitLabWebhook } from '../controllers/webhookController.js';

const router = express.Router();

// Webhook handlers
router.post('/github/:serviceId', handleGitHubWebhook);
router.post('/gitlab/:serviceId', handleGitLabWebhook);

export default router;

