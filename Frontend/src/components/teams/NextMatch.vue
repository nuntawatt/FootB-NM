<template>
  <div class="card p-4 space-y-3">
    <div class="flex items-center justify-between">
      <h2 class="text-sm font-bold" :class="isDark ? 'text-white' : 'text-gray-900'">แมตช์ถัดไป</h2>
      <span class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        {{ nextMatch?.league?.name || 'Thai League 1' }}
      </span>
    </div>

    <RouterLink
      v-if="nextMatch"
      :to="`/matches/${nextMatch.id}`"
      class="flex items-center justify-between gap-3 py-2 rounded-lg transition-colors"
      :class="isDark ? 'hover:bg-sport-hover' : 'hover:bg-gray-50'"
    >
      <div class="flex flex-col items-center gap-1 flex-1 min-w-0">
        <img :src="nextMatch.home_team?.logo_url" :alt="nextMatch.home_team?.name" class="w-9 h-9 object-contain" />
        <span class="text-xs font-semibold truncate max-w-full">{{ nextMatch.home_team?.name }}</span>
      </div>

      <div class="text-center shrink-0">
        <div class="text-lg font-extrabold text-sport-red">{{ formatTime(nextMatch.match_date) }}</div>
        <div class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          {{ formatDate(nextMatch.match_date) }}
        </div>
      </div>

      <div class="flex flex-col items-center gap-1 flex-1 min-w-0">
        <img :src="nextMatch.away_team?.logo_url" :alt="nextMatch.away_team?.name" class="w-9 h-9 object-contain" />
        <span class="text-xs font-semibold truncate max-w-full">{{ nextMatch.away_team?.name }}</span>
      </div>
    </RouterLink>

    <div v-else class="text-xs text-center py-4" :class="isDark ? 'text-neutral-500' : 'text-gray-400'">
      ยังไม่มีโปรแกรมแข่งขันถัดไป
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { Match } from '@/types'

defineProps<{
  nextMatch: Match | null
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit', hour12: false })
}
function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('th-TH', { day: 'numeric', month: 'short' })
}
</script>
