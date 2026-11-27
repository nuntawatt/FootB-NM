<template>
  <div class="card p-5 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <span class="w-1.5 h-4 bg-sport-red rounded-full" />
        <h2 class="text-base font-bold">สถิติการแข่งขัน</h2>
      </div>
      <span class="text-[11px] px-2 py-0.5 rounded bg-neutral-700 text-neutral-300">Mock Data</span>
    </div>

    <div class="space-y-3.5">
      <div v-for="stat in stats" :key="stat.label" class="space-y-1">
        <div class="flex items-center justify-between text-xs sm:text-sm font-bold tabular-nums">
          <span class="text-sport-red">{{ stat.home }}{{ stat.unit || '' }}</span>
          <span class="font-semibold" :class="isDark ? 'text-neutral-300' : 'text-gray-700'">{{ stat.label }}</span>
          <span>{{ stat.away }}{{ stat.unit || '' }}</span>
        </div>

        <div class="h-2 rounded-full overflow-hidden flex gap-0.5" :class="isDark ? 'bg-sport-bg' : 'bg-gray-200'">
          <div
            class="bg-sport-red transition-all duration-300 rounded-l-full"
            :style="{ width: `${percent(stat.home, stat.away)}%` }"
          />
          <div
            class="transition-all duration-300 rounded-r-full flex-1"
            :class="isDark ? 'bg-neutral-500' : 'bg-gray-400'"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { MatchStatItem } from '@/types'

defineProps<{
  stats: MatchStatItem[]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function percent(home: number, away: number) {
  const total = home + away
  if (total === 0) return 50
  return Math.round((home / total) * 100)
}
</script>
