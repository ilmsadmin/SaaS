// Tenant Management API Types

export interface Tenant {
  id: string
  name: string
  subdomain: string
  domain?: string
  logo?: string
  status: 'active' | 'suspended' | 'trial'
  settings: Record<string, unknown>
  created_at: string
  updated_at: string
  // Extended properties for UI
  user_count?: number
  subscription?: {
    plan_name: string
    status: string
    expires_at: string
  }
}

export interface Plan {
  id: string
  name: string
  description?: string
  price: number
  currency: string
  billing_cycle: 'monthly' | 'yearly'
  max_users?: number
  max_storage?: number
  features: string[]
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Subscription {
  id: string
  tenant_id: string
  plan_id: string
  status: 'active' | 'cancelled' | 'expired'
  trial_end_at?: string
  current_period_start: string
  current_period_end: string
  cancel_at_period_end: boolean
  created_at: string
  updated_at: string
  tenant?: Tenant
  plan?: Plan
}

export interface CreateTenantRequest {
  name: string
  subdomain: string
  domain?: string
  plan_id?: string
}

export interface UpdateTenantRequest {
  name?: string
  domain?: string
  logo?: string
  settings?: Record<string, unknown>
}

export interface CreatePlanRequest {
  name: string
  description?: string
  price: number
  currency: string
  billing_cycle: 'monthly' | 'yearly'
  max_users?: number
  max_storage?: number
  features?: string[]
}

export interface UpdatePlanRequest {
  name?: string
  description?: string
  price?: number
  currency?: string
  billing_cycle?: 'monthly' | 'yearly'
  max_users?: number
  max_storage?: number
  features?: string[]
  is_active?: boolean
}

export interface ListResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
  total_pages: number
}

// API Endpoints
export const TENANT_API_ENDPOINTS = {
  TENANTS: {
    LIST: '/api/v1/tenants',
    CREATE: '/api/v1/tenants',
    GET: (id: string) => `/api/v1/tenants/${id}`,
    UPDATE: (id: string) => `/api/v1/tenants/${id}`,
    DELETE: (id: string) => `/api/v1/tenants/${id}`,
    ACTIVATE: (id: string) => `/api/v1/tenants/${id}/activate`,
    SUSPEND: (id: string) => `/api/v1/tenants/${id}/suspend`,
  },
  PLANS: {
    LIST: '/api/v1/plans',
    CREATE: '/api/v1/plans',
    GET: (id: string) => `/api/v1/plans/${id}`,
    UPDATE: (id: string) => `/api/v1/plans/${id}`,
    DELETE: (id: string) => `/api/v1/plans/${id}`,
  },
  SUBSCRIPTIONS: {
    GET: (tenantId: string) => `/api/v1/subscriptions/tenant/${tenantId}`,
    CREATE: (tenantId: string) => `/api/v1/subscriptions/tenant/${tenantId}`,
    UPDATE: (tenantId: string) => `/api/v1/subscriptions/tenant/${tenantId}`,
  }
} as const
