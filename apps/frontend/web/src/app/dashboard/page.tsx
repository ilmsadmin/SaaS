'use client'

import React from 'react'
import DashboardLayout from '@/components/ui/dashboard-layout'
import { StatCard, ActivityFeed, ChartCard, QuickActions } from '@/components/ui/dashboard-components'
import { useAuth } from '@/lib/auth-context'
import { useRequireAuth } from '@/lib/auth-context'
import { 
  Users, 
  Building, 
  TrendingUp, 
  Clock,
  Settings,
  BarChart3,
  FileText,
  Bell,
  Calendar
} from 'lucide-react'

export default function DashboardPage() {
  const { isAuthenticated, isLoading } = useRequireAuth()
  const { user } = useAuth()

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

  // Mock data - in real app, this would come from API
  const dashboardStats = {
    totalUsers: user?.role === 'admin' || user?.role === 'super_admin' ? 1247 : 47,
    activeModules: 5,
    tasksCompleted: 23,
    upcomingEvents: 3
  }

  const recentActivity = [
    { type: 'user_login', message: 'New user joined your organization', time: '15 minutes ago' },
    { type: 'module_access', message: 'CRM module updated with new features', time: '1 hour ago' },
    { type: 'user_login', message: 'Weekly report generated successfully', time: '2 hours ago' },
    { type: 'module_access', message: 'HRM attendance records updated', time: '4 hours ago' },
  ]

  const quickActions = [
    { label: 'View Analytics', onClick: () => console.log('Analytics'), icon: BarChart3 },
    { label: 'Generate Report', onClick: () => console.log('Report'), icon: FileText },
    { label: 'Notifications', onClick: () => console.log('Notifications'), icon: Bell },
    { label: 'Calendar', onClick: () => console.log('Calendar'), icon: Calendar },
    { label: 'Settings', onClick: () => console.log('Settings'), icon: Settings },
  ]

  const modules = [
    { name: 'CRM', status: 'Available', users: 12, description: 'Customer Relationship Management' },
    { name: 'HRM', status: 'Available', users: 8, description: 'Human Resource Management' },
    { name: 'POS', status: 'Available', users: 5, description: 'Point of Sale System' },
    { name: 'LMS', status: 'Available', users: 15, description: 'Learning Management System' },
    { name: 'Check-in', status: 'Available', users: 20, description: 'Attendance Tracking' },
    { name: 'Payment', status: 'Development', users: 0, description: 'Payment Processing' },
    { name: 'Accounting', status: 'Development', users: 0, description: 'Financial Management' },
    { name: 'E-commerce', status: 'Planned', users: 0, description: 'Online Store Platform' },
  ]

  return (
    <DashboardLayout 
      title="Dashboard Overview" 
      description="Welcome to your Zplus SaaS platform dashboard"
    >
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title={user?.role === 'admin' || user?.role === 'super_admin' ? "Total Users" : "Team Members"}
          value={dashboardStats.totalUsers}
          change={user?.role === 'admin' || user?.role === 'super_admin' ? "+8.2% this month" : "In your organization"}
          trend={user?.role === 'admin' || user?.role === 'super_admin' ? "up" : undefined}
          icon={Users}
          color="blue"
        />
        <StatCard
          title="Active Modules"
          value={`${dashboardStats.activeModules}/8`}
          change="Available to use"
          icon={Building}
          color="green"
        />
        <StatCard
          title="Tasks Completed"
          value={dashboardStats.tasksCompleted}
          change="+5 this week"
          trend="up"
          icon={TrendingUp}
          color="indigo"
        />
        <StatCard
          title="Upcoming Events"
          value={dashboardStats.upcomingEvents}
          change="This week"
          icon={Clock}
          color="yellow"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-8">
          {/* Modules Overview */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Available Modules</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {modules.map((module, index) => (
                  <div key={index} className="border rounded-lg p-4 hover:bg-gray-50 cursor-pointer">
                    <div className="flex items-center justify-between mb-2">
                      <h4 className="font-medium text-gray-900">{module.name}</h4>
                      <span className={`text-xs px-2 py-1 rounded-full ${
                        module.status === 'Available' ? 'bg-green-100 text-green-800' :
                        module.status === 'Development' ? 'bg-yellow-100 text-yellow-800' :
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {module.status}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 mb-2">{module.description}</p>
                    <p className="text-sm text-gray-500">
                      {module.users > 0 ? `${module.users} active users` : 'Not yet deployed'}
                    </p>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* User Profile Card */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Profile Information</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-3">
                  <div>
                    <label className="text-sm font-medium text-gray-600">Full Name</label>
                    <p className="text-gray-900">{user?.first_name} {user?.last_name}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-600">Email</label>
                    <p className="text-gray-900">{user?.email}</p>
                  </div>
                </div>
                <div className="space-y-3">
                  <div>
                    <label className="text-sm font-medium text-gray-600">Role</label>
                    <p className="text-gray-900">{user?.role}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-600">Status</label>
                    <p className={user?.is_verified ? 'text-green-600' : 'text-yellow-600'}>
                      {user?.is_verified ? 'Verified' : 'Pending Verification'}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Performance Overview */}
          <ChartCard title="Performance Overview">
            <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
              <div className="text-center">
                <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">Performance analytics will be displayed here</p>
                <p className="text-sm text-gray-400">Charts showing usage, productivity, and growth metrics</p>
              </div>
            </div>
          </ChartCard>
        </div>

        {/* Sidebar */}
        <div className="space-y-8">
          {/* Quick Actions */}
          <QuickActions actions={quickActions} />

          {/* Recent Activity */}
          <ActivityFeed activities={recentActivity} />
        </div>
      </div>
    </DashboardLayout>
  )
}
