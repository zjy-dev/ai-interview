import client from './client'

export interface LoginResponse {
  id: number
  token: string
  nickname: string
}

export interface UserProfile {
  id: number
  email: string
  nickname: string
  created_at: string
}

export interface UserSettings {
  llm_provider: string
  llm_api_key_set: boolean
  llm_base_url: string
  llm_model: string
  tts_provider: string
  tts_api_key_set: boolean
  tts_voice: string
  tts_enabled: boolean
  stt_provider: string
  stt_api_key_set: boolean
}

export interface UpdateSettingsPayload {
  llm_provider?: string
  llm_api_key?: string
  llm_base_url?: string
  llm_model?: string
  tts_provider?: string
  tts_api_key?: string
  tts_voice?: string
  tts_enabled?: boolean
  stt_provider?: string
  stt_api_key?: string
}

export const authApi = {
  register(email: string, password: string, nickname: string) {
    return client.post<LoginResponse>('/auth/register', { email, password, nickname })
  },
  login(email: string, password: string) {
    return client.post<LoginResponse>('/auth/login', { email, password })
  },
  getProfile() {
    return client.get<UserProfile>('/auth/profile')
  },
  getSettings() {
    return client.get<UserSettings>('/auth/settings')
  },
  updateSettings(payload: UpdateSettingsPayload) {
    return client.put<{ success: boolean }>('/auth/settings', payload)
  },
}
