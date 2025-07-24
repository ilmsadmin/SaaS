'use client'

import React, { createContext, useContext, useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'

interface AdminUser {
  id: string
  email: string
  name: string
  role: 'super_admin' | 'admin'
  permissions: string[]
}

interface AdminContextType {
  user: AdminUser | null
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  isLoading: boolean
  isAuthenticated: boolean
}

const AdminContext = createContext<AdminContextType | undefined>(undefined)

export function AdminProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AdminUser | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const router = useRouter()

  const isAuthenticated = !!user

  useEffect(() => {
    // Check if user is authenticated on mount
    const token = localStorage.getItem('admin_token')
    if (token) {
      // Validate token and get user info
      validateToken(token)
    } else {
      setIsLoading(false)
    }
  }, [])

  const validateToken = async (token: string) => {
    console.log('Validating token...', token ? 'Token exists' : 'No token')
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/auth/validate`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      
      console.log('Validate response status:', response.status)
      
      if (response.ok) {
        const userData = await response.json()
        console.log('User data received:', userData)
        
        // The API returns the user object directly, not wrapped in a 'user' field
        const adminUser: AdminUser = {
          id: userData.id,
          email: userData.email,
          name: `${userData.first_name} ${userData.last_name}`,
          role: userData.role === 'super_admin' ? 'super_admin' : 'admin',
          permissions: [] // Will be populated based on role
        }
        setUser(adminUser)
        console.log('User set successfully:', adminUser)
      } else {
        console.log('Token validation failed, removing tokens')
        localStorage.removeItem('admin_token')
        localStorage.removeItem('admin_refresh_token')
      }
    } catch (error) {
      console.error('Token validation failed:', error)
      // Don't remove tokens on network errors, only on auth failures
      if (error instanceof TypeError && error.message.includes('fetch')) {
        console.log('Network error - keeping tokens for retry')
      } else {
        localStorage.removeItem('admin_token')
        localStorage.removeItem('admin_refresh_token')
      }
    } finally {
      setIsLoading(false)
    }
  }

  const login = async (email: string, password: string) => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      })

      if (!response.ok) {
        throw new Error('Login failed')
      }

      const data = await response.json()
      
      // Store both access and refresh tokens
      localStorage.setItem('admin_token', data.access_token)
      localStorage.setItem('admin_refresh_token', data.refresh_token)
      
      // Map the API response to our AdminUser interface
      const adminUser: AdminUser = {
        id: data.user.id,
        email: data.user.email,
        name: `${data.user.first_name} ${data.user.last_name}`,
        role: data.user.role === 'super_admin' ? 'super_admin' : 'admin',
        permissions: [] // Will be populated based on role
      }
      
      setUser(adminUser)
      router.push('/dashboard')
    } catch (error) {
      throw error
    }
  }

  const logout = () => {
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_refresh_token')
    setUser(null)
    router.push('/login')
  }

  return (
    <AdminContext.Provider value={{
      user,
      login,
      logout,
      isLoading,
      isAuthenticated
    }}>
      {children}
    </AdminContext.Provider>
  )
}

export function useAdmin() {
  const context = useContext(AdminContext)
  if (context === undefined) {
    throw new Error('useAdmin must be used within an AdminProvider')
  }
  return context
}
