<template>
  <div class="card overflow-hidden border-t-2 border-t-sport-red">
    <div
      class="px-4 py-3 border-b flex items-center justify-between"
      :class="isDark ? 'border-sport-border' : 'border-gray-200'"
    >
      <div class="flex items-center gap-2">
        <span class="w-1.5 h-4 bg-sport-red rounded-full" />
        <h3 class="font-extrabold text-sm">ดาวซัลโวไทยลีก</h3>
      </div>
      <span class="text-[10px] px-2 py-0.5 rounded bg-neutral-700/70 text-neutral-300">Mock</span>
    </div>

    <div class="divide-y" :class="isDark ? 'divide-sport-border/60' : 'divide-gray-100'">
      <div
        v-for="scorer in scorers"
        :key="scorer.rank"
        class="flex items-center justify-between px-4 py-2.5 text-xs sm:text-sm"
      >
        <div class="flex items-center gap-2.5 min-w-0">
          <span
            class="w-5 text-center font-extrabold tabular-nums"
            :class="scorer.rank === 1 ? 'text-sport-red' : (isDark ? 'text-neutral-400' : 'text-gray-500')"
          >
            {{ scorer.rank }}
          </span>
          <img :src="scorer.team_logo" :alt="scorer.team" class="w-5 h-5 object-contain shrink-0" />
          <div class="min-w-0">
            <p class="font-semibold truncate">{{ scorer.name }}</p>
            <p class="text-[11px] truncate" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              {{ scorer.team }}
            </p>
          </div>
        </div>

        <div class="shrink-0 flex items-center gap-1">
          <span class="text-base font-black text-sport-red tabular-nums">{{ scorer.goals }}</span>
          <span class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">ประตู</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { mockTopScorers } from '@/mock/thaileague'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
const scorers = mockTopScorers
</script>
