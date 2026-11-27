<template>
  <MainLayout>
    <div class="space-y-6">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-xl sm:text-2xl font-extrabold flex items-center gap-2.5">
            <span class="w-1.5 h-6 bg-sport-red rounded-full" />
            <span>ทายผลบอลไทยลีก</span>
          </h1>
          <p class="text-xs mt-1" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            อิงโมเดล Prediction (บันทึกในเครื่องชั่วคราว — รอ Backend Prediction API)
          </p>
        </div>

        <div class="card px-4 py-2 flex items-center gap-3 text-xs">
          <span>คะแนนสะสม: <strong class="text-sport-red text-base tabular-nums">{{ totalPoints }}</strong> แต้ม</span>
        </div>
      </div>

      <div v-if="loading" class="card p-4 space-y-3">
        <div v-for="i in 4" :key="i" class="skeleton h-16 w-full" />
      </div>

      <div v-else-if="matches.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div
          v-for="m in matches"
          :key="m.id"
          class="card p-4 space-y-4 border-t-2 border-t-sport-red"
        >
          <div class="flex items-center justify-between text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            <span>{{ m.league?.name }} • {{ m.round }}</span>
            <span>{{ formatDate(m.match_date) }}</span>
          </div>

          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2 flex-1 min-w-0">
              <img :src="m.home_team?.logo_url" :alt="m.home_team?.name" class="w-8 h-8 object-contain shrink-0" />
              <span class="text-xs sm:text-sm font-bold truncate">{{ m.home_team?.name }}</span>
            </div>

            <div class="flex items-center gap-1.5 shrink-0">
              <input
                v-model.number="drafts[m.id].home"
                type="number"
                min="0"
                max="20"
                class="w-12 text-center py-1 rounded-lg border font-black text-sm outline-none focus:ring-2 focus:ring-sport-red"
                :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-gray-50 border-gray-300 text-gray-900'"
              />
              <span class="font-bold text-sport-red">-</span>
              <input
                v-model.number="drafts[m.id].away"
                type="number"
                min="0"
                max="20"
                class="w-12 text-center py-1 rounded-lg border font-black text-sm outline-none focus:ring-2 focus:ring-sport-red"
                :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-gray-50 border-gray-300 text-gray-900'"
              />
            </div>

            <div class="flex items-center justify-end gap-2 flex-1 min-w-0 text-right">
              <span class="text-xs sm:text-sm font-bold truncate">{{ m.away_team?.name }}</span>
              <img :src="m.away_team?.logo_url" :alt="m.away_team?.name" class="w-8 h-8 object-contain shrink-0" />
            </div>
          </div>

          <div class="flex items-center justify-between gap-2 pt-2 border-t" :class="isDark ? 'border-sport-border/60' : 'border-gray-200'">
            <div class="text-xs">
              <span v-if="savedPredictions[m.id]" class="text-emerald-400 font-semibold">
                ✓ ทายไว้: {{ savedPredictions[m.id].predicted_home_score }}-{{ savedPredictions[m.id].predicted_away_score }}
                ({{ resultLabel(savedPredictions[m.id].predicted_result) }})
              </span>
              <span v-else :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
                ยังไม่ทายผล
              </span>
            </div>

            <button
              type="button"
              @click="savePrediction(m.id)"
              class="btn-primary py-1.5 px-4 text-xs"
            >
              บันทึกคำทาย
            </button>
          </div>
        </div>
      </div>

      <div v-else class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ไม่พบแมตช์สำหรับทายผลในขณะนี้
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchMatches } from '@/api/matches'
import type { Match } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'

interface LocalPrediction {
  match_id: string
  predicted_home_score: number
  predicted_away_score: number
  predicted_result: 'home_win' | 'draw' | 'away_win'
  points_earned: number
}

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const matches = ref<Match[]>([])
const loading = ref(false)
const drafts = reactive<Record<string, { home: number; away: number }>>({})
const savedPredictions = ref<Record<string, LocalPrediction>>(
  JSON.parse(localStorage.getItem('footb_predictions') || '{}')
)

const totalPoints = computed(() =>
  Object.values(savedPredictions.value).reduce((sum, p) => sum + p.points_earned, 0)
)

onMounted(async () => {
  loading.value = true
  try {
    const data = await fetchMatches({ limit: 12 })
    matches.value = data
    for (const m of data) {
      const existing = savedPredictions.value[m.id]
      drafts[m.id] = {
        home: existing ? existing.predicted_home_score : 1,
        away: existing ? existing.predicted_away_score : 0,
      }
    }
  } catch {
    matches.value = []
  } finally {
    loading.value = false
  }
})

function savePrediction(matchId: string) {
  const d = drafts[matchId]
  const hs = Math.max(0, Number(d?.home ?? 0))
  const as = Math.max(0, Number(d?.away ?? 0))
  const result: LocalPrediction['predicted_result'] =
    hs > as ? 'home_win' : hs < as ? 'away_win' : 'draw'

  savedPredictions.value[matchId] = {
    match_id: matchId,
    predicted_home_score: hs,
    predicted_away_score: as,
    predicted_result: result,
    points_earned: 0,
  }
  localStorage.setItem('footb_predictions', JSON.stringify(savedPredictions.value))
}

function resultLabel(r: LocalPrediction['predicted_result']) {
  if (r === 'home_win') return 'เจ้าบ้านชนะ'
  if (r === 'away_win') return 'ทีมเยือนชนะ'
  return 'เสมอ'
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('th-TH', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}
</script>
