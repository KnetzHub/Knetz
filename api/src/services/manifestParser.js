import logger from '../utils/logger.js';

class ManifestParser {
  /**
   * Parse package.json (Node.js)
   */
  parsePackageJson(content) {
    try {
      const data = JSON.parse(content);
      const dependencies = [];

      // Parse dependencies
      const depSections = ['dependencies', 'devDependencies', 'peerDependencies'];
      
      for (const section of depSections) {
        if (data[section]) {
          for (const [name, version] of Object.entries(data[section])) {
            dependencies.push({
              name,
              requiredVersion: this.cleanVersion(version),
              currentVersion: this.cleanVersion(version),
              status: 'external'
            });
          }
        }
      }

      return {
        language: 'javascript',
        packageManager: 'npm',
        manifestFile: 'package.json',
        currentVersion: data.version || '0.0.0',
        dependencies
      };
    } catch (error) {
      logger.error('Failed to parse package.json:', error);
      return null;
    }
  }

  /**
   * Parse go.mod (Go)
   */
  parseGoMod(content) {
    try {
      const dependencies = [];
      const lines = content.split('\n');
      
      let inRequireBlock = false;
      
      for (const line of lines) {
        const trimmed = line.trim();
        
        // Check for require block
        if (trimmed.startsWith('require (')) {
          inRequireBlock = true;
          continue;
        }
        
        if (trimmed === ')') {
          inRequireBlock = false;
          continue;
        }
        
        // Parse require line
        if (inRequireBlock || trimmed.startsWith('require ')) {
          const match = trimmed.match(/^(?:require\s+)?([^\s]+)\s+v?([^\s]+)/);
          if (match) {
            const [, name, version] = match;
            dependencies.push({
              name,
              requiredVersion: version.replace(/\+incompatible$/, ''),
              currentVersion: version.replace(/\+incompatible$/, ''),
              status: 'external'
            });
          }
        }
      }

      // Extract module version
      const moduleMatch = content.match(/^module\s+(.+)$/m);
      const moduleName = moduleMatch ? moduleMatch[1] : 'unknown';

      return {
        language: 'go',
        packageManager: 'go mod',
        manifestFile: 'go.mod',
        moduleName,
        dependencies
      };
    } catch (error) {
      logger.error('Failed to parse go.mod:', error);
      return null;
    }
  }

  /**
   * Parse requirements.txt (Python)
   */
  parseRequirementsTxt(content) {
    try {
      const dependencies = [];
      const lines = content.split('\n');
      
      for (const line of lines) {
        const trimmed = line.trim();
        
        // Skip comments and empty lines
        if (!trimmed || trimmed.startsWith('#')) continue;
        
        // Parse package==version or package>=version
        const match = trimmed.match(/^([a-zA-Z0-9_-]+)([>=<~!]+)(.+)$/);
        if (match) {
          const [, name, , version] = match;
          dependencies.push({
            name,
            requiredVersion: this.cleanVersion(version),
            currentVersion: this.cleanVersion(version),
            status: 'external'
          });
        } else {
          // Just package name without version
          dependencies.push({
            name: trimmed,
            requiredVersion: '*',
            currentVersion: '*',
            status: 'external'
          });
        }
      }

      return {
        language: 'python',
        packageManager: 'pip',
        manifestFile: 'requirements.txt',
        dependencies
      };
    } catch (error) {
      logger.error('Failed to parse requirements.txt:', error);
      return null;
    }
  }

  /**
   * Parse Pipfile (Python)
   */
  parsePipfile(content) {
    try {
      const dependencies = [];
      const lines = content.split('\n');
      
      let inPackagesBlock = false;
      
      for (const line of lines) {
        const trimmed = line.trim();
        
        if (trimmed === '[packages]' || trimmed === '[dev-packages]') {
          inPackagesBlock = true;
          continue;
        }
        
        if (trimmed.startsWith('[') && trimmed !== '[packages]' && trimmed !== '[dev-packages]') {
          inPackagesBlock = false;
          continue;
        }
        
        if (inPackagesBlock && trimmed.includes('=')) {
          const match = trimmed.match(/^(\S+)\s*=\s*"([^"]+)"/);
          if (match) {
            const [, name, version] = match;
            dependencies.push({
              name,
              requiredVersion: this.cleanVersion(version),
              currentVersion: this.cleanVersion(version),
              status: 'external'
            });
          }
        }
      }

      return {
        language: 'python',
        packageManager: 'pipenv',
        manifestFile: 'Pipfile',
        dependencies
      };
    } catch (error) {
      logger.error('Failed to parse Pipfile:', error);
      return null;
    }
  }

  /**
   * Parse Cargo.toml (Rust)
   */
  parseCargoToml(content) {
    try {
      const dependencies = [];
      const lines = content.split('\n');
      
      let inDependenciesBlock = false;
      let currentVersion = '0.0.0';
      
      // Get package version
      const versionMatch = content.match(/^\[package\][^[]*version\s*=\s*"([^"]+)"/ms);
      if (versionMatch) {
        currentVersion = versionMatch[1];
      }
      
      for (const line of lines) {
        const trimmed = line.trim();
        
        if (trimmed === '[dependencies]' || trimmed === '[dev-dependencies]') {
          inDependenciesBlock = true;
          continue;
        }
        
        if (trimmed.startsWith('[') && !trimmed.includes('dependencies]')) {
          inDependenciesBlock = false;
          continue;
        }
        
        if (inDependenciesBlock && trimmed.includes('=')) {
          const match = trimmed.match(/^(\S+)\s*=\s*"([^"]+)"/);
          if (match) {
            const [, name, version] = match;
            dependencies.push({
              name,
              requiredVersion: this.cleanVersion(version),
              currentVersion: this.cleanVersion(version),
              status: 'external'
            });
          }
        }
      }

      return {
        language: 'rust',
        packageManager: 'cargo',
        manifestFile: 'Cargo.toml',
        currentVersion,
        dependencies
      };
    } catch (error) {
      logger.error('Failed to parse Cargo.toml:', error);
      return null;
    }
  }

  /**
   * Auto-detect and parse manifest file
   */
  parse(filename, content) {
    switch (filename) {
      case 'package.json':
        return this.parsePackageJson(content);
      case 'go.mod':
        return this.parseGoMod(content);
      case 'requirements.txt':
        return this.parseRequirementsTxt(content);
      case 'Pipfile':
        return this.parsePipfile(content);
      case 'Cargo.toml':
        return this.parseCargoToml(content);
      default:
        logger.warn(`Unknown manifest file: ${filename}`);
        return null;
    }
  }

  /**
   * Clean version string
   */
  cleanVersion(version) {
    return version
      .replace(/^[~^>=<]+/, '') // Remove version prefixes
      .replace(/\s+.*$/, '') // Remove comments
      .trim();
  }
}

export default new ManifestParser();

