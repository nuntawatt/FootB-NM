<template>
  <div class="card overflow-hidden">
    <div class="bg-white text-gray-900 px-4 py-3 text-center font-extrabold text-base border-b border-gray-200 flex items-center justify-center gap-2">
      <span>11 ตัวจริงล่าสุด</span>
      <span class="text-[10px] font-semibold px-2 py-0.5 rounded bg-gray-200 text-gray-600">Mock</span>
    </div>

    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b text-xs" :class="isDark ? 'border-sport-border text-neutral-400' : 'border-gray-200 text-gray-500'">
            <th class="py-3 px-4 text-left">ตำแหน่ง</th>
            <th class="py-3 px-3 text-center">หมายเลข</th>
            <th class="py-3 px-4 text-left">ชื่อผู้เล่น</th>
            <th class="py-3 px-4 text-center">เรทติ้ง</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="p in squad"
            :key="p.number"
            class="border-b last:border-b-0"
            :class="isDark ? 'border-sport-border/50 hover:bg-sport-hover/50' : 'border-gray-100 hover:bg-gray-50'"
          >
            <td class="py-2.5 px-4 font-bold text-sport-red">{{ p.position }}</td>
            <td class="py-2.5 px-3 text-center tabular-nums">{{ p.number }}</td>
            <td class="py-2.5 px-4 font-medium">{{ p.name }}</td>
            <td class="py-2.5 px-4 text-center">
              <span
                class="inline-block w-7 py-0.5 rounded text-xs font-bold text-white"
                :class="(p.rating ?? 6) >= 7 ? 'bg-emerald-600' : 'bg-neutral-600'"
              >
                {{ p.rating ?? '-' }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { SquadPlayer } from '@/types'

defineProps<{
  squad: SquadPlayer[]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
</script>
