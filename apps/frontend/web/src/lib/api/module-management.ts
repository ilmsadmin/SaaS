import { api } from './client'

export interface Module {
  id: string
  name: string
  display_name: string
  description?: string
  version: string
  category: string
  icon?: string
  is_active: boolean
  dependencies: string
  permissions: string
}

export interface TenantModule {
  id: string
  tenant_id: string
  module_id: string
  is_enabled: boolean
  config: string
  installed_at: string
  updated_at: string
}

export interface InstallModuleRequest {
  module_id: string
  version: string
  config?: string
}

export interface UpdateModuleConfigRequest {
  config: string
}

export interface TenantConfiguration {
  id: string
  tenant_id: string
  custom_domain?: string
  ssl_enabled: boolean
  custom_css?: string
  custom_javascript?: string
  branding_config: string
  security_config: string
  notification_config: string
  integration_config: string
  feature_flags: string
  data_retention_days: number
  allowed_ips?: string
  two_factor_required: boolean
  password_policy: string
  session_timeout_mins: number
  created_at: string
  updated_at: string
}

export interface UpdateTenantConfigRequest {
  custom_domain?: string
  ssl_enabled?: boolean
  custom_css?: string
  custom_javascript?: string
  branding_config?: string
  security_config?: string
  notification_config?: string
  integration_config?: string
  feature_flags?: string
  data_retention_days?: number
  allowed_ips?: string
  two_factor_required?: boolean
  password_policy?: string
  session_timeout_mins?: number
}

// Module Management API
export const getAvailableModules = async (): Promise<Module[]> => {
  const response = await api.get('/modules')
  return response.data.modules
}

export const getTenantModules = async (): Promise<TenantModule[]> => {
  const response = await api.get('/modules/tenant')
  return response.data.modules
}

export const installModule = async (data: InstallModuleRequest) => {
  const response = await api.post('/modules/install', data)
  return response.data
}

export const uninstallModule = async (moduleId: string) => {
  const response = await api.delete(`/modules/${moduleId}`)
  return response.data
}

export const enableModule = async (moduleId: string) => {
  const response = await api.post(`/modules/${moduleId}/enable`)
  return response.data
}

export const disableModule = async (moduleId: string) => {
  const response = await api.post(`/modules/${moduleId}/disable`)
  return response.data
}

export const updateModuleConfig = async (moduleId: string, data: UpdateModuleConfigRequest) => {
  const response = await api.put(`/modules/${moduleId}/config`, data)
  return response.data
}

// Tenant Configuration API
export const getTenantConfiguration = async (): Promise<TenantConfiguration> => {
  const response = await api.get('/tenant/config')
  return response.data.configuration
}

export const updateTenantConfiguration = async (data: UpdateTenantConfigRequest) => {
  const response = await api.put('/tenant/config', data)
  return response.data
}

export const setupCustomDomain = async (domain: string, sslEnabled: boolean) => {
  const response = await api.post('/tenant/config/domain', { domain, ssl_enabled: sslEnabled })
  return response.data
}

export const removeCustomDomain = async () => {
  const response = await api.delete('/tenant/config/domain')
  return response.data
}

export const getFeatureFlags = async () => {
  const response = await api.get('/tenant/config/features')
  return response.data.feature_flags
}

export const updateFeatureFlag = async (flagName: string, value: any) => {
  const response = await api.put(`/tenant/config/features/${flagName}`, { value })
  return response.data
}
