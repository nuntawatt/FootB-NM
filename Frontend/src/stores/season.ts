import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { SeasonItem } from '@/types'
import { fetchSeasons } from '@/api/seasons'

export const useSeasonStore = defineStore('season', () => {
  const seasons = ref<SeasonItem[]>([])
  const selectedSeason = ref<number>(new Date().getFullYear())
  const loading = ref(false)

  async function loadSeasons() {
    loading.value = true
    try {
      seasons.value = await fetchSeasons()
      // เลือก season ที่มี matches มากสุด เป็น default
      const withMatches = seasons.value.filter(s => s.matches_count > 0)
      if (withMatches.length > 0) {
        selectedSeason.value = withMatches[withMatches.length - 1].season
      }
    } catch (e) {
      console.error('Failed to load seasons:', e)
    } finally {
      loading.value = false
    }
  }

  function setSeason(season: number) {
    selectedSeason.value = season
  }

  return { seasons, selectedSeason, loading, loadSeasons, setSeason }
})
