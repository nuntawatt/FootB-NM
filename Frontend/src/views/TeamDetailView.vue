<template>
  <MainLayout>
    <div v-if="loading" class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
      กำลังโหลดข้อมูลทีม...
    </div>

    <div v-else-if="team" class="space-y-5">
      <!-- Team Header + Coach -->
      <TeamHeader
        :team="team"
        :coach-name="coach?.name || mockDetail.coach_name"
        :coach-image="coach?.photo_url || mockDetail.coach_image"
        :coach-nationality="coach?.nationality"
      />

      <!-- Form & Next Match -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <TeamForm :form="form" />
        <NextMatch :next-match="nextMatch" />
      </div>

      <!-- Squad List -->
      <SquadList :players="players" />
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchTeam, fetchTeamCoach, fetchTeamForm, fetchTeamSquad } from '@/api/teams'
import { fetchMatches } from '@/api/matches'
import { mockTeamDetail } from '@/mock/teamDetail'
import { mockFavorites } from '@/mock/favorites'
import type { ApiFormMatch, Coach, Match, Player, Team } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import TeamHeader from '@/components/teams/TeamHeader.vue'
import TeamForm from '@/components/teams/TeamForm.vue'
import NextMatch from '@/components/teams/NextMatch.vue'
import SquadList from '@/components/teams/SquadList.vue'

const route = useRoute()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const team = ref<Team | null>(null)
const coach = ref<Coach | null>(null)
const form = ref<ApiFormMatch[]>([])
const players = ref<Player[]>([])
const nextMatch = ref<Match | null>(null)
const loading = ref(false)
const mockDetail = mockTeamDetail

async function loadTeamData() {
  const id = String(route.params.id)
  loading.value = true
  try {
    const [teamRes, coachRes, formRes, squadRes, matchesRes] = await Promise.allSettled([
      fetchTeam(id),
      fetchTeamCoach(id),
      fetchTeamForm(id),
      fetchTeamSquad(id),
      fetchMatches({ status: 'scheduled', limit: 50 }),
    ])

    if (teamRes.status === 'fulfilled') {
      team.value = teamRes.value
    } else {
      team.value = mockFavorites.find((f) => f.team.id === id)?.team ?? mockFavorites[0].team
    }

    coach.value = coachRes.status === 'fulfilled' ? coachRes.value : null
    form.value = formRes.status === 'fulfilled' && Array.isArray(formRes.value) ? formRes.value : []

    if (squadRes.status === 'fulfilled' && squadRes.value.length > 0) {
      players.value = squadRes.value
    } else {
      players.value = mockDetail.squad.map((s, idx) => ({
        id: String(idx),
        team_id: id,
        name: s.name,
        number: s.number,
        position: s.position,
        age: null,
        nationality: 'Thailand',
        photo_url: '',
      }))
    }

    if (matchesRes.status === 'fulfilled') {
      nextMatch.value =
        matchesRes.value.find((m) => m.home_team_id === id || m.away_team_id === id) ?? null
    }
  } finally {
    loading.value = false
  }
}

watch(() => route.params.id, loadTeamData)
onMounted(loadTeamData)
</script>
