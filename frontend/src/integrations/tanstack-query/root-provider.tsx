import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React from "react";

const queryClient = new QueryClient()

export function getContext() {
  return {
    queryClient,
  }
}

export function Provider({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}
