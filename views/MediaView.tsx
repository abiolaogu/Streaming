
import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { getMediaPipelineStatus } from '../services/geminiService';
import { fetchTranscodeJobs, fetchLiveIngests } from '../services/mockApiService';
import { TranscodeJob, LiveStreamIngest } from '../types';

const getTranscodeStatusColor = (status: TranscodeJob['status']) => {
    switch (status) {
        case 'In Progress': return 'text-brand-accent';
        case 'Completed': return 'text-brand-success';
        case 'Queued': return 'text-brand-warning';
        case 'Failed': return 'text-brand-danger';
        default: return 'text-brand-text-secondary';
    }
};

const getIngestStatusColor = (status: LiveStreamIngest['status']) => {
    switch(status) {
        case 'Stable': return 'text-brand-success';
        case 'Packet Loss': return 'text-brand-warning';
        case 'Critical': return 'text-brand-danger animate-pulse';
        case 'Offline': return 'text-brand-danger';
        default: return 'text-brand-text-secondary';
    }
};

const PipelineStage: React.FC<{ name: string, status: 'ok' | 'warning' | 'critical' }> = ({ name, status }) => {
    const colorClasses = {
        ok: 'border-brand-success bg-green-900/50',
        warning: 'border-brand-warning bg-yellow-900/50',
        critical: 'border-brand-danger bg-red-900/50',
    };
    return (
        <div className={`flex-1 text-center p-4 rounded-lg border-2 ${colorClasses[status]}`}>
            <p className="font-semibold text-brand-text-primary">{name}</p>
        </div>
    );
};

const MediaView: React.FC = () => {
    const [report, setReport] = useState<string>('');
    const [isLoadingReport, setIsLoadingReport] = useState<boolean>(true);
    const [transcodeJobs, setTranscodeJobs] = useState<TranscodeJob[]>([]);
    const [liveIngests, setLiveIngests] = useState<LiveStreamIngest[]>([]);
    const [isLoadingData, setIsLoadingData] = useState<boolean>(true);

    useEffect(() => {
        const fetchAllData = async () => {
            setIsLoadingReport(true);
            setIsLoadingData(true);
            
            const reportContent = getMediaPipelineStatus();
            const jobsData = fetchTranscodeJobs();
            const ingestsData = fetchLiveIngests();

            const [report, jobs, ingests] = await Promise.all([reportContent, jobsData, ingestsData]);

            setReport(report);
            setTranscodeJobs(jobs);
            setLiveIngests(ingests);

            setIsLoadingReport(false);
            setIsLoadingData(false);
        };
        fetchAllData();
    }, []);

    return (
        <div className="space-y-6">
            <DashboardCard title="Live Media Pipeline Status">
                <div className="flex items-center space-x-2 md:space-x-4 text-sm md:text-base">
                    <PipelineStage name="Ingest" status="critical" />
                    <div className="text-brand-text-secondary font-bold text-2xl">&rarr;</div>
                    <PipelineStage name="Transcode" status="warning" />
                    <div className="text-brand-text-secondary font-bold text-2xl">&rarr;</div>
                    <PipelineStage name="Package & DRM" status="ok" />
                    <div className="text-brand-text-secondary font-bold text-2xl">&rarr;</div>
                    <PipelineStage name="Origin Store" status="ok" />
                </div>
            </DashboardCard>
            
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className="lg:col-span-2 space-y-6">
                    {isLoadingData ? <DashboardCard title="Live Stream Ingest Status"><LoadingSpinner /></DashboardCard> : (
                        <DashboardCard title="Live Stream Ingest Status">
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-brand-border">
                                    <thead className="bg-brand-bg">
                                        <tr>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Stream ID</th>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Protocol</th>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Status</th>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Connections</th>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Bitrate</th>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Issues</th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-brand-surface divide-y divide-brand-border">
                                        {liveIngests.map((stream) => (
                                            <tr key={stream.id}>
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-primary font-mono">{stream.id}</td>
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-secondary">{stream.protocol}</td>
                                                <td className={`px-4 py-3 whitespace-nowrap text-sm font-semibold ${getIngestStatusColor(stream.status)}`}>{stream.status}</td>
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-center text-brand-text-secondary">{stream.connections}</td>
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-secondary">{stream.bitrateMb.toFixed(1)} Mbps</td>
                                                <td className={`px-4 py-3 whitespace-nowrap text-sm ${stream.issues !== 'None' ? 'text-brand-warning' : 'text-brand-text-secondary'}`}>{stream.issues}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </DashboardCard>
                    )}
                    {isLoadingData ? <DashboardCard title="Active Transcode Jobs"><LoadingSpinner /></DashboardCard> : (
                        <DashboardCard title="Active Transcode Jobs">
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-brand-border">
                                    <thead className="bg-brand-bg">
                                        <tr>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Source</th>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Profile</th>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Status</th>
                                            <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Progress</th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-brand-surface divide-y divide-brand-border">
                                        {transcodeJobs.map((job) => (
                                            <tr key={job.id}>
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-primary">{job.source}</td>
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-secondary font-mono">{job.profile}</td>
                                                <td className={`px-4 py-3 whitespace-nowrap text-sm font-semibold ${getTranscodeStatusColor(job.status)}`}>{job.status}</td>
                                                <td className="px-4 py-3 whitespace-nowrap text-sm text-brand-text-secondary">
                                                    <div className="w-full bg-brand-bg rounded-full h-2.5">
                                                        <div className="bg-brand-accent h-2.5 rounded-full" style={{ width: `${job.progress}%` }}></div>
                                                    </div>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </DashboardCard>
                    )}
                </div>
                <div className="lg:col-span-1">
                    <DashboardCard title="Pipeline Status Report (AI)">
                        {isLoadingReport ? <LoadingSpinner /> : <GeminiResponseDisplay content={report} />}
                    </DashboardCard>
                </div>
            </div>
        </div>
    );
};

export default MediaView;
