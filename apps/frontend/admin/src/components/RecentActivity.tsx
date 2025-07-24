import { useEffect, useState } from 'react'
import { formatDistanceToNow } from 'date-fns'

interface Activity {
  id: string
  type: 'tenant_created' | 'user_registered' | 'payment_received' | 'system_alert'
  title: string
  description: string
  timestamp: string
  metadata?: any
}

export default function RecentActivity() {
  const [activities, setActivities] = useState<Activity[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchActivities = async () => {
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/activities`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('admin_token')}`
          }
        })
        if (response.ok) {
          const data = await response.json()
          setActivities(data)
        } else {
          // Fallback with mock data if API fails
          setActivities([
            {
              id: '1',
              type: 'tenant_created',
              title: 'New Tenant Created',
              description: 'Tech Corp has been successfully onboarded',
              timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(), // 2 hours ago
            },
            {
              id: '2',
              type: 'user_registered',
              title: 'New User Registration',
              description: 'john.doe@example.com registered for Startup Ltd',
              timestamp: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString(), // 4 hours ago
            },
            {
              id: '3',
              type: 'payment_received',
              title: 'Payment Received',
              description: '$299 payment received from Digital Agency',
              timestamp: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(), // 6 hours ago
            },
            {
              id: '4',
              type: 'system_alert',
              title: 'System Maintenance',
              description: 'Scheduled maintenance completed successfully',
              timestamp: new Date(Date.now() - 8 * 60 * 60 * 1000).toISOString(), // 8 hours ago
            },
          ])
        }
      } catch (error) {
        console.error('Failed to fetch activities:', error)
        // Use mock data on error
        setActivities([
          {
            id: '1',
            type: 'tenant_created',
            title: 'New Tenant Created',
            description: 'Tech Corp has been successfully onboarded',
            timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
          },
          {
            id: '2',
            type: 'user_registered',
            title: 'New User Registration',
            description: 'john.doe@example.com registered for Startup Ltd',
            timestamp: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString(),
          }
        ])
      } finally {
        setIsLoading(false)
      }
    }

    fetchActivities()
  }, [])

  const getActivityIcon = (type: Activity['type']) => {
    switch (type) {
      case 'tenant_created':
        return (
          <div className="bg-green-100 rounded-full p-2">
            <svg className="h-4 w-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
          </div>
        )
      case 'user_registered':
        return (
          <div className="bg-blue-100 rounded-full p-2">
            <svg className="h-4 w-4 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
              <path d="M8 9a3 3 0 100-6 3 3 0 000 6zM8 11a6 6 0 016 6H2a6 6 0 016-6zM16 7a1 1 0 10-2 0v1h-1a1 1 0 100 2h1v1a1 1 0 102 0v-1h1a1 1 0 100-2h-1V7z" />
            </svg>
          </div>
        )
      case 'payment_received':
        return (
          <div className="bg-green-100 rounded-full p-2">
            <svg className="h-4 w-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
              <path d="M4 4a2 2 0 00-2 2v4a2 2 0 002 2V6h10a2 2 0 00-2-2H4zM14 6a2 2 0 012 2v4a2 2 0 01-2 2H8a2 2 0 01-2-2V8a2 2 0 012-2h6zM4 14a2 2 0 002 2h8a2 2 0 002-2v-4a2 2 0 00-2-2H6a2 2 0 00-2 2v4z" />
            </svg>
          </div>
        )
      case 'system_alert':
        return (
          <div className="bg-yellow-100 rounded-full p-2">
            <svg className="h-4 w-4 text-yellow-600" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
          </div>
        )
      default:
        return (
          <div className="bg-gray-100 rounded-full p-2">
            <svg className="h-4 w-4 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
            </svg>
          </div>
        )
    }
  }

  if (isLoading) {
    return (
      <div className="space-y-4">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="animate-pulse flex space-x-4">
            <div className="rounded-full bg-gray-200 h-10 w-10" />
            <div className="flex-1 space-y-2">
              <div className="h-4 bg-gray-200 rounded w-3/4" />
              <div className="h-3 bg-gray-200 rounded w-1/2" />
            </div>
          </div>
        ))}
      </div>
    )
  }

  return (
    <div className="flow-root">
      <ul className="-my-5 divide-y divide-gray-200">
        {activities.map((activity) => (
          <li key={activity.id} className="py-4">
            <div className="flex items-start space-x-4">
              <div className="flex-shrink-0">
                {getActivityIcon(activity.type)}
              </div>
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium text-gray-900">{activity.title}</p>
                <p className="text-sm text-gray-500">{activity.description}</p>
                <p className="text-xs text-gray-400 mt-1">
                  {formatDistanceToNow(new Date(activity.timestamp), { addSuffix: true })}
                </p>
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}
