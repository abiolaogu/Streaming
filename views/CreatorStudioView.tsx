import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { fetchCreatorContent, fetchRevenueAnalytics } from '../services/mockApiService';
import { getCreatorAnalyticsSummary } from '../services/geminiService';
import { CreatorContent, RevenueAnalytics } from '../types';
import { DocumentArrowUpIcon, CheckCircleIcon, ClockIcon, XCircleIcon, ChartBarIcon } from '@heroicons/react/24/solid';

const statusIcons = {
    Live: <CheckCircleIcon className="h-5 w-5 text-brand-success" />,
    'In Review': <ClockIcon className="h-5 w-5 text-brand-warning" />,
    Rejected: <XCircleIcon className="h-5 w-5 text-brand-danger" />,
    Processing: <LoadingSpinner />,
};

const CreatorStudioView: React.FC = () => {
    const [content, setContent] = useState<CreatorContent[]>([]);
    const [analytics, setAnalytics] = useState<RevenueAnalytics[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [selectedContent, setSelectedContent] = useState<CreatorContent | null>(null);
    const [aiSummary, setAiSummary] = useState('');
    const [isSummaryLoading, setIsSummaryLoading] = useState(false);

    useEffect(() => {
        const loadData = async () => {
            setIsLoading(true);
            const [contentData, analyticsData] = await Promise.all([fetchCreatorContent(), fetchRevenueAnalytics()]);
            setContent(contentData);
            setAnalytics(analyticsData);
            setSelectedContent(contentData[0] || null);
            setIsLoading(false);
        };
        loadData();
    }, []);

    useEffect(() => {
        if (selectedContent) {
            const fetchSummary = async () => {
                setIsSummaryLoading(true);
                const summary = await getCreatorAnalyticsSummary(selectedContent.title);
                setAiSummary(summary);
                setIsSummaryLoading(false);
            };
            fetchSummary();
        }
    }, [selectedContent]);

    if (isLoading) {
        return <LoadingSpinner />;
    }

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 space-y-6">
                <DashboardCard title="Content Management & Upload">
                     <div className="flex justify-end mb-4">
                        <button className="flex items-center gap-2 bg-brand-accent hover:bg-brand-accent-hover text-white font-bold py-2 px-4 rounded-md transition">
                           <DocumentArrowUpIcon className="h-5 w-5" /> Upload New Content
                        </button>
                     </div>
                     <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-brand-border">
                            <thead className="bg-brand-bg">
                                <tr>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase">Title</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase">Status</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase">Upload Date</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase">Monetization</th>
                                </tr>
                            </thead>
                            <tbody className="bg-brand-surface divide-y divide-brand-border">
                                {content.map(item => (
                                    <tr key={item.id} onClick={() => setSelectedContent(item)} className={`cursor-pointer ${selectedContent?.id === item.id ? 'bg-brand-accent/20' : 'hover:bg-brand-surface/50'}`}>
                                        <td className="px-4 py-3 text-sm font-medium text-brand-text-primary">{item.title}</td>
                                        <td className="px-4 py-3 text-sm"><span className="flex items-center gap-2">{statusIcons[item.status]} {item.status}</span></td>
                                        <td className="px-4 py-3 text-sm text-brand-text-secondary">{item.uploadDate}</td>
                                        <td className="px-4 py-3 text-sm text-brand-text-secondary">{item.monetization.join(', ')}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </DashboardCard>

                <DashboardCard title={`Analytics for: ${selectedContent?.title || 'Select Content'}`}>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                        {analytics.map(item => (
                            <div key={item.period} className="bg-brand-bg p-4 rounded-lg">
                                <p className="text-sm text-brand-text-secondary">{item.period} Revenue</p>
                                <p className="text-2xl font-bold text-brand-accent">${item.totalRevenue.toLocaleString()}</p>
                                <p className="text-xs text-brand-text-secondary">{item.watchHours.toLocaleString()} watch hours</p>
                            </div>
                        ))}
                    </div>
                     <p className="text-xs text-center text-brand-text-secondary mt-4">Revenue is calculated via watch-time attribution from SVOD, AVOD, and TVOD sources. 60/40 split.</p>
                </DashboardCard>

            </div>

            <div className="lg:col-span-1">
                <DashboardCard title="AI Performance Insights">
                    {isSummaryLoading ? <LoadingSpinner /> : <GeminiResponseDisplay content={aiSummary} />}
                </DashboardCard>
            </div>
        </div>
    );
};

export default CreatorStudioView;