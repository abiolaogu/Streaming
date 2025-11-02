import React, { useState, useRef, useEffect } from 'react';
import { ChatIcon } from './IconComponents';
import { PaperAirplaneIcon, XMarkIcon } from '@heroicons/react/24/solid';
import { ChatMessage } from '../types';
import { startChat } from '../services/geminiService';
import GeminiResponseDisplay from './GeminiResponseDisplay';
import { Chat } from '@google/genai';

const ChatBot: React.FC = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [messages, setMessages] = useState<ChatMessage[]>([
        { role: 'model', content: "Hi, I'm Vera, your Virtual Entertainment Assistant! How can I make your StreamVerse experience amazing today? You can ask me to find movies, explain platform features, or even help with technical questions." }
    ]);
    const [input, setInput] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const chatRef = useRef<Chat | null>(null);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);
    
    useEffect(() => {
        if(isOpen) {
             const history = messages.slice(0, -1).map(msg => ({
                role: msg.role,
                parts: [{ text: msg.content }]
             })) as { role: "user" | "model"; parts: { text: string }[] }[];
             chatRef.current = startChat(history);
        }
    }, [isOpen]);

    const handleSendMessage = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!input.trim() || isLoading) return;

        const userMessage: ChatMessage = { role: 'user', content: input };
        setMessages(prev => [...prev, userMessage]);
        setInput('');
        setIsLoading(true);

        try {
            if (!chatRef.current) {
                 const history = messages.map(msg => ({
                    role: msg.role,
                    parts: [{ text: msg.content }]
                 })) as { role: "user" | "model"; parts: { text: string }[] }[];
                chatRef.current = startChat(history);
            }

            const stream = await chatRef.current.sendMessageStream({ message: input });
            
            let modelResponse = '';
            setMessages(prev => [...prev, { role: 'model', content: '' }]);

            for await (const chunk of stream) {
                modelResponse += chunk.text;
                setMessages(prev => {
                    const newMessages = [...prev];
                    newMessages[newMessages.length - 1] = { role: 'model', content: modelResponse };
                    return newMessages;
                });
            }
        } catch (error) {
            console.error('Chat error:', error);
            setMessages(prev => [...prev, { role: 'model', content: 'Sorry, I encountered an error. Please try again.' }]);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <>
            <button
                onClick={() => setIsOpen(!isOpen)}
                className="fixed bottom-6 right-6 bg-brand-accent hover:bg-brand-accent-hover text-white rounded-full p-4 shadow-lg z-50 transition-transform hover:scale-110"
                aria-label="Toggle chat"
            >
                {isOpen ? <XMarkIcon className="h-8 w-8" /> : <ChatIcon className="h-8 w-8" />}
            </button>

            {isOpen && (
                <div className="fixed bottom-24 right-6 w-96 h-[600px] bg-brand-surface border border-brand-border rounded-lg shadow-2xl z-40 flex flex-col">
                    <div className="p-4 border-b border-brand-border flex justify-between items-center">
                        <h3 className="font-bold text-lg text-brand-text-primary">Vera AI Assistant</h3>
                    </div>
                    <div className="flex-1 p-4 overflow-y-auto space-y-4">
                        {messages.map((msg, index) => (
                            <div key={index} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                                <div className={`max-w-xs p-3 rounded-lg ${msg.role === 'user' ? 'bg-brand-accent text-white' : 'bg-brand-bg text-brand-text-primary'}`}>
                                    {msg.role === 'model' ? <GeminiResponseDisplay content={msg.content} /> : <p>{msg.content}</p>}
                                </div>
                            </div>
                        ))}
                         {isLoading && messages[messages.length - 1].role === 'user' && (
                             <div className="flex justify-start">
                                <div className="max-w-xs p-3 rounded-lg bg-brand-bg text-brand-text-primary">
                                    <div className="animate-pulse flex space-x-2">
                                        <div className="rounded-full bg-gray-500 h-2 w-2"></div>
                                        <div className="rounded-full bg-gray-500 h-2 w-2"></div>
                                        <div className="rounded-full bg-gray-500 h-2 w-2"></div>
                                    </div>
                                </div>
                            </div>
                        )}
                        <div ref={messagesEndRef} />
                    </div>
                    <form onSubmit={handleSendMessage} className="p-4 border-t border-brand-border flex items-center">
                        <input
                            type="text"
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            placeholder="Ask Vera..."
                            className="flex-1 bg-brand-bg border border-brand-border rounded-l-md p-2 text-sm focus:ring-brand-accent focus:border-brand-accent"
                            disabled={isLoading}
                        />
                        <button type="submit" className="bg-brand-accent text-white p-2 rounded-r-md disabled:bg-gray-500" disabled={isLoading}>
                            <PaperAirplaneIcon className="h-5 w-5" />
                        </button>
                    </form>
                </div>
            )}
        </>
    );
};

export default ChatBot;