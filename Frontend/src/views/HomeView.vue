<template>
  <MainLayout>
    <template #hero>
      <HeroHighlight />
    </template>

    <div class="space-y-4">
      <!-- Highlight / Date Bar -->
      <DateNavigator
        label="Highlight"
        :selected-date="selectedDate"
        show-date-picker
        @prev="shiftDate(-1)"
        @next="shiftDate(1)"
        @update:date="onDatePick"
      />

      <!-- Selected League Header -->
      <div
        class="card px-4 py-3 flex items-center justify-between border-2 border-white/80"
        :class="isDark ? 'bg-sport-card' : 'bg-white !border-sport-red'"
      >
        <div class="flex items-center gap-3">
          <img :src="selectedLeague.logo_url" :alt="selectedLeague.name" class="w-8 h-8 object-contain" />
          <h2 class="text-base sm:text-lg font-extrabold tracking-wide">{{ selectedLeague.name }}</h2>
        </div>

        <RouterLink to="/fixtures" class="text-xs font-semibold text-sport-red hover:underline">
          ดูทั้งหมด →
        </RouterLink>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        กำลังโหลดข้อมูลการแข่งขัน...
      </div>

      <!-- Error -->
      <div v-else-if="error" class="card p-6 text-center space-y-3">
        <p class="text-sm text-sport-red font-medium">{{ error }}</p>
        <button type="button" class="btn-primary text-xs" @click="loadMatches">ลองใหม่อีกครั้ง</button>
      </div>

      <!-- Match List -->
      <div v-else-if="matches.length > 0" class="card overflow-hidden">
        <FixtureGroup :matches="matches" />
      </div>

      <!-- Empty -->
      <div v-else class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ไม่พบข้อมูลการแข่งขันสำหรับ {{ selectedLeague.name }}
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useLeagueStore } from '@/stores/league'
import { useSeasonStore } from '@/stores/season'
import { useThemeStore } from '@/stores/theme'
import { fetchLeagues } from '@/api/leagues'
import { fetchMatches } from '@/api/matches'
import type { League, Match } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import HeroHighlight from '@/components/home/HeroHighlight.vue'
import DateNavigator from '@/components/fixtures/DateNavigator.vue'
import FixtureGroup from '@/components/fixtures/FixtureGroup.vue'

const leagueStore = useLeagueStore()
const seasonStore = useSeasonStore()
const themeStore = useThemeStore()

const { selectedLeague } = storeToRefs(leagueStore)
const { selectedSeason } = storeToRefs(seasonStore)
const { isDark } = storeToRefs(themeStore)

const apiLeagues = ref<League[]>([])
const matches = ref<Match[]>([])
const loading = ref(false)
const error = ref('')
const selectedDate = ref('')

async function resolveLeagueId(): Promise<string | undefined> {
  if (apiLeagues.value.length === 0) {
    try {
      apiLeagues.value = await fetchLeagues()
    } catch {
      return undefined
    }
  }
  const target = selectedLeague.value.name.toLowerCase()
  const found = apiLeagues.value.find(
    (l) => l.name.toLowerCase() === target || l.name.toLowerCase().includes(target)
  )
  return found?.id
}

async function loadMatches() {
  loading.value = true
  error.value = ''
  try {
    const leagueId = await resolveLeagueId()
    matches.value = await fetchMatches({
      league_id: leagueId,
      season: selectedDate.value ? undefined : selectedSeason.value,
      date: selectedDate.value || undefined,
      limit: 20,
    })
  } catch {
    error.value = 'ไม่สามารถเชื่อมต่อ Backend API ได้ กรุณาตรวจสอบว่า Backend กำลังทำงานอยู่ (http://localhost:3000)'
  } finally {
    loading.value = false
  }
}

function shiftDate(days: number) {
  const base = selectedDate.value ? new Date(selectedDate.value) : new Date()
  base.setDate(base.getDate() + days)
  selectedDate.value = base.toISOString().slice(0, 10)
  loadMatches()
}

function onDatePick(val: string) {
  selectedDate.value = val
  loadMatches()
}

watch([selectedLeague, selectedSeason], () => {
  loadMatches()
})

onMounted(async () => {
  if (seasonStore.seasons.length === 0) {
    await seasonStore.loadSeasons()
  }
  await loadMatches()
})
</script>
