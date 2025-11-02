
import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { getSatelliteRolloutSummary } from '../services/geminiService';
import { SatelliteStatus } from '../types';

const mockUplinkStatus: SatelliteStatus[] = [
    { component: 'DVB-S2X Modulator 1', status: 'Nominal', telemetry: 'C/N: 15.2 dB' },
    { component: 'DVB-S2X Modulator 2', status: 'Nominal', telemetry: 'C/N: 15.1 dB' },
    { component: 'LEO Gateway Partner', status: 'Nominal', telemetry: 'Handover success: 100%' },
];

const mockCarouselStatus: SatelliteStatus[] = [
    { component: 'DVB-NIP Carousel', status: 'Nominal', telemetry: '2.1 Gbps' },
    { component: 'DVB-I Service Catalog', status: 'Warning', telemetry: 'Sync delay: 2.5s' },
    { component: 'DVB-MABR Stream', status: 'Nominal', telemetry: 'All fragments nominal' },
];

const StatusIndicator: React.FC<{ status: 'Nominal' | 'Warning' | 'Error' }> = ({ status }) => {
    const color = {
        Nominal: 'bg-brand-success',
        Warning: 'bg-brand-warning',
        Error: 'bg-brand-danger',
    }[status];
    return <span className={`inline-block h-3 w-3 rounded-full ${color}`}></span>;
}


const SatelliteView: React.FC = () => {
    const [summary, setSummary] = useState<string>('');
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchSummary = async () => {
            setIsLoading(true);
            const content = await getSatelliteRolloutSummary();
            setSummary(content);
            setIsLoading(false);
        };
        fetchSummary();
    }, []);

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-6">
                <DashboardCard title="Headend Uplink Status (PoC)">
                    <ul className="space-y-3">
                        {mockUplinkStatus.map(item => (
                            <li key={item.component} className="p-2 rounded-md bg-brand-bg">
                                <div className="flex justify-between items-center text-sm">
                                    <div className="flex items-center">
                                         <StatusIndicator status={item.status} />
                                         <span className="ml-3 text-brand-text-primary font-semibold">{item.component}</span>
                                    </div>
                                    <span className="font-bold text-brand-text-primary">{item.status}</span>
                                </div>
                                <p className="text-xs text-brand-text-secondary mt-1 ml-6 font-mono">{item.telemetry}</p>
                            </li>
                        ))}
                    </ul>
                </DashboardCard>
                <DashboardCard title="DVB-NIP/I/MABR Carousel Monitor">
                    <ul className="space-y-3">
                        {mockCarouselStatus.map(item => (
                            <li key={item.component} className="p-2 rounded-md bg-brand-bg">
                                <div className="flex justify-between items-center text-sm">
                                    <div className="flex items-center">
                                         <StatusIndicator status={item.status} />
                                         <span className="ml-3 text-brand-text-primary font-semibold">{item.component}</span>
                                    </div>
                                    <span className="font-bold text-brand-text-primary">{item.status}</span>
                                </div>
                                <p className="text-xs text-brand-text-secondary mt-1 ml-6 font-mono">{item.telemetry}</p>
                            </li>
                        ))}
                    </ul>
                </DashboardCard>
            </div>
            <div className="lg:col-span-1">
                <DashboardCard title="Rollout Summary (AI Generated)">
                    {isLoading ? <LoadingSpinner /> : <GeminiResponseDisplay content={summary} />}
                </DashboardCard>
            </div>
        </div>
    );
};

export default SatelliteView;
