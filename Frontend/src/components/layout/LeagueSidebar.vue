<template>
  <aside class="w-full lg:w-64 shrink-0 space-y-4">
    <!-- League List Card -->
    <div class="card p-3 space-y-1.5">
      <button
        v-for="league in leagues"
        :key="league.id"
        type="button"
        @click="leagueStore.selectLeague(league)"
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-all duration-150 border"
        :class="selectedLeague.id === league.id
          ? 'bg-sport-red/15 border-sport-red text-sport-red font-bold shadow-sm'
          : (isDark
              ? 'border-transparent text-neutral-200 hover:bg-sport-hover'
              : 'border-transparent text-gray-700 hover:bg-gray-100')"
      >
        <div class="w-9 h-9 rounded-lg flex items-center justify-center shrink-0 overflow-hidden bg-white/10 p-1">
          <img :src="league.logo_url" :alt="league.name" class="w-full h-full object-contain" />
        </div>
        <span class="text-sm truncate">{{ league.name }}</span>
      </button>
    </div>

    <!-- Bottom Info / Placeholder Card matching design -->
    <div
      class="hidden lg:block card p-4 text-xs space-y-2"
      :class="isDark ? 'text-neutral-400' : 'text-gray-500'"
    >
      <div class="flex items-center gap-2 font-semibold text-sport-red uppercase tracking-wider">
        <span class="w-2 h-2 rounded-full bg-sport-red animate-pulse" />
        <span>ฤดูกาล {{ seasonStore.selectedSeason }}/{{ seasonStore.selectedSeason + 1 }}</span>
      </div>
      <p class="leading-relaxed">
        ติดตามผลบอลสด โปรแกรมแข่งขัน และตารางคะแนนฟุตบอลไทยครบทุกรายการ
      </p>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useLeagueStore } from '@/stores/league'
import { useSeasonStore } from '@/stores/season'
import { useThemeStore } from '@/stores/theme'

const leagueStore = useLeagueStore()
const seasonStore = useSeasonStore()
const themeStore = useThemeStore()

const { leagues, selectedLeague } = storeToRefs(leagueStore)
const { isDark } = storeToRefs(themeStore)
</script>
