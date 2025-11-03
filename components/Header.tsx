
import React, { useState, useRef, useEffect } from 'react';
import { View, ViewState, Theme } from '../types';
import { ArrowRightOnRectangleIcon, UserCircleIcon, ChevronDownIcon, ChevronLeftIcon } from '@heroicons/react/24/solid';
import ThemeToggle from './ThemeToggle';

interface HeaderProps {
  currentViewState: ViewState;
  onLogout: () => void;
  setCurrentView: (view: View) => void;
  theme: Theme;
  setTheme: (theme: Theme) => void;
}

const viewTitles: Record<View, string> = {
  // Consumer Views
  streamverse_home: 'StreamVerse Home',
  watch: 'Now Playing',

  // Creator Views
  creator_studio: 'Creator Studio',

  // User Views
  user_profile: 'User Profile',

  // Admin Views
  roadmap: 'Project Implementation Roadmap',
  overview: 'Platform Overview',
  bi: 'Business Intelligence',
  ai_ops: 'AIOps Co-pilot',
  neural_engine: 'Neural Content Engine',
  media: 'Media Pipelines',
  cdn: 'Global CDN Topology',
  broadcast_ops: 'Broadcast Operations (DVB)',
  gpu_fabric: 'GPU Fabric Management',
  data: 'Data Platform',
  telecom: 'Telecom Core',
  satellite: 'Satellite Overlay',
  security: 'Security & Compliance',
  drm: 'DRM Management',
};

const Header: React.FC<HeaderProps> = ({ currentViewState, onLogout, setCurrentView, theme, setTheme }) => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  const showBackButton = ['watch'].includes(currentViewState.name);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <header className="flex-shrink-0 bg-brand-surface border-b border-brand-border px-4 sm:px-6 lg:px-8">
      <div className="flex items-center justify-between h-16">
        <div className="flex items-center space-x-4">
            {showBackButton && (
                <button 
                    onClick={() => setCurrentView('streamverse_home')} 
                    className="p-1 rounded-full text-brand-text-secondary hover:bg-brand-bg hover:text-brand-text-primary transition-colors"
                    aria-label="Go back"
                >
                    <ChevronLeftIcon className="h-6 w-6" />
                </button>
            )}
            <h1 className="text-xl font-semibold text-brand-text-primary">{viewTitles[currentViewState.name]}</h1>
        </div>
        <div className="flex items-center space-x-4">
          <ThemeToggle theme={theme} setTheme={setTheme} />
          <div className="relative" ref={menuRef}>
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="flex items-center space-x-2 text-brand-text-secondary hover:text-brand-text-primary transition-colors"
              aria-label="User menu"
            >
              <UserCircleIcon className="h-8 w-8" />
              <span className="text-sm hidden sm:inline">Platform Admin</span>
              <ChevronDownIcon className={`h-4 w-4 transition-transform ${isMenuOpen ? 'rotate-180' : ''}`} />
            </button>
            {isMenuOpen && (
              <div className="absolute right-0 mt-2 w-48 bg-brand-surface rounded-md shadow-lg border border-brand-border z-10">
                <div className="py-1">
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      setCurrentView('user_profile');
                      setIsMenuOpen(false);
                    }}
                    className="block px-4 py-2 text-sm text-brand-text-primary hover:bg-brand-bg"
                  >
                    Your Profile
                  </a>
                  <button
                    onClick={() => {
                      onLogout();
                      setIsMenuOpen(false);
                    }}
                    className="w-full text-left flex items-center px-4 py-2 text-sm text-brand-danger hover:bg-brand-bg"
                  >
                    <ArrowRightOnRectangleIcon className="h-5 w-5 mr-2" />
                    Sign Out
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
