import React, { useState, useCallback, useEffect } from 'react';
import Sidebar from './components/Sidebar';
import Header from './components/Header';
import OverviewView from './views/OverviewView';
import GpuFabricView from './views/GpuFabricView';
import CdnView from './views/CdnView';
import SecurityView from './views/SecurityView';
import MediaView from './views/MediaView';
import StreamVerseHomeView from './views/StreamVerseHomeView';
import AiOpsView from './views/AiOpsView';
import NeuralContentEngineView from './views/NeuralContentEngineView';
import ChatBot from './components/ChatBot';
import DataView from './views/DataView';
import TelecomView from './views/TelecomView';
import SatelliteView from './views/SatelliteView';
import DrmView from './views/DrmView';
import CreatorStudioView from './views/CreatorStudioView';
import BusinessIntelligenceView from './views/BusinessIntelligenceView';
import BroadcastOpsView from './views/BroadcastOpsView';
import LoginView from './views/LoginView';
import UserProfileView from './views/UserProfileView';
import WatchView from './views/WatchView';
import { View, ViewState, Theme } from './types';

const App: React.FC = () => {
  const [currentViewState, setCurrentViewState] = useState<ViewState>({ name: 'streamverse_home' });
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [theme, setTheme] = useState<Theme>(() => (localStorage.getItem('theme') as Theme) || 'dark');

  useEffect(() => {
    if (theme === 'light') {
      document.documentElement.classList.remove('dark');
    } else {
      document.documentElement.classList.add('dark');
    }
    localStorage.setItem('theme', theme);
  }, [theme]);


  const handleLogin = () => {
    setIsAuthenticated(true);
    setCurrentViewState({ name: 'overview' }); // Default to overview after login
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    setCurrentViewState({ name: 'streamverse_home' });
  };

  const handleSetView = (view: View) => {
    setCurrentViewState({ name: view });
  };

  const handleWatchContent = (contentId: string) => {
    setCurrentViewState({ name: 'watch', params: { contentId } });
  };

  const renderView = useCallback(() => {
    switch (currentViewState.name) {
      // Consumer Views
      case 'streamverse_home':
        return <StreamVerseHomeView onWatch={handleWatchContent} />;
      case 'watch':
        return <WatchView contentId={currentViewState.params?.contentId || ''} />;
      
      // Creator View
      case 'creator_studio':
        return <CreatorStudioView />;

      // User View
      case 'user_profile':
        return <UserProfileView />;

      // Admin Views
      case 'overview':
        return <OverviewView />;
      case 'gpu_fabric':
        return <GpuFabricView />;
      case 'cdn':
        return <CdnView />;
      case 'security':
        return <SecurityView />;
      case 'media':
        return <MediaView />;
      case 'data':
        return <DataView />;
      case 'telecom':
        return <TelecomView />;
      case 'satellite':
        return <SatelliteView />;
      case 'drm':
        return <DrmView />;
      case 'ai_ops':
        return <AiOpsView />;
      case 'neural_engine':
        return <NeuralContentEngineView />;
      case 'bi':
        return <BusinessIntelligenceView />;
      case 'broadcast_ops':
        return <BroadcastOpsView />;
      default:
        return <StreamVerseHomeView onWatch={handleWatchContent} />;
    }
  }, [currentViewState]);

  if (!isAuthenticated) {
    return <LoginView onLogin={handleLogin} />;
  }

  return (
    <div className="flex h-screen bg-brand-bg font-sans">
      <Sidebar currentView={currentViewState.name} setCurrentView={handleSetView} />
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header 
          currentViewState={currentViewState} 
          onLogout={handleLogout} 
          setCurrentView={handleSetView}
          theme={theme}
          setTheme={setTheme}
        />
        <main className="flex-1 overflow-x-hidden overflow-y-auto bg-brand-bg p-4 sm:p-6 lg:p-8">
          {renderView()}
        </main>
      </div>
      <ChatBot />
    </div>
  );
};

export default App;