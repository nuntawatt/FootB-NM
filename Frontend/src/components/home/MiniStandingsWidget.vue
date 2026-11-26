<template>
  <div class="card overflow-hidden border-t-2 border-t-sport-red">
    <div
      class="px-4 py-3 border-b flex items-center justify-between"
      :class="isDark ? 'border-sport-border' : 'border-gray-200'"
    >
      <div class="flex items-center gap-2">
        <span class="w-1.5 h-4 bg-sport-red rounded-full" />
        <h3 class="font-extrabold text-sm">ตารางคะแนน (Top 5)</h3>
      </div>
      <RouterLink to="/standings" class="text-xs font-semibold text-sport-red hover:underline">
        ดูทั้งหมด →
      </RouterLink>
    </div>

    <div class="divide-y" :class="isDark ? 'divide-sport-border/60' : 'divide-gray-100'">
      <div
        v-for="row in topFive"
        :key="row.position"
        class="flex items-center justify-between px-4 py-2.5 text-xs sm:text-sm"
      >
        <div class="flex items-center gap-2.5 min-w-0">
          <span
            class="w-5 text-center font-extrabold tabular-nums"
            :class="row.position <= 2 ? 'text-sport-red' : (isDark ? 'text-neutral-400' : 'text-gray-500')"
          >
            {{ row.position }}
          </span>
          <img :src="row.team.logo_url" :alt="row.team.name" class="w-5 h-5 object-contain shrink-0" />
          <span class="font-semibold truncate">{{ row.team.name }}</span>
        </div>

        <div class="flex items-center gap-4 shrink-0 tabular-nums">
          <span :class="isDark ? 'text-neutral-400' : 'text-gray-500'">{{ row.played }} นัด</span>
          <span class="font-extrabold text-sport-red w-6 text-right">{{ row.points }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { mockStandings } from '@/mock/standings'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const topFive = computed(() => mockStandings['thai-league-1'].slice(0, 5))
</script>
