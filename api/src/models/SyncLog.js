import mongoose from 'mongoose';

const errorSchema = new mongoose.Schema({
  message: String,
  code: String,
  timestamp: {
    type: Date,
    default: Date.now
  }
}, { _id: false });

const metadataSchema = new mongoose.Schema({
  source: String,
  trigger: String,
  rateLimitRemaining: Number
}, { _id: false });

const syncLogSchema = new mongoose.Schema({
  serviceId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Service',
    required: true,
    index: true
  },
  serviceName: {
    type: String,
    required: true
  },
  syncType: {
    type: String,
    enum: ['cron', 'webhook', 'manual', 'cli'],
    required: true
  },
  status: {
    type: String,
    enum: ['success', 'failed', 'partial', 'pending'],
    default: 'pending',
    index: true
  },
  triggeredBy: {
    type: String,
    enum: ['system', 'user', 'webhook'],
    default: 'system',
    index: true
  },
  userId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User'
  },
  versionsAdded: {
    type: Number,
    default: 0
  },
  versionsUpdated: {
    type: Number,
    default: 0
  },
  errors: [errorSchema],
  duration: Number, // milliseconds
  startedAt: {
    type: Date,
    default: Date.now,
    index: -1 // Descending index for recent logs
  },
  completedAt: Date,
  metadata: metadataSchema
}, {
  timestamps: true,
  collection: 'sync_logs'
});

// Indexes
syncLogSchema.index({ serviceId: 1, startedAt: -1 });

// Methods
syncLogSchema.methods.markComplete = function(success = true, data = {}) {
  this.completedAt = new Date();
  this.status = success ? 'success' : 'failed';
  this.duration = this.completedAt - this.startedAt;
  
  if (data.versionsAdded) this.versionsAdded = data.versionsAdded;
  if (data.versionsUpdated) this.versionsUpdated = data.versionsUpdated;
  if (data.errors) this.errors = data.errors;
  
  return this;
};

// Static methods
syncLogSchema.statics.getRecentLogs = function(serviceId, limit = 10) {
  return this.find({ serviceId })
    .sort({ startedAt: -1 })
    .limit(limit);
};

syncLogSchema.statics.getFailedLogs = function(limit = 50) {
  return this.find({ status: 'failed' })
    .sort({ startedAt: -1 })
    .limit(limit)
    .populate('serviceId', 'name repoURL');
};

const SyncLog = mongoose.model('SyncLog', syncLogSchema);

export default SyncLog;

