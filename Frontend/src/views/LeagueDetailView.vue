<template>
  <MainLayout>
    <div v-if="loading" class="card p-4 space-y-3">
      <div v-for="i in 4" :key="i" class="skeleton h-12 w-full" />
    </div>

    <div v-else-if="error" class="card p-6 text-center space-y-3">
      <p class="text-sm text-sport-red font-medium">{{ error }}</p>
      <RouterLink to="/" class="btn-primary inline-block text-xs">กลับหน้าแรก</RouterLink>
    </div>

    <div v-else-if="league" class="space-y-6">
      <!-- League Header Card (GET /leagues/:id) -->
      <div class="card p-6 flex flex-wrap items-center justify-between gap-4 border-t-4 border-t-sport-red">
        <div class="flex items-center gap-5">
          <img :src="league.logo_url" :alt="league.name" class="w-16 h-16 sm:w-20 sm:h-20 object-contain shrink-0" />
          <div class="space-y-1">
            <div class="flex flex-wrap items-center gap-2">
              <h1 class="text-xl sm:text-2xl font-black">{{ league.name }}</h1>
              <span v-if="league.short_name" class="px-2 py-0.5 text-xs font-extrabold rounded bg-sport-red text-white">
                {{ league.short_name }}
              </span>
            </div>
            <p class="text-sm" :class="isDark ? 'text-neutral-300' : 'text-gray-600'">
              {{ league.country || 'Thailand' }} • {{ league.type === 'cup' ? 'บอลถ้วย (Cup)' : 'ลีก (League)' }}
            </p>
            <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              ลำดับการแสดงผล: #{{ league.sort_order }}
            </p>
          </div>
        </div>

        <RouterLink to="/standings" class="btn-primary text-xs">
          ดูตารางคะแนน →
        </RouterLink>
      </div>

      <!-- Teams in League -->
      <div class="card p-5 space-y-4">
        <h2 class="text-base font-extrabold flex items-center gap-2">
          <span class="w-1.5 h-4 bg-sport-red rounded-full" />
          <span>สโมสรในรายการ ({{ teams.length }})</span>
        </h2>

        <div v-if="teams.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
          <RouterLink
            v-for="t in teams"
            :key="t.id"
            :to="`/teams/${t.id}`"
            class="flex items-center gap-2.5 p-2.5 rounded-lg border transition-colors"
            :class="isDark ? 'border-sport-border hover:bg-sport-hover hover:border-sport-red' : 'border-gray-200 hover:bg-gray-50 hover:border-sport-red'"
          >
            <img :src="t.logo_url" :alt="t.name" class="w-7 h-7 object-contain shrink-0" />
            <span class="text-xs font-semibold truncate">{{ t.name }}</span>
          </RouterLink>
        </div>
        <p v-else class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          ยังไม่มีข้อมูลทีมในลีกนี้
        </p>
      </div>

      <!-- Recent Matches -->
      <div class="card overflow-hidden">
        <div class="px-5 py-3.5 border-b flex items-center justify-between" :class="isDark ? 'border-sport-border' : 'border-gray-200'">
          <h2 class="text-base font-extrabold flex items-center gap-2">
            <span class="w-1.5 h-4 bg-sport-red rounded-full" />
            <span>แมตช์ล่าสุดในรายการ</span>
          </h2>
          <RouterLink to="/fixtures" class="text-xs font-semibold text-sport-red hover:underline">
            ดูทั้งหมด →
          </RouterLink>
        </div>

        <FixtureGroup v-if="matches.length > 0" :matches="matches" />
        <p v-else class="p-6 text-center text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          ยังไม่มีข้อมูลแมตช์ในรายการนี้
        </p>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchLeague } from '@/api/leagues'
import { fetchTeams } from '@/api/teams'
import { fetchMatches } from '@/api/matches'
import type { League, Match, Team } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import FixtureGroup from '@/components/fixtures/FixtureGroup.vue'

const route = useRoute()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const league = ref<League | null>(null)
const teams = ref<Team[]>([])
const matches = ref<Match[]>([])
const loading = ref(false)
const error = ref('')

async function loadData() {
  const id = String(route.params.id)
  loading.value = true
  error.value = ''
  try {
    const [lgRes, tmRes, mtRes] = await Promise.allSettled([
      fetchLeague(id),
      fetchTeams(id),
      fetchMatches({ league_id: id, limit: 15 }),
    ])
    if (lgRes.status === 'fulfilled') {
      league.value = lgRes.value
    } else {
      error.value = 'ไม่พบข้อมูลลีกหรือรายการแข่งขัน'
      return
    }
    teams.value = tmRes.status === 'fulfilled' ? tmRes.value : []
    matches.value = mtRes.status === 'fulfilled' ? mtRes.value : []
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, loadData)
onMounted(loadData)
</script>
