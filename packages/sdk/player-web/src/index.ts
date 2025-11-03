/**
 * StreamVerse Web Video Player SDK
 * Issue #28: Video Player SDK Development
 */

import Hls from 'hls.js';
import dashjs from 'dashjs';
import shaka from 'shaka-player';

export interface PlayerConfig {
  container: string | HTMLElement;
  src: string; // HLS or DASH manifest URL
  drm?: DRMConfig;
  streaming?: StreamingConfig;
  captions?: CaptionsConfig;
  logging?: LoggingConfig;
}

export interface DRMConfig {
  widevine?: {
    licenseUrl: string;
    serverCertificate?: string;
  };
  fairplay?: {
    licenseUrl: string;
    certificateUrl: string;
    processSpcUrl: string;
  };
  playready?: {
    licenseUrl: string;
  };
}

export interface StreamingConfig {
  bufferingGoal?: number; // seconds of buffer (default: 8)
  rebufferGoal?: number; // seconds before rebuffering (default: 2)
  abrEnabled?: boolean; // adaptive bitrate (default: true)
  maxBitrate?: number; // max bitrate in bps
  minBitrate?: number; // min bitrate in bps
}

export interface CaptionsConfig {
  enabled?: boolean;
  defaultLanguage?: string;
  availableLanguages?: string[];
}

export interface LoggingConfig {
  enabled?: boolean;
  analyticsUrl?: string; // Analytics Service endpoint
  userId?: string;
  contentId?: string;
}

export interface Quality {
  id: string;
  width: number;
  height: number;
  bitrate: number;
}

export type PlayerEvent = 'play' | 'pause' | 'seek' | 'ended' | 'error' | 'qualitychange' | 'buffering' | 'ready';

export class StreamVersePlayer {
  private videoElement: HTMLVideoElement;
  private hls?: Hls;
  private dashPlayer?: dashjs.MediaPlayer;
  private shakaPlayer?: shaka.Player;
  private config: PlayerConfig;
  private eventListeners: Map<PlayerEvent, Array<(data?: any) => void>>;
  private analyticsUrl?: string;
  private userId?: string;
  private contentId?: string;
  private qoeMetrics: {
    startupTime?: number;
    rebufferingEvents: number;
    rebufferingDuration: number;
    qualityChanges: number;
    errors: number;
  };

  constructor(config: PlayerConfig) {
    this.config = config;
    this.eventListeners = new Map();
    this.qoeMetrics = {
      rebufferingEvents: 0,
      rebufferingDuration: 0,
      qualityChanges: 0,
      errors: 0,
    };

    // Initialize analytics
    if (config.logging?.enabled) {
      this.analyticsUrl = config.logging.analyticsUrl || 'http://localhost:8080/analytics/events';
      this.userId = config.logging.userId;
      this.contentId = config.logging.contentId;
    }

    // Create video element
    const container = typeof config.container === 'string' 
      ? document.querySelector(config.container) 
      : config.container;
    
    if (!container) {
      throw new Error('Container element not found');
    }

    this.videoElement = document.createElement('video');
    this.videoElement.controls = true;
    container.appendChild(this.videoElement);

    // Initialize player based on manifest type
    this.initializePlayer();

    // Setup event listeners
    this.setupEventListeners();
  }

  private initializePlayer(): void {
    const src = this.config.src;
    const isHLS = src.endsWith('.m3u8') || src.includes('.m3u8');
    const isDASH = src.endsWith('.mpd') || src.includes('.mpd');

    if (isHLS) {
      this.initializeHLS();
    } else if (isDASH) {
      this.initializeDASH();
    } else {
      // Fallback: try HLS
      this.initializeHLS();
    }
  }

  private initializeHLS(): void {
    if (Hls.isSupported()) {
      this.hls = new Hls({
        enableWorker: true,
        lowLatencyMode: true,
        backBufferLength: 90,
        maxBufferLength: 30,
        maxBufferSize: 60 * 1000 * 1000, // 60MB
        maxMaxBufferLength: 600,
        startLevel: -1, // Auto
        capLevelToPlayerSize: true,
        abrEwmaDefaultEstimate: 500000, // 500kbps default
      });

      // DRM configuration
      if (this.config.drm?.widevine) {
        this.hls.config.widevineLicenseUrl = this.config.drm.widevine.licenseUrl;
      }

      // Quality selection
      this.hls.on(Hls.Events.MANIFEST_PARSED, () => {
        const qualities = this.hls!.levels.map((level, index) => ({
          id: `level_${index}`,
          width: level.width || 0,
          height: level.height || 0,
          bitrate: level.bitrate || 0,
        }));
        this.emit('qualitychange', { qualities, current: this.hls!.currentLevel });
      });

      this.hls.loadSource(this.config.src);
      this.hls.attachMedia(this.videoElement);
    } else if (this.videoElement.canPlayType('application/vnd.apple.mpegurl')) {
      // Native HLS support (Safari)
      this.videoElement.src = this.config.src;
    } else {
      throw new Error('HLS not supported');
    }
  }

  private initializeDASH(): void {
    if (shaka.Player.isBrowserSupported()) {
      this.shakaPlayer = new shaka.Player(this.videoElement);

      // DRM configuration
      if (this.config.drm) {
        const drmConfig: any = {};
        if (this.config.drm.widevine) {
          drmConfig.widevine = {
            licenseServerUrl: this.config.drm.widevine.licenseUrl,
          };
        }
        if (this.config.drm.playready) {
          drmConfig.playReady = {
            licenseServerUrl: this.config.drm.playready.licenseUrl,
          };
        }
        this.shakaPlayer.configure({
          drm: drmConfig,
        });
      }

      // Streaming configuration
      if (this.config.streaming) {
        this.shakaPlayer.configure({
          streaming: {
            bufferingGoal: this.config.streaming.bufferingGoal || 8,
            rebufferingGoal: this.config.streaming.rebufferGoal || 2,
            bufferBehind: 30,
          },
          abr: {
            enabled: this.config.streaming.abrEnabled !== false,
            restrictions: {
              maxWidth: this.config.streaming.maxBitrate ? undefined : 1920,
              maxHeight: this.config.streaming.maxBitrate ? undefined : 1080,
            },
          },
        });
      }

      this.shakaPlayer.load(this.config.src).then(() => {
        this.emit('ready');
      }).catch((error) => {
        this.emit('error', error);
      });
    } else {
      // Fallback to dash.js
      this.dashPlayer = dashjs.MediaPlayer().create();
      this.dashPlayer.initialize(this.videoElement, this.config.src, true);
    }
  }

  private setupEventListeners(): void {
    const video = this.videoElement;
    const startTime = Date.now();

    video.addEventListener('play', () => {
      this.emit('play');
      if (this.qoeMetrics.startupTime === undefined) {
        this.qoeMetrics.startupTime = Date.now() - startTime;
        this.submitQoE('play');
      }
    });

    video.addEventListener('pause', () => {
      this.emit('pause');
      this.submitQoE('pause');
    });

    video.addEventListener('seeked', () => {
      this.emit('seek', { currentTime: video.currentTime });
      this.submitQoE('seek');
    });

    video.addEventListener('ended', () => {
      this.emit('ended');
      this.submitQoE('ended');
    });

    video.addEventListener('waiting', () => {
      const startRebuffer = Date.now();
      this.qoeMetrics.rebufferingEvents++;
      this.emit('buffering', { start: true });
      
      video.addEventListener('playing', () => {
        const rebufferDuration = (Date.now() - startRebuffer) / 1000;
        this.qoeMetrics.rebufferingDuration += rebufferDuration;
        this.emit('buffering', { start: false, duration: rebufferDuration });
        this.submitQoE('buffering', { duration: rebufferDuration });
      }, { once: true });
    });

    video.addEventListener('error', (e) => {
      this.qoeMetrics.errors++;
      this.emit('error', { error: video.error });
      this.submitQoE('error', { error: video.error?.message });
    });
  }

  // Public API
  play(): void {
    this.videoElement.play();
  }

  pause(): void {
    this.videoElement.pause();
  }

  seek(time: number): void {
    this.videoElement.currentTime = time;
  }

  setQuality(qualityId: string): void {
    if (this.hls) {
      const level = parseInt(qualityId.split('_')[1]);
      if (!isNaN(level)) {
        this.hls.currentLevel = level;
        this.qoeMetrics.qualityChanges++;
        this.emit('qualitychange', { quality: qualityId });
      }
    } else if (this.shakaPlayer) {
      const variants = this.shakaPlayer.getVariantTracks();
      const variant = variants.find((v: any) => v.id === qualityId);
      if (variant) {
        this.shakaPlayer.selectVariantTrack(variant, true);
        this.qoeMetrics.qualityChanges++;
        this.emit('qualitychange', { quality: qualityId });
      }
    }
  }

  getQualities(): Quality[] {
    if (this.hls) {
      return this.hls.levels.map((level, index) => ({
        id: `level_${index}`,
        width: level.width || 0,
        height: level.height || 0,
        bitrate: level.bitrate || 0,
      }));
    }
    return [];
  }

  on(event: PlayerEvent, callback: (data?: any) => void): void {
    if (!this.eventListeners.has(event)) {
      this.eventListeners.set(event, []);
    }
    this.eventListeners.get(event)!.push(callback);
  }

  off(event: PlayerEvent, callback: (data?: any) => void): void {
    const listeners = this.eventListeners.get(event);
    if (listeners) {
      const index = listeners.indexOf(callback);
      if (index > -1) {
        listeners.splice(index, 1);
      }
    }
  }

  destroy(): void {
    if (this.hls) {
      this.hls.destroy();
    }
    if (this.shakaPlayer) {
      this.shakaPlayer.destroy();
    }
    if (this.dashPlayer) {
      this.dashPlayer.destroy();
    }
    this.videoElement.remove();
  }

  private emit(event: PlayerEvent, data?: any): void {
    const listeners = this.eventListeners.get(event);
    if (listeners) {
      listeners.forEach(callback => callback(data));
    }
  }

  private submitQoE(event: string, metadata?: any): void {
    if (!this.analyticsUrl || !this.userId || !this.contentId) {
      return;
    }

    const payload = {
      user_id: this.userId,
      content_id: this.contentId,
      event_type: event,
      metadata: {
        ...metadata,
        startup_time: this.qoeMetrics.startupTime,
        rebuffering_events: this.qoeMetrics.rebufferingEvents,
        rebuffering_duration: this.qoeMetrics.rebufferingDuration,
        quality_changes: this.qoeMetrics.qualityChanges,
        errors: this.qoeMetrics.errors,
      },
      timestamp: new Date().toISOString(),
    };

    // Send to Analytics Service (async, non-blocking)
    fetch(this.analyticsUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    }).catch(() => {
      // Silently fail - don't block playback
    });
  }
}

// Export default
export default StreamVersePlayer;

