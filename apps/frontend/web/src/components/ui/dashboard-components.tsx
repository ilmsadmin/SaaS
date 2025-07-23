'use client'

import React from 'react'
import { BarChart3, TrendingUp, Users, DollarSign, Activity, AlertCircle } from 'lucide-react'

interface StatCardProps {
  title: string
  value: string | number
  change?: string
  trend?: 'up' | 'down' | 'neutral'
  icon: React.ElementType
  color?: 'blue' | 'green' | 'yellow' | 'red' | 'indigo'
}

export function StatCard({ title, value, change, trend, icon: Icon, color = 'blue' }: StatCardProps) {
  const colorClasses = {
    blue: 'text-blue-600',
    green: 'text-green-600',
    yellow: 'text-yellow-600',
    red: 'text-red-600',
    indigo: 'text-indigo-600'
  }

  const trendClasses = {
    up: 'text-green-600',
    down: 'text-red-600',
    neutral: 'text-gray-600'
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-600">{title}</p>
          <p className="text-2xl font-bold text-gray-900">{value}</p>
          {change && (
            <p className={`text-sm ${trend ? trendClasses[trend] : 'text-gray-600'} flex items-center mt-1`}>
              {trend === 'up' && <TrendingUp className="h-4 w-4 mr-1" />}
              {trend === 'down' && <TrendingUp className="h-4 w-4 mr-1 rotate-180" />}
              {change}
            </p>
          )}
        </div>
        <Icon className={`h-8 w-8 ${colorClasses[color]}`} />
      </div>
    </div>
  )
}

interface ActivityItemProps {
  type: string
  message: string
  time: string
}

export function ActivityFeed({ activities }: { activities: ActivityItemProps[] }) {
  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'user_login':
        return <Users className="h-5 w-5 text-green-600" />
      case 'module_access':
        return <Activity className="h-5 w-5 text-blue-600" />
      case 'error':
        return <AlertCircle className="h-5 w-5 text-red-600" />
      default:
        return <Activity className="h-5 w-5 text-gray-600" />
    }
  }

  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Activity</h3>
        <div className="space-y-4">
          {activities.map((activity, index) => (
            <div key={index} className="flex items-start space-x-3">
              <div className="flex-shrink-0 mt-0.5">
                {getActivityIcon(activity.type)}
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-sm text-gray-900">{activity.message}</p>
                <p className="text-xs text-gray-500">{activity.time}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

interface ChartCardProps {
  title: string
  children: React.ReactNode
}

export function ChartCard({ title, children }: ChartCardProps) {
  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">{title}</h3>
        {children}
      </div>
    </div>
  )
}

export function QuickActions({ actions }: { actions: Array<{ label: string; onClick: () => void; icon: React.ElementType }> }) {
  return (
    <div className="bg-white rounded-lg shadow">
      <div className="p-6">
        <h3 className="text-lg font-medium text-gray-900 mb-4">Quick Actions</h3>
        <div className="space-y-2">
          {actions.map((action, index) => (
            <button
              key={index}
              onClick={action.onClick}
              className="w-full flex items-center p-3 text-sm text-left text-indigo-600 hover:bg-indigo-50 rounded-md transition-colors"
            >
              <action.icon className="h-4 w-4 mr-3" />
              {action.label}
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}
