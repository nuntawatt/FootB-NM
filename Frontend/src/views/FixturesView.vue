<template>
  <MainLayout>
    <div class="space-y-4">
      <!-- Date Navigator & Season/Status Filters -->
      <DateNavigator
        :label="dateRangeLabel"
        :selected-date="selectedDate"
        show-date-picker
        @prev="shiftDate(-7)"
        @next="shiftDate(7)"
        @update:date="onDatePick"
      />

      <!-- Filters bar -->
      <div class="card p-3 flex flex-wrap items-center justify-between gap-3">
        <!-- Status tabs -->
        <div class="flex flex-wrap items-center gap-1.5">
          <button
            v-for="st in statusOptions"
            :key="st.value"
            type="button"
            @click="selectedStatus = st.value; loadMatches()"
            class="px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors"
            :class="selectedStatus === st.value
              ? 'bg-sport-red text-white'
              : (isDark ? 'bg-sport-bg text-neutral-300 hover:bg-sport-hover' : 'bg-gray-100 text-gray-700 hover:bg-gray-200')"
          >
            {{ st.label }}
          </button>
        </div>

        <!-- Season Dropdown -->
        <div class="flex items-center gap-2">
          <span class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">ฤดูกาล:</span>
          <select
            :value="selectedSeason"
            @change="onSeasonChange"
            class="px-3 py-1.5 rounded-lg text-xs font-semibold border outline-none"
            :class="isDark
              ? 'bg-sport-bg border-sport-border text-white'
              : 'bg-gray-50 border-gray-300 text-gray-900'"
          >
            <option v-for="s in seasonStore.seasons" :key="s.season" :value="s.season">
              {{ s.label }} ({{ s.matches_count }})
            </option>
          </select>
          <button
            v-if="selectedDate"
            type="button"
            @click="selectedDate = ''; loadMatches()"
            class="text-xs text-sport-red hover:underline ml-1"
          >
            ล้างวันที่
          </button>
        </div>
      </div>

      <!-- League Banner -->
      <div
        class="card px-4 py-3 flex items-center gap-3 border-2 border-white/80"
        :class="isDark ? 'bg-sport-card' : 'bg-white !border-sport-red'"
      >
        <img :src="selectedLeague.logo_url" :alt="selectedLeague.name" class="w-8 h-8 object-contain" />
        <h1 class="text-base sm:text-lg font-extrabold">{{ selectedLeague.name }}</h1>
      </div>

      <!-- Content -->
      <div v-if="loading" class="card p-4 space-y-3">
        <div v-for="i in 6" :key="i" class="skeleton h-12 w-full" />
      </div>

      <div v-else-if="error" class="card p-6 text-center space-y-3">
        <p class="text-sm text-sport-red font-medium">{{ error }}</p>
        <button type="button" class="btn-primary text-xs" @click="loadMatches">ลองใหม่อีกครั้ง</button>
      </div>

      <div v-else-if="matches.length > 0" class="card overflow-hidden">
        <FixtureGroup :matches="matches" />
      </div>

      <div v-else class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ไม่พบโปรแกรมแข่งขันตามเงื่อนไขที่เลือก
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useLeagueStore } from '@/stores/league'
import { useSeasonStore } from '@/stores/season'
import { useThemeStore } from '@/stores/theme'
import { fetchLeagues } from '@/api/leagues'
import { fetchMatches } from '@/api/matches'
import type { League, Match } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import DateNavigator from '@/components/fixtures/DateNavigator.vue'
import FixtureGroup from '@/components/fixtures/FixtureGroup.vue'

const leagueStore = useLeagueStore()
const seasonStore = useSeasonStore()
const themeStore = useThemeStore()

const { selectedLeague } = storeToRefs(leagueStore)
const { selectedSeason } = storeToRefs(seasonStore)
const { isDark } = storeToRefs(themeStore)

const statusOptions = [
  { label: 'ทั้งหมด', value: '' },
  { label: 'กำลังแข่ง', value: 'live' },
  { label: 'ยังไม่แข่ง', value: 'scheduled' },
  { label: 'จบแล้ว', value: 'finished' },
]

const apiLeagues = ref<League[]>([])
const matches = ref<Match[]>([])
const loading = ref(false)
const error = ref('')
const selectedStatus = ref('')
const selectedDate = ref('')

const dateRangeLabel = computed(() => {
  if (selectedDate.value) {
    return new Date(selectedDate.value).toLocaleDateString('th-TH', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    })
  }
  return `ฤดูกาล ${selectedSeason.value}/${selectedSeason.value + 1}`
})

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
      status: selectedStatus.value || undefined,
      date: selectedDate.value || undefined,
      limit: 50,
    })
  } catch {
    error.value = 'ไม่สามารถเชื่อมต่อ Backend API ได้ กรุณาตรวจสอบการเชื่อมต่อ'
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

function onSeasonChange(e: Event) {
  const yr = Number((e.target as HTMLSelectElement).value)
  seasonStore.setSeason(yr)
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
