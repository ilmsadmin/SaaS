'use client'

import { useAuth } from '@/lib/auth-context'
import { usePathname } from 'next/navigation'
import { useState, useEffect } from 'react'

export default function AuthDebugger() {
  const { user, isLoading, isAuthenticated } = useAuth()
  const pathname = usePathname()
  const [isMounted, setIsMounted] = useState(false)
  const [token, setToken] = useState<string | null>(null)

  useEffect(() => {
    setIsMounted(true)
    setToken(localStorage.getItem('access_token'))
  }, [])

  // Only show in development and after client-side mounting
  if (process.env.NODE_ENV !== 'development' || !isMounted) {
    return null
  }

  return (
    <div className="fixed bottom-4 right-4 bg-black bg-opacity-80 text-white p-4 rounded-lg text-xs max-w-sm z-50">
      <h3 className="font-bold mb-2">Auth Debug</h3>
      <div className="space-y-1">
        <div>
          <strong>Path:</strong> {pathname}
        </div>
        <div>
          <strong>Loading:</strong> {isLoading ? 'Yes' : 'No'}
        </div>
        <div>
          <strong>Authenticated:</strong> {isAuthenticated ? 'Yes' : 'No'}
        </div>
        <div>
          <strong>User:</strong> {user ? `${user.first_name} ${user.last_name}` : 'None'}
        </div>
        <div>
          <strong>Token:</strong> {token ? 'Present' : 'None'}
        </div>
      </div>
    </div>
  )
}
