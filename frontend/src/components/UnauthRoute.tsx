import { useNavigate } from '@tanstack/react-router'
import React from "react";
import { useAuth } from '../context/AuthContext'

export default function UnauthRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  useNavigate();
  return !isAuthenticated ? <>{children}</> : null
} 