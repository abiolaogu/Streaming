
import React from 'react';

interface DashboardCardProps {
  title: string;
  children: React.ReactNode;
  className?: string;
}

const DashboardCard: React.FC<DashboardCardProps> = ({ title, children, className = '' }) => {
  return (
    <div className={`bg-brand-surface rounded-lg border border-brand-border shadow-md ${className}`}>
      <div className="p-4 border-b border-brand-border">
        <h3 className="text-lg font-semibold text-brand-text-primary">{title}</h3>
      </div>
      <div className="p-4">
        {children}
      </div>
    </div>
  );
};

export default DashboardCard;
