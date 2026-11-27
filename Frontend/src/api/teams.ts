import apiClient from './client'
import type { Team, ApiFormMatch, Coach, Player } from '@/types'

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

export async function fetchTeamForm(id: string): Promise<ApiFormMatch[]> {
  const { data } = await apiClient.get<ApiFormMatch[]>(`/teams/${id}/form`)
  return data
}

export async function fetchTeamCoach(id: string): Promise<Coach> {
  const { data } = await apiClient.get<Coach>(`/teams/${id}/coach`)
  return data
}

export async function fetchTeamSquad(id: string): Promise<Player[]> {
  const { data } = await apiClient.get<Player[]>(`/teams/${id}/squad`)
  return data
}
