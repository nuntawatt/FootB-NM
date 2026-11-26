<template>
  <MainLayout>
    <div v-if="loading" class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
      กำลังโหลดข้อมูลทีม...
    </div>

    <div v-else-if="team" class="space-y-5">
      <!-- Team Header + Coach -->
      <TeamHeader
        :team="team"
        :coach-name="mockDetail.coach_name"
        :coach-image="mockDetail.coach_image"
      />

      <!-- Form & Next Match -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <TeamForm :form="mockDetail.form" />
        <NextMatch :next-match="nextMatch" />
      </div>

      <!-- Squad List -->
      <SquadList :squad="mockDetail.squad" />
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchTeam } from '@/api/teams'
import { fetchMatches } from '@/api/matches'
import { mockTeamDetail } from '@/mock/teamDetail'
import { mockFavorites } from '@/mock/favorites'
import type { Match, Team } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import TeamHeader from '@/components/teams/TeamHeader.vue'
import TeamForm from '@/components/teams/TeamForm.vue'
import NextMatch from '@/components/teams/NextMatch.vue'
import SquadList from '@/components/teams/SquadList.vue'

const route = useRoute()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const team = ref<Team | null>(null)
const nextMatch = ref<Match | null>(null)
const loading = ref(false)
const mockDetail = mockTeamDetail

async function loadTeamData() {
  const id = String(route.params.id)
  loading.value = true
  try {
    team.value = await fetchTeam(id)
    const allScheduled = await fetchMatches({ status: 'scheduled', limit: 50 })
    nextMatch.value =
      allScheduled.find((m) => m.home_team_id === id || m.away_team_id === id) ?? null
  } catch {
    // Fallback to mock if UUID isn't in DB or backend is offline
    const fallback = mockFavorites.find((f) => f.team.id === id)?.team ?? mockFavorites[0].team
    team.value = fallback
    nextMatch.value = null
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, loadTeamData)
onMounted(loadTeamData)
</script>
