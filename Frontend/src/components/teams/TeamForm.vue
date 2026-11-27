<template>
  <div class="card p-4 space-y-3">
    <div class="flex items-center justify-between">
      <h2 class="text-sm font-bold" :class="isDark ? 'text-white' : 'text-gray-900'">ฟอร์มการเล่น 5 นัดล่าสุด</h2>
    </div>

    <div v-if="form.length > 0" class="grid grid-cols-5 gap-2">
      <RouterLink
        v-for="item in form"
        :key="item.match_id"
        :to="`/matches/${item.match_id}`"
        class="flex flex-col items-center gap-1 p-2 rounded-lg transition-colors text-center"
        :class="isDark ? 'hover:bg-sport-hover' : 'hover:bg-gray-50'"
        :title="`${item.opponent} (${item.date})`"
      >
        <span
          class="px-2 py-0.5 rounded-md text-xs font-extrabold text-white"
          :class="badgeColor(item.result)"
        >
          {{ item.result }}
        </span>
        <span class="text-xs font-bold tabular-nums">{{ item.score }}</span>
        <span
          class="text-[9px] font-bold px-1.5 py-0.2 rounded"
          :class="item.is_home
            ? 'bg-sport-red/15 text-sport-red'
            : (isDark ? 'bg-neutral-700 text-neutral-300' : 'bg-gray-200 text-gray-600')"
        >
          {{ item.is_home ? 'เหย้า' : 'เยือน' }}
        </span>
        <span class="text-[10px] truncate w-full" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          {{ item.opponent }}
        </span>
      </RouterLink>
    </div>

    <div v-else class="text-xs text-center py-4" :class="isDark ? 'text-neutral-500' : 'text-gray-400'">
      ยังไม่มีประวัติการแข่งขันที่จบแล้ว
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { ApiFormMatch } from '@/types'

defineProps<{
  form: ApiFormMatch[]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function badgeColor(result: ApiFormMatch['result']) {
  if (result === 'W') return 'bg-emerald-600'
  if (result === 'L') return 'bg-sport-red'
  return 'bg-neutral-500'
}
</script>
