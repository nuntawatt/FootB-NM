<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
    @click.self="$emit('close')"
  >
    <div class="card w-full max-w-md p-6 space-y-5 border-t-4 border-t-sport-red">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <h2 class="text-xl font-extrabold">
          {{ mode === 'login' ? 'เข้าสู่ระบบ' : 'สมัครสมาชิก' }}
        </h2>
        <button
          type="button"
          class="text-sm px-2 py-1 rounded hover:bg-sport-hover"
          :class="isDark ? 'text-neutral-400 hover:text-white' : 'text-gray-500 hover:text-gray-900'"
          @click="$emit('close')"
        >
          ✕
        </button>
      </div>

      <!-- Mode Tabs -->
      <div class="grid grid-cols-2 gap-2 p-1 rounded-lg" :class="isDark ? 'bg-sport-bg' : 'bg-gray-100'">
        <button
          type="button"
          @click="mode = 'login'; errorMsg = ''"
          class="py-1.5 rounded-md text-xs font-bold transition-colors"
          :class="mode === 'login' ? 'bg-sport-red text-white' : (isDark ? 'text-neutral-400' : 'text-gray-600')"
        >
          เข้าสู่ระบบ
        </button>
        <button
          type="button"
          @click="mode = 'register'; errorMsg = ''"
          class="py-1.5 rounded-md text-xs font-bold transition-colors"
          :class="mode === 'register' ? 'bg-sport-red text-white' : (isDark ? 'text-neutral-400' : 'text-gray-600')"
        >
          สมัครสมาชิก
        </button>
      </div>

      <!-- Error -->
      <div v-if="errorMsg" class="p-3 rounded-lg bg-sport-red/15 border border-sport-red text-xs text-sport-red font-medium">
        {{ errorMsg }}
      </div>

      <!-- Form -->
      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div v-if="mode === 'register'" class="space-y-1">
          <label class="text-xs font-semibold">ชื่อผู้ใช้ (Username)</label>
          <input
            v-model="username"
            type="text"
            required
            placeholder="a-z, 0-9, _ (3-30 ตัวอักษร)"
            class="w-full px-3 py-2 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
            :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-gray-50 border-gray-300 text-gray-900'"
          />
        </div>

        <div class="space-y-1">
          <label class="text-xs font-semibold">
            {{ mode === 'login' ? 'อีเมล หรือ ชื่อผู้ใช้' : 'อีเมล' }}
          </label>
          <input
            v-model="email"
            :type="mode === 'register' ? 'email' : 'text'"
            required
            placeholder="example@email.com"
            class="w-full px-3 py-2 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
            :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-gray-50 border-gray-300 text-gray-900'"
          />
        </div>

        <div class="space-y-1">
          <label class="text-xs font-semibold">รหัสผ่าน</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="อย่างน้อย 6 ตัว มีตัวอักษรและตัวเลข"
            class="w-full px-3 py-2 rounded-lg text-sm border outline-none focus:ring-2 focus:ring-sport-red"
            :class="isDark ? 'bg-sport-bg border-sport-border text-white' : 'bg-gray-50 border-gray-300 text-gray-900'"
          />
        </div>

        <button
          type="submit"
          :disabled="authStore.loading"
          class="w-full btn-primary py-2.5 text-sm disabled:opacity-50"
        >
          {{ authStore.loading ? 'กำลังดำเนินการ...' : (mode === 'login' ? 'เข้าสู่ระบบ' : 'สมัครสมาชิก') }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'

defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const authStore = useAuthStore()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const mode = ref<'login' | 'register'>('login')
const username = ref('')
const email = ref('')
const password = ref('')
const errorMsg = ref('')

async function handleSubmit() {
  errorMsg.value = ''
  try {
    if (mode.value === 'login') {
      await authStore.login({ email: email.value, password: password.value })
    } else {
      await authStore.register({
        username: username.value,
        email: email.value,
        password: password.value,
      })
    }
    emit('close')
    password.value = ''
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: string } } }
    errorMsg.value = axiosErr.response?.data?.error || 'ไม่สามารถดำเนินการได้ กรุณาลองใหม่'
  }
}
</script>
