import mongoose from 'mongoose';

const versionSchema = new mongoose.Schema({
  version: {
    type: String,
    required: true
  },
  type: {
    type: String,
    enum: ['release', 'tag', 'commit', 'image'],
    default: 'release'
  },
  releasedAt: {
    type: Date,
    default: Date.now
  },
  changelog: String,
  artifacts: [String],
  isMajor: {
    type: Boolean,
    default: false
  },
  isPrerelease: {
    type: Boolean,
    default: false
  }
}, { _id: false });

const dependencySchema = new mongoose.Schema({
  name: {
    type: String,
    required: true
  },
  currentVersion: String,
  requiredVersion: String,
  serviceId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Service'
  },
  status: {
    type: String,
    enum: ['up-to-date', 'outdated', 'deprecated', 'external'],
    default: 'external'
  }
}, { _id: false });

const metadataSchema = new mongoose.Schema({
  language: String,
  framework: String,
  packageManager: String,
  manifestFile: String,
  manifestPath: String
}, { _id: false });

const serviceSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
    trim: true
  },
  alias: {
    type: String,
    trim: true,
    index: true
  },
  repository: {
    type: String,
    required: true,
    enum: ['github', 'gitlab', 'docker', 'quay'],
    index: true
  },
  repoURL: {
    type: String,
    required: true,
    unique: true,
    index: true
  },
  repoOwner: {
    type: String,
    required: true
  },
  repoName: {
    type: String,
    required: true
  },
  visibility: {
    type: String,
    enum: ['public', 'private'],
    default: 'public'
  },
  trackingMethod: {
    type: String,
    enum: ['cron', 'webhook', 'manual'],
    default: 'manual'
  },
  lastSyncedAt: {
    type: Date,
    index: true
  },
  userId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
    index: true
  },
  versions: [versionSchema],
  dependencies: [dependencySchema],
  metadata: metadataSchema
}, {
  timestamps: true,
  collection: 'dependencies_hub'
});

// Indexes
serviceSchema.index({ name: 1 });
serviceSchema.index({ 'versions.version': 1 });

// Methods
serviceSchema.methods.getLatestVersion = function() {
  if (this.versions.length === 0) return null;
  return this.versions.sort((a, b) => 
    new Date(b.releasedAt) - new Date(a.releasedAt)
  )[0];
};

serviceSchema.methods.addVersion = function(versionData) {
  const exists = this.versions.find(v => v.version === versionData.version);
  if (!exists) {
    this.versions.push(versionData);
    this.lastSyncedAt = new Date();
  }
  return this;
};

// Static methods
serviceSchema.statics.findByURL = function(repoURL) {
  return this.findOne({ repoURL });
};

serviceSchema.statics.findByAlias = function(alias) {
  return this.findOne({ alias });
};

const Service = mongoose.model('Service', serviceSchema);

export default Service;

