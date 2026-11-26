<template>
  <div class="card p-6">
    <div class="flex items-center justify-between gap-4">
      <!-- Home -->
      <RouterLink
        :to="`/teams/${match.home_team_id}`"
        class="flex flex-col sm:flex-row items-center justify-end gap-3 flex-1 min-w-0 hover:opacity-80 transition-opacity"
      >
        <span class="text-sm sm:text-lg font-extrabold text-center sm:text-right truncate">
          {{ match.home_team?.name }}
        </span>
        <img :src="match.home_team?.logo_url" :alt="match.home_team?.name" class="w-12 h-12 sm:w-14 sm:h-14 object-contain shrink-0" />
      </RouterLink>

      <!-- Score / Kickoff -->
      <div class="px-3 sm:px-6 text-center shrink-0">
        <div
          v-if="match.status === 'live' || match.status === 'finished'"
          class="text-2xl sm:text-4xl font-black tracking-tight"
          :class="match.status === 'live' ? 'text-sport-red' : (isDark ? 'text-white' : 'text-gray-900')"
        >
          {{ match.home_score ?? 0 }} - {{ match.away_score ?? 0 }}
        </div>
        <div v-else class="text-2xl sm:text-3xl font-black text-sport-red">
          {{ formatTime(match.match_date) }}
        </div>

        <p class="text-xs mt-1" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          {{ formatDate(match.match_date) }}
        </p>
        <span
          v-if="match.round"
          class="inline-block mt-1.5 px-2 py-0.5 text-[11px] rounded bg-sport-red/15 text-sport-red font-semibold"
        >
          {{ match.round }}
        </span>
      </div>

      <!-- Away -->
      <RouterLink
        :to="`/teams/${match.away_team_id}`"
        class="flex flex-col-reverse sm:flex-row items-center justify-start gap-3 flex-1 min-w-0 hover:opacity-80 transition-opacity"
      >
        <img :src="match.away_team?.logo_url" :alt="match.away_team?.name" class="w-12 h-12 sm:w-14 sm:h-14 object-contain shrink-0" />
        <span class="text-sm sm:text-lg font-extrabold text-center sm:text-left truncate">
          {{ match.away_team?.name }}
        </span>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { Match } from '@/types'

defineProps<{
  match: Match
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit', hour12: false })
}
function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('th-TH', { weekday: 'long', day: 'numeric', month: 'long' })
}
</script>
