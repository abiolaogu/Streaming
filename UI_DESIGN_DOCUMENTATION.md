# StreamVerse UI/UX Design Documentation

## 🎨 Design Philosophy

StreamVerse follows Netflix's design principles with a focus on:
- **Content-First**: Minimal UI chrome to maximize content visibility
- **Dark Theme**: Black (#000000) and dark gray (#141414) backgrounds
- **Bold Accents**: Netflix red (#E50914) for CTAs and highlights
- **Smooth Interactions**: Hover effects, transitions, and animations
- **Responsive**: Adapts seamlessly across all screen sizes

---

## 💻 WEB APPLICATION UI

### 1. Login Page
```
┌─────────────────────────────────────────────────────────┐
│                    STREAMVERSE                          │
│                                                         │
│              ┌─────────────────────┐                    │
│              │   Sign In           │                    │
│              │                     │                    │
│              │  ┌───────────────┐  │                    │
│              │  │ Email         │  │                    │
│              │  └───────────────┘  │                    │
│              │  ┌───────────────┐  │                    │
│              │  │ Password      │  │                    │
│              │  └───────────────┘  │                    │
│              │                     │                    │
│              │  ┌───────────────┐  │                    │
│              │  │  Sign In  [Red]│  │                    │
│              │  └───────────────┘  │                    │
│              │                     │                    │
│              │  New to StreamVerse?│                    │
│              │  Sign up now        │                    │
│              └─────────────────────┘                    │
│                                                         │
│         © 2025 StreamVerse. All rights reserved.       │
└─────────────────────────────────────────────────────────┘
```

**Features:**
- Semi-transparent dark overlay on background
- Email/password authentication
- Toggle between sign-in/sign-up
- OAuth support (Google, Facebook, Apple)
- Remember me option
- Responsive design (mobile-first)

**Color Scheme:**
- Background: Linear gradient with dark overlay
- Form background: rgba(0,0,0,0.75)
- Input fields: #333
- Primary button: #E50914 (Netflix Red)
- Text: #fff, #737373

---

### 2. Home Page

```
┌─────────────────────────────────────────────────────────────────┐
│ STREAMVERSE    Home  Browse                         🔍  👤     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ╔═══════════════════════════════════════════════════════╗    │
│   ║                                                       ║    │
│   ║  STRANGER THINGS                                      ║    │
│   ║                                                       ║    │
│   ║  When a young boy vanishes, a small town uncovers    ║    │
│   ║  a mystery involving secret experiments...           ║    │
│   ║                                                       ║    │
│   ║  ▶ Play     ℹ More Info                             ║    │
│   ║                                                       ║    │
│   ╚═══════════════════════════════════════════════════════╝    │
│                                                                 │
│   Trending Now                                                  │
│   ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ →        │
│   │ Movie │ │ Movie │ │ Movie │ │ Movie │ │ Movie │          │
│   │   1   │ │   2   │ │   3   │ │   4   │ │   5   │          │
│   └───────┘ └───────┘ └───────┘ └───────┘ └───────┘          │
│                                                                 │
│   New Releases                                                  │
│   ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ →        │
│   │ Show  │ │ Show  │ │ Show  │ │ Show  │ │ Show  │          │
│   │   1   │ │   2   │ │   3   │ │   4   │ │   5   │          │
│   └───────┘ └───────┘ └───────┘ └───────┘ └───────┘          │
│                                                                 │
│   Popular on StreamVerse                                        │
│   ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ →        │
│   │ Film  │ │ Film  │ │ Film  │ │ Film  │ │ Film  │          │
│   │   1   │ │   2   │ │   3   │ │   4   │ │   5   │          │
│   └───────┘ └───────┘ └───────┘ └───────┘ └───────┘          │
└─────────────────────────────────────────────────────────────────┘
```

**Features:**
- Fixed navigation bar with gradient fade
- Large hero banner (80vh) with featured content
- Auto-playing trailer preview (muted)
- Multiple horizontal content rows
- Smooth scroll with lazy loading
- Hover effects: Scale 1.05, show metadata
- Infinite scroll capability

**Interactions:**
- Hover on content card → Scale up, show title/rating
- Click content → Navigate to detail/watch page
- Horizontal scroll with mouse or touch
- Keyboard navigation support

---

### 3. Video Player Page

```
┌─────────────────────────────────────────────────────────────────┐
│  ← Back                                                         │
│                                                                 │
│  ╔═══════════════════════════════════════════════════════════╗ │
│  ║                                                           ║ │
│  ║                      ▶                                    ║ │
│  ║                                                           ║ │
│  ║  ────────●──────────────────────────────────  1:23:45    ║ │
│  ║  🔊 ⚙️ CC 🖥️                                              ║ │
│  ╚═══════════════════════════════════════════════════════════╝ │
│                                                                 │
│  The Last Guardian                                              │
│  ⭐ 8.5    2h 15m    2024                                      │
│                                                                 │
│  An epic adventure about courage, friendship, and discovering  │
│  what it truly means to be a hero in a world that needs       │
│  saving. Follow the journey of a young warrior...             │
│                                                                 │
│  [Action] [Adventure] [Sci-Fi] [Fantasy]                       │
│                                                                 │
│  Similar Titles                                                 │
│  ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐                     │
│  │ Movie │ │ Movie │ │ Movie │ │ Movie │                     │
│  └───────┘ └───────┘ └───────┘ └───────┘                     │
└─────────────────────────────────────────────────────────────────┘
```

**Features:**
- Full-width video player (HTML5)
- Custom controls overlay
- Progress bar with preview thumbnails
- Quality selector (Auto, 4K, 1080p, 720p, 480p)
- Playback speed control
- Subtitles/Closed captions
- Picture-in-Picture mode
- Auto-save watch progress
- Next episode countdown (for series)

**Player Controls:**
- Play/Pause, Rewind 10s, Forward 10s
- Volume control with mute
- Fullscreen toggle
- Audio/subtitle track selection
- Chromecast button

---

### 4. Browse/Search Page

```
┌─────────────────────────────────────────────────────────────────┐
│ STREAMVERSE    Home  Browse                         🔍  👤     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌────────────────────────────────────────┐                    │
│  │ 🔍 Search for movies, series...        │                    │
│  └────────────────────────────────────────┘                    │
│                                                                 │
│  Trending Now                                                   │
│                                                                 │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐       │
│  │Movie │ │Movie │ │Movie │ │Movie │ │Movie │ │Movie │       │
│  │  ⭐8.5│ │  ⭐7.8│ │  ⭐9.1│ │  ⭐8.2│ │  ⭐7.5│ │  ⭐8.8│       │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘       │
│                                                                 │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐       │
│  │Series│ │Series│ │Series│ │Series│ │Series│ │Series│       │
│  │  ⭐8.3│ │  ⭐7.9│ │  ⭐8.7│ │  ⭐9.0│ │  ⭐8.1│ │  ⭐7.6│       │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘       │
│                                                                 │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐       │
│  │ Docs │ │ Docs │ │ Docs │ │ Docs │ │ Docs │ │ Docs │       │
│  │  ⭐8.9│ │  ⭐8.4│ │  ⭐7.7│ │  ⭐9.2│ │  ⭐8.0│ │  ⭐8.6│       │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘ └──────┘       │
└─────────────────────────────────────────────────────────────────┘
```

**Features:**
- Real-time search with debouncing
- Search suggestions/autocomplete
- Filter by genre, year, rating, type
- Sort by: Trending, Newest, Rating, A-Z
- Grid view (default) or List view
- Infinite scroll pagination
- Empty state for no results

---

### 5. Profile Selection Page

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│                     Who's watching?                             │
│                                                                 │
│         ┌────────┐    ┌────────┐    ┌────────┐                │
│         │  👤    │    │  👶    │    │   ➕   │                │
│         │        │    │        │    │        │                │
│         │  John  │    │  Kids  │    │  Add   │                │
│         └────────┘    └────────┘    └────────┘                │
│                                                                 │
│                  Manage Profiles                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Features:**
- Multiple user profiles (up to 5)
- Kids profile with content filtering
- Custom avatar selection
- Profile editing
- Profile-specific watch history
- Parental controls per profile

---

## 📱 MOBILE APP UI (Flutter)

### Main Features:
```
┌──────────────────────┐
│ ≡  STREAMVERSE    👤 │
├──────────────────────┤
│                      │
│ ╔══════════════════╗ │
│ ║ Featured Content ║ │
│ ║                  ║ │
│ ║  ▶ Play          ║ │
│ ╚══════════════════╝ │
│                      │
│ Continue Watching    │
│ ← ▓▓ ▓▓ ▓▓ ▓▓ ▓▓ → │
│                      │
│ Trending Now         │
│ ← ▓▓ ▓▓ ▓▓ ▓▓ ▓▓ → │
│                      │
│ New Releases         │
│ ← ▓▓ ▓▓ ▓▓ ▓▓ ▓▓ → │
│                      │
├──────────────────────┤
│ 🏠  🔍  📥  👤      │
└──────────────────────┘
```

**Navigation:**
- 🏠 Home: Featured + content rows
- 🔍 Search: Search + browse
- 📥 Downloads: Offline content
- 👤 Profile: Settings + account

**Mobile-Specific Features:**
- Offline downloads (up to 100 titles)
- Picture-in-Picture mode
- Chromecast support
- Push notifications
- Biometric login (Face ID/Touch ID)
- Swipe gestures
- Adaptive streaming based on connection
- Data saver mode

---

## 📺 TV APP UI (10 Platforms)

### Smart TV Interface (D-Pad Navigation)

```
┌─────────────────────────────────────────────────────────────────┐
│  STREAMVERSE     [Home] [Movies] [Series] [My List]        👤  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ╔══════════════════════════════════════════════════════════╗  │
│  ║                  FEATURED CONTENT                        ║  │
│  ║                                                          ║  │
│  ║                      ▶ Play                              ║  │
│  ║                    + My List                             ║  │
│  ╚══════════════════════════════════════════════════════════╝  │
│                                                                 │
│  ╔═══════╗  ┌───────┐  ┌───────┐  ┌───────┐  ┌───────┐       │
│  ║Movie 1║  │Movie 2│  │Movie 3│  │Movie 4│  │Movie 5│       │
│  ║  ⭐8.5 ║  │  ⭐7.8 │  │  ⭐9.1 │  │  ⭐8.2 │  │  ⭐7.5 │       │
│  ╚═══════╝  └───────┘  └───────┘  └───────┘  └───────┘       │
│   ^ FOCUSED                                                     │
│                                                                 │
│  ┌───────┐  ┌───────┐  ┌───────┐  ┌───────┐  ┌───────┐       │
│  │Series1│  │Series2│  │Series3│  │Series4│  │Series5│       │
│  │  ⭐8.3 │  │  ⭐7.9 │  │  ⭐8.7 │  │  ⭐9.0 │  │  ⭐8.1 │       │
│  └───────┘  └───────┘  └───────┘  └───────┘  └───────┘       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**TV-Specific Features:**
- **D-Pad Navigation**: Arrow keys + OK button
- **Focus Indicators**: Bold border (#E50914) + scale 1.1
- **10-Foot UI**: Large text, high contrast
- **Voice Search**: Platform-specific voice input
- **Screensaver**: Auto-play trailers when idle
- **Remote Control**: Volume, play/pause, channel keys

**Platform Coverage:**
1. **Android TV / Google TV** - Kotlin, Android SDK
2. **Samsung Tizen** - HTML5/JavaScript
3. **LG webOS** - HTML5/JavaScript + webOS SDK
4. **Roku** - BrightScript + SceneGraph
5. **Amazon Fire TV** - Android (Fire TV SDK)
6. **Apple tvOS** - Swift, tvOS SDK
7. **Vizio SmartCast** - HTML5/JavaScript
8. **Hisense VIDAA** - HTML5/JavaScript
9. **Panasonic My Home Screen** - HTML5/JavaScript
10. **Huawei HarmonyOS** - ArkTS (HarmonyOS SDK)

**Navigation Pattern:**
- UP/DOWN: Navigate between rows
- LEFT/RIGHT: Navigate within row
- OK/SELECT: Play content
- BACK: Return to previous screen
- HOME: Return to main menu
- MENU: Open settings

---

## ⚙️ ADMIN DASHBOARD UI

### Dashboard Overview

```
┌─────────────────────────────────────────────────────────────────┐
│ ╔═══════════╗ │ Dashboard Overview         Last updated: 2m ago│
│ ║   ADMIN   ║ │                                                 │
│ ╚═══════════╝ ├─────────────────────────────────────────────────┤
│ ┌───────────┐ │                                                 │
│ │📊Dashboard│ │  ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│ └───────────┘ │  │  1.2M    │ │  45.8K   │ │  8,547   │       │
│ ┌───────────┐ │  │  Users   │ │  Active  │ │  Content │       │
│ │👥 Users   │ │  └──────────┘ └──────────┘ └──────────┘       │
│ └───────────┘ │                                                 │
│ ┌───────────┐ │  ┌──────────┐                                  │
│ │🎬 Content │ │  │  $2.4M   │                                  │
│ └───────────┘ │  │  Revenue │                                  │
│ ┌───────────┐ │  └──────────┘                                  │
│ │💰 Revenue │ │                                                 │
│ └───────────┘ │  Recent Content                                │
│ ┌───────────┐ │  ┌─────────────────────────────────────────┐  │
│ │📈Analytics│ │  │ Title        Type  Views   Status Actions│  │
│ └───────────┘ │  ├─────────────────────────────────────────┤  │
│ ┌───────────┐ │  │ Stranger T.  Series 2.4M  ✓Live   Edit │  │
│ │⚙️Settings │ │  │ Last Guard.  Movie  1.8M  ✓Live   Edit │  │
│ └───────────┘ │  │ Nature Doc.  Doc    890K  ⚠Proc   Edit │  │
│ ┌───────────┐ │  │ Kids Adv.    Series 650K  ⚠Draft  Edit │  │
│ │🔐Security │ │  └─────────────────────────────────────────┘  │
│ └───────────┘ │                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Admin Features:**

### 1. Dashboard
- Real-time statistics
- User growth charts
- Revenue analytics
- Content performance
- System health status

### 2. User Management
- User list with search/filter
- Ban/unban users
- Subscription management
- User activity logs
- Support ticket management

### 3. Content Management
- Upload new content
- Edit metadata (title, description, genres)
- Transcoding status
- Content moderation queue
- Scheduled releases
- Geographic restrictions

### 4. Revenue & Analytics
- Revenue charts (daily/monthly/yearly)
- Subscription metrics
- Churn rate analysis
- Content performance
- User engagement metrics
- A/B test results

### 5. Platform Settings
- CDN configuration
- Streaming quality settings
- DRM settings
- Email templates
- Feature flags
- API rate limits

### 6. Security & Compliance
- Security logs
- Failed login attempts
- API access logs
- GDPR data requests
- Content rating settings
- Parental control config

---

## 🎨 Component Library

### Buttons
```css
/* Primary Button */
.btn-primary {
  background: #E50914;
  color: #fff;
  padding: 14px 32px;
  border-radius: 4px;
  font-weight: bold;
}
.btn-primary:hover {
  background: #f40612;
}

/* Secondary Button */
.btn-secondary {
  background: rgba(109, 109, 110, 0.7);
  color: #fff;
  padding: 14px 32px;
  border-radius: 4px;
}
```

### Cards
```css
.content-card {
  border-radius: 4px;
  overflow: hidden;
  transition: transform 0.2s;
  cursor: pointer;
}
.content-card:hover {
  transform: scale(1.05);
  z-index: 10;
}
```

### Typography
```css
/* Heading 1 */
h1 { font-size: 3.5rem; font-weight: bold; }

/* Heading 2 */
h2 { font-size: 2rem; font-weight: bold; }

/* Body */
body { font-size: 1rem; line-height: 1.5; }
```

---

## 📊 Responsive Breakpoints

```css
/* Mobile First */
@media (max-width: 768px) {
  /* Mobile styles */
  .hero-title { font-size: 2rem; }
  .navbar { padding: 15px 30px; }
}

@media (min-width: 769px) and (max-width: 1024px) {
  /* Tablet styles */
}

@media (min-width: 1025px) {
  /* Desktop styles */
}

@media (min-width: 1920px) {
  /* Large desktop / 4K */
}
```

---

## ✨ Animations & Transitions

### Hover Effects
- Scale: `transform: scale(1.05)`
- Duration: `0.2s ease-in-out`
- Fade: `opacity: 0 → 1`

### Page Transitions
- Fade in: 300ms
- Slide in from bottom: 400ms
- Loading spinner: Rotate 360°

### Scroll Effects
- Lazy load: Fade in when visible
- Parallax: Hero background
- Sticky nav: Fade background on scroll

---

## 📱 Platform-Specific Adaptations

### iOS (Flutter)
- Native navigation bar
- Safe area insets
- Haptic feedback
- Face ID integration

### Android (Flutter)
- Material Design 3
- Back button handling
- Bottom navigation
- Google Sign-In

### TV Platforms
- Large touch targets (48px minimum)
- High contrast text
- Focus management
- Voice search integration

---

## 🎯 Accessibility

- **WCAG 2.1 AA** compliance
- **Keyboard navigation** full support
- **Screen readers** ARIA labels
- **Closed captions** available
- **Audio descriptions** for visually impaired
- **High contrast mode**
- **Font size adjustment**
- **Color blind friendly** palette

---

## 📐 Grid System

```
Desktop (1920px):
┌─────────────────────────────────────────────────────────┐
│ 60px │        Content Area (1800px)              │ 60px │
└─────────────────────────────────────────────────────────┘

Tablet (768px):
┌───────────────────────────────────────────────┐
│ 40px │  Content Area (688px)         │ 40px  │
└───────────────────────────────────────────────┘

Mobile (375px):
┌─────────────────────────────────┐
│ 20px │  Content (335px)  │ 20px │
└─────────────────────────────────┘
```

---

## 🎬 Video Player Features

### Controls
- Play/Pause
- Progress bar with scrubbing
- Volume control
- Quality selector (Auto/4K/1080p/720p/480p)
- Playback speed (0.5x - 2x)
- Subtitles (20+ languages)
- Audio tracks
- Next episode (auto-play)
- Skip intro button
- 10s forward/backward

### Advanced Features
- Thumbnail preview on hover
- Chapter markers
- Watch party (sync viewing)
- Offline viewing
- Chromecast/AirPlay
- Picture-in-Picture

---

**Open `UI_SHOWCASE.html` in your browser to see interactive mockups!**
