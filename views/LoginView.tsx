import React, { useState } from 'react';
import { SparklesIcon, GoogleIcon, AppleIcon } from '../components/IconComponents';

interface LoginViewProps {
  onLogin: () => void;
}

const LoginView: React.FC<LoginViewProps> = ({ onLogin }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      setError('Please enter both email and password.');
      return;
    }
    setError('');
    // In a real app, you'd call an API here.
    // For this simulation, any non-empty fields will log in.
    onLogin();
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-brand-bg p-4">
      <div className="w-full max-w-md p-8 space-y-8 bg-brand-surface rounded-2xl border border-brand-border shadow-2xl">
        <div className="text-center">
            <div className="flex items-center justify-center mb-4">
                <SparklesIcon className="h-12 w-12 text-brand-accent"/>
                <span className="ml-3 text-3xl font-bold text-brand-text-primary tracking-wider">StreamVerse</span>
            </div>
          <h2 className="text-xl text-brand-text-secondary">Sign in to the Platform</h2>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="email-address" className="sr-only">Email address</label>
              <input
                id="email-address"
                name="email"
                type="email"
                autoComplete="email"
                required
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-brand-border bg-brand-bg placeholder-brand-text-secondary text-brand-text-primary rounded-t-md focus:outline-none focus:ring-brand-accent focus:border-brand-accent focus:z-10 sm:text-sm"
                placeholder="Email address"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="password" className="sr-only">Password</label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-brand-border bg-brand-bg placeholder-brand-text-secondary text-brand-text-primary rounded-b-md focus:outline-none focus:ring-brand-accent focus:border-brand-accent focus:z-10 sm:text-sm"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
          </div>

          {error && <p className="text-sm text-brand-danger text-center">{error}</p>}

          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <input id="remember-me" name="remember-me" type="checkbox" className="h-4 w-4 text-brand-accent bg-brand-bg border-brand-border rounded focus:ring-brand-accent" />
              <label htmlFor="remember-me" className="ml-2 block text-sm text-brand-text-secondary">Remember me</label>
            </div>
            <div className="text-sm">
              <a href="#" className="font-medium text-brand-accent hover:text-brand-accent-hover">Forgot your password?</a>
            </div>
          </div>

          <div>
            <button type="submit" className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-brand-accent hover:bg-brand-accent-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-brand-bg focus:ring-brand-accent">
              Sign in
            </button>
          </div>
        </form>

        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-brand-border" />
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-brand-surface text-brand-text-secondary">Or continue with</span>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-3">
            <button className="w-full inline-flex justify-center py-2 px-4 border border-brand-border rounded-md shadow-sm bg-brand-bg text-sm font-medium text-brand-text-primary hover:bg-brand-border">
                <GoogleIcon />
            </button>
            <button className="w-full inline-flex justify-center py-2 px-4 border border-brand-border rounded-md shadow-sm bg-brand-bg text-sm font-medium text-brand-text-primary hover:bg-brand-border">
                <AppleIcon />
            </button>
        </div>

      </div>
    </div>
  );
};

export default LoginView;
