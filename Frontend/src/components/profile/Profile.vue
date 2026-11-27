<template>
  <div class="w-full max-w-3xl mx-auto space-y-6">
    <!-- Not Logged In State -->
    <div v-if="!authStore.isAuthenticated || !authStore.user" class="card p-8 text-center space-y-4 border-t-4 border-t-sport-red">
      <div class="w-14 h-14 rounded-2xl bg-sport-red/15 text-sport-red flex items-center justify-center mx-auto text-2xl">
        👤
      </div>
      <h1 class="text-xl font-extrabold">กรุณาเข้าสู่ระบบ</h1>
      <p class="text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        เข้าสู่ระบบเพื่อดูและแก้ไขข้อมูลโปรไฟล์ของคุณ รวมถึงจัดการทีมโปรด
      </p>
      <div class="flex items-center justify-center gap-3 pt-2">
        <RouterLink to="/login?redirect=/profile" class="btn-primary text-sm">
          เข้าสู่ระบบ
        </RouterLink>
        <RouterLink
          to="/register"
          class="px-4 py-2 rounded-lg text-sm font-semibold border transition-colors"
          :class="isDark
            ? 'border-sport-border text-neutral-200 hover:border-sport-red'
            : 'border-gray-300 text-gray-700 hover:border-sport-red'"
        >
          สมัครสมาชิก
        </RouterLink>
      </div>
    </div>

    <!-- Authenticated Profile -->
    <template v-else>
      <!-- Profile Summary Card -->
      <div class="card p-6 flex flex-col sm:flex-row items-center justify-between gap-5 border-t-4 border-t-sport-red">
        <div class="flex flex-col sm:flex-row items-center gap-5 text-center sm:text-left">
          <!-- Clickable Avatar -->
          <button
            type="button"
            @click="openAvatarModal"
            title="คลิกเพื่อเปลี่ยนรูปโปรไฟล์"
            class="relative group rounded-full focus:outline-none focus:ring-2 focus:ring-sport-red shrink-0"
          >
            <img
              :src="avatarPreview"
              :alt="authStore.user.username"
              class="w-24 h-24 rounded-full object-cover border-2 border-sport-red shadow-lg bg-neutral-800 group-hover:opacity-75 transition-opacity"
              @error="onAvatarError"
            />
            <span
              class="absolute inset-0 rounded-full bg-black/55 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center text-white text-[11px] font-bold"
            >
              <svg class="w-5 h-5 mb-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              เปลี่ยนรูป
            </span>
            <span
              class="absolute bottom-0 right-0 w-7 h-7 rounded-full bg-sport-red text-white flex items-center justify-center border-2 shadow"
              :class="isDark ? 'border-sport-card' : 'border-white'"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
              </svg>
            </span>
          </button>

          <div class="space-y-1">
            <div class="flex flex-wrap items-center justify-center sm:justify-start gap-2">
              <h1 class="text-2xl font-black">{{ authStore.user.username }}</h1>
              <span class="px-2 py-0.5 text-[10px] font-bold uppercase rounded bg-sport-red text-white">
                {{ authStore.user.role }}
              </span>
            </div>
            <p class="text-sm" :class="isDark ? 'text-neutral-300' : 'text-gray-600'">
              {{ authStore.user.email }}
            </p>
            <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              สมาชิกตั้งแต่: {{ formatDate(authStore.user.created_at) }}
            </p>
          </div>
        </div>

        <button
          type="button"
          @click="showLogoutModal = true"
          class="px-4 py-2 rounded-lg text-xs font-bold border border-sport-red text-sport-red hover:bg-sport-red hover:text-white transition-colors shrink-0"
        >
          ออกจากระบบ
        </button>
      </div>

      <!-- Edit Username Form -->
      <div class="card p-6 space-y-5">
        <div class="flex items-center gap-2">
          <span class="w-1.5 h-5 bg-sport-red rounded-full" />
          <h2 class="text-lg font-extrabold">แก้ไขชื่อผู้ใช้</h2>
        </div>

        <div
          v-if="successMsg"
          class="p-3 rounded-lg bg-emerald-500/15 border border-emerald-500 text-xs text-emerald-400 font-semibold"
        >
          {{ successMsg }}
        </div>
        <div
          v-if="errorMsg"
          class="p-3 rounded-lg bg-sport-red/15 border border-sport-red text-xs text-sport-red font-semibold"
        >
          {{ errorMsg }}
        </div>

        <form class="space-y-4" @submit.prevent="handleSaveUsername">
          <div class="space-y-1.5">
            <label class="text-xs font-semibold block">ชื่อผู้ใช้ (Username)</label>
            <input
              v-model="usernameInput"
              type="text"
              required
              placeholder="a-z, 0-9, _ (3-30 ตัวอักษร)"
              class="w-full px-3.5 py-2.5 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
              :class="isDark
                ? 'bg-sport-bg border-sport-border text-white'
                : 'bg-gray-50 border-gray-300 text-gray-900'"
            />
          </div>

          <div class="flex justify-end">
            <button
              type="submit"
              :disabled="authStore.loading"
              class="btn-primary text-sm px-6 disabled:opacity-50"
            >
              {{ authStore.loading ? 'กำลังบันทึก...' : 'บันทึกชื่อผู้ใช้' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Favorite Teams Summary -->
      <div class="card p-6 space-y-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="w-1.5 h-5 bg-sport-red rounded-full" />
            <h2 class="text-lg font-extrabold">ทีมโปรดของฉัน ({{ favorites.length }})</h2>
          </div>
          <RouterLink to="/favorites" class="text-xs font-semibold text-sport-red hover:underline">
            จัดการทีมโปรด →
          </RouterLink>
        </div>

        <div v-if="favorites.length > 0" class="grid grid-cols-2 sm:grid-cols-3 gap-3">
          <RouterLink
            v-for="team in favorites"
            :key="team.id"
            :to="`/teams/${team.id}`"
            class="flex items-center gap-2.5 p-2.5 rounded-lg border transition-colors"
            :class="isDark
              ? 'border-sport-border hover:bg-sport-hover hover:border-sport-red'
              : 'border-gray-200 hover:bg-gray-50 hover:border-sport-red'"
          >
            <img :src="team.logo_url" :alt="team.name" class="w-7 h-7 object-contain shrink-0" />
            <span class="text-xs font-semibold truncate">{{ team.name }}</span>
          </RouterLink>
        </div>

        <p v-else class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          ยังไม่มีทีมโปรด สามารถเพิ่มได้จากหน้าทีมโปรดหรือหน้ารายละเอียดทีม
        </p>
      </div>

      <!-- Avatar Picker / Edit Modal -->
      <div
        v-if="showAvatarModal"
        class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
        @click.self="showAvatarModal = false"
      >
        <div class="card w-full max-w-md p-6 space-y-5 border-t-4 border-t-sport-red">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-extrabold">เปลี่ยนรูปโปรไฟล์</h2>
            <button
              type="button"
              @click="showAvatarModal = false"
              class="text-sm px-2 py-1 rounded hover:bg-sport-hover"
              :class="isDark ? 'text-neutral-400 hover:text-white' : 'text-gray-500 hover:text-gray-900'"
            >
              ✕
            </button>
          </div>

          <!-- Preview -->
          <div class="flex justify-center">
            <img
              :src="modalPreviewSrc"
              alt="Avatar preview"
              class="w-24 h-24 rounded-full object-cover border-2 border-sport-red bg-neutral-800"
              @error="onAvatarError"
            />
          </div>

          <!-- Upload Image File (FormData) -->
          <div class="space-y-1.5">
            <label class="text-xs font-semibold block">อัปโหลดรูปจากเครื่อง (FormData)</label>
            <label
              class="flex items-center justify-center gap-2 px-4 py-3 rounded-lg border border-dashed cursor-pointer transition-colors text-xs font-semibold"
              :class="isDark
                ? 'border-sport-border hover:border-sport-red bg-sport-bg/50 text-neutral-200'
                : 'border-gray-300 hover:border-sport-red bg-gray-50 text-gray-700'"
            >
              <span>📁 {{ selectedAvatarFile ? selectedAvatarFile.name : 'เลือกไฟล์รูปภาพ (JPG, PNG, WEBP)' }}</span>
              <input
                type="file"
                accept="image/*"
                class="hidden"
                @change="onAvatarFileChange"
              />
            </label>
          </div>

          <!-- Preset Avatars -->
          <div class="space-y-2">
            <p class="text-xs font-semibold" :class="isDark ? 'text-neutral-300' : 'text-gray-700'">
              หรือเลือกรูปแนะนำ
            </p>
            <div class="grid grid-cols-4 gap-3">
              <button
                v-for="preset in presetAvatars"
                :key="preset"
                type="button"
                @click="selectPresetAvatar(preset)"
                class="p-1 rounded-full border-2 transition-all flex items-center justify-center"
                :class="!selectedAvatarFile && modalAvatarUrl === preset ? 'border-sport-red scale-105' : 'border-transparent hover:border-neutral-500'"
              >
                <img :src="preset" alt="Preset" class="w-12 h-12 rounded-full object-cover" />
              </button>
            </div>
          </div>

          <!-- Custom URL -->
          <!-- <div class="space-y-1.5">
            <label class="text-xs font-semibold block">หรือวาง URL รูปภาพ (http:// หรือ https://)</label>
            <input
              v-model="modalAvatarUrl"
              type="url"
              placeholder="https://example.com/avatar.png"
              @input="clearSelectedFile"
              class="w-full px-3.5 py-2 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
              :class="isDark
                ? 'bg-sport-bg border-sport-border text-white'
                : 'bg-gray-50 border-gray-300 text-gray-900'"
            />
          </div> -->

          <div
            v-if="avatarError"
            class="p-2.5 rounded-lg bg-sport-red/15 border border-sport-red text-xs text-sport-red font-semibold"
          >
            {{ avatarError }}
          </div>

          <div class="flex justify-end gap-2">
            <button
              type="button"
              @click="showAvatarModal = false"
              class="px-4 py-2 rounded-lg text-xs font-semibold border"
              :class="isDark ? 'border-sport-border text-neutral-300' : 'border-gray-300 text-gray-700'"
            >
              ยกเลิก
            </button>
            <button
              type="button"
              :disabled="authStore.loading"
              @click="handleSaveAvatar"
              class="btn-primary text-xs px-5 disabled:opacity-50"
            >
              {{ authStore.loading ? 'กำลังบันทึก...' : 'บันทึกรูปโปรไฟล์' }}
            </button>
          </div>
        </div>
      </div>
      <!-- Logout Confirm Modal -->
      <LogoutConfirmModal
        :open="showLogoutModal"
        @cancel="showLogoutModal = false"
        @confirm="handleLogout"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { fetchFavorites } from '@/api/favorites'
import type { Team } from '@/types'
import LogoutConfirmModal from '@/components/layout/LogoutConfirmModal.vue'

const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const usernameInput = ref('')
const successMsg = ref('')
const errorMsg = ref('')
const favorites = ref<Team[]>([])

const showAvatarModal = ref(false)
const showLogoutModal = ref(false)
const modalAvatarUrl = ref('')
const selectedAvatarFile = ref<File | null>(null)
const selectedAvatarPreviewUrl = ref('')
const avatarError = ref('')

const modalPreviewSrc = computed(() => {
  return selectedAvatarPreviewUrl.value || modalAvatarUrl.value || avatarPreview.value
})

const presetAvatars = computed(() => {
  const name = encodeURIComponent(authStore.user?.username || 'Footb')
  return [
    `https://ui-avatars.com/api/?name=${name}&background=dc2626&color=ffffff&size=256`,
    `https://ui-avatars.com/api/?name=${name}&background=0a0a0a&color=dc2626&size=256`,
    `https://ui-avatars.com/api/?name=${name}&background=1e293b&color=38bdf8&size=256`,
    `https://ui-avatars.com/api/?name=${name}&background=14532d&color=4ade80&size=256`,
  ]
})

const avatarPreview = computed(() => {
  return authStore.user?.avatar_url || 'https://ui-avatars.com/api/?name=User&background=dc2626&color=fff'
})

function syncInputs() {
  if (authStore.user) {
    usernameInput.value = authStore.user.username
  }
}

async function loadFavorites() {
  if (!authStore.isAuthenticated) return
  try {
    favorites.value = await fetchFavorites()
  } catch {
    favorites.value = []
  }
}

watch(() => authStore.user, () => {
  syncInputs()
  loadFavorites()
}, { immediate: true })

onMounted(async () => {
  if (authStore.token && !authStore.user) {
    await authStore.loadProfile()
  }
  syncInputs()
  await loadFavorites()
})

function clearSelectedFile() {
  if (selectedAvatarPreviewUrl.value) {
    URL.revokeObjectURL(selectedAvatarPreviewUrl.value)
  }
  selectedAvatarFile.value = null
  selectedAvatarPreviewUrl.value = ''
}

function selectPresetAvatar(url: string) {
  clearSelectedFile()
  modalAvatarUrl.value = url
}

function onAvatarFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  clearSelectedFile()
  selectedAvatarFile.value = file
  selectedAvatarPreviewUrl.value = URL.createObjectURL(file)
}

function openAvatarModal() {
  avatarError.value = ''
  clearSelectedFile()
  modalAvatarUrl.value = authStore.user?.avatar_url || ''
  showAvatarModal.value = true
}

async function handleSaveAvatar() {
  avatarError.value = ''
  successMsg.value = ''

  const formData = new FormData()
  if (selectedAvatarFile.value) {
    formData.append('avatar', selectedAvatarFile.value)
  } else {
    const av = modalAvatarUrl.value.trim()
    if (!av || !/^https?:\/\//i.test(av)) {
      avatarError.value = 'กรุณาเลือกไฟล์รูปภาพ หรือระบุ URL ที่ขึ้นต้นด้วย http:// หรือ https://'
      return
    }
    formData.append('avatar_url', av)
  }

  try {
    await authStore.saveProfile(formData)
    clearSelectedFile()
    showAvatarModal.value = false
    successMsg.value = 'เปลี่ยนรูปโปรไฟล์เรียบร้อยแล้ว'
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: string } } }
    avatarError.value = axiosErr.response?.data?.error || 'ไม่สามารถเปลี่ยนรูปโปรไฟล์ได้'
  }
}

async function handleSaveUsername() {
  successMsg.value = ''
  errorMsg.value = ''

  const u = usernameInput.value.trim()
  if (!/^[a-zA-Z0-9_]{3,30}$/.test(u)) {
    errorMsg.value = 'ชื่อผู้ใช้ต้องยาว 3-30 ตัวอักษร และใช้ได้เฉพาะภาษาอังกฤษ ตัวเลข และ _'
    return
  }

  try {
    await authStore.saveProfile({ username: u })
    successMsg.value = 'บันทึกชื่อผู้ใช้เรียบร้อยแล้ว'
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: string } } }
    errorMsg.value = axiosErr.response?.data?.error || 'ไม่สามารถอัปเดตชื่อผู้ใช้ได้'
  }
}

function handleLogout() {
  showLogoutModal.value = false
  authStore.logout()
  router.push('/')
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('th-TH', { day: 'numeric', month: 'long', year: 'numeric' })
}

function onAvatarError(e: Event) {
  (e.target as HTMLImageElement).src = 'https://ui-avatars.com/api/?name=User&background=dc2626&color=fff'
}
</script>
