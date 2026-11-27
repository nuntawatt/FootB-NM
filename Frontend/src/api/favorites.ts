import apiClient from './client'
import type { Team } from '@/types'

export async function fetchFavorites(): Promise<Team[]> {
  const { data } = await apiClient.get<Team[]>('/favorites')
  return data
}

export async function addFavorite(teamId: string): Promise<void> {
  await apiClient.post(`/favorites/${teamId}`)
}

export async function removeFavorite(teamId: string): Promise<void> {
  await apiClient.delete(`/favorites/${teamId}`)
}
