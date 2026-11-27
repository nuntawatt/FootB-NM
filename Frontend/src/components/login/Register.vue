<template>
  <div class="card w-full max-w-md p-6 sm:p-8 space-y-6 border-t-4 border-t-sport-red shadow-2xl">
    <!-- Header -->
    <div class="text-center space-y-1.5">
      <div class="w-12 h-12 rounded-xl bg-sport-red flex items-center justify-center mx-auto shadow-lg shadow-sport-red/30">
        <span class="text-white font-black text-xl tracking-tighter">FB</span>
      </div>
      <h1 class="text-2xl font-black tracking-tight pt-2">สมัครสมาชิก</h1>
      <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        สร้างบัญชีเพื่อบันทึกทีมโปรดและติดตามผลบอล
      </p>
    </div>

    <!-- Error Alert -->
    <div
      v-if="errorMsg"
      class="p-3 rounded-lg bg-sport-red/15 border border-sport-red text-xs text-sport-red font-semibold"
    >
      {{ errorMsg }}
    </div>

    <!-- Form -->
    <form class="space-y-4" @submit.prevent="handleSubmit">
      <div class="space-y-1.5">
        <label class="text-xs font-semibold block">ชื่อผู้ใช้ (Username)</label>
        <input
          v-model="username"
          type="text"
          required
          autocomplete="username"
          placeholder="somchai_99 (ภาษาอังกฤษ, ตัวเลข, _ 3-30 ตัว)"
          class="w-full px-3.5 py-2.5 rounded-lg text-sm border outline-none transition-all focus:ring-2 focus:ring-sport-red"
          :class="isDark
            ? 'bg-sport-bg border-sport-border text-white placeholder-neutral-500'
            : 'bg-gray-50 border-gray-300 text-gray-900 placeholder-gray-400'"
        />
      </div>

      <div class="space-y-1.5">
        <label class="text-xs font-semibold block">อีเมล</label>
        <input
          v-model="email"
          type="email"
          required
          autocomplete="email"
          placeholder="somchai@example.com"
          class="w-full px-3.5 py-2.5 rounded-lg text-sm border outline-none transition-all focus:ring-2 focus:ring-sport-red"
          :class="isDark
            ? 'bg-sport-bg border-sport-border text-white placeholder-neutral-500'
            : 'bg-gray-50 border-gray-300 text-gray-900 placeholder-gray-400'"
        />
      </div>

      <div class="space-y-1.5">
        <label class="text-xs font-semibold block">รหัสผ่าน</label>
        <div class="relative">
          <input
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            required
            autocomplete="new-password"
            placeholder="อย่างน้อย 6 ตัว มีตัวอักษรและตัวเลข"
            class="w-full px-3.5 py-2.5 pr-16 rounded-lg text-sm border outline-none transition-all focus:ring-2 focus:ring-sport-red"
            :class="isDark
              ? 'bg-sport-bg border-sport-border text-white placeholder-neutral-500'
              : 'bg-gray-50 border-gray-300 text-gray-900 placeholder-gray-400'"
          />
          <button
            type="button"
            @click="showPassword = !showPassword"
            class="absolute right-3 top-1/2 -translate-y-1/2 text-xs font-semibold"
            :class="isDark ? 'text-neutral-400 hover:text-white' : 'text-gray-500 hover:text-gray-800'"
          >
            {{ showPassword ? 'ซ่อน' : 'แสดง' }}
          </button>
        </div>
      </div>

      <div class="space-y-1.5">
        <label class="text-xs font-semibold block">ยืนยันรหัสผ่าน</label>
        <input
          v-model="confirmPassword"
          :type="showPassword ? 'text' : 'password'"
          required
          autocomplete="new-password"
          placeholder="กรอกรหัสผ่านอีกครั้ง"
          class="w-full px-3.5 py-2.5 rounded-lg text-sm border outline-none transition-all focus:ring-2 focus:ring-sport-red"
          :class="isDark
            ? 'bg-sport-bg border-sport-border text-white placeholder-neutral-500'
            : 'bg-gray-50 border-gray-300 text-gray-900 placeholder-gray-400'"
        />
      </div>

      <div class="space-y-1.5">
        <label class="text-xs font-semibold block">
          URL รูปโปรไฟล์ <span class="font-normal opacity-70">(ไม่บังคับ)</span>
        </label>
        <input
          v-model="avatarUrl"
          type="url"
          placeholder="https://example.com/avatar.png"
          class="w-full px-3.5 py-2.5 rounded-lg text-sm border outline-none transition-all focus:ring-2 focus:ring-sport-red"
          :class="isDark
            ? 'bg-sport-bg border-sport-border text-white placeholder-neutral-500'
            : 'bg-gray-50 border-gray-300 text-gray-900 placeholder-gray-400'"
        />
      </div>

      <button
        type="submit"
        :disabled="authStore.loading"
        class="w-full btn-primary py-2.5 text-sm font-bold disabled:opacity-50"
      >
        {{ authStore.loading ? 'กำลังสมัครสมาชิก...' : 'สมัครสมาชิก' }}
      </button>
    </form>

    <!-- Footer Link -->
    <div
      class="pt-4 border-t text-center text-xs"
      :class="isDark ? 'border-sport-border text-neutral-400' : 'border-gray-200 text-gray-600'"
    >
      มีบัญชีอยู่แล้ว?
      <RouterLink to="/login" class="font-bold text-sport-red hover:underline ml-1">
        เข้าสู่ระบบ
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const avatarUrl = ref('')
const showPassword = ref(false)
const errorMsg = ref('')

function validateForm(): string | null {
  const u = username.value.trim()
  if (!/^[a-zA-Z0-9_]{3,30}$/.test(u)) {
    return 'ชื่อผู้ใช้ต้องยาว 3-30 ตัวอักษร และใช้ได้เฉพาะภาษาอังกฤษ ตัวเลข และ _'
  }
  if (password.value.length < 6 || password.value.length > 72) {
    return 'รหัสผ่านต้องมีความยาว 6-72 ตัวอักษร'
  }
  if (!/[a-zA-Z]/.test(password.value) || !/[0-9]/.test(password.value)) {
    return 'รหัสผ่านต้องมีทั้งตัวอักษรภาษาอังกฤษและตัวเลขผสมกัน'
  }
  if (password.value !== confirmPassword.value) {
    return 'รหัสผ่านและยืนยันรหัสผ่านไม่ตรงกัน'
  }
  const av = avatarUrl.value.trim()
  if (av && !/^https?:\/\//i.test(av)) {
    return 'URL รูปโปรไฟล์ต้องขึ้นต้นด้วย http:// หรือ https://'
  }
  return null
}

async function handleSubmit() {
  errorMsg.value = ''
  const validationError = validateForm()
  if (validationError) {
    errorMsg.value = validationError
    return
  }

  try {
    await authStore.register({
      username: username.value.trim(),
      email: email.value.trim(),
      password: password.value,
      confirm_password: confirmPassword.value,
      avatar_url: avatarUrl.value.trim() || undefined,
    })
    router.push('/')
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: string } } }
    errorMsg.value = axiosErr.response?.data?.error || 'ไม่สามารถสมัครสมาชิกได้ กรุณาลองใหม่'
  }
}
</script>
