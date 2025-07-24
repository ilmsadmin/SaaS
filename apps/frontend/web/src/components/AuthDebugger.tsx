'use client'

import { useAuth } from '@/lib/auth-context'
import { usePathname } from 'next/navigation'

export default function AuthDebugger() {
  const { user, isLoading, isAuthenticated } = useAuth()
  const pathname = usePathname()

  // Only show in development
  if (process.env.NODE_ENV !== 'development') {
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
          <strong>Token:</strong> {typeof window !== 'undefined' && localStorage.getItem('access_token') ? 'Present' : 'None'}
        </div>
      </div>
    </div>
  )
}
