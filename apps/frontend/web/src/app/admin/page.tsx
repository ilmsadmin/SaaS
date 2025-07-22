'use client'

import React from 'react'
import { Building, Users, CreditCard, TrendingUp, Activity, DollarSign } from 'lucide-react'

export default function AdminDashboard() {
  // Mock data for development
  const stats = {
    totalTenants: 156,
    activeTenants: 142,
    totalUsers: 2847,
    monthlyRevenue: 45670,
    growth: 12.5
  }

  const recentActivity = [
    { id: 1, type: 'tenant_created', message: 'New tenant "TechCorp" created', time: '2 hours ago' },
    { id: 2, type: 'subscription_updated', message: 'Acme Corp upgraded to Enterprise plan', time: '4 hours ago' },
    { id: 3, type: 'tenant_suspended', message: 'Tenant "OldCompany" suspended for non-payment', time: '1 day ago' },
    { id: 4, type: 'user_registered', message: '45 new users registered today', time: '1 day ago' },
  ]

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
          <p className="text-gray-600 mt-2">
            Overview of your SaaS platform
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Tenants</p>
                <p className="text-2xl font-bold text-gray-900">{stats.totalTenants}</p>
              </div>
              <Building className="h-8 w-8 text-blue-600" />
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Active Tenants</p>
                <p className="text-2xl font-bold text-green-600">{stats.activeTenants}</p>
              </div>
              <Activity className="h-8 w-8 text-green-600" />
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Users</p>
                <p className="text-2xl font-bold text-gray-900">{stats.totalUsers.toLocaleString()}</p>
              </div>
              <Users className="h-8 w-8 text-purple-600" />
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Monthly Revenue</p>
                <p className="text-2xl font-bold text-gray-900">${stats.monthlyRevenue.toLocaleString()}</p>
              </div>
              <DollarSign className="h-8 w-8 text-yellow-600" />
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Growth Rate</p>
                <p className="text-2xl font-bold text-green-600">+{stats.growth}%</p>
              </div>
              <TrendingUp className="h-8 w-8 text-green-600" />
            </div>
          </div>
        </div>

        {/* Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Recent Activity */}
          <div className="bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200">
              <h2 className="text-lg font-semibold text-gray-900">Recent Activity</h2>
            </div>
            <div className="p-6">
              <div className="space-y-4">
                {recentActivity.map((activity) => (
                  <div key={activity.id} className="flex items-start space-x-3">
                    <div className="flex-shrink-0">
                      <div className="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center">
                        <Activity className="h-4 w-4 text-blue-600" />
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm text-gray-900">{activity.message}</p>
                      <p className="text-xs text-gray-500">{activity.time}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Quick Actions */}
          <div className="bg-white rounded-lg shadow">
            <div className="px-6 py-4 border-b border-gray-200">
              <h2 className="text-lg font-semibold text-gray-900">Quick Actions</h2>
            </div>
            <div className="p-6">
              <div className="grid grid-cols-2 gap-4">
                <a
                  href="/admin/tenants"
                  className="flex items-center p-4 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors"
                >
                  <Building className="h-6 w-6 text-blue-600 mr-3" />
                  <span className="text-sm font-medium text-blue-900">Manage Tenants</span>
                </a>
                <a
                  href="/admin/users"
                  className="flex items-center p-4 bg-green-50 rounded-lg hover:bg-green-100 transition-colors"
                >
                  <Users className="h-6 w-6 text-green-600 mr-3" />
                  <span className="text-sm font-medium text-green-900">Manage Users</span>
                </a>
                <a
                  href="/admin/plans"
                  className="flex items-center p-4 bg-purple-50 rounded-lg hover:bg-purple-100 transition-colors"
                >
                  <CreditCard className="h-6 w-6 text-purple-600 mr-3" />
                  <span className="text-sm font-medium text-purple-900">Manage Plans</span>
                </a>
                <a
                  href="/admin/analytics"
                  className="flex items-center p-4 bg-yellow-50 rounded-lg hover:bg-yellow-100 transition-colors"
                >
                  <TrendingUp className="h-6 w-6 text-yellow-600 mr-3" />
                  <span className="text-sm font-medium text-yellow-900">View Analytics</span>
                </a>
              </div>
            </div>
          </div>
        </div>

        {/* System Status */}
        <div className="mt-8 bg-white rounded-lg shadow">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">System Status</h2>
          </div>
          <div className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-green-400 rounded-full"></div>
                <span className="text-sm text-gray-900">API Gateway</span>
                <span className="text-xs text-green-600">Healthy</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-green-400 rounded-full"></div>
                <span className="text-sm text-gray-900">Database</span>
                <span className="text-xs text-green-600">Healthy</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-green-400 rounded-full"></div>
                <span className="text-sm text-gray-900">Auth Service</span>
                <span className="text-xs text-green-600">Healthy</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-yellow-400 rounded-full"></div>
                <span className="text-sm text-gray-900">Tenant Service</span>
                <span className="text-xs text-yellow-600">Starting</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-green-400 rounded-full"></div>
                <span className="text-sm text-gray-900">Redis</span>
                <span className="text-xs text-green-600">Healthy</span>
              </div>
              <div className="flex items-center space-x-3">
                <div className="h-3 w-3 bg-green-400 rounded-full"></div>
                <span className="text-sm text-gray-900">File Storage</span>
                <span className="text-xs text-green-600">Healthy</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
