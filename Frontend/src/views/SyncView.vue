<template>
  <MainLayout>
    <div class="space-y-6">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center gap-2.5">
          <span class="w-1.5 h-6 bg-sport-red rounded-full" />
          <h1 class="text-xl sm:text-2xl font-extrabold">ซิงก์ข้อมูลจาก API-Football (Admin)</h1>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Sync Fixtures (POST /sync/fixtures) -->
        <div class="card p-6 space-y-4 border-t-4 border-t-sport-red">
          <h2 class="text-lg font-extrabold">1. ซิงก์ตารางแข่งขันและผลบอล</h2>
          <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            เรียก <code class="text-sport-red">POST /api/v1/sync/fixtures</code> เพื่อดึงแมตช์จาก API-Football ลงฐานข้อมูล
          </p>

          <form class="space-y-3" @submit.prevent="runSyncFixtures">
            <div class="space-y-1">
              <label class="text-xs font-semibold block">ฤดูกาล (Season) *</label>
              <input
                v-model.number="season"
                type="number"
                required
                min="2020"
                max="2030"
                class="w-full px-3 py-2 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
                :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-gray-50 border-gray-300 text-gray-900'"
              />
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div class="space-y-1">
                <label class="text-xs font-semibold block">ตั้งแต่วันที่ (Optional)</label>
                <input
                  v-model="fromDate"
                  type="date"
                  class="w-full px-3 py-2 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
                  :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-gray-50 border-gray-300 text-gray-900'"
                />
              </div>
              <div class="space-y-1">
                <label class="text-xs font-semibold block">ถึงวันที่ (Optional)</label>
                <input
                  v-model="toDate"
                  type="date"
                  class="w-full px-3 py-2 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
                  :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-gray-50 border-gray-300 text-gray-900'"
                />
              </div>
            </div>

            <button
              type="submit"
              :disabled="syncingFixtures"
              class="w-full btn-primary py-2.5 text-xs font-bold disabled:opacity-50"
            >
              {{ syncingFixtures ? 'กำลังซิงก์ Fixtures...' : 'ซิงก์ Fixtures ทันที' }}
            </button>
          </form>

          <div v-if="fixturesResult" class="p-3.5 rounded-xl border border-emerald-500/50 bg-emerald-500/10 text-xs space-y-1">
            <p class="font-bold text-emerald-400">✓ ซิงก์ Fixtures สำเร็จ (Season {{ fixturesResult.season }})</p>
            <p>ลีกที่ซิงก์: {{ fixturesResult.leagues_synced }} | ทีมใหม่: {{ fixturesResult.teams_created }}</p>
            <p>แมตช์ใหม่: {{ fixturesResult.matches_created }} | อัปเดต: {{ fixturesResult.matches_updated }} | ข้าม: {{ fixturesResult.matches_skipped }}</p>
            <p v-if="fixturesResult.quota_remaining" class="text-neutral-300">
              Quota คงเหลือ — วันนี้: {{ fixturesResult.quota_remaining.day }}, นาทีนี้: {{ fixturesResult.quota_remaining.minute }}
            </p>
          </div>
          <div v-if="fixturesError" class="p-3 rounded-lg bg-sport-red/15 border border-sport-red text-xs text-sport-red font-semibold">
            {{ fixturesError }}
          </div>
        </div>

        <!-- Sync Squads (POST /sync/squads) -->
        <div class="card p-6 space-y-4 border-t-4 border-t-sport-red flex flex-col justify-between">
          <div class="space-y-3">
            <h2 class="text-lg font-extrabold">2. ซิงก์โค้ชและรายชื่อนักเตะ</h2>
            <p class="text-xs leading-relaxed" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              เรียก <code class="text-sport-red">POST /api/v1/sync/squads</code> เพื่อดึงข้อมูลหัวหน้าผู้ฝึกสอนและนักเตะของทุกทีมในลีกลงฐานข้อมูลอัตโนมัติ
            </p>
          </div>

          <div class="space-y-3">
            <button
              type="button"
              :disabled="syncingSquads"
              @click="runSyncSquads"
              class="w-full btn-primary py-2.5 text-xs font-bold disabled:opacity-50"
            >
              {{ syncingSquads ? 'กำลังซิงก์ Squads...' : 'ซิงก์ Squads & Coaches ทันที' }}
            </button>

            <div v-if="squadsResult" class="p-3.5 rounded-xl border border-emerald-500/50 bg-emerald-500/10 text-xs">
              <p class="font-bold text-emerald-400 mb-1">✓ ซิงก์ Squads สำเร็จ</p>
              <pre class="overflow-x-auto text-[11px]">{{ JSON.stringify(squadsResult, null, 2) }}</pre>
            </div>
            <div v-if="squadsError" class="p-3 rounded-lg bg-sport-red/15 border border-sport-red text-xs text-sport-red font-semibold">
              {{ squadsError }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { syncFixtures, syncSquads, type SyncFixturesResponse } from '@/api/sync'
import MainLayout from '@/components/layout/MainLayout.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const season = ref(2024)
const fromDate = ref('')
const toDate = ref('')

const syncingFixtures = ref(false)
const fixturesResult = ref<SyncFixturesResponse | null>(null)
const fixturesError = ref('')

const syncingSquads = ref(false)
const squadsResult = ref<Record<string, unknown> | null>(null)
const squadsError = ref('')

async function runSyncFixtures() {
  syncingFixtures.value = true
  fixturesError.value = ''
  fixturesResult.value = null
  try {
    fixturesResult.value = await syncFixtures({
      season: season.value,
      from: fromDate.value || undefined,
      to: toDate.value || undefined,
    })
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: string } } }
    fixturesError.value = axiosErr.response?.data?.error || 'ไม่สามารถซิงก์ Fixtures ได้'
  } finally {
    syncingFixtures.value = false
  }
}

async function runSyncSquads() {
  syncingSquads.value = true
  squadsError.value = ''
  squadsResult.value = null
  try {
    squadsResult.value = await syncSquads()
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: string } } }
    squadsError.value = axiosErr.response?.data?.error || 'ไม่สามารถซิงก์ Squads ได้'
  } finally {
    syncingSquads.value = false
  }
}
</script>
