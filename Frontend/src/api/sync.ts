import apiClient from './client'

export interface SyncFixturesPayload {
  season: number
  from?: string
  to?: string
}

export interface SyncFixturesResponse {
  status: string
  season: number
  leagues_synced: number
  teams_created: number
  matches_created: number
  matches_updated: number
  matches_skipped: number
  quota_remaining?: {
    day: number
    minute: number
  }
}

export async function syncFixtures(payload: SyncFixturesPayload): Promise<SyncFixturesResponse> {
  const { data } = await apiClient.post<SyncFixturesResponse>('/sync/fixtures', payload)
  return data
}

export async function syncSquads(): Promise<Record<string, unknown>> {
  const { data } = await apiClient.post<Record<string, unknown>>('/sync/squads')
  return data
}
