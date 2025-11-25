import Service from '../models/Service.js';
import User from '../models/User.js';
import Event from '../models/Event.js';
import githubService from '../services/githubService.js';
import syncService from '../services/syncService.js';
import { AppError } from '../middleware/errorHandler.js';
import logger from '../utils/logger.js';

/**
 * OAuth callback handler
 * GET /api/v1/integrations/github/callback
 */
export const githubCallback = async (req, res, next) => {
  try {
    const { code } = req.query;

    if (!code) {
      throw new AppError('Authorization code missing', 400, 'CODE_MISSING');
    }

    // Exchange code for access token
    const accessToken = await githubService.exchangeCodeForToken(code);

    // Get user info from GitHub
    const githubUser = await githubService.getUserInfo(accessToken);

    // Find or create user
    let user = await User.findOne({ authProviderId: githubUser.id.toString() });

    if (!user) {
      user = new User({
        email: githubUser.email || `${githubUser.login}@github.com`,
        username: githubUser.login,
        authProvider: 'github',
        authProviderId: githubUser.id.toString()
      });
    }

    // Store GitHub token
    user.tokens.github = {
      accessToken,
      expiresAt: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000) // 1 year
    };
    user.lastLoginAt = new Date();

    await user.save();

    logger.info(`GitHub OAuth successful for user: ${user.email}`);

    // In production, redirect to frontend with token
    res.json({
      success: true,
      message: 'GitHub authentication successful',
      user: {
        id: user._id,
        email: user.email,
        username: user.username
      }
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Get user's accessible repositories
 * GET /api/v1/integrations/github/repos
 */
export const getUserRepos = async (req, res, next) => {
  try {
    if (!req.user) {
      throw new AppError('Authentication required', 401, 'AUTH_REQUIRED');
    }

    const accessToken = req.user.tokens?.github?.accessToken;

    if (!accessToken) {
      throw new AppError('GitHub not connected', 400, 'GITHUB_NOT_CONNECTED');
    }

    const { page = 1, perPage = 30 } = req.query;

    const repos = await githubService.getUserRepositories(
      accessToken,
      parseInt(page),
      parseInt(perPage)
    );

    res.json({
      success: true,
      data: repos,
      page: parseInt(page),
      perPage: parseInt(perPage)
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Import GitHub repository as service
 * POST /api/v1/integrations/github/import
 * 
 * Supports both authenticated and unauthenticated requests:
 * - Public repositories: No authentication required
 * - Private repositories: Requires authentication via API key and GitHub OAuth
 */
export const importRepository = async (req, res, next) => {
  try {
    const { repoURL, alias, trackingMethod } = req.body;

    if (!repoURL) {
      throw new AppError('Repository URL is required', 400, 'REPO_URL_REQUIRED');
    }

    // Parse repository URL
    const { owner, repo } = githubService.parseRepoURL(repoURL);

    // Get access token if user is authenticated (optional for public repos)
    const accessToken = req.user?.tokens?.github?.accessToken;

    // Fetch repository info (works without token for public repos)
    const repoInfo = await githubService.getRepository(owner, repo, accessToken);

    // Check if service already exists
    const existingService = await Service.findByURL(repoInfo.html_url);
    if (existingService) {
      throw new AppError('Service already exists', 409, 'SERVICE_EXISTS');
    }

    // Create service
    const service = new Service({
      name: repoInfo.name,
      alias: alias || repoInfo.name,
      repository: 'github',
      repoURL: repoInfo.html_url,
      repoOwner: owner,
      repoName: repo,
      visibility: repoInfo.private ? 'private' : 'public',
      trackingMethod: trackingMethod || 'manual',
      userId: req.user?._id,
      metadata: {
        language: repoInfo.language
      }
    });

    await service.save();

    // Trigger initial sync
    try {
      await syncService.syncService(service._id, {
        triggeredBy: req.user ? 'user' : 'system',
        userId: req.user?._id,
        accessToken
      });
    } catch (syncError) {
      logger.warn(`Initial sync failed for ${service.name}:`, syncError);
    }

    // Create event
    await Event.createEvent({
      type: 'service.created',
      severity: 'info',
      serviceId: service._id,
      userId: req.user?._id,
      message: `GitHub repository ${owner}/${repo} imported`,
      details: { repoURL: repoInfo.html_url }
    });

    logger.info(`Imported GitHub repository: ${owner}/${repo}`);

    res.status(201).json({
      success: true,
      data: service
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Trigger manual sync for a service
 * POST /api/v1/services/:id/sync
 * 
 * Supports both authenticated and unauthenticated requests:
 * - Public repositories: No authentication required
 * - Private repositories: Requires authentication
 */
export const triggerSync = async (req, res, next) => {
  try {
    const { id } = req.params;

    const service = await Service.findById(id);

    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    // Verify access for private repositories
    if (service.visibility === 'private' && !req.user) {
      throw new AppError('Authentication required for private repositories', 401, 'AUTH_REQUIRED');
    }

    // Get access token if available (required for private repos)
    const accessToken = req.user?.tokens?.github?.accessToken;

    // Trigger sync
    const result = await syncService.syncService(id, {
      triggeredBy: req.user ? 'user' : 'system',
      userId: req.user?._id,
      syncType: 'manual',
      accessToken
    });

    logger.info(`Manual sync triggered for service: ${service.name}`);

    res.json({
      success: true,
      message: 'Sync completed successfully',
      data: result
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Setup webhook for a repository
 * POST /api/v1/integrations/github/webhook
 */
export const setupWebhook = async (req, res, next) => {
  try {
    const { serviceId, events } = req.body;

    if (!req.user) {
      throw new AppError('Authentication required', 401, 'AUTH_REQUIRED');
    }

    const accessToken = req.user.tokens?.github?.accessToken;

    if (!accessToken) {
      throw new AppError('GitHub not connected', 400, 'GITHUB_NOT_CONNECTED');
    }

    const service = await Service.findById(serviceId);

    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    if (service.repository !== 'github') {
      throw new AppError('Service is not a GitHub repository', 400, 'INVALID_REPOSITORY_TYPE');
    }

    // Generate webhook secret
    const crypto = await import('crypto');
    const secret = crypto.randomBytes(32).toString('hex');

    // Setup webhook on GitHub
    const callbackURL = `${process.env.API_BASE_URL || 'http://localhost:3000'}/api/v1/webhooks/github/${serviceId}`;
    
    const webhook = await githubService.createWebhook(
      service.repoOwner,
      service.repoName,
      callbackURL,
      secret,
      accessToken
    );

    // Store webhook info
    const Webhook = (await import('../models/Webhook.js')).default;
    const webhookDoc = new Webhook({
      serviceId: service._id,
      userId: req.user._id,
      provider: 'github',
      webhookId: webhook.id.toString(),
      events: events || ['push', 'release', 'create'],
      secret,
      active: true
    });

    await webhookDoc.save();

    // Update service tracking method
    service.trackingMethod = 'webhook';
    await service.save();

    logger.info(`Webhook created for ${service.name}`);

    res.json({
      success: true,
      message: 'Webhook setup successful',
      data: {
        webhookId: webhook.id,
        events: webhookDoc.events
      }
    });
  } catch (error) {
    next(error);
  }
};

