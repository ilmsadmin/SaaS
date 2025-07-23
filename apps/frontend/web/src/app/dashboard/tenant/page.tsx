'use client'

import React from 'react'
import DashboardLayout from '@/components/ui/dashboard-layout'
import { StatCard, ActivityFeed, ChartCard, QuickActions } from '@/components/ui/dashboard-components'
import { 
  Building, 
  Users, 
  TrendingUp, 
  DollarSign, 
  Settings, 
  Plus, 
  BarChart3,
  CreditCard,
  Shield
} from 'lucide-react'

export default function TenantDashboard() {
  // Mock data - in real app, this would come from API
  const tenantStats = {
    totalUsers: 47,
    activeModules: 5,
    monthlyUsage: 89.5,
    storageUsed: 2.3,
    subscription: 'Enterprise'
  }

  const recentActivity = [
    { type: 'user_login', message: 'John Doe logged in to CRM module', time: '2 minutes ago' },
    { type: 'module_access', message: 'New employee accessed HRM module', time: '15 minutes ago' },
    { type: 'user_login', message: 'Sarah Wilson completed LMS course', time: '1 hour ago' },
    { type: 'module_access', message: 'POS transaction processed', time: '2 hours ago' },
  ]

  const quickActions = [
    { label: 'Add New User', onClick: () => console.log('Add user'), icon: Plus },
    { label: 'Manage Modules', onClick: () => console.log('Manage modules'), icon: Settings },
    { label: 'View Analytics', onClick: () => console.log('View analytics'), icon: BarChart3 },
    { label: 'Billing Settings', onClick: () => console.log('Billing'), icon: CreditCard },
    { label: 'Security Settings', onClick: () => console.log('Security'), icon: Shield },
  ]

  const modules = [
    { name: 'CRM', status: 'active', users: 12, lastUsed: '2 minutes ago' },
    { name: 'HRM', status: 'active', users: 8, lastUsed: '15 minutes ago' },
    { name: 'POS', status: 'active', users: 5, lastUsed: '2 hours ago' },
    { name: 'LMS', status: 'active', users: 15, lastUsed: '1 hour ago' },
    { name: 'Check-in', status: 'active', users: 20, lastUsed: '30 minutes ago' },
    { name: 'Payment', status: 'inactive', users: 0, lastUsed: 'Never' },
  ]

  return (
    <DashboardLayout 
      title="Tenant Dashboard" 
      description="Manage your organization's modules, users, and settings"
    >
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title="Total Users"
          value={tenantStats.totalUsers}
          change="+12% from last month"
          trend="up"
          icon={Users}
          color="blue"
        />
        <StatCard
          title="Active Modules"
          value={`${tenantStats.activeModules}/8`}
          icon={Building}
          color="green"
        />
        <StatCard
          title="Monthly Usage"
          value={`${tenantStats.monthlyUsage}%`}
          change="+5% from last month"
          trend="up"
          icon={TrendingUp}
          color="indigo"
        />
        <StatCard
          title="Storage Used"
          value={`${tenantStats.storageUsed} GB`}
          change="of 10 GB limit"
          icon={DollarSign}
          color="yellow"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Modules Overview */}
        <div className="lg:col-span-2">
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Module Status</h3>
              <div className="space-y-4">
                {modules.map((module, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center space-x-4">
                      <div className={`w-3 h-3 rounded-full ${
                        module.status === 'active' ? 'bg-green-400' : 'bg-gray-300'
                      }`} />
                      <div>
                        <p className="font-medium text-gray-900">{module.name}</p>
                        <p className="text-sm text-gray-500">
                          {module.users} users • Last used: {module.lastUsed}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        module.status === 'active' 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-gray-100 text-gray-800'
                      }`}>
                        {module.status}
                      </span>
                      <button className="text-indigo-600 hover:text-indigo-900 text-sm">
                        Configure
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Usage Chart */}
          <div className="mt-8">
            <ChartCard title="Usage Analytics">
              <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
                <div className="text-center">
                  <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                  <p className="text-gray-500">Usage chart will be implemented here</p>
                  <p className="text-sm text-gray-400">Connect to analytics service</p>
                </div>
              </div>
            </ChartCard>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-8">
          {/* Subscription Info */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Subscription</h3>
            <div className="space-y-3">
              <div>
                <p className="text-sm text-gray-600">Current Plan</p>
                <p className="text-lg font-semibold text-indigo-600">{tenantStats.subscription}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Next Billing</p>
                <p className="text-sm font-medium">January 15, 2025</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Monthly Cost</p>
                <p className="text-lg font-semibold">$299/month</p>
              </div>
              <button className="w-full bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 transition-colors">
                Upgrade Plan
              </button>
            </div>
          </div>

          {/* Quick Actions */}
          <QuickActions actions={quickActions} />

          {/* Recent Activity */}
          <ActivityFeed activities={recentActivity} />
        </div>
      </div>
    </DashboardLayout>
  )
}
