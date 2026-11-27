<template>
  <div
    class="border-b transition-colors duration-200"
    :class="isDark ? 'bg-[#111111] border-sport-border' : 'bg-white border-gray-200'"
  >
    <div class="max-w-[1400px] mx-auto px-4 py-2 flex items-center gap-3 overflow-x-auto no-scrollbar">
      <span
        class="text-[11px] font-extrabold uppercase tracking-wider shrink-0 px-2.5 py-0.5 rounded bg-sport-red text-white"
      >
        สโมสร • {{ selectedLeague.name }}
      </span>

      <div v-if="loading" class="flex items-center gap-2">
        <div v-for="i in 8" :key="i" class="skeleton w-7 h-7 rounded-full shrink-0" />
      </div>

      <div v-else-if="clubs.length > 0" class="flex items-center gap-2 sm:gap-3">
        <RouterLink
          v-for="team in clubs"
          :key="team.id"
          :to="`/teams/${team.id}`"
          :title="team.name"
          class="group shrink-0 p-1 rounded-lg transition-all hover:scale-110"
          :class="isDark ? 'hover:bg-sport-hover' : 'hover:bg-gray-100'"
        >
          <img
            :src="team.logo_url"
            :alt="team.name"
            class="w-7 h-7 object-contain opacity-85 group-hover:opacity-100 transition-opacity"
          />
        </RouterLink>
      </div>

      <span v-else class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ยังไม่มีข้อมูลสโมสรในรายการนี้
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { useLeagueStore } from '@/stores/league'
import { fetchLeagues } from '@/api/leagues'
import { fetchTeams } from '@/api/teams'
import { fetchMatches } from '@/api/matches'
import type { League, Team } from '@/types'

const themeStore = useThemeStore()
const leagueStore = useLeagueStore()
const { isDark } = storeToRefs(themeStore)
const { selectedLeague } = storeToRefs(leagueStore)

const apiLeagues = ref<League[]>([])
const clubs = ref<Team[]>([])
const loading = ref(false)

async function resolveLeagueId(name: string): Promise<string | undefined> {
  if (apiLeagues.value.length === 0) {
    try {
      apiLeagues.value = await fetchLeagues()
    } catch {
      return undefined
    }
  }
  const target = name.toLowerCase()
  const found = apiLeagues.value.find(
    (l) => l.name.toLowerCase() === target || l.name.toLowerCase().includes(target)
  )
  return found?.id
}

async function loadClubsForLeague() {
  loading.value = true
  try {
    const leagueId = await resolveLeagueId(selectedLeague.value.name)
    if (leagueId) {
      const teams = await fetchTeams(leagueId)
      if (teams.length > 0) {
        clubs.value = teams
        return
      }

      const matches = await fetchMatches({ league_id: leagueId, limit: 100 })
      const teamMap = new Map<string, Team>()
      for (const m of matches) {
        if (m.home_team) teamMap.set(m.home_team.id, m.home_team)
        if (m.away_team) teamMap.set(m.away_team.id, m.away_team)
      }
      if (teamMap.size > 0) {
        clubs.value = Array.from(teamMap.values()).sort((a, b) => a.name.localeCompare(b.name))
        return
      }
    }

    const allTeams = await fetchTeams()
    clubs.value = allTeams.slice(0, 16)
  } catch {
    clubs.value = []
  } finally {
    loading.value = false
  }
}

watch(selectedLeague, loadClubsForLeague)
onMounted(loadClubsForLeague)
</script>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
