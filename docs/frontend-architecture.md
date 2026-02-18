# Frontend Architecture — StreamVerse Streaming Platform

## 1. Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| Framework | React | 18.2+ |
| Language | TypeScript | 5.3+ |
| Build Tool | Vite | 7.3+ |
| Styling | Tailwind CSS | via brand-* utility classes |
| Linting | ESLint | 9.39+ |
| Type Checking | TypeScript Compiler | `tsc --noEmit` |
| Package Manager | npm | via `package-lock.json` |

### Build Commands (from `package.json`)
```json
{
  "dev": "vite",
  "build": "tsc && vite build",
  "preview": "vite preview",
  "lint": "eslint \"src/**/*.{ts,tsx}\"",
  "lint:fix": "npm run lint -- --fix",
  "type-check": "tsc --noEmit",
  "check": "npm run type-check && npm run lint && npm run build"
}
```

---

## 2. Application Entry Points

### `index.html`
Static HTML shell with a root div. Vite injects the compiled JavaScript bundle.

### `index.tsx`
React root render mounting `<App />` into the DOM.

### `App.tsx` — Main Application Component
The root component manages:
- **API key gate**: Checks `VITE_GEMINI_API_KEY` env var or `window.aistudio` for Gemini API access
- **Authentication gate**: Shows `LoginView` until `isAuthenticated` is true
- **Theme management**: Dark/light mode persisted in localStorage, toggled via `<html>` class
- **View routing**: State-based routing via `currentViewState: ViewState`
- **Layout**: Sidebar + Header + Main content + ChatBot overlay

---

## 3. Type System

All TypeScript types are centralized in `types.ts`:

### View Types
```typescript
type AdminView = 'overview' | 'gpu_fabric' | 'cdn' | 'security' | 'media'
  | 'data' | 'telecom' | 'satellite' | 'drm' | 'ai_ops'
  | 'neural_engine' | 'broadcast_ops' | 'bi' | 'user_profile' | 'roadmap';
type CreatorView = 'creator_studio';
type ConsumerView = 'streamverse_home' | 'watch';
type View = AdminView | CreatorView | ConsumerView;
```

### Content Types
- `MediaContent`: Movies and series with metadata, thumbnails, monetization model
- `LiveChannel`: Live broadcast channels with current program schedule
- `MediaContentWithProgress`: Extends MediaContent with watch progress percentage
- `ContentCategory`: Named content row (e.g., "Trending Now") with item array

### Infrastructure Types
- `GpuInstance`: GPU node (Local/RunPod, status, cost/hour)
- `TranscodeJob`: Transcoding job (source, profile, progress, GPU assignment)
- `LiveStreamIngest`: Ingest stream (SRT/RIST/RTMP, bitrate, stability)
- `DrmLicenseServer`: DRM server (region, provider, p95 latency)
- `DrmKey`: Content encryption key (rotation status, licenses issued)
- `CdnStat`, `CdnAlert`, `CacheHitRatioDataPoint`, `LatencyByRegion`: CDN analytics

### Business Types
- `ChurnRiskProfile`: Subscriber churn prediction (risk score, recommended offer)
- `ContentRoiPrediction`: Content investment ROI prediction
- `RevenueAnalytics`: Revenue breakdown by period and monetization model
- `CreatorContent`: Creator's uploaded content (status, monetization)

---

## 4. Component Architecture

### Layout Components (`components/`)

| Component | File | Responsibility |
|-----------|------|----------------|
| Sidebar | `Sidebar.tsx` | Navigation menu with view switching, admin/consumer/creator sections |
| Header | `Header.tsx` | Current view title, user actions, theme toggle, logout |
| ChatBot | `ChatBot.tsx` | AI assistant overlay using Gemini API |
| DashboardCard | `DashboardCard.tsx` | Reusable card container for dashboard metrics |
| LoadingSpinner | `LoadingSpinner.tsx` | Loading state indicator |
| WorldMap | `WorldMap.tsx` | Geographic visualization for CDN and user distribution |
| ThemeToggle | `ThemeToggle.tsx` | Dark/light mode switch |
| DayNightOverlay | `DayNightOverlay.tsx` | Visual day/night overlay for world map |
| GeminiResponseDisplay | `GeminiResponseDisplay.tsx` | Formatted AI response rendering |
| IconComponents | `IconComponents.tsx` | SVG icon library |

### View Components (`views/`)

#### Consumer Views
| View | File | Features |
|------|------|----------|
| StreamVerseHomeView | `StreamVerseHomeView.tsx` | Hero banner, content rows, FAST channels, trending |
| WatchView | `WatchView.tsx` | Video player, metadata, related content, comments |

#### Creator Views
| View | File | Features |
|------|------|----------|
| CreatorStudioView | `CreatorStudioView.tsx` | Upload management, analytics, monetization settings |

#### Admin Dashboard Views
| View | File | Features |
|------|------|----------|
| OverviewView | `OverviewView.tsx` | Platform health, SLO metrics, service status |
| GpuFabricView | `GpuFabricView.tsx` | GPU utilization, transcode jobs, Runpod.io scaling |
| CdnView | `CdnView.tsx` | Cache hit ratios, bandwidth, latency by region |
| SecurityView | `SecurityView.tsx` | Threat monitoring, vulnerability scanning |
| MediaView | `MediaView.tsx` | Transcoding pipeline, ingest streams, media processing |
| DataView | `DataView.tsx` | Database metrics, Kafka lag, data pipeline health |
| TelecomView | `TelecomView.tsx` | Telecom integration status |
| SatelliteView | `SatelliteView.tsx` | Satellite delivery telemetry |
| DrmView | `DrmView.tsx` | License server status, key rotation, provider health |
| AiOpsView | `AiOpsView.tsx` | AI-driven operational actions (CDN, Media, Security, Cost) |
| NeuralContentEngineView | `NeuralContentEngineView.tsx` | Content analysis ML pipeline |
| BusinessIntelligenceView | `BusinessIntelligenceView.tsx` | Churn prediction, content ROI, revenue forecasting |
| BroadcastOpsView | `BroadcastOpsView.tsx` | DVB-NIP, DVB-I, DVB-IP component status |
| RoadmapView | `RoadmapView.tsx` | GitHub issues integration, phase tracking |
| UserProfileView | `UserProfileView.tsx` | User settings, devices, preferences |

---

## 5. Service Layer

### `services/geminiService.ts`
API client for Google Gemini AI integration. Used by the ChatBot component for:
- Platform operation questions
- Content recommendation queries
- Troubleshooting assistance
- Analytics interpretation

### `services/mockApiService.ts`
Mock data service providing development-time data for all dashboard views. Returns typed mock data matching the production API shapes.

### API Communication Pattern
```typescript
// Standard pattern for API calls
const response = await fetch(`${API_BASE_URL}/v1/endpoint`, {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json'
  }
});
```

---

## 6. Styling Architecture

### Theme System
- Dark mode (default): `document.documentElement.classList.add('dark')`
- Light mode: class removed
- Persisted in `localStorage.getItem('theme')`
- Brand utility classes: `bg-brand-bg`, `text-brand-text`, etc.

### Design Tokens
- Background: Deep dark (#0a0a0f range for dark mode)
- Primary accent: Brand purple/blue gradients
- Cards: Semi-transparent with backdrop blur
- Typography: System font stack via `font-sans`
- Spacing: Tailwind responsive scale (p-4 sm:p-6 lg:p-8)

---

## 7. Build and Deployment

### Development
```bash
npm install
npm run dev          # Vite dev server at http://localhost:5173
```

### Production Build
```bash
npm run check        # Type check + lint + build
npm run build        # Outputs to dist/
```

### Docker Build (from `Dockerfile`)
```dockerfile
# Stage 1: Build with Node 18
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

### Nginx Configuration
- Gzip compression for text assets
- Security headers (X-Frame-Options, X-Content-Type-Options, X-XSS-Protection)
- Static asset caching: 1 year with immutable flag
- SPA fallback: `try_files $uri $uri/ /index.html`
- API proxy: `/api/` forwards to `api-gateway:8080`
- Health check: `/health` returns 200

---

## 8. Multi-Platform Client Architecture

### Mobile (Flutter)
Located in `apps/clients/mobile-flutter/` and `mobile-app/`:
- Dart language with Flutter framework
- Platform: iOS 13+, Android 7.0+
- Features: Offline downloads, PiP, Chromecast, biometric auth
- Service layer in `lib/services/`
- Screen navigation in `lib/screens/`

### TV Apps
Each TV platform has a dedicated implementation:

| Platform | Directory | Technology |
|----------|-----------|-----------|
| Android TV | `apps/clients/tv-apps/android-tv/` | Kotlin, Leanback, ExoPlayer |
| Samsung Tizen | `apps/clients/tv-apps/samsung-tizen/` | HTML5/JS/CSS |
| LG webOS | `apps/clients/tv-apps/lg-webos/` | Enact (React-based) |
| Roku | `apps/clients/tv-apps/roku/` | BrightScript, SceneGraph |
| Apple tvOS | `apps/clients/tv-apps/apple-tvos/` | Swift, UIKit |
| Amazon Fire TV | `apps/clients/tv-apps/amazon-fire-tv/` | Kotlin, Fire UI |

### Player SDKs
Located in `packages/sdk/`:
- `player-web/`: Web video player SDK with DRM support
- `player-flutter/`: Flutter video player SDK with DRM support

---

## 9. State Management Patterns

### Component-Level State
React `useState` hooks for local UI state:
- `currentViewState`: Active view and parameters
- `isAuthenticated`: Auth gate boolean
- `theme`: Dark/light mode
- `isKeySelected`: API key initialization

### Callback-Based Navigation
View changes propagated via callback props:
```typescript
const handleSetView = (view: View) => {
  setCurrentViewState({ name: view });
};

const handleWatchContent = (contentId: string) => {
  setCurrentViewState({ name: 'watch', params: { contentId } });
};
```

### Data Fetching
Each view component manages its own data fetching lifecycle using `useEffect` and `useCallback` hooks. Mock data service provides development fallbacks.

---

## 10. Configuration

### `config.ts`
Exports API base URL configuration, typically reading from `import.meta.env.VITE_API_URL`.

### `vite.config.ts`
Vite configuration with React plugin, build optimization, and development server settings.

### `tsconfig.json`
TypeScript compiler configuration with strict mode enabled, JSX React transform, and path aliases.

### `eslint.config.js`
ESLint flat config with TypeScript parser, React Hooks plugin, and React Refresh plugin.

---

**Document Version**: 2.0
**Last Updated**: 2026-02-17
