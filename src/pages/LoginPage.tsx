import { useState } from 'react'
import { api } from '@/services/api'
import './LoginPage.css'

interface LoginPageProps {
  onLogin: () => void
}

export default function LoginPage({ onLogin }: LoginPageProps) {
  const [isSignUp, setIsSignUp] = useState(false)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [username, setUsername] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      if (isSignUp) {
        await api.register(email, password, username)
      } else {
        await api.login(email, password)
      }
      onLogin()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-header">
        <h1 className="logo">STREAMVERSE</h1>
      </div>

      <div className="login-content">
        <div className="login-form-container">
          <h2>{isSignUp ? 'Sign Up' : 'Sign In'}</h2>

          <form onSubmit={handleSubmit}>
            {isSignUp && (
              <input
                type="text"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            )}

            <input
              type="email"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />

            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={6}
            />

            {error && <div className="error-message">{error}</div>}

            <button type="submit" disabled={loading} className="submit-button">
              {loading ? 'Loading...' : isSignUp ? 'Sign Up' : 'Sign In'}
            </button>
          </form>

          <div className="form-switch">
            <span>
              {isSignUp ? 'Already have an account?' : 'New to StreamVerse?'}
            </span>
            <button onClick={() => setIsSignUp(!isSignUp)} className="switch-button">
              {isSignUp ? 'Sign in now' : 'Sign up now'}
            </button>
          </div>
        </div>
      </div>

      <div className="login-footer">
        <p>&copy; 2025 StreamVerse. All rights reserved.</p>
      </div>
    </div>
  )
}
