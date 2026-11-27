import apiClient from './client'
import type { AuthResponse, User } from '@/types'

export interface RegisterPayload {
  username: string
  email: string
  password: string
  confirm_password: string
  avatar_url?: string
}

export interface LoginPayload {
  email: string
  password: string
}

export interface UpdateProfilePayload {
  username?: string
  avatar_url?: string
}

export async function registerUser(payload: RegisterPayload): Promise<AuthResponse> {
  const { data } = await apiClient.post<AuthResponse>('/register', payload)
  return data
}

export async function loginUser(payload: LoginPayload): Promise<AuthResponse> {
  const { data } = await apiClient.post<AuthResponse>('/login', payload)
  return data
}

export async function fetchProfile(): Promise<User> {
  const { data } = await apiClient.get<User>('/profile')
  return data
}

export async function updateProfile(
  payload: UpdateProfilePayload | FormData
): Promise<{ status: string; message: string; user: User }> {
  const isFormData = payload instanceof FormData
  const { data } = await apiClient.put<{ status: string; message: string; user: User }>(
    '/profile',
    payload,
    isFormData ? { headers: { 'Content-Type': 'multipart/form-data' } } : undefined
  )
  return data
}
