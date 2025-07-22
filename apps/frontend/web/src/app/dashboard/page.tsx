'use client'

import { useAuth } from '@/lib/auth-context'
import { useRequireAuth } from '@/lib/auth-context'

export default function DashboardPage() {
  const { isAuthenticated, isLoading } = useRequireAuth()
  const { user, logout } = useAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Loading...</div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return null // Will redirect to login
  }

  const handleLogout = async () => {
    try {
      await logout()
      window.location.href = '/login'
    } catch (error) {
      console.error('Logout failed:', error)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-semibold">Zplus SaaS Dashboard</h1>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-gray-700">
                Welcome, {user?.first_name} {user?.last_name}
              </span>
              <button
                onClick={handleLogout}
                className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* User Info Card */}
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-5">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                  Profile Information
                </h3>
                <div className="mt-4 space-y-2">
                  <p className="text-sm text-gray-600">
                    <span className="font-medium">Email:</span> {user?.email}
                  </p>
                  <p className="text-sm text-gray-600">
                    <span className="font-medium">Role:</span> {user?.role}
                  </p>
                  <p className="text-sm text-gray-600">
                    <span className="font-medium">Status:</span>{' '}
                    <span className={user?.is_verified ? 'text-green-600' : 'text-yellow-600'}>
                      {user?.is_verified ? 'Verified' : 'Pending Verification'}
                    </span>
                  </p>
                  <p className="text-sm text-gray-600">
                    <span className="font-medium">Member since:</span>{' '}
                    {user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}
                  </p>
                </div>
              </div>
            </div>

            {/* Modules Card */}
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-5">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                  Available Modules
                </h3>
                <div className="mt-4 space-y-2">
                  <div className="flex items-center">
                    <span className="w-2 h-2 bg-green-400 rounded-full mr-2"></span>
                    <span className="text-sm text-gray-600">CRM</span>
                  </div>
                  <div className="flex items-center">
                    <span className="w-2 h-2 bg-green-400 rounded-full mr-2"></span>
                    <span className="text-sm text-gray-600">HRM</span>
                  </div>
                  <div className="flex items-center">
                    <span className="w-2 h-2 bg-green-400 rounded-full mr-2"></span>
                    <span className="text-sm text-gray-600">POS</span>
                  </div>
                  <div className="flex items-center">
                    <span className="w-2 h-2 bg-yellow-400 rounded-full mr-2"></span>
                    <span className="text-sm text-gray-600">LMS</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Quick Actions */}
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="p-5">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                  Quick Actions
                </h3>
                <div className="mt-4 space-y-2">
                  <button className="w-full text-left p-2 text-sm text-indigo-600 hover:bg-indigo-50 rounded">
                    View Analytics
                  </button>
                  <button className="w-full text-left p-2 text-sm text-indigo-600 hover:bg-indigo-50 rounded">
                    Manage Users
                  </button>
                  <button className="w-full text-left p-2 text-sm text-indigo-600 hover:bg-indigo-50 rounded">
                    Settings
                  </button>
                  <button className="w-full text-left p-2 text-sm text-indigo-600 hover:bg-indigo-50 rounded">
                    Support
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
