import { useState } from 'react'
import './NavBar.css'

interface NavBarProps {
  onNavigate: (page: string) => void
  onLogout: () => void
}

export default function NavBar({ onNavigate, onLogout }: NavBarProps) {
  const [showMenu, setShowMenu] = useState(false)

  return (
    <nav className="navbar">
      <div className="navbar-left">
        <h1 className="logo" onClick={() => onNavigate('home')}>
          STREAMVERSE
        </h1>

        <div className="nav-links">
          <button onClick={() => onNavigate('home')}>Home</button>
          <button onClick={() => onNavigate('browse')}>Browse</button>
        </div>
      </div>

      <div className="navbar-right">
        <div className="user-menu">
          <button
            className="user-avatar"
            onClick={() => setShowMenu(!showMenu)}
          >
            <span>👤</span>
          </button>

          {showMenu && (
            <div className="dropdown-menu">
              <button onClick={() => onNavigate('profiles')}>
                Switch Profile
              </button>
              <button onClick={onLogout}>Sign Out</button>
            </div>
          )}
        </div>
      </div>
    </nav>
  )
}
