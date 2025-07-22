// API configuration
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/api/v1/auth/login',
    REGISTER: '/api/v1/auth/register',
    REFRESH: '/api/v1/auth/refresh',
    LOGOUT: '/api/v1/auth/logout',
    PROFILE: '/api/v1/auth/profile',
  },
} as const

// Types
export interface User {
  id: string
  tenant_id: string
  email: string
  first_name: string
  last_name: string
  role: string
  is_active: boolean
  is_verified: boolean
  last_login_at: string | null
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
  tenant_id?: string
}

export interface TokenResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: User
}

export interface ApiError {
  error: string
  message: string
}

// API Client
class ApiClient {
  private baseURL: string
  private accessToken: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    this.loadTokenFromStorage()
  }

  private loadTokenFromStorage() {
    if (typeof window !== 'undefined') {
      this.accessToken = localStorage.getItem('access_token')
    }
  }

  setAccessToken(token: string | null) {
    this.accessToken = token
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem('access_token', token)
      } else {
        localStorage.removeItem('access_token')
      }
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    }

    if (this.accessToken) {
      headers.Authorization = `Bearer ${this.accessToken}`
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error: ApiError = await response.json()
      throw new Error(error.message || 'An error occurred')
    }

    return response.json()
  }

  // Auth methods
  async login(data: LoginRequest): Promise<TokenResponse> {
    const response = await this.request<TokenResponse>(
      API_ENDPOINTS.AUTH.LOGIN,
      {
        method: 'POST',
        body: JSON.stringify(data),
      }
    )

    this.setAccessToken(response.access_token)
    if (typeof window !== 'undefined') {
      localStorage.setItem('refresh_token', response.refresh_token)
    }

    return response
  }

  async register(data: RegisterRequest): Promise<TokenResponse> {
    const response = await this.request<TokenResponse>(
      API_ENDPOINTS.AUTH.REGISTER,
      {
        method: 'POST',
        body: JSON.stringify(data),
      }
    )

    this.setAccessToken(response.access_token)
    if (typeof window !== 'undefined') {
      localStorage.setItem('refresh_token', response.refresh_token)
    }

    return response
  }

  async logout(): Promise<void> {
    try {
      await this.request(API_ENDPOINTS.AUTH.LOGOUT, {
        method: 'POST',
      })
    } finally {
      this.setAccessToken(null)
      if (typeof window !== 'undefined') {
        localStorage.removeItem('refresh_token')
      }
    }
  }

  async getProfile(): Promise<User> {
    return this.request<User>(API_ENDPOINTS.AUTH.PROFILE)
  }

  async refreshToken(): Promise<TokenResponse> {
    const refreshToken = typeof window !== 'undefined' 
      ? localStorage.getItem('refresh_token')
      : null

    if (!refreshToken) {
      throw new Error('No refresh token available')
    }

    const response = await this.request<TokenResponse>(
      API_ENDPOINTS.AUTH.REFRESH,
      {
        method: 'POST',
        body: JSON.stringify({ refresh_token: refreshToken }),
      }
    )

    this.setAccessToken(response.access_token)
    if (typeof window !== 'undefined') {
      localStorage.setItem('refresh_token', response.refresh_token)
    }

    return response
  }
}

export const apiClient = new ApiClient(API_BASE_URL)
