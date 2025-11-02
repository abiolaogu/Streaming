
import React, { useEffect, useState } from 'react';
import DashboardCard from '../components/DashboardCard';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import LoadingSpinner from '../components/LoadingSpinner';
import { getSecuritySummary } from '../services/geminiService';
import { ShieldCheckIcon, ExclamationTriangleIcon } from '@heroicons/react/24/solid';


const SecurityView: React.FC = () => {
    const [summary, setSummary] = useState<string>('');
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchSummary = async () => {
            setIsLoading(true);
            const content = await getSecuritySummary();
            setSummary(content);
            setIsLoading(false);
        };
        fetchSummary();
    }, []);
    
    const mockScanResults = [
        { service: 'media/transcoder:v1.2', vulnerabilities: { critical: 1, high: 3}, signed: true },
        { service: 'cdn/ats-cache:v3.1', vulnerabilities: { critical: 0, high: 0}, signed: true },
        { service: 'control/autoscaler:v0.9', vulnerabilities: { critical: 0, high: 1}, signed: true },
        { service: 'data/scylla-proxy:latest', vulnerabilities: { critical: 2, high: 5}, signed: false },
    ];

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-1">
                <DashboardCard title="DevSecOps Summary (AI Generated)">
                     {isLoading ? <LoadingSpinner /> : <GeminiResponseDisplay content={summary} />}
                </DashboardCard>
            </div>
            <div className="lg:col-span-2">
                 <DashboardCard title="Container Image Scan Results (Trivy + Cosign)">
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-brand-border">
                            <thead className="bg-brand-bg">
                                <tr>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Service Image</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Critical Vulns</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">High Vulns</th>
                                    <th scope="col" className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase tracking-wider">Signature</th>
                                </tr>
                            </thead>
                            <tbody className="bg-brand-surface divide-y divide-brand-border">
                                {mockScanResults.map((result) => (
                                    <tr key={result.service}>
                                        <td className="px-4 py-2 whitespace-nowrap text-sm text-brand-text-primary font-mono">{result.service}</td>
                                        <td className={`px-4 py-2 whitespace-nowrap text-sm text-center font-bold ${result.vulnerabilities.critical > 0 ? 'text-brand-danger' : 'text-brand-success'}`}>{result.vulnerabilities.critical}</td>
                                        <td className={`px-4 py-2 whitespace-nowrap text-sm text-center font-bold ${result.vulnerabilities.high > 0 ? 'text-brand-warning' : 'text-brand-success'}`}>{result.vulnerabilities.high}</td>
                                        <td className="px-4 py-2 whitespace-nowrap text-sm">
                                            {result.signed ? 
                                                <span className="flex items-center text-brand-success"><ShieldCheckIcon className="h-5 w-5 mr-1" /> Verified</span> : 
                                                <span className="flex items-center text-brand-danger"><ExclamationTriangleIcon className="h-5 w-5 mr-1"/> Unsigned</span>
                                            }
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                 </DashboardCard>
            </div>
        </div>
    );
};

export default SecurityView;
