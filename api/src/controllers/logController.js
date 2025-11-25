import SyncLog from '../models/SyncLog.js';
import { AppError } from '../middleware/errorHandler.js';

/**
 * Get sync logs
 * GET /api/v1/logs
 */
export const getLogs = async (req, res, next) => {
  try {
    const {
      serviceId,
      status,
      syncType,
      limit = 50,
      page = 1
    } = req.query;

    // Build query
    const query = {};
    if (serviceId) query.serviceId = serviceId;
    if (status) query.status = status;
    if (syncType) query.syncType = syncType;
    if (req.user) query.userId = req.user._id;

    // Pagination
    const skip = (page - 1) * limit;

    const logs = await SyncLog.find(query)
      .sort({ startedAt: -1 })
      .skip(skip)
      .limit(parseInt(limit))
      .populate('serviceId', 'name repoURL alias')
      .populate('userId', 'email username');

    const total = await SyncLog.countDocuments(query);

    res.json({
      success: true,
      data: logs,
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
 * Get log by ID
 * GET /api/v1/logs/:id
 */
export const getLogById = async (req, res, next) => {
  try {
    const { id } = req.params;

    const log = await SyncLog.findById(id)
      .populate('serviceId', 'name repoURL alias repository')
      .populate('userId', 'email username');

    if (!log) {
      throw new AppError('Log not found', 404, 'LOG_NOT_FOUND');
    }

    res.json({
      success: true,
      data: log
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Get logs for a specific service
 * GET /api/v1/services/:serviceId/logs
 */
export const getServiceLogs = async (req, res, next) => {
  try {
    const { serviceId } = req.params;
    const { limit = 10 } = req.query;

    const logs = await SyncLog.getRecentLogs(serviceId, parseInt(limit));

    res.json({
      success: true,
      data: logs
    });
  } catch (error) {
    next(error);
  }
};

