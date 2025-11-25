import Joi from 'joi';
import { AppError } from './errorHandler.js';

/**
 * Validate service creation/update
 */
export const validateService = (req, res, next) => {
  const schema = Joi.object({
    name: Joi.string().required().trim().min(1).max(100),
    alias: Joi.string().trim().max(50),
    repository: Joi.string().required().valid('github', 'gitlab', 'docker', 'quay'),
    repoURL: Joi.string().required().uri().max(500),
    repoOwner: Joi.string().required().trim().max(100),
    repoName: Joi.string().required().trim().max(100),
    visibility: Joi.string().valid('public', 'private').default('public'),
    trackingMethod: Joi.string().valid('cron', 'webhook', 'manual').default('manual'),
    metadata: Joi.object({
      language: Joi.string().max(50),
      framework: Joi.string().max(50),
      packageManager: Joi.string().max(50),
      manifestFile: Joi.string().max(100),
      manifestPath: Joi.string().max(500)
    })
  });

  const { error, value } = schema.validate(req.body, { abortEarly: false });

  if (error) {
    const messages = error.details.map(detail => detail.message).join(', ');
    throw new AppError(messages, 400, 'VALIDATION_ERROR');
  }

  req.body = value;
  next();
};

/**
 * Validate version data
 */
export const validateVersion = (req, res, next) => {
  const schema = Joi.object({
    version: Joi.string().required().pattern(/^\d+\.\d+\.\d+/),
    type: Joi.string().valid('release', 'tag', 'commit', 'image').default('release'),
    releasedAt: Joi.date().default(Date.now),
    changelog: Joi.string().max(5000),
    artifacts: Joi.array().items(Joi.string()),
    isMajor: Joi.boolean().default(false),
    isPrerelease: Joi.boolean().default(false)
  });

  const { error, value } = schema.validate(req.body, { abortEarly: false });

  if (error) {
    const messages = error.details.map(detail => detail.message).join(', ');
    throw new AppError(messages, 400, 'VALIDATION_ERROR');
  }

  req.body = value;
  next();
};

/**
 * Validate pagination parameters
 */
export const validatePagination = (req, res, next) => {
  const schema = Joi.object({
    page: Joi.number().integer().min(1).default(1),
    limit: Joi.number().integer().min(1).max(100).default(20)
  });

  const { error, value } = schema.validate(req.query, { abortEarly: false });

  if (error) {
    const messages = error.details.map(detail => detail.message).join(', ');
    throw new AppError(messages, 400, 'VALIDATION_ERROR');
  }

  req.query = { ...req.query, ...value };
  next();
};

