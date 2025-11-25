import axios from 'axios';
import logger from '../utils/logger.js';
import { AppError } from '../middleware/errorHandler.js';

class GitHubService {
  constructor() {
    this.baseURL = 'https://api.github.com';
    this.clientId = process.env.GITHUB_CLIENT_ID;
    this.clientSecret = process.env.GITHUB_CLIENT_SECRET;
  }

  /**
   * Create authenticated axios instance
   */
  createClient(accessToken = null) {
    const headers = {
      'Accept': 'application/vnd.github.v3+json',
      'User-Agent': 'KnetZ-App'
    };

    if (accessToken) {
      headers['Authorization'] = `token ${accessToken}`;
    }

    return axios.create({
      baseURL: this.baseURL,
      headers,
      timeout: 10000
    });
  }

  /**
   * Exchange OAuth code for access token
   */
  async exchangeCodeForToken(code) {
    try {
      const response = await axios.post(
        'https://github.com/login/oauth/access_token',
        {
          client_id: this.clientId,
          client_secret: this.clientSecret,
          code
        },
        {
          headers: {
            'Accept': 'application/json'
          }
        }
      );

      if (response.data.error) {
        throw new AppError(response.data.error_description, 400, 'GITHUB_OAUTH_ERROR');
      }

      return response.data.access_token;
    } catch (error) {
      logger.error('GitHub OAuth error:', error);
      throw new AppError('Failed to authenticate with GitHub', 500, 'GITHUB_AUTH_FAILED');
    }
  }

  /**
   * Get authenticated user info
   */
  async getUserInfo(accessToken) {
    try {
      const client = this.createClient(accessToken);
      const response = await client.get('/user');
      return response.data;
    } catch (error) {
      logger.error('Failed to fetch GitHub user info:', error);
      throw new AppError('Failed to fetch user info', 500, 'GITHUB_API_ERROR');
    }
  }

  /**
   * Fetch repository info
   */
  async getRepository(owner, repo, accessToken = null) {
    try {
      const client = this.createClient(accessToken);
      const response = await client.get(`/repos/${owner}/${repo}`);
      return response.data;
    } catch (error) {
      if (error.response?.status === 404) {
        throw new AppError('Repository not found', 404, 'REPO_NOT_FOUND');
      }
      logger.error('Failed to fetch repository:', error);
      throw new AppError('Failed to fetch repository info', 500, 'GITHUB_API_ERROR');
    }
  }

  /**
   * Fetch releases from a repository
   */
  async getReleases(owner, repo, accessToken = null) {
    try {
      const client = this.createClient(accessToken);
      const response = await client.get(`/repos/${owner}/${repo}/releases`, {
        params: {
          per_page: 100
        }
      });

      logger.info(`Fetched ${response.data.length} releases for ${owner}/${repo}`);
      
      return response.data.map(release => ({
        version: release.tag_name.replace(/^v/, ''), // Remove 'v' prefix
        type: 'release',
        releasedAt: new Date(release.published_at),
        changelog: release.body || '',
        artifacts: release.assets.map(asset => asset.browser_download_url),
        isMajor: this.isMajorVersion(release.tag_name),
        isPrerelease: release.prerelease
      }));
    } catch (error) {
      if (error.response?.status === 404) {
        logger.info(`No releases found for ${owner}/${repo}`);
        return [];
      }
      logger.error('Failed to fetch releases:', error);
      throw new AppError('Failed to fetch releases', 500, 'GITHUB_API_ERROR');
    }
  }

  /**
   * Fetch tags from a repository
   */
  async getTags(owner, repo, accessToken = null) {
    try {
      const client = this.createClient(accessToken);
      const response = await client.get(`/repos/${owner}/${repo}/tags`, {
        params: {
          per_page: 100
        }
      });

      logger.info(`Fetched ${response.data.length} tags for ${owner}/${repo}`);
      
      return response.data.map(tag => ({
        version: tag.name.replace(/^v/, ''),
        type: 'tag',
        releasedAt: new Date(),
        changelog: '',
        artifacts: [],
        isMajor: this.isMajorVersion(tag.name),
        isPrerelease: false
      }));
    } catch (error) {
      if (error.response?.status === 404) {
        logger.info(`No tags found for ${owner}/${repo}`);
        return [];
      }
      logger.error('Failed to fetch tags:', error);
      throw new AppError('Failed to fetch tags', 500, 'GITHUB_API_ERROR');
    }
  }

  /**
   * Fetch file content from repository
   */
  async getFileContent(owner, repo, path, accessToken = null) {
    try {
      const client = this.createClient(accessToken);
      const response = await client.get(`/repos/${owner}/${repo}/contents/${path}`);
      
      // Decode base64 content
      const content = Buffer.from(response.data.content, 'base64').toString('utf-8');
      return content;
    } catch (error) {
      if (error.response?.status === 404) {
        throw new AppError(`File ${path} not found`, 404, 'FILE_NOT_FOUND');
      }
      logger.error('Failed to fetch file content:', error);
      throw new AppError('Failed to fetch file content', 500, 'GITHUB_API_ERROR');
    }
  }

  /**
   * Search for manifest files in repository
   */
  async findManifestFile(owner, repo, accessToken = null) {
    const manifestFiles = [
      'package.json',
      'go.mod',
      'requirements.txt',
      'Pipfile',
      'pom.xml',
      'build.gradle',
      'Cargo.toml',
      'composer.json'
    ];

    for (const file of manifestFiles) {
      try {
        const content = await this.getFileContent(owner, repo, file, accessToken);
        return { file, content };
      } catch (error) {
        // Continue to next file
        continue;
      }
    }

    return null;
  }

  /**
   * Get list of accessible repositories for user
   */
  async getUserRepositories(accessToken, page = 1, perPage = 30) {
    try {
      const client = this.createClient(accessToken);
      const response = await client.get('/user/repos', {
        params: {
          sort: 'updated',
          per_page: perPage,
          page
        }
      });

      return response.data.map(repo => ({
        id: repo.id,
        name: repo.name,
        fullName: repo.full_name,
        owner: repo.owner.login,
        description: repo.description,
        url: repo.html_url,
        language: repo.language,
        visibility: repo.private ? 'private' : 'public',
        stars: repo.stargazers_count,
        forks: repo.forks_count,
        updatedAt: repo.updated_at
      }));
    } catch (error) {
      logger.error('Failed to fetch user repositories:', error);
      throw new AppError('Failed to fetch repositories', 500, 'GITHUB_API_ERROR');
    }
  }

  /**
   * Create webhook for repository
   */
  async createWebhook(owner, repo, callbackURL, secret, accessToken) {
    try {
      const client = this.createClient(accessToken);
      const response = await client.post(`/repos/${owner}/${repo}/hooks`, {
        name: 'web',
        active: true,
        events: ['push', 'release', 'create'],
        config: {
          url: callbackURL,
          content_type: 'json',
          secret: secret,
          insecure_ssl: '0'
        }
      });

      return response.data;
    } catch (error) {
      logger.error('Failed to create webhook:', error);
      throw new AppError('Failed to create webhook', 500, 'WEBHOOK_CREATE_FAILED');
    }
  }

  /**
   * Delete webhook
   */
  async deleteWebhook(owner, repo, hookId, accessToken) {
    try {
      const client = this.createClient(accessToken);
      await client.delete(`/repos/${owner}/${repo}/hooks/${hookId}`);
      return true;
    } catch (error) {
      logger.error('Failed to delete webhook:', error);
      throw new AppError('Failed to delete webhook', 500, 'WEBHOOK_DELETE_FAILED');
    }
  }

  /**
   * Check if version is a major release
   */
  isMajorVersion(version) {
    const cleaned = version.replace(/^v/, '');
    const parts = cleaned.split('.');
    
    if (parts.length >= 2) {
      // Major version changed or breaking change indicators
      return parts[1] === '0' && parts[0] !== '0';
    }
    
    return false;
  }

  /**
   * Parse repository URL
   */
  parseRepoURL(url) {
    // Support formats:
    // - https://github.com/owner/repo
    // - github.com/owner/repo
    // - owner/repo
    
    const cleaned = url.replace(/^https?:\/\//, '').replace(/\.git$/, '');
    const parts = cleaned.split('/');
    
    if (parts.length >= 2) {
      const owner = parts[parts.length - 2];
      const repo = parts[parts.length - 1];
      return { owner, repo };
    }
    
    throw new AppError('Invalid repository URL format', 400, 'INVALID_REPO_URL');
  }
}

export default new GitHubService();

