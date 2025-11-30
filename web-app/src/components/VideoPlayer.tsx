/**
 * StreamVerse Video Player with Smallpixel Integration
 * Features:
 * - HLS/DASH playback
 * - Client-side AI upscaling (Smallpixel SDK)
 * - Adaptive bitrate streaming
 * - DRM support
 * - Analytics tracking
 */

import React, { useEffect, useRef, useState } from 'react';
import Hls from 'hls.js';
import dashjs from 'dashjs';
import SmallpixelSDK, { SmallpixelConfig, UpscalingStats } from '../utils/smallpixel-sdk';

interface VideoPlayerProps {
  src: string;
  type: 'hls' | 'dash' | 'mp4';
  poster?: string;
  autoplay?: boolean;
  controls?: boolean;
  enableSmallpixel?: boolean;
  smallpixelConfig?: Partial<SmallpixelConfig>;
  onProgress?: (progress: number) => void;
  onEnded?: () => void;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({
  src,
  type,
  poster,
  autoplay = false,
  controls = true,
  enableSmallpixel = true,
  smallpixelConfig,
  onProgress,
  onEnded,
}) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [hls, setHls] = useState<Hls | null>(null);
  const [dash, setDash] = useState<any>(null);
  const [smallpixel, setSmallpixel] = useState<SmallpixelSDK | null>(null);
  const [upscalingStats, setUpscalingStats] = useState<UpscalingStats | null>(null);
  const [bandwidthSaved, setBandwidthSaved] = useState<number>(0);

  useEffect(() => {
    if (!videoRef.current) return;

    const video = videoRef.current;

    // Initialize video player based on type
    if (type === 'hls' && Hls.isSupported()) {
      const hlsInstance = new Hls({
        enableWorker: true,
        lowLatencyMode: true,
        backBufferLength: 90,
      });

      hlsInstance.loadSource(src);
      hlsInstance.attachMedia(video);

      hlsInstance.on(Hls.Events.MANIFEST_PARSED, () => {
        if (autoplay) {
          video.play();
        }
      });

      setHls(hlsInstance);
    } else if (type === 'dash') {
      const dashPlayer = dashjs.MediaPlayer().create();
      dashPlayer.initialize(video, src, autoplay);
      setDash(dashPlayer);
    } else {
      video.src = src;
      if (autoplay) {
        video.play();
      }
    }

    // Initialize Smallpixel SDK
    if (enableSmallpixel) {
      initializeSmallpixel(video);
    }

    // Track progress
    const handleTimeUpdate = () => {
      if (onProgress && video.duration) {
        const progress = (video.currentTime / video.duration) * 100;
        onProgress(progress);
      }
    };

    video.addEventListener('timeupdate', handleTimeUpdate);
    video.addEventListener('ended', () => {
      if (onEnded) onEnded();
    });

    return () => {
      video.removeEventListener('timeupdate', handleTimeUpdate);
      if (hls) hls.destroy();
      if (dash) dash.reset();
      if (smallpixel) smallpixel.destroy();
    };
  }, [src, type, enableSmallpixel]);

  /**
   * Initialize Smallpixel SDK for bandwidth savings
   */
  const initializeSmallpixel = async (video: HTMLVideoElement) => {
    const defaultConfig: SmallpixelConfig = {
      apiKey: process.env.REACT_APP_SMALLPIXEL_API_KEY || '',
      targetResolution: 'auto',
      enableUpscaling: true,
      quality: 'high',
      gpuAcceleration: true,
      bandwidthSavings: true,
      ...smallpixelConfig,
    };

    const sdk = new SmallpixelSDK(defaultConfig);
    await sdk.initialize(video);
    sdk.startUpscaling();

    setSmallpixel(sdk);

    // Update stats every 5 seconds
    const statsInterval = setInterval(() => {
      const stats = sdk.getStats();
      setUpscalingStats(stats);
      setBandwidthSaved((prev) => prev + stats.bandwidthSaved / 12); // per 5 seconds
    }, 5000);

    return () => clearInterval(statsInterval);
  };

  /**
   * Toggle Smallpixel upscaling
   */
  const toggleSmallpixel = () => {
    if (smallpixel) {
      smallpixel.destroy();
      setSmallpixel(null);
      setUpscalingStats(null);
    } else if (videoRef.current) {
      initializeSmallpixel(videoRef.current);
    }
  };

  return (
    <div ref={containerRef} className="video-player-container" style={{ position: 'relative' }}>
      <video
        ref={videoRef}
        className="video-player"
        poster={poster}
        controls={controls}
        playsInline
        style={{ width: '100%', height: '100%' }}
      />

      {/* Smallpixel Stats Overlay */}
      {upscalingStats && (
        <div
          className="smallpixel-stats"
          style={{
            position: 'absolute',
            top: '10px',
            right: '10px',
            background: 'rgba(0, 0, 0, 0.7)',
            color: 'white',
            padding: '10px',
            borderRadius: '8px',
            fontSize: '12px',
            fontFamily: 'monospace',
          }}
        >
          <div style={{ fontWeight: 'bold', marginBottom: '5px' }}>
            🔺 Smallpixel Active
          </div>
          <div>Source: {upscalingStats.originalResolution}</div>
          <div>Target: {upscalingStats.targetResolution}</div>
          <div>Latency: {upscalingStats.upscalingLatency.toFixed(1)}ms</div>
          <div>Saved: {bandwidthSaved.toFixed(2)} MB</div>
          <div>FPS: {upscalingStats.frameRate}</div>
          <button
            onClick={toggleSmallpixel}
            style={{
              marginTop: '8px',
              padding: '4px 8px',
              background: '#ff6b6b',
              border: 'none',
              borderRadius: '4px',
              color: 'white',
              cursor: 'pointer',
              fontSize: '11px',
            }}
          >
            Disable
          </button>
        </div>
      )}

      {/* Enable Smallpixel Button (if disabled) */}
      {!smallpixel && enableSmallpixel && (
        <button
          onClick={toggleSmallpixel}
          style={{
            position: 'absolute',
            top: '10px',
            right: '10px',
            padding: '8px 16px',
            background: '#4ecdc4',
            border: 'none',
            borderRadius: '8px',
            color: 'white',
            cursor: 'pointer',
            fontWeight: 'bold',
          }}
        >
          Enable Smallpixel (Save 60% Bandwidth)
        </button>
      )}
    </div>
  );
};

export default VideoPlayer;
