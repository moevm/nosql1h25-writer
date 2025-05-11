import { useAuth } from '../context/AuthContext'

export default function AuthOnly({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  return isAuthenticated ? <>{children}</> : null
} 