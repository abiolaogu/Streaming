
import React from 'react';

// A very simple markdown to HTML converter for this use case.
// A more robust solution would use a library like 'marked' or 'react-markdown'.
const SimpleMarkdown: React.FC<{ text: string }> = ({ text }) => {
  const lines = text.split('\n');
  const elements = lines.map((line, index) => {
    if (line.startsWith('### ')) {
      return <h3 key={index} className="text-lg font-semibold mt-4 mb-2">{line.substring(4)}</h3>;
    }
    if (line.startsWith('## ')) {
      return <h2 key={index} className="text-xl font-bold mt-6 mb-3">{line.substring(3)}</h2>;
    }
    if (line.startsWith('# ')) {
      return <h1 key={index} className="text-2xl font-bold mt-8 mb-4">{line.substring(2)}</h1>;
    }
    if (line.startsWith('* ')) {
        const boldRegex = /\*\*(.*?)\*\*/g;
        const parts = line.substring(2).split(boldRegex);
        return (
            <li key={index} className="ml-5 list-disc">
                 {parts.map((part, i) => i % 2 === 1 ? <strong key={i}>{part}</strong> : part)}
            </li>
        );
    }
    if (line.trim() === '') {
      return <br key={index} />;
    }

    const boldRegex = /\*\*(.*?)\*\*/g;
    const parts = line.split(boldRegex);

    return (
      <p key={index} className="text-brand-text-secondary leading-relaxed">
        {parts.map((part, i) => i % 2 === 1 ? <strong key={i} className="text-brand-text-primary">{part}</strong> : part)}
      </p>
    );
  });
  return <div className="prose prose-invert max-w-none">{elements}</div>;
};


interface GeminiResponseDisplayProps {
    content: string;
}

const GeminiResponseDisplay: React.FC<GeminiResponseDisplayProps> = ({ content }) => {
    return (
        <div className="space-y-4">
           <SimpleMarkdown text={content} />
        </div>
    );
};

export default GeminiResponseDisplay;
