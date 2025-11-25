import Service from '../models/Service.js';
import Event from '../models/Event.js';
import { AppError } from '../middleware/errorHandler.js';
import logger from '../utils/logger.js';

/**
 * Register a new service
 * POST /api/v1/services
 */
export const createService = async (req, res, next) => {
  try {
    const { name, alias, repository, repoURL, repoOwner, repoName, visibility, trackingMethod } = req.body;

    // Check if service already exists
    const existingService = await Service.findByURL(repoURL);
    if (existingService) {
      throw new AppError('Service with this repository URL already exists', 409, 'SERVICE_EXISTS');
    }

    // Create new service
    const service = new Service({
      name,
      alias,
      repository,
      repoURL,
      repoOwner,
      repoName,
      visibility: visibility || 'public',
      trackingMethod: trackingMethod || 'manual',
      userId: req.user?._id
    });

    await service.save();

    // Create event
    await Event.createEvent({
      type: 'service.created',
      severity: 'info',
      serviceId: service._id,
      userId: req.user?._id,
      message: `Service ${name} created successfully`,
      details: { repoURL, repository }
    });

    logger.info('Service created:', { serviceId: service._id, name, repoURL });

    res.status(201).json({
      success: true,
      data: service
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Get all services
 * GET /api/v1/services
 */
export const getServices = async (req, res, next) => {
  try {
    const { repository, trackingMethod, page = 1, limit = 20, search } = req.query;

    // Build query
    const query = {};
    if (repository) query.repository = repository;
    if (trackingMethod) query.trackingMethod = trackingMethod;
    if (req.user) query.userId = req.user._id;
    if (search) {
      query.$or = [
        { name: { $regex: search, $options: 'i' } },
        { alias: { $regex: search, $options: 'i' } },
        { repoURL: { $regex: search, $options: 'i' } }
      ];
    }

    // Pagination
    const skip = (page - 1) * limit;
    
    const services = await Service.find(query)
      .sort({ updatedAt: -1 })
      .skip(skip)
      .limit(parseInt(limit))
      .select('-dependencies -versions'); // Exclude large fields for list view

    const total = await Service.countDocuments(query);

    res.json({
      success: true,
      data: services,
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
 * Get service by ID
 * GET /api/v1/services/:id
 */
export const getServiceById = async (req, res, next) => {
  try {
    const { id } = req.params;

    const service = await Service.findById(id).populate('userId', 'email username');

    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    res.json({
      success: true,
      data: service
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Update service
 * PUT /api/v1/services/:id
 */
export const updateService = async (req, res, next) => {
  try {
    const { id } = req.params;
    const updates = req.body;

    // Prevent updating certain fields
    delete updates._id;
    delete updates.repoURL;
    delete updates.versions;
    delete updates.createdAt;

    const service = await Service.findByIdAndUpdate(
      id,
      { $set: updates },
      { new: true, runValidators: true }
    );

    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    // Create event
    await Event.createEvent({
      type: 'service.updated',
      severity: 'info',
      serviceId: service._id,
      userId: req.user?._id,
      message: `Service ${service.name} updated`,
      details: updates
    });

    logger.info('Service updated:', { serviceId: id, updates });

    res.json({
      success: true,
      data: service
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Delete service
 * DELETE /api/v1/services/:id
 */
export const deleteService = async (req, res, next) => {
  try {
    const { id } = req.params;

    const service = await Service.findByIdAndDelete(id);

    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    // Create event
    await Event.createEvent({
      type: 'service.deleted',
      severity: 'warning',
      userId: req.user?._id,
      message: `Service ${service.name} deleted`,
      details: { repoURL: service.repoURL }
    });

    logger.info('Service deleted:', { serviceId: id, name: service.name });

    res.json({
      success: true,
      message: 'Service deleted successfully'
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Get service versions
 * GET /api/v1/services/:id/versions
 */
export const getServiceVersions = async (req, res, next) => {
  try {
    const { id } = req.params;
    const { limit = 50, type } = req.query;

    const service = await Service.findById(id);

    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    let versions = service.versions;

    // Filter by type
    if (type) {
      versions = versions.filter(v => v.type === type);
    }

    // Sort by release date (newest first)
    versions = versions.sort((a, b) => new Date(b.releasedAt) - new Date(a.releasedAt));

    // Limit results
    versions = versions.slice(0, parseInt(limit));

    res.json({
      success: true,
      data: {
        serviceId: service._id,
        serviceName: service.name,
        totalVersions: service.versions.length,
        versions
      }
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Get latest version
 * GET /api/v1/services/:id/latest
 */
export const getLatestVersion = async (req, res, next) => {
  try {
    const { id } = req.params;

    const service = await Service.findById(id);

    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    const latestVersion = service.getLatestVersion();

    if (!latestVersion) {
      throw new AppError('No versions found for this service', 404, 'NO_VERSIONS');
    }

    res.json({
      success: true,
      data: {
        serviceId: service._id,
        serviceName: service.name,
        latestVersion
      }
    });
  } catch (error) {
    next(error);
  }
};

/**
 * Get service dependencies
 * GET /api/v1/services/:id/dependencies
 */
export const getServiceDependencies = async (req, res, next) => {
  try {
    const { id } = req.params;

    const service = await Service.findById(id).populate('dependencies.serviceId', 'name versions');

    if (!service) {
      throw new AppError('Service not found', 404, 'SERVICE_NOT_FOUND');
    }

    res.json({
      success: true,
      data: {
        serviceId: service._id,
        serviceName: service.name,
        currentVersion: service.getLatestVersion()?.version,
        dependencies: service.dependencies,
        metadata: service.metadata
      }
    });
  } catch (error) {
    next(error);
  }
};

