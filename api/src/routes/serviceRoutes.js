import express from 'express';
import {
  createService,
  getServices,
  getServiceById,
  updateService,
  deleteService,
  getServiceVersions,
  getLatestVersion,
  getServiceDependencies
} from '../controllers/serviceController.js';
import { triggerSync } from '../controllers/githubController.js';
import { optionalAuth } from '../middleware/auth.js';
import { validateService } from '../middleware/validators.js';

const router = express.Router();

// Service CRUD
router.post('/', optionalAuth, validateService, createService);
router.get('/', optionalAuth, getServices);
router.get('/:id', getServiceById);
router.put('/:id', optionalAuth, updateService);
router.delete('/:id', optionalAuth, deleteService);

// Version endpoints
router.get('/:id/versions', getServiceVersions);
router.get('/:id/latest', getLatestVersion);

// Dependency endpoints
router.get('/:id/dependencies', getServiceDependencies);

// Sync endpoint
router.post('/:id/sync', optionalAuth, triggerSync);

export default router;

