import React, { createContext, useContext, useEffect, useState } from 'react'
import { api, isAuthenticated as checkAuth, clearTokens, setTokens } from '../integrations/auth'

interface User {
  id: string;
  displayName: string;
  email: string;
  systemRole: string;
  balance: number;
  createdAt: string;
}

interface AuthContextType {
  isAuthenticated: boolean;
  user: User | null;
  loading: boolean;
  login: (accessToken: string, refreshToken: string) => void;
  logout: () => void;
  check: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(checkAuth())
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const update = () => setIsAuthenticated(checkAuth())
    window.addEventListener('auth-changed', update)
    return () => window.removeEventListener('auth-changed', update)
  }, [])

  useEffect(() => {
    setLoading(true);
    if (isAuthenticated) {
      api.get('/admin')
        .then(response => {
          setUser(response.data)
        })
        .catch(() => {
          setUser(null)
        })
        .finally(() => {
          setLoading(false);
        });
    } else {
      setUser(null)
      setLoading(false);
    }
  }, [isAuthenticated])

  const login = (accessToken: string, refreshToken: string) => {
    setTokens(accessToken, refreshToken)
    setIsAuthenticated(true)
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout', { refreshToken: localStorage.getItem('refreshToken') })
    } catch {}
    clearTokens()
    setIsAuthenticated(false)
    setUser(null);
  }

  const check = () => setIsAuthenticated(checkAuth())

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, loading, login, logout, check }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
} 