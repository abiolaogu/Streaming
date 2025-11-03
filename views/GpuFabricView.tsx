
import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { getGpuScalingRecommendation } from '../services/geminiService';
import { fetchGpuInstances } from '../services/mockApiService';
import { GpuInstance } from '../types';

const GpuFabricView: React.FC = () => {
    const [workload, setWorkload] = useState('High VOD ingest queue depth for 4K AV1 encodes. 2 major live events starting in 15 minutes.');
    const [recommendation, setRecommendation] = useState<string>('');
    const [isLoadingRecommendation, setIsLoadingRecommendation] = useState<boolean>(false);
    const [gpuInstances, setGpuInstances] = useState<GpuInstance[]>([]);
    const [isLoadingInstances, setIsLoadingInstances] = useState<boolean>(true);

    useEffect(() => {
        const loadInstances = async () => {
            setIsLoadingInstances(true);
            const data = await fetchGpuInstances();
            setGpuInstances(data);
            setIsLoadingInstances(false);
        };
        loadInstances();
    }, []);

    const handleGetRecommendation = async () => {
        setIsLoadingRecommendation(true);
        setRecommendation('');
        const content = await getGpuScalingRecommendation(workload);
        setRecommendation(content);
        setIsLoadingRecommendation(false);
    };

    const getStatusColor = (status: GpuInstance['status']) => {
        switch (status) {
            case 'Processing': return 'text-brand-success';
            case 'Idle': return 'text-brand-warning';
            case 'Terminating': return 'text-brand-danger';
            default: return 'text-brand-text-secondary';
        }
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <DashboardCard title="GPU Instance Status">
                {isLoadingInstances ? <LoadingSpinner /> : (
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-brand-border">
                            <thead className="bg-brand-bg">
                                <tr>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Provider</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Type</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Status</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Cost/hr</th>
                                </tr>
                            </thead>
                            <tbody className="bg-brand-surface divide-y divide-brand-border">
                                {gpuInstances.map((gpu) => (
                                    <tr key={gpu.id}>
                                        <td className="px-4 py-2 whitespace-nowrap text-sm text-brand-text-primary">{gpu.provider} {gpu.spot && '(Spot)'}</td>
                                        <td className="px-4 py-2 whitespace-nowrap text-sm text-brand-text-secondary">{gpu.type}</td>
                                        <td className={`px-4 py-2 whitespace-nowrap text-sm font-semibold ${getStatusColor(gpu.status)}`}>{gpu.status}</td>
                                        <td className="px-4 py-2 whitespace-nowrap text-sm text-brand-text-secondary">${gpu.costPerHour.toFixed(2)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </DashboardCard>

            <DashboardCard title="AI Scaling Recommendation">
                <div className="space-y-4">
                    <div>
                        <label htmlFor="workload" className="block text-sm font-medium text-brand-text-secondary">
                            Current Workload
                        </label>
                        <textarea
                            id="workload"
                            rows={3}
                            className="mt-1 block w-full bg-brand-bg border border-brand-border rounded-md shadow-sm p-2 text-sm text-brand-text-primary focus:ring-brand-accent focus:border-brand-accent"
                            value={workload}
                            onChange={(e) => setWorkload(e.target.value)}
                        />
                    </div>
                    <button
                        onClick={handleGetRecommendation}
                        disabled={isLoadingRecommendation}
                        className="w-full bg-brand-accent hover:bg-brand-accent-hover text-white font-bold py-2 px-4 rounded-md transition duration-200 disabled:bg-gray-500 disabled:cursor-not-allowed"
                    >
                        {isLoadingRecommendation ? 'Generating...' : 'Get Recommendation'}
                    </button>
                    <div className="mt-4 p-4 min-h-[100px] bg-brand-bg rounded-md border border-dashed border-brand-border">
                        {isLoadingRecommendation && <LoadingSpinner />}
                        {recommendation && <GeminiResponseDisplay content={recommendation} />}
                        {!isLoadingRecommendation && !recommendation && <p className="text-brand-text-secondary text-center">AI recommendation will appear here.</p>}
                    </div>
                </div>
            </DashboardCard>
        </div>
    );
};

export default GpuFabricView;
