import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type UserProfile, type UserSettings, type UpdateSettingsPayload } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<UserProfile | null>(null)
  const settings = ref<UserSettings | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  async function login(email: string, password: string) {
    const { data } = await authApi.login(email, password)
    token.value = data.token
    localStorage.setItem('token', data.token)
    user.value = { id: data.id, email, nickname: data.nickname, created_at: '' }
  }

  async function register(email: string, password: string, nickname: string) {
    const { data } = await authApi.register(email, password, nickname)
    token.value = data.token
    localStorage.setItem('token', data.token)
    user.value = { id: data.id, email, nickname, created_at: '' }
  }

  function logout() {
    token.value = ''
    user.value = null
    settings.value = null
    localStorage.removeItem('token')
  }

  async function fetchProfile() {
    const { data } = await authApi.getProfile()
    user.value = data
  }

  async function fetchSettings() {
    const { data } = await authApi.getSettings()
    settings.value = data
  }

  async function updateSettings(payload: UpdateSettingsPayload) {
    await authApi.updateSettings(payload)
    await fetchSettings()
  }

  return {
    token,
    user,
    settings,
    isAuthenticated,
    login,
    register,
    logout,
    fetchProfile,
    fetchSettings,
    updateSettings,
  }
})
