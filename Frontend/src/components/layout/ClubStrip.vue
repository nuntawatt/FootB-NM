<template>
  <div
    class="border-b transition-colors duration-200"
    :class="isDark ? 'bg-[#111111] border-sport-border' : 'bg-white border-gray-200'"
  >
    <div class="max-w-[1400px] mx-auto px-4 py-2 flex items-center gap-3 overflow-x-auto no-scrollbar">
      <span
        class="text-[11px] font-extrabold uppercase tracking-wider shrink-0 px-2 py-0.5 rounded bg-sport-red text-white"
      >
        สโมสร
      </span>
      <div class="flex items-center gap-2 sm:gap-3">
        <RouterLink
          v-for="team in clubs"
          :key="team.id"
          :to="`/teams/${team.id}`"
          :title="team.name"
          class="group shrink-0 p-1 rounded-lg transition-all hover:scale-110"
          :class="isDark ? 'hover:bg-sport-hover' : 'hover:bg-gray-100'"
        >
          <img
            :src="team.logo_url"
            :alt="team.name"
            class="w-7 h-7 object-contain opacity-85 group-hover:opacity-100 transition-opacity"
          />
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchTeams } from '@/api/teams'
import { mockStandings } from '@/mock/standings'
import type { Team } from '@/types'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const clubs = ref<Team[]>(mockStandings['thai-league-1'].map((r) => r.team))

onMounted(async () => {
  try {
    const apiTeams = await fetchTeams()
    if (apiTeams.length > 0) {
      clubs.value = apiTeams.slice(0, 16)
    }
  } catch {
    // keep mock club crests if backend is offline
  }
})
</script>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
