import { useState, useEffect } from 'react'
import { api } from '@/services/api'
import type { Content } from '@/types'
import type { Page } from '@/types/navigation'
import NavBar from '@/components/NavBar'
import HeroBanner from '@/components/HeroBanner'
import ContentRow from '@/components/ContentRow'
import './HomePage.css'

interface HomePageProps {
  onWatchContent: (contentId: string) => void
  onNavigate: (page: Page) => void
  onLogout: () => void
}

export default function HomePage({ onWatchContent, onNavigate, onLogout }: HomePageProps) {
  const [categories, setCategories] = useState<{ name: string; content: Content[] }[]>([])
  const [featuredContent, setFeaturedContent] = useState<Content | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadContent()
  }, [])

  const loadContent = async () => {
    try {
      const data = await api.getHomeContent()
      setCategories(data.categories)
      if (data.categories.length > 0 && data.categories[0].content.length > 0) {
        setFeaturedContent(data.categories[0].content[0])
      }
    } catch (error) {
      console.error('Failed to load content:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="loading">Loading...</div>
  }

  return (
    <div className="home-page">
      <NavBar onNavigate={onNavigate} onLogout={onLogout} />

      {featuredContent && (
        <HeroBanner content={featuredContent} onPlay={() => onWatchContent(featuredContent.id)} />
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
