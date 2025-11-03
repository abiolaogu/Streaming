import React, { useState, useEffect, useRef } from 'react';
import { fetchContentDetails } from '../services/mockApiService';
import { AnyContent, LiveChannel, MediaContent } from '../types';
import LoadingSpinner from '../components/LoadingSpinner';
import { PlusCircleIcon, ShareIcon, PlayIcon, PauseIcon, SpeakerWaveIcon, SpeakerXMarkIcon, ArrowsPointingOutIcon, ArrowsPointingInIcon } from '../components/IconComponents';

interface WatchViewProps {
    contentId: string;
}

const formatTime = (timeInSeconds: number) => {
    const minutes = Math.floor(timeInSeconds / 60);
    const seconds = Math.floor(timeInSeconds % 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
};

const WatchView: React.FC<WatchViewProps> = ({ contentId }) => {
    const [content, setContent] = useState<AnyContent | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    // Player State
    const [isPlaying, setIsPlaying] = useState(false);
    const [progress, setProgress] = useState(0); // 0-100
    const [volume, setVolume] = useState(0.75); // 0-1
    const [isMuted, setIsMuted] = useState(false);
    const [duration] = useState(300); // Mock duration: 5 minutes
    const [isFullscreen, setIsFullscreen] = useState(false);
    const [showControls, setShowControls] = useState(true);
    const playerRef = useRef<HTMLDivElement>(null);
    const controlsTimeoutRef = useRef<number | null>(null);


    useEffect(() => {
        const loadContent = async () => {
            if (!contentId) return;
            setIsLoading(true);
            const details = await fetchContentDetails(contentId);
            setContent(details || null);
            setIsLoading(false);
        };
        loadContent();
    }, [contentId]);

    // Mock video progress
    useEffect(() => {
        let interval: number;
        if (isPlaying) {
            interval = window.setInterval(() => {
                setProgress(p => {
                    const newProgress = p + (100 / duration);
                    if (newProgress >= 100) {
                        setIsPlaying(false);
                        return 100;
                    }
                    return newProgress;
                });
            }, 1000);
        }
        return () => clearInterval(interval);
    }, [isPlaying, duration]);

    // Handlers
    const togglePlay = () => setIsPlaying(!isPlaying);
    const handleSeek = (e: React.ChangeEvent<HTMLInputElement>) => setProgress(Number(e.target.value));
    const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setVolume(Number(e.target.value));
        if (Number(e.target.value) > 0) setIsMuted(false);
    };
    const toggleMute = () => {
        setIsMuted(!isMuted);
        if (isMuted) setVolume(0.75); else setVolume(0);
    };

    const toggleFullscreen = () => {
        if (!playerRef.current) return;
        if (!document.fullscreenElement) {
            playerRef.current.requestFullscreen();
            setIsFullscreen(true);
        } else {
            document.exitFullscreen();
            setIsFullscreen(false);
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>, step: number) => {
        const target = e.target as HTMLInputElement;
        let value = Number(target.value);
        if (e.key === 'ArrowLeft' || e.key === 'ArrowDown') {
            value = Math.max(0, value - step);
        } else if (e.key === 'ArrowRight' || e.key === 'ArrowUp') {
            value = Math.min(100, value + step);
        }
        
        if (target.id === 'seek-slider') {
             setProgress(value);
        } else if (target.id === 'volume-slider') {
            setVolume(value / 100);
            if (value > 0) setIsMuted(false);
        }
    };
    
    const hideControls = () => {
        if(controlsTimeoutRef.current) clearTimeout(controlsTimeoutRef.current);
        controlsTimeoutRef.current = window.setTimeout(() => setShowControls(false), 3000);
    };

    const handlePlayerInteraction = () => {
        setShowControls(true);
        hideControls();
    };
    
    useEffect(() => {
      hideControls();
      return () => {
         if(controlsTimeoutRef.current) clearTimeout(controlsTimeoutRef.current);
      }
    },[]);

    if (isLoading) {
        return <div className="flex justify-center items-center h-full"><LoadingSpinner /></div>;
    }

    if (!content) {
        return <div className="text-center text-brand-text-secondary">Content not found.</div>;
    }
    
    const isLive = content.type === 'Live';
    const title = isLive ? (content as LiveChannel).name : (content as MediaContent).title;
    const description = content.description;
    const imageUrl = isLive ? (content as LiveChannel).logoUrl.replace('/100/50', '/800/450') : (content as MediaContent).thumbnailUrl.replace('/400/225', '/800/450');
    const currentTime = (progress / 100) * duration;
    const currentTimeText = formatTime(currentTime);
    const durationText = formatTime(duration);

    return (
        <div className="max-w-6xl mx-auto">
            <div 
                ref={playerRef} 
                className="relative group aspect-video bg-black rounded-lg border border-brand-border flex items-center justify-center mb-6 shadow-2xl overflow-hidden"
                onMouseMove={handlePlayerInteraction}
                onMouseLeave={() => { if(isPlaying) hideControls() }}
                onFocus={() => setShowControls(true)}
                onBlur={() => setShowControls(false)}
            >
                <img src={imageUrl} alt={title} className="w-full h-full object-cover" />
                
                {/* Controls Overlay */}
                <div className={`absolute inset-0 bg-black/40 transition-opacity duration-300 ${showControls ? 'opacity-100' : 'opacity-0'}`}>
                    <div 
                        role="toolbar" 
                        aria-label="Video player controls" 
                        className="absolute bottom-0 left-0 right-0 p-3 text-white flex flex-col gap-2 bg-gradient-to-t from-black/80 to-transparent"
                    >
                        {/* Seek Bar */}
                        <div className="relative w-full">
                           <input
                                id="seek-slider"
                                type="range"
                                min="0"
                                max="100"
                                step="0.1"
                                value={progress}
                                onChange={handleSeek}
                                onKeyDown={(e) => handleKeyDown(e, 5)}
                                className="w-full h-1 bg-white/30 rounded-lg appearance-none cursor-pointer accent-brand-accent"
                                aria-label="Video progress"
                                aria-valuemin={0}
                                aria-valuemax={100}
                                aria-valuenow={progress}
                                aria-valuetext={`${currentTimeText} of ${durationText}`}
                            />
                        </div>

                        {/* Bottom Controls */}
                        <div className="flex justify-between items-center">
                            <div className="flex items-center gap-4">
                                <button onClick={togglePlay} aria-label={isPlaying ? 'Pause' : 'Play'}>
                                    {isPlaying ? <PauseIcon className="w-7 h-7" /> : <PlayIcon className="w-7 h-7" />}
                                </button>
                                <div className="flex items-center gap-2 group/volume">
                                     <button onClick={toggleMute} aria-label={isMuted || volume === 0 ? 'Unmute' : 'Mute'}>
                                        {isMuted || volume === 0 ? <SpeakerXMarkIcon className="w-6 h-6" /> : <SpeakerWaveIcon className="w-6 h-6" />}
                                    </button>
                                    <input
                                        id="volume-slider"
                                        type="range"
                                        min="0"
                                        max="1"
                                        step="0.05"
                                        value={isMuted ? 0 : volume}
                                        onChange={handleVolumeChange}
                                        onKeyDown={(e) => handleKeyDown(e, 5)}
                                        className="w-20 h-1 bg-white/30 rounded-lg appearance-none cursor-pointer accent-brand-accent transition-all duration-300 w-0 group-hover/volume:w-20"
                                        aria-label="Volume"
                                        aria-valuemin={0}
                                        aria-valuemax={100}
                                        aria-valuenow={isMuted ? 0 : volume * 100}
                                        aria-valuetext={`Volume at ${isMuted ? 0 : Math.round(volume * 100)}%`}
                                    />
                                </div>
                            </div>
                            <div className="flex items-center gap-4">
                                <span className="text-xs font-mono">{currentTimeText} / {durationText}</span>
                                <button onClick={toggleFullscreen} aria-label={isFullscreen ? 'Exit fullscreen' : 'Enter fullscreen'}>
                                    {isFullscreen ? <ArrowsPointingInIcon className="w-6 h-6" /> : <ArrowsPointingOutIcon className="w-6 h-6" />}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="md:col-span-2">
                    <h1 className="text-3xl font-bold text-brand-text-primary mb-2">{title}</h1>
                    <div className="flex items-center space-x-4 text-sm text-brand-text-secondary mb-4">
                        <span>{content.type}</span>
                        <span>&bull;</span>
                        <span>{content.category}</span>
                        {!isLive && <span>&bull;</span>}
                        {!isLive && <span>{ (content as MediaContent).monetizationModel}</span>}
                    </div>
                    <p className="text-brand-text-secondary leading-relaxed">{description}</p>
                </div>
                <div className="md:col-span-1">
                    <div className="space-y-3">
                         <button className="w-full flex justify-center items-center gap-2 bg-brand-accent hover:bg-brand-accent-hover text-white font-bold py-3 px-4 rounded-md transition">
                           <PlayIcon className="h-6 w-6" /> Play
                        </button>
                         <button className="w-full flex justify-center items-center gap-2 bg-brand-surface hover:bg-brand-bg text-brand-text-primary font-bold py-3 px-4 rounded-md border border-brand-border transition">
                           <PlusCircleIcon className="h-6 w-6" /> Add to Watchlist
                        </button>
                         <button className="w-full flex justify-center items-center gap-2 bg-brand-surface hover:bg-brand-bg text-brand-text-primary font-bold py-3 px-4 rounded-md border border-brand-border transition">
                           <ShareIcon className="h-6 w-6" /> Share
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default WatchView;