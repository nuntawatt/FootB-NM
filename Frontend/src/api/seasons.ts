import apiClient from './client'
import type { SeasonItem } from '@/types'

export async function fetchSeasons(): Promise<SeasonItem[]> {
  const { data } = await apiClient.get<SeasonItem[]>('/seasons')
  return data
}
