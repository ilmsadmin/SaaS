'use client'

import { useAdmin } from '@/contexts/AdminContext'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'

interface ProtectedRouteProps {
  children: React.ReactNode
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAdmin()
  const router = useRouter()

  useEffect(() => {
    console.log('ProtectedRoute - isLoading:', isLoading, 'isAuthenticated:', isAuthenticated)
    
    if (!isLoading && !isAuthenticated) {
      console.log('Not authenticated, redirecting to login')
      router.push('/login')
    }
  }, [isAuthenticated, isLoading, router])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return null // Will redirect in useEffect
  }

  return <>{children}</>
}
