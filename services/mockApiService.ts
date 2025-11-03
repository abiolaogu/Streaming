import { GpuInstance, TranscodeJob, LiveStreamIngest, CreatorContent, RevenueAnalytics, ChurnRiskProfile, ContentRoiPrediction, DvbComponent, UserProfile, UserPreferences, UserDevice, CdnStat, CdnAlert, AnyContent, MediaContent, MediaContentWithProgress, CacheHitRatioDataPoint, LatencyByRegion, BandwidthByTier } from '../types';

// This file simulates a real API client.
// In a real application, these functions would make network requests (e.g., using fetch).
// The timeout simulates network latency.

const MOCK_API_LATENCY = 300; // 300ms

// --- MOCK DATA ---
const now = Date.now();
const oneHour = 60 * 60 * 1000;

const mockLiveChannels: AnyContent[] = [
    { id: 'live-news-1', name: 'Nollywood News Live', description: "Your 24/7 source for news and gossip from the heart of Nollywood.", logoUrl: 'https://picsum.photos/seed/nnl/100/50', category: 'Nollywood', type: 'Live', currentProgram: { title: 'Industry Insider', startTime: now - oneHour * 0.25, endTime: now + oneHour * 0.75 } },
    { id: 'live-movies-1', name: 'StreamVerse Premiere', description: "Exclusive premieres and classic movies, playing non-stop.", logoUrl: 'https://picsum.photos/seed/svlogo/100/50', category: 'Movies', type: 'Live', currentProgram: { title: 'The Wedding Party', startTime: now - oneHour * 0.8, endTime: now + oneHour * 0.2 } },
];

// FIX: Changed the type to allow content with progress property.
const mockVODContent: (MediaContent | MediaContentWithProgress)[] = [
     { id: 'movie-1', title: 'Omo Ghetto: The Saga', description: "A story of two identical twins living separate lives, one a ghetto-raised hustler and the other a wealthy sophisticate.", thumbnailUrl: 'https://picsum.photos/seed/omoghetto/400/225', type: 'Movie', category: 'Nollywood', monetizationModel: 'SVOD', cast: ['Funke Akindele-Bello', 'Nancy Isime'], director: 'Funke Akindele', releaseYear: 2020 },
     { id: 'movie-2', title: 'King of Boys', description: "A businesswoman and philanthropist with a checkered past is drawn into a power struggle that threatens everything she holds dear.", thumbnailUrl: 'https://picsum.photos/seed/kingofboys/400/225', type: 'Movie', category: 'Nollywood', monetizationModel: 'SVOD', cast: ['Sola Sobowale', 'Adesua Etomi-Wellington'], director: 'Kemi Adetiba', releaseYear: 2018 },
     { id: 'series-1', title: 'Shanty Town', description: "A group of courtesans attempts to escape the grasp of a notorious kingpin, but political corruption and blood ties make freedom a near-impossible goal.", thumbnailUrl: 'https://picsum.photos/seed/shantytown/400/225', type: 'Series', category: 'Nollywood', monetizationModel: 'SVOD', cast: ['Chidi Mokeme', 'Ini Edo', 'Richard Mofe-Damijo'], director: 'Dimeji Ajibola', releaseYear: 2023 },
     { id: 'movie-3', title: 'A Tribe Called Judah', description: "A single mother's five sons from five different fathers unite to rob a small mall to save her life, but their plan takes an unexpected turn.", thumbnailUrl: 'https://picsum.photos/seed/judah/400/225', type: 'Movie', category: 'Nollywood', monetizationModel: 'TVOD', price: { rent: 3.99, buy: 9.99}, cast: ['Funke Akindele', 'Jide Kene Achufusi'], director: 'Funke Akindele', releaseYear: 2023 },
     { id: 'doc-1', title: 'Journey of an African Colony', description: "An in-depth look at the history of Nigeria, from its colonial past to its vibrant present.", thumbnailUrl: 'https://picsum.photos/seed/colony/400/225', type: 'Movie', category: 'Documentaries', monetizationModel: 'AVOD', cast: ['N/A'], director: 'Olasupo Shasore', releaseYear: 2019 },
     { id: 'drama-1', title: 'Blood Sisters', description: "Bound by a dangerous secret, best friends Sarah and Kemi are forced to go on the run after a wealthy groom disappears during his engagement party.", thumbnailUrl: 'https://picsum.photos/seed/bloodsisters/400/225', type: 'Series', category: 'Drama', monetizationModel: 'SVOD', cast: ['Ini Dima-Okojie', 'Nancy Isime'], director: 'Biyi Bandele', releaseYear: 2022 },
     { id: 'series-4', title: 'Far From Home', description: "A financially struggling teen finds himself in the world of luxury after a prestigious scholarship sends him to an exclusive school for the one percent.", thumbnailUrl: 'https://picsum.photos/seed/farfromhome/400/225', type: 'Series', progress: 75, category: 'Drama', monetizationModel: 'SVOD', cast: ['Mike Afolarin', 'Richard Mofe-Damijo'], director: 'Catherine Stewart', releaseYear: 2022 }
];

const allContent: AnyContent[] = [...mockLiveChannels, ...mockVODContent];

const mockWatchHistory: AnyContent[] = [
    mockVODContent[2], // Shanty Town
    mockVODContent[0], // Omo Ghetto
    mockVODContent[4], // Journey of an African Colony
];

const mockWatchlist: AnyContent[] = [
    mockVODContent[1], // King of Boys
    mockVODContent[3], // A Tribe Called Judah
    mockVODContent[5], // Blood Sisters
];


const mockGpuInstances: GpuInstance[] = [
    { id: 'local-gpu-01', provider: 'Local', type: 'NVIDIA A6000', status: 'Processing', spot: false, costPerHour: 1.50 },
    { id: 'local-gpu-02', provider: 'Local', type: 'NVIDIA A6000', status: 'Idle', spot: false, costPerHour: 1.50 },
    { id: 'runpod-spot-1', provider: 'RunPod', type: 'RTX 4090', status: 'Processing', spot: true, costPerHour: 0.79 },
    { id: 'runpod-spot-2', provider: 'RunPod', type: 'RTX 4090', status: 'Processing', spot: true, costPerHour: 0.79 },
    { id: 'runpod-spot-3', provider: 'RunPod', type: 'RTX 4090', status: 'Terminating', spot: true, costPerHour: 0.79 },
];

const mockTranscodeJobs: TranscodeJob[] = [
    { id: 'job-a1b2', source: 'Live Event 1 (SRT)', profile: '4K HEVC 10-bit', status: 'In Progress', gpuId: 'local-gpu-01', progress: 82 },
    { id: 'job-c3d4', source: 'VOD Ingest (S3)', profile: '1080p H.264 Main', status: 'In Progress', gpuId: 'runpod-spot-1', progress: 65 },
    { id: 'job-e5f6', source: 'Live Event 2 (RIST)', profile: '720p H.264 High', status: 'Queued', gpuId: 'N/A', progress: 0 },
    { id: 'job-g7h8', source: 'VOD Ingest (S3)', profile: '480p AV1 Low', status: 'Completed', gpuId: 'runpod-spot-2', progress: 100 },
    { id: 'job-i9j0', source: 'Live Event 3 (RTMP)', profile: '1080p HEVC Main', status: 'Failed', gpuId: 'local-gpu-02', progress: 43 },
];

const mockLiveIngests: LiveStreamIngest[] = [
    { id: 'EU-CENTRAL-01A', protocol: 'RIST', status: 'Critical', bitrateMb: 1.2, connections: 1, issues: 'Clock Sync Error' },
    { id: 'EU-WEST-04C', protocol: 'SRT', status: 'Packet Loss', bitrateMb: 6.2, connections: 1, issues: 'High Jitter Detected' },
    { id: 'US-EAST-01A', protocol: 'SRT', status: 'Stable', bitrateMb: 8.5, connections: 1, issues: 'None' },
];

const mockCreatorContent: CreatorContent[] = [
    { id: 'creator-mov-1', title: 'Chronicles of Lagos', status: 'Live', uploadDate: '2024-05-10', monetization: ['SVOD', 'TVOD'] },
    { id: 'creator-mov-2', title: 'The Abuja Connection', status: 'In Review', uploadDate: '2024-05-20', monetization: ['SVOD'] },
    { id: 'creator-doc-1', title: 'A Day in the Market', status: 'Live', uploadDate: '2024-04-25', monetization: ['AVOD'] },
    { id: 'creator-mov-3', title: 'Rejected Film', status: 'Rejected', uploadDate: '2024-05-18', monetization: [] },
];

const mockRevenueAnalytics: RevenueAnalytics[] = [
    { period: 'Daily', watchHours: 12500, svodRevenue: 1800, avodRevenue: 450, tvodRevenue: 950, totalRevenue: 3200 },
    { period: 'Weekly', watchHours: 89000, svodRevenue: 12600, avodRevenue: 3150, tvodRevenue: 6650, totalRevenue: 22400 },
    { period: 'Monthly', watchHours: 380000, svodRevenue: 54000, avodRevenue: 13500, tvodRevenue: 28500, totalRevenue: 96000 },
    { period: 'Yearly', watchHours: 1500000, svodRevenue: 216000, avodRevenue: 54000, tvodRevenue: 114000, totalRevenue: 384000 },
];

const mockChurnRiskProfiles: ChurnRiskProfile[] = [
    { subscriberId: 'sub-a1b2c3d4', riskScore: 94, primaryReason: 'Low Engagement', recommendedOffer: '30% off 3 months' },
    { subscriberId: 'sub-e5f6g7h8', riskScore: 81, primaryReason: 'Payment Declined', recommendedOffer: 'Update Payment Method + 1 week free' },
    { subscriberId: 'sub-i9j0k1l2', riskScore: 75, primaryReason: 'Content Gap (Sports)', recommendedOffer: 'Free Sports Plus Add-On' },
];

const mockContentRoi: ContentRoiPrediction[] = [
    { contentTitle: 'Acquisition Target A', predictedViewership: 2500000, predictedRoi: 2.1, confidence: 'High' },
    { contentTitle: 'Indie Film Festival Winner', predictedViewership: 450000, predictedRoi: 1.3, confidence: 'Medium' },
    { contentTitle: 'Foreign Language Series', predictedViewership: 800000, predictedRoi: 0.9, confidence: 'High' },
];

const mockDvbComponents: DvbComponent[] = [
    { id: 'sat-uplink-eu', type: 'DVB-NIP', status: 'Nominal', details: 'SES-1, 450 Mbps throughput' },
    { id: 'service-list-main', type: 'DVB-I', status: 'Nominal', details: 'Last updated: 2 mins ago' },
    { id: 'iptv-verizon-us', type: 'DVB-IP', status: 'Warning', details: 'High multicast packet loss on CH 102' },
];

const mockUserProfile: UserProfile = {
    name: 'Platform Admin',
    email: 'admin@streamverse.io',
    memberSince: '2024-01-01',
    subscriptionPlan: 'StreamVerse Premium (Internal)',
};

const mockUserPreferences: UserPreferences = {
    language: 'English',
    playbackQuality: '4K',
    subtitles: true,
};

const mockUserDevices: UserDevice[] = [
    { id: 'dev-1', type: 'Web', os: 'macOS Sonoma', lastSeen: '2 minutes ago' },
    { id: 'dev-2', type: 'Mobile', os: 'iOS 17.5', lastSeen: '3 hours ago' },
    { id: 'dev-3', type: 'Smart TV', os: 'LG webOS 23', lastSeen: 'Yesterday' },
];

const mockCdnStatsBase = {
    bandwidth: 4.82,
    requests: 1.8,
    hitRatio: 96.3,
    latency: 48,
};

const allPossibleAlerts: Omit<CdnAlert, 'id' | 'timestamp'>[] = [
    { severity: 'Critical', message: 'BGP session lost with peer in `lon1`' },
    { severity: 'Warning', message: 'Cache hit ratio dropped below 85% in `sao1`' },
    { severity: 'Info', message: 'Tier 3 PoP `jnb3` scaled up successfully' },
    { severity: 'Critical', message: 'Power outage detected at `fra2` data center.'},
    { severity: 'Warning', message: 'High packet loss on backbone link between `ash1` and `lon1`.' },
    { severity: 'Info', message: 'Software update `v3.2.1` rolled out to all Tier 2 PoPs.' },
    { severity: 'Warning', message: 'Origin shield latency > 200ms for `mum2`.' },
];
let alertIndex = 0;


const mockCacheHitRatioData: CacheHitRatioDataPoint[] = [
    { time: '12:00', ratio: 95.8 },
    { time: '13:00', ratio: 96.1 },
    { time: '14:00', ratio: 96.5 },
    { time: '15:00', ratio: 96.2 },
    { time: '16:00', ratio: 96.3 },
    { time: '17:00', ratio: 96.8 },
    { time: '18:00', ratio: 97.1 },
];

const mockLatencyByRegionData: LatencyByRegion[] = [
    { region: 'Americas', latency: 55 },
    { region: 'EMEA', latency: 42 },
    { region: 'APAC', latency: 89 },
];

const mockBandwidthByTierData: BandwidthByTier[] = [
    { tier: 'Tier 1', bandwidth: 3.1 },
    { tier: 'Tier 2', bandwidth: 1.72 },
];


// --- EXPORTED API FUNCTIONS ---
export const fetchVODContent = (): Promise<AnyContent[]> => new Promise(resolve => setTimeout(() => resolve(mockVODContent as MediaContent[]), MOCK_API_LATENCY));
export const fetchLiveChannels = (): Promise<AnyContent[]> => new Promise(resolve => setTimeout(() => resolve(mockLiveChannels), MOCK_API_LATENCY));
export const fetchContinueWatching = (): Promise<AnyContent[]> => new Promise(resolve => setTimeout(() => resolve(allContent.filter(c => 'progress' in c)), MOCK_API_LATENCY));

export const fetchContentDetails = (id: string): Promise<AnyContent | undefined> => new Promise(resolve => setTimeout(() => resolve(allContent.find(c => c.id === id)), MOCK_API_LATENCY));

export const fetchGpuInstances = (): Promise<GpuInstance[]> => new Promise(resolve => setTimeout(() => resolve(mockGpuInstances), MOCK_API_LATENCY));
export const fetchTranscodeJobs = (): Promise<TranscodeJob[]> => new Promise(resolve => setTimeout(() => resolve(mockTranscodeJobs), MOCK_API_LATENCY));
export const fetchLiveIngests = (): Promise<LiveStreamIngest[]> => new Promise(resolve => setTimeout(() => resolve(mockLiveIngests), MOCK_API_LATENCY));
export const fetchCreatorContent = (): Promise<CreatorContent[]> => new Promise(resolve => setTimeout(() => resolve(mockCreatorContent), MOCK_API_LATENCY));
export const fetchRevenueAnalytics = (): Promise<RevenueAnalytics[]> => new Promise(resolve => setTimeout(() => resolve(mockRevenueAnalytics), MOCK_API_LATENCY));
export const fetchChurnRiskProfiles = (): Promise<ChurnRiskProfile[]> => new Promise(resolve => setTimeout(() => resolve(mockChurnRiskProfiles), MOCK_API_LATENCY));
export const fetchContentRoi = (): Promise<ContentRoiPrediction[]> => new Promise(resolve => setTimeout(() => resolve(mockContentRoi), MOCK_API_LATENCY));
export const fetchDvbComponents = (): Promise<DvbComponent[]> => new Promise(resolve => setTimeout(() => resolve(mockDvbComponents), MOCK_API_LATENCY));
export const fetchUserProfile = (): Promise<UserProfile> => new Promise(resolve => setTimeout(() => resolve(mockUserProfile), MOCK_API_LATENCY));
export const fetchUserPreferences = (): Promise<UserPreferences> => new Promise(resolve => setTimeout(() => resolve(mockUserPreferences), MOCK_API_LATENCY));
export const fetchUserDevices = (): Promise<UserDevice[]> => new Promise(resolve => setTimeout(() => resolve(mockUserDevices), MOCK_API_LATENCY));

export const fetchCdnStats = (): Promise<CdnStat[]> => {
    return new Promise(resolve => {
        // Simulate real-time fluctuations
        const newBandwidth = mockCdnStatsBase.bandwidth + (Math.random() - 0.5) * 0.2;
        const newRequests = mockCdnStatsBase.requests + (Math.random() - 0.5) * 0.1;
        const newHitRatio = mockCdnStatsBase.hitRatio + (Math.random() - 0.5) * 0.5;
        const newLatency = mockCdnStatsBase.latency + (Math.random() - 0.5) * 5;

        const formatChange = (oldVal: number, newVal: number) => {
            const change = ((newVal - oldVal) / oldVal) * 100;
            return `${change > 0 ? '+' : ''}${change.toFixed(1)}%`;
        };

        setTimeout(() => resolve([
            { label: 'Global Bandwidth', value: `${newBandwidth.toFixed(2)} Tbps`, change: formatChange(4.82, newBandwidth), changeType: newBandwidth > 4.82 ? 'positive' : 'negative' },
            { label: 'Requests/sec', value: `${newRequests.toFixed(2)}M`, change: formatChange(1.8, newRequests), changeType: newRequests > 1.8 ? 'positive' : 'negative' },
            { label: 'Cache Hit Ratio', value: `${newHitRatio.toFixed(1)}%`, change: formatChange(96.3, newHitRatio), changeType: newHitRatio > 96.3 ? 'positive' : 'negative' },
            { label: 'P99 Latency', value: `${Math.round(newLatency)}ms`, change: formatChange(48, newLatency), changeType: newLatency < 48 ? 'positive' : 'negative' },
        ]), MOCK_API_LATENCY);
    });
};

export const fetchCdnAlerts = (): Promise<CdnAlert[]> => {
    return new Promise(resolve => {
        // Cycle through a larger list of alerts to simulate new ones appearing
        const newAlerts = [
            allPossibleAlerts[alertIndex % allPossibleAlerts.length],
            allPossibleAlerts[(alertIndex + 1) % allPossibleAlerts.length],
            allPossibleAlerts[(alertIndex + 2) % allPossibleAlerts.length],
        ].map((alert, index) => ({
            ...alert,
            id: `alert-${alertIndex}-${index}`,
            timestamp: new Date(Date.now() - (index * 5 + 2) * 60 * 1000).toISOString(),
        }));
        alertIndex++;
        setTimeout(() => resolve(newAlerts.reverse()), MOCK_API_LATENCY);
    });
};


export const fetchWatchHistory = (): Promise<AnyContent[]> => new Promise(resolve => setTimeout(() => resolve(mockWatchHistory), MOCK_API_LATENCY));
export const fetchWatchlist = (): Promise<AnyContent[]> => new Promise(resolve => setTimeout(() => resolve(mockWatchlist), MOCK_API_LATENCY));

export const fetchCacheHitRatioData = (): Promise<CacheHitRatioDataPoint[]> => new Promise(resolve => setTimeout(() => resolve(mockCacheHitRatioData), MOCK_API_LATENCY + 200));
export const fetchLatencyByRegionData = (): Promise<LatencyByRegion[]> => new Promise(resolve => setTimeout(() => resolve(mockLatencyByRegionData), MOCK_API_LATENCY + 200));
export const fetchBandwidthByTierData = (): Promise<BandwidthByTier[]> => new Promise(resolve => setTimeout(() => resolve(mockBandwidthByTierData), MOCK_API_LATENCY + 200));