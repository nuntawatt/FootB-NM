<template>
  <div
    v-if="tickerMatches.length > 0"
    class="border-b text-xs transition-colors duration-200"
    :class="isDark ? 'bg-[#080808] border-sport-border text-neutral-200' : 'bg-gray-50 border-gray-200 text-gray-800'"
  >
    <div class="max-w-[1400px] mx-auto px-4 py-1.5 flex items-center gap-3 overflow-x-auto no-scrollbar">
      <span class="shrink-0 inline-flex items-center gap-1.5 px-2 py-0.5 rounded bg-sport-red/15 text-sport-red font-extrabold uppercase tracking-wider text-[10px]">
        <span class="w-1.5 h-1.5 rounded-full bg-sport-red animate-ping" />
        SCOREBOARD
      </span>

      <div class="flex items-center gap-3">
        <RouterLink
          v-for="m in tickerMatches"
          :key="m.id"
          :to="`/matches/${m.id}`"
          class="shrink-0 flex items-center gap-2 px-2.5 py-1 rounded-md border transition-colors"
          :class="m.status === 'live'
            ? 'border-sport-red/70 bg-sport-red/10 hover:bg-sport-red/20'
            : (isDark ? 'border-sport-border/70 hover:bg-sport-card' : 'border-gray-200 hover:bg-white')"
        >
          <span class="font-semibold truncate max-w-[90px]">{{ m.home_team?.short_name || m.home_team?.name }}</span>
          <span
            class="font-black tabular-nums px-1.5 py-0.5 rounded text-[11px]"
            :class="m.status === 'live' ? 'bg-sport-red text-white' : (isDark ? 'bg-sport-card text-white' : 'bg-gray-200 text-gray-900')"
          >
            <template v-if="m.status === 'live' || m.status === 'finished'">
              {{ m.home_score ?? 0 }}-{{ m.away_score ?? 0 }}
            </template>
            <template v-else>
              {{ formatTime(m.match_date) }}
            </template>
          </span>
          <span class="font-semibold truncate max-w-[90px]">{{ m.away_team?.short_name || m.away_team?.name }}</span>
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchMatches } from '@/api/matches'
import type { Match } from '@/types'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const tickerMatches = ref<Match[]>([])

onMounted(async () => {
  try {
    const data = await fetchMatches({ limit: 8 })
    tickerMatches.value = data
  } catch {
    tickerMatches.value = []
  }
})

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit', hour12: false })
}
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
