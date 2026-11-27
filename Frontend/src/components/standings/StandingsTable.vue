<template>
  <div class="card overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-sm text-left">
        <thead>
          <tr
            class="border-b text-xs uppercase"
            :class="isDark ? 'border-sport-border text-neutral-400 bg-sport-bg/50' : 'border-gray-200 text-gray-500 bg-gray-50'"
          >
            <th class="py-3 px-4 w-12">ลำดับ</th>
            <th class="py-3 px-2">ทีม</th>
            <th class="py-3 px-2 text-center">แข่ง</th>
            <th class="py-3 px-2 text-center">ชนะ</th>
            <th class="py-3 px-2 text-center">เสมอ</th>
            <th class="py-3 px-2 text-center">แพ้</th>
            <th class="py-3 px-2 text-center hidden sm:table-cell">+/-</th>
            <th class="py-3 px-2 text-center">ผลต่าง</th>
            <th class="py-3 px-2 text-center hidden md:table-cell">ฟอร์ม 5 นัด</th>
            <th class="py-3 px-4 text-center font-bold">คะแนน</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="row in rows"
            :key="row.position"
            class="border-b last:border-b-0 transition-colors"
            :class="isDark ? 'border-sport-border/60 hover:bg-sport-hover/60' : 'border-gray-100 hover:bg-gray-50'"
          >
            <td class="py-2.5 px-4 font-bold relative">
              <span
                v-if="row.zone"
                class="absolute left-0 top-1.5 bottom-1.5 w-1 rounded-r"
                :class="zoneColor(row.zone)"
              />
              {{ row.position }}
            </td>
            <td class="py-2.5 px-2">
              <RouterLink :to="`/teams/${row.team.id}`" class="flex items-center gap-2.5 hover:text-sport-red transition-colors">
                <img :src="row.team.logo_url" :alt="row.team.name" class="w-6 h-6 object-contain shrink-0" />
                <span class="font-semibold truncate">{{ row.team.name }}</span>
              </RouterLink>
            </td>
            <td class="py-2.5 px-2 text-center tabular-nums">{{ row.played }}</td>
            <td class="py-2.5 px-2 text-center tabular-nums">{{ row.won }}</td>
            <td class="py-2.5 px-2 text-center tabular-nums">{{ row.drawn }}</td>
            <td class="py-2.5 px-2 text-center tabular-nums">{{ row.lost }}</td>
            <td class="py-2.5 px-2 text-center tabular-nums hidden sm:table-cell">
              {{ row.goals_for }}-{{ row.goals_against }}
            </td>
            <td class="py-2.5 px-2 text-center tabular-nums">
              {{ row.goal_difference > 0 ? `+${row.goal_difference}` : row.goal_difference }}
            </td>
            <td class="py-2.5 px-2 hidden md:table-cell">
              <div v-if="row.form?.length" class="flex items-center justify-center gap-1">
                <span
                  v-for="(f, i) in row.form"
                  :key="i"
                  class="w-5 h-5 rounded text-[10px] font-extrabold text-white flex items-center justify-center"
                  :class="f === 'W' ? 'bg-emerald-600' : f === 'L' ? 'bg-sport-red' : 'bg-neutral-500'"
                >
                  {{ f }}
                </span>
              </div>
              <span v-else class="text-xs text-neutral-500 block text-center">-</span>
            </td>
            <td class="py-2.5 px-4 text-center font-extrabold text-sport-red tabular-nums">
              {{ row.points }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Legend -->
    <div
      class="px-4 py-3 border-t flex flex-wrap items-center gap-5 text-xs"
      :class="isDark ? 'border-sport-border text-neutral-300' : 'border-gray-200 text-gray-600'"
    >
      <div class="flex items-center gap-2">
        <span class="w-3 h-3 rounded-full bg-yellow-400 inline-block" />
        <span>เอเอฟซีแชมเปียนส์ลีกอีลิต</span>
      </div>
      <div class="flex items-center gap-2">
        <span class="w-3 h-3 rounded-full bg-blue-500 inline-block" />
        <span>เอเอฟซีแชมเปียนส์ลีกลีกสอง</span>
      </div>
      <div class="flex items-center gap-2">
        <span class="w-3 h-3 rounded-full bg-sport-red inline-block" />
        <span>ตกชั้นสู่ Thai League 2</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { StandingRow } from '@/types'

defineProps<{
  rows: StandingRow[]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function zoneColor(zone: StandingRow['zone']) {
  if (zone === 'champion' || zone === 'acl') return 'bg-yellow-400'
  if (zone === 'acl2') return 'bg-blue-500'
  if (zone === 'relegation') return 'bg-sport-red'
  return ''
}
</script>
