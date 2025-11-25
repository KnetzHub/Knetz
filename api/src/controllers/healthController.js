import mongoose from 'mongoose';
import Service from '../models/Service.js';
import User from '../models/User.js';
import SyncLog from '../models/SyncLog.js';
import Event from '../models/Event.js';

/**
 * Health check endpoint
 * GET /api/v1/health
 */
export const healthCheck = async (req, res) => {
  try {
    // Check MongoDB connection
    const dbStatus = mongoose.connection.readyState === 1 ? 'connected' : 'disconnected';
    
    // Basic ping to database
    if (dbStatus === 'connected') {
      await mongoose.connection.db.admin().ping();
    }

    const health = {
      status: dbStatus === 'connected' ? 'healthy' : 'unhealthy',
      timestamp: new Date().toISOString(),
      uptime: process.uptime(),
      environment: process.env.NODE_ENV || 'development',
      version: '1.0.0',
      services: {
        database: {
          status: dbStatus,
          name: mongoose.connection.name
        },
        api: {
          status: 'running'
        }
      }
    };

    const statusCode = health.status === 'healthy' ? 200 : 503;
    res.status(statusCode).json(health);
  } catch (error) {
    res.status(503).json({
      status: 'unhealthy',
      timestamp: new Date().toISOString(),
      error: error.message
    });
  }
};

/**
 * Get platform statistics
 * GET /api/v1/stats
 */
export const getStats = async (req, res, next) => {
  try {
    // Get counts
    const [
      totalServices,
      totalUsers,
      totalSyncs,
      recentEvents,
      syncSuccess,
      syncFailed
    ] = await Promise.all([
      Service.countDocuments(),
      User.countDocuments(),
      SyncLog.countDocuments(),
      Event.countDocuments({ createdAt: { $gte: new Date(Date.now() - 24 * 60 * 60 * 1000) } }),
      SyncLog.countDocuments({ status: 'success', completedAt: { $gte: new Date(Date.now() - 24 * 60 * 60 * 1000) } }),
      SyncLog.countDocuments({ status: 'failed', completedAt: { $gte: new Date(Date.now() - 24 * 60 * 60 * 1000) } })
    ]);

    // Get services by repository
    const servicesByRepo = await Service.aggregate([
      {
        $group: {
          _id: '$repository',
          count: { $sum: 1 }
        }
      }
    ]);

    // Get services by tracking method
    const servicesByTracking = await Service.aggregate([
      {
        $group: {
          _id: '$trackingMethod',
          count: { $sum: 1 }
        }
      }
    ]);

    // Get recent syncs
    const recentSyncs = await SyncLog.find()
      .sort({ startedAt: -1 })
      .limit(10)
      .populate('serviceId', 'name repoURL')
      .select('serviceName status startedAt completedAt duration versionsAdded');

    // Calculate success rate
    const totalRecentSyncs = syncSuccess + syncFailed;
    const successRate = totalRecentSyncs > 0 ? ((syncSuccess / totalRecentSyncs) * 100).toFixed(2) : 0;

    res.json({
      success: true,
      data: {
        overview: {
          totalServices,
          totalUsers,
          totalSyncs,
          recentEvents
        },
        sync: {
          last24Hours: {
            success: syncSuccess,
            failed: syncFailed,
            total: totalRecentSyncs,
            successRate: `${successRate}%`
          },
          recent: recentSyncs
        },
        services: {
          byRepository: servicesByRepo.reduce((acc, item) => {
            acc[item._id] = item.count;
            return acc;
          }, {}),
          byTrackingMethod: servicesByTracking.reduce((acc, item) => {
            acc[item._id] = item.count;
            return acc;
          }, {})
        },
        timestamp: new Date().toISOString()
      }
    });
  } catch (error) {
    next(error);
  }
};

