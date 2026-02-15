import client from './client'

export interface Interview {
  id: number
  title: string
  position: string
  status: string
  language: string
  created_at: string
  websocket_url?: string
}

export interface InterviewMessage {
  id: number
  role: 'user' | 'assistant' | 'system'
  content: string
  created_at: string
}

export interface Evaluation {
  id: number
  interview_id: number
  overall_score: number
  summary: string
  categories: { category: string; score: number; comment: string }[]
  strengths: string
  weaknesses: string
  suggestions: string
  created_at: string
}

export interface CreateInterviewPayload {
  title: string
  position: string
  language: string
  llm_provider?: string
  llm_model?: string
  tts_provider?: string
  tts_voice?: string
  resume?: string
}

export const interviewApi = {
  create(payload: CreateInterviewPayload) {
    return client.post<Interview & { websocket_url: string }>('/interviews', payload)
  },
  list(page = 1, pageSize = 20) {
    return client.get<{ interviews: Interview[]; total: number }>('/interviews', {
      params: { page, page_size: pageSize },
    })
  },
  get(id: number) {
    return client.get<Interview & { messages: InterviewMessage[] }>(`/interviews/${id}`)
  },
  sendMessage(id: number, content: string) {
    return client.post<{
      user_message: InterviewMessage
      assistant_message: { role: string; content: string }
    }>(`/interviews/${id}/messages`, { content })
  },
  end(id: number) {
    return client.post<{ status: string; evaluation_summary: string }>(
      `/interviews/${id}/end`,
    )
  },
  getEvaluation(id: number) {
    return client.get<Evaluation>(`/interviews/${id}/evaluation`)
  },
}
