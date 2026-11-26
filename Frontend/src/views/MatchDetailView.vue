<template>
  <MainLayout>
    <div v-if="loading" class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
      กำลังโหลดข้อมูลแมตช์...
    </div>

    <div v-else-if="error" class="card p-6 text-center space-y-3">
      <p class="text-sm text-sport-red font-medium">{{ error }}</p>
      <RouterLink to="/" class="btn-primary inline-block text-xs">กลับหน้าแรก</RouterLink>
    </div>

    <div v-else-if="match" class="space-y-5">
      <MatchHeader :match="match" />

      <LineupTable
        :home-team="match.home_team"
        :away-team="match.away_team"
        :home-lineup="mockLineup.home"
        :away-lineup="mockLineup.away"
      />
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchMatch } from '@/api/matches'
import { mockLineup } from '@/mock/squad'
import type { Match } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import MatchHeader from '@/components/matches/MatchHeader.vue'
import LineupTable from '@/components/matches/LineupTable.vue'

const route = useRoute()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const match = ref<Match | null>(null)
const loading = ref(false)
const error = ref('')

async function loadMatch() {
  const id = String(route.params.id)
  loading.value = true
  error.value = ''
  try {
    match.value = await fetchMatch(id)
  } catch {
    error.value = 'ไม่พบข้อมูลการแข่งขัน หรือไม่สามารถเชื่อมต่อ Backend ได้'
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, loadMatch)
onMounted(loadMatch)
</script>
