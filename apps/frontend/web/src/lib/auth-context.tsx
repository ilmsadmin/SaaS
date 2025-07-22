'use client'

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient, User, LoginRequest, RegisterRequest } from './api'

interface AuthContextType {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (data: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
  refetch: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  const queryClient = useQueryClient()
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  // Query to get user profile
  const {
    data: user,
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ['auth', 'profile'],
    queryFn: async () => {
      try {
        return await apiClient.getProfile()
      } catch (error) {
        // If profile fetch fails, user is not authenticated
        setIsAuthenticated(false)
        return null
      }
    },
    enabled: isAuthenticated,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: false,
  })

  // Login mutation
  const loginMutation = useMutation({
    mutationFn: apiClient.login.bind(apiClient),
    onSuccess: () => {
      setIsAuthenticated(true)
      queryClient.invalidateQueries({ queryKey: ['auth'] })
    },
    onError: (error) => {
      console.error('Login failed:', error)
      throw error
    },
  })

  // Register mutation
  const registerMutation = useMutation({
    mutationFn: apiClient.register.bind(apiClient),
    onSuccess: () => {
      setIsAuthenticated(true)
      queryClient.invalidateQueries({ queryKey: ['auth'] })
    },
    onError: (error) => {
      console.error('Registration failed:', error)
      throw error
    },
  })

  // Logout mutation
  const logoutMutation = useMutation({
    mutationFn: apiClient.logout.bind(apiClient),
    onSuccess: () => {
      setIsAuthenticated(false)
      queryClient.clear() // Clear all cached data
    },
    onError: (error) => {
      console.error('Logout failed:', error)
      // Still set as not authenticated even if logout API fails
      setIsAuthenticated(false)
      queryClient.clear()
    },
  })

  // Check if user is authenticated on mount
  useEffect(() => {
    const token = typeof window !== 'undefined' 
      ? localStorage.getItem('access_token')
      : null

    if (token) {
      setIsAuthenticated(true)
    }
  }, [])

  const login = async (data: LoginRequest) => {
    await loginMutation.mutateAsync(data)
  }

  const register = async (data: RegisterRequest) => {
    await registerMutation.mutateAsync(data)
  }

  const logout = async () => {
    await logoutMutation.mutateAsync()
  }

  const value: AuthContextType = {
    user: user || null,
    isLoading: isLoading || loginMutation.isPending || registerMutation.isPending,
    isAuthenticated,
    login,
    register,
    logout,
    refetch,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

// Hook for protected routes
export function useRequireAuth() {
  const { isAuthenticated, isLoading } = useAuth()
  
  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      // Redirect to login page
      window.location.href = '/login'
    }
  }, [isAuthenticated, isLoading])

  return { isAuthenticated, isLoading }
}
