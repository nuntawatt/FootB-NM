<template>
  <div class="card p-4 space-y-3">
    <div class="flex items-center justify-between">
      <h2 class="text-sm font-bold" :class="isDark ? 'text-white' : 'text-gray-900'">ฟอร์มการเล่นของทีม</h2>
      <span class="text-[10px] px-2 py-0.5 rounded bg-neutral-700 text-neutral-300">Mock</span>
    </div>

    <div class="flex items-center justify-between gap-2">
      <div
        v-for="(item, idx) in form"
        :key="idx"
        class="flex flex-col items-center gap-2 flex-1"
      >
        <span
          class="px-2.5 py-0.5 rounded-md text-xs font-bold text-white"
          :class="badgeColor(item.result)"
        >
          {{ item.score }}
        </span>
        <img :src="item.opponent_logo" alt="Opponent" class="w-8 h-8 object-contain" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { FormMatch } from '@/types'

defineProps<{
  form: FormMatch[]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function badgeColor(result: FormMatch['result']) {
  if (result === 'win') return 'bg-emerald-600'
  if (result === 'loss') return 'bg-sport-red'
  return 'bg-neutral-500'
}
</script>
