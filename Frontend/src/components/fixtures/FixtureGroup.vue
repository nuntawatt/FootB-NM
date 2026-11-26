<template>
  <div class="space-y-4">
    <div v-for="group in groupedMatches" :key="group.dateKey">
      <h3
        class="text-sm font-bold px-4 py-2 border-b"
        :class="isDark ? 'text-neutral-300 border-sport-border bg-sport-bg/40' : 'text-gray-700 border-gray-200 bg-gray-50'"
      >
        {{ group.dateLabel }}
      </h3>
      <div>
        <MatchCard v-for="match in group.matches" :key="match.id" :match="match" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { Match } from '@/types'
import MatchCard from '@/components/home/MatchCard.vue'

const props = defineProps<{
  matches: Match[]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const groupedMatches = computed(() => {
  const groups = new Map<string, { dateKey: string; dateLabel: string; matches: Match[] }>()

  for (const m of props.matches) {
    const d = new Date(m.match_date)
    const key = d.toISOString().slice(0, 10)
    if (!groups.has(key)) {
      const label = d.toLocaleDateString('th-TH', {
        weekday: 'long',
        day: 'numeric',
        month: 'long',
      })
      groups.set(key, { dateKey: key, dateLabel: label, matches: [] })
    }
    groups.get(key)!.matches.push(m)
  }

  return Array.from(groups.values())
})
</script>
