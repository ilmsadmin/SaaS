// API service for tenant management
import { 
  Tenant, 
  Plan, 
  Subscription, 
  CreateTenantRequest, 
  UpdateTenantRequest,
  ListResponse 
} from '@/types/tenant'

interface CreateSubscriptionRequest {
  plan_id: string
  billing_cycle: 'monthly' | 'yearly'
}

interface UpdateSubscriptionRequest {
  plan_id?: string
  billing_cycle?: 'monthly' | 'yearly'
  status?: 'active' | 'suspended' | 'cancelled'
}

type TenantListResponse = ListResponse<Tenant>
type PlanListResponse = ListResponse<Plan>

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

class TenantService {
  private async request<T>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}/api${endpoint}`
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    // Add authorization header if token exists
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers = {
        ...config.headers,
        'Authorization': `Bearer ${token}`,
      }
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      return await response.json()
    } catch (error) {
      console.error('API request failed:', error)
      throw error
    }
  }

  // Tenant Management
  async getTenants(): Promise<TenantListResponse> {
    return this.request<TenantListResponse>('/tenants')
  }

  async getTenant(id: string): Promise<Tenant> {
    return this.request<Tenant>(`/tenants/${id}`)
  }

  async createTenant(data: CreateTenantRequest): Promise<Tenant> {
    return this.request<Tenant>('/tenants', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateTenant(id: string, data: UpdateTenantRequest): Promise<Tenant> {
    return this.request<Tenant>(`/tenants/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteTenant(id: string): Promise<void> {
    await this.request(`/tenants/${id}`, {
      method: 'DELETE',
    })
  }

  async activateTenant(id: string): Promise<Tenant> {
    return this.request<Tenant>(`/tenants/${id}/activate`, {
      method: 'POST',
    })
  }

  async suspendTenant(id: string): Promise<Tenant> {
    return this.request<Tenant>(`/tenants/${id}/suspend`, {
      method: 'POST',
    })
  }

  // Subscription Management
  async getTenantSubscription(tenantId: string): Promise<Subscription> {
    return this.request<Subscription>(`/subscriptions/tenant/${tenantId}`)
  }

  async createSubscription(tenantId: string, data: CreateSubscriptionRequest): Promise<Subscription> {
    return this.request<Subscription>(`/subscriptions/tenant/${tenantId}`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateSubscription(tenantId: string, data: UpdateSubscriptionRequest): Promise<Subscription> {
    return this.request<Subscription>(`/subscriptions/tenant/${tenantId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  // Plan Management
  async getPlans(): Promise<PlanListResponse> {
    return this.request<PlanListResponse>('/plans')
  }

  async getPlan(id: string): Promise<Plan> {
    return this.request<Plan>(`/plans/${id}`)
  }

  async createPlan(data: Partial<Plan>): Promise<Plan> {
    return this.request<Plan>('/plans', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updatePlan(id: string, data: Partial<Plan>): Promise<Plan> {
    return this.request<Plan>(`/plans/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deletePlan(id: string): Promise<void> {
    await this.request(`/plans/${id}`, {
      method: 'DELETE',
    })
  }
}

export const tenantService = new TenantService()
export default tenantService
