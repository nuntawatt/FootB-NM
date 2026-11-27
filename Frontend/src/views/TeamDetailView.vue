<template>
  <MainLayout>
    <div v-if="loading" class="card p-4 space-y-3">
      <div v-for="i in 4" :key="i" class="skeleton h-16 w-full" />
    </div>

    <div v-else-if="error" class="card p-6 text-center space-y-3">
      <p class="text-sm text-sport-red font-medium">{{ error }}</p>
      <RouterLink to="/" class="btn-primary inline-block text-xs">กลับหน้าแรก</RouterLink>
    </div>

    <div v-else-if="team" class="space-y-5">
      <!-- Team Header + Coach -->
      <TeamHeader
        :team="team"
        :coach-name="coach?.name || 'ไม่ระบุ'"
        :coach-image="coach?.photo_url || ''"
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
import { useRoute, RouterLink } from 'vue-router'
import { fetchTeam, fetchTeamCoach, fetchTeamForm, fetchTeamSquad } from '@/api/teams'
import { fetchMatches } from '@/api/matches'
import type { ApiFormMatch, Coach, Match, Player, Team } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import TeamHeader from '@/components/teams/TeamHeader.vue'
import TeamForm from '@/components/teams/TeamForm.vue'
import NextMatch from '@/components/teams/NextMatch.vue'
import SquadList from '@/components/teams/SquadList.vue'

const route = useRoute()

const team = ref<Team | null>(null)
const coach = ref<Coach | null>(null)
const form = ref<ApiFormMatch[]>([])
const players = ref<Player[]>([])
const nextMatch = ref<Match | null>(null)
const loading = ref(false)
const error = ref('')

async function loadTeamData() {
  const id = String(route.params.id)
  loading.value = true
  error.value = ''
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
      team.value = null
      error.value = 'ไม่พบข้อมูลทีม หรือไม่สามารถเชื่อมต่อ Backend ได้'
      return
    }

    coach.value = coachRes.status === 'fulfilled' ? coachRes.value : null
    form.value = formRes.status === 'fulfilled' && Array.isArray(formRes.value) ? formRes.value : []
    players.value = squadRes.status === 'fulfilled' && Array.isArray(squadRes.value) ? squadRes.value : []

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
