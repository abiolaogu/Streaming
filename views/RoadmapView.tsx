
import React, { useState, useMemo } from 'react';
import { roadmapIssues } from '../data/roadmapData';
import { RoadmapIssue } from '../types';
import { ChevronDownIcon, ChevronUpIcon } from '@heroicons/react/24/solid';

const StatusPill: React.FC<{ status: RoadmapIssue['status'] }> = ({ status }) => {
    const colors = {
        'Done': 'bg-brand-success/20 text-brand-success',
        'In Progress': 'bg-brand-accent/20 text-brand-accent',
        'Not Started': 'bg-brand-secondary/20 text-brand-text-secondary',
    };
    return <span className={`px-2 py-1 text-xs font-semibold rounded-full ${colors[status]}`}>{status}</span>;
};

const PriorityPill: React.FC<{ priority: string }> = ({ priority }) => {
    const colors = {
        'P0': 'bg-red-500/80 text-white',
        'P1': 'bg-yellow-500/80 text-black',
    };
    const colorClass = priority === 'P0' ? colors.P0 : colors.P1;
    return <span className={`px-2 py-0.5 text-xs font-bold rounded-full ${colorClass}`}>{priority}</span>;
};


const RoadmapView: React.FC = () => {
    const [expandedIssue, setExpandedIssue] = useState<string | null>(null);

    const toggleIssue = (id: string) => {
        setExpandedIssue(expandedIssue === id ? null : id);
    };

    const phases = useMemo(() => {
        const grouped: Record<string, RoadmapIssue[]> = {};
        roadmapIssues.forEach(issue => {
            if (!grouped[issue.phase]) {
                grouped[issue.phase] = [];
            }
            grouped[issue.phase].push(issue);
        });
        return Object.entries(grouped);
    }, [roadmapIssues]);

    return (
        <div className="space-y-8">
            {phases.map(([phaseName, issues]) => (
                <div key={phaseName}>
                    <h2 className="text-2xl font-bold text-brand-text-primary mb-4 border-b border-brand-border pb-2">{phaseName}</h2>
                    <div className="space-y-2">
                        {issues.map(issue => (
                             <div key={issue.id} className="bg-brand-surface border border-brand-border rounded-lg">
                                <button
                                    onClick={() => toggleIssue(issue.id)}
                                    className="w-full p-4 text-left flex items-center justify-between"
                                >
                                    <div className="flex items-center gap-4">
                                        <div className="font-mono text-sm text-brand-text-secondary">{issue.id}</div>
                                        <div className="font-semibold text-brand-text-primary">{issue.title}</div>
                                    </div>
                                    <div className="flex items-center gap-4">
                                        <PriorityPill priority={issue.priority} />
                                        <StatusPill status={issue.status} />
                                        {expandedIssue === issue.id ? <ChevronUpIcon className="h-5 w-5"/> : <ChevronDownIcon className="h-5 w-5"/>}
                                    </div>
                                </button>
                                {expandedIssue === issue.id && (
                                    <div className="p-4 border-t border-brand-border bg-brand-bg">
                                        <div className="grid grid-cols-3 gap-4 text-sm mb-4">
                                            <div><span className="font-semibold">Owner:</span> {issue.owner}</div>
                                            <div><span className="font-semibold">Effort:</span> {issue.effort}</div>
                                            <div><span className="font-semibold">Priority:</span> {issue.priority}</div>
                                        </div>
                                        <p className="text-sm text-brand-text-secondary whitespace-pre-wrap">{issue.description}</p>
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>
                </div>
            ))}
        </div>
    );
};

export default RoadmapView;
