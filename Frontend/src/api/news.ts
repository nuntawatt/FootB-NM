import apiClient from './client'
import type { NewsArticle } from '@/types'

export async function fetchNews(category?: string): Promise<NewsArticle[]> {
  const params: Record<string, string> = {}
  if (category) params.category = category
  const { data } = await apiClient.get<NewsArticle[]>('/news', { params })
  return data
}

export async function fetchNewsById(id: string): Promise<NewsArticle> {
  const { data } = await apiClient.get<NewsArticle>(`/news/${id}`)
  return data
}
