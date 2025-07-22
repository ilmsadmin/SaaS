import { api } from './client'

export interface InviteUserRequest {
  email: string
  role: string
}

export interface UserInvitation {
  id: string
  tenant_id: string
  email: string
  role: string
  status: string
  expires_at: string
  created_at: string
  updated_at: string
}

export interface AcceptInvitationRequest {
  token: string
  password: string
  first_name: string
  last_name: string
}

export interface ResetPasswordRequest {
  email: string
}

export interface ConfirmResetPasswordRequest {
  token: string
  new_password: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export interface UpdateProfileRequest {
  first_name?: string
  last_name?: string
  avatar?: string
  phone?: string
  address?: string
  city?: string
  country?: string
  postal_code?: string
  date_of_birth?: string
  bio?: string
  language?: string
  timezone?: string
}

// User Invitation API
export const inviteUser = async (data: InviteUserRequest) => {
  const response = await api.post('/invitations', data)
  return response.data
}

export const getInvitations = async (): Promise<UserInvitation[]> => {
  const response = await api.get('/invitations')
  return response.data.invitations
}

export const acceptInvitation = async (data: AcceptInvitationRequest) => {
  const response = await api.post('/invitations/accept', data)
  return response.data
}

export const revokeInvitation = async (invitationId: string) => {
  const response = await api.delete(`/invitations/${invitationId}`)
  return response.data
}

export const resendInvitation = async (invitationId: string) => {
  const response = await api.post(`/invitations/${invitationId}/resend`)
  return response.data
}

// Password Management API
export const requestPasswordReset = async (data: ResetPasswordRequest) => {
  const response = await api.post('/password/reset/request', data)
  return response.data
}

export const confirmPasswordReset = async (data: ConfirmResetPasswordRequest) => {
  const response = await api.post('/password/reset/confirm', data)
  return response.data
}

export const changePassword = async (data: ChangePasswordRequest) => {
  const response = await api.post('/password/change', data)
  return response.data
}

// Email Verification API
export const sendVerificationEmail = async () => {
  const response = await api.post('/email/verify/send')
  return response.data
}

export const verifyEmail = async (token: string) => {
  const response = await api.post(`/email/verify?token=${token}`)
  return response.data
}

export const resendEmailVerification = async (data: { email: string }) => {
  const response = await api.post('/email/verify/resend', data)
  return response.data
}

// Profile Management API
export const getProfile = async () => {
  const response = await api.get('/profile')
  return response.data
}

export const updateProfile = async (data: UpdateProfileRequest) => {
  const response = await api.put('/profile', data)
  return response.data
}
