import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import { fetchDvbComponents } from '../services/mockApiService';
import { DvbComponent } from '../types';
import { SatelliteIcon, CdnIcon, TelecomIcon } from '../components/IconComponents';

const componentIcons = {
    'DVB-NIP': <SatelliteIcon className="h-6 w-6 text-blue-400" />,
    'DVB-I': <CdnIcon className="h-6 w-6 text-green-400" />,
    'DVB-IP': <TelecomIcon className="h-6 w-6 text-purple-400" />,
};

const statusColors = {
    Nominal: 'text-brand-success',
    Warning: 'text-brand-warning',
    Critical: 'text-brand-danger',
};


const BroadcastOpsView: React.FC = () => {
    const [dvbComponents, setDvbComponents] = useState<DvbComponent[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            setIsLoading(true);
            const data = await fetchDvbComponents();
            setDvbComponents(data);
            setIsLoading(false);
        };
        loadData();
    }, []);

    return (
        <DashboardCard title="Hybrid Broadcast-Streaming Operations (DVB)">
            {isLoading ? <LoadingSpinner /> : (
                <div className="space-y-4">
                    {dvbComponents.map(component => (
                        <div key={component.id} className="p-4 bg-brand-bg rounded-lg border border-brand-border flex items-start space-x-4">
                            <div className="flex-shrink-0 mt-1">
                                {componentIcons[component.type]}
                            </div>
                            <div className="flex-1">
                                <div className="flex justify-between items-center">
                                    <p className="text-sm font-semibold text-brand-text-primary">{component.type}</p>
                                    <p className={`text-sm font-bold ${statusColors[component.status]}`}>{component.status}</p>
                                </div>
                                <p className="mt-1 text-sm text-brand-text-secondary font-mono">{component.id}</p>
                                <p className="mt-1 text-xs text-brand-text-secondary">{component.details}</p>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </DashboardCard>
    );
};

export default BroadcastOpsView;