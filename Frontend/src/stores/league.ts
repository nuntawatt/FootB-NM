import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { SidebarLeague } from '@/types'

// Hardcoded 4 leagues ตามดีไซน์
export const SIDEBAR_LEAGUES: SidebarLeague[] = [
  {
    id: 'thai-league-1',
    name: 'Thai League 1',
    logo_url: 'https://media.api-sports.io/football/leagues/296.png',
    type: 'league',
  },
  {
    id: 'thai-league-2',
    name: 'Thai League 2',
    logo_url: 'https://media.api-sports.io/football/leagues/297.png',
    type: 'league',
  },
  {
    id: 'fa-cup',
    name: 'FA CUP',
    logo_url: 'https://media.api-sports.io/football/leagues/298.png',
    type: 'cup',
  },
  {
    id: 'thai-cup',
    name: 'THAI CUP',
    logo_url: 'https://media.api-sports.io/football/leagues/871.png',
    type: 'cup',
  },
]

export const useLeagueStore = defineStore('league', () => {
  const selectedLeague = ref<SidebarLeague>(SIDEBAR_LEAGUES[0])
  const leagues = ref<SidebarLeague[]>(SIDEBAR_LEAGUES)

  function selectLeague(league: SidebarLeague) {
    selectedLeague.value = league
  }

  return { selectedLeague, leagues, selectLeague }
})
