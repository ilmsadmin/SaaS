'use client'

import { useState, useEffect, useCallback } from 'react'
import { Tenant, Plan, Subscription } from '@/types/tenant'
import tenantService from '@/lib/tenant-service'

export const useTenants = () => {
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTenants = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await tenantService.getTenants()
      setTenants(response.data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch tenants')
      console.error('Error fetching tenants:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  const createTenant = async (data: any) => {
    try {
      const newTenant = await tenantService.createTenant(data)
      setTenants(prev => [...prev, newTenant])
      return newTenant
    } catch (err) {
      const error = err instanceof Error ? err.message : 'Failed to create tenant'
      setError(error)
      throw new Error(error)
    }
  }

  const updateTenant = async (id: string, data: any) => {
    try {
      const updatedTenant = await tenantService.updateTenant(id, data)
      setTenants(prev => prev.map(tenant => 
        tenant.id === id ? updatedTenant : tenant
      ))
      return updatedTenant
    } catch (err) {
      const error = err instanceof Error ? err.message : 'Failed to update tenant'
      setError(error)
      throw new Error(error)
    }
  }

  const deleteTenant = async (id: string) => {
    try {
      await tenantService.deleteTenant(id)
      setTenants(prev => prev.filter(tenant => tenant.id !== id))
    } catch (err) {
      const error = err instanceof Error ? err.message : 'Failed to delete tenant'
      setError(error)
      throw new Error(error)
    }
  }

  const activateTenant = async (id: string) => {
    try {
      const updatedTenant = await tenantService.activateTenant(id)
      setTenants(prev => prev.map(tenant => 
        tenant.id === id ? updatedTenant : tenant
      ))
      return updatedTenant
    } catch (err) {
      const error = err instanceof Error ? err.message : 'Failed to activate tenant'
      setError(error)
      throw new Error(error)
    }
  }

  const suspendTenant = async (id: string) => {
    try {
      const updatedTenant = await tenantService.suspendTenant(id)
      setTenants(prev => prev.map(tenant => 
        tenant.id === id ? updatedTenant : tenant
      ))
      return updatedTenant
    } catch (err) {
      const error = err instanceof Error ? err.message : 'Failed to suspend tenant'
      setError(error)
      throw new Error(error)
    }
  }

  useEffect(() => {
    fetchTenants()
  }, [fetchTenants])

  return {
    tenants,
    loading,
    error,
    refetch: fetchTenants,
    createTenant,
    updateTenant,
    deleteTenant,
    activateTenant,
    suspendTenant
  }
}

export const usePlans = () => {
  const [plans, setPlans] = useState<Plan[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchPlans = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await tenantService.getPlans()
      setPlans(response.data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch plans')
      console.error('Error fetching plans:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  const createPlan = async (data: Partial<Plan>) => {
    try {
      const newPlan = await tenantService.createPlan(data)
      setPlans(prev => [...prev, newPlan])
      return newPlan
    } catch (err) {
      const error = err instanceof Error ? err.message : 'Failed to create plan'
      setError(error)
      throw new Error(error)
    }
  }

  const updatePlan = async (id: string, data: Partial<Plan>) => {
    try {
      const updatedPlan = await tenantService.updatePlan(id, data)
      setPlans(prev => prev.map(plan => 
        plan.id === id ? updatedPlan : plan
      ))
      return updatedPlan
    } catch (err) {
      const error = err instanceof Error ? err.message : 'Failed to update plan'
      setError(error)
      throw new Error(error)
    }
  }

  useEffect(() => {
    fetchPlans()
  }, [fetchPlans])

  return {
    plans,
    loading,
    error,
    refetch: fetchPlans,
    createPlan,
    updatePlan
  }
}

export const useTenantSubscription = (tenantId: string | null) => {
  const [subscription, setSubscription] = useState<Subscription | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchSubscription = useCallback(async () => {
    if (!tenantId) return

    try {
      setLoading(true)
      setError(null)
      const response = await tenantService.getTenantSubscription(tenantId)
      setSubscription(response)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch subscription')
      console.error('Error fetching subscription:', err)
    } finally {
      setLoading(false)
    }
  }, [tenantId])

  useEffect(() => {
    fetchSubscription()
  }, [fetchSubscription])

  return {
    subscription,
    loading,
    error,
    refetch: fetchSubscription
  }
}
