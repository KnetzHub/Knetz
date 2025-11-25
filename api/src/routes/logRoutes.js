import express from 'express';
import {
  getLogs,
  getLogById,
  getServiceLogs
} from '../controllers/logController.js';
import { optionalAuth } from '../middleware/auth.js';

const router = express.Router();

router.get('/', optionalAuth, getLogs);
router.get('/:id', getLogById);

export default router;

