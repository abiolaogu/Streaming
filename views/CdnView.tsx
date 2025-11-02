import React, { useEffect, useState } from 'react';
import DashboardCard from '../components/DashboardCard';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import LoadingSpinner from '../components/LoadingSpinner';
import { getCdnStatusReport } from '../services/geminiService';
import { fetchCdnStats, fetchCdnAlerts } from '../services/mockApiService';
import DayNightOverlay from '../components/DayNightOverlay';
import WorldMap from '../components/WorldMap';
import { CdnStat, CdnAlert } from '../types';
import { ArrowUpIcon, ArrowDownIcon, BellAlertIcon, ExclamationTriangleIcon, InformationCircleIcon } from '@heroicons/react/24/solid';

const pops = {
    // Tier 1
    ash1: { x: 25.4, y: 32.5, tier: 1, name: 'Ashburn, USA' },
    lon1: { x: 49.7, y: 23, tier: 1, name: 'London, UK' },
    sgp1: { x: 79.8, y: 53.5, tier: 1, name: 'Singapore' },
    sao1: { x: 35.5, y: 69.5, tier: 1, name: 'São Paulo, Brazil' },
    los1: { x: 50.8, y: 49.5, tier: 1, name: 'Lagos, Nigeria' },
    // Tier 2
    lax2: { x: 16.5, y: 35.5, tier: 2, name: 'Los Angeles, USA' },
    mia2: { x: 24.3, y: 40, tier: 2, name: 'Miami, USA' },
    fra2: { x: 52.4, y: 24.5, tier: 2, name: 'Frankfurt, DE' },
    mum2: { x: 69.4, y: 42.5, tier: 2, name: 'Mumbai, India' },
    syd2: { x: 91.5, y: 72.5, tier: 2, name: 'Sydney, Australia' },
    tyo2: { x: 88, y: 34, tier: 2, name: 'Tokyo, Japan' },
};

const StatCard: React.FC<{ stat: CdnStat }> = ({ stat }) => {
    const isPositive = stat.changeType === 'positive';
    const isNegative = stat.changeType === 'negative';
    const changeColor = isPositive ? 'text-brand-success' : isNegative ? 'text-brand-danger' : 'text-brand-text-secondary';

    return (
        <DashboardCard title={stat.label} className="flex-1">
            <div className="flex justify-between items-center">
                <p className="text-3xl font-bold text-brand-text-primary">{stat.value}</p>
                <div className={`flex items-center text-sm font-semibold ${changeColor}`}>
                    {isPositive && <ArrowUpIcon className="h-4 w-4 mr-1" />}
                    {isNegative && <ArrowDownIcon className="h-4 w-4 mr-1" />}
                    <span>{stat.change}</span>
                </div>
            </div>
        </DashboardCard>
    );
};

const alertIcons = {
    Critical: <BellAlertIcon className="h-5 w-5 text-brand-danger" />,
    Warning: <ExclamationTriangleIcon className="h-5 w-5 text-brand-warning" />,
    Info: <InformationCircleIcon className="h-5 w-5 text-brand-accent" />,
};

const CdnView: React.FC = () => {
    const [report, setReport] = useState<string>('');
    const [isLoadingReport, setIsLoadingReport] = useState<boolean>(true);
    const [stats, setStats] = useState<CdnStat[]>([]);
    const [alerts, setAlerts] = useState<CdnAlert[]>([]);
    const [isLoadingData, setIsLoadingData] = useState<boolean>(true);

    useEffect(() => {
        const fetchAllData = async () => {
            setIsLoadingReport(true);
            setIsLoadingData(true);
            const [reportData, statsData, alertsData] = await Promise.all([
                getCdnStatusReport(),
                fetchCdnStats(),
                fetchCdnAlerts(),
            ]);
            setReport(reportData);
            setStats(statsData);
            setAlerts(alertsData);
            setIsLoadingReport(false);
            setIsLoadingData(false);
        };
        fetchAllData();
    }, []);

    return (
        <div className="space-y-6">
            <div className="flex flex-col md:flex-row gap-6">
                {isLoadingData ? <LoadingSpinner /> : stats.map(stat => <StatCard key={stat.label} stat={stat} />)}
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className="lg:col-span-2">
                    <DashboardCard title="Global PoP Status (CDN v2.0)">
                        <div className="relative aspect-video bg-brand-bg rounded-md overflow-hidden border border-brand-border">
                            <WorldMap />
                            <DayNightOverlay />
                             <div className="absolute inset-0">
                                {Object.entries(pops).map(([id, pop]) => (
                                    <div key={id} className="absolute group" style={{ left: `${pop.x}%`, top: `${pop.y}%` }}>
                                        {pop.tier === 1 ? (
                                            <svg height="20" width="20" viewBox="0 0 20 20" className="-translate-x-1/2 -translate-y-1/2">
                                                <circle cx="10" cy="10" r="7" className="fill-brand-accent stroke-white dark:stroke-brand-surface stroke-2" />
                                                <circle cx="10" cy="10" r="7" className="fill-brand-accent animate-pulse" />
                                            </svg>
                                        ) : (
                                            <svg height="12" width="12" viewBox="0 0 12 12" className="-translate-x-1/2 -translate-y-1/2">
                                                <circle cx="6" cy="6" r="4" className="fill-green-400 stroke-brand-surface stroke-2" />
                                            </svg>
                                        )}
                                        <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 hidden group-hover:block px-2 py-1 bg-brand-surface text-brand-text-primary text-xs rounded-md shadow-lg whitespace-nowrap z-10 border border-brand-border">
                                            <strong className={pop.tier === 1 ? 'text-brand-accent' : 'text-green-400'}>Tier {pop.tier}</strong> - {pop.name}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </DashboardCard>
                </div>
                <div className="space-y-6">
                    <DashboardCard title="Recent Alerts">
                        {isLoadingData ? <LoadingSpinner /> : (
                            <ul className="space-y-3">
                                {alerts.map(alert => (
                                    <li key={alert.id} className="flex items-start space-x-3 text-sm p-2 rounded-md bg-brand-bg">
                                        <div className="flex-shrink-0 mt-0.5">{alertIcons[alert.severity]}</div>
                                        <div>
                                            <p className="text-brand-text-primary">{alert.message}</p>
                                            <p className="text-xs text-brand-text-secondary">{new Date(alert.timestamp).toLocaleTimeString()}</p>
                                        </div>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </DashboardCard>
                     <DashboardCard title="AI Status Report">
                        {isLoadingReport ? <LoadingSpinner /> : <GeminiResponseDisplay content={report} />}
                    </DashboardCard>
                </div>
            </div>
        </div>
    );
};

export default CdnView;