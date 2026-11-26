<template>
  <MainLayout>
    <div class="space-y-4">
      <DateNavigator :label="`ฤดูกาล ${selectedSeason}/${selectedSeason + 1}`" />

      <!-- League Header -->
      <div
        class="card px-4 py-3 flex items-center justify-between border-2 border-white/80"
        :class="isDark ? 'bg-sport-card' : 'bg-white !border-sport-red'"
      >
        <div class="flex items-center gap-3">
          <img :src="selectedLeague.logo_url" :alt="selectedLeague.name" class="w-8 h-8 object-contain" />
          <h1 class="text-base sm:text-lg font-extrabold">{{ selectedLeague.name }}</h1>
        </div>
        <span class="text-xs px-2.5 py-1 rounded bg-amber-500/20 text-amber-400 font-semibold">
          Mock Data (รอ API ตารางคะแนน)
        </span>
      </div>

      <StandingsTable v-if="rows.length > 0" :rows="rows" />

      <div v-else class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ยังไม่มีข้อมูลตารางคะแนนสำหรับ {{ selectedLeague.name }}
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useLeagueStore } from '@/stores/league'
import { useSeasonStore } from '@/stores/season'
import { useThemeStore } from '@/stores/theme'
import { mockStandings } from '@/mock/standings'
import MainLayout from '@/components/layout/MainLayout.vue'
import DateNavigator from '@/components/fixtures/DateNavigator.vue'
import StandingsTable from '@/components/standings/StandingsTable.vue'

const leagueStore = useLeagueStore()
const seasonStore = useSeasonStore()
const themeStore = useThemeStore()

const { selectedLeague } = storeToRefs(leagueStore)
const { selectedSeason } = storeToRefs(seasonStore)
const { isDark } = storeToRefs(themeStore)

const rows = computed(() => mockStandings[selectedLeague.value.id] ?? mockStandings['thai-league-1'])
</script>
