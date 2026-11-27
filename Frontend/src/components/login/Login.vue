<template>
  <div class="card w-full max-w-md p-6 sm:p-8 space-y-6 border-t-4 border-t-sport-red shadow-2xl">
    <!-- Header -->
    <div class="text-center space-y-1.5">
      <div class="w-12 h-12 rounded-xl bg-sport-red flex items-center justify-center mx-auto shadow-lg shadow-sport-red/30">
        <span class="text-white font-black text-xl tracking-tighter">FB</span>
      </div>
      <h1 class="text-2xl font-black tracking-tight pt-2">เข้าสู่ระบบ</h1>
      <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        เข้าสู่ระบบ FOOTB เพื่อจัดการทีมโปรดและติดตามไทยลีก
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
        <label class="text-xs font-semibold block">อีเมล หรือ ชื่อผู้ใช้</label>
        <input
          v-model="email"
          type="text"
          required
          autocomplete="username"
          placeholder="somchai@example.com หรือ somchai99"
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
            autocomplete="current-password"
            placeholder="กรอกรหัสผ่านของคุณ"
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

      <button
        type="submit"
        :disabled="authStore.loading"
        class="w-full btn-primary py-2.5 text-sm font-bold disabled:opacity-50"
      >
        {{ authStore.loading ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ' }}
      </button>
    </form>

    <!-- Footer Link -->
    <div
      class="pt-4 border-t text-center text-xs"
      :class="isDark ? 'border-sport-border text-neutral-400' : 'border-gray-200 text-gray-600'"
    >
      ยังไม่มีบัญชีใช่ไหม?
      <RouterLink to="/register" class="font-bold text-sport-red hover:underline ml-1">
        สมัครสมาชิกเลย
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const errorMsg = ref('')

async function handleSubmit() {
  errorMsg.value = ''
  try {
    await authStore.login({ email: email.value.trim(), password: password.value })
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.push(redirect)
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: string } } }
    errorMsg.value = axiosErr.response?.data?.error || 'อีเมล/ชื่อผู้ใช้ หรือรหัสผ่านไม่ถูกต้อง'
  }
}
</script>
