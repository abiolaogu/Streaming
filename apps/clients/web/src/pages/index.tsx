import { useEffect, useRef, useState } from 'react';
import Head from 'next/head';
import shaka from 'shaka-player/dist/shaka-player.ui.js';

export default function Home() {
  const videoRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<shaka.Player | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);

  useEffect(() => {
    if (videoRef.current && !playerRef.current) {
      const player = new shaka.Player(videoRef.current);
      player.configure({
        streaming: {
          retryParameters: {
            timeout: 30000,
            maxAttempts: 3,
            baseDelay: 1000,
            backoffFactor: 2,
            fuzzFactor: 0.5
          },
          bufferingGoal: 30,
          rebufferingGoal: 5,
          bufferBehind: 30
        },
        manifest: {
          retryParameters: {
            timeout: 30000,
            maxAttempts: 3,
            baseDelay: 1000,
            backoffFactor: 2
          },
          availabilityWindowOverride: 60
        }
      });

      playerRef.current = player;

      player.addEventListener('error', (event) => {
        console.error('Shaka error:', event.detail);
      });

      player.addEventListener('adaptation', () => {
        console.log('Adaptation changed');
      });

      return () => {
        player.destroy();
        playerRef.current = null;
      };
    }
  }, []);

  const loadVideo = async (url: string) => {
    if (playerRef.current) {
      try {
        await playerRef.current.load(url);
      } catch (error) {
        console.error('Load error:', error);
      }
    }
  };

  const handlePlay = () => {
    if (videoRef.current) {
      if (isPlaying) {
        videoRef.current.pause();
      } else {
        videoRef.current.play();
      }
      setIsPlaying(!isPlaying);
    }
  };

  const handleTimeUpdate = () => {
    if (videoRef.current) {
      setCurrentTime(videoRef.current.currentTime);
    }
  };

  return (
    <div className="container">
      <Head>
        <title>Streaming Platform</title>
        <meta name="description" content="Video streaming platform" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <h1>Streaming Platform</h1>
        
        <div className="player-container">
          <video
            ref={videoRef}
            width="960"
            height="540"
            onTimeUpdate={handleTimeUpdate}
            controls
          />
        </div>

        <div className="controls">
          <button onClick={handlePlay}>
            {isPlaying ? 'Pause' : 'Play'}
          </button>
          
          <button onClick={() => loadVideo('/test/manifest.mpd')}>
            Load Test Video
          </button>
          
          <button onClick={() => loadVideo('/fast/channel-1.m3u8')}>
            Load FAST Channel
          </button>
        </div>

        <div className="info">
          <p>Time: {Math.floor(currentTime)}s</p>
        </div>
      </main>

      <style jsx>{`
        .container {
          min-height: 100vh;
          padding: 0 0.5rem;
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
        }

        main {
          padding: 5rem 0;
          flex: 1;
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
        }

        h1 {
          margin: 0;
          line-height: 1.15;
          font-size: 4rem;
        }

        .player-container {
          margin: 2rem 0;
        }

        .controls {
          display: flex;
          gap: 1rem;
          margin: 1rem 0;
        }

        button {
          padding: 0.5rem 1rem;
          font-size: 1rem;
          cursor: pointer;
        }

        .info {
          margin-top: 2rem;
        }
      `}</style>
    </div>
  );
}

