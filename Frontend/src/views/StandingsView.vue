<template>
  <MainLayout>
    <div class="space-y-4">
      <DateNavigator :label="`ฤดูกาล ${selectedSeason}/${selectedSeason + 1}`" />

      <!-- League Header -->
      <div
        class="card px-4 py-3 flex flex-wrap items-center justify-between gap-2 border-2 border-white/80"
        :class="isDark ? 'bg-sport-card' : 'bg-white !border-sport-red'"
      >
        <div class="flex items-center gap-3">
          <img :src="selectedLeague.logo_url" :alt="selectedLeague.name" class="w-8 h-8 object-contain" />
          <h1 class="text-base sm:text-lg font-extrabold">{{ selectedLeague.name }}</h1>
        </div>
        <span
          v-if="usingMock"
          class="text-xs px-2.5 py-1 rounded bg-amber-500/20 text-amber-400 font-semibold"
        >
          ข้อมูลตัวอย่าง (Fallback)
        </span>
      </div>

      <div v-if="loading" class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        กำลังคำนวณตารางคะแนน...
      </div>

      <StandingsTable v-else-if="rows.length > 0" :rows="rows" />

      <div v-else class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ยังไม่มีข้อมูลตารางคะแนนสำหรับ {{ selectedLeague.name }}
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useLeagueStore } from '@/stores/league'
import { useSeasonStore } from '@/stores/season'
import { useThemeStore } from '@/stores/theme'
import { fetchLeagues } from '@/api/leagues'
import { fetchStandings } from '@/api/standings'
import { mockStandings } from '@/mock/standings'
import type { League, StandingRow } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import DateNavigator from '@/components/fixtures/DateNavigator.vue'
import StandingsTable from '@/components/standings/StandingsTable.vue'

const leagueStore = useLeagueStore()
const seasonStore = useSeasonStore()
const themeStore = useThemeStore()

const { selectedLeague } = storeToRefs(leagueStore)
const { selectedSeason } = storeToRefs(seasonStore)
const { isDark } = storeToRefs(themeStore)

const apiLeagues = ref<League[]>([])
const rows = ref<StandingRow[]>([])
const loading = ref(false)
const usingMock = ref(false)

async function resolveLeagueId(): Promise<string | undefined> {
  if (apiLeagues.value.length === 0) {
    try {
      apiLeagues.value = await fetchLeagues()
    } catch {
      return undefined
    }
  }
  const target = selectedLeague.value.name.toLowerCase()
  const found = apiLeagues.value.find(
    (l) => l.name.toLowerCase() === target || l.name.toLowerCase().includes(target)
  )
  return found?.id
}

async function loadStandings() {
  loading.value = true
  usingMock.value = false
  try {
    const leagueId = await resolveLeagueId()
    const apiRows = await fetchStandings(selectedSeason.value || 2024, leagueId)
    if (apiRows.length > 0) {
      const total = apiRows.length
      rows.value = apiRows.map((r) => ({
        position: r.rank,
        team: {
          id: r.team_id,
          league_id: leagueId || '',
          name: r.team_name,
          short_name: '',
          logo_url: r.team_logo,
          created_at: '',
          updated_at: '',
        },
        played: r.played,
        won: r.won,
        drawn: r.drawn,
        lost: r.lost,
        goals_for: r.goals_for,
        goals_against: r.goals_against,
        goal_difference: r.goal_diff,
        points: r.points,
        form: r.form,
        zone:
          r.rank <= 2
            ? 'acl'
            : r.rank === 3
              ? 'acl2'
              : total >= 10 && r.rank >= total - 1
                ? 'relegation'
                : undefined,
      }))
    } else {
      rows.value = mockStandings[selectedLeague.value.id] ?? mockStandings['thai-league-1']
      usingMock.value = true
    }
  } catch {
    rows.value = mockStandings[selectedLeague.value.id] ?? mockStandings['thai-league-1']
    usingMock.value = true
  } finally {
    loading.value = false
  }
}

watch([selectedLeague, selectedSeason], loadStandings)

onMounted(async () => {
  if (seasonStore.seasons.length === 0) {
    await seasonStore.loadSeasons()
  }
  await loadStandings()
})
</script>
