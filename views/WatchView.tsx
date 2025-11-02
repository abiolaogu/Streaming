import React, { useState, useEffect } from 'react';
import { fetchContentDetails } from '../services/mockApiService';
import { AnyContent, LiveChannel, MediaContent } from '../types';
import LoadingSpinner from '../components/LoadingSpinner';
import { PlayCircleIcon, PlusCircleIcon, ShareIcon } from '@heroicons/react/24/solid';

interface WatchViewProps {
    contentId: string;
}

const WatchView: React.FC<WatchViewProps> = ({ contentId }) => {
    const [content, setContent] = useState<AnyContent | null>(null);
    const [isLoading, setIsLoading] = useState(true);

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

    return (
        <div className="max-w-6xl mx-auto">
            <div className="aspect-video bg-brand-surface rounded-lg border border-brand-border flex items-center justify-center mb-6 shadow-2xl overflow-hidden">
                {/* Video Player Placeholder */}
                <div className="w-full h-full bg-black flex items-center justify-center relative">
                    <img src={imageUrl} alt={title} className="w-full h-full object-cover opacity-30" />
                    <div className="absolute inset-0 flex items-center justify-center">
                        <PlayCircleIcon className="w-24 h-24 text-white/70 hover:text-white transition-colors cursor-pointer" />
                    </div>
                    <div className="absolute bottom-4 left-4 text-white text-xs">Player SDK (ISSUE-028) - Mock Player</div>
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
                           <PlayCircleIcon className="h-6 w-6" /> Play
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