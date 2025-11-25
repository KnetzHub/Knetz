import mongoose from 'mongoose';

const webhookSchema = new mongoose.Schema({
  serviceId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Service',
    required: true,
    index: true
  },
  userId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
    required: true
  },
  provider: {
    type: String,
    required: true,
    enum: ['github', 'gitlab', 'docker', 'quay']
  },
  webhookId: {
    type: String,
    required: true,
    index: true
  },
  events: [{
    type: String,
    enum: ['push', 'release', 'tag', 'pull_request', 'create', 'delete']
  }],
  secret: {
    type: String,
    required: true
  },
  active: {
    type: Boolean,
    default: true,
    index: true
  },
  lastTriggered: Date
}, {
  timestamps: true,
  collection: 'webhooks'
});

// Indexes
webhookSchema.index({ serviceId: 1, provider: 1 });

// Methods
webhookSchema.methods.updateLastTriggered = function() {
  this.lastTriggered = new Date();
  return this;
};

webhookSchema.methods.disable = function() {
  this.active = false;
  return this;
};

webhookSchema.methods.enable = function() {
  this.active = true;
  return this;
};

// Static methods
webhookSchema.statics.findByService = function(serviceId) {
  return this.find({ serviceId, active: true });
};

webhookSchema.statics.findByWebhookId = function(webhookId) {
  return this.findOne({ webhookId, active: true });
};

const Webhook = mongoose.model('Webhook', webhookSchema);

export default Webhook;

