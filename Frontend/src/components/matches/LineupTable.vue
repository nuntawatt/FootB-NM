<template>
  <div class="card p-5">
    <div class="flex flex-wrap items-center justify-between gap-2 mb-4">
      <h2 class="text-base font-bold">รายชื่อผู้เล่น (Lineups)</h2>
      <div class="flex items-center gap-3 text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        <span v-if="homeFormation">แผน {{ homeTeam?.name }}: <strong class="text-sport-red">{{ homeFormation }}</strong></span>
        <span v-if="awayFormation">แผน {{ awayTeam?.name }}: <strong class="text-sport-red">{{ awayFormation }}</strong></span>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 divide-y md:divide-y-0 md:divide-x" :class="isDark ? 'divide-sport-border' : 'divide-gray-200'">
      <!-- Home Lineup -->
      <div class="space-y-3">
        <div class="flex items-center justify-between pb-2 border-b" :class="isDark ? 'border-sport-border' : 'border-gray-200'">
          <div class="flex items-center gap-2">
            <img :src="homeTeam?.logo_url" :alt="homeTeam?.name" class="w-6 h-6 object-contain" />
            <h3 class="font-bold text-sm">{{ homeTeam?.name }}</h3>
          </div>
          <span v-if="homeCoach" class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            โค้ช: {{ homeCoach }}
          </span>
        </div>
        <div class="space-y-1.5">
          <div
            v-for="(p, idx) in homeLineup"
            :key="idx"
            class="flex items-center justify-between text-sm py-1 px-2 rounded"
            :class="isDark ? 'hover:bg-sport-hover/50' : 'hover:bg-gray-50'"
          >
            <div class="flex items-center gap-3">
              <span class="w-8 font-bold text-xs text-sport-red">{{ p.position }}</span>
              <span class="w-6 text-xs tabular-nums" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">#{{ p.number ?? '-' }}</span>
              <span class="font-medium">{{ p.name }}</span>
            </div>
            <span class="text-xs font-bold px-2 py-0.5 rounded bg-sport-red/15 text-sport-red">
              {{ p.rating ?? '-' }}
            </span>
          </div>
        </div>
      </div>

      <!-- Away Lineup -->
      <div class="space-y-3 pt-4 md:pt-0 md:pl-6">
        <div class="flex items-center justify-between pb-2 border-b" :class="isDark ? 'border-sport-border' : 'border-gray-200'">
          <div class="flex items-center gap-2">
            <img :src="awayTeam?.logo_url" :alt="awayTeam?.name" class="w-6 h-6 object-contain" />
            <h3 class="font-bold text-sm">{{ awayTeam?.name }}</h3>
          </div>
          <span v-if="awayCoach" class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            โค้ช: {{ awayCoach }}
          </span>
        </div>
        <div class="space-y-1.5">
          <div
            v-for="(p, idx) in awayLineup"
            :key="idx"
            class="flex items-center justify-between text-sm py-1 px-2 rounded"
            :class="isDark ? 'hover:bg-sport-hover/50' : 'hover:bg-gray-50'"
          >
            <div class="flex items-center gap-3">
              <span class="w-8 font-bold text-xs text-sport-red">{{ p.position }}</span>
              <span class="w-6 text-xs tabular-nums" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">#{{ p.number ?? '-' }}</span>
              <span class="font-medium">{{ p.name }}</span>
            </div>
            <span class="text-xs font-bold px-2 py-0.5 rounded bg-sport-red/15 text-sport-red">
              {{ p.rating ?? '-' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { LineupPlayer, Team } from '@/types'

defineProps<{
  homeTeam?: Team
  awayTeam?: Team
  homeLineup: LineupPlayer[]
  awayLineup: LineupPlayer[]
  homeFormation?: string
  awayFormation?: string
  homeCoach?: string
  awayCoach?: string
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
</script>
