import mongoose from 'mongoose';

const eventSchema = new mongoose.Schema({
  type: {
    type: String,
    required: true,
    index: true,
    enum: [
      'service.created',
      'service.updated',
      'service.deleted',
      'version.detected',
      'version.added',
      'dependency.outdated',
      'dependency.updated',
      'sync.started',
      'sync.completed',
      'sync.failed',
      'user.login',
      'apikey.created',
      'webhook.triggered'
    ]
  },
  severity: {
    type: String,
    enum: ['info', 'warning', 'error', 'critical'],
    default: 'info',
    index: true
  },
  serviceId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Service'
  },
  userId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
    index: true
  },
  message: {
    type: String,
    required: true
  },
  details: {
    type: mongoose.Schema.Types.Mixed,
    default: {}
  },
  readBy: [{
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User'
  }],
  isArchived: {
    type: Boolean,
    default: false,
    index: true
  }
}, {
  timestamps: true,
  collection: 'events'
});

// Indexes
eventSchema.index({ createdAt: -1 });
eventSchema.index({ type: 1, createdAt: -1 });
eventSchema.index({ severity: 1, isArchived: 1 });

// Methods
eventSchema.methods.markAsRead = function(userId) {
  if (!this.readBy.includes(userId)) {
    this.readBy.push(userId);
  }
  return this;
};

eventSchema.methods.archive = function() {
  this.isArchived = true;
  return this;
};

// Static methods
eventSchema.statics.getRecentEvents = function(userId, limit = 20) {
  const query = userId ? { userId } : {};
  return this.find({ ...query, isArchived: false })
    .sort({ createdAt: -1 })
    .limit(limit)
    .populate('serviceId', 'name alias')
    .populate('userId', 'email username');
};

eventSchema.statics.getUnreadCount = function(userId) {
  return this.countDocuments({
    userId,
    readBy: { $ne: userId },
    isArchived: false
  });
};

eventSchema.statics.createEvent = async function(eventData) {
  const event = new this(eventData);
  await event.save();
  return event;
};

const Event = mongoose.model('Event', eventSchema);

export default Event;

