// Mock API service for tenant management (for development)
import { 
  Tenant, 
  Plan, 
  Subscription,
  CreateTenantRequest, 
  UpdateTenantRequest,
  ListResponse 
} from '@/types/tenant'

// Mock data
const mockTenants: Tenant[] = [
  {
    id: '1',
    name: 'Acme Corporation',
    subdomain: 'acme',
    domain: 'acme.com',
    logo: 'https://via.placeholder.com/100',
    status: 'active',
    settings: { theme: 'light', timezone: 'UTC' },
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-20T10:00:00Z',
    user_count: 25,
    subscription: {
      plan_name: 'Professional',
      status: 'active',
      expires_at: '2024-12-31T23:59:59Z'
    }
  },
  {
    id: '2',
    name: 'Tech Startup Inc',
    subdomain: 'techstartup',
    domain: 'techstartup.io',
    status: 'trial',
    settings: { theme: 'dark', timezone: 'America/New_York' },
    created_at: '2024-02-01T08:30:00Z',
    updated_at: '2024-02-01T08:30:00Z',
    user_count: 5,
    subscription: {
      plan_name: 'Starter',
      status: 'trial',
      expires_at: '2024-03-01T23:59:59Z'
    }
  },
  {
    id: '3',
    name: 'Global Enterprises',
    subdomain: 'global',
    domain: 'globalent.com',
    status: 'suspended',
    settings: { theme: 'light', timezone: 'Europe/London' },
    created_at: '2024-01-01T12:00:00Z',
    updated_at: '2024-02-15T09:00:00Z',
    user_count: 150,
    subscription: {
      plan_name: 'Enterprise',
      status: 'suspended',
      expires_at: '2024-06-30T23:59:59Z'
    }
  }
]

const mockPlans: Plan[] = [
  {
    id: '1',
    name: 'Starter',
    description: 'Perfect for small teams getting started',
    price: 29,
    currency: 'USD',
    billing_cycle: 'monthly',
    max_users: 10,
    max_storage: 5000,
    features: ['Basic Dashboard', 'User Management', 'Email Support', '5GB Storage'],
    is_active: true,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z'
  },
  {
    id: '2',
    name: 'Professional',
    description: 'Ideal for growing businesses',
    price: 79,
    currency: 'USD',
    billing_cycle: 'monthly',
    max_users: 50,
    max_storage: 25000,
    features: ['Advanced Analytics', 'API Access', 'Priority Support', '25GB Storage', 'Custom Integrations'],
    is_active: true,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z'
  },
  {
    id: '3',
    name: 'Enterprise',
    description: 'For large organizations with advanced needs',
    price: 199,
    currency: 'USD',
    billing_cycle: 'monthly',
    max_users: 500,
    max_storage: 100000,
    features: ['White-label Solution', 'Dedicated Support', 'Custom Development', 'Unlimited Storage', 'SSO Integration'],
    is_active: true,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z'
  }
]

class MockTenantService {
  private isMounted = false
  
  constructor() {
    // Only run client-side code after mounting
    if (typeof window !== 'undefined') {
      this.isMounted = true
    }
  }

  private delay(ms: number = 500): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms))
  }

  private generateId(): string {
    // Use a more predictable ID generation for SSR consistency
    return this.isMounted 
      ? Math.random().toString(36).substr(2, 9)
      : `temp-${Date.now()}-${Math.floor(Math.random() * 1000)}`
  }

  private getTimestamp(): string {
    // Return consistent timestamp during SSR
    return this.isMounted 
      ? new Date().toISOString()
      : '2024-02-20T10:00:00.000Z'
  }

  async getTenants(): Promise<ListResponse<Tenant>> {
    await this.delay()
    return {
      data: [...mockTenants],
      total: mockTenants.length,
      page: 1,
      per_page: 20,
      total_pages: 1
    }
  }

  async getTenant(id: string): Promise<Tenant> {
    await this.delay()
    const tenant = mockTenants.find(t => t.id === id)
    if (!tenant) {
      throw new Error('Tenant not found')
    }
    return tenant
  }

  async createTenant(data: CreateTenantRequest): Promise<Tenant> {
    await this.delay()
    const newTenant: Tenant = {
      id: this.generateId(),
      name: data.name,
      subdomain: data.subdomain,
      domain: data.domain,
      status: 'trial',
      settings: {},
      created_at: this.getTimestamp(),
      updated_at: this.getTimestamp(),
      user_count: 0
    }
    mockTenants.push(newTenant)
    return newTenant
  }

  async updateTenant(id: string, data: UpdateTenantRequest): Promise<Tenant> {
    await this.delay()
    const index = mockTenants.findIndex(t => t.id === id)
    if (index === -1) {
      throw new Error('Tenant not found')
    }
    
    mockTenants[index] = {
      ...mockTenants[index],
      ...data,
      updated_at: this.getTimestamp()
    }
    
    return mockTenants[index]
  }

  async deleteTenant(id: string): Promise<void> {
    await this.delay()
    const index = mockTenants.findIndex(t => t.id === id)
    if (index === -1) {
      throw new Error('Tenant not found')
    }
    mockTenants.splice(index, 1)
  }

  async activateTenant(id: string): Promise<Tenant> {
    await this.delay()
    const index = mockTenants.findIndex(t => t.id === id)
    if (index === -1) {
      throw new Error('Tenant not found')
    }
    
    mockTenants[index].status = 'active'
    mockTenants[index].updated_at = this.getTimestamp()
    
    return mockTenants[index]
  }

  async suspendTenant(id: string): Promise<Tenant> {
    await this.delay()
    const index = mockTenants.findIndex(t => t.id === id)
    if (index === -1) {
      throw new Error('Tenant not found')
    }
    
    mockTenants[index].status = 'suspended'
    mockTenants[index].updated_at = this.getTimestamp()
    
    return mockTenants[index]
  }

  // Plans
  async getPlans(): Promise<ListResponse<Plan>> {
    await this.delay()
    return {
      data: [...mockPlans],
      total: mockPlans.length,
      page: 1,
      per_page: 20,
      total_pages: 1
    }
  }

  async createPlan(data: Partial<Plan>): Promise<Plan> {
    await this.delay()
    const newPlan: Plan = {
      id: this.generateId(),
      name: data.name || '',
      description: data.description || '',
      price: data.price || 0,
      currency: data.currency || 'USD',
      billing_cycle: data.billing_cycle || 'monthly',
      max_users: data.max_users || 10,
      max_storage: data.max_storage || 1000,
      features: data.features || [],
      is_active: data.is_active ?? true,
      created_at: this.getTimestamp(),
      updated_at: this.getTimestamp()
    }
    mockPlans.push(newPlan)
    return newPlan
  }

  async updatePlan(id: string, data: Partial<Plan>): Promise<Plan> {
    await this.delay()
    const index = mockPlans.findIndex(p => p.id === id)
    if (index === -1) {
      throw new Error('Plan not found')
    }
    
    mockPlans[index] = {
      ...mockPlans[index],
      ...data,
      updated_at: this.getTimestamp()
    }
    
    return mockPlans[index]
  }

  async deletePlan(id: string): Promise<void> {
    await this.delay()
    const index = mockPlans.findIndex(p => p.id === id)
    if (index === -1) {
      throw new Error('Plan not found')
    }
    mockPlans.splice(index, 1)
  }

  // Subscriptions
  async getTenantSubscription(tenantId: string): Promise<Subscription | null> {
    await this.delay()
    const tenant = mockTenants.find(t => t.id === tenantId)
    if (!tenant || !tenant.subscription) {
      return null
    }
    // Convert the subscription data to match the Subscription interface
    return {
      id: this.generateId(),
      tenant_id: tenantId,
      plan_id: '1', // Default plan ID
      status: tenant.subscription.status as 'active' | 'cancelled' | 'expired',
      current_period_start: this.getTimestamp(),
      current_period_end: tenant.subscription.expires_at,
      cancel_at_period_end: false,
      created_at: tenant.created_at,
      updated_at: tenant.updated_at
    }
  }
}

const mockTenantService = new MockTenantService()
export default mockTenantService
