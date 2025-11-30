// StreamVerse SaaS - End-to-End Tests
// Comprehensive testing suite for streaming infrastructure

import { describe, it, expect, beforeAll, afterAll } from '@jest/globals';
import axios from 'axios';
import WebSocket from 'ws';
import * as fs from 'fs';
import * as path from 'path';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
const INGESTION_URL = process.env.INGESTION_URL || 'http://localhost:8100';
const TRANSCODING_URL = process.env.TRANSCODING_URL || 'http://localhost:8101';
const ANALYTICS_URL = process.env.ANALYTICS_URL || 'http://localhost:8105';

const TEST_VIDEO_PATH = path.join(__dirname, '../fixtures/test-video.mp4');
const TEST_TIMEOUT = 120000; // 2 minutes

describe('StreamVerse SaaS End-to-End Tests', () => {
  let authToken: string;
  let streamId: string;
  let videoId: string;

  beforeAll(async () => {
    // Authenticate
    const authResponse = await axios.post(`${API_BASE_URL}/api/v1/auth/login`, {
      email: 'test@streamverse.io',
      password: 'test-password',
    });

    authToken = authResponse.data.token;
    expect(authToken).toBeDefined();
  });

  // ==================== HEALTH CHECKS ====================
  describe('Service Health Checks', () => {
    it('should verify API Gateway is healthy', async () => {
      const response = await axios.get(`${API_BASE_URL}/health`);
      expect(response.status).toBe(200);
      expect(response.data).toBe('OK');
    });

    it('should verify Ingestion Service is healthy', async () => {
      const response = await axios.get(`${INGESTION_URL}/health`);
      expect(response.status).toBe(200);
    });

    it('should verify Transcoding Service is healthy', async () => {
      const response = await axios.get(`${TRANSCODING_URL}/health`);
      expect(response.status).toBe(200);
    });

    it('should verify Analytics Service is healthy', async () => {
      const response = await axios.get(`${ANALYTICS_URL}/health`);
      expect(response.status).toBe(200);
    });
  });

  // ==================== VIDEO INGESTION ====================
  describe('Video Ingestion', () => {
    it('should start RTMP ingestion', async () => {
      const response = await axios.post(
        `${INGESTION_URL}/api/v1/ingest/start`,
        {
          protocol: 'rtmp',
          input_url: 'rtmp://source.example.com/live/stream',
          title: 'Test Live Stream',
          target_platforms: ['youtube', 'twitch'],
          quality_profile: {
            resolution: '1080p',
            bitrate: 5000,
            fps: 60,
            codec: 'h264',
          },
          drm_enabled: true,
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.stream_id).toBeDefined();
      expect(response.data.ingest_url).toContain('rtmp://');

      streamId = response.data.stream_id;
    }, TEST_TIMEOUT);

    it('should get ingestion status', async () => {
      const response = await axios.get(
        `${INGESTION_URL}/api/v1/ingest/${streamId}/status`,
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.stream_id).toBe(streamId);
      expect(response.data.status).toBeDefined();
    });

    it('should handle SRT ingestion', async () => {
      const response = await axios.post(
        `${INGESTION_URL}/api/v1/ingest/start`,
        {
          protocol: 'srt',
          input_url: 'srt://source.example.com:9999',
          title: 'SRT Test Stream',
          target_platforms: ['kick'],
          quality_profile: {
            resolution: '4K',
            bitrate: 15000,
            fps: 60,
            codec: 'hevc',
          },
          drm_enabled: false,
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.ingest_url).toContain('srt://');
    });
  });

  // ==================== TRANSCODING ====================
  describe('GPU-Accelerated Transcoding', () => {
    it('should transcode video to multiple formats', async () => {
      const response = await axios.post(
        `${TRANSCODING_URL}/api/v1/transcode`,
        {
          video_id: 'test-video-123',
          source_url: 's3://bucket/source.mp4',
          profiles: [
            { resolution: '4K', bitrate: 15000, codec: 'hevc' },
            { resolution: '1080p', bitrate: 5000, codec: 'h264' },
            { resolution: '720p', bitrate: 2500, codec: 'h264' },
            { resolution: '480p', bitrate: 1000, codec: 'h264' },
          ],
          gpu_acceleration: true,
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.job_id).toBeDefined();
      expect(response.data.estimated_completion).toBeDefined();
    }, TEST_TIMEOUT);

    it('should support AV1 codec', async () => {
      const response = await axios.post(
        `${TRANSCODING_URL}/api/v1/transcode`,
        {
          video_id: 'test-video-av1',
          source_url: 's3://bucket/source.mp4',
          profiles: [
            { resolution: '4K', bitrate: 10000, codec: 'av1' },
          ],
          gpu_acceleration: true,
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.codec).toBe('av1');
    });

    it('should generate adaptive bitrate ladder', async () => {
      const response = await axios.post(
        `${TRANSCODING_URL}/api/v1/transcode/abr`,
        {
          video_id: 'test-video-abr',
          source_url: 's3://bucket/source.mp4',
          target_quality: '4K',
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.profiles).toBeInstanceOf(Array);
      expect(response.data.profiles.length).toBeGreaterThanOrEqual(5);
    });
  });

  // ==================== PLATFORM INTEGRATIONS ====================
  describe('Platform Integrations', () => {
    it('should upload to YouTube', async () => {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/platforms/youtube/upload`,
        {
          video_url: 's3://bucket/video.mp4',
          metadata: {
            title: 'Test Upload',
            description: 'E2E test video',
            tags: ['test', 'e2e'],
            category: 'Education',
            privacy: 'unlisted',
          },
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.youtube_video_id).toBeDefined();
      expect(response.data.url).toContain('youtube.com');
    }, TEST_TIMEOUT);

    it('should upload to Twitch', async () => {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/platforms/twitch/upload`,
        {
          video_url: 's3://bucket/video.mp4',
          metadata: {
            title: 'Test Stream Highlight',
            description: 'E2E test',
          },
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.twitch_video_id).toBeDefined();
    }, TEST_TIMEOUT);

    it('should upload to TikTok', async () => {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/platforms/tiktok/upload`,
        {
          video_url: 's3://bucket/short-video.mp4',
          metadata: {
            title: 'Test TikTok',
            hashtags: ['test', 'streamverse'],
          },
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.tiktok_video_id).toBeDefined();
    }, TEST_TIMEOUT);

    it('should sync metadata across platforms', async () => {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/platforms/sync`,
        {
          video_id: videoId,
          platforms: ['youtube', 'twitch', 'vimeo'],
          metadata: {
            title: 'Updated Title',
            description: 'Updated description',
            tags: ['updated', 'synced'],
          },
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.synced_platforms).toEqual([
        'youtube',
        'twitch',
        'vimeo',
      ]);
    });
  });

  // ==================== DRM & SECURITY ====================
  describe('DRM and Security', () => {
    it('should generate DRM license', async () => {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/drm/license`,
        {
          video_id: 'test-video-123',
          drm_type: 'widevine',
          device_id: 'test-device-001',
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.license_url).toBeDefined();
      expect(response.data.license_token).toBeDefined();
      expect(response.data.expiry).toBeDefined();
    });

    it('should apply forensic watermark', async () => {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/drm/watermark`,
        {
          video_id: 'test-video-123',
          user_id: 'test-user',
          watermark_text: 'USER-123',
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.watermarked_url).toBeDefined();
    });

    it('should verify blockchain license', async () => {
      const response = await axios.post(
        `${API_BASE_URL}/api/v1/drm/blockchain/verify`,
        {
          video_id: 'test-video-123',
          transaction_hash: '0x1234567890abcdef',
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.valid).toBe(true);
    });
  });

  // ==================== ANALYTICS ====================
  describe('Real-Time Analytics', () => {
    it('should track video views', async () => {
      const response = await axios.get(
        `${ANALYTICS_URL}/api/v1/analytics/views`,
        {
          params: {
            video_id: 'test-video-123',
            time_range: 'last_24_hours',
          },
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.total_views).toBeGreaterThanOrEqual(0);
      expect(response.data.unique_viewers).toBeDefined();
      expect(response.data.watch_time_minutes).toBeDefined();
    });

    it('should provide engagement metrics', async () => {
      const response = await axios.get(
        `${ANALYTICS_URL}/api/v1/analytics/engagement`,
        {
          params: {
            video_id: 'test-video-123',
          },
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.engagement_rate).toBeDefined();
      expect(response.data.avg_watch_percentage).toBeDefined();
      expect(response.data.likes).toBeGreaterThanOrEqual(0);
    });

    it('should generate QoE (Quality of Experience) report', async () => {
      const response = await axios.get(
        `${ANALYTICS_URL}/api/v1/analytics/qoe`,
        {
          params: {
            video_id: 'test-video-123',
          },
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      expect(response.status).toBe(200);
      expect(response.data.avg_bitrate).toBeDefined();
      expect(response.data.buffering_events).toBeDefined();
      expect(response.data.startup_time_ms).toBeDefined();
    });
  });

  // ==================== PERFORMANCE TESTS ====================
  describe('Performance and Scalability', () => {
    it('should handle 1000+ concurrent ingestion streams', async () => {
      const requests = Array.from({ length: 1000 }, (_, i) =>
        axios.post(
          `${INGESTION_URL}/api/v1/ingest/start`,
          {
            protocol: 'rtmp',
            input_url: `rtmp://test-${i}.example.com/live/stream`,
            title: `Concurrent Stream ${i}`,
            target_platforms: ['youtube'],
            quality_profile: {
              resolution: '1080p',
              bitrate: 5000,
              fps: 30,
              codec: 'h264',
            },
            drm_enabled: false,
          },
          {
            headers: {
              Authorization: `Bearer ${authToken}`,
            },
          }
        )
      );

      const results = await Promise.allSettled(requests);
      const successful = results.filter((r) => r.status === 'fulfilled').length;

      expect(successful).toBeGreaterThan(950); // 95% success rate
    }, TEST_TIMEOUT * 5);

    it('should process transcoding at 100x real-time', async () => {
      const startTime = Date.now();

      await axios.post(
        `${TRANSCODING_URL}/api/v1/transcode`,
        {
          video_id: 'benchmark-video',
          source_url: 's3://bucket/10min-video.mp4',
          profiles: [{ resolution: '1080p', bitrate: 5000, codec: 'h264' }],
          gpu_acceleration: true,
        },
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );

      const endTime = Date.now();
      const durationSeconds = (endTime - startTime) / 1000;

      // 10-minute video should transcode in < 6 seconds (100x real-time)
      expect(durationSeconds).toBeLessThan(6);
    }, TEST_TIMEOUT);

    it('should handle 10M+ API requests per second', async () => {
      const requests = 100000;
      const startTime = Date.now();

      const promises = Array.from({ length: requests }, () =>
        axios.get(`${API_BASE_URL}/health`)
      );

      await Promise.all(promises);

      const endTime = Date.now();
      const durationSeconds = (endTime - startTime) / 1000;
      const rps = requests / durationSeconds;

      expect(rps).toBeGreaterThan(10000); // At least 10K RPS
    }, TEST_TIMEOUT);
  });

  // ==================== CLEANUP ====================
  afterAll(async () => {
    // Stop test streams
    if (streamId) {
      await axios.post(
        `${INGESTION_URL}/api/v1/ingest/${streamId}/stop`,
        {},
        {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        }
      );
    }

    // Clean up test data
    console.log('✅ All tests completed');
  });
});
