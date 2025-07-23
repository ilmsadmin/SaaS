'use client'

import React from 'react'
import DashboardLayout from '@/components/ui/dashboard-layout'
import { StatCard, ActivityFeed, ChartCard, QuickActions } from '@/components/ui/dashboard-components'
import { useAuth } from '@/lib/auth-context'
import { 
  User, 
  Clock, 
  BookOpen, 
  Award, 
  Calendar,
  Settings, 
  Bell, 
  FileText,
  Target,
  TrendingUp
} from 'lucide-react'

export default function UserDashboard() {
  const { user } = useAuth()

  // Mock data - in real app, this would come from API
  const userStats = {
    tasksCompleted: 23,
    hoursWorked: 38.5,
    coursesCompleted: 5,
    achievements: 12
  }

  const recentActivity = [
    { type: 'user_login', message: 'Completed "Advanced Sales Techniques" course', time: '2 hours ago' },
    { type: 'module_access', message: 'Updated customer profile in CRM', time: '4 hours ago' },
    { type: 'user_login', message: 'Clocked in for the day', time: '8 hours ago' },
    { type: 'module_access', message: 'Processed 5 POS transactions', time: '1 day ago' },
  ]

  const quickActions = [
    { label: 'Update Profile', onClick: () => console.log('Update profile'), icon: User },
    { label: 'Notifications', onClick: () => console.log('Notifications'), icon: Bell },
    { label: 'My Tasks', onClick: () => console.log('Tasks'), icon: FileText },
    { label: 'Time Tracking', onClick: () => console.log('Time tracking'), icon: Clock },
    { label: 'Settings', onClick: () => console.log('Settings'), icon: Settings },
  ]

  const modules = [
    { name: 'CRM', lastAccessed: '2 hours ago', usage: '85%', status: 'active' },
    { name: 'HRM', lastAccessed: '1 day ago', usage: '45%', status: 'active' },
    { name: 'LMS', lastAccessed: '2 hours ago', usage: '92%', status: 'active' },
    { name: 'POS', lastAccessed: '4 hours ago', usage: '67%', status: 'active' },
    { name: 'Check-in', lastAccessed: '8 hours ago', usage: '100%', status: 'active' },
  ]

  const upcomingTasks = [
    { title: 'Complete Q4 Sales Report', due: 'Today, 5:00 PM', priority: 'high' },
    { title: 'Customer Follow-up Call', due: 'Tomorrow, 10:00 AM', priority: 'medium' },
    { title: 'Team Meeting Preparation', due: 'Dec 24, 2:00 PM', priority: 'low' },
    { title: 'Training Module Review', due: 'Dec 25, 9:00 AM', priority: 'medium' },
  ]

  const achievements = [
    { title: 'Sales Champion', description: 'Exceeded monthly sales target', date: 'Dec 20, 2024' },
    { title: 'Learning Enthusiast', description: 'Completed 5 courses this month', date: 'Dec 18, 2024' },
    { title: 'Team Player', description: 'Helped 10 colleagues this week', date: 'Dec 15, 2024' },
  ]

  return (
    <DashboardLayout 
      title={`Welcome back, ${user?.first_name}!`}
      description="Track your progress, manage tasks, and access your modules"
    >
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title="Tasks Completed"
          value={userStats.tasksCompleted}
          change="+5 this week"
          trend="up"
          icon={Target}
          color="green"
        />
        <StatCard
          title="Hours Worked"
          value={`${userStats.hoursWorked}h`}
          change="This week"
          icon={Clock}
          color="blue"
        />
        <StatCard
          title="Courses Completed"
          value={userStats.coursesCompleted}
          change="+2 this month"
          trend="up"
          icon={BookOpen}
          color="indigo"
        />
        <StatCard
          title="Achievements"
          value={userStats.achievements}
          change="+3 this month"
          trend="up"
          icon={Award}
          color="yellow"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-8">
          {/* My Modules */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">My Modules</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {modules.map((module, index) => (
                  <div key={index} className="border rounded-lg p-4 hover:bg-gray-50 cursor-pointer">
                    <div className="flex items-center justify-between mb-2">
                      <h4 className="font-medium text-gray-900">{module.name}</h4>
                      <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full">
                        {module.status}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 mb-2">Last accessed: {module.lastAccessed}</p>
                    <div className="flex items-center">
                      <div className="flex-1 bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-indigo-600 h-2 rounded-full" 
                          style={{ width: module.usage }}
                        />
                      </div>
                      <span className="ml-2 text-xs text-gray-600">{module.usage}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Tasks */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Upcoming Tasks</h3>
              <div className="space-y-3">
                {upcomingTasks.map((task, index) => (
                  <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center space-x-3">
                      <div className={`w-3 h-3 rounded-full ${
                        task.priority === 'high' ? 'bg-red-400' :
                        task.priority === 'medium' ? 'bg-yellow-400' :
                        'bg-green-400'
                      }`} />
                      <div>
                        <p className="font-medium text-gray-900">{task.title}</p>
                        <p className="text-sm text-gray-500">{task.due}</p>
                      </div>
                    </div>
                    <button className="text-indigo-600 hover:text-indigo-900 text-sm">
                      View
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Performance Chart */}
          <ChartCard title="Performance Overview">
            <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
              <div className="text-center">
                <TrendingUp className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">Performance chart will be implemented here</p>
                <p className="text-sm text-gray-400">Shows productivity trends and metrics</p>
              </div>
            </div>
          </ChartCard>
        </div>

        {/* Sidebar */}
        <div className="space-y-8">
          {/* Profile Summary */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Profile</h3>
            <div className="text-center">
              <div className="w-20 h-20 bg-indigo-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <User className="h-10 w-10 text-indigo-600" />
              </div>
              <h4 className="font-medium text-gray-900">{user?.first_name} {user?.last_name}</h4>
              <p className="text-sm text-gray-600">{user?.role}</p>
              <p className="text-sm text-gray-600">{user?.email}</p>
              <div className="mt-4 pt-4 border-t">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Completion Rate</span>
                  <span className="font-medium">87%</span>
                </div>
                <div className="mt-1 bg-gray-200 rounded-full h-2">
                  <div className="bg-green-600 h-2 rounded-full" style={{ width: '87%' }} />
                </div>
              </div>
            </div>
          </div>

          {/* Quick Actions */}
          <QuickActions actions={quickActions} />

          {/* Recent Achievements */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Achievements</h3>
            <div className="space-y-3">
              {achievements.map((achievement, index) => (
                <div key={index} className="flex items-start space-x-3">
                  <Award className="h-5 w-5 text-yellow-500 mt-0.5" />
                  <div>
                    <p className="font-medium text-gray-900">{achievement.title}</p>
                    <p className="text-sm text-gray-600">{achievement.description}</p>
                    <p className="text-xs text-gray-500">{achievement.date}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Recent Activity */}
          <ActivityFeed activities={recentActivity} />
        </div>
      </div>
    </DashboardLayout>
  )
}
