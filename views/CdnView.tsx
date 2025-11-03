import React, { useEffect, useState } from 'react';
import DashboardCard from '../components/DashboardCard';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import LoadingSpinner from '../components/LoadingSpinner';
import { getCdnStatusReport } from '../services/geminiService';
import { fetchCdnStats, fetchCdnAlerts, fetchCacheHitRatioData, fetchLatencyByRegionData, fetchBandwidthByTierData } from '../services/mockApiService';
import DayNightOverlay from '../components/DayNightOverlay';
import WorldMap from '../components/WorldMap';
import { CdnStat, CdnAlert, CacheHitRatioDataPoint, LatencyByRegion, BandwidthByTier } from '../types';
import { ArrowUpIcon, ArrowDownIcon, BellAlertIcon, ExclamationTriangleIcon, InformationCircleIcon } from '@heroicons/react/24/solid';

const pops = {
    // Tier 1
    ash1: { x: 25.4, y: 32.5, tier: 1, name: 'Ashburn, USA', region: 'americas' as const, status: 'Active' },
    lon1: { x: 49.7, y: 23, tier: 1, name: 'London, UK', region: 'emea' as const, status: 'Active' },
    sgp1: { x: 79.8, y: 53.5, tier: 1, name: 'Singapore', region: 'apac' as const, status: 'Active' },
    sao1: { x: 35.5, y: 69.5, tier: 1, name: 'São Paulo, Brazil', region: 'americas' as const, status: 'Degraded' },
    los1: { x: 50.8, y: 49.5, tier: 1, name: 'Lagos, Nigeria', region: 'emea' as const, status: 'Active' },
    // Tier 2
    lax2: { x: 16.5, y: 35.5, tier: 2, name: 'Los Angeles, USA', region: 'americas' as const, status: 'Active' },
    mia2: { x: 24.3, y: 40, tier: 2, name: 'Miami, USA', region: 'americas' as const, status: 'Active' },
    fra2: { x: 52.4, y: 24.5, tier: 2, name: 'Frankfurt, DE', region: 'emea' as const, status: 'Active' },
    mum2: { x: 69.4, y: 42.5, tier: 2, name: 'Mumbai, India', region: 'apac' as const, status: 'Active' },
    syd2: { x: 91.5, y: 72.5, tier: 2, name: 'Sydney, Australia', region: 'apac' as const, status: 'Active' },
    tyo2: { x: 88, y: 34, tier: 2, name: 'Tokyo, Japan', region: 'apac' as const, status: 'Degraded' },
};

type Region = 'all' | 'americas' | 'emea' | 'apac';
type Tier = 'all' | 1 | 2;

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

const CacheHitRatioChart: React.FC<{ data: CacheHitRatioDataPoint[] }> = ({ data }) => {
    if (!data || data.length < 2) return <p className="text-brand-text-secondary text-center">Not enough data to display chart.</p>;

    const width = 500;
    const height = 150;
    const padding = 20;

    const ratios = data.map(p => p.ratio);
    const maxRatio = Math.ceil(Math.max(...ratios));
    const minRatio = Math.floor(Math.min(...ratios));
    const yRange = maxRatio - minRatio;

    const effectiveYRange = yRange === 0 ? 1 : yRange;
    const effectiveMinRatio = yRange === 0 ? minRatio - 0.5 : minRatio;

    const getX = (index: number) => padding + (index / (data.length - 1)) * (width - 2 * padding);
    const getY = (ratio: number) => height - padding - ((ratio - effectiveMinRatio) / effectiveYRange) * (height - 2 * padding);
    
    const points = data.map((point, i) => `${getX(i)},${getY(point.ratio)}`).join(' ');

    return (
        <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-auto" role="figure" aria-label="Cache hit ratio over the last 6 hours">
            {[maxRatio, minRatio].map(val => (
                <g key={`y-axis-${val}`}>
                    <text x="5" y={getY(val) + 5} className="text-xs fill-current text-brand-text-secondary">{val}%</text>
                    <line x1={padding} y1={getY(val)} x2={width - padding} y2={getY(val)} className="stroke-current text-brand-border" strokeWidth="0.5" strokeDasharray="2" />
                </g>
            ))}
            
            {data.map((point, i) => (
                <text key={point.time} x={getX(i)} y={height - 5} textAnchor="middle" className="text-xs fill-current text-brand-text-secondary">
                    {point.time}
                </text>
            ))}

            <polyline points={points} className="fill-none stroke-brand-accent" strokeWidth="2" />
            
            {data.map((point, i) => (
                 <g key={`g-${i}`} className="group" tabIndex={0} aria-label={`${point.time}: ${point.ratio}%`}>
                    <circle cx={getX(i)} cy={getY(point.ratio)} r="8" className="fill-transparent" />
                    <circle cx={getX(i)} cy={getY(point.ratio)} r="3" className="fill-brand-accent group-hover:r-4 transition-all" />
                    <title>{`${point.time}: ${point.ratio}%`}</title>
                 </g>
            ))}
        </svg>
    );
};

const LatencyByRegionChart: React.FC<{ data: LatencyByRegion[] }> = ({ data }) => {
    if (!data || data.length === 0) return null;
    const maxLatency = Math.max(...data.map(d => d.latency)) * 1.1;
    const regionColors: Record<LatencyByRegion['region'], string> = {
        'Americas': 'bg-blue-500',
        'EMEA': 'bg-green-500',
        'APAC': 'bg-purple-500',
    };

    return (
        <div className="space-y-3 pt-2">
            {data.map(item => (
                <div key={item.region} className="flex items-center">
                    <span className="w-20 text-sm text-brand-text-secondary">{item.region}</span>
                    <div className="flex-1 bg-brand-bg rounded-full h-4">
                        <div
                            className={`h-4 rounded-full ${regionColors[item.region]}`}
                            style={{ width: `${(item.latency / maxLatency) * 100}%` }}
                        />
                    </div>
                    <span className="w-16 text-right text-sm font-semibold text-brand-text-primary">{item.latency}ms</span>
                </div>
            ))}
        </div>
    );
};

const BandwidthByTierChart: React.FC<{ data: BandwidthByTier[] }> = ({ data }) => {
    if (!data || data.length === 0) return null;
    const totalBandwidth = data.reduce((sum, item) => sum + item.bandwidth, 0);
    
    return (
        <div className="space-y-4 pt-2">
            {data.map(item => (
                <div key={item.tier}>
                    <div className="flex justify-between items-center mb-1">
                         <span className="text-sm font-semibold text-brand-text-primary">{item.tier}</span>
                         <span className="text-sm text-brand-text-secondary">{item.bandwidth.toFixed(2)} Tbps</span>
                    </div>
                    <div className="w-full bg-brand-bg rounded-full h-5">
                        <div
                            className={`h-5 rounded-full ${item.tier === 'Tier 1' ? 'bg-brand-accent' : 'bg-green-400'}`}
                            style={{ width: `${(item.bandwidth / totalBandwidth) * 100}%` }}
                         />
                    </div>
                </div>
            ))}
        </div>
    );
};


const CdnView: React.FC = () => {
    const [report, setReport] = useState<string>('');
    const [isLoadingReport, setIsLoadingReport] = useState<boolean>(true);
    const [stats, setStats] = useState<CdnStat[]>([]);
    const [alerts, setAlerts] = useState<CdnAlert[]>([]);
    const [isLoadingData, setIsLoadingData] = useState<boolean>(true);
    const [cacheHitRatioData, setCacheHitRatioData] = useState<CacheHitRatioDataPoint[]>([]);
    const [latencyByRegion, setLatencyByRegion] = useState<LatencyByRegion[]>([]);
    const [bandwidthByTier, setBandwidthByTier] = useState<BandwidthByTier[]>([]);
    const [isLoadingAnalytics, setIsLoadingAnalytics] = useState<boolean>(true);
    
    const [regionFilter, setRegionFilter] = useState<Region>('all');
    const [tierFilter, setTierFilter] = useState<Tier>('all');

    useEffect(() => {
        const fetchStaticData = async () => {
            setIsLoadingReport(true);
            setIsLoadingAnalytics(true);
            const [reportData, cacheData, latencyData, bandwidthData] = await Promise.all([
                getCdnStatusReport(),
                fetchCacheHitRatioData(),
                fetchLatencyByRegionData(),
                fetchBandwidthByTierData(),
            ]);
            setReport(reportData);
            setCacheHitRatioData(cacheData);
            setLatencyByRegion(latencyData);
            setBandwidthByTier(bandwidthData);
            
            setIsLoadingReport(false);
            setIsLoadingAnalytics(false);
        };
        
        const refreshDynamicData = async () => {
            const [statsData, alertsData] = await Promise.all([
                fetchCdnStats(),
                fetchCdnAlerts(),
            ]);
            setStats(statsData);
            setAlerts(alertsData);
            setIsLoadingData(false);
        };

        fetchStaticData();
        refreshDynamicData(); // Initial fetch
        
        const intervalId = setInterval(refreshDynamicData, 30000); // Refresh every 30 seconds

        return () => clearInterval(intervalId); // Cleanup on unmount
    }, []);

    const getButtonClass = (isActive: boolean) => 
        `px-3 py-1 text-xs font-semibold rounded-full transition-colors duration-200 ${
            isActive 
            ? 'bg-brand-accent text-white' 
            : 'bg-brand-bg text-brand-text-secondary hover:bg-brand-surface hover:text-brand-text-primary'
        }`;

    const filteredPops = Object.entries(pops).filter(([, pop]) => {
        const regionMatch = regionFilter === 'all' || pop.region === regionFilter;
        const tierMatch = tierFilter === 'all' || pop.tier === tierFilter;
        return regionMatch && tierMatch;
    });

    const handleMapClick = (event: React.MouseEvent<HTMLDivElement>) => {
        const rect = event.currentTarget.getBoundingClientRect();
        const x = event.clientX - rect.left;
        const width = rect.width;
        const xPercent = (x / width) * 100;

        if (xPercent < 45) {
            setRegionFilter(prev => prev === 'americas' ? 'all' : 'americas');
        } else if (xPercent < 65) {
            setRegionFilter(prev => prev === 'emea' ? 'all' : 'emea');
        } else {
            setRegionFilter(prev => prev === 'apac' ? 'all' : 'apac');
        }
    };

    const getPopColor = (status: string) => status === 'Degraded' ? 'fill-yellow-400' : 'fill-green-400';
    const getTier1PopColor = (status: string) => status === 'Degraded' ? 'fill-yellow-400' : 'fill-brand-accent';

    return (
        <div className="space-y-6">
            <div className="flex flex-col md:flex-row gap-6">
                {isLoadingData ? <LoadingSpinner /> : stats.map(stat => <StatCard key={stat.label} stat={stat} />)}
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className="lg:col-span-2">
                    <DashboardCard title="Global PoP Status (CDN v2.0)">
                        <div className="flex flex-col sm:flex-row flex-wrap gap-4 justify-between items-center mb-4">
                            <div className="flex items-center space-x-2" role="group" aria-label="Filter by region">
                                <span className="text-sm font-medium text-brand-text-secondary">Region:</span>
                                <button onClick={() => setRegionFilter('all')} className={getButtonClass(regionFilter === 'all')}>All</button>
                                <button onClick={() => setRegionFilter('americas')} className={getButtonClass(regionFilter === 'americas')}>Americas</button>
                                <button onClick={() => setRegionFilter('emea')} className={getButtonClass(regionFilter === 'emea')}>EMEA</button>
                                <button onClick={() => setRegionFilter('apac')} className={getButtonClass(regionFilter === 'apac')}>APAC</button>
                            </div>
                            <div className="flex items-center space-x-2" role="group" aria-label="Filter by tier">
                                <span className="text-sm font-medium text-brand-text-secondary">Tier:</span>
                                <button onClick={() => setTierFilter('all')} className={getButtonClass(tierFilter === 'all')}>All</button>
                                <button onClick={() => setTierFilter(1)} className={getButtonClass(tierFilter === 1)}>Tier 1</button>
                                <button onClick={() => setTierFilter(2)} className={getButtonClass(tierFilter === 2)}>Tier 2</button>
                            </div>
                        </div>
                        <div 
                            className="relative aspect-video bg-brand-bg rounded-md overflow-hidden border border-brand-border cursor-pointer group"
                            onClick={handleMapClick}
                            role="button"
                            aria-label="Filter map by clicking on a region. Click again to clear."
                            tabIndex={0}
                            onKeyDown={(e) => { if(e.key === 'Enter' || e.key === ' ') handleMapClick(e as any); }}
                        >
                             <div className="absolute inset-0 z-10 hidden group-hover:flex items-center justify-center bg-black/30 transition-all duration-300">
                                <p className="text-white font-semibold text-lg">Click to toggle region filter</p>
                            </div>
                            <WorldMap activeRegion={regionFilter} />
                            <DayNightOverlay />
                             <div className="absolute inset-0">
                                {filteredPops.map(([id, pop]) => (
                                    <div key={id} className="absolute group/pop" style={{ left: `${pop.x}%`, top: `${pop.y}%` }}>
                                        {pop.tier === 1 ? (
                                            <svg height="20" width="20" viewBox="0 0 20 20" className="-translate-x-1/2 -translate-y-1/2">
                                                <circle cx="10" cy="10" r="7" className={`${getTier1PopColor(pop.status)} stroke-white dark:stroke-brand-surface stroke-2`} />
                                                {pop.status === 'Degraded' && <circle cx="10" cy="10" r="7" className="fill-yellow-400 animate-pulse" />}
                                            </svg>
                                        ) : (
                                            <svg height="12" width="12" viewBox="0 0 12 12" className="-translate-x-1/2 -translate-y-1/2">
                                                 <circle cx="6" cy="6" r="4" className={`${getPopColor(pop.status)} stroke-brand-surface stroke-2`} />
                                                 {pop.status === 'Degraded' && <circle cx="6" cy="6" r="4" className="fill-yellow-400 animate-pulse" />}
                                            </svg>
                                        )}
                                        <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 hidden group-hover/pop:block px-2 py-1 bg-brand-surface text-brand-text-primary text-xs rounded-md shadow-lg whitespace-nowrap z-20 border border-brand-border">
                                            <p><strong className={pop.tier === 1 ? 'text-brand-accent' : 'text-green-400'}>Tier {pop.tier}</strong> - {pop.name}</p>
                                            <p>Status: <span className={pop.status === 'Degraded' ? 'text-brand-warning' : 'text-brand-success'}>{pop.status}</span></p>
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
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <DashboardCard title="Cache Hit Ratio (Last 6 Hours)">
                    {isLoadingAnalytics ? <LoadingSpinner /> : <CacheHitRatioChart data={cacheHitRatioData} />}
                </DashboardCard>
                <DashboardCard title="Average Latency by Region">
                    {isLoadingAnalytics ? <LoadingSpinner /> : <LatencyByRegionChart data={latencyByRegion} />}
                </DashboardCard>
                <DashboardCard title="Bandwidth Utilization by Tier">
                    {isLoadingAnalytics ? <LoadingSpinner /> : <BandwidthByTierChart data={bandwidthByTier} />}
                </DashboardCard>
            </div>
        </div>
    );
};

export default CdnView;