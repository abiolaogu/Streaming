// StreamVerse Unified Platform SDK
// Integrate with 10+ video platforms from a single API

import { EventEmitter } from 'eventemitter3';
import { YouTubeClient } from './platforms/youtube';
import { TwitchClient } from './platforms/twitch';
import { TikTokClient } from './platforms/tiktok';
import { VimeoClient } from './platforms/vimeo';
import { FacebookClient } from './platforms/facebook';
import { InstagramClient } from './platforms/instagram';
import { RumbleClient } from './platforms/rumble';
import { OdyseeClient } from './platforms/odysee';
import { KickClient } from './platforms/kick';
import { DailymotionClient } from './platforms/dailymotion';

export type Platform =
  | 'youtube'
  | 'twitch'
  | 'tiktok'
  | 'vimeo'
  | 'facebook'
  | 'instagram'
  | 'rumble'
  | 'odysee'
  | 'kick'
  | 'dailymotion';

export interface PlatformConfig {
  apiKey?: string;
  apiSecret?: string;
  accessToken?: string;
  refreshToken?: string;
  clientId?: string;
  clientSecret?: string;
}

export interface VideoMetadata {
  title: string;
  description: string;
  tags: string[];
  category?: string;
  privacy?: 'public' | 'private' | 'unlisted';
  thumbnailUrl?: string;
  language?: string;
}

export interface UploadOptions {
  filePath?: string;
  streamUrl?: string;
  metadata: VideoMetadata;
  platforms: Platform[];
  monetization?: {
    enabled: boolean;
    adSettings?: any;
  };
}

export interface UploadResult {
  platform: Platform;
  videoId: string;
  url: string;
  status: 'success' | 'failed' | 'processing';
  error?: string;
}

export interface LiveStreamOptions {
  title: string;
  description: string;
  platforms: Platform[];
  lowLatency?: boolean;
  resolution?: '720p' | '1080p' | '4K';
  privacy?: 'public' | 'private' | 'unlisted';
}

export interface LiveStreamResult {
  platform: Platform;
  streamKey: string;
  ingestUrl: string;
  watchUrl: string;
  status: 'live' | 'offline' | 'starting';
}

export class StreamVersePlatformSDK extends EventEmitter {
  private clients: Map<Platform, any>;
  private config: Map<Platform, PlatformConfig>;

  constructor() {
    super();
    this.clients = new Map();
    this.config = new Map();
  }

  /**
   * Configure authentication for a platform
   */
  configurePlatform(platform: Platform, config: PlatformConfig): void {
    this.config.set(platform, config);

    // Initialize client
    switch (platform) {
      case 'youtube':
        this.clients.set(platform, new YouTubeClient(config));
        break;
      case 'twitch':
        this.clients.set(platform, new TwitchClient(config));
        break;
      case 'tiktok':
        this.clients.set(platform, new TikTokClient(config));
        break;
      case 'vimeo':
        this.clients.set(platform, new VimeoClient(config));
        break;
      case 'facebook':
        this.clients.set(platform, new FacebookClient(config));
        break;
      case 'instagram':
        this.clients.set(platform, new InstagramClient(config));
        break;
      case 'rumble':
        this.clients.set(platform, new RumbleClient(config));
        break;
      case 'odysee':
        this.clients.set(platform, new OdyseeClient(config));
        break;
      case 'kick':
        this.clients.set(platform, new KickClient(config));
        break;
      case 'dailymotion':
        this.clients.set(platform, new DailymotionClient(config));
        break;
    }
  }

  /**
   * Upload video to multiple platforms simultaneously
   */
  async uploadVideo(options: UploadOptions): Promise<UploadResult[]> {
    const results: UploadResult[] = [];

    // Upload to all platforms in parallel
    const uploads = options.platforms.map(async (platform) => {
      try {
        const client = this.clients.get(platform);
        if (!client) {
          throw new Error(`Platform ${platform} not configured`);
        }

        this.emit('upload:start', { platform, metadata: options.metadata });

        const result = await client.uploadVideo({
          filePath: options.filePath,
          streamUrl: options.streamUrl,
          metadata: options.metadata,
        });

        this.emit('upload:success', { platform, result });

        return {
          platform,
          videoId: result.id,
          url: result.url,
          status: 'success' as const,
        };
      } catch (error) {
        this.emit('upload:error', { platform, error });

        return {
          platform,
          videoId: '',
          url: '',
          status: 'failed' as const,
          error: error instanceof Error ? error.message : 'Unknown error',
        };
      }
    });

    const settled = await Promise.allSettled(uploads);

    settled.forEach((result) => {
      if (result.status === 'fulfilled') {
        results.push(result.value);
      }
    });

    return results;
  }

  /**
   * Start live stream on multiple platforms
   */
  async startLiveStream(options: LiveStreamOptions): Promise<LiveStreamResult[]> {
    const results: LiveStreamResult[] = [];

    const streams = options.platforms.map(async (platform) => {
      try {
        const client = this.clients.get(platform);
        if (!client) {
          throw new Error(`Platform ${platform} not configured`);
        }

        this.emit('stream:start', { platform, options });

        const stream = await client.createLiveStream({
          title: options.title,
          description: options.description,
          lowLatency: options.lowLatency,
          resolution: options.resolution,
        });

        this.emit('stream:ready', { platform, stream });

        return {
          platform,
          streamKey: stream.streamKey,
          ingestUrl: stream.ingestUrl,
          watchUrl: stream.watchUrl,
          status: 'starting' as const,
        };
      } catch (error) {
        this.emit('stream:error', { platform, error });

        return {
          platform,
          streamKey: '',
          ingestUrl: '',
          watchUrl: '',
          status: 'offline' as const,
        };
      }
    });

    const settled = await Promise.allSettled(streams);

    settled.forEach((result) => {
      if (result.status === 'fulfilled') {
        results.push(result.value);
      }
    });

    return results;
  }

  /**
   * Stop live stream on platforms
   */
  async stopLiveStream(platforms: Platform[]): Promise<void> {
    const stops = platforms.map(async (platform) => {
      const client = this.clients.get(platform);
      if (client && client.stopLiveStream) {
        await client.stopLiveStream();
        this.emit('stream:stopped', { platform });
      }
    });

    await Promise.all(stops);
  }

  /**
   * Get analytics from all platforms
   */
  async getAnalytics(videoId: string, platforms: Platform[]): Promise<any> {
    const analytics = await Promise.all(
      platforms.map(async (platform) => {
        const client = this.clients.get(platform);
        if (!client) return null;

        try {
          const data = await client.getAnalytics(videoId);
          return { platform, data };
        } catch {
          return null;
        }
      })
    );

    return analytics.filter(Boolean);
  }

  /**
   * Sync metadata across platforms
   */
  async syncMetadata(
    videoId: string,
    metadata: VideoMetadata,
    platforms: Platform[]
  ): Promise<void> {
    await Promise.all(
      platforms.map(async (platform) => {
        const client = this.clients.get(platform);
        if (client && client.updateMetadata) {
          await client.updateMetadata(videoId, metadata);
          this.emit('metadata:synced', { platform, videoId });
        }
      })
    );
  }
}

// Example usage:
/*
const sdk = new StreamVersePlatformSDK();

// Configure platforms
sdk.configurePlatform('youtube', {
  apiKey: 'YOUR_YOUTUBE_API_KEY',
  clientId: 'YOUR_CLIENT_ID',
  clientSecret: 'YOUR_CLIENT_SECRET',
});

sdk.configurePlatform('twitch', {
  clientId: 'YOUR_TWITCH_CLIENT_ID',
  clientSecret: 'YOUR_TWITCH_CLIENT_SECRET',
  accessToken: 'YOUR_ACCESS_TOKEN',
});

// Upload to multiple platforms
const results = await sdk.uploadVideo({
  filePath: './video.mp4',
  metadata: {
    title: 'My Awesome Video',
    description: 'This is a test upload',
    tags: ['gaming', 'tutorial'],
    privacy: 'public',
  },
  platforms: ['youtube', 'twitch', 'tiktok'],
});

console.log('Upload results:', results);

// Start live stream
const streams = await sdk.startLiveStream({
  title: 'Live Gaming Session',
  description: 'Playing Fortnite',
  platforms: ['youtube', 'twitch', 'kick'],
  lowLatency: true,
});

console.log('Stream URLs:', streams);
*/

export * from './platforms/youtube';
export * from './platforms/twitch';
export * from './platforms/tiktok';
