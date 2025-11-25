import semver from 'semver';

/**
 * Compare semantic versions
 */
export const compareVersions = (v1, v2) => {
  try {
    const sv1 = semver.coerce(v1);
    const sv2 = semver.coerce(v2);
    
    if (!sv1 || !sv2) return 0;
    
    return semver.compare(sv1, sv2);
  } catch (error) {
    return 0;
  }
};

/**
 * Check if version is outdated
 */
export const isOutdated = (currentVersion, latestVersion) => {
  try {
    const sv1 = semver.coerce(currentVersion);
    const sv2 = semver.coerce(latestVersion);
    
    if (!sv1 || !sv2) return false;
    
    return semver.lt(sv1, sv2);
  } catch (error) {
    return false;
  }
};

/**
 * Parse semantic version
 */
export const parseVersion = (versionString) => {
  const cleaned = versionString.replace(/^v/, '');
  const coerced = semver.coerce(cleaned);
  return coerced ? coerced.version : null;
};

/**
 * Check if version is valid semver
 */
export const isValidSemver = (versionString) => {
  return semver.valid(semver.coerce(versionString)) !== null;
};

/**
 * Get major version number
 */
export const getMajorVersion = (versionString) => {
  const coerced = semver.coerce(versionString);
  return coerced ? semver.major(coerced) : null;
};

/**
 * Sanitize URL
 */
export const sanitizeURL = (url) => {
  return url
    .replace(/^https?:\/\//, '')
    .replace(/\.git$/, '')
    .replace(/\/$/, '')
    .toLowerCase();
};

/**
 * Generate unique identifier
 */
export const generateId = () => {
  return Date.now().toString(36) + Math.random().toString(36).substring(2);
};

/**
 * Format duration in milliseconds to human readable
 */
export const formatDuration = (ms) => {
  if (ms < 1000) return `${ms}ms`;
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
  if (ms < 3600000) return `${(ms / 60000).toFixed(1)}m`;
  return `${(ms / 3600000).toFixed(1)}h`;
};

/**
 * Delay execution
 */
export const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

/**
 * Retry function with exponential backoff
 */
export const retry = async (fn, maxRetries = 3, delayMs = 1000) => {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await fn();
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      await delay(delayMs * Math.pow(2, i));
    }
  }
};

/**
 * Chunk array into smaller arrays
 */
export const chunk = (array, size) => {
  const chunks = [];
  for (let i = 0; i < array.length; i += size) {
    chunks.push(array.slice(i, i + size));
  }
  return chunks;
};

/**
 * Remove duplicates from array
 */
export const unique = (array) => [...new Set(array)];

/**
 * Deep clone object
 */
export const deepClone = (obj) => JSON.parse(JSON.stringify(obj));

/**
 * Capitalize first letter
 */
export const capitalize = (str) => str.charAt(0).toUpperCase() + str.slice(1);

