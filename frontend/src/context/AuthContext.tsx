import React, { createContext, useContext, useEffect, useState } from 'react'
import { api, isAuthenticated as checkAuth, clearTokens, setTokens } from '../integrations/auth'

interface AuthContextType {
  isAuthenticated: boolean
  login: (accessToken: string, refreshToken: string) => void
  logout: () => void
  check: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(checkAuth())

  useEffect(() => {
    const update = () => setIsAuthenticated(checkAuth())
    window.addEventListener('auth-changed', update)
    return () => window.removeEventListener('auth-changed', update)
  }, [])

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
  }

  const check = () => setIsAuthenticated(checkAuth())

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout, check }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
} 