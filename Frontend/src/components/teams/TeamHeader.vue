<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <!-- Team Info Card -->
    <div class="md:col-span-2 card p-5 flex items-center gap-5">
      <img
        :src="team.logo_url"
        :alt="team.name"
        class="w-20 h-20 sm:w-24 sm:h-24 object-contain shrink-0"
      />
      <div class="space-y-1 min-w-0">
        <h1 class="text-xl sm:text-2xl font-extrabold truncate" :class="isDark ? 'text-white' : 'text-gray-900'">
          {{ team.name }}
        </h1>
        <p class="text-sm" :class="isDark ? 'text-neutral-300' : 'text-gray-600'">
          {{ team.league?.name || 'Thai League 1' }}
        </p>
        <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          {{ team.league?.country || 'Thailand' }}
        </p>
        <span class="inline-block mt-1 px-2.5 py-0.5 text-xs font-bold rounded bg-sport-red text-white">
          #5 ตารางคะแนน
        </span>
      </div>
    </div>

    <!-- Coach Card (Mock) -->
    <div class="card p-4 flex flex-col items-center justify-center text-center">
      <img
        :src="coachImage"
        :alt="coachName"
        class="w-20 h-20 rounded-xl object-cover bg-neutral-800 mb-2 border border-sport-border"
        @error="onImgError"
      />
      <p class="font-bold text-sm" :class="isDark ? 'text-white' : 'text-gray-900'">{{ coachName }}</p>
      <p class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">ผู้จัดการทีม (Mock)</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { Team } from '@/types'

defineProps<{
  team: Team
  coachName: string
  coachImage: string
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function onImgError(e: Event) {
  (e.target as HTMLImageElement).src = 'https://via.placeholder.com/80?text=Coach'
}
</script>
