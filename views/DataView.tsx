
import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { getDataPlatformHealth } from '../services/geminiService';
import { DataPlatformMetric } from '../types';

const mockKafkaMetrics: DataPlatformMetric[] = [
    { name: 'Ingest Topic Throughput', value: '1.2 GB/s', status: 'ok' },
    { name: 'Analytics Topic Lag', value: '5,230', status: 'warning' },
    { name: 'MM2 Replication (US-EU)', value: '350 ms', status: 'ok' },
];

const mockDbMetrics: DataPlatformMetric[] = [
    { name: 'YugabyteDB Cluster Health', value: 'Nominal', status: 'ok' },
    { name: 'YugabyteDB Replication Lag', value: '< 100ms', status: 'ok' },
    { name: 'ScyllaDB p99 Latency', value: '12ms', status: 'ok' },
    { name: 'DragonflyDB Hit Rate', value: '98.7%', status: 'ok' },
];

const StatusIndicator: React.FC<{ status: 'ok' | 'warning' | 'critical' }> = ({ status }) => {
    const color = {
        ok: 'bg-brand-success',
        warning: 'bg-brand-warning',
        critical: 'bg-brand-danger',
    }[status];
    return <span className={`inline-block h-3 w-3 rounded-full ${color}`}></span>;
}

const DataView: React.FC = () => {
    const [report, setReport] = useState<string>('');
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchReport = async () => {
            setIsLoading(true);
            const content = await getDataPlatformHealth();
            setReport(content);
            setIsLoading(false);
        };
        fetchReport();
    }, []);

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-6">
                <DashboardCard title="Kafka Cluster Health">
                    <ul className="space-y-3">
                        {mockKafkaMetrics.map(metric => (
                            <li key={metric.name} className="flex justify-between items-center text-sm p-2 rounded-md bg-brand-bg">
                                <div className="flex items-center">
                                    <StatusIndicator status={metric.status} />
                                    <span className="ml-3 text-brand-text-secondary">{metric.name}</span>
                                </div>
                                <span className="font-semibold text-brand-text-primary">{metric.value}</span>
                            </li>
                        ))}
                    </ul>
                </DashboardCard>
                <DashboardCard title="Data Stores (YugabyteDB, ScyllaDB, DragonflyDB)">
                    <ul className="space-y-3">
                        {mockDbMetrics.map(metric => (
                            <li key={metric.name} className="flex justify-between items-center text-sm p-2 rounded-md bg-brand-bg">
                                <div className="flex items-center">
                                    <StatusIndicator status={metric.status} />
                                    <span className="ml-3 text-brand-text-secondary">{metric.name}</span>
                                </div>
                                <span className="font-semibold text-brand-text-primary">{metric.value}</span>
                            </li>
                        ))}
                    </ul>
                </DashboardCard>
            </div>
            <div className="lg:col-span-1">
                <DashboardCard title="Data Platform Report (AI Generated)">
                    {isLoading ? <LoadingSpinner /> : <GeminiResponseDisplay content={report} />}
                </DashboardCard>
            </div>
        </div>
    );
};

export default DataView;