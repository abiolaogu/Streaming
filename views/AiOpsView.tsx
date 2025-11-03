
import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { getAiOpsSummaryAndActions } from '../services/geminiService';
import { AiAction } from '../types';
import { LightBulbIcon, ExclamationTriangleIcon, ShieldCheckIcon, CogIcon } from '@heroicons/react/24/solid';

const mockActions: AiAction[] = [
    { id: 'cdn-01', category: 'CDN', description: 'High p95 latency detected in EU-West. Recommend deploying 2 additional ATS caches in FRA to handle load.', impact: 'High', status: 'Pending' },
    { id: 'media-01', category: 'Media', description: 'Transcode queue for 4K profiles has a lag of >500 jobs. Recommend scaling up RunPod GPU workers by 5 nodes.', impact: 'High', status: 'Pending' },
    { id: 'security-01', category: 'Security', description: 'Critical CVE-2024-XXXX found in DRM proxy. Recommend initiating immediate blue/green deployment to patch.', impact: 'High', status: 'Pending' },
    { id: 'cost-01', category: 'Cost', description: 'RunPod spot prices are peaking. Recommend shifting non-urgent VOD jobs to local Tier-1 GPUs for 4 hours.', impact: 'Medium', status: 'Pending' },
];

const categoryIcons = {
    CDN: <CogIcon className="h-6 w-6 text-blue-400" />,
    Media: <CogIcon className="h-6 w-6 text-purple-400" />,
    Security: <ShieldCheckIcon className="h-6 w-6 text-red-400" />,
    Cost: <ExclamationTriangleIcon className="h-6 w-6 text-yellow-400" />,
};

const impactColors = {
    Low: 'bg-green-500',
    Medium: 'bg-yellow-500',
    High: 'bg-red-500',
};

const ToggleSwitch: React.FC<{ checked: boolean, onChange: (e: React.ChangeEvent<HTMLInputElement>) => void }> = ({ checked, onChange }) => (
    <label className="relative inline-flex items-center cursor-pointer">
        <input type="checkbox" checked={checked} onChange={onChange} className="sr-only peer" />
        <div className="w-11 h-6 bg-gray-600 rounded-full peer peer-focus:ring-4 peer-focus:ring-blue-800 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-brand-accent"></div>
        <span className="ml-3 text-sm font-medium text-brand-text-primary">Autopilot Mode</span>
    </label>
);

const AiOpsView: React.FC = () => {
    const [summary, setSummary] = useState<string>('');
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [actions, setActions] = useState<AiAction[]>(mockActions);
    const [isAutopilotEngaged, setIsAutopilotEngaged] = useState(false);
    const [countdowns, setCountdowns] = useState<Record<string, number>>({});
    const [auditLog, setAuditLog] = useState<AiAction[]>([]);

    useEffect(() => {
        const fetchSummary = async () => {
            setIsLoading(true);
            const content = await getAiOpsSummaryAndActions();
            setSummary(content);
            setIsLoading(false);
        };
        fetchSummary();
    }, []);

    const updateActionAndLog = (id: string, newStatus: AiAction['status']) => {
        const action = actions.find(a => a.id === id);
        if (action) {
            const updatedAction = { ...action, status: newStatus, timestamp: new Date().toISOString() };
            setActions(prev => prev.map(a => a.id === id ? updatedAction : a));
            setAuditLog(prev => [updatedAction, ...prev]);
            return updatedAction;
        }
        return null;
    };

    useEffect(() => {
        if (!isAutopilotEngaged) {
            setCountdowns({});
            return;
        }

        const initialCountdowns: Record<string, number> = {};
        actions.forEach(action => {
            if (action.status === 'Pending') initialCountdowns[action.id] = 15;
        });
        setCountdowns(initialCountdowns);

        const timer = setInterval(() => {
            setCountdowns(prev => {
                const newCountdowns = { ...prev };
                let hasChanged = false;
                Object.keys(newCountdowns).forEach(id => {
                    if (newCountdowns[id] > 0) {
                        newCountdowns[id] -= 1;
                        hasChanged = true;
                    } else {
                        updateActionAndLog(id, 'Auto-Executed');
                        delete newCountdowns[id];
                        hasChanged = true;
                    }
                });
                return hasChanged ? newCountdowns : prev;
            });
        }, 1000);

        return () => clearInterval(timer);

    }, [isAutopilotEngaged]);

    const handleManualAction = (id: string, newStatus: 'Approved' | 'Denied') => {
        updateActionAndLog(id, newStatus);
    };

    const handleCancelAutopilot = (id: string) => {
        setCountdowns(prev => {
            const newCountdowns = { ...prev };
            delete newCountdowns[id];
            return newCountdowns;
        });
        updateActionAndLog(id, 'Canceled by Operator');
    };

    const formatTimestamp = (isoString?: string) => {
        if (!isoString) return '';
        return new Date(isoString).toLocaleTimeString();
    };


    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-1 space-y-6">
                <DashboardCard title="AIOps Co-pilot Summary">
                    {isLoading ? <LoadingSpinner /> : <GeminiResponseDisplay content={summary} />}
                </DashboardCard>
                <DashboardCard title="Autopilot Audit Log">
                    {auditLog.length === 0 ? (
                         <p className="text-brand-text-secondary text-center text-sm p-4">No autonomous actions taken yet.</p>
                    ) : (
                        <ul className="space-y-2 max-h-96 overflow-y-auto">
                           {auditLog.map(log => (
                               <li key={`${log.id}-${log.timestamp}`} className="text-xs p-2 bg-brand-bg rounded-md">
                                   <p className="font-mono text-brand-text-secondary">[{formatTimestamp(log.timestamp)}]</p>
                                   <p className="font-semibold text-brand-text-primary">{log.status}: <span className="font-normal text-brand-text-secondary">{log.description}</span></p>
                               </li>
                           ))}
                        </ul>
                    )}
                </DashboardCard>
            </div>
            <div className="lg:col-span-2">
                 <DashboardCard title="Actionable Insights & Recommendations">
                    <div className="flex justify-end mb-4 border-b border-brand-border pb-4">
                        <ToggleSwitch checked={isAutopilotEngaged} onChange={() => setIsAutopilotEngaged(!isAutopilotEngaged)} />
                    </div>
                    <div className="space-y-4">
                        {actions.map(action => (
                            <div key={action.id} className="p-4 bg-brand-bg rounded-lg border border-brand-border flex items-start space-x-4">
                                <div className="flex-shrink-0">{categoryIcons[action.category]}</div>
                                <div className="flex-1">
                                    <div className="flex justify-between items-center">
                                        <p className="text-sm font-semibold text-brand-text-primary">{action.category} Task</p>
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${impactColors[action.impact]}`}>{action.impact} Impact</span>
                                    </div>
                                    <p className="mt-1 text-sm text-brand-text-secondary">{action.description}</p>
                                </div>
                                <div className="flex flex-col space-y-2 w-32 text-center">
                                    {action.status === 'Pending' ? (
                                        isAutopilotEngaged ? (
                                            <>
                                                <p className="text-sm text-brand-accent animate-pulse">Auto-executing in {countdowns[action.id] || 15}s</p>
                                                <div className="w-full bg-brand-surface rounded-full h-1.5 my-1"><div className="bg-brand-accent h-1.5 rounded-full" style={{width: `${((countdowns[action.id] || 15) / 15) * 100}%`}}></div></div>
                                                <button onClick={() => handleCancelAutopilot(action.id)} className="text-xs bg-brand-danger text-white font-bold py-1 px-3 rounded-md hover:bg-red-600 transition">Cancel</button>
                                            </>
                                        ) : (
                                            <>
                                                <button onClick={() => handleManualAction(action.id, 'Approved')} className="text-xs bg-brand-success text-white font-bold py-1 px-3 rounded-md hover:bg-green-600 transition">Approve</button>
                                                <button onClick={() => handleManualAction(action.id, 'Denied')} className="text-xs bg-brand-danger text-white font-bold py-1 px-3 rounded-md hover:bg-red-600 transition">Deny</button>
                                            </>
                                        )
                                    ) : (
                                        <span className={`text-sm font-bold ${action.status === 'Approved' || action.status === 'Auto-Executed' ? 'text-brand-success' : 'text-brand-text-secondary'}`}>{action.status}</span>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                 </DashboardCard>
            </div>
        </div>
    );
};

export default AiOpsView;