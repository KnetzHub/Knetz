import mongoose from 'mongoose';

const apiKeySchema = new mongoose.Schema({
  key: {
    type: String,
    required: true
  },
  name: {
    type: String,
    required: true
  },
  createdAt: {
    type: Date,
    default: Date.now
  },
  expiresAt: Date,
  lastUsed: Date
}, { _id: false });

const tokenSchema = new mongoose.Schema({
  accessToken: String,
  refreshToken: String,
  expiresAt: Date
}, { _id: false });

const settingsSchema = new mongoose.Schema({
  defaultTrackingMethod: {
    type: String,
    enum: ['cron', 'webhook', 'manual'],
    default: 'manual'
  },
  notifications: {
    type: Boolean,
    default: true
  },
  webhookURL: String
}, { _id: false });

const userSchema = new mongoose.Schema({
  email: {
    type: String,
    required: true,
    unique: true,
    lowercase: true,
    trim: true,
    index: true
  },
  username: {
    type: String,
    unique: true,
    sparse: true,
    trim: true,
    index: true
  },
  authProvider: {
    type: String,
    enum: ['gmail', 'github', 'sso', 'local'],
    default: 'local'
  },
  authProviderId: {
    type: String,
    index: true
  },
  apiKeys: [apiKeySchema],
  tokens: {
    github: tokenSchema,
    gitlab: tokenSchema,
    docker: tokenSchema
  },
  settings: {
    type: settingsSchema,
    default: () => ({})
  },
  lastLoginAt: Date
}, {
  timestamps: true,
  collection: 'users'
});

// Methods
userSchema.methods.addAPIKey = function(keyData) {
  this.apiKeys.push(keyData);
  return this;
};

userSchema.methods.removeAPIKey = function(keyName) {
  this.apiKeys = this.apiKeys.filter(k => k.name !== keyName);
  return this;
};

userSchema.methods.getActiveAPIKeys = function() {
  const now = new Date();
  return this.apiKeys.filter(k => !k.expiresAt || k.expiresAt > now);
};

// Static methods
userSchema.statics.findByEmail = function(email) {
  return this.findOne({ email: email.toLowerCase() });
};

userSchema.statics.findByAPIKey = function(hashedKey) {
  return this.findOne({ 'apiKeys.key': hashedKey });
};

const User = mongoose.model('User', userSchema);

export default User;

