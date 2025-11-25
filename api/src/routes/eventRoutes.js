import express from 'express';
import {
  getEvents,
  getEventById,
  markEventAsRead,
  archiveEvent
} from '../controllers/eventController.js';
import { optionalAuth } from '../middleware/auth.js';

const router = express.Router();

router.get('/', optionalAuth, getEvents);
router.get('/:id', getEventById);
router.post('/:id/read', optionalAuth, markEventAsRead);
router.post('/:id/archive', optionalAuth, archiveEvent);

export default router;

