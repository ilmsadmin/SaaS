'use client'

import { useAdmin } from '@/contexts/AdminContext'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import DashboardLayout from '@/components/DashboardLayout'
import StatsCard from '@/components/StatsCard'
import { Users, Building2, CreditCard, TrendingUp } from 'lucide-react'
import RecentActivity from '@/components/RecentActivity'
import SystemHealth from '@/components/SystemHealth'

export default function DashboardPage() {
  const { isAuthenticated, isLoading, user } = useAdmin()
  const router = useRouter()
  const [stats, setStats] = useState({
    totalTenants: 0,
    totalUsers: 0,
    totalRevenue: 0,
    activeSubscriptions: 0
  })

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, isLoading, router])

  useEffect(() => {
    // Fetch dashboard stats
    const fetchStats = async () => {
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/stats`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('admin_token')}`
          }
        })
        if (response.ok) {
          const data = await response.json()
          setStats(data)
        }
      } catch (error) {
        console.error('Failed to fetch stats:', error)
      }
    }

    if (isAuthenticated) {
      fetchStats()
    }
  }, [isAuthenticated])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return null
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="bg-white shadow">
          <div className="px-4 py-6 sm:px-6 lg:px-8">
            <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
            <p className="mt-2 text-sm text-gray-700">
              Welcome back, {user?.name}. Here's what's happening with your platform.
            </p>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard
            title="Total Tenants"
            value={stats.totalTenants.toString()}
            icon={Building2}
            change="+12%"
            changeType="increase"
          />
          <StatsCard
            title="Total Users"
            value={stats.totalUsers.toString()}
            icon={Users}
            change="+8%"
            changeType="increase"
          />
          <StatsCard
            title="Monthly Revenue"
            value={`$${stats.totalRevenue.toLocaleString()}`}
            icon={CreditCard}
            change="+23%"
            changeType="increase"
          />
          <StatsCard
            title="Active Subscriptions"
            value={stats.activeSubscriptions.toString()}
            icon={TrendingUp}
            change="+5%"
            changeType="increase"
          />
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          {/* Recent Activity */}
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
              <RecentActivity />
            </div>
          </div>

          {/* System Health */}
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">System Health</h3>
              <SystemHealth />
            </div>
          </div>
        </div>
      </div>
    </DashboardLayout>
  )
}
