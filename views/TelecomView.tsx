
import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { getTelecomStatusReport } from '../services/geminiService';
import { TelecomServiceStatus } from '../types';

const mockVoiceServices: TelecomServiceStatus[] = [
    { name: 'Kamailio Registrar (US-East)', status: 'Online', details: '2.1M registrations' },
    { name: 'FreeSWITCH IVR (EU-West)', status: 'Online', details: '3,400 active channels' },
    { name: 'RTPengine (APAC-NE)', status: 'Degraded', details: 'High packet loss detected' },
];

const mockMobileCoreServices: TelecomServiceStatus[] = [
    { name: 'Open5GS AMF/SMF', status: 'Online', details: 'All core functions nominal' },
    { name: 'UPF (Tier-1 Edges)', status: 'Online', details: '1.2 Tbps throughput' },
    { name: 'Subscribers (Private 5G)', status: 'Online', details: '15,780 active UEs' },
];

const StatusIndicator: React.FC<{ status: 'Online' | 'Degraded' | 'Offline' }> = ({ status }) => {
    const color = {
        Online: 'bg-brand-success',
        Degraded: 'bg-brand-warning',
        Offline: 'bg-brand-danger',
    }[status];
    return <span className={`inline-block h-3 w-3 rounded-full ${color}`}></span>;
}

const TelecomView: React.FC = () => {
    const [report, setReport] = useState<string>('');
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchReport = async () => {
            setIsLoading(true);
            const content = await getTelecomStatusReport();
            setReport(content);
            setIsLoading(false);
        };
        fetchReport();
    }, []);

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-6">
                <DashboardCard title="Voice/IMS Service Status">
                     <ul className="space-y-3">
                        {mockVoiceServices.map(service => (
                            <li key={service.name} className="p-2 rounded-md bg-brand-bg">
                                <div className="flex justify-between items-center text-sm">
                                    <div className="flex items-center">
                                         <StatusIndicator status={service.status} />
                                         <span className="ml-3 text-brand-text-primary font-semibold">{service.name}</span>
                                    </div>
                                    <span className="font-bold text-brand-text-primary">{service.status}</span>
                                </div>
                                <p className="text-xs text-brand-text-secondary mt-1 ml-6">{service.details}</p>
                            </li>
                        ))}
                    </ul>
                </DashboardCard>
                <DashboardCard title="Open5GS Mobile Core">
                    <ul className="space-y-3">
                        {mockMobileCoreServices.map(service => (
                            <li key={service.name} className="p-2 rounded-md bg-brand-bg">
                                <div className="flex justify-between items-center text-sm">
                                    <div className="flex items-center">
                                         <StatusIndicator status={service.status} />
                                         <span className="ml-3 text-brand-text-primary font-semibold">{service.name}</span>
                                    </div>
                                    <span className="font-bold text-brand-text-primary">{service.status}</span>
                                </div>
                                <p className="text-xs text-brand-text-secondary mt-1 ml-6">{service.details}</p>
                            </li>
                        ))}
                    </ul>
                </DashboardCard>
            </div>
            <div className="lg:col-span-1">
                <DashboardCard title="Telecom Status Report (AI)">
                    {isLoading ? <LoadingSpinner /> : <GeminiResponseDisplay content={report} />}
                </DashboardCard>
            </div>
        </div>
    );
};

export default TelecomView;
