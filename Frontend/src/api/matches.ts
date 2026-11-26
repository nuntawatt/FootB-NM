import apiClient from './client'
import type { Match } from '@/types'

export interface MatchFilters {
  season?: number
  year?: number
  status?: string
  league_id?: string
  date?: string
  limit?: number
  page?: number
}

export async function fetchMatches(filters: MatchFilters = {}): Promise<Match[]> {
  const params: Record<string, string | number> = {}
  if (filters.season) params.season = filters.season
  if (filters.year) params.year = filters.year
  if (filters.status) params.status = filters.status
  if (filters.league_id) params.league_id = filters.league_id
  if (filters.date) params.date = filters.date
  if (filters.limit) params.limit = filters.limit
  if (filters.page) params.page = filters.page
  const { data } = await apiClient.get<Match[]>('/matches', { params })
  return data
}

export async function fetchMatch(id: string): Promise<Match> {
  const { data } = await apiClient.get<Match>(`/matches/${id}`)
  return data
}
