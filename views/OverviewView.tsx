import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { getDailyBriefing } from '../services/geminiService';
import { SloMetric } from '../types';

const mockSloMetrics: SloMetric[] = [
    { name: 'Video Startup Time (p95)', value: '280ms', target: '< 400ms', status: 'ok' },
    { name: 'Rebuffer Rate', value: '0.12%', target: '< 0.2%', status: 'ok' },
    { name: 'Deploy Frequency', value: '3.2/day', target: '> 1/day', status: 'ok' },
    { name: 'Change Fail Rate', value: '4.5%', target: '< 15%', status: 'ok' },
    { name: 'On-Demand CPU Nodes', value: '3', target: '≤ 3', status: 'ok' },
    // FIX: Changed property 'warning' to 'status' to match the SloMetric type definition.
    { name: 'Spot GPU Availability', value: '92%', target: '> 90%', status: 'warning' },
];

const StatusIndicator: React.FC<{ status: 'ok' | 'warning' | 'critical' }> = ({ status }) => {
    const color = {
        ok: 'bg-brand-success',
        warning: 'bg-brand-warning',
        critical: 'bg-brand-danger',
    }[status];
    return <span className={`inline-block h-3 w-3 rounded-full ${color}`}></span>;
}

const OverviewView: React.FC = () => {
    const [briefing, setBriefing] = useState<string>('');
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchBriefing = async () => {
            setIsLoading(true);
            const content = await getDailyBriefing();
            setBriefing(content);
            setIsLoading(false);
        };
        fetchBriefing();
    }, []);

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2">
                <DashboardCard title="Architect's Daily Briefing (AI Generated)">
                    {isLoading ? <LoadingSpinner /> : <GeminiResponseDisplay content={briefing} />}
                </DashboardCard>
            </div>
            <div className="lg:col-span-1">
                 <DashboardCard title="Key SLOs & Metrics">
                    <ul className="space-y-3">
                        {mockSloMetrics.map(metric => (
                            <li key={metric.name} className="flex justify-between items-center text-sm p-2 rounded-md bg-brand-bg">
                                <div className="flex items-center">
                                    <StatusIndicator status={metric.status} />
                                    <span className="ml-3 text-brand-text-secondary">{metric.name}</span>
                                </div>
                                <div>
                                    <span className="font-semibold text-brand-text-primary">{metric.value}</span>
                                    <span className="text-brand-text-secondary"> / {metric.target}</span>
                                </div>
                            </li>
                        ))}
                    </ul>
                </DashboardCard>
            </div>
        </div>
    );
};

export default OverviewView;