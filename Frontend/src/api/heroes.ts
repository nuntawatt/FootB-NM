import apiClient from './client'
import type { HeroHighlight } from '@/types'

export async function fetchHeroes(): Promise<HeroHighlight[]> {
  const { data } = await apiClient.get<HeroHighlight[]>('/heroes')
  return data
}
