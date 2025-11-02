import React, { useState, useEffect } from 'react';
import { ContentCategory, MediaContent, MediaContentWithProgress, LiveChannel, ContentCategoryString, AnyContent } from '../types';
import { PlayCircleIcon, LockClosedIcon, CurrencyDollarIcon } from '@heroicons/react/24/solid';
import { fetchVODContent, fetchLiveChannels, fetchContinueWatching } from '../services/mockApiService';
import LoadingSpinner from '../components/LoadingSpinner';

// --- COMPONENTS ---

const LiveProgramProgress: React.FC<{ startTime: number; endTime: number }> = ({ startTime, endTime }) => {
    const [progress, setProgress] = useState(0);
    useEffect(() => {
        const calculateProgress = () => {
            const now = Date.now();
            const duration = endTime - startTime;
            const elapsed = now - startTime;
            setProgress(Math.min(100, Math.max(0, (elapsed / duration) * 100)));
        };
        calculateProgress();
        const interval = setInterval(calculateProgress, 60000);
        return () => clearInterval(interval);
    }, [startTime, endTime]);
    return <div className="absolute bottom-0 left-0 h-1.5 w-full bg-gray-500/70"><div className="h-full bg-brand-danger" style={{ width: `${progress}%` }} /></div>;
};

const ContentCard: React.FC<{ item: AnyContent, onWatch: (id: string) => void }> = ({ item, onWatch }) => {
    const { type } = item;
    const isLive = type === 'Live';
    const isSvod = !isLive && (item as MediaContent).monetizationModel === 'SVOD';
    const isTvod = !isLive && (item as MediaContent).monetizationModel === 'TVOD';
    const isAvod = !isLive && (item as MediaContent).monetizationModel === 'AVOD';

    const renderOverlay = () => {
        if (isSvod) return <div className="text-center"><LockClosedIcon className="w-12 h-12 text-white/80 mx-auto"/><p className="text-xs mt-1">Subscribe to Watch</p></div>;
        if (isTvod) return <div className="text-center"><CurrencyDollarIcon className="w-12 h-12 text-white/80 mx-auto"/><p className="text-xs mt-1">Rent or Buy</p></div>;
        return <PlayCircleIcon className="w-16 h-16 text-white/80" />;
    };

    const imageUrl = isLive ? (item as LiveChannel).logoUrl.replace('/100/50', '/400/225') : (item as MediaContent).thumbnailUrl;
    const title = isLive ? (item as LiveChannel).name : (item as MediaContent).title;

    return (
        <div 
            onClick={() => onWatch(item.id)}
            className="flex-shrink-0 w-64 group cursor-pointer transform-gpu transition-transform hover:scale-105"
        >
            <div className="relative aspect-video bg-brand-surface rounded-lg overflow-hidden border border-brand-border shadow-lg">
                <img src={imageUrl} alt={title} className="w-full h-full object-cover" />
                <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center p-2">{renderOverlay()}</div>
                {isLive && <div className="absolute top-2 left-2 bg-brand-danger text-white text-xs font-bold uppercase px-2 py-1 rounded">LIVE</div>}
                {isAvod && <div className="absolute top-2 right-2 bg-brand-warning text-black text-xs font-bold uppercase px-2 py-1 rounded">Ad-Supported</div>}
                {isLive && <LiveProgramProgress startTime={(item as LiveChannel).currentProgram.startTime} endTime={(item as LiveChannel).currentProgram.endTime} />}
                {'progress' in item && <div className="absolute bottom-0 left-0 h-1.5 w-full bg-gray-500/70"><div className="h-full bg-brand-accent" style={{ width: `${(item as MediaContentWithProgress).progress}%` }} /></div>}
            </div>
            <div className="mt-2">
                <p className="text-sm font-semibold text-brand-text-primary truncate">{title}</p>
                {isLive && <p className="text-xs text-brand-text-secondary truncate">{(item as LiveChannel).currentProgram.title}</p>}
                {isTvod && <p className="text-xs text-brand-text-secondary">Rent from ${(item as MediaContent).price?.rent.toFixed(2)}</p>}
            </div>
        </div>
    );
};


const genres: ContentCategoryString[] = ['Nollywood', 'Movies', 'Drama', 'Documentaries', 'Lifestyle'];

interface StreamVerseHomeViewProps {
    onWatch: (id: string) => void;
}

const StreamVerseHomeView: React.FC<StreamVerseHomeViewProps> = ({ onWatch }) => {
    const [selectedGenre, setSelectedGenre] = useState<ContentCategoryString | 'All'>('All');
    const [contentRows, setContentRows] = useState<ContentCategory[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const generateContentRows = async (genre: ContentCategoryString | 'All') => {
            setIsLoading(true);
            const vodContent = await fetchVODContent();
            const liveChannels = await fetchLiveChannels();
            const continueWatching = await fetchContinueWatching();

            const allContent = [...vodContent, ...liveChannels];
            let categories: ContentCategory[] = [];

            if (genre === 'All') {
                if(continueWatching.length > 0) categories.push({ title: 'Continue Watching', items: continueWatching });
                if(liveChannels.length > 0) categories.push({ title: 'Featured Live', items: liveChannels });
                categories.push({ title: 'Popular in Nollywood', items: allContent.filter(c => c.category === 'Nollywood') });
                categories.push({ title: 'Top Movies', items: allContent.filter(c => c.type === 'Movie' && c.category !== 'Nollywood') });
            } else {
                const liveFiltered = liveChannels.filter(c => c.category === genre);
                const vodFiltered = vodContent.filter(c => c.category === genre);
                if(liveFiltered.length > 0) categories.push({ title: `Live ${genre}`, items: liveFiltered });
                if(vodFiltered.length > 0) categories.push({ title: `${genre} On Demand`, items: vodFiltered });
            }
            setContentRows(categories);
            setIsLoading(false);
        };
        
        generateContentRows(selectedGenre);
    }, [selectedGenre]);

    return (
        <div className="flex h-full">
            <nav className="w-48 flex-shrink-0 pr-4">
                <div className="space-y-1">
                    <p className="px-3 pt-2 pb-2 text-xs font-semibold text-brand-text-secondary uppercase">Discover</p>
                    <button onClick={() => setSelectedGenre('All')} className={`w-full text-left px-3 py-2 text-sm font-medium rounded-md ${selectedGenre === 'All' ? 'bg-brand-accent text-white' : 'hover:bg-brand-surface'}`}>All</button>
                    {genres.map(genre => (
                        <button key={genre} onClick={() => setSelectedGenre(genre)} className={`w-full text-left px-3 py-2 text-sm font-medium rounded-md ${selectedGenre === genre ? 'bg-brand-accent text-white' : 'hover:bg-brand-surface'}`}>{genre}</button>
                    ))}
                </div>
            </nav>
            <div className="flex-1 overflow-y-auto pl-4">
                {isLoading ? <LoadingSpinner /> : (
                    <div className="space-y-8">
                        {contentRows.map(category => (
                            <div key={category.title}>
                                <h2 className="text-xl font-bold text-brand-text-primary mb-3">{category.title}</h2>
                                <div className="flex space-x-4 overflow-x-auto pb-4 -mb-4">
                                    {category.items.map(item => <ContentCard key={item.id} item={item} onWatch={onWatch} />)}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default StreamVerseHomeView;