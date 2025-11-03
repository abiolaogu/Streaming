import React, { useState, useEffect } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { fetchChurnRiskProfiles, fetchContentRoi } from '../services/mockApiService';
import { getBiSummary } from '../services/geminiService';
import { ChurnRiskProfile, ContentRoiPrediction } from '../types';

const BusinessIntelligenceView: React.FC = () => {
    const [summary, setSummary] = useState('');
    const [isLoadingSummary, setIsLoadingSummary] = useState(true);
    const [churnRisks, setChurnRisks] = useState<ChurnRiskProfile[]>([]);
    const [contentRoi, setContentRoi] = useState<ContentRoiPrediction[]>([]);
    const [isLoadingData, setIsLoadingData] = useState(true);

    useEffect(() => {
        const loadAllData = async () => {
            setIsLoadingSummary(true);
            setIsLoadingData(true);
            const [summaryData, churnData, roiData] = await Promise.all([
                getBiSummary(),
                fetchChurnRiskProfiles(),
                fetchContentRoi(),
            ]);
            setSummary(summaryData);
            setChurnRisks(churnData);
            setContentRoi(roiData);
            setIsLoadingSummary(false);
            setIsLoadingData(false);
        };
        loadAllData();
    }, []);

    const getRiskColor = (score: number) => {
        if (score > 90) return 'text-brand-danger';
        if (score > 70) return 'text-brand-warning';
        return 'text-brand-text-secondary';
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 space-y-6">
                <DashboardCard title="Subscriber Churn Risk (Top 3)">
                    {isLoadingData ? <LoadingSpinner/> : (
                         <ul className="space-y-3">
                            {churnRisks.map(profile => (
                                <li key={profile.subscriberId} className="p-3 bg-brand-bg rounded-lg">
                                    <div className="flex justify-between items-center">
                                        <p className="font-mono text-sm text-brand-text-primary">{profile.subscriberId}</p>
                                        <p className={`text-xl font-bold ${getRiskColor(profile.riskScore)}`}>{profile.riskScore}% Risk</p>
                                    </div>
                                    <p className="text-xs text-brand-text-secondary mt-1"><strong>Reason:</strong> {profile.primaryReason}</p>
                                    <p className="text-xs text-brand-accent mt-1"><strong>Suggested Offer:</strong> {profile.recommendedOffer}</p>
                                </li>
                            ))}
                        </ul>
                    )}
                </DashboardCard>
                 <DashboardCard title="Content Investment AI (Top Opportunities)">
                    {isLoadingData ? <LoadingSpinner/> : (
                         <div className="overflow-x-auto">
                            <table className="min-w-full divide-y divide-brand-border">
                                <thead className="bg-brand-bg">
                                    <tr>
                                        <th className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase">Content Title</th>
                                        <th className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase">Predicted Viewership</th>
                                        <th className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase">Predicted ROI</th>
                                        <th className="px-4 py-2 text-left text-xs font-medium text-brand-text-secondary uppercase">Confidence</th>
                                    </tr>
                                </thead>
                                <tbody className="bg-brand-surface divide-y divide-brand-border">
                                    {contentRoi.map(item => (
                                        <tr key={item.contentTitle}>
                                            <td className="px-4 py-3 text-sm font-medium text-brand-text-primary">{item.contentTitle}</td>
                                            <td className="px-4 py-3 text-sm text-brand-text-secondary">{item.predictedViewership.toLocaleString()}</td>
                                            <td className="px-4 py-3 text-sm font-bold text-brand-success">{item.predictedRoi}x</td>
                                            <td className="px-4 py-3 text-sm text-brand-text-secondary">{item.confidence}</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </DashboardCard>
            </div>
            <div className="lg:col-span-1">
                <DashboardCard title="BI Summary (AI Generated)">
                    {isLoadingSummary ? <LoadingSpinner /> : <GeminiResponseDisplay content={summary} />}
                </DashboardCard>
            </div>
        </div>
    );
};

export default BusinessIntelligenceView;