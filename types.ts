// FIX: Import React to provide types for component definitions.
import type React from 'react';

export type AdminView = 
    'overview' | 'gpu_fabric' | 'cdn' | 'security' | 'media' 
    | 'data' | 'telecom' | 'satellite' | 'drm' | 'ai_ops' 
    | 'neural_engine' | 'broadcast_ops' | 'bi' | 'user_profile';

export type CreatorView = 'creator_studio';

export type ConsumerView = 'streamverse_home' | 'watch';

export type View = AdminView | CreatorView | ConsumerView;

export type Theme = 'light' | 'dark';

export interface ViewState {
    name: View;
    params?: {
        contentId?: string;
    };
}

export interface NavItem {
  id: View;
  label: string;
  icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
}

export interface GpuInstance {
  id: string;
  provider: 'Local' | 'RunPod';
  type: string;
  status: 'Idle' | 'Processing' | 'Terminating';
  spot: boolean;
  costPerHour: number;
}

export interface SloMetric {
  name: string;
  value: string;
  target: string;
  status: 'ok' | 'warning' | 'critical';
}

export interface TranscodeJob {
    id: string;
    source: string;
    profile: string;
    status: 'Queued' | 'In Progress' | 'Completed' | 'Failed';
    gpuId: string;
    progress: number; // e.g., 75 for 75%
}

export interface LiveStreamIngest {
  id: string;
  protocol: 'SRT' | 'RIST' | 'RTMP';
  status: 'Stable' | 'Packet Loss' | 'Offline' | 'Critical';
  bitrateMb: number;
  connections: number;
  issues: string;
}

export interface DataPlatformMetric {
    name: string;
    status: 'ok' | 'warning' | 'critical';
    value: string;
}

export interface TelecomServiceStatus {
    name: string;
    status: 'Online' | 'Degraded' | 'Offline';
    details: string;
}

export interface SatelliteStatus {
    component: string;
    status: 'Nominal' | 'Warning' | 'Error';
    telemetry: string;
}

export interface DrmLicenseServer {
    id: string;
    region: string;
    provider: 'Widevine' | 'PlayReady' | 'FairPlay';
    status: 'Active' | 'Degraded' | 'Offline';
    p95LatencyMs: number;
}

export interface DrmKey {
    contentId: string;
    status: 'Active' | 'Rotated' | 'Compromised';
    lastRotation: string;
    licensesIssued: number;
}

// --- StreamVerse Consumer Content Types ---
export type ContentCategoryString = 'News' | 'Sports' | 'Movies' | 'Lifestyle' | 'Drama' | 'Documentaries' | 'Educational' | 'Nollywood';

export interface MediaContent {
  id: string;
  title: string;
  description: string;
  thumbnailUrl: string;
  type: 'Movie' | 'Series';
  category: ContentCategoryString;
  monetizationModel: 'AVOD' | 'SVOD' | 'TVOD';
  price?: {
      rent: number;
      buy: number;
  };
}

export interface LiveChannel {
    id: string;
    name: string;
    description: string;
    logoUrl: string;
    category: ContentCategoryString;
    type: 'Live';
    currentProgram: {
        title: string;
        startTime: number; // Unix timestamp
        endTime: number; // Unix timestamp
    };
}


export interface MediaContentWithProgress extends MediaContent {
    progress: number; // e.g. 80 for 80% watched
}

export type AnyContent = MediaContent | MediaContentWithProgress | LiveChannel;

export interface ContentCategory {
  title: string;
  items: AnyContent[];
}

export interface ChatMessage {
  role: 'user' | 'model';
  content: string;
}

// --- Platform & AI Types ---
export interface AiAction {
    id: string;
    category: 'CDN' | 'Media' | 'Security' | 'Cost';
    description: string;
    impact: 'Low' | 'Medium' | 'High';
    status: 'Pending' | 'Approved' | 'Denied' | 'Complete' | 'Auto-Executed' | 'Canceled by Operator';
    timestamp?: string;
}


// --- Creator Economy Types ---
export interface CreatorContent {
    id: string;
    title: string;
    status: 'Processing' | 'In Review' | 'Live' | 'Rejected';
    uploadDate: string;
    monetization: ('AVOD' | 'SVOD' | 'TVOD')[];
}

export interface RevenueAnalytics {
    period: 'Daily' | 'Weekly' | 'Monthly' | 'Yearly';
    watchHours: number;
    svodRevenue: number;
    avodRevenue: number;
    tvodRevenue: number;
    totalRevenue: number;
}

// --- Business Intelligence Types ---
export interface ChurnRiskProfile {
    subscriberId: string;
    riskScore: number; // 0-100
    primaryReason: string;
    recommendedOffer: string;
}

export interface ContentRoiPrediction {
    contentTitle: string;
    predictedViewership: number;
    predictedRoi: number; // e.g., 1.5 for 150%
    confidence: 'High' | 'Medium' | 'Low';
}

// --- Broadcast Ops Types ---
export interface DvbComponent {
    id: string;
    type: 'DVB-NIP' | 'DVB-I' | 'DVB-IP';
    status: 'Nominal' | 'Warning' | 'Critical';
    details: string;
}

// --- User Profile Types ---
export interface UserProfile {
    name: string;
    email: string;
    memberSince: string;
    subscriptionPlan: string;
}

export interface UserPreferences {
    language: string;
    playbackQuality: 'Auto' | 'HD' | '4K';
    subtitles: boolean;
}

export interface UserDevice {
    id: string;
    type: 'Web' | 'Mobile' | 'Smart TV';
    os: string;
    lastSeen: string;
}

// --- CDN v2.0 Types ---
export interface CdnStat {
    label: string;
    value: string;
    change: string;
    changeType: 'positive' | 'negative' | 'neutral';
}

export interface CdnAlert {
    id: string;
    severity: 'Critical' | 'Warning' | 'Info';
    message: string;
    timestamp: string;
}