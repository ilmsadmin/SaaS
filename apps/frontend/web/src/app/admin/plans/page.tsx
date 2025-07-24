'use client'

import React, { useState } from 'react'
import { Plus, Edit2, Trash2, CheckCircle, XCircle, CreditCard, Users, DollarSign, Calendar } from 'lucide-react'
import { usePlans } from '@/hooks/useTenants'
import { Plan } from '@/types/tenant'

interface CreatePlanForm {
  name: string
  description: string
  price: number
  currency: string
  billing_cycle: 'monthly' | 'yearly'
  max_users: number
  max_storage: number
  features: string[]
}

export default function PlansManagement() {
  const { plans, loading, error, createPlan, updatePlan, deletePlan } = usePlans()
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null)
  const [showEditModal, setShowEditModal] = useState(false)
  const [formData, setFormData] = useState<CreatePlanForm>({
    name: '',
    description: '',
    price: 0,
    currency: 'USD',
    billing_cycle: 'monthly',
    max_users: 10,
    max_storage: 1000,
    features: []
  })
  const [featureInput, setFeatureInput] = useState('')

  const formatPrice = (price: number, currency: string, billing_cycle: string) => {
    const formatter = new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    })
    return `${formatter.format(price)}/${billing_cycle === 'monthly' ? 'mo' : 'yr'}`
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }

  const handleCreatePlan = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createPlan(formData)
      setShowCreateModal(false)
      resetForm()
    } catch (error) {
      console.error('Failed to create plan:', error)
    }
  }

  const handleUpdatePlan = async (e: React.FormEvent) => {
    e.preventDefault()
    if (selectedPlan) {
      try {
        await updatePlan(selectedPlan.id, formData)
        setShowEditModal(false)
        resetForm()
      } catch (error) {
        console.error('Failed to update plan:', error)
      }
    }
  }

  const handleDeletePlan = async (planId: string) => {
    if (confirm('Are you sure you want to delete this plan?')) {
      try {
        await deletePlan(planId)
      } catch (error) {
        console.error('Failed to delete plan:', error)
      }
    }
  }

  const resetForm = () => {
    setFormData({
      name: '',
      description: '',
      price: 0,
      currency: 'USD',
      billing_cycle: 'monthly',
      max_users: 10,
      max_storage: 1000,
      features: []
    })
    setFeatureInput('')
    setSelectedPlan(null)
  }

  const addFeature = () => {
    if (featureInput.trim() && !formData.features.includes(featureInput.trim())) {
      setFormData({
        ...formData,
        features: [...formData.features, featureInput.trim()]
      })
      setFeatureInput('')
    }
  }

  const removeFeature = (index: number) => {
    setFormData({
      ...formData,
      features: formData.features.filter((_, i) => i !== index)
    })
  }

  const openEditModal = (plan: Plan) => {
    setSelectedPlan(plan)
    setFormData({
      name: plan.name,
      description: plan.description || '',
      price: plan.price,
      currency: plan.currency,
      billing_cycle: plan.billing_cycle,
      max_users: plan.max_users || 10,
      max_storage: plan.max_storage || 1000,
      features: plan.features || []
    })
    setShowEditModal(true)
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
                <CreditCard className="h-8 w-8 text-blue-600" />
                Plans Management
              </h1>
              <p className="text-gray-600 mt-2">
                Manage subscription plans, pricing, and features
              </p>
            </div>
            <button
              onClick={() => setShowCreateModal(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg flex items-center gap-2 transition-colors"
            >
              <Plus className="h-5 w-5" />
              Create Plan
            </button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Plans</p>
                <p className="text-2xl font-bold text-gray-900">{plans.length}</p>
              </div>
              <CreditCard className="h-8 w-8 text-blue-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Active Plans</p>
                <p className="text-2xl font-bold text-gray-900">
                  {plans.filter(plan => plan.is_active).length}
                </p>
              </div>
              <CheckCircle className="h-8 w-8 text-green-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Avg Price</p>
                <p className="text-2xl font-bold text-gray-900">
                  ${plans.length > 0 ? Math.round(plans.reduce((sum, plan) => sum + plan.price, 0) / plans.length) : 0}
                </p>
              </div>
              <DollarSign className="h-8 w-8 text-yellow-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Subscribers</p>
                <p className="text-2xl font-bold text-gray-900">0</p>
              </div>
              <Users className="h-8 w-8 text-purple-600" />
            </div>
          </div>
        </div>

        {/* Plans Table */}
        <div className="bg-white rounded-lg shadow">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">All Plans</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Plan Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Price
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Limits
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Created
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {plans.map((plan) => (
                  <tr key={plan.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div>
                        <div className="text-sm font-medium text-gray-900">{plan.name}</div>
                        <div className="text-sm text-gray-500">{plan.description}</div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900">
                        {formatPrice(plan.price, plan.currency, plan.billing_cycle)}
                      </div>
                      <div className="text-sm text-gray-500">
                        {plan.billing_cycle}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="text-sm text-gray-900">
                        {plan.max_users} users
                      </div>
                      <div className="text-sm text-gray-500">
                        {plan.max_storage}MB storage
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        plan.is_active 
                          ? 'text-green-800 bg-green-100'
                          : 'text-red-800 bg-red-100'
                      }`}>
                        {plan.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatDate(plan.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <div className="flex items-center space-x-2">
                        <button
                          onClick={() => openEditModal(plan)}
                          className="text-indigo-600 hover:text-indigo-900"
                        >
                          <Edit2 className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleDeletePlan(plan.id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Create Plan Modal */}
        {showCreateModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Create New Plan</h3>
              <form onSubmit={handleCreatePlan}>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Plan Name
                    </label>
                    <input
                      type="text"
                      value={formData.name}
                      onChange={(e) => setFormData({...formData, name: e.target.value})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Price
                    </label>
                    <input
                      type="number"
                      value={formData.price}
                      onChange={(e) => setFormData({...formData, price: Number(e.target.value)})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      required
                    />
                  </div>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Currency
                    </label>
                    <select
                      value={formData.currency}
                      onChange={(e) => setFormData({...formData, currency: e.target.value})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="USD">USD</option>
                      <option value="EUR">EUR</option>
                      <option value="GBP">GBP</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Billing Cycle
                    </label>
                    <select
                      value={formData.billing_cycle}
                      onChange={(e) => setFormData({...formData, billing_cycle: e.target.value as 'monthly' | 'yearly'})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="monthly">Monthly</option>
                      <option value="yearly">Yearly</option>
                    </select>
                  </div>
                </div>

                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Description
                  </label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows={3}
                  />
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Max Users
                    </label>
                    <input
                      type="number"
                      value={formData.max_users}
                      onChange={(e) => setFormData({...formData, max_users: Number(e.target.value)})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Max Storage (MB)
                    </label>
                    <input
                      type="number"
                      value={formData.max_storage}
                      onChange={(e) => setFormData({...formData, max_storage: Number(e.target.value)})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                </div>

                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Features
                  </label>
                  <div className="flex gap-2 mb-2">
                    <input
                      type="text"
                      value={featureInput}
                      onChange={(e) => setFeatureInput(e.target.value)}
                      placeholder="Add a feature"
                      className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addFeature())}
                    />
                    <button
                      type="button"
                      onClick={addFeature}
                      className="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
                    >
                      Add
                    </button>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {formData.features.map((feature, index) => (
                      <span
                        key={index}
                        className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800"
                      >
                        {feature}
                        <button
                          type="button"
                          onClick={() => removeFeature(index)}
                          className="ml-2 text-blue-600 hover:text-blue-800"
                        >
                          ×
                        </button>
                      </span>
                    ))}
                  </div>
                </div>

                <div className="flex justify-end space-x-3">
                  <button
                    type="button"
                    onClick={() => {setShowCreateModal(false); resetForm()}}
                    className="px-4 py-2 text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                  >
                    Create Plan
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Edit Plan Modal */}
        {showEditModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Edit Plan</h3>
              <form onSubmit={handleUpdatePlan}>
                {/* Same form fields as create modal */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Plan Name
                    </label>
                    <input
                      type="text"
                      value={formData.name}
                      onChange={(e) => setFormData({...formData, name: e.target.value})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Price
                    </label>
                    <input
                      type="number"
                      value={formData.price}
                      onChange={(e) => setFormData({...formData, price: Number(e.target.value)})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      required
                    />
                  </div>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Currency
                    </label>
                    <select
                      value={formData.currency}
                      onChange={(e) => setFormData({...formData, currency: e.target.value})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="USD">USD</option>
                      <option value="EUR">EUR</option>
                      <option value="GBP">GBP</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Billing Cycle
                    </label>
                    <select
                      value={formData.billing_cycle}
                      onChange={(e) => setFormData({...formData, billing_cycle: e.target.value as 'monthly' | 'yearly'})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="monthly">Monthly</option>
                      <option value="yearly">Yearly</option>
                    </select>
                  </div>
                </div>

                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Description
                  </label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({...formData, description: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows={3}
                  />
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Max Users
                    </label>
                    <input
                      type="number"
                      value={formData.max_users}
                      onChange={(e) => setFormData({...formData, max_users: Number(e.target.value)})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Max Storage (MB)
                    </label>
                    <input
                      type="number"
                      value={formData.max_storage}
                      onChange={(e) => setFormData({...formData, max_storage: Number(e.target.value)})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                </div>

                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Features
                  </label>
                  <div className="flex gap-2 mb-2">
                    <input
                      type="text"
                      value={featureInput}
                      onChange={(e) => setFeatureInput(e.target.value)}
                      placeholder="Add a feature"
                      className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addFeature())}
                    />
                    <button
                      type="button"
                      onClick={addFeature}
                      className="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
                    >
                      Add
                    </button>
                  </div>
                  <div className="flex flex-wrap gap-2">
                    {formData.features.map((feature, index) => (
                      <span
                        key={index}
                        className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-blue-100 text-blue-800"
                      >
                        {feature}
                        <button
                          type="button"
                          onClick={() => removeFeature(index)}
                          className="ml-2 text-blue-600 hover:text-blue-800"
                        >
                          ×
                        </button>
                      </span>
                    ))}
                  </div>
                </div>

                <div className="flex justify-end space-x-3">
                  <button
                    type="button"
                    onClick={() => {setShowEditModal(false); resetForm()}}
                    className="px-4 py-2 text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
                  >
                    Update Plan
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
