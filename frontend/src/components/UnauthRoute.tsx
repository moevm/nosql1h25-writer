import { useNavigate } from '@tanstack/react-router'
import React, { useEffect } from "react";
import { useAuth } from '../context/AuthContext'

export default function UnauthRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated) {
      navigate({ to: '/profile' });
    }
  }, [isAuthenticated, navigate]);

  return !isAuthenticated ? <>{children}</> : null;
} 