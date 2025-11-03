import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import { fetchUserProfile, fetchUserPreferences, fetchUserDevices, fetchWatchHistory, fetchWatchlist } from '../services/mockApiService';
import { UserProfile, UserPreferences, UserDevice, AnyContent, MediaContent } from '../types';
import { DevicePhoneMobileIcon, ComputerDesktopIcon, TvIcon, TrashIcon } from '@heroicons/react/24/solid';

const deviceIcons: Record<UserDevice['type'], React.ReactNode> = {
    'Mobile': <DevicePhoneMobileIcon className="h-6 w-6 text-brand-text-secondary" />,
    'Web': <ComputerDesktopIcon className="h-6 w-6 text-brand-text-secondary" />,
    'Smart TV': <TvIcon className="h-6 w-6 text-brand-text-secondary" />,
};

const ContentItemCard: React.FC<{ item: AnyContent, onRemove: (id: string) => void }> = ({ item, onRemove }) => {
    const mediaItem = item as MediaContent;
    return (
        <div className="flex items-center space-x-4 p-2 bg-brand-bg rounded-lg hover:bg-brand-surface transition-colors">
            <img src={mediaItem.thumbnailUrl} alt={mediaItem.title} className="w-24 h-14 object-cover rounded-md flex-shrink-0" />
            <div className="flex-1 min-w-0">
                <p className="font-semibold text-sm text-brand-text-primary truncate">{mediaItem.title}</p>
                <p className="text-xs text-brand-text-secondary">{mediaItem.category}</p>
            </div>
            <button 
                onClick={() => onRemove(item.id)}
                className="text-brand-text-secondary hover:text-brand-danger p-2 rounded-full"
                aria-label={`Remove ${mediaItem.title}`}
            >
                <TrashIcon className="h-5 w-5" />
            </button>
        </div>
    );
};


const UserProfileView: React.FC = () => {
    const [profile, setProfile] = useState<UserProfile | null>(null);
    const [preferences, setPreferences] = useState<UserPreferences | null>(null);
    const [devices, setDevices] = useState<UserDevice[]>([]);
    const [watchHistory, setWatchHistory] = useState<AnyContent[]>([]);
    const [watchlist, setWatchlist] = useState<AnyContent[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            setIsLoading(true);
            const [profileData, preferencesData, devicesData, historyData, listData] = await Promise.all([
                fetchUserProfile(),
                fetchUserPreferences(),
                fetchUserDevices(),
                fetchWatchHistory(),
                fetchWatchlist(),
            ]);
            setProfile(profileData);
            setPreferences(preferencesData);
            setDevices(devicesData);
            setWatchHistory(historyData);
            setWatchlist(listData);
            setIsLoading(false);
        };
        loadData();
    }, []);

    const handleRemoveHistory = (id: string) => setWatchHistory(prev => prev.filter(item => item.id !== id));
    const handleRemoveWatchlist = (id: string) => setWatchlist(prev => prev.filter(item => item.id !== id));
    
    const handleClearAllHistory = () => {
        if (window.confirm('Are you sure you want to clear your entire watch history? This action cannot be undone.')) {
            setWatchHistory([]);
        }
    };

    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-full">
                <LoadingSpinner />
            </div>
        );
    }
    
    if(!profile || !preferences) {
        return <p>Error loading user data.</p>;
    }

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-1 space-y-6">
                <DashboardCard title="Profile">
                   <div className="space-y-3 text-sm">
                        <div className="flex justify-between">
                            <span className="text-brand-text-secondary">Name</span>
                            <span className="font-semibold text-brand-text-primary">{profile.name}</span>
                        </div>
                         <div className="flex justify-between">
                            <span className="text-brand-text-secondary">Email</span>
                            <span className="font-semibold text-brand-text-primary">{profile.email}</span>
                        </div>
                         <div className="flex justify-between">
                            <span className="text-brand-text-secondary">Member Since</span>
                            <span className="font-semibold text-brand-text-primary">{profile.memberSince}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-brand-text-secondary">Plan</span>
                            <span className="font-semibold text-brand-accent">{profile.subscriptionPlan}</span>
                        </div>
                   </div>
                </DashboardCard>
                <DashboardCard title="Preferences">
                   <div className="space-y-3 text-sm">
                        <div className="flex justify-between items-center">
                            <label htmlFor="language" className="text-brand-text-secondary">Language</label>
                            <select id="language" defaultValue={preferences.language} className="bg-brand-bg border border-brand-border rounded-md p-1 text-sm">
                                <option>English</option>
                                <option>Español</option>
                                <option>Français</option>
                            </select>
                        </div>
                        <div className="flex justify-between items-center">
                            <label htmlFor="quality" className="text-brand-text-secondary">Playback Quality</label>
                            <select id="quality" defaultValue={preferences.playbackQuality} className="bg-brand-bg border border-brand-border rounded-md p-1 text-sm">
                                <option>Auto</option>
                                <option>HD</option>
                                <option>4K</option>
                            </select>
                        </div>
                        <div className="flex justify-between items-center">
                            <label htmlFor="subtitles" className="text-brand-text-secondary">Enable Subtitles</label>
                            <input id="subtitles" type="checkbox" defaultChecked={preferences.subtitles} className="h-4 w-4 text-brand-accent bg-brand-bg border-brand-border rounded focus:ring-brand-accent" />
                        </div>
                   </div>
                   <button className="w-full mt-4 bg-brand-accent hover:bg-brand-accent-hover text-white font-bold py-2 px-4 rounded-md transition">
                       Save Preferences
                   </button>
                </DashboardCard>
            </div>
            <div className="lg:col-span-2 space-y-6">
                <DashboardCard title="Active Device Sessions">
                    <p className="text-xs text-brand-text-secondary mb-4">You can have a maximum of 5 active devices. Deregister a device to free up a slot.</p>
                    <ul className="space-y-3">
                        {devices.map(device => (
                             <li key={device.id} className="p-3 bg-brand-bg rounded-lg flex items-center justify-between">
                                <div className="flex items-center space-x-4">
                                    {deviceIcons[device.type]}
                                    <div>
                                        <p className="font-semibold text-brand-text-primary">{device.os}</p>
                                        <p className="text-xs text-brand-text-secondary">Last seen: {device.lastSeen}</p>
                                    </div>
                                </div>
                                <button className="text-xs bg-brand-danger/20 text-brand-danger font-bold py-1 px-3 rounded-md hover:bg-brand-danger/40 transition">Deregister</button>
                            </li>
                        ))}
                    </ul>
                </DashboardCard>
                 <DashboardCard title={
                    <div className="flex justify-between items-center">
                        <span>Watch History</span>
                        {watchHistory.length > 0 && (
                             <button 
                                onClick={handleClearAllHistory}
                                className="text-xs font-semibold text-brand-text-secondary hover:text-brand-danger transition-colors"
                            >
                                Clear History
                            </button>
                        )}
                    </div>
                 }>
                    <div className="space-y-3 max-h-80 overflow-y-auto pr-2">
                        {watchHistory.length > 0 ? (
                            watchHistory.map(item => <ContentItemCard key={`hist-${item.id}`} item={item} onRemove={handleRemoveHistory} />)
                        ) : (
                            <p className="text-sm text-brand-text-secondary text-center py-4">Your watch history is empty.</p>
                        )}
                    </div>
                </DashboardCard>
                <DashboardCard title="My Watchlist">
                    <div className="space-y-3 max-h-80 overflow-y-auto pr-2">
                        {watchlist.length > 0 ? (
                            watchlist.map(item => <ContentItemCard key={`list-${item.id}`} item={item} onRemove={handleRemoveWatchlist} />)
                        ) : (
                             <p className="text-sm text-brand-text-secondary text-center py-4">Your watchlist is empty.</p>
                        )}
                    </div>
                </DashboardCard>
            </div>
        </div>
    );
};

export default UserProfileView;