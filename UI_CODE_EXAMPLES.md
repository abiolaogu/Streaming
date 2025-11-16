# StreamVerse UI Code Examples

This document shows actual code snippets from each platform to demonstrate the technical implementation.

---

## 💻 WEB UI (React + TypeScript)

### Home Page Component
```tsx
// src/pages/HomePage.tsx
import { useState, useEffect } from 'react'
import { api } from '@/services/api'
import type { Content } from '@/types'
import NavBar from '@/components/NavBar'
import HeroBanner from '@/components/HeroBanner'
import ContentRow from '@/components/ContentRow'

export default function HomePage({ onWatchContent, onNavigate, onLogout }) {
  const [categories, setCategories] = useState([])
  const [featuredContent, setFeaturedContent] = useState(null)

  useEffect(() => {
    loadContent()
  }, [])

  return (
    <div className="home-page">
      <NavBar onNavigate={onNavigate} onLogout={onLogout} />
      
      {featuredContent && (
        <HeroBanner 
          content={featuredContent} 
          onPlay={() => onWatchContent(featuredContent.id)} 
        />
      )}

      <div className="content-rows">
        {categories.map((category) => (
          <ContentRow
            key={category.name}
            title={category.name}
            content={category.content}
            onContentClick={onWatchContent}
          />
        ))}
      </div>
    </div>
  )
}
```

### Styling (CSS)
```css
/* src/pages/HomePage.css */
.home-page {
  min-height: 100vh;
  background: #141414;
}

.content-rows {
  position: relative;
  z-index: 2;
  margin-top: -150px;  /* Overlap with hero */
  padding-bottom: 50px;
}
```

### Hero Banner Component
```tsx
// src/components/HeroBanner.tsx
export default function HeroBanner({ content, onPlay }) {
  return (
    <div 
      className="hero-banner" 
      style={{ backgroundImage: `url(${content.thumbnail_url})` }}
    >
      <div className="hero-overlay">
        <div className="hero-content">
          <h1 className="hero-title">{content.title}</h1>
          <p className="hero-description">{content.description}</p>

          <div className="hero-actions">
            <button className="play-button" onClick={onPlay}>
              <span className="play-icon">▶</span>
              Play
            </button>
            <button className="info-button">
              <span className="info-icon">ℹ</span>
              More Info
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
```

---

## 📱 MOBILE APP (Flutter/Dart)

### Main App Structure
```dart
// apps/clients/mobile-flutter/lib/main.dart
void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp();
  await Hive.initFlutter();
  
  runApp(const ProviderScope(child: StreamVerseApp()));
}

class StreamVerseApp extends StatelessWidget {
  const StreamVerseApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'StreamVerse',
      theme: ThemeData.dark().copyWith(
        primaryColor: const Color(0xFFE50914),
        scaffoldBackgroundColor: Colors.black,
        colorScheme: const ColorScheme.dark(
          primary: Color(0xFFE50914),
          secondary: Color(0xFF141414),
        ),
      ),
      home: const SplashScreen(),
      routes: {
        '/home': (context) => const HomeScreen(),
        '/watch': (context) => const WatchScreen(),
        '/browse': (context) => const BrowseScreen(),
        '/profile': (context) => const ProfileScreen(),
      },
    );
  }
}
```

### Home Screen Widget
```dart
// lib/screens/home_screen.dart
class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            floating: true,
            backgroundColor: Colors.black.withOpacity(0.7),
            title: const Text('STREAMVERSE'),
            actions: [
              IconButton(
                icon: const Icon(Icons.search),
                onPressed: () => Navigator.pushNamed(context, '/browse'),
              ),
              IconButton(
                icon: const Icon(Icons.account_circle),
                onPressed: () => Navigator.pushNamed(context, '/profile'),
              ),
            ],
          ),
          SliverToBoxAdapter(
            child: Column(
              children: [
                const HeroBanner(),
                ContentRow(title: 'Trending Now'),
                ContentRow(title: 'New Releases'),
                ContentRow(title: 'Popular on StreamVerse'),
              ],
            ),
          ),
        ],
      ),
      bottomNavigationBar: BottomNavigationBar(
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
          BottomNavigationBarItem(icon: Icon(Icons.search), label: 'Search'),
          BottomNavigationBarItem(icon: Icon(Icons.download), label: 'Downloads'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
        ],
      ),
    );
  }
}
```

---

## 📺 TV APP - Roku (BrightScript)

### Main Scene
```brightscript
' apps/clients/tv-apps/roku/components/MainScene.brs
sub init()
    print "MainScene: init"
    
    ' Set background
    m.top.backgroundURI = ""
    m.top.backgroundColor = "0x000000"
    
    ' Get UI components
    m.navigationMenu = m.top.findNode("navigationMenu")
    m.contentRowList = m.top.findNode("contentRowList")
    m.videoPlayer = m.top.findNode("videoPlayer")
    
    ' Set up event handlers
    m.contentRowList.observeField("itemFocused", "onContentItemFocused")
    m.contentRowList.observeField("itemSelected", "onContentItemSelected")
    
    ' Initialize
    m.currentScreen = "home"
    m.isLoggedIn = false
    
    ' Load content
    if m.isLoggedIn then
        loadHomeContent()
    else
        showLoginScreen()
    end if
end sub

sub loadHomeContent()
    m.loadingSpinner.visible = true
    
    urlTransfer = CreateObject("roUrlTransfer")
    urlTransfer.SetUrl("https://api.streamverse.io/v1/content/home")
    urlTransfer.AddHeader("Authorization", "Bearer " + m.authToken)
    
    if urlTransfer.AsyncGetToString() then
        msg = wait(5000, port)
        if type(msg) = "roUrlEvent" then
            if msg.GetResponseCode() = 200 then
                response = ParseJson(msg.GetString())
                displayHomeContent(response)
            end if
        end if
    end if
    
    m.loadingSpinner.visible = false
end sub

sub onContentItemFocused()
    ' Show content details when focused
    focusedItem = m.contentRowList.focusedContent
    if focusedItem <> invalid then
        showContentPreview(focusedItem)
    end if
end sub

sub onContentItemSelected()
    ' Play content when selected
    selectedItem = m.contentRowList.selectedContent
    if selectedItem <> invalid then
        playContent(selectedItem)
    end if
end sub
```

### Scene XML Layout
```xml
<!-- apps/clients/tv-apps/roku/components/MainScene.xml -->
<?xml version="1.0" encoding="utf-8" ?>
<component name="MainScene" extends="Scene">
  <script type="text/brightscript" uri="pkg:/components/MainScene.brs" />
  
  <children>
    <!-- Navigation Menu -->
    <LabelList id="navigationMenu"
      translation="[60, 60]"
      itemSize="[300, 60]">
      <ContentNode role="content">
        <ContentNode title="Home" />
        <ContentNode title="Movies" />
        <ContentNode title="Series" />
        <ContentNode title="My List" />
      </ContentNode>
    </LabelList>
    
    <!-- Content Rows -->
    <RowList id="contentRowList"
      translation="[0, 300]"
      itemComponentName="ContentRow"
      focusXOffset="[100]"
      numRows="3" />
    
    <!-- Video Player -->
    <Video id="videoPlayer"
      width="1920"
      height="1080"
      visible="false" />
    
    <!-- Loading Spinner -->
    <BusySpinner id="loadingSpinner"
      translation="[900, 500]"
      visible="false" />
  </children>
</component>
```

---

## 📺 TV APP - LG webOS (HTML5/JavaScript)

### Main Application
```html
<!-- apps/clients/tv-apps/lg-webos/index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>StreamVerse</title>
    <link rel="stylesheet" href="assets/css/app.css">
    <script src="https://webostv.sdk.lgappstv.com/webOS.js"></script>
</head>
<body>
    <div id="app">
        <nav class="tv-nav">
            <div class="nav-item focused">Home</div>
            <div class="nav-item">Movies</div>
            <div class="nav-item">Series</div>
            <div class="nav-item">My List</div>
        </nav>
        
        <div class="hero-banner">
            <h1 id="hero-title"></h1>
            <p id="hero-description"></p>
            <button class="play-btn">Play</button>
        </div>
        
        <div id="content-rows"></div>
    </div>
    
    <script src="dist/bundle.js"></script>
</body>
</html>
```

### JavaScript Controller
```javascript
// apps/clients/tv-apps/lg-webos/src/app.js
class StreamVerseApp {
    constructor() {
        this.currentFocus = { row: 0, col: 0 };
        this.init();
    }

    init() {
        // Initialize webOS
        if (window.webOS) {
            webOS.platformBack = () => this.handleBack();
        }

        // Set up keyboard navigation
        document.addEventListener('keydown', (e) => this.handleKeyDown(e));

        // Load content
        this.loadHomeContent();
    }

    handleKeyDown(event) {
        switch(event.keyCode) {
            case 37: // LEFT
                this.moveFocus('left');
                break;
            case 38: // UP
                this.moveFocus('up');
                break;
            case 39: // RIGHT
                this.moveFocus('right');
                break;
            case 40: // DOWN
                this.moveFocus('down');
                break;
            case 13: // OK/ENTER
                this.selectContent();
                break;
            case 10009: // BACK
                this.handleBack();
                break;
        }
    }

    moveFocus(direction) {
        // Update focus position
        const previousFocus = document.querySelector('.focused');
        if (previousFocus) {
            previousFocus.classList.remove('focused');
        }

        // Calculate new focus
        if (direction === 'left') {
            this.currentFocus.col = Math.max(0, this.currentFocus.col - 1);
        } else if (direction === 'right') {
            this.currentFocus.col++;
        } else if (direction === 'up') {
            this.currentFocus.row = Math.max(0, this.currentFocus.row - 1);
        } else if (direction === 'down') {
            this.currentFocus.row++;
        }

        // Apply new focus
        const newFocus = this.getElementAtPosition(this.currentFocus);
        if (newFocus) {
            newFocus.classList.add('focused');
            newFocus.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    }

    async loadHomeContent() {
        try {
            const response = await fetch('https://api.streamverse.io/v1/content/home', {
                headers: {
                    'Authorization': `Bearer ${this.authToken}`
                }
            });
            const data = await response.json();
            this.renderContent(data);
        } catch (error) {
            console.error('Failed to load content:', error);
        }
    }

    renderContent(data) {
        // Render hero banner
        document.getElementById('hero-title').textContent = data.featured.title;
        document.getElementById('hero-description').textContent = data.featured.description;

        // Render content rows
        const contentRows = document.getElementById('content-rows');
        data.categories.forEach(category => {
            const row = this.createContentRow(category);
            contentRows.appendChild(row);
        });
    }
}

// Initialize app
const app = new StreamVerseApp();
```

### CSS Styling
```css
/* apps/clients/tv-apps/lg-webos/assets/css/app.css */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    background: #000;
    color: #fff;
    font-family: 'Helvetica Neue', Arial, sans-serif;
    overflow: hidden;
}

.tv-nav {
    display: flex;
    gap: 40px;
    padding: 40px 60px;
    font-size: 28px;
}

.nav-item {
    padding: 10px 20px;
    border-bottom: 4px solid transparent;
    transition: all 0.3s;
}

.nav-item.focused {
    border-bottom-color: #E50914;
    transform: scale(1.1);
}

.content-card {
    width: 400px;
    height: 225px;
    border-radius: 8px;
    border: 4px solid transparent;
    transition: all 0.3s;
}

.content-card.focused {
    border-color: #E50914;
    transform: scale(1.15);
    box-shadow: 0 0 30px rgba(229, 9, 20, 0.6);
}

.play-btn {
    padding: 20px 50px;
    font-size: 32px;
    background: #E50914;
    color: #fff;
    border: none;
    border-radius: 8px;
    cursor: pointer;
}

.play-btn.focused {
    background: #f40612;
    transform: scale(1.1);
}
```

---

## 📺 TV APP - Android TV (Kotlin)

### Main Activity
```kotlin
// apps/clients/tv-apps/android-tv/app/src/main/java/com/streamverse/tv/MainActivity.kt
class MainActivity : FragmentActivity() {
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        
        if (savedInstanceState == null) {
            supportFragmentManager.beginTransaction()
                .replace(R.id.main_frame, MainFragment())
                .commit()
        }
    }
}

class MainFragment : BrowseSupportFragment() {
    
    private lateinit var rowsAdapter: ArrayObjectAdapter
    
    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        
        setupUI()
        loadContent()
    }
    
    private fun setupUI() {
        // Set branding
        headersState = HEADERS_ENABLED
        isHeadersTransitionOnBackEnabled = true
        brandColor = ContextCompat.getColor(requireContext(), R.color.netflix_red)
        
        // Set up rows adapter
        rowsAdapter = ArrayObjectAdapter(ListRowPresenter())
        adapter = rowsAdapter
        
        // Set up click listener
        onItemViewClickedListener = ItemViewClickedListener()
    }
    
    private fun loadContent() {
        viewLifecycleOwner.lifecycleScope.launch {
            try {
                val content = ApiService.getHomeContent()
                
                content.categories.forEach { category ->
                    val listRowAdapter = ArrayObjectAdapter(CardPresenter())
                    category.items.forEach { item ->
                        listRowAdapter.add(item)
                    }
                    
                    val header = HeaderItem(category.name)
                    rowsAdapter.add(ListRow(header, listRowAdapter))
                }
            } catch (e: Exception) {
                Log.e("MainFragment", "Failed to load content", e)
            }
        }
    }
    
    private inner class ItemViewClickedListener : OnItemViewClickedListener {
        override fun onItemClicked(
            itemViewHolder: Presenter.ViewHolder,
            item: Any,
            rowViewHolder: RowPresenter.ViewHolder,
            row: Row
        ) {
            if (item is Content) {
                val intent = Intent(requireContext(), PlayerActivity::class.java)
                intent.putExtra("content_id", item.id)
                startActivity(intent)
            }
        }
    }
}
```

### Card Presenter
```kotlin
class CardPresenter : Presenter() {
    
    override fun onCreateViewHolder(parent: ViewGroup): ViewHolder {
        val cardView = ImageCardView(parent.context).apply {
            setMainImageDimensions(CARD_WIDTH, CARD_HEIGHT)
            isFocusable = true
            isFocusableInTouchMode = true
        }
        return ViewHolder(cardView)
    }
    
    override fun onBindViewHolder(viewHolder: ViewHolder, item: Any) {
        val content = item as Content
        val cardView = viewHolder.view as ImageCardView
        
        cardView.titleText = content.title
        cardView.contentText = "★ ${content.rating}"
        
        Glide.with(cardView.context)
            .load(content.thumbnailUrl)
            .into(cardView.mainImageView)
    }
    
    override fun onUnbindViewHolder(viewHolder: ViewHolder) {
        val cardView = viewHolder.view as ImageCardView
        cardView.mainImage = null
    }
    
    companion object {
        private const val CARD_WIDTH = 400
        private const val CARD_HEIGHT = 225
    }
}
```

---

## ⚙️ ADMIN DASHBOARD

### Dashboard Component
```tsx
// Admin dashboard example
import { useState, useEffect } from 'react'

export default function AdminDashboard() {
  const [stats, setStats] = useState({
    totalUsers: 0,
    activeUsers: 0,
    totalContent: 0,
    monthlyRevenue: 0
  })

  return (
    <div className="admin-dashboard">
      <aside className="sidebar">
        <h2>ADMIN</h2>
        <nav>
          <NavItem icon="📊" label="Dashboard" active />
          <NavItem icon="👥" label="Users" />
          <NavItem icon="🎬" label="Content" />
          <NavItem icon="💰" label="Revenue" />
        </nav>
      </aside>

      <main className="main-content">
        <header>
          <h1>Dashboard Overview</h1>
        </header>

        <div className="stats-grid">
          <StatCard 
            label="Total Users" 
            value="1.2M" 
            color="#E50914" 
          />
          <StatCard 
            label="Active Now" 
            value="45.8K" 
            color="#16a34a" 
          />
          <StatCard 
            label="Content Items" 
            value="8,547" 
            color="#3b82f6" 
          />
          <StatCard 
            label="Monthly Revenue" 
            value="$2.4M" 
            color="#d97706" 
          />
        </div>

        <ContentTable />
      </main>
    </div>
  )
}
```

---

## 🎨 Shared Styling

### Global CSS Variables
```css
:root {
  /* Colors */
  --netflix-red: #E50914;
  --netflix-black: #141414;
  --netflix-dark: #000000;
  --netflix-gray: #808080;
  --netflix-light-gray: #E5E5E5;
  
  /* Spacing */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;
  
  /* Typography */
  --font-size-xs: 0.75rem;
  --font-size-sm: 0.875rem;
  --font-size-base: 1rem;
  --font-size-lg: 1.125rem;
  --font-size-xl: 1.25rem;
  --font-size-2xl: 1.5rem;
  --font-size-3xl: 2rem;
  
  /* Transitions */
  --transition-fast: 0.15s ease-in-out;
  --transition-base: 0.2s ease-in-out;
  --transition-slow: 0.3s ease-in-out;
}
```

---

**This shows the actual code structure across all platforms!**
