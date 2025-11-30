// YouTube API Integration Client

import { google, youtube_v3 } from 'googleapis';
import { PlatformConfig } from '../index';
import * as fs from 'fs';

export class YouTubeClient {
  private youtube: youtube_v3.Youtube;
  private oauth2Client: any;

  constructor(config: PlatformConfig) {
    this.oauth2Client = new google.auth.OAuth2(
      config.clientId,
      config.clientSecret,
      'http://localhost:3000/oauth/callback'
    );

    if (config.accessToken) {
      this.oauth2Client.setCredentials({
        access_token: config.accessToken,
        refresh_token: config.refreshToken,
      });
    }

    this.youtube = google.youtube({
      version: 'v3',
      auth: this.oauth2Client,
    });
  }

  async uploadVideo(options: {
    filePath?: string;
    streamUrl?: string;
    metadata: any;
  }): Promise<{ id: string; url: string }> {
    const requestBody = {
      snippet: {
        title: options.metadata.title,
        description: options.metadata.description,
        tags: options.metadata.tags,
        categoryId: options.metadata.category || '22', // People & Blogs
      },
      status: {
        privacyStatus: options.metadata.privacy || 'public',
        selfDeclaredMadeForKids: false,
      },
    };

    const media = {
      body: options.filePath ? fs.createReadStream(options.filePath) : undefined,
    };

    const response = await this.youtube.videos.insert({
      part: ['snippet', 'status'],
      requestBody,
      media,
    });

    const videoId = response.data.id!;

    return {
      id: videoId,
      url: `https://www.youtube.com/watch?v=${videoId}`,
    };
  }

  async createLiveStream(options: {
    title: string;
    description: string;
    lowLatency?: boolean;
    resolution?: string;
  }): Promise<{
    streamKey: string;
    ingestUrl: string;
    watchUrl: string;
  }> {
    // Create broadcast
    const broadcastResponse = await this.youtube.liveBroadcasts.insert({
      part: ['snippet', 'contentDetails', 'status'],
      requestBody: {
        snippet: {
          title: options.title,
          description: options.description,
          scheduledStartTime: new Date().toISOString(),
        },
        contentDetails: {
          enableLowLatency: options.lowLatency || false,
        },
        status: {
          privacyStatus: 'public',
        },
      },
    });

    // Create stream
    const streamResponse = await this.youtube.liveStreams.insert({
      part: ['snippet', 'cdn'],
      requestBody: {
        snippet: {
          title: `${options.title} - Stream`,
        },
        cdn: {
          frameRate: '60fps',
          resolution: options.resolution || '1080p',
          ingestionType: 'rtmp',
        },
      },
    });

    const broadcastId = broadcastResponse.data.id!;
    const streamId = streamResponse.data.id!;

    // Bind broadcast to stream
    await this.youtube.liveBroadcasts.bind({
      part: ['id'],
      id: broadcastId,
      streamId: streamId,
    });

    const streamKey = streamResponse.data.cdn?.ingestionInfo?.streamName!;
    const ingestUrl = streamResponse.data.cdn?.ingestionInfo?.ingestionAddress!;

    return {
      streamKey,
      ingestUrl,
      watchUrl: `https://www.youtube.com/watch?v=${broadcastId}`,
    };
  }

  async stopLiveStream(): Promise<void> {
    // Implementation for stopping live stream
  }

  async getAnalytics(videoId: string): Promise<any> {
    const response = await this.youtube.videos.list({
      part: ['statistics', 'contentDetails'],
      id: [videoId],
    });

    return response.data.items?.[0];
  }

  async updateMetadata(videoId: string, metadata: any): Promise<void> {
    await this.youtube.videos.update({
      part: ['snippet'],
      requestBody: {
        id: videoId,
        snippet: {
          title: metadata.title,
          description: metadata.description,
          tags: metadata.tags,
        },
      },
    });
  }
}
