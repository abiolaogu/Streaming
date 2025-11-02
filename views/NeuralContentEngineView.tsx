import React, { useState, useRef } from 'react';
import DashboardCard from '../components/DashboardCard';
import LoadingSpinner from '../components/LoadingSpinner';
import GeminiResponseDisplay from '../components/GeminiResponseDisplay';
import { generateImage, analyzeVideo, transcribeAudio } from '../services/geminiService';
import { PhotoIcon, VideoCameraIcon, MicrophoneIcon, StopIcon, DocumentArrowUpIcon } from '@heroicons/react/24/solid';

const NeuralContentEngineView: React.FC = () => {
    // Image Generation State
    const [imagePrompt, setImagePrompt] = useState<string>('A photorealistic image of a futuristic broadcast studio on Mars');
    const [aspectRatio, setAspectRatio] = useState<string>('16:9');
    const [generatedImage, setGeneratedImage] = useState<string>('');
    const [isGeneratingImage, setIsGeneratingImage] = useState<boolean>(false);

    // Video Analysis State
    const [videoFile, setVideoFile] = useState<File | null>(null);
    const [videoPrompt, setVideoPrompt] = useState<string>('Analyze this video scene-by-scene. Identify emotional tone, objects, characters, and key actions.');
    const [videoAnalysis, setVideoAnalysis] = useState<string>('');
    const [isAnalyzingVideo, setIsAnalyzingVideo] = useState<boolean>(false);
    const videoInputRef = useRef<HTMLInputElement>(null);

    // Audio Transcription (Media File) State
    const [mediaFile, setMediaFile] = useState<File | null>(null);
    const [transcriptionEngine, setTranscriptionEngine] = useState<'on-premise' | 'gemini'>('on-premise');
    const [mediaTranscription, setMediaTranscription] = useState<string>('');
    const [isTranscribingMedia, setIsTranscribingMedia] = useState<boolean>(false);
    const mediaInputRef = useRef<HTMLInputElement>(null);

    const handleGenerateImage = async () => {
        setIsGeneratingImage(true);
        setGeneratedImage('');
        const result = await generateImage(imagePrompt, aspectRatio);
        if (result.startsWith('Error:')) {
            alert(result);
        } else {
            setGeneratedImage(`data:image/jpeg;base64,${result}`);
        }
        setIsGeneratingImage(false);
    };

    const handleAnalyzeVideo = async () => {
        if (!videoFile) {
            alert('Please select a video file first.');
            return;
        }
        setIsAnalyzingVideo(true);
        setVideoAnalysis('');
        
        const reader = new FileReader();
        reader.onloadend = async () => {
            const base64Data = (reader.result as string).split(',')[1];
            const result = await analyzeVideo(videoPrompt, base64Data, videoFile.type);
            setVideoAnalysis(result);
            setIsAnalyzingVideo(false);
        };
        reader.readAsDataURL(videoFile);
    };
    
    const handleTranscribeMediaFile = async () => {
        if (!mediaFile) {
            alert('Please select a media file first.');
            return;
        }
        setIsTranscribingMedia(true);
        setMediaTranscription('');

        if (transcriptionEngine === 'on-premise') {
            setTimeout(() => {
                setMediaTranscription(`[Simulated On-Premise Result]\nTranscription successful for file: ${mediaFile.name}.\nThis process ran on local GPU resources within our secure network.`);
                setIsTranscribingMedia(false);
            }, 2000);
        } else {
            const reader = new FileReader();
            reader.onloadend = async () => {
                const base64Data = (reader.result as string).split(',')[1];
                const result = await transcribeAudio(base64Data, mediaFile.type);
                setMediaTranscription(result);
                setIsTranscribingMedia(false);
            };
            reader.readAsDataURL(mediaFile);
        }
    };

    return (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <DashboardCard title="Generative Imagery (Imagen 4)" className="lg:col-span-2">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 items-start">
                    <div className="space-y-4">
                        <div>
                            <label className="block text-sm font-medium text-brand-text-secondary mb-1">Prompt</label>
                            <textarea value={imagePrompt} onChange={(e) => setImagePrompt(e.target.value)} rows={3} className="w-full bg-brand-bg border border-brand-border rounded-md p-2 text-sm" />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-brand-text-secondary mb-1">Aspect Ratio</label>
                            <select value={aspectRatio} onChange={(e) => setAspectRatio(e.target.value)} className="w-full bg-brand-bg border border-brand-border rounded-md p-2 text-sm">
                                <option>16:9</option><option>9:16</option><option>4:3</option><option>3:4</option><option>1:1</option>
                            </select>
                        </div>
                        <button onClick={handleGenerateImage} disabled={isGeneratingImage} className="w-full flex justify-center items-center gap-2 bg-brand-accent hover:bg-brand-accent-hover text-white font-bold py-2 px-4 rounded-md transition disabled:bg-gray-500">
                           <PhotoIcon className="h-5 w-5" /> {isGeneratingImage ? 'Generating...' : 'Generate Image'}
                        </button>
                    </div>
                    <div className="aspect-video bg-brand-bg rounded-md border border-dashed border-brand-border flex items-center justify-center">
                        {isGeneratingImage ? <LoadingSpinner /> : (generatedImage ? <img src={generatedImage} alt="Generated content" className="w-full h-full object-contain" /> : <p className="text-brand-text-secondary">Image will appear here</p>)}
                    </div>
                </div>
            </DashboardCard>
            
            <DashboardCard title="Deep Content Understanding (Video Analysis)">
                 <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-brand-text-secondary mb-1">Video File</label>
                         <button onClick={() => videoInputRef.current?.click()} className="w-full p-4 border-2 border-dashed border-brand-border rounded-md text-brand-text-secondary hover:bg-brand-bg">
                            {videoFile ? `Selected: ${videoFile.name}` : 'Click to select a video file'}
                        </button>
                        <input type="file" accept="video/*" ref={videoInputRef} onChange={(e) => setVideoFile(e.target.files ? e.target.files[0] : null)} className="hidden" />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-brand-text-secondary mb-1">Analysis Prompt</label>
                        <textarea value={videoPrompt} onChange={(e) => setVideoPrompt(e.target.value)} rows={2} className="w-full bg-brand-bg border border-brand-border rounded-md p-2 text-sm" />
                    </div>
                    <button onClick={handleAnalyzeVideo} disabled={isAnalyzingVideo || !videoFile} className="w-full flex justify-center items-center gap-2 bg-brand-accent hover:bg-brand-accent-hover text-white font-bold py-2 px-4 rounded-md transition disabled:bg-gray-500">
                        <VideoCameraIcon className="h-5 w-5" /> {isAnalyzingVideo ? 'Analyzing...' : 'Analyze Video'}
                    </button>
                    <div className="mt-4 p-4 min-h-[150px] bg-brand-bg rounded-md border border-dashed border-brand-border">
                       {isAnalyzingVideo ? <LoadingSpinner /> : (videoAnalysis ? <GeminiResponseDisplay content={videoAnalysis} /> : <p className="text-brand-text-secondary">Analysis will appear here.</p>)}
                    </div>
                </div>
            </DashboardCard>

            <DashboardCard title="Multi-Engine Transcription">
                 <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-brand-text-secondary mb-1">Media File (VOD, Live Stream)</label>
                         <button onClick={() => mediaInputRef.current?.click()} className="w-full p-4 border-2 border-dashed border-brand-border rounded-md text-brand-text-secondary hover:bg-brand-bg">
                            {mediaFile ? `Selected: ${mediaFile.name}` : 'Click to select an audio or video file'}
                        </button>
                        <input type="file" accept="audio/*,video/*" ref={mediaInputRef} onChange={(e) => setMediaFile(e.target.files ? e.target.files[0] : null)} className="hidden" />
                    </div>
                     <div>
                        <label className="block text-sm font-medium text-brand-text-secondary mb-1">Transcription Engine</label>
                        <select value={transcriptionEngine} onChange={(e) => setTranscriptionEngine(e.target.value as 'gemini' | 'on-premise')} className="w-full bg-brand-bg border border-brand-border rounded-md p-2 text-sm">
                            <option value="on-premise">On-Premise Engine (Secure & Cost-Effective)</option>
                            <option value="gemini">Cloud AI (Highest Accuracy)</option>
                        </select>
                    </div>
                    <button onClick={handleTranscribeMediaFile} disabled={isTranscribingMedia || !mediaFile} className="w-full flex justify-center items-center gap-2 bg-brand-accent hover:bg-brand-accent-hover text-white font-bold py-2 px-4 rounded-md transition disabled:bg-gray-500">
                       <DocumentArrowUpIcon className="h-5 w-5" /> {isTranscribingMedia ? 'Transcribing...' : 'Transcribe Media File'}
                    </button>
                    <div className="mt-4 p-4 min-h-[150px] bg-brand-bg rounded-md border border-dashed border-brand-border">
                        {isTranscribingMedia ? <LoadingSpinner /> : (mediaTranscription ? <p className="text-brand-text-primary whitespace-pre-wrap text-sm">{mediaTranscription}</p> : <p className="text-brand-text-secondary text-center">Transcription will appear here.</p>)}
                    </div>
                </div>
            </DashboardCard>
        </div>
    );
};

export default NeuralContentEngineView;