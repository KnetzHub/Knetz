import Service from '../models/Service.js';
import SyncLog from '../models/SyncLog.js';
import Event from '../models/Event.js';
import githubService from './githubService.js';
import manifestParser from './manifestParser.js';
import logger from '../utils/logger.js';
import semver from 'semver';

class SyncService {
  /**
   * Sync a service - fetch versions and update database
   */
  async syncService(serviceId, options = {}) {
    const { triggeredBy = 'manual', userId = null, accessToken = null } = options;

    // Create sync log
    const syncLog = new SyncLog({
      serviceId,
      serviceName: '',
      syncType: options.syncType || 'manual',
      status: 'pending',
      triggeredBy,
      userId,
      startedAt: new Date()
    });

    try {
      // Get service
      const service = await Service.findById(serviceId);
      if (!service) {
        throw new Error('Service not found');
      }

      syncLog.serviceName = service.name;
      await syncLog.save();

      logger.info(`Starting sync for service: ${service.name} (${service.repository})`);

      // Sync based on repository type
      let versions = [];
      let manifestData = null;

      switch (service.repository) {
        case 'github':
          versions = await this.syncGitHub(service, accessToken);
          manifestData = await this.syncManifest(service, accessToken);
          break;
        case 'gitlab':
          // TODO: Implement GitLab sync
          throw new Error('GitLab sync not yet implemented');
        case 'docker':
          // TODO: Implement Docker sync
          throw new Error('Docker sync not yet implemented');
        default:
          throw new Error(`Unsupported repository type: ${service.repository}`);
      }

      // Update service with new versions
      const { added, updated } = await this.updateServiceVersions(service, versions);

      // Update manifest metadata if found
      if (manifestData) {
        service.metadata = {
          language: manifestData.language,
          packageManager: manifestData.packageManager,
          manifestFile: manifestData.manifestFile,
          manifestPath: manifestData.manifestPath || manifestData.manifestFile
        };
        service.dependencies = manifestData.dependencies || [];
      }

      service.lastSyncedAt = new Date();
      await service.save();

      // Mark sync as successful
      syncLog.markComplete(true, {
        versionsAdded: added,
        versionsUpdated: updated
      });
      await syncLog.save();

      // Create events for new versions
      if (added > 0) {
        await Event.createEvent({
          type: 'version.detected',
          severity: 'info',
          serviceId: service._id,
          userId,
          message: `${added} new version(s) detected for ${service.name}`,
          details: { versionsAdded: added }
        });
      }

      logger.info(`Sync completed for ${service.name}: ${added} added, ${updated} updated`);

      return {
        success: true,
        versionsAdded: added,
        versionsUpdated: updated,
        service
      };

    } catch (error) {
      logger.error(`Sync failed for service ${serviceId}:`, error);

      syncLog.markComplete(false, {
        errors: [{
          message: error.message,
          code: error.code || 'SYNC_ERROR',
          timestamp: new Date()
        }]
      });
      await syncLog.save();

      await Event.createEvent({
        type: 'sync.failed',
        severity: 'error',
        serviceId,
        userId,
        message: `Sync failed: ${error.message}`,
        details: { error: error.message }
      });

      throw error;
    }
  }

  /**
   * Sync GitHub repository
   */
  async syncGitHub(service, accessToken = null) {
    const { repoOwner, repoName } = service;

    // Fetch releases and tags
    const [releases, tags] = await Promise.all([
      githubService.getReleases(repoOwner, repoName, accessToken),
      githubService.getTags(repoOwner, repoName, accessToken)
    ]);

    // Combine and deduplicate versions
    const allVersions = [...releases, ...tags];
    const uniqueVersions = this.deduplicateVersions(allVersions);

    // Filter valid semantic versions
    const validVersions = uniqueVersions.filter(v => 
      semver.valid(semver.coerce(v.version))
    );

    logger.info(`Found ${validVersions.length} valid versions for ${repoOwner}/${repoName}`);

    return validVersions;
  }

  /**
   * Sync manifest file and parse dependencies
   */
  async syncManifest(service, accessToken = null) {
    try {
      const { repoOwner, repoName } = service;

      // Find manifest file
      const manifest = await githubService.findManifestFile(repoOwner, repoName, accessToken);

      if (!manifest) {
        logger.info(`No manifest file found for ${repoOwner}/${repoName}`);
        return null;
      }

      logger.info(`Found manifest file: ${manifest.file}`);

      // Parse manifest
      const parsed = manifestParser.parse(manifest.file, manifest.content);

      if (!parsed) {
        logger.warn(`Failed to parse manifest file: ${manifest.file}`);
        return null;
      }

      // Link dependencies to tracked services
      if (parsed.dependencies && parsed.dependencies.length > 0) {
        parsed.dependencies = await this.linkDependencies(parsed.dependencies);
      }

      return {
        ...parsed,
        manifestPath: manifest.file
      };

    } catch (error) {
      logger.error('Failed to sync manifest:', error);
      return null;
    }
  }

  /**
   * Link dependencies to tracked services in database
   */
  async linkDependencies(dependencies) {
    const linkedDeps = [];

    for (const dep of dependencies) {
      try {
        // Try to find service by name
        const trackedService = await Service.findByName(dep.name);

        if (trackedService) {
          // Link to tracked service
          dep.serviceId = trackedService._id;
          dep.status = await this.compareDependencyVersion(
            dep.requiredVersion,
            trackedService.versions
          );
          
          logger.info(`Linked dependency ${dep.name} to tracked service`);
        } else {
          // External dependency (not tracked)
          dep.status = 'external';
        }

        linkedDeps.push(dep);
      } catch (error) {
        logger.warn(`Failed to link dependency ${dep.name}:`, error);
        linkedDeps.push({ ...dep, status: 'external' });
      }
    }

    return linkedDeps;
  }

  /**
   * Compare dependency version with available versions
   */
  async compareDependencyVersion(requiredVersion, availableVersions) {
    if (!availableVersions || availableVersions.length === 0) {
      return 'unknown';
    }

    // Get latest version
    const latestVersion = availableVersions[0]?.version;

    if (!latestVersion) {
      return 'unknown';
    }

    try {
      // Clean version strings
      const required = semver.coerce(requiredVersion);
      const latest = semver.coerce(latestVersion);

      if (!required || !latest) {
        return 'unknown';
      }

      // Compare versions
      if (semver.satisfies(latest, requiredVersion)) {
        return 'up-to-date';
      } else if (semver.gt(latest, required)) {
        return 'outdated';
      } else {
        return 'ahead';
      }
    } catch (error) {
      logger.warn('Version comparison failed:', error);
      return 'unknown';
    }
  }

  /**
   * Update service versions
   */
  async updateServiceVersions(service, newVersions) {
    let added = 0;
    let updated = 0;

    for (const version of newVersions) {
      const exists = service.versions.find(v => v.version === version.version);

      if (!exists) {
        service.versions.push(version);
        added++;
      } else if (this.shouldUpdateVersion(exists, version)) {
        Object.assign(exists, version);
        updated++;
      }
    }

    return { added, updated };
  }

  /**
   * Check if version should be updated
   */
  shouldUpdateVersion(existingVersion, newVersion) {
    // Update if new version has more information
    return (
      (!existingVersion.changelog && newVersion.changelog) ||
      (existingVersion.artifacts.length === 0 && newVersion.artifacts.length > 0) ||
      (existingVersion.type === 'tag' && newVersion.type === 'release')
    );
  }

  /**
   * Deduplicate versions (prefer releases over tags)
   */
  deduplicateVersions(versions) {
    const versionMap = new Map();

    for (const version of versions) {
      const key = version.version;

      if (!versionMap.has(key)) {
        versionMap.set(key, version);
      } else {
        const existing = versionMap.get(key);
        // Prefer releases over tags
        if (version.type === 'release' && existing.type === 'tag') {
          versionMap.set(key, version);
        }
      }
    }

    return Array.from(versionMap.values());
  }

  /**
   * Sync all services with cron tracking method
   */
  async syncAllCronServices() {
    const services = await Service.find({
      trackingMethod: 'cron',
      visibility: 'public'
    });

    logger.info(`Syncing ${services.length} cron-tracked services`);

    const results = [];

    for (const service of services) {
      try {
        const result = await this.syncService(service._id, {
          triggeredBy: 'system',
          syncType: 'cron'
        });
        results.push({ serviceId: service._id, success: true, ...result });
      } catch (error) {
        results.push({ 
          serviceId: service._id, 
          success: false, 
          error: error.message 
        });
      }
    }

    return results;
  }

  /**
   * Check if service needs sync
   */
  shouldSync(service, intervalMinutes = 60) {
    if (!service.lastSyncedAt) return true;

    const now = new Date();
    const lastSync = new Date(service.lastSyncedAt);
    const diffMinutes = (now - lastSync) / (1000 * 60);

    return diffMinutes >= intervalMinutes;
  }
}

export default new SyncService();

