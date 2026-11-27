<template>
  <MainLayout>
    <div v-if="loading" class="card p-4 space-y-3">
      <div v-for="i in 4" :key="i" class="skeleton h-16 w-full" />
    </div>

    <div v-else-if="error" class="card p-6 text-center space-y-3">
      <p class="text-sm text-sport-red font-medium">{{ error }}</p>
      <RouterLink to="/" class="btn-primary inline-block text-xs">กลับหน้าแรก</RouterLink>
    </div>

    <div v-else-if="match" class="space-y-5">
      <MatchHeader
        :match="match"
        :venue="venue"
      />

      <!-- Venue & Broadcast Info -->
      <div v-if="venue" class="card p-4 flex flex-wrap items-center justify-between gap-4 text-xs sm:text-sm">
        <div class="flex flex-wrap items-center gap-4">
          <span>🏟️ <strong>{{ venue.stadium }}</strong> ({{ venue.city }})</span>
          <span :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            ความจุ {{ venue.capacity.toLocaleString() }} ที่นั่ง
          </span>
        </div>
        <div v-if="venue.broadcast.length > 0" class="flex flex-wrap items-center gap-1.5">
          <span :class="isDark ? 'text-neutral-400' : 'text-gray-500'">📺 ถ่ายทอดสด:</span>
          <span
            v-for="ch in venue.broadcast"
            :key="ch"
            class="px-2 py-0.5 rounded text-[11px] font-semibold bg-sport-red/15 text-sport-red"
          >
            {{ ch }}
          </span>
        </div>
      </div>

      <!-- Lineups -->
      <LineupTable
        v-if="lineups"
        :home-team="match.home_team"
        :away-team="match.away_team"
        :home-lineup="lineups.home.starting"
        :away-lineup="lineups.away.starting"
        :home-subs="lineups.home.substitutes"
        :away-subs="lineups.away.substitutes"
        :home-formation="lineups.home.formation"
        :away-formation="lineups.away.formation"
        :home-coach="lineups.home.coach_name"
        :away-coach="lineups.away.coach_name"
      />
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchMatch, fetchMatchLineups, fetchMatchVenue } from '@/api/matches'
import type { Match, MatchLineup, VenueInfo } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import MatchHeader from '@/components/matches/MatchHeader.vue'
import LineupTable from '@/components/matches/LineupTable.vue'

const route = useRoute()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const match = ref<Match | null>(null)
const lineups = ref<MatchLineup | null>(null)
const venue = ref<VenueInfo | null>(null)
const loading = ref(false)
const error = ref('')

async function loadMatch() {
  const id = String(route.params.id)
  loading.value = true
  error.value = ''
  try {
    const [matchRes, lineupRes, venueRes] = await Promise.allSettled([
      fetchMatch(id),
      fetchMatchLineups(id),
      fetchMatchVenue(id),
    ])

    if (matchRes.status === 'fulfilled') {
      match.value = matchRes.value
    } else {
      error.value = 'ไม่พบข้อมูลการแข่งขัน หรือไม่สามารถเชื่อมต่อ Backend ได้'
      return
    }

    lineups.value = lineupRes.status === 'fulfilled' ? lineupRes.value : null
    venue.value = venueRes.status === 'fulfilled' ? venueRes.value : null
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, loadMatch)
onMounted(loadMatch)
</script>
