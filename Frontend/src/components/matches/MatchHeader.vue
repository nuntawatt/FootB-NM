<template>
  <div class="card overflow-hidden border-t-2 border-t-sport-red">
    <!-- Top Meta Bar: League, Season, Round, Status -->
    <div
      class="px-5 py-3 border-b flex flex-wrap items-center justify-between gap-2 text-xs"
      :class="isDark ? 'border-sport-border bg-sport-bg/60' : 'border-gray-200 bg-gray-50'"
    >
      <div class="flex flex-wrap items-center gap-2">
        <img
          v-if="match.league?.logo_url"
          :src="match.league.logo_url"
          :alt="match.league.name"
          class="w-5 h-5 object-contain"
        />
        <span class="font-extrabold">{{ match.league?.name || 'Thai League' }}</span>
        <span :class="isDark ? 'text-neutral-400' : 'text-gray-500'">• ฤดูกาล {{ seasonLabel }}</span>
        <span
          v-if="match.round"
          class="px-2 py-0.5 rounded bg-sport-red/15 text-sport-red font-semibold"
        >
          {{ match.round }}
        </span>
      </div>

      <!-- Status Badge -->
      <span
        class="px-2.5 py-0.5 rounded-full font-bold text-[11px]"
        :class="statusBadge.classes"
      >
        {{ statusBadge.label }}
      </span>
    </div>

    <!-- Main Scoreboard -->
    <div class="p-6">
      <div class="flex items-center justify-between gap-4">
        <!-- Home -->
        <RouterLink
          :to="`/teams/${match.home_team_id}`"
          class="flex flex-col sm:flex-row items-center justify-end gap-3 flex-1 min-w-0 hover:opacity-80 transition-opacity"
        >
          <span class="text-sm sm:text-lg font-extrabold text-center sm:text-right truncate">
            {{ match.home_team?.name }}
          </span>
          <img :src="match.home_team?.logo_url" :alt="match.home_team?.name" class="w-12 h-12 sm:w-16 sm:h-16 object-contain shrink-0" />
        </RouterLink>

        <!-- Score / Kickoff -->
        <div class="px-3 sm:px-6 text-center shrink-0 space-y-1">
          <div
            v-if="match.status === 'live' || match.status === 'finished'"
            class="text-3xl sm:text-5xl font-black tracking-tight tabular-nums"
            :class="match.status === 'live' ? 'text-sport-red' : (isDark ? 'text-white' : 'text-gray-900')"
          >
            {{ match.home_score ?? 0 }} - {{ match.away_score ?? 0 }}
          </div>
          <div v-else class="text-2xl sm:text-4xl font-black text-sport-red tabular-nums">
            {{ formatTime(match.match_date) }} น.
          </div>

          <!-- Penalty -->
          <div v-if="(match.status === 'finished' || match.status === 'live') && match.home_penalty_score !== null && match.away_penalty_score !== null" class="flex flex-wrap items-center justify-center gap-2 text-xs">
            <span class="px-2 py-0.5 rounded bg-neutral-700 text-white font-bold">
              จุดโทษ: {{ match.home_penalty_score }} - {{ match.away_penalty_score }}
            </span>
          </div>

          <!-- Date & Time -->
          <p class="text-xs pt-1" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            📅 {{ formatDate(match.match_date) }} • เวลา {{ formatTime(match.match_date) }} น.
          </p>
          <p v-if="venue" class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            🏟️ {{ venue.stadium }}, {{ venue.city }}
          </p>
        </div>

        <!-- Away -->
        <RouterLink
          :to="`/teams/${match.away_team_id}`"
          class="flex flex-col-reverse sm:flex-row items-center justify-start gap-3 flex-1 min-w-0 hover:opacity-80 transition-opacity"
        >
          <img :src="match.away_team?.logo_url" :alt="match.away_team?.name" class="w-12 h-12 sm:w-16 sm:h-16 object-contain shrink-0" />
          <span class="text-sm sm:text-lg font-extrabold text-center sm:text-left truncate">
            {{ match.away_team?.name }}
          </span>
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { Match, VenueInfo } from '@/types'

const props = defineProps<{
  match: Match
  venue?: VenueInfo | null
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const seasonLabel = computed(() => {
  const d = new Date(props.match.match_date)
  const yr = d.getMonth() >= 6 ? d.getFullYear() : d.getFullYear() - 1
  return `${yr}/${yr + 1}`
})

const statusBadge = computed(() => {
  switch (props.match.status) {
    case 'live':
      return {
        label: props.match.minute ? `กำลังแข่งขัน (${props.match.minute}')` : 'กำลังแข่งขัน',
        classes: 'bg-sport-red text-white animate-pulse',
      }
    case 'finished':
      return { label: 'จบการแข่งขัน', classes: 'bg-emerald-600 text-white' }
    case 'postponed':
      return { label: 'เลื่อนการแข่งขัน', classes: 'bg-amber-500 text-black' }
    default:
      return { label: 'ยังไม่เริ่ม', classes: 'bg-neutral-700 text-neutral-200' }
  }
})

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit', hour12: false })
}
function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('th-TH', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}
</script>
