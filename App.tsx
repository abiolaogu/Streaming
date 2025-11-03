
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
import RoadmapView from './views/RoadmapView';
import { View, ViewState, Theme } from './types';
import ApiKeyPromptView from './views/ApiKeyPromptView';
import LoadingSpinner from './components/LoadingSpinner';

const App: React.FC = () => {
  const [currentViewState, setCurrentViewState] = useState<ViewState>({ name: 'streamverse_home' });
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [theme, setTheme] = useState<Theme>(() => (localStorage.getItem('theme') as Theme) || 'dark');
  
  const [isKeySelected, setIsKeySelected] = useState(false);
  const [hasCheckedKey, setHasCheckedKey] = useState(false);

  useEffect(() => {
    if (theme === 'light') {
      document.documentElement.classList.remove('dark');
    } else {
      document.documentElement.classList.add('dark');
    }
    localStorage.setItem('theme', theme);
  }, [theme]);

  useEffect(() => {
    const checkKey = async () => {
      // Check for local dev environment variable first
      if (import.meta.env.VITE_GEMINI_API_KEY) {
        setIsKeySelected(true);
      } 
      // If no env var, check for AI Studio environment
      else if ((window as any).aistudio && await (window as any).aistudio.hasSelectedApiKey()) {
        setIsKeySelected(true);
      }
      setHasCheckedKey(true);
    };
    checkKey();

    const handleApiKeyError = () => {
        console.warn('API key error detected. Prompting for new key.');
        // For local dev, this indicates the key is bad. For AI Studio, it will re-trigger the prompt.
        if (!import.meta.env.VITE_GEMINI_API_KEY) {
          setIsKeySelected(false);
        }
    };

    window.addEventListener('apiKeyError', handleApiKeyError);

    return () => {
        window.removeEventListener('apiKeyError', handleApiKeyError);
    };
  }, []);

  const handleKeySelected = () => {
    setIsKeySelected(true);
  };

  const handleLogin = () => {
    setIsAuthenticated(true);
    setCurrentViewState({ name: 'roadmap' }); // Default to roadmap after login
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
      case 'roadmap':
        return <RoadmapView />;
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

  if (!hasCheckedKey) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-brand-bg">
        <LoadingSpinner />
      </div>
    );
  }

  // Only show the prompt if no key is selected and we are not in a local dev environment with a key
  if (!isKeySelected) {
    return <ApiKeyPromptView onKeySelected={handleKeySelected} />;
  }

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
