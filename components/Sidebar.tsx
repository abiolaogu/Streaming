
import React, { useState } from 'react';
import { NavItem, View } from '../types';
import { 
    OverviewIcon, GpuIcon, CdnIcon, SecurityIcon, MediaIcon, 
    DataIcon, TelecomIcon, SatelliteIcon, AiOpsIcon, NeuralEngineIcon,
    HomeIcon, CreatorStudioIcon, SparklesIcon, DrmIcon,
    BusinessIntelligenceIcon, BroadcastIcon, RoadmapIcon
} from './IconComponents';
import { ChevronDownIcon } from '@heroicons/react/24/solid';

interface SidebarProps {
  currentView: View;
  setCurrentView: (view: View) => void;
}

const consumerNavItems: NavItem[] = [
    { id: 'streamverse_home', label: 'Home', icon: HomeIcon },
];

const creatorNavItems: NavItem[] = [
    { id: 'creator_studio', label: 'Creator Studio', icon: CreatorStudioIcon },
];

const adminNavItems: NavItem[] = [
  { id: 'roadmap', label: 'Project Roadmap', icon: RoadmapIcon },
  { id: 'overview', label: 'Overview', icon: OverviewIcon },
  { id: 'bi', label: 'Business Intelligence', icon: BusinessIntelligenceIcon },
  { id: 'ai_ops', label: 'AIOps Co-pilot', icon: AiOpsIcon },
  { id: 'neural_engine', label: 'Neural Engine', icon: NeuralEngineIcon },
  { id: 'media', label: 'Media Pipeline', icon: MediaIcon },
  { id: 'cdn', label: 'CDN', icon: CdnIcon },
  { id: 'broadcast_ops', label: 'Broadcast Ops', icon: BroadcastIcon },
  { id: 'gpu_fabric', label: 'GPU Fabric', icon: GpuIcon },
  { id: 'data', label: 'Data Platform', icon: DataIcon },
  { id: 'telecom', label: 'Telecom Core', icon: TelecomIcon },
  { id: 'satellite', label: 'Satellite Overlay', icon: SatelliteIcon },
  { id: 'security', label: 'Security', icon: SecurityIcon },
  { id: 'drm', label: 'DRM Management', icon: DrmIcon },
];

const Sidebar: React.FC<SidebarProps> = ({ currentView, setCurrentView }) => {
  const [isAdminOpen, setIsAdminOpen] = useState(true);

  const NavLink: React.FC<{item: NavItem, isChild?: boolean}> = ({ item, isChild = false }) => (
    <a
      key={item.id}
      href="#"
      onClick={(e) => {
        e.preventDefault();
        setCurrentView(item.id);
      }}
      className={`flex items-center px-4 py-2 text-sm font-medium rounded-md transition-colors duration-200 ${
        currentView === item.id
          ? 'bg-brand-accent text-white'
          : 'text-brand-text-secondary hover:bg-brand-surface hover:text-brand-text-primary'
      } ${isChild ? 'pl-11' : ''}`}
    >
      <item.icon className="h-5 w-5 mr-3" />
      {item.label}
    </a>
  );

  return (
    <div className="flex flex-col w-64 bg-brand-bg border-r border-brand-border">
      <div className="flex items-center justify-center h-16 border-b border-brand-border">
        <SparklesIcon className="h-8 w-8 text-brand-accent"/>
        <span className="ml-2 text-xl font-bold text-brand-text-primary tracking-wider">StreamVerse</span>
      </div>
      <div className="flex-1 overflow-y-auto">
        <nav className="flex-1 px-2 py-4 space-y-1">
          {consumerNavItems.map((item) => <NavLink key={item.id} item={item} />)}
          
          <div className="pt-2">
             <p className="px-4 pt-2 pb-1 text-xs font-semibold text-brand-text-secondary uppercase">For Creators</p>
             {creatorNavItems.map((item) => <NavLink key={item.id} item={item} />)}
          </div>
          
          <div className="pt-2">
            <button 
              onClick={() => setIsAdminOpen(!isAdminOpen)}
              className="w-full flex items-center justify-between px-4 py-2 text-xs font-semibold text-brand-text-secondary uppercase transition-colors duration-200 hover:bg-brand-surface rounded-md"
            >
              Platform Admin
              <ChevronDownIcon className={`w-4 h-4 transition-transform ${isAdminOpen ? 'rotate-180' : ''}`} />
            </button>
            {isAdminOpen && (
              <div className="mt-1 space-y-1">
                {adminNavItems.map((item) => <NavLink key={item.id} item={item} isChild />)}
              </div>
            )}
          </div>
        </nav>
      </div>
      <div className="p-4 border-t border-brand-border text-xs text-brand-text-secondary">
        <p>© {new Date().getFullYear()} StreamVerse</p>
        <p>Version 3.0.0</p>
      </div>
    </div>
  );
};

export default Sidebar;
