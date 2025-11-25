import express from 'express';
import { healthCheck, getStats } from '../controllers/healthController.js';

const router = express.Router();

router.get('/', healthCheck);
router.get('/stats', getStats);

export default router;

