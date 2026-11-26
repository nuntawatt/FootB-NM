import apiClient from './client'
import type { Team } from '@/types'

export async function fetchTeams(leagueId?: string): Promise<Team[]> {
  const params: Record<string, string> = {}
  if (leagueId) params.league_id = leagueId
  const { data } = await apiClient.get<Team[]>('/teams', { params })
  return data
}

export async function fetchTeam(id: string): Promise<Team> {
  const { data } = await apiClient.get<Team>(`/teams/${id}`)
  return data
}
