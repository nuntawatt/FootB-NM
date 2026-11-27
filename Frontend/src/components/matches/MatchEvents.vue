<template>
  <div class="card p-5 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <span class="w-1.5 h-4 bg-sport-red rounded-full" />
        <h2 class="text-base font-bold">เหตุการณ์สำคัญ (ประตู &amp; ใบเหลือง/ใบแดง)</h2>
      </div>
      <span class="text-[11px] px-2 py-0.5 rounded bg-neutral-700 text-neutral-300">Mock Data</span>
    </div>

    <div class="space-y-2">
      <div
        v-for="ev in events"
        :key="ev.id"
        class="flex items-center justify-between gap-3 py-2 px-3 rounded-lg border text-xs sm:text-sm"
        :class="isDark ? 'border-sport-border/60 bg-sport-bg/40' : 'border-gray-100 bg-gray-50'"
      >
        <!-- Home side -->
        <div class="flex-1 flex items-center justify-end gap-2 text-right min-w-0">
          <template v-if="ev.team === 'home'">
            <div class="min-w-0">
              <p class="font-bold truncate">{{ ev.player_name }}</p>
              <p v-if="ev.assist_name" class="text-[11px] truncate" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
                แอสซิสต์: {{ ev.assist_name }}
              </p>
            </div>
            <span class="text-base shrink-0">{{ eventIcon(ev.type) }}</span>
          </template>
        </div>

        <!-- Minute -->
        <span class="px-2.5 py-1 rounded-full font-extrabold text-xs bg-sport-red text-white shrink-0 tabular-nums">
          {{ ev.minute }}'
        </span>

        <!-- Away side -->
        <div class="flex-1 flex items-center justify-start gap-2 text-left min-w-0">
          <template v-if="ev.team === 'away'">
            <span class="text-base shrink-0">{{ eventIcon(ev.type) }}</span>
            <div class="min-w-0">
              <p class="font-bold truncate">{{ ev.player_name }}</p>
              <p v-if="ev.assist_name" class="text-[11px] truncate" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
                แอสซิสต์: {{ ev.assist_name }}
              </p>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { MatchEvent } from '@/types'

defineProps<{
  events: MatchEvent[]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function eventIcon(type: MatchEvent['type']) {
  if (type === 'goal') return '⚽'
  if (type === 'yellow_card') return '🟨'
  return '🟥'
}
</script>
