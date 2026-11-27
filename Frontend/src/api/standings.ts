import apiClient from './client'
import type { ApiStandingsRow } from '@/types'

export async function fetchStandings(season?: number, leagueId?: string): Promise<ApiStandingsRow[]> {
  const params: Record<string, string | number> = {}
  if (season) params.season = season
  if (leagueId) params.league_id = leagueId
  const { data } = await apiClient.get<ApiStandingsRow[]>('/standings', { params })
  return data
}
