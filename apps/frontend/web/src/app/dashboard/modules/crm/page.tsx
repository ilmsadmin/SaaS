'use client'

import React from 'react'
import DashboardLayout from '@/components/ui/dashboard-layout'
import { StatCard, ActivityFeed, ChartCard, QuickActions } from '@/components/ui/dashboard-components'
import { 
  Users, 
  Phone, 
  Mail, 
  DollarSign, 
  TrendingUp,
  Plus,
  Filter,
  Search,
  Calendar,
  FileText
} from 'lucide-react'

export default function CRMDashboard() {
  // Mock CRM data
  const crmStats = {
    totalCustomers: 1247,
    activeLeads: 89,
    dealsWon: 23,
    revenue: 145670
  }

  const recentActivity = [
    { type: 'user_login', message: 'New lead "TechCorp" added by John Doe', time: '15 minutes ago' },
    { type: 'module_access', message: 'Deal "Enterprise License" closed - $25,000', time: '1 hour ago' },
    { type: 'user_login', message: 'Follow-up call scheduled with ABC Company', time: '2 hours ago' },
    { type: 'module_access', message: 'Customer "XYZ Ltd" updated contact information', time: '3 hours ago' },
  ]

  const quickActions = [
    { label: 'Add New Lead', onClick: () => console.log('Add lead'), icon: Plus },
    { label: 'Schedule Call', onClick: () => console.log('Schedule call'), icon: Phone },
    { label: 'Send Email', onClick: () => console.log('Send email'), icon: Mail },
    { label: 'Create Report', onClick: () => console.log('Create report'), icon: FileText },
    { label: 'View Calendar', onClick: () => console.log('View calendar'), icon: Calendar },
  ]

  const recentLeads = [
    { name: 'TechCorp Solutions', status: 'new', value: '$15,000', contact: 'John Smith', date: '2024-12-22' },
    { name: 'Global Industries', status: 'qualified', value: '$35,000', contact: 'Sarah Johnson', date: '2024-12-21' },
    { name: 'Digital Dynamics', status: 'proposal', value: '$22,000', contact: 'Mike Wilson', date: '2024-12-20' },
    { name: 'Innovation Labs', status: 'negotiation', value: '$45,000', contact: 'Lisa Brown', date: '2024-12-19' },
  ]

  const topCustomers = [
    { name: 'Enterprise Corp', revenue: '$125,000', deals: 8, lastContact: '2 days ago' },
    { name: 'MegaTech Ltd', revenue: '$89,000', deals: 5, lastContact: '1 week ago' },
    { name: 'StartupXYZ', revenue: '$67,000', deals: 12, lastContact: '3 days ago' },
    { name: 'Global Solutions', revenue: '$54,000', deals: 6, lastContact: '5 days ago' },
  ]

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'new': return 'bg-blue-100 text-blue-800'
      case 'qualified': return 'bg-green-100 text-green-800'
      case 'proposal': return 'bg-yellow-100 text-yellow-800'
      case 'negotiation': return 'bg-orange-100 text-orange-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  return (
    <DashboardLayout 
      title="CRM Dashboard" 
      description="Manage customers, leads, and sales pipeline"
    >
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title="Total Customers"
          value={crmStats.totalCustomers.toLocaleString()}
          change="+8.2% from last month"
          trend="up"
          icon={Users}
          color="blue"
        />
        <StatCard
          title="Active Leads"
          value={crmStats.activeLeads}
          change="+12 this week"
          trend="up"
          icon={TrendingUp}
          color="green"
        />
        <StatCard
          title="Deals Won"
          value={crmStats.dealsWon}
          change="This month"
          icon={FileText}
          color="indigo"
        />
        <StatCard
          title="Revenue"
          value={`$${crmStats.revenue.toLocaleString()}`}
          change="+15.3% from last month"
          trend="up"
          icon={DollarSign}
          color="yellow"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-8">
          {/* Recent Leads */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Recent Leads</h3>
                <div className="flex space-x-2">
                  <button className="flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm">
                    <Filter className="h-4 w-4 mr-2" />
                    Filter
                  </button>
                  <button className="flex items-center px-3 py-2 border border-gray-300 rounded-md text-sm">
                    <Search className="h-4 w-4 mr-2" />
                    Search
                  </button>
                </div>
              </div>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Company
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Value
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Contact
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Date
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {recentLeads.map((lead, index) => (
                      <tr key={index} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="font-medium text-gray-900">{lead.name}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 py-1 text-xs rounded-full ${getStatusColor(lead.status)}`}>
                            {lead.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {lead.value}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {lead.contact}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {lead.date}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          {/* Sales Pipeline Chart */}
          <ChartCard title="Sales Pipeline">
            <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
              <div className="text-center">
                <TrendingUp className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">Sales pipeline chart will be implemented here</p>
                <p className="text-sm text-gray-400">Shows lead progression through sales stages</p>
              </div>
            </div>
          </ChartCard>

          {/* Top Customers */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Top Customers</h3>
              <div className="space-y-4">
                {topCustomers.map((customer, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                    <div>
                      <h4 className="font-medium text-gray-900">{customer.name}</h4>
                      <p className="text-sm text-gray-600">
                        {customer.deals} deals • Last contact: {customer.lastContact}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-gray-900">{customer.revenue}</p>
                      <button className="text-indigo-600 hover:text-indigo-900 text-sm">
                        View Details
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-8">
          {/* Quick Actions */}
          <QuickActions actions={quickActions} />

          {/* Today's Schedule */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Today's Schedule</h3>
            <div className="space-y-3">
              <div className="flex items-center space-x-3 p-2 bg-blue-50 rounded">
                <Calendar className="h-4 w-4 text-blue-600" />
                <div>
                  <p className="text-sm font-medium">Call with TechCorp</p>
                  <p className="text-xs text-gray-600">2:00 PM - 2:30 PM</p>
                </div>
              </div>
              <div className="flex items-center space-x-3 p-2 bg-green-50 rounded">
                <Calendar className="h-4 w-4 text-green-600" />
                <div>
                  <p className="text-sm font-medium">Demo for StartupXYZ</p>
                  <p className="text-xs text-gray-600">4:00 PM - 5:00 PM</p>
                </div>
              </div>
              <div className="flex items-center space-x-3 p-2 bg-yellow-50 rounded">
                <Calendar className="h-4 w-4 text-yellow-600" />
                <div>
                  <p className="text-sm font-medium">Follow-up emails</p>
                  <p className="text-xs text-gray-600">5:30 PM - 6:00 PM</p>
                </div>
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
