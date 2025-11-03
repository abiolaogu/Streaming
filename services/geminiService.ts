
import { GoogleGenAI, GenerateContentResponse, Chat } from "@google/genai";
import { GEMINI_API_KEY } from "../config";

const getApiClient = () => {
  if (!GEMINI_API_KEY || GEMINI_API_KEY === "your_api_key_here") {
    throw new Error("API_KEY is not configured. Please add your key to config.ts");
  }
  return new GoogleGenAI({ apiKey: GEMINI_API_KEY });
};

const handleApiError = (error: unknown): string => {
  console.error("Error calling Gemini API:", error);

  // Check for the specific API key error and dispatch an event for the UI to handle.
  const errorObj = error as { error?: { message?: string, status?: string } };
  if (errorObj?.error?.message?.includes("API key not valid") || errorObj?.error?.status === 'PERMISSION_DENIED') {
    window.dispatchEvent(new CustomEvent('apiKeyError'));
    return "Error: Your API Key is no longer valid. Please select a new key to continue.";
  }
  
  if (error instanceof Error) {
    if (error.message.includes('xhr error')) {
      return "Error: A network error occurred after multiple retries. Please check your connection and try again.";
    }
    return `Error: ${error.message}`;
  }
  return "An unknown error occurred.";
};

// Helper function to add retry logic with exponential backoff and jitter
const makeApiCallWithRetry = async <T>(
  apiCall: () => Promise<T>,
  retries = 3,
  delay = 500
): Promise<T> => {
  let lastError: unknown;
  for (let i = 0; i < retries; i++) {
    try {
      return await apiCall();
    } catch (error) {
      lastError = error;
      // Only retry on specific, transient network errors
      if (error instanceof Error && error.message.includes('xhr error')) {
        const jitter = Math.random() * 100; // Add jitter to prevent thundering herd
        const backoffDelay = delay * Math.pow(2, i) + jitter;
        console.warn(`Attempt ${i + 1} failed with network error. Retrying in ${backoffDelay.toFixed(0)}ms...`);
        await new Promise(res => setTimeout(res, backoffDelay));
      } else {
        // Don't retry on other errors (e.g., malformed request, auth issues)
        throw error;
      }
    }
  }
  // If all retries fail, throw the last captured error
  throw lastError;
};


export const getDailyBriefing = async (): Promise<string> => {
  const prompt = `
    Act as the Chief Architect of the StreamVerse platform.
    Provide a concise, markdown-formatted daily briefing for the executive team.
    Covers: System Status, QoE SLOs, Cost Management, Security Posture, Active Incidents.
  `;
  try {
    const response = await makeApiCallWithRetry(async () => {
        const ai = getApiClient();
        return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
    });
    return response.text;
  } catch (error) {
    return handleApiError(error);
  }
};

export const getGpuScalingRecommendation = async (workload: string): Promise<string> => {
  const prompt = `
    Act as an intelligent GPU autoscaler for a media transcoding platform with local and RunPod GPUs.
    Current workload: "${workload}"
    Provide a scaling recommendation (Analysis, Local GPU Action, RunPod Action, Rationale).
  `;
   try {
    const response = await makeApiCallWithRetry(async () => {
        const ai = getApiClient();
        return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
    });
    return response.text;
  } catch (error) {
    return handleApiError(error);
  }
};

export const getCdnStatusReport = async (): Promise<string> => {
    const prompt = `Act as a Lead CDN Engineer for our new Global Cloud CDN v2.0. Provide a brief, markdown-formatted status report. Mention our multi-tier architecture (ATS edges, Varnish shields), our 10 Tbps global capacity, and comment on the current cache hit ratio, global latency, and any active alerts from our Prometheus/Alertmanager stack.`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
}

export const getSecuritySummary = async (): Promise<string> => {
    const prompt = `Act as a DevSecOps Lead. Provide a markdown summary of security posture covering Trivy, OpenSCAP, OPA, Cosign, and Threat Intelligence.`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
}

export const getMediaPipelineStatus = async (): Promise<string> => {
    const prompt = `Act as a Media Operations Engineer. Provide a concise, markdown-formatted status report on the media pipeline (Ingest, Transcode, Packaging, Storage).`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
}

export const getDataPlatformHealth = async (): Promise<string> => {
    const prompt = `Act as a Data Platform Engineer for StreamVerse. Provide a concise, markdown-formatted health summary of our data plane. Cover our primary SQL DB (PostgreSQL), Kafka event bus, ScyllaDB time-series, and DragonflyDB cache. Highlight any replication lag or query performance issues.`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

export const getTelecomStatusReport = async (): Promise<string> => {
    const prompt = `Act as a Telecom Core Engineer. Provide a brief, markdown-formatted status report on the Voice/IMS services (Kamailio/FreeSWITCH) and the Open5GS Mobile Core (AMF/SMF/UPF). Note any registration or session setup failures.`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

export const getSatelliteRolloutSummary = async (): Promise<string> => {
    const prompt = `Act as the lead for the Satellite Overlay project. Provide a concise, markdown-formatted executive summary of the T+2 year rollout plan. Cover the status of the DVB-NIP headend PoC, DVB-I service catalog integration, and partnership with LEO gateway providers.`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-pro', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

export const getDrmStatusSummary = async (): Promise<string> => {
    const prompt = `Act as a Content Security Engineer. Provide a concise, markdown-formatted summary of the DRM platform. Cover the health of Widevine, PlayReady, and FairPlay license servers, any recent key rotation failures, and unusual spikes in license requests.`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

export const getAiOpsSummaryAndActions = async (): Promise<string> => {
    const prompt = `
    Act as an autonomous AIOps agent for the StreamVerse platform. Analyze simulated real-time telemetry (CDN latency, media queue, security CVE, cost spikes) and provide a concise, markdown-formatted summary and a list of actionable recommendations.
    Output: ### Overall Status and ### Actionable Insights.
    `;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({
                model: 'gemini-2.5-pro',
                contents: prompt,
                config: {
                    thinkingConfig: { thinkingBudget: 32768 }
                }
            });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

export const getBiSummary = async (): Promise<string> => {
    const prompt = `Act as a Business Intelligence analyst for StreamVerse. Provide a markdown summary covering Churn Prediction insights (top risk factors, campaign effectiveness) and Content Investment AI (top acquisition opportunities, ROI forecast trends).`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-pro', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

export const getCreatorAnalyticsSummary = async (contentTitle: string): Promise<string> => {
    const prompt = `Act as a Creator Success analyst for StreamVerse. For the title "${contentTitle}", analyze its simulated performance data (watch hours, audience demographics, revenue). Provide a concise, markdown-formatted summary with actionable insights for the creator to improve engagement and earnings.`;
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateContent({ model: 'gemini-2.5-flash', contents: prompt });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

export const generateImage = async (prompt: string, aspectRatio: string): Promise<string> => {
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            return ai.models.generateImages({
                model: 'imagen-4.0-generate-001',
                prompt: prompt,
                config: {
                    numberOfImages: 1,
                    aspectRatio: aspectRatio as "1:1" | "3:4" | "4:3" | "9:16" | "16:9",
                },
            });
        });
        return response.generatedImages[0].image.imageBytes;
    } catch (error) {
        return handleApiError(error);
    }
};

export const analyzeVideo = async (prompt: string, videoBase64: string, mimeType: string): Promise<string> => {
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            const videoPart = { inlineData: { data: videoBase64, mimeType } };
            const textPart = { text: prompt };
            return ai.models.generateContent({
                model: 'gemini-2.5-pro',
                contents: { parts: [videoPart, textPart] },
            });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

export const transcribeAudio = async (audioBase64: string, mimeType: string): Promise<string> => {
    try {
        const response = await makeApiCallWithRetry(async () => {
            const ai = getApiClient();
            const audioPart = { inlineData: { data: audioBase64, mimeType } };
            return ai.models.generateContent({
                model: 'gemini-2.5-flash',
                contents: { parts: [audioPart, { text: "Transcribe this audio." }] },
            });
        });
        return response.text;
    } catch (error) {
        return handleApiError(error);
    }
};

// Chat functionality is streaming, so a simple retry isn't appropriate.
// We'll leave this as-is, as the error seems to be with unary calls.
export const startChat = (history: { role: "user" | "model"; parts: { text: string }[] }[]) => {
    const ai = getApiClient();
    const chat: Chat = ai.chats.create({
        model: 'gemini-2.5-flash-lite', // Use flash-lite for low-latency streaming
        // System instruction to define Vera's personality
        config: {
          systemInstruction: "You are Vera, the Virtual Entertainment & Recommendations Assistant for StreamVerse. You are friendly, knowledgeable, and non-intrusive. Your goal is to help users find content and use the platform effortlessly."
        },
        history: history,
    });
    return chat;
};
