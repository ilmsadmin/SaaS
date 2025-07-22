import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { 
  getAvailableModules, 
  getTenantModules,
  installModule,
  uninstallModule,
  enableModule,
  disableModule,
  type Module,
  type TenantModule 
} from '@/lib/api/module-management'

interface ModuleListProps {
  className?: string
}

export function ModuleList({ className }: ModuleListProps) {
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const queryClient = useQueryClient()
  
  const { data: availableModules, isLoading: loadingAvailable } = useQuery({
    queryKey: ['available-modules'],
    queryFn: getAvailableModules,
  })

  const { data: tenantModules, isLoading: loadingTenant } = useQuery({
    queryKey: ['tenant-modules'],
    queryFn: getTenantModules,
  })

  const installMutation = useMutation({
    mutationFn: (data: { moduleId: string; version: string }) =>
      installModule({ module_id: data.moduleId, version: data.version }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenant-modules'] })
    },
  })

  const uninstallMutation = useMutation({
    mutationFn: uninstallModule,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenant-modules'] })
    },
  })

  const enableMutation = useMutation({
    mutationFn: enableModule,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenant-modules'] })
    },
  })

  const disableMutation = useMutation({
    mutationFn: disableModule,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenant-modules'] })
    },
  })

  // Get unique categories
  const categories = ['all']
  if (availableModules) {
    const uniqueCategories = Array.from(new Set(availableModules.map(m => m.category)))
    categories.push(...uniqueCategories)
  }

  // Filter modules by category
  const filteredModules = availableModules?.filter(module => 
    selectedCategory === 'all' || module.category === selectedCategory
  ) || []

  // Check if module is installed
  const isInstalled = (moduleId: string) => {
    return tenantModules?.some(tm => tm.module_id === moduleId)
  }

  // Check if module is enabled
  const isEnabled = (moduleId: string) => {
    return tenantModules?.find(tm => tm.module_id === moduleId)?.is_enabled || false
  }

  const getCategoryIcon = (category: string) => {
    const icons: Record<string, string> = {
      crm: '👥',
      hrm: '👔',
      pos: '🏪',
      lms: '📚',
      checkin: '⏰',
      payment: '💳',
      accounting: '📊',
      ecommerce: '🛒',
    }
    return icons[category] || '📦'
  }

  const getStatusColor = (module: Module) => {
    if (!isInstalled(module.id)) {
      return 'bg-gray-100 text-gray-800'
    }
    if (isEnabled(module.id)) {
      return 'bg-green-100 text-green-800'
    }
    return 'bg-yellow-100 text-yellow-800'
  }

  const getStatusText = (module: Module) => {
    if (!isInstalled(module.id)) {
      return 'Available'
    }
    if (isEnabled(module.id)) {
      return 'Enabled'
    }
    return 'Disabled'
  }

  if (loadingAvailable || loadingTenant) {
    return (
      <div className={`space-y-4 ${className}`}>
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="h-48 bg-gray-200 rounded-lg"></div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={className}>
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900 mb-4">Module Management</h2>
        
        {/* Category Filter */}
        <div className="flex flex-wrap gap-2">
          {categories.map((category) => (
            <button
              key={category}
              onClick={() => setSelectedCategory(category)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                selectedCategory === category
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {category === 'all' ? 'All Modules' : category.toUpperCase()}
            </button>
          ))}
        </div>
      </div>

      {/* Modules Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredModules.map((module) => (
          <div
            key={module.id}
            className="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow"
          >
            <div className="p-6">
              {/* Header */}
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center">
                  <span className="text-2xl mr-3">
                    {getCategoryIcon(module.category)}
                  </span>
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900">
                      {module.display_name}
                    </h3>
                    <p className="text-sm text-gray-500">v{module.version}</p>
                  </div>
                </div>
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(module)}`}>
                  {getStatusText(module)}
                </span>
              </div>

              {/* Description */}
              <p className="text-gray-600 text-sm mb-4 line-clamp-3">
                {module.description || 'No description available.'}
              </p>

              {/* Category */}
              <div className="mb-4">
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  {module.category.toUpperCase()}
                </span>
              </div>

              {/* Actions */}
              <div className="flex gap-2">
                {!isInstalled(module.id) ? (
                  <button
                    onClick={() => installMutation.mutate({ moduleId: module.id, version: module.version })}
                    disabled={installMutation.isPending}
                    className="flex-1 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium disabled:opacity-50"
                  >
                    {installMutation.isPending ? 'Installing...' : 'Install'}
                  </button>
                ) : (
                  <>
                    {isEnabled(module.id) ? (
                      <button
                        onClick={() => disableMutation.mutate(module.id)}
                        disabled={disableMutation.isPending}
                        className="flex-1 bg-yellow-600 hover:bg-yellow-700 text-white px-4 py-2 rounded-lg text-sm font-medium disabled:opacity-50"
                      >
                        Disable
                      </button>
                    ) : (
                      <button
                        onClick={() => enableMutation.mutate(module.id)}
                        disabled={enableMutation.isPending}
                        className="flex-1 bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg text-sm font-medium disabled:opacity-50"
                      >
                        Enable
                      </button>
                    )}
                    <button
                      onClick={() => uninstallMutation.mutate(module.id)}
                      disabled={uninstallMutation.isPending}
                      className="px-4 py-2 border border-red-300 text-red-700 hover:bg-red-50 rounded-lg text-sm font-medium disabled:opacity-50"
                    >
                      Uninstall
                    </button>
                  </>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {filteredModules.length === 0 && (
        <div className="text-center py-12">
          <div className="text-gray-500">No modules found in this category</div>
        </div>
      )}
    </div>
  )
}
