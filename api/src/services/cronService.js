import cron from 'node-cron';
import syncService from './syncService.js';
import logger from '../utils/logger.js';

class CronService {
  constructor() {
    this.jobs = new Map();
  }

  /**
   * Start all cron jobs
   */
  start() {
    logger.info('Starting cron jobs...');

    // Medium priority: Sync every 4 hours
    this.scheduleJob('sync-medium', '0 */4 * * *', async () => {
      logger.info('Running medium priority sync...');
      try {
        const results = await syncService.syncAllCronServices();
        const success = results.filter(r => r.success).length;
        const failed = results.filter(r => !r.success).length;
        logger.info(`Medium priority sync completed: ${success} succeeded, ${failed} failed`);
      } catch (error) {
        logger.error('Medium priority sync error:', error);
      }
    });

    // High priority: Sync every 30 minutes (for critical services)
    this.scheduleJob('sync-high', '*/30 * * * *', async () => {
      logger.info('Running high priority sync...');
      // TODO: Implement priority-based sync
    });

    // Low priority: Sync daily at midnight
    this.scheduleJob('sync-low', '0 0 * * *', async () => {
      logger.info('Running low priority sync...');
      // TODO: Implement low priority sync for archived services
    });

    // Cleanup old logs: Run daily at 2 AM
    this.scheduleJob('cleanup-logs', '0 2 * * *', async () => {
      logger.info('Running log cleanup...');
      try {
        await this.cleanupOldLogs();
      } catch (error) {
        logger.error('Log cleanup error:', error);
      }
    });

    logger.info(`${this.jobs.size} cron jobs started`);
  }

  /**
   * Schedule a cron job
   */
  scheduleJob(name, schedule, task) {
    if (this.jobs.has(name)) {
      logger.warn(`Cron job ${name} already exists, skipping`);
      return;
    }

    const job = cron.schedule(schedule, task, {
      scheduled: true,
      timezone: 'UTC'
    });

    this.jobs.set(name, job);
    logger.info(`Scheduled cron job: ${name} (${schedule})`);
  }

  /**
   * Stop all cron jobs
   */
  stop() {
    logger.info('Stopping all cron jobs...');
    
    for (const [name, job] of this.jobs.entries()) {
      job.stop();
      logger.info(`Stopped cron job: ${name}`);
    }

    this.jobs.clear();
  }

  /**
   * Stop a specific cron job
   */
  stopJob(name) {
    const job = this.jobs.get(name);
    if (job) {
      job.stop();
      this.jobs.delete(name);
      logger.info(`Stopped cron job: ${name}`);
    }
  }

  /**
   * Clean up old logs (older than 30 days)
   */
  async cleanupOldLogs() {
    const SyncLog = (await import('../models/SyncLog.js')).default;
    const Event = (await import('../models/Event.js')).default;

    const thirtyDaysAgo = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000);

    // Delete old sync logs
    const syncLogsDeleted = await SyncLog.deleteMany({
      createdAt: { $lt: thirtyDaysAgo }
    });

    // Archive old events
    const eventsArchived = await Event.updateMany(
      { 
        createdAt: { $lt: thirtyDaysAgo },
        isArchived: false 
      },
      { $set: { isArchived: true } }
    );

    logger.info(`Cleanup completed: ${syncLogsDeleted.deletedCount} sync logs deleted, ${eventsArchived.modifiedCount} events archived`);
  }

  /**
   * Get status of all cron jobs
   */
  getStatus() {
    const status = {};
    
    for (const [name, job] of this.jobs.entries()) {
      status[name] = {
        running: job.options.scheduled
      };
    }

    return status;
  }
}

export default new CronService();

