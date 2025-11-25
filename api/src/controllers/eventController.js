import Event from '../models/Event.js';
import { AppError } from '../middleware/errorHandler.js';

/**
 * Get events feed
 * GET /api/v1/events
 */
export const getEvents = async (req, res, next) => {
  try {
    const { 
      severity, 
      type, 
      limit = 20, 
      page = 1,
      includeArchived = false 
    } = req.query;

    // Build query
    const query = {};
    if (severity) query.severity = severity;
    if (type) query.type = type;
    if (req.user) query.userId = req.user._id;
    if (!includeArchived) query.isArchived = false;

    // Pagination
    const skip = (page - 1) * limit;

    const events = await Event.find(query)
      .sort({ createdAt: -1 })
      .skip(skip)
      .limit(parseInt(limit))
      .populate('serviceId', 'name alias')
      .populate('userId', 'email username');

    const total = await Event.countDocuments(query);
    const unreadCount = req.user ? await Event.getUnreadCount(req.user._id) : 0;

    res.json({
      success: true,
      data: events,
      unreadCount,
      pagination: {
        total,
        page: parseInt(page),
        limit: parseInt(limit),
        pages: Math.ceil(total / limit)
      }
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Get event by ID
 * GET /api/v1/events/:id
 */
export const getEventById = async (req, res, next) => {
  try {
    const { id } = req.params;

    const event = await Event.findById(id)
      .populate('serviceId', 'name alias repoURL')
      .populate('userId', 'email username');

    if (!event) {
      throw new AppError('Event not found', 404, 'EVENT_NOT_FOUND');
    }

    res.json({
      success: true,
      data: event
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Mark event as read
 * POST /api/v1/events/:id/read
 */
export const markEventAsRead = async (req, res, next) => {
  try {
    const { id } = req.params;

    if (!req.user) {
      throw new AppError('Authentication required', 401, 'AUTH_REQUIRED');
    }

    const event = await Event.findById(id);

    if (!event) {
      throw new AppError('Event not found', 404, 'EVENT_NOT_FOUND');
    }

    event.markAsRead(req.user._id);
    await event.save();

    res.json({
      success: true,
      message: 'Event marked as read'
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Archive event
 * POST /api/v1/events/:id/archive
 */
export const archiveEvent = async (req, res, next) => {
  try {
    const { id } = req.params;

    const event = await Event.findById(id);

    if (!event) {
      throw new AppError('Event not found', 404, 'EVENT_NOT_FOUND');
    }

    event.archive();
    await event.save();

    res.json({
      success: true,
      message: 'Event archived successfully'
    });
  } catch (error) {
    next(error);
  }
};

