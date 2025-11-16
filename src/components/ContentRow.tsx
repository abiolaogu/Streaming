import type { Content } from '@/types'
import './ContentRow.css'

interface ContentRowProps {
  title: string
  content: Content[]
  onContentClick: (contentId: string) => void
}

export default function ContentRow({ title, content, onContentClick }: ContentRowProps) {
  return (
    <div className="content-row">
      <h2 className="row-title">{title}</h2>

      <div className="row-container">
        <div className="row-content">
          {content.map((item) => (
            <div
              key={item.id}
              className="content-item"
              onClick={() => onContentClick(item.id)}
            >
              <img src={item.thumbnail_url} alt={item.title} />
              <div className="item-overlay">
                <h3>{item.title}</h3>
                <div className="item-info">
                  <span className="rating">★ {item.rating.toFixed(1)}</span>
                  <span className="duration">{Math.floor(item.duration / 60)}m</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
