import crypto from 'crypto';
import Webhook from '../models/Webhook.js';
import Service from '../models/Service.js';
import Event from '../models/Event.js';
import syncService from '../services/syncService.js';
import logger from '../utils/logger.js';
import { AppError } from '../middleware/errorHandler.js';

/**
 * Verify GitHub webhook signature
 */
const verifyGitHubSignature = (payload, signature, secret) => {
  if (!signature) return false;

  const hmac = crypto.createHmac('sha256', secret);
  const digest = 'sha256=' + hmac.update(JSON.stringify(payload)).digest('hex');
  
  return crypto.timingSafeEqual(
    Buffer.from(signature),
    Buffer.from(digest)
  );
};

/**
 * Handle GitHub webhook
 * POST /api/v1/webhooks/github/:serviceId
 */
export const handleGitHubWebhook = async (req, res, next) => {
  try {
    const { serviceId } = req.params;
    const signature = req.headers['x-hub-signature-256'];
    const event = req.headers['x-github-event'];
    const payload = req.body;

    logger.info(`Received GitHub webhook: ${event} for service ${serviceId}`);

    // Find webhook configuration
    const service = await Service.findById(serviceId);
    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    const webhook = await Webhook.findOne({ 
      serviceId, 
      provider: 'github',
      active: true 
    });

    if (!webhook) {
      throw new AppError('Webhook not found or inactive', 404, 'WEBHOOK_NOT_FOUND');
    }

    // Verify signature
    if (!verifyGitHubSignature(payload, signature, webhook.secret)) {
      logger.warn(`Invalid webhook signature for service ${serviceId}`);
      throw new AppError('Invalid signature', 401, 'INVALID_SIGNATURE');
    }

    // Update last triggered
    webhook.lastTriggered = new Date();
    await webhook.save();

    // Handle different event types
    let shouldSync = false;
    let eventMessage = '';

    switch (event) {
      case 'release':
        if (payload.action === 'published') {
          shouldSync = true;
          eventMessage = `New release published: ${payload.release.tag_name}`;
        }
        break;

      case 'create':
        if (payload.ref_type === 'tag') {
          shouldSync = true;
          eventMessage = `New tag created: ${payload.ref}`;
        }
        break;

      case 'push':
        // Only sync on main/master branch pushes
        if (payload.ref === 'refs/heads/main' || payload.ref === 'refs/heads/master') {
          shouldSync = true;
          eventMessage = `Push to ${payload.ref}`;
        }
        break;

      default:
        logger.info(`Ignoring webhook event: ${event}`);
    }

    // Trigger sync if needed
    if (shouldSync) {
      // Create event
      await Event.createEvent({
        type: 'webhook.triggered',
        severity: 'info',
        serviceId: service._id,
        message: eventMessage,
        details: { event, payload: payload }
      });

      // Trigger sync in background
      syncService.syncService(serviceId, {
        triggeredBy: 'webhook',
        syncType: 'webhook'
      }).catch(error => {
        logger.error(`Webhook sync failed for ${serviceId}:`, error);
      });
    }

    // Always respond quickly to GitHub
    res.status(200).json({ success: true, message: 'Webhook received' });

  } catch (error) {
    logger.error('Webhook handler error:', error);
    next(error);
  }
};

/**
 * Handle GitLab webhook (placeholder for Phase 2)
 * POST /api/v1/webhooks/gitlab/:serviceId
 */
export const handleGitLabWebhook = async (req, res, next) => {
  res.status(501).json({
    error: {
      code: 'NOT_IMPLEMENTED',
      message: 'GitLab webhooks not yet implemented'
    }
  });
};

