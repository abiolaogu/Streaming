import { useState, useEffect, useRef } from 'react'
import { api } from '@/services/api'
import type { Content } from '@/types'
import './WatchPage.css'

interface WatchPageProps {
  contentId: string
  onBack: () => void
}

export default function WatchPage({ contentId, onBack }: WatchPageProps) {
  const [content, setContent] = useState<Content | null>(null)
  const [loading, setLoading] = useState(true)
  const videoRef = useRef<HTMLVideoElement>(null)

  useEffect(() => {
    loadContent()
  }, [contentId])

  useEffect(() => {
    const video = videoRef.current
    if (!video) return

    const handleTimeUpdate = () => {
      const progress = (video.currentTime / video.duration) * 100
      api.updateWatchProgress(contentId, progress).catch(console.error)
    }

    video.addEventListener('timeupdate', handleTimeUpdate)
    return () => video.removeEventListener('timeupdate', handleTimeUpdate)
  }, [contentId])

  const loadContent = async () => {
    try {
      const data = await api.getContent(contentId)
      setContent(data)
    } catch (error) {
      console.error('Failed to load content:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="loading">Loading...</div>
  }

  if (!content) {
    return <div className="error">Content not found</div>
  }

  return (
    <div className="watch-page">
      <button className="back-button" onClick={onBack}>
        ← Back
      </button>

      <div className="video-container">
        <video
          ref={videoRef}
          controls
          autoPlay
          className="video-player"
          poster={content.thumbnail_url}
        >
          <source src={content.video_url} type="video/mp4" />
          Your browser does not support the video tag.
        </video>
      </div>

      <div className="content-details">
        <h1>{content.title}</h1>
        <div className="content-meta">
          <span className="rating">★ {content.rating.toFixed(1)}</span>
          <span className="views">{content.views.toLocaleString()} views</span>
          <span className="duration">{Math.floor(content.duration / 60)}m</span>
        </div>
        <p className="description">{content.description}</p>
        <div className="genres">
          {content.genres.map((genre) => (
            <span key={genre} className="genre-tag">
              {genre}
            </span>
          ))}
        </div>
      </div>
    </div>
  )
}
