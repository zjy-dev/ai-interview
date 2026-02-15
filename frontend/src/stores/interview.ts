import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  interviewApi,
  type Interview,
  type InterviewMessage,
  type Evaluation,
  type CreateInterviewPayload,
} from '@/api'

export const useInterviewStore = defineStore('interview', () => {
  const interviews = ref<Interview[]>([])
  const total = ref(0)
  const current = ref<Interview | null>(null)
  const messages = ref<InterviewMessage[]>([])
  const evaluation = ref<Evaluation | null>(null)
  const loading = ref(false)

  async function fetchList(page = 1, pageSize = 20) {
    loading.value = true
    try {
      const { data } = await interviewApi.list(page, pageSize)
      interviews.value = data.interviews
      total.value = data.total
    } finally {
      loading.value = false
    }
  }

  async function fetchInterview(id: number) {
    loading.value = true
    try {
      const { data } = await interviewApi.get(id)
      current.value = data
      messages.value = data.messages ?? []
    } finally {
      loading.value = false
    }
  }

  async function create(payload: CreateInterviewPayload) {
    const { data } = await interviewApi.create(payload)
    return data
  }

  async function sendMessage(id: number, content: string) {
    const { data } = await interviewApi.sendMessage(id, content)
    messages.value.push({
      id: data.user_message.id,
      role: 'user',
      content: data.user_message.content,
      created_at: new Date().toISOString(),
    })
    messages.value.push({
      id: 0,
      role: 'assistant',
      content: data.assistant_message.content,
      created_at: new Date().toISOString(),
    })
    return data
  }

  async function endInterview(id: number) {
    const { data } = await interviewApi.end(id)
    if (current.value) {
      current.value.status = 'completed'
    }
    return data
  }

  async function fetchEvaluation(id: number) {
    const { data } = await interviewApi.getEvaluation(id)
    evaluation.value = data
  }

  return {
    interviews,
    total,
    current,
    messages,
    evaluation,
    loading,
    fetchList,
    fetchInterview,
    create,
    sendMessage,
    endInterview,
    fetchEvaluation,
  }
})
