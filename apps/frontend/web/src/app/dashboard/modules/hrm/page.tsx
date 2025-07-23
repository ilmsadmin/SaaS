'use client'

import React from 'react'
import DashboardLayout from '@/components/ui/dashboard-layout'
import { StatCard, ActivityFeed, ChartCard, QuickActions } from '@/components/ui/dashboard-components'
import { 
  Users, 
  Clock, 
  Calendar, 
  Award, 
  UserPlus,
  FileText,
  TrendingUp,
  AlertCircle,
  CheckCircle,
  XCircle
} from 'lucide-react'

export default function HRMDashboard() {
  // Mock HRM data
  const hrmStats = {
    totalEmployees: 142,
    presentToday: 128,
    onLeave: 8,
    pendingRequests: 12
  }

  const recentActivity = [
    { type: 'user_login', message: 'John Doe submitted leave request for Dec 25-26', time: '30 minutes ago' },
    { type: 'module_access', message: 'Sarah Wilson completed performance review', time: '1 hour ago' },
    { type: 'user_login', message: 'New employee Michael Johnson onboarded', time: '2 hours ago' },
    { type: 'module_access', message: 'Team meeting scheduled for Dec 24', time: '3 hours ago' },
  ]

  const quickActions = [
    { label: 'Add Employee', onClick: () => console.log('Add employee'), icon: UserPlus },
    { label: 'Approve Requests', onClick: () => console.log('Approve requests'), icon: CheckCircle },
    { label: 'Schedule Meeting', onClick: () => console.log('Schedule meeting'), icon: Calendar },
    { label: 'Generate Report', onClick: () => console.log('Generate report'), icon: FileText },
    { label: 'Performance Review', onClick: () => console.log('Performance review'), icon: Award },
  ]

  const leaveRequests = [
    { employee: 'John Doe', type: 'Annual Leave', dates: 'Dec 25-26', status: 'pending', submitted: '2 hours ago' },
    { employee: 'Sarah Wilson', type: 'Sick Leave', dates: 'Dec 23', status: 'approved', submitted: '1 day ago' },
    { employee: 'Mike Johnson', type: 'Personal', dates: 'Dec 30', status: 'pending', submitted: '3 hours ago' },
    { employee: 'Lisa Brown', type: 'Annual Leave', dates: 'Jan 2-5', status: 'rejected', submitted: '2 days ago' },
  ]

  const upcomingEvents = [
    { title: 'Team Meeting', date: 'Today, 2:00 PM', type: 'meeting', attendees: 12 },
    { title: 'Performance Reviews Due', date: 'Dec 24, 5:00 PM', type: 'deadline', attendees: 8 },
    { title: 'Holiday Party', date: 'Dec 25, 6:00 PM', type: 'event', attendees: 142 },
    { title: 'Q4 Planning Meeting', date: 'Dec 26, 10:00 AM', type: 'meeting', attendees: 6 },
  ]

  const departmentStats = [
    { name: 'Engineering', employees: 45, present: 42, onLeave: 2 },
    { name: 'Sales', employees: 28, present: 26, onLeave: 1 },
    { name: 'Marketing', employees: 22, present: 20, onLeave: 2 },
    { name: 'HR', employees: 8, present: 8, onLeave: 0 },
    { name: 'Finance', employees: 12, present: 11, onLeave: 1 },
    { name: 'Operations', employees: 15, present: 14, onLeave: 1 },
  ]

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'approved': return 'bg-green-100 text-green-800'
      case 'pending': return 'bg-yellow-100 text-yellow-800'
      case 'rejected': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'meeting': return <Calendar className="h-4 w-4 text-blue-600" />
      case 'deadline': return <AlertCircle className="h-4 w-4 text-red-600" />
      case 'event': return <Award className="h-4 w-4 text-green-600" />
      default: return <Calendar className="h-4 w-4 text-gray-600" />
    }
  }

  return (
    <DashboardLayout 
      title="HRM Dashboard" 
      description="Manage employees, attendance, and HR processes"
    >
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title="Total Employees"
          value={hrmStats.totalEmployees}
          change="+3 this month"
          trend="up"
          icon={Users}
          color="blue"
        />
        <StatCard
          title="Present Today"
          value={hrmStats.presentToday}
          change="90.1% attendance"
          icon={CheckCircle}
          color="green"
        />
        <StatCard
          title="On Leave"
          value={hrmStats.onLeave}
          change="5.6% of workforce"
          icon={XCircle}
          color="yellow"
        />
        <StatCard
          title="Pending Requests"
          value={hrmStats.pendingRequests}
          change="Require action"
          icon={AlertCircle}
          color="red"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-8">
          {/* Department Overview */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Department Overview</h3>
              <div className="space-y-4">
                {departmentStats.map((dept, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                    <div>
                      <h4 className="font-medium text-gray-900">{dept.name}</h4>
                      <p className="text-sm text-gray-600">
                        {dept.employees} employees • {dept.present} present • {dept.onLeave} on leave
                      </p>
                    </div>
                    <div className="flex items-center space-x-4">
                      <div className="text-right">
                        <p className="text-sm font-medium text-green-600">
                          {Math.round((dept.present / dept.employees) * 100)}%
                        </p>
                        <p className="text-xs text-gray-500">Attendance</p>
                      </div>
                      <button className="text-indigo-600 hover:text-indigo-900 text-sm">
                        View
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Leave Requests */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Leave Requests</h3>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Employee
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Type
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Dates
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Action
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {leaveRequests.map((request, index) => (
                      <tr key={index} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="font-medium text-gray-900">{request.employee}</div>
                          <div className="text-sm text-gray-500">Submitted {request.submitted}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {request.type}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {request.dates}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 py-1 text-xs rounded-full ${getStatusColor(request.status)}`}>
                            {request.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          {request.status === 'pending' && (
                            <div className="flex space-x-2">
                              <button className="text-green-600 hover:text-green-900">Approve</button>
                              <button className="text-red-600 hover:text-red-900">Reject</button>
                            </div>
                          )}
                          {request.status !== 'pending' && (
                            <button className="text-indigo-600 hover:text-indigo-900">View</button>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          {/* Attendance Trends */}
          <ChartCard title="Attendance Trends">
            <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
              <div className="text-center">
                <TrendingUp className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">Attendance trends chart will be implemented here</p>
                <p className="text-sm text-gray-400">Shows daily/weekly/monthly attendance patterns</p>
              </div>
            </div>
          </ChartCard>
        </div>

        {/* Sidebar */}
        <div className="space-y-8">
          {/* Quick Actions */}
          <QuickActions actions={quickActions} />

          {/* Upcoming Events */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Upcoming Events</h3>
            <div className="space-y-3">
              {upcomingEvents.map((event, index) => (
                <div key={index} className="flex items-start space-x-3 p-3 border rounded-lg">
                  <div className="flex-shrink-0 mt-0.5">
                    {getEventIcon(event.type)}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium text-gray-900">{event.title}</p>
                    <p className="text-xs text-gray-600">{event.date}</p>
                    <p className="text-xs text-gray-500">{event.attendees} attendees</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Quick Stats */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Quick Stats</h3>
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Average Work Hours</span>
                <span className="text-sm font-medium">8.2h/day</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Employee Satisfaction</span>
                <span className="text-sm font-medium">87%</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Retention Rate</span>
                <span className="text-sm font-medium">94.5%</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Training Completion</span>
                <span className="text-sm font-medium">76%</span>
              </div>
            </div>
          </div>

          {/* Recent Activity */}
          <ActivityFeed activities={recentActivity} />
        </div>
      </div>
    </DashboardLayout>
  )
}
