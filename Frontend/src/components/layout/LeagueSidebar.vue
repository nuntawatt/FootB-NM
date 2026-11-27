<template>
  <aside class="w-full lg:w-64 shrink-0 space-y-4">
    <!-- League List Card -->
    <div class="card p-3 space-y-1.5">
      <button
        v-for="league in leagues"
        :key="league.id"
        type="button"
        @click="leagueStore.selectLeague(league)"
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-all duration-150 border"
        :class="selectedLeague.id === league.id
          ? 'bg-sport-red/15 border-sport-red text-sport-red font-bold shadow-sm'
          : (isDark
              ? 'border-transparent text-neutral-200 hover:bg-sport-hover'
              : 'border-transparent text-gray-700 hover:bg-gray-100')"
      >
        <div class="w-9 h-9 rounded-lg flex items-center justify-center shrink-0 overflow-hidden bg-white/10 p-1">
          <img :src="league.logo_url" :alt="league.name" class="w-full h-full object-contain" />
        </div>
        <span class="text-sm truncate flex-1">{{ league.name }}</span>
        <RouterLink
          v-if="resolvedId(league.name)"
          :to="`/leagues/${resolvedId(league.name)}`"
          @click.stop
          title="ดูรายละเอียดลีก"
          class="text-[10px] px-1.5 py-0.5 rounded bg-sport-red/15 text-sport-red hover:bg-sport-red hover:text-white transition-colors shrink-0"
        >
          {{ resolvedShortName(league.name) || 'ข้อมูล' }}
        </RouterLink>
      </button>
    </div>

    <!-- Bottom Info / Placeholder Card matching design -->
    <div
      class="hidden lg:block card p-4 text-xs space-y-2"
      :class="isDark ? 'text-neutral-400' : 'text-gray-500'"
    >
      <div class="flex items-center gap-2 font-semibold text-sport-red uppercase tracking-wider">
        <span class="w-2 h-2 rounded-full bg-sport-red animate-pulse" />
        <span>ฤดูกาล {{ seasonStore.selectedSeason }}/{{ seasonStore.selectedSeason + 1 }}</span>
      </div>
      <p class="leading-relaxed">
        ติดตามผลบอลสด โปรแกรมแข่งขัน และตารางคะแนนฟุตบอลไทยครบทุกรายการ
      </p>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useLeagueStore } from '@/stores/league'
import { useSeasonStore } from '@/stores/season'
import { useThemeStore } from '@/stores/theme'
import { fetchLeagues } from '@/api/leagues'
import type { League } from '@/types'

const leagueStore = useLeagueStore()
const seasonStore = useSeasonStore()
const themeStore = useThemeStore()

const { leagues, selectedLeague } = storeToRefs(leagueStore)
const { isDark } = storeToRefs(themeStore)

const apiLeagues = ref<League[]>([])

onMounted(async () => {
  try {
    const list = await fetchLeagues()
    apiLeagues.value = list.sort((a, b) => a.sort_order - b.sort_order)
  } catch {
    // ignore if offline
  }
})

function findApiLeague(name: string): League | undefined {
  const target = name.toLowerCase()
  return apiLeagues.value.find((l) => l.name.toLowerCase() === target || l.name.toLowerCase().includes(target))
}

function resolvedId(name: string) {
  return findApiLeague(name)?.id
}

function resolvedShortName(name: string) {
  return findApiLeague(name)?.short_name
}
</script>
