import apiClient from './client'
import type { League } from '@/types'

export async function fetchLeagues(): Promise<League[]> {
  const { data } = await apiClient.get<League[]>('/leagues')
  return data
}

export async function fetchLeague(id: string): Promise<League> {
  const { data } = await apiClient.get<League>(`/leagues/${id}`)
  return data
}
