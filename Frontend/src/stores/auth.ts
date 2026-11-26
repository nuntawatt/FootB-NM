import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { loginUser, registerUser, fetchProfile, type LoginPayload, type RegisterPayload } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('auth_token'))
  const user = ref<User | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => Boolean(token.value && user.value))

  function setToken(newToken: string | null) {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('auth_token', newToken)
    } else {
      localStorage.removeItem('auth_token')
    }
  }

  async function login(payload: LoginPayload) {
    loading.value = true
    try {
      const res = await loginUser(payload)
      setToken(res.token)
      user.value = res.user
      return res
    } finally {
      loading.value = false
    }
  }

  async function register(payload: RegisterPayload) {
    loading.value = true
    try {
      const res = await registerUser(payload)
      setToken(res.token)
      user.value = res.user
      return res
    } finally {
      loading.value = false
    }
  }

  async function loadProfile() {
    if (!token.value) return
    try {
      user.value = await fetchProfile()
    } catch {
      // Token expired or invalid
      logout()
    }
  }

  function logout() {
    setToken(null)
    user.value = null
  }

  return { token, user, loading, isAuthenticated, login, register, loadProfile, logout }
})
