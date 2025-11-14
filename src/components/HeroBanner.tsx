import type { Content } from '@/types'
import './HeroBanner.css'

interface HeroBannerProps {
  content: Content
  onPlay: () => void
}

export default function HeroBanner({ content, onPlay }: HeroBannerProps) {
  return (
    <div className="hero-banner" style={{ backgroundImage: `url(${content.thumbnail_url})` }}>
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

          <div className="hero-meta">
            <span className="rating">★ {content.rating.toFixed(1)}</span>
            <span className="year">{new Date(content.release_date).getFullYear()}</span>
            <span className="duration">{Math.floor(content.duration / 60)}m</span>
          </div>
        </div>
      </div>
    </div>
  )
}
