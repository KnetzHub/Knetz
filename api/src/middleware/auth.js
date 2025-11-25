import { AppError } from './errorHandler.js';
import User from '../models/User.js';
import logger from '../utils/logger.js';
import crypto from 'crypto';

/**
 * Authenticate request using API key
 */
export const authenticateAPIKey = async (req, res, next) => {
  try {
    const apiKey = req.headers['x-api-key'] || req.query.apiKey;

    if (!apiKey) {
      throw new AppError('API key is required', 401, 'API_KEY_MISSING');
    }

    // Hash the provided API key
    const hashedKey = hashAPIKey(apiKey);

    // Find user with matching API key
    const user = await User.findOne({
      'apiKeys.key': hashedKey
    });

    if (!user) {
      throw new AppError('Invalid API key', 401, 'INVALID_API_KEY');
    }

    // Update last used timestamp
    const apiKeyDoc = user.apiKeys.find(k => k.key === hashedKey);
    if (apiKeyDoc) {
      // Check if key is expired
      if (apiKeyDoc.expiresAt && new Date() > apiKeyDoc.expiresAt) {
        throw new AppError('API key has expired', 401, 'API_KEY_EXPIRED');
      }
      
      apiKeyDoc.lastUsed = new Date();
      await user.save();
    }

    // Attach user to request
    req.user = user;
    next();
  } catch (error) {
    next(error);
  }
};

/**
 * Optional authentication - doesn't fail if no auth provided
 */
export const optionalAuth = async (req, res, next) => {
  try {
    const apiKey = req.headers['x-api-key'] || req.query.apiKey;

    if (apiKey) {
      const hashedKey = hashAPIKey(apiKey);
      const user = await User.findOne({ 'apiKeys.key': hashedKey });
      
      if (user) {
        req.user = user;
      }
    }
    
    next();
  } catch (error) {
    // Don't fail on optional auth errors
    logger.warn('Optional auth failed:', error.message);
    next();
  }
};

/**
 * Hash API key using SHA256
 */
export const hashAPIKey = (apiKey) => {
  return crypto
    .createHash('sha256')
    .update(apiKey + (process.env.API_KEY_SECRET || 'default-secret'))
    .digest('hex');
};

/**
 * Generate a new API key
 */
export const generateAPIKey = () => {
  return 'knetz_' + crypto.randomBytes(32).toString('hex');
};

