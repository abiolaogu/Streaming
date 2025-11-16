import { useState, useEffect } from 'react'
import HomePage from './pages/HomePage'
import LoginPage from './pages/LoginPage'
import WatchPage from './pages/WatchPage'
import BrowsePage from './pages/BrowsePage'
import ProfilesPage from './pages/ProfilesPage'
import { api } from './services/api'

type Page = 'home' | 'login' | 'watch' | 'browse' | 'profiles'

function App() {
  const [currentPage, setCurrentPage] = useState<Page>('login')
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [selectedContentId, setSelectedContentId] = useState<string | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      setIsAuthenticated(true)
      setCurrentPage('home')
    }
  }, [])

  const handleLogin = () => {
    setIsAuthenticated(true)
    setCurrentPage('profiles')
  }

  const handleLogout = async () => {
    await api.logout()
    setIsAuthenticated(false)
    setCurrentPage('login')
  }

  const handleWatchContent = (contentId: string) => {
    setSelectedContentId(contentId)
    setCurrentPage('watch')
  }

  const handleNavigate = (page: Page) => {
    setCurrentPage(page)
  }

  if (!isAuthenticated) {
    return <LoginPage onLogin={handleLogin} />
  }

  return (
    <div className="app">
      {currentPage === 'home' && (
        <HomePage
          onWatchContent={handleWatchContent}
          onNavigate={handleNavigate}
          onLogout={handleLogout}
        />
      )}
      {currentPage === 'browse' && (
        <BrowsePage
          onWatchContent={handleWatchContent}
          onNavigate={handleNavigate}
          onLogout={handleLogout}
        />
      )}
      {currentPage === 'watch' && selectedContentId && (
        <WatchPage
          contentId={selectedContentId}
          onBack={() => setCurrentPage('home')}
        />
      )}
      {currentPage === 'profiles' && (
        <ProfilesPage onSelectProfile={() => setCurrentPage('home')} />
      )}
    </div>
  )
}

export default App
