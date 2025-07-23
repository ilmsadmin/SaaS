'use client'

import React from 'react'
import DashboardLayout from '@/components/ui/dashboard-layout'
import { StatCard, ActivityFeed, ChartCard, QuickActions } from '@/components/ui/dashboard-components'
import { 
  ShoppingCart, 
  DollarSign, 
  Package, 
  TrendingUp,
  Plus,
  Scan,
  Receipt,
  BarChart3,
  AlertTriangle,
  CheckCircle
} from 'lucide-react'

export default function POSDashboard() {
  // Mock POS data
  const posStats = {
    todaySales: 12450,
    transactionsToday: 89,
    averageOrder: 139.89,
    lowStockItems: 12
  }

  const recentActivity = [
    { type: 'module_access', message: 'Transaction #TX-2024-1245 completed - $89.50', time: '2 minutes ago' },
    { type: 'user_login', message: 'Product "Wireless Headphones" added to inventory', time: '15 minutes ago' },
    { type: 'module_access', message: 'Refund processed for transaction #TX-2024-1240', time: '30 minutes ago' },
    { type: 'user_login', message: 'Daily sales report generated', time: '1 hour ago' },
  ]

  const quickActions = [
    { label: 'New Sale', onClick: () => console.log('New sale'), icon: Plus },
    { label: 'Scan Product', onClick: () => console.log('Scan product'), icon: Scan },
    { label: 'View Receipt', onClick: () => console.log('View receipt'), icon: Receipt },
    { label: 'Add Product', onClick: () => console.log('Add product'), icon: Package },
    { label: 'Sales Report', onClick: () => console.log('Sales report'), icon: BarChart3 },
  ]

  const recentTransactions = [
    { id: 'TX-2024-1245', amount: 89.50, items: 3, customer: 'John Doe', time: '2 mins ago', status: 'completed' },
    { id: 'TX-2024-1244', amount: 245.00, items: 5, customer: 'Sarah Wilson', time: '8 mins ago', status: 'completed' },
    { id: 'TX-2024-1243', amount: 67.25, items: 2, customer: 'Mike Johnson', time: '12 mins ago', status: 'completed' },
    { id: 'TX-2024-1242', amount: 156.75, items: 4, customer: 'Lisa Brown', time: '18 mins ago', status: 'refunded' },
  ]

  const topProducts = [
    { name: 'Wireless Headphones', sales: 45, revenue: '$2,250', stock: 23, trend: 'up' },
    { name: 'Smartphone Case', sales: 67, revenue: '$1,340', stock: 156, trend: 'up' },
    { name: 'USB Cable', sales: 89, revenue: '$890', stock: 89, trend: 'neutral' },
    { name: 'Power Bank', sales: 34, revenue: '$1,700', stock: 12, trend: 'down' },
  ]

  const lowStockProducts = [
    { name: 'Wireless Mouse', currentStock: 5, minStock: 20, category: 'Electronics' },
    { name: 'Notebook', currentStock: 8, minStock: 50, category: 'Stationery' },
    { name: 'Phone Charger', currentStock: 3, minStock: 15, category: 'Accessories' },
    { name: 'Bluetooth Speaker', currentStock: 2, minStock: 10, category: 'Electronics' },
  ]

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed': return 'bg-green-100 text-green-800'
      case 'pending': return 'bg-yellow-100 text-yellow-800'
      case 'refunded': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return <TrendingUp className="h-4 w-4 text-green-600" />
      case 'down': return <TrendingUp className="h-4 w-4 text-red-600 rotate-180" />
      default: return <div className="h-4 w-4" />
    }
  }

  return (
    <DashboardLayout 
      title="POS Dashboard" 
      description="Manage sales, inventory, and transactions"
    >
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title="Today's Sales"
          value={`$${posStats.todaySales.toLocaleString()}`}
          change="+18.2% from yesterday"
          trend="up"
          icon={DollarSign}
          color="green"
        />
        <StatCard
          title="Transactions"
          value={posStats.transactionsToday}
          change="+12 from yesterday"
          trend="up"
          icon={ShoppingCart}
          color="blue"
        />
        <StatCard
          title="Average Order"
          value={`$${posStats.averageOrder}`}
          change="+$15.50 from yesterday"
          trend="up"
          icon={Receipt}
          color="indigo"
        />
        <StatCard
          title="Low Stock Items"
          value={posStats.lowStockItems}
          change="Require attention"
          icon={AlertTriangle}
          color="red"
        />
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-8">
          {/* Recent Transactions */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Recent Transactions</h3>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Transaction ID
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Amount
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Items
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Customer
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {recentTransactions.map((transaction, index) => (
                      <tr key={index} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="font-medium text-gray-900">{transaction.id}</div>
                          <div className="text-sm text-gray-500">{transaction.time}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          ${transaction.amount}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {transaction.items} items
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {transaction.customer}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 py-1 text-xs rounded-full ${getStatusColor(transaction.status)}`}>
                            {transaction.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          {/* Top Products */}
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Top Selling Products</h3>
              <div className="space-y-4">
                {topProducts.map((product, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                    <div className="flex items-center space-x-4">
                      <div className="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center">
                        <Package className="h-6 w-6 text-gray-600" />
                      </div>
                      <div>
                        <h4 className="font-medium text-gray-900">{product.name}</h4>
                        <p className="text-sm text-gray-600">
                          {product.sales} sales • Stock: {product.stock}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center space-x-4">
                      <div className="text-right">
                        <p className="font-semibold text-gray-900">{product.revenue}</p>
                        <div className="flex items-center">
                          {getTrendIcon(product.trend)}
                        </div>
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

          {/* Sales Chart */}
          <ChartCard title="Sales Trends">
            <div className="h-64 flex items-center justify-center bg-gray-50 rounded-lg">
              <div className="text-center">
                <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">Sales trends chart will be implemented here</p>
                <p className="text-sm text-gray-400">Shows hourly/daily/weekly sales patterns</p>
              </div>
            </div>
          </ChartCard>
        </div>

        {/* Sidebar */}
        <div className="space-y-8">
          {/* Quick Actions */}
          <QuickActions actions={quickActions} />

          {/* Low Stock Alert */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4 flex items-center">
              <AlertTriangle className="h-5 w-5 text-red-500 mr-2" />
              Low Stock Alert
            </h3>
            <div className="space-y-3">
              {lowStockProducts.map((product, index) => (
                <div key={index} className="p-3 bg-red-50 border border-red-200 rounded-lg">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="font-medium text-gray-900">{product.name}</p>
                      <p className="text-sm text-gray-600">{product.category}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-medium text-red-600">
                        {product.currentStock} left
                      </p>
                      <p className="text-xs text-gray-500">
                        Min: {product.minStock}
                      </p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            <button className="w-full mt-4 bg-red-600 text-white py-2 px-4 rounded-md hover:bg-red-700 transition-colors">
              Reorder All
            </button>
          </div>

          {/* Daily Summary */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Today's Summary</h3>
            <div className="space-y-4">
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Total Sales</span>
                <span className="text-sm font-medium">$12,450</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Transactions</span>
                <span className="text-sm font-medium">89</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Refunds</span>
                <span className="text-sm font-medium">3 ($245)</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Net Sales</span>
                <span className="text-sm font-medium text-green-600">$12,205</span>
              </div>
              <div className="pt-2 border-t">
                <div className="flex items-center text-sm">
                  <CheckCircle className="h-4 w-4 text-green-500 mr-2" />
                  <span className="text-green-600">+18% from yesterday</span>
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
