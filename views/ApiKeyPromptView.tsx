
import React from 'react';
import { SparklesIcon } from '../components/IconComponents';

interface ApiKeyPromptViewProps {
  onKeySelected: () => void;
}

const ApiKeyPromptView: React.FC<ApiKeyPromptViewProps> = ({ onKeySelected }) => {
  const handleSelectKey = async () => {
    await (window as any).aistudio.openSelectKey();
    onKeySelected(); // Optimistically update the UI to show the app
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-brand-bg p-4">
      <div className="w-full max-w-lg p-8 text-center bg-brand-surface rounded-2xl border border-brand-border shadow-2xl">
        <SparklesIcon className="h-16 w-16 text-brand-accent mx-auto mb-4"/>
        <h1 className="text-2xl font-bold text-brand-text-primary mb-2">Welcome to StreamVerse Platform</h1>
        <p className="text-brand-text-secondary mb-6">To use the AI-powered features of this application, you need to select an API key. Your key is stored securely and is not shared.</p>
        <button
          onClick={handleSelectKey}
          className="w-full bg-brand-accent hover:bg-brand-accent-hover text-white font-bold py-3 px-6 rounded-md transition duration-200"
        >
          Select API Key
        </button>
        <p className="text-xs text-brand-text-secondary mt-4">
          For information on pricing, please see the <a href="https://ai.google.dev/gemini-api/docs/billing" target="_blank" rel="noopener noreferrer" className="text-brand-accent hover:underline">billing documentation</a>.
        </p>
      </div>
    </div>
  );
};

export default ApiKeyPromptView;
