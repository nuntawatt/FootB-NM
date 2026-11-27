import apiClient from './client'
import type { HeroHighlight, NewsArticle } from '@/types'

export function mapHeroToNewsArticle(item: HeroHighlight): NewsArticle {
  return {
    id: item.id,
    title: item.title,
    summary: item.subtitle || '',
    content: item.subtitle || '',
    image_url: item.image_url,
    category: item.tag || 'NEWS',
    published_at: '',
    redirect_url: item.redirect_url || (item.match_id ? `/matches/${item.match_id}` : undefined),
  }
}

export async function fetchNews(limit = 12): Promise<NewsArticle[]> {
  const { data } = await apiClient.get<HeroHighlight[]>('/news', { params: { limit } })
  return data.map(mapHeroToNewsArticle)
}
