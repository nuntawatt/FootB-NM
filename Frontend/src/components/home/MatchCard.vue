<template>
  <RouterLink
    :to="`/matches/${match.id}`"
    class="group flex items-center justify-between gap-2 sm:gap-4 px-3 sm:px-5 py-3.5 transition-all border-b last:border-b-0 border-l-2"
    :class="[
      match.status === 'live' ? 'border-l-sport-red bg-sport-red/5' : 'border-l-transparent',
      isDark
        ? 'border-b-sport-border/70 hover:bg-sport-hover/70'
        : 'border-b-gray-100 hover:bg-gray-50'
    ]"
  >
    <!-- Home Team -->
    <div class="flex items-center justify-end gap-2 sm:gap-3 flex-1 min-w-0 text-right">
      <span class="text-xs sm:text-sm font-semibold truncate group-hover:text-sport-red transition-colors" :class="isDark ? 'text-white' : 'text-gray-900'">
        {{ match.home_team?.name || 'Home Team' }}
      </span>
      <img
        v-if="match.home_team?.logo_url"
        :src="match.home_team.logo_url"
        :alt="match.home_team.name"
        class="w-7 h-7 sm:w-8 sm:h-8 object-contain shrink-0 group-hover:scale-110 transition-transform"
      />
    </div>

    <!-- Center Score / Time -->
    <div class="w-20 sm:w-24 shrink-0 flex flex-col items-center justify-center text-center">
      <template v-if="match.status === 'live'">
        <div class="text-base sm:text-lg font-black text-sport-red tracking-tight">
          {{ match.home_score ?? 0 }} - {{ match.away_score ?? 0 }}
        </div>
        <span class="text-[11px] font-bold px-1.5 py-0.5 rounded bg-sport-red text-white animate-pulse">
          {{ match.minute ? `${match.minute}'` : 'LIVE' }}
        </span>
      </template>

      <template v-else-if="match.status === 'finished'">
        <div class="text-base sm:text-lg font-extrabold" :class="isDark ? 'text-white' : 'text-gray-900'">
          {{ match.home_score ?? 0 }} - {{ match.away_score ?? 0 }}
        </div>
        <span
          v-if="match.home_penalty_score !== null && match.away_penalty_score !== null"
          class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-neutral-600 text-white"
        >
          Pen ({{ match.home_penalty_score }}-{{ match.away_penalty_score }})
        </span>
        <span v-else class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          จบเกม
        </span>
      </template>

      <template v-else-if="match.status === 'postponed'">
        <span class="text-xs font-bold text-amber-500">เลื่อนแข่ง</span>
      </template>

      <template v-else>
        <div class="text-sm sm:text-base font-extrabold text-sport-red">
          {{ formatTime(match.match_date) }}
        </div>
        <span class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          {{ countdownLabel(match.match_date) }}
        </span>
      </template>
    </div>

    <!-- Away Team -->
    <div class="flex items-center justify-start gap-2 sm:gap-3 flex-1 min-w-0 text-left">
      <img
        v-if="match.away_team?.logo_url"
        :src="match.away_team.logo_url"
        :alt="match.away_team.name"
        class="w-7 h-7 sm:w-8 sm:h-8 object-contain shrink-0 group-hover:scale-110 transition-transform"
      />
      <span class="text-xs sm:text-sm font-semibold truncate group-hover:text-sport-red transition-colors" :class="isDark ? 'text-white' : 'text-gray-900'">
        {{ match.away_team?.name || 'Away Team' }}
      </span>
    </div>
  </RouterLink>
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

function formatTime(iso: string): string {
  const d = new Date(iso)
  return d.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function countdownLabel(iso: string): string {
  const diffDays = Math.ceil((new Date(iso).getTime() - Date.now()) / (1000 * 60 * 60 * 24))
  if (diffDays <= 0) return 'วันนี้'
  return `${diffDays} วัน`
}
</script>
