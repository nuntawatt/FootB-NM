<template>
  <div class="w-full max-w-4xl mx-auto space-y-6">
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
      <!-- Sport Cover + Profile Header Card -->
      <div class="card overflow-hidden shadow-xl">
        <!-- Red/Black Sport Cover Banner -->
        <div class="h-28 sm:h-36 bg-gradient-to-r from-black via-sport-red-dark to-sport-red relative overflow-hidden">
          <div class="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_top_right,_#ffffff,_transparent_60%)]" />
          <div class="absolute top-3 right-4 flex items-center gap-2">
            <button
              type="button"
              @click="refreshProfile"
              :disabled="refreshing"
              title="รีเฟรชข้อมูลจากเซิร์ฟเวอร์"
              class="px-3 py-1.5 rounded-lg text-xs font-bold bg-black/50 hover:bg-black/80 text-white backdrop-blur-sm transition-colors disabled:opacity-50"
            >
              {{ refreshing ? 'กำลังโหลด...' : '↻ รีเฟรช' }}
            </button>
            <button
              type="button"
              @click="showLogoutModal = true"
              class="px-3 py-1.5 rounded-lg text-xs font-bold bg-black/50 hover:bg-sport-red text-white backdrop-blur-sm transition-colors"
            >
              ออกจากระบบ
            </button>
          </div>
        </div>

        <!-- Avatar & User Summary -->
        <div class="px-6 pb-6">
          <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-4 -mt-12 sm:-mt-14">
            <div class="flex flex-col sm:flex-row items-center sm:items-end gap-4 text-center sm:text-left">
              <!-- Clickable Avatar -->
              <button
                type="button"
                @click="openEditProfile"
                title="คลิกเพื่อแก้ไขโปรไฟล์ / เปลี่ยนรูป"
                class="relative group rounded-full focus:outline-none focus:ring-2 focus:ring-sport-red shrink-0"
              >
                <img
                  :src="avatarPreview"
                  :alt="authStore.user.username"
                  class="w-24 h-24 sm:w-28 sm:h-28 rounded-full object-cover border-4 shadow-xl bg-neutral-900 group-hover:opacity-80 transition-opacity"
                  :class="isDark ? 'border-sport-card' : 'border-white'"
                  @error="onAvatarError"
                />
                <span
                  class="absolute inset-0 rounded-full bg-black/55 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center text-white text-[11px] font-bold"
                >
                  <svg class="w-5 h-5 mb-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  แก้ไขรูป
                </span>
                <span
                  class="absolute bottom-1 right-1 w-7 h-7 rounded-full bg-sport-red text-white flex items-center justify-center border-2 shadow"
                  :class="isDark ? 'border-sport-card' : 'border-white'"
                >
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                  </svg>
                </span>
              </button>

              <div class="pb-1 space-y-1">
                <div class="flex flex-wrap items-center justify-center sm:justify-start gap-2">
                  <h1 class="text-2xl font-black">{{ authStore.user.username }}</h1>
                  <span class="px-2.5 py-0.5 text-[10px] font-extrabold uppercase rounded-full bg-sport-red text-white">
                    {{ authStore.user.role === 'admin' ? 'ผู้ดูแลระบบ' : 'สมาชิก' }}
                  </span>
                </div>
                <p class="text-sm" :class="isDark ? 'text-neutral-300' : 'text-gray-600'">
                  {{ authStore.user.email }}
                </p>
              </div>
            </div>

            <!-- Quick Stats Strip -->
            <div class="flex items-center justify-center gap-3 pt-2 sm:pt-0">
              <div class="text-center px-3 py-1.5 rounded-lg border" :class="isDark ? 'border-sport-border bg-sport-bg/60' : 'border-gray-200 bg-gray-50'">
                <p class="text-lg font-black text-sport-red tabular-nums">{{ favorites.length }}</p>
                <p class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">ทีมโปรด</p>
              </div>
              <div class="text-center px-3 py-1.5 rounded-lg border" :class="isDark ? 'border-sport-border bg-sport-bg/60' : 'border-gray-200 bg-gray-50'">
                <p class="text-lg font-black text-emerald-500">Active</p>
                <p class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">สถานะบัญชี</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Alerts -->
      <div
        v-if="successMsg"
        class="p-3.5 rounded-xl bg-emerald-500/15 border border-emerald-500 text-xs text-emerald-400 font-semibold flex items-center justify-between"
      >
        <span>✓ {{ successMsg }}</span>
        <button type="button" @click="successMsg = ''" class="opacity-70 hover:opacity-100">✕</button>
      </div>
      <div
        v-if="errorMsg"
        class="p-3.5 rounded-xl bg-sport-red/15 border border-sport-red text-xs text-sport-red font-semibold flex items-center justify-between"
      >
        <span>{{ errorMsg }}</span>
        <button type="button" @click="errorMsg = ''" class="opacity-70 hover:opacity-100">✕</button>
      </div>

      <!-- Account Info Card (GET /profile) -->
      <div class="card p-6 space-y-5 border-t-2 border-t-sport-red">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div class="flex items-center gap-2">
            <span class="w-1.5 h-5 bg-sport-red rounded-full" />
            <h2 class="text-lg font-extrabold">ข้อมูลบัญชีผู้ใช้</h2>
          </div>

          <button
            type="button"
            @click="isEditing ? cancelEditProfile() : openEditProfile()"
            class="px-4 py-1.5 rounded-lg text-xs font-bold border transition-colors"
            :class="isEditing
              ? (isDark ? 'border-sport-border text-neutral-300 hover:bg-sport-hover' : 'border-gray-300 text-gray-700 hover:bg-gray-100')
              : 'border-sport-red bg-sport-red text-white hover:bg-sport-red-hover'"
          >
            {{ isEditing ? '✕ ปิดการแก้ไข' : '✎ แก้ไขโปรไฟล์' }}
          </button>
        </div>

        <!-- View Mode -->
        <div v-if="!isEditing" class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
          <div class="p-3.5 rounded-xl border space-y-1" :class="isDark ? 'border-sport-border bg-sport-bg/40' : 'border-gray-200 bg-gray-50'">
            <p class="text-[11px] font-semibold uppercase tracking-wider" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              ชื่อผู้ใช้ (Username)
            </p>
            <p class="font-bold">{{ authStore.user.username }}</p>
          </div>

          <div class="p-3.5 rounded-xl border space-y-1" :class="isDark ? 'border-sport-border bg-sport-bg/40' : 'border-gray-200 bg-gray-50'">
            <p class="text-[11px] font-semibold uppercase tracking-wider" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              อีเมล (Email)
            </p>
            <p class="font-bold truncate">{{ authStore.user.email }}</p>
          </div>

          <div class="p-3.5 rounded-xl border space-y-1" :class="isDark ? 'border-sport-border bg-sport-bg/40' : 'border-gray-200 bg-gray-50'">
            <p class="text-[11px] font-semibold uppercase tracking-wider" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              สิทธิ์การใช้งาน (Role)
            </p>
            <p class="font-bold uppercase text-sport-red">{{ authStore.user.role }}</p>
          </div>

          <div class="p-3.5 rounded-xl border space-y-1" :class="isDark ? 'border-sport-border bg-sport-bg/40' : 'border-gray-200 bg-gray-50'">
            <p class="text-[11px] font-semibold uppercase tracking-wider" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              วันที่สมัครสมาชิก
            </p>
            <p class="font-bold">{{ formatDate(authStore.user.created_at) }}</p>
          </div>

          <div class="sm:col-span-2 p-3.5 rounded-xl border flex flex-wrap items-center justify-between gap-2" :class="isDark ? 'border-sport-border bg-sport-bg/40' : 'border-gray-200 bg-gray-50'">
            <div class="space-y-0.5 min-w-0">
              <p class="text-[11px] font-semibold uppercase tracking-wider" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
                รหัสสมาชิก (User ID)
              </p>
              <p class="font-mono text-xs truncate">{{ authStore.user.id }}</p>
            </div>
            <button
              type="button"
              @click="copyUserId"
              class="px-3 py-1 rounded-lg text-xs font-semibold border transition-colors shrink-0"
              :class="copied
                ? 'border-emerald-500 text-emerald-400 bg-emerald-500/10'
                : (isDark ? 'border-sport-border text-neutral-300 hover:border-sport-red' : 'border-gray-300 text-gray-700 hover:border-sport-red')"
            >
              {{ copied ? '✓ คัดลอกแล้ว' : 'คัดลอก ID' }}
            </button>
          </div>
        </div>

        <!-- Edit Profile UI -->
        <form v-else class="space-y-6" @submit.prevent="requestConfirmUpdate">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 items-start">
            <!-- Left: Avatar Picker -->
            <div class="flex flex-col items-center gap-3 p-4 rounded-xl border" :class="isDark ? 'border-sport-border bg-sport-bg/40' : 'border-gray-200 bg-gray-50'">
              <img
                :src="editAvatarPreviewSrc"
                alt="Avatar preview"
                class="w-24 h-24 rounded-full object-cover border-2 border-sport-red bg-neutral-800 shadow-md"
                @error="onAvatarError"
              />
              <label
                class="w-full text-center px-3 py-2 rounded-lg border border-dashed cursor-pointer transition-colors text-xs font-semibold"
                :class="isDark
                  ? 'border-sport-border hover:border-sport-red bg-sport-card text-neutral-200'
                  : 'border-gray-300 hover:border-sport-red bg-white text-gray-700'"
              >
                <span>📁 {{ selectedAvatarFile ? selectedAvatarFile.name : 'อัปโหลดรูปจากเครื่อง' }}</span>
                <input
                  type="file"
                  accept="image/*"
                  class="hidden"
                  @change="onAvatarFileChange"
                />
              </label>

              <div class="w-full space-y-1.5">
                <p class="text-[11px] font-semibold text-center" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
                  หรือเลือกรูปแนะนำ
                </p>
                <div class="grid grid-cols-4 gap-2">
                  <button
                    v-for="preset in presetAvatars"
                    :key="preset"
                    type="button"
                    @click="selectPresetAvatar(preset)"
                    class="p-0.5 rounded-full border-2 transition-all flex items-center justify-center"
                    :class="!selectedAvatarFile && editAvatarUrl === preset ? 'border-sport-red scale-105' : 'border-transparent hover:border-neutral-500'"
                  >
                    <img :src="preset" alt="Preset" class="w-9 h-9 rounded-full object-cover" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Right: Fields -->
            <div class="md:col-span-2 space-y-4">
              <div class="space-y-1.5">
                <label class="text-xs font-semibold block">ชื่อผู้ใช้ (Username)</label>
                <input
                  v-model="editUsername"
                  type="text"
                  required
                  placeholder="a-z, 0-9, _ (3-30 ตัวอักษร)"
                  class="w-full px-3.5 py-2.5 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
                  :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-white border-gray-300 text-gray-900'"
                />
                <p class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
                  อนุญาตเฉพาะภาษาอังกฤษ ตัวเลข และขีดล่าง (_) ความยาว 3-30 ตัวอักษร
                </p>
              </div>

              <div class="space-y-1.5">
                <label class="text-xs font-semibold block">URL รูปโปรไฟล์ (ถ้าไม่ใช้อัปโหลดไฟล์)</label>
                <input
                  v-model="editAvatarUrl"
                  type="url"
                  placeholder="https://example.com/avatar.png"
                  @input="clearSelectedFile"
                  class="w-full px-3.5 py-2.5 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
                  :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-white border-gray-300 text-gray-900'"
                />
              </div>

              <div class="space-y-1.5">
                <label class="text-xs font-semibold block">อีเมล (ไม่สามารถแก้ไขได้)</label>
                <input
                  :value="authStore.user.email"
                  type="email"
                  disabled
                  class="w-full px-3.5 py-2.5 rounded-lg text-sm border opacity-60 cursor-not-allowed"
                  :class="isDark ? 'bg-sport-bg border-sport-border text-neutral-400' : 'bg-gray-100 border-gray-300 text-gray-500'"
                />
              </div>
            </div>
          </div>

          <div class="flex items-center justify-end gap-3 pt-2 border-t" :class="isDark ? 'border-sport-border' : 'border-gray-200'">
            <button
              type="button"
              @click="cancelEditProfile"
              class="px-4 py-2 rounded-lg text-xs font-bold border transition-colors"
              :class="isDark ? 'border-sport-border text-neutral-300 hover:bg-sport-hover' : 'border-gray-300 text-gray-700 hover:bg-gray-100'"
            >
              ยกเลิก
            </button>
            <button
              type="submit"
              :disabled="authStore.loading"
              class="btn-primary text-xs px-6 py-2.5 disabled:opacity-50"
            >
              ตรวจสอบและบันทึกข้อมูล
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

      <!-- Update Confirmation Card Modal -->
      <div
        v-if="showConfirmModal"
        class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
        @click.self="showConfirmModal = false"
      >
        <div class="card w-full max-w-md p-6 space-y-5 border-t-4 border-t-sport-red shadow-2xl">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="w-1.5 h-5 bg-sport-red rounded-full" />
              <h3 class="text-lg font-extrabold">ยืนยันข้อมูลที่จะอัปเดต</h3>
            </div>
            <button
              type="button"
              @click="showConfirmModal = false"
              class="text-sm px-2 py-1 rounded hover:bg-sport-hover"
              :class="isDark ? 'text-neutral-400 hover:text-white' : 'text-gray-500 hover:text-gray-900'"
            >
              ✕
            </button>
          </div>

          <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            กรุณาตรวจสอบข้อมูลใหม่ของคุณก่อนกดยืนยันการบันทึก
          </p>

          <!-- Summary of Changes -->
          <div class="space-y-3">
            <!-- Avatar Change Preview -->
            <div
              v-if="hasAvatarChanged"
              class="p-3.5 rounded-xl border space-y-2"
              :class="isDark ? 'border-sport-border bg-sport-bg/60' : 'border-gray-200 bg-gray-50'"
            >
              <p class="text-[11px] font-bold uppercase tracking-wider text-sport-red">รูปโปรไฟล์</p>
              <div class="flex items-center justify-center gap-4">
                <div class="text-center space-y-1">
                  <img
                    :src="avatarPreview"
                    alt="Old avatar"
                    class="w-14 h-14 rounded-full object-cover border opacity-70 mx-auto"
                    @error="onAvatarError"
                  />
                  <span class="text-[10px] block" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">เดิม</span>
                </div>
                <span class="text-sport-red font-black text-lg">➔</span>
                <div class="text-center space-y-1">
                  <img
                    :src="editAvatarPreviewSrc"
                    alt="New avatar"
                    class="w-14 h-14 rounded-full object-cover border-2 border-sport-red mx-auto shadow"
                    @error="onAvatarError"
                  />
                  <span class="text-[10px] font-bold text-sport-red block">ใหม่</span>
                </div>
              </div>
              <p v-if="selectedAvatarFile" class="text-[11px] text-center truncate" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
                ไฟล์: {{ selectedAvatarFile.name }}
              </p>
            </div>

            <!-- Username Change Preview -->
            <div
              v-if="hasUsernameChanged"
              class="p-3.5 rounded-xl border space-y-1.5"
              :class="isDark ? 'border-sport-border bg-sport-bg/60' : 'border-gray-200 bg-gray-50'"
            >
              <p class="text-[11px] font-bold uppercase tracking-wider text-sport-red">ชื่อผู้ใช้ (Username)</p>
              <div class="flex items-center justify-between gap-2 text-sm">
                <span class="line-through opacity-60 truncate">{{ authStore.user.username }}</span>
                <span class="text-sport-red font-black shrink-0">➔</span>
                <span class="font-extrabold truncate">{{ editUsername.trim() }}</span>
              </div>
            </div>
          </div>

          <div class="flex items-center gap-3 pt-1">
            <button
              type="button"
              @click="showConfirmModal = false"
              class="flex-1 px-4 py-2.5 rounded-lg text-xs font-bold border transition-colors"
              :class="isDark
                ? 'border-sport-border text-neutral-300 hover:bg-sport-hover'
                : 'border-gray-300 text-gray-700 hover:bg-gray-100'"
            >
              กลับไปแก้ไข
            </button>
            <button
              type="button"
              :disabled="authStore.loading"
              @click="confirmUpdateProfile"
              class="flex-1 btn-primary py-2.5 text-xs font-bold disabled:opacity-50"
            >
              {{ authStore.loading ? 'กำลังอัปเดต...' : 'ยืนยันการอัปเดต' }}
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

const isEditing = ref(false)
const editUsername = ref('')
const editAvatarUrl = ref('')
const selectedAvatarFile = ref<File | null>(null)
const selectedAvatarPreviewUrl = ref('')

const refreshing = ref(false)
const copied = ref(false)
const successMsg = ref('')
const errorMsg = ref('')
const favorites = ref<Team[]>([])

const showConfirmModal = ref(false)
const showLogoutModal = ref(false)

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

const editAvatarPreviewSrc = computed(() => {
  return selectedAvatarPreviewUrl.value || editAvatarUrl.value || avatarPreview.value
})

const hasUsernameChanged = computed(() => {
  if (!authStore.user) return false
  return editUsername.value.trim() !== authStore.user.username
})

const hasAvatarChanged = computed(() => {
  if (!authStore.user) return false
  if (selectedAvatarFile.value) return true
  return editAvatarUrl.value.trim() !== '' && editAvatarUrl.value.trim() !== authStore.user.avatar_url
})

function syncInputs() {
  if (authStore.user) {
    editUsername.value = authStore.user.username
    editAvatarUrl.value = authStore.user.avatar_url
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

async function refreshProfile() {
  refreshing.value = true
  try {
    await authStore.loadProfile()
    await loadFavorites()
    successMsg.value = 'อัปเดตข้อมูลล่าสุดจากเซิร์ฟเวอร์แล้ว'
  } finally {
    refreshing.value = false
  }
}

function clearSelectedFile() {
  if (selectedAvatarPreviewUrl.value) {
    URL.revokeObjectURL(selectedAvatarPreviewUrl.value)
  }
  selectedAvatarFile.value = null
  selectedAvatarPreviewUrl.value = ''
}

function selectPresetAvatar(url: string) {
  clearSelectedFile()
  editAvatarUrl.value = url
}

function onAvatarFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  clearSelectedFile()
  selectedAvatarFile.value = file
  selectedAvatarPreviewUrl.value = URL.createObjectURL(file)
}

function openEditProfile() {
  errorMsg.value = ''
  successMsg.value = ''
  clearSelectedFile()
  syncInputs()
  isEditing.value = true
}

function cancelEditProfile() {
  clearSelectedFile()
  syncInputs()
  isEditing.value = false
}

function requestConfirmUpdate() {
  errorMsg.value = ''
  successMsg.value = ''

  const u = editUsername.value.trim()
  if (!/^[a-zA-Z0-9_]{3,30}$/.test(u)) {
    errorMsg.value = 'ชื่อผู้ใช้ต้องยาว 3-30 ตัวอักษร และใช้ได้เฉพาะภาษาอังกฤษ ตัวเลข และ _'
    return
  }

  const av = editAvatarUrl.value.trim()
  if (!selectedAvatarFile.value && av && !/^https?:\/\//i.test(av)) {
    errorMsg.value = 'URL รูปโปรไฟล์ต้องขึ้นต้นด้วย http:// หรือ https://'
    return
  }

  if (!hasUsernameChanged.value && !hasAvatarChanged.value) {
    errorMsg.value = 'ยังไม่มีข้อมูลที่เปลี่ยนแปลง'
    return
  }

  showConfirmModal.value = true
}

async function confirmUpdateProfile() {
  errorMsg.value = ''
  successMsg.value = ''

  const formData = new FormData()
  if (hasUsernameChanged.value) {
    formData.append('username', editUsername.value.trim())
  }
  if (selectedAvatarFile.value) {
    formData.append('avatar', selectedAvatarFile.value)
  } else if (hasAvatarChanged.value) {
    formData.append('avatar_url', editAvatarUrl.value.trim())
  }

  try {
    await authStore.saveProfile(formData)
    showConfirmModal.value = false
    isEditing.value = false
    clearSelectedFile()
    successMsg.value = 'อัปเดตข้อมูลโปรไฟล์เรียบร้อยแล้ว'
  } catch (err: unknown) {
    showConfirmModal.value = false
    const axiosErr = err as { response?: { data?: { error?: string } } }
    errorMsg.value = axiosErr.response?.data?.error || 'ไม่สามารถอัปเดตโปรไฟล์ได้'
  }
}

async function copyUserId() {
  if (!authStore.user?.id) return
  try {
    await navigator.clipboard.writeText(authStore.user.id)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    // ignore clipboard permission errors
  }
}

function handleLogout() {
  showLogoutModal.value = false
  authStore.logout()
  router.push('/')
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('th-TH', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function onAvatarError(e: Event) {
  (e.target as HTMLImageElement).src = 'https://ui-avatars.com/api/?name=User&background=dc2626&color=fff'
}
</script>
