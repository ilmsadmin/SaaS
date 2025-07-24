'use client'

import React, { useState, useEffect } from 'react'
import { Plus, Edit2, Trash2, CheckCircle, XCircle, Users, Mail, Calendar, Shield } from 'lucide-react'

interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  role: 'super_admin' | 'admin' | 'manager' | 'user'
  status: 'active' | 'inactive' | 'pending'
  email_verified: boolean
  last_login?: string
  created_at: string
  updated_at: string
  tenant?: {
    id: string
    name: string
  }
}

interface CreateUserForm {
  email: string
  first_name: string
  last_name: string
  role: 'super_admin' | 'admin' | 'manager' | 'user'
  password: string
}

export default function UsersManagement() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [showEditModal, setShowEditModal] = useState(false)
  const [formData, setFormData] = useState<CreateUserForm>({
    email: '',
    first_name: '',
    last_name: '',
    role: 'user',
    password: ''
  })

  // Mock data for development
  useEffect(() => {
    const mockUsers: User[] = [
      {
        id: '1',
        email: 'admin@zplus.com',
        first_name: 'Super',
        last_name: 'Admin',
        role: 'super_admin',
        status: 'active',
        email_verified: true,
        last_login: '2024-12-22T10:30:00Z',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-12-22T10:30:00Z'
      },
      {
        id: '2',
        email: 'john.doe@techcorp.com',
        first_name: 'John',
        last_name: 'Doe',
        role: 'admin',
        status: 'active',
        email_verified: true,
        last_login: '2024-12-21T15:20:00Z',
        created_at: '2024-02-15T09:00:00Z',
        updated_at: '2024-12-21T15:20:00Z',
        tenant: {
          id: 'tenant1',
          name: 'TechCorp'
        }
      },
      {
        id: '3',
        email: 'jane.smith@acme.com',
        first_name: 'Jane',
        last_name: 'Smith',
        role: 'manager',
        status: 'active',
        email_verified: true,
        last_login: '2024-12-20T12:45:00Z',
        created_at: '2024-03-10T14:30:00Z',
        updated_at: '2024-12-20T12:45:00Z',
        tenant: {
          id: 'tenant2',
          name: 'Acme Corp'
        }
      },
      {
        id: '4',
        email: 'bob.wilson@startup.com',
        first_name: 'Bob',
        last_name: 'Wilson',
        role: 'user',
        status: 'pending',
        email_verified: false,
        created_at: '2024-12-15T16:20:00Z',
        updated_at: '2024-12-15T16:20:00Z',
        tenant: {
          id: 'tenant3',
          name: 'Startup Inc'
        }
      },
      {
        id: '5',
        email: 'alice.brown@enterprise.com',
        first_name: 'Alice',
        last_name: 'Brown',
        role: 'user',
        status: 'inactive',
        email_verified: true,
        last_login: '2024-11-30T09:15:00Z',
        created_at: '2024-04-05T11:00:00Z',
        updated_at: '2024-11-30T09:15:00Z',
        tenant: {
          id: 'tenant4',
          name: 'Enterprise Solutions'
        }
      }
    ]
    
    setUsers(mockUsers)
    setLoading(false)
  }, [])

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'super_admin':
        return 'text-purple-600 bg-purple-100'
      case 'admin':
        return 'text-blue-600 bg-blue-100'
      case 'manager':
        return 'text-green-600 bg-green-100'
      case 'user':
        return 'text-gray-600 bg-gray-100'
      default:
        return 'text-gray-600 bg-gray-100'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'text-green-600 bg-green-100'
      case 'inactive':
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
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      // Mock create - in real app, this would call an API
      const newUser: User = {
        id: (users.length + 1).toString(),
        ...formData,
        status: 'pending',
        email_verified: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
      setUsers([...users, newUser])
      setShowCreateModal(false)
      resetForm()
    } catch (error) {
      console.error('Failed to create user:', error)
    }
  }

  const handleUpdateUser = async (e: React.FormEvent) => {
    e.preventDefault()
    if (selectedUser) {
      try {
        // Mock update - in real app, this would call an API
        const updatedUser = {
          ...selectedUser,
          ...formData,
          updated_at: new Date().toISOString()
        }
        setUsers(users.map(user => user.id === selectedUser.id ? updatedUser : user))
        setShowEditModal(false)
        resetForm()
      } catch (error) {
        console.error('Failed to update user:', error)
      }
    }
  }

  const handleDeleteUser = async (userId: string) => {
    if (confirm('Are you sure you want to delete this user?')) {
      try {
        // Mock delete - in real app, this would call an API
        setUsers(users.filter(user => user.id !== userId))
      } catch (error) {
        console.error('Failed to delete user:', error)
      }
    }
  }

  const handleToggleUserStatus = async (userId: string) => {
    try {
      // Mock status toggle - in real app, this would call an API
      setUsers(users.map(user => 
        user.id === userId 
          ? { ...user, status: user.status === 'active' ? 'inactive' : 'active' as any }
          : user
      ))
    } catch (error) {
      console.error('Failed to toggle user status:', error)
    }
  }

  const resetForm = () => {
    setFormData({
      email: '',
      first_name: '',
      last_name: '',
      role: 'user',
      password: ''
    })
    setSelectedUser(null)
  }

  const openEditModal = (user: User) => {
    setSelectedUser(user)
    setFormData({
      email: user.email,
      first_name: user.first_name,
      last_name: user.last_name,
      role: user.role,
      password: '' // Don't pre-fill password
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
                <Users className="h-8 w-8 text-blue-600" />
                Users Management
              </h1>
              <p className="text-gray-600 mt-2">
                Manage system users, roles, and permissions
              </p>
            </div>
            <button
              onClick={() => setShowCreateModal(true)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg flex items-center gap-2 transition-colors"
            >
              <Plus className="h-5 w-5" />
              Create User
            </button>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Users</p>
                <p className="text-2xl font-bold text-gray-900">{users.length}</p>
              </div>
              <Users className="h-8 w-8 text-blue-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Active Users</p>
                <p className="text-2xl font-bold text-gray-900">
                  {users.filter(user => user.status === 'active').length}
                </p>
              </div>
              <CheckCircle className="h-8 w-8 text-green-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Pending Users</p>
                <p className="text-2xl font-bold text-gray-900">
                  {users.filter(user => user.status === 'pending').length}
                </p>
              </div>
              <Calendar className="h-8 w-8 text-yellow-600" />
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Admins</p>
                <p className="text-2xl font-bold text-gray-900">
                  {users.filter(user => user.role === 'admin' || user.role === 'super_admin').length}
                </p>
              </div>
              <Shield className="h-8 w-8 text-purple-600" />
            </div>
          </div>
        </div>

        {/* Users Table */}
        <div className="bg-white rounded-lg shadow">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">All Users</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    User
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Role
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Tenant
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Last Login
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
                {users.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <div className="flex-shrink-0 h-10 w-10">
                          <div className="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
                            <span className="text-sm font-medium text-gray-700">
                              {user.first_name[0]}{user.last_name[0]}
                            </span>
                          </div>
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-gray-900">
                            {user.first_name} {user.last_name}
                          </div>
                          <div className="text-sm text-gray-500 flex items-center">
                            <Mail className="h-3 w-3 mr-1" />
                            {user.email}
                            {user.email_verified && (
                              <CheckCircle className="h-3 w-3 ml-1 text-green-500" />
                            )}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getRoleColor(user.role)}`}>
                        {user.role.replace('_', ' ')}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(user.status)}`}>
                        {user.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {user.tenant?.name || 'System User'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {user.last_login ? formatDate(user.last_login) : 'Never'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {formatDate(user.created_at)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <div className="flex items-center space-x-2">
                        <button
                          onClick={() => openEditModal(user)}
                          className="text-indigo-600 hover:text-indigo-900"
                        >
                          <Edit2 className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleToggleUserStatus(user.id)}
                          className={user.status === 'active' ? 'text-red-600 hover:text-red-900' : 'text-green-600 hover:text-green-900'}
                        >
                          {user.status === 'active' ? <XCircle className="h-4 w-4" /> : <CheckCircle className="h-4 w-4" />}
                        </button>
                        <button
                          onClick={() => handleDeleteUser(user.id)}
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

        {/* Create User Modal */}
        {showCreateModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Create New User</h3>
              <form onSubmit={handleCreateUser}>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Email
                  </label>
                  <input
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({...formData, email: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                  />
                </div>
                
                <div className="grid grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      First Name
                    </label>
                    <input
                      type="text"
                      value={formData.first_name}
                      onChange={(e) => setFormData({...formData, first_name: e.target.value})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Last Name
                    </label>
                    <input
                      type="text"
                      value={formData.last_name}
                      onChange={(e) => setFormData({...formData, last_name: e.target.value})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      required
                    />
                  </div>
                </div>

                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Role
                  </label>
                  <select
                    value={formData.role}
                    onChange={(e) => setFormData({...formData, role: e.target.value as 'super_admin' | 'admin' | 'manager' | 'user'})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="user">User</option>
                    <option value="manager">Manager</option>
                    <option value="admin">Admin</option>
                    <option value="super_admin">Super Admin</option>
                  </select>
                </div>

                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Password
                  </label>
                  <input
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({...formData, password: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                  />
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
                    Create User
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Edit User Modal */}
        {showEditModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Edit User</h3>
              <form onSubmit={handleUpdateUser}>
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Email
                  </label>
                  <input
                    type="email"
                    value={formData.email}
                    onChange={(e) => setFormData({...formData, email: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                  />
                </div>
                
                <div className="grid grid-cols-2 gap-4 mb-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      First Name
                    </label>
                    <input
                      type="text"
                      value={formData.first_name}
                      onChange={(e) => setFormData({...formData, first_name: e.target.value})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Last Name
                    </label>
                    <input
                      type="text"
                      value={formData.last_name}
                      onChange={(e) => setFormData({...formData, last_name: e.target.value})}
                      className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                      required
                    />
                  </div>
                </div>

                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Role
                  </label>
                  <select
                    value={formData.role}
                    onChange={(e) => setFormData({...formData, role: e.target.value as 'super_admin' | 'admin' | 'manager' | 'user'})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="user">User</option>
                    <option value="manager">Manager</option>
                    <option value="admin">Admin</option>
                    <option value="super_admin">Super Admin</option>
                  </select>
                </div>

                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    New Password (leave blank to keep current)
                  </label>
                  <input
                    type="password"
                    value={formData.password}
                    onChange={(e) => setFormData({...formData, password: e.target.value})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
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
                    Update User
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
