'use client'

import React, { useState } from 'react'
import { Plus, Edit2, Trash2, CheckCircle, XCircle, Building, Users, Settings } from 'lucide-react'
import { useTenants, usePlans } from '@/hooks/useTenants'
import { Tenant, CreateTenantRequest, UpdateTenantRequest } from '@/types/tenant'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'

export default function TenantManagement() {
  const { tenants, loading: tenantsLoading, error: tenantsError, createTenant, updateTenant, deleteTenant, activateTenant, suspendTenant } = useTenants()
  const { plans, loading: plansLoading } = usePlans()
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [selectedTenant, setSelectedTenant] = useState<Tenant | null>(null)
  const [createFormData, setCreateFormData] = useState<CreateTenantRequest>({
    name: '',
    subdomain: '',
    domain: ''
  })
  const [editFormData, setEditFormData] = useState<UpdateTenantRequest>({
    name: '',
    domain: '',
    logo: '',
    settings: {}
  })

  const loading = tenantsLoading || plansLoading
  const error = tenantsError

  const handleCreateTenant = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createTenant(createFormData)
      setShowCreateModal(false)
      resetCreateForm()
    } catch (error) {
      console.error('Failed to create tenant:', error)
    }
  }

  const handleUpdateTenant = async (e: React.FormEvent) => {
    e.preventDefault()
    if (selectedTenant) {
      try {
        await updateTenant(selectedTenant.id, editFormData)
        setShowEditModal(false)
        resetEditForm()
      } catch (error) {
        console.error('Failed to update tenant:', error)
      }
    }
  }

  const handleDeleteTenant = async (tenantId: string) => {
    if (confirm('Are you sure you want to delete this tenant? This action cannot be undone.')) {
      try {
        await deleteTenant(tenantId)
      } catch (error) {
        console.error('Failed to delete tenant:', error)
      }
    }
  }

  const resetCreateForm = () => {
    setCreateFormData({
      name: '',
      subdomain: '',
      domain: ''
    })
  }

  const resetEditForm = () => {
    setEditFormData({
      name: '',
      domain: '',
      logo: '',
      settings: {}
    })
    setSelectedTenant(null)
  }

  const openEditModal = (tenant: Tenant) => {
    setSelectedTenant(tenant)
    setEditFormData({
      name: tenant.name,
      domain: tenant.domain || '',
      logo: tenant.logo || '',
      settings: tenant.settings || {}
    })
    setShowEditModal(true)
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'text-green-600 bg-green-100'
      case 'suspended':
        return 'text-red-600 bg-red-100'
      case 'pending':
        return 'text-yellow-600 bg-yellow-100'
      default:
        return 'text-gray-600 bg-gray-100'
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }

  const handleSuspendTenant = async (tenantId: string) => {
    try {
      await suspendTenant(tenantId)
    } catch (error) {
      console.error('Failed to suspend tenant:', error)
    }
  }

  const handleActivateTenant = async (tenantId: string) => {
    try {
      await activateTenant(tenantId)
    } catch (error) {
      console.error('Failed to activate tenant:', error)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
                <Building className="h-8 w-8 text-blue-600" />
                Tenant Management
              </h1>
              <p className="text-gray-600 mt-2">
                Manage your SaaS tenants, subscriptions, and access
              </p>
            </div>
            <button
              onClick={() => setShowCreateModal(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg flex items-center gap-2 transition-colors"
            >
              <Plus className="h-5 w-5" />
              Create Tenant
            </button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Tenants</p>
                <p className="text-2xl font-bold text-gray-900">{tenants.length}</p>
              </div>
              <Building className="h-8 w-8 text-blue-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Active Tenants</p>
                <p className="text-2xl font-bold text-green-600">
                  {tenants.filter(t => t.status === 'active').length}
                </p>
              </div>
              <CheckCircle className="h-8 w-8 text-green-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Suspended</p>
                <p className="text-2xl font-bold text-red-600">
                  {tenants.filter(t => t.status === 'suspended').length}
                </p>
              </div>
              <XCircle className="h-8 w-8 text-red-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Users</p>
                <p className="text-2xl font-bold text-gray-900">
                  {tenants.reduce((sum, t) => sum + (t.user_count || 0), 0)}
                </p>
              </div>
              <Users className="h-8 w-8 text-purple-600" />
            </div>
          </div>
        </div>

        {/* Tenants Table */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">All Tenants</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Tenant
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Domain
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Plan
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Users
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Created
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {tenants.map((tenant) => (
                  <tr key={tenant.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center">
                          <Building className="h-5 w-5 text-blue-600" />
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-gray-900">
                            {tenant.name}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">{tenant.domain}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {tenant.subscription?.plan_name || 'No Plan'}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">{tenant.user_count || 0}</div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(tenant.status)}`}>
                        {tenant.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {formatDate(tenant.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div className="flex items-center justify-end space-x-2">
                        <button
                          onClick={() => openEditModal(tenant)}
                          className="text-blue-600 hover:text-blue-900"
                          title="Edit Tenant"
                        >
                          <Edit2 className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleDeleteTenant(tenant.id)}
                          className="text-red-600 hover:text-red-900"
                          title="Delete Tenant"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                        {tenant.status === 'active' ? (
                          <button
                            onClick={() => handleSuspendTenant(tenant.id)}
                            className="text-yellow-600 hover:text-yellow-900"
                            title="Suspend Tenant"
                          >
                            <XCircle className="h-4 w-4" />
                          </button>
                        ) : (
                          <button
                            onClick={() => handleActivateTenant(tenant.id)}
                            className="text-green-600 hover:text-green-900"
                            title="Activate Tenant"
                          >
                            <CheckCircle className="h-4 w-4" />
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Plans Section */}
        <div className="mt-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">Available Plans</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {plans.map((plan) => (
              <div key={plan.id} className="bg-white rounded-lg shadow-lg overflow-hidden">
                <div className="px-6 py-8">
                  <h3 className="text-lg font-semibold text-gray-900">{plan.name}</h3>
                  <p className="text-gray-600 mt-2">{plan.description}</p>
                  <div className="mt-4">
                    <span className="text-3xl font-bold text-gray-900">${plan.price}</span>
                    <span className="text-gray-600">/{plan.billing_cycle}</span>
                  </div>
                  <ul className="mt-6 space-y-2">
                    {plan.features.map((feature, index) => (
                      <li key={index} className="flex items-center text-sm text-gray-600">
                        <CheckCircle className="h-4 w-4 text-green-600 mr-2" />
                        {feature}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Create Tenant Modal */}
        <Dialog open={showCreateModal} onOpenChange={setShowCreateModal}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <Building className="h-5 w-5 text-blue-600" />
                Create New Tenant
              </DialogTitle>
            </DialogHeader>
            <form onSubmit={handleCreateTenant} className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="name">Tenant Name</Label>
                  <Input
                    id="name"
                    type="text"
                    value={createFormData.name}
                    onChange={(e) => setCreateFormData({ ...createFormData, name: e.target.value })}
                    placeholder="e.g., Acme Corporation"
                    required
                  />
                </div>
                <div>
                  <Label htmlFor="subdomain">Subdomain</Label>
                  <Input
                    id="subdomain"
                    type="text"
                    value={createFormData.subdomain}
                    onChange={(e) => setCreateFormData({ ...createFormData, subdomain: e.target.value })}
                    placeholder="e.g., acme"
                    required
                  />
                </div>
              </div>
              
              <div>
                <Label htmlFor="domain">Custom Domain (Optional)</Label>
                <Input
                  id="domain"
                  type="text"
                  value={createFormData.domain || ''}
                  onChange={(e) => setCreateFormData({ ...createFormData, domain: e.target.value })}
                  placeholder="e.g., acme.com"
                />
              </div>

              <div className="flex justify-end space-x-4">
                <Button type="button" variant="outline" onClick={() => setShowCreateModal(false)}>
                  Cancel
                </Button>
                <Button type="submit">
                  Create Tenant
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>

        {/* Edit Tenant Modal */}
        <Dialog open={showEditModal} onOpenChange={setShowEditModal}>
          <DialogContent className="sm:max-w-[600px]">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <Settings className="h-5 w-5 text-blue-600" />
                Edit Tenant
              </DialogTitle>
            </DialogHeader>
            <form onSubmit={handleUpdateTenant} className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="edit-name">Tenant Name</Label>
                  <Input
                    id="edit-name"
                    type="text"
                    value={editFormData.name || ''}
                    onChange={(e) => setEditFormData({ ...editFormData, name: e.target.value })}
                    placeholder="e.g., Acme Corporation"
                    required
                  />
                </div>
                <div>
                  <Label htmlFor="edit-domain">Custom Domain (Optional)</Label>
                  <Input
                    id="edit-domain"
                    type="text"
                    value={editFormData.domain || ''}
                    onChange={(e) => setEditFormData({ ...editFormData, domain: e.target.value })}
                    placeholder="e.g., acme.com"
                  />
                </div>
              </div>

              <div>
                <Label htmlFor="edit-logo">Logo URL (Optional)</Label>
                <Input
                  id="edit-logo"
                  type="url"
                  value={editFormData.logo || ''}
                  onChange={(e) => setEditFormData({ ...editFormData, logo: e.target.value })}
                  placeholder="https://example.com/logo.png"
                />
              </div>

              <div className="flex justify-end space-x-4">
                <Button type="button" variant="outline" onClick={() => setShowEditModal(false)}>
                  Cancel
                </Button>
                <Button type="submit">
                  Update Tenant
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  )
}
