'use client'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactNode, useState } from 'react'
import { AdminProvider } from '@/contexts/AdminContext'

interface QueryProvidersProps {
  children: ReactNode
}

export function QueryProviders({ children }: QueryProvidersProps) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000,
            refetchOnWindowFocus: false,
          },
        },
      })
  )

  return (
    <QueryClientProvider client={queryClient}>
      <AdminProvider>
        {children}
      </AdminProvider>
    </QueryClientProvider>
  )
}
