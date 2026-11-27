<template>
  <div class="card overflow-hidden border-t-2 border-t-sport-red">
    <div
      class="px-4 py-3 border-b flex items-center justify-between"
      :class="isDark ? 'border-sport-border' : 'border-gray-200'"
    >
      <div class="flex items-center gap-2">
        <span class="w-1.5 h-4 bg-sport-red rounded-full" />
        <h3 class="font-extrabold text-sm">ตารางคะแนน (Top 5)</h3>
      </div>
      <RouterLink to="/standings" class="text-xs font-semibold text-sport-red hover:underline">
        ดูทั้งหมด →
      </RouterLink>
    </div>

    <div class="divide-y" :class="isDark ? 'divide-sport-border/60' : 'divide-gray-100'">
      <RouterLink
        v-for="row in topFive"
        :key="row.team_id"
        :to="`/teams/${row.team_id}`"
        class="flex items-center justify-between px-4 py-2.5 text-xs sm:text-sm transition-colors"
        :class="isDark ? 'hover:bg-sport-hover/50' : 'hover:bg-gray-50'"
      >
        <div class="flex items-center gap-2.5 min-w-0">
          <span
            class="w-5 text-center font-extrabold tabular-nums"
            :class="row.rank <= 2 ? 'text-sport-red' : (isDark ? 'text-neutral-400' : 'text-gray-500')"
          >
            {{ row.rank }}
          </span>
          <img :src="row.team_logo" :alt="row.team_name" class="w-5 h-5 object-contain shrink-0" />
          <span class="font-semibold truncate">{{ row.team_name }}</span>
        </div>

        <div class="flex items-center gap-4 shrink-0 tabular-nums">
          <span :class="isDark ? 'text-neutral-400' : 'text-gray-500'">{{ row.played }} นัด</span>
          <span class="font-extrabold text-sport-red w-6 text-right">{{ row.points }}</span>
        </div>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { useSeasonStore } from '@/stores/season'
import { fetchStandings } from '@/api/standings'
import { mockStandings } from '@/mock/standings'
import type { ApiStandingsRow } from '@/types'

const themeStore = useThemeStore()
const seasonStore = useSeasonStore()
const { isDark } = storeToRefs(themeStore)

const fallback: ApiStandingsRow[] = mockStandings['thai-league-1'].slice(0, 5).map((r) => ({
  rank: r.position,
  team_id: r.team.id,
  team_name: r.team.name,
  team_logo: r.team.logo_url,
  played: r.played,
  won: r.won,
  drawn: r.drawn,
  lost: r.lost,
  goals_for: r.goals_for,
  goals_against: r.goals_against,
  goal_diff: r.goal_difference,
  points: r.points,
  form: [],
}))

const topFive = ref<ApiStandingsRow[]>(fallback)

onMounted(async () => {
  try {
    const data = await fetchStandings(seasonStore.selectedSeason || 2024)
    if (data.length > 0) {
      topFive.value = data.slice(0, 5)
    }
  } catch {
    // keep fallback
  }
})
</script>
