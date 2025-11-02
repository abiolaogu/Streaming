
import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { getDrmStatusSummary } from '../services/geminiService';
import { DrmLicenseServer, DrmKey } from '../types';

const mockLicenseServers: DrmLicenseServer[] = [
    { id: 'wv-us-east-1', region: 'us-east-1', provider: 'Widevine', status: 'Active', p95LatencyMs: 45 },
    { id: 'pr-eu-west-1', region: 'eu-west-1', provider: 'PlayReady', status: 'Active', p95LatencyMs: 62 },
    { id: 'fp-ap-ne-1', region: 'ap-northeast-1', provider: 'FairPlay', status: 'Degraded', p95LatencyMs: 150 },
    { id: 'wv-sa-east-1', region: 'sa-east-1', provider: 'Widevine', status: 'Active', p95LatencyMs: 88 },
    { id: 'pr-us-west-2', region: 'us-west-2', provider: 'PlayReady', status: 'Offline', p95LatencyMs: 0 },
];

const mockKeys: DrmKey[] = [
    { contentId: 'movie-1', status: 'Active', lastRotation: '2024-05-10', licensesIssued: 152340 },
    { contentId: 'series-1', status: 'Active', lastRotation: '2024-05-08', licensesIssued: 893421 },
    { contentId: 'live-event-sports', status: 'Rotated', lastRotation: '2024-05-20', licensesIssued: 54032 },
    { contentId: 'movie-old-promo', status: 'Compromised', lastRotation: '2023-11-20', licensesIssued: 1023 },
];

const getStatusColor = (status: DrmLicenseServer['status'] | DrmKey['status']) => {
    switch(status) {
        case 'Active': return 'text-brand-success';
        case 'Degraded':
        case 'Rotated': return 'text-brand-warning';
        case 'Offline':
        case 'Compromised': return 'text-brand-danger';
        default: return 'text-brand-text-secondary';
    }
};

const DrmView: React.FC = () => {
    const [summary, setSummary] = useState<string>('');
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchSummary = async () => {
            setIsLoading(true);
            const content = await getDrmStatusSummary();
            setSummary(content);
            setIsLoading(false);
        };
        fetchSummary();
    }, []);

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 space-y-6">
                <DashboardCard title="License Server Status">
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-brand-border">
                            <thead className="bg-brand-bg">
                                <tr>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Region</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Provider</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Status</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">p95 Latency</th>
                                </tr>
                            </thead>
                            <tbody className="bg-brand-surface divide-y divide-brand-border">
                                {mockLicenseServers.map((server) => (
                                    <tr key={server.id}>
                                        <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-primary">{server.region}</td>
                                        <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-secondary">{server.provider}</td>
                                        <td className={`px-4 py-3 whitespace-nowrap text-sm font-semibold ${getStatusColor(server.status)}`}>{server.status}</td>
                                        <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-secondary">{server.status === 'Offline' ? 'N/A' : `${server.p95LatencyMs}ms`}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </DashboardCard>
                <DashboardCard title="Key Rotation & Usage">
                     <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-brand-border">
                            <thead className="bg-brand-bg">
                                <tr>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Content ID</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Status</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Last Rotation</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Licenses Issued</th>
                                </tr>
                            </thead>
                            <tbody className="bg-brand-surface divide-y divide-brand-border">
                                {mockKeys.map((key) => (
                                    <tr key={key.contentId}>
                                        <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-primary font-mono">{key.contentId}</td>
                                        <td className={`px-4 py-3 whitespace-nowrap text-sm font-semibold ${getStatusColor(key.status)}`}>{key.status}</td>
                                        <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-secondary">{key.lastRotation}</td>
                                        <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-secondary">{key.licensesIssued.toLocaleString()}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </DashboardCard>
            </div>
            <div className="lg:col-span-1">
                <DashboardCard title="DRM Status Summary (AI)">
                    {isLoading ? <LoadingSpinner /> : <GeminiResponseDisplay content={summary} />}
                </DashboardCard>
            </div>
        </div>
    );
};

export default DrmView;
