import { Navigate } from '@tanstack/react-router'
import { roleUtils } from '../utils/role'
import { useAuth } from '../context/AuthContext'
import type { ReactNode } from 'react'

type UserRole = 'client' | 'freelancer'

interface ProtectedRouteProps {
  children: ReactNode
  allowedRoles?: Array<UserRole>
  fallbackPath?: string
}

export default function ProtectedRoute({ 
  children, 
  allowedRoles = [],
  fallbackPath = '/'
}: ProtectedRouteProps) {
  const { isAuthenticated } = useAuth()
  const currentRole = roleUtils.getRole()
  
  if (!isAuthenticated) {
    return <Navigate to="/login" />
  }
  
  if (allowedRoles.length > 0 && !allowedRoles.includes(currentRole)) {
    return <Navigate to={fallbackPath} />
  }

  return <>{children}</>
} 