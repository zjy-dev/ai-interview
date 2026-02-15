<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useInterviewStore } from "@/stores/interview";
import { useAudioPlayer } from "@/composables/useAudioPlayer";
import { useSpeechRecognition } from "@/composables/useSpeechRecognition";
import { useWebSocket } from "@/composables/useWebSocket";

const route = useRoute();
const router = useRouter();
const store = useInterviewStore();
const audio = useAudioPlayer();
const speech = useSpeechRecognition();
const ws = useWebSocket();

const interviewId = computed(() => Number(route.params.id));
const inputText = ref("");
const inputMode = ref<"text" | "voice">("text");
const sending = ref(false);
const chatContainer = ref<HTMLElement | null>(null);

onMounted(async () => {
  await store.fetchInterview(interviewId.value);
  scrollToBottom();
});

function scrollToBottom() {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
    }
  });
}

async function sendMessage() {
  const content = inputText.value.trim();
  if (!content || sending.value) return;

  sending.value = true;
  inputText.value = "";
  try {
    await store.sendMessage(interviewId.value, content);
    scrollToBottom();
  } catch (e: any) {
    console.error("Send failed:", e);
  } finally {
    sending.value = false;
  }
}

function toggleVoice() {
  if (speech.isListening.value) {
    speech.stop();
    if (speech.transcript.value) {
      inputText.value = speech.transcript.value;
    }
  } else {
    const lang = store.current?.language === "en" ? "en-US" : "zh-CN";
    speech.start(lang);
  }
}

async function endInterview() {
  if (!confirm("确定结束面试吗？结束后将生成评估报告。")) return;
  try {
    await store.endInterview(interviewId.value);
    router.push(`/interviews/${interviewId.value}/report`);
  } catch (e: any) {
    console.error("End failed:", e);
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault();
    sendMessage();
  }
}
</script>

<template>
  <div class="interview-view">
    <!-- Header -->
    <div class="iv-header">
      <div class="iv-header-left">
        <button
          class="btn btn-ghost btn-sm"
          @click="router.push('/interviews')"
        >
          ← 返回
        </button>
        <div class="iv-title">
          <h2>{{ store.current?.title || "面试" }}</h2>
          <span class="iv-position">{{ store.current?.position }}</span>
        </div>
      </div>
      <div class="iv-header-right">
        <button
          class="btn-icon-toggle"
          @click="audio.enabled.value = !audio.enabled.value"
          :title="audio.enabled.value ? '关闭语音' : '开启语音'"
        >
          <svg
            v-if="audio.enabled.value"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" />
            <path d="M19.07 4.93a10 10 0 0 1 0 14.14" />
            <path d="M15.54 8.46a5 5 0 0 1 0 7.07" />
          </svg>
          <svg
            v-else
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" />
            <line x1="23" y1="9" x2="17" y2="15" />
            <line x1="17" y1="9" x2="23" y2="15" />
          </svg>
        </button>
        <button
          v-if="store.current?.status === 'active'"
          class="btn btn-danger btn-sm"
          @click="endInterview"
        >
          结束面试
        </button>
        <button
          v-else
          class="btn btn-secondary btn-sm"
          @click="router.push(`/interviews/${interviewId}/report`)"
        >
          查看报告
        </button>
      </div>
    </div>

    <!-- Chat area -->
    <div ref="chatContainer" class="chat-area">
      <div
        v-for="(msg, index) in store.messages"
        :key="msg.id || index"
        :class="['msg', `msg-${msg.role}`]"
        :style="{ animationDelay: `${Math.min(index * 30, 300)}ms` }"
      >
        <div class="msg-avatar" :class="`avatar-${msg.role}`">
          <svg
            v-if="msg.role === 'user'"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
            <circle cx="12" cy="7" r="4" />
          </svg>
          <svg
            v-else
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M12 2a4 4 0 0 1 4 4v2a4 4 0 0 1-8 0V6a4 4 0 0 1 4-4z" />
            <path d="M6 10a6 6 0 0 0 12 0" />
            <rect x="9" y="16" width="6" height="6" rx="1" />
          </svg>
        </div>
        <div class="msg-body">
          <div class="msg-meta">
            <span class="msg-sender">{{
              msg.role === "user" ? "你" : "面试官"
            }}</span>
          </div>
          <div class="msg-bubble">{{ msg.content }}</div>
        </div>
      </div>

      <!-- Typing indicator -->
      <div v-if="sending" class="msg msg-assistant">
        <div class="msg-avatar avatar-assistant">
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M12 2a4 4 0 0 1 4 4v2a4 4 0 0 1-8 0V6a4 4 0 0 1 4-4z" />
            <path d="M6 10a6 6 0 0 0 12 0" />
            <rect x="9" y="16" width="6" height="6" rx="1" />
          </svg>
        </div>
        <div class="msg-body">
          <div class="msg-meta">
            <span class="msg-sender">面试官</span>
          </div>
          <div class="msg-bubble typing-bubble">
            <span class="dot"></span>
            <span class="dot"></span>
            <span class="dot"></span>
          </div>
        </div>
      </div>
    </div>

    <!-- Input area -->
    <div v-if="store.current?.status === 'active'" class="input-area">
      <div class="input-card">
        <textarea
          v-model="inputText"
          class="chat-input"
          :placeholder="
            speech.isListening.value ? '正在录音...' : '输入你的回答...'
          "
          rows="1"
          @keydown="handleKeydown"
        />
        <div class="input-actions">
          <button
            v-if="speech.isSupported.value"
            :class="['btn-mic', { 'is-recording': speech.isListening.value }]"
            @click="toggleVoice"
            :title="speech.isListening.value ? '停止录音' : '语音输入'"
          >
            <svg
              v-if="!speech.isListening.value"
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z" />
              <path d="M19 10v2a7 7 0 0 1-14 0v-2" />
              <line x1="12" y1="19" x2="12" y2="23" />
              <line x1="8" y1="23" x2="16" y2="23" />
            </svg>
            <svg
              v-else
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <rect x="6" y="6" width="12" height="12" rx="2" />
            </svg>
          </button>
          <button
            class="btn-send"
            :disabled="!inputText.trim() || sending"
            @click="sendMessage"
          >
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <line x1="22" y1="2" x2="11" y2="13" />
              <polygon points="22 2 15 22 11 13 2 9 22 2" />
            </svg>
          </button>
        </div>
      </div>
      <p v-if="speech.isListening.value" class="voice-status">
        <span class="rec-dot"></span>
        录音中 — {{ speech.transcript.value || "等待语音..." }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.interview-view {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 96px);
  max-width: 860px;
  margin: 0 auto;
  width: 100%;
}

/* Header */
.iv-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  margin-bottom: 8px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.iv-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.iv-title h2 {
  font-size: 18px;
  font-weight: 700;
}

.iv-position {
  font-size: 13px;
  color: var(--text-muted);
}

.iv-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-sm {
  padding: 8px 14px;
  font-size: 13px;
}

.btn-icon-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-icon-toggle:hover {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--primary-soft);
}

/* Chat area */
.chat-area {
  flex: 1;
  overflow-y: auto;
  padding: 20px 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.msg {
  display: flex;
  gap: 12px;
  max-width: 85%;
  animation: msgIn 0.35s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes msgIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.msg-user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.msg-assistant {
  align-self: flex-start;
}

.msg-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.3s;
}

.avatar-user {
  background: var(--primary);
  color: #fff;
}

.avatar-assistant {
  background: var(--surface-solid);
  border: 1.5px solid var(--border-strong);
  color: var(--primary);
}

.msg-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.msg-user .msg-body {
  align-items: flex-end;
}

.msg-meta {
  padding: 0 4px;
}

.msg-sender {
  font-family: "Outfit", system-ui, sans-serif;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: var(--text-muted);
}

.msg-bubble {
  padding: 14px 18px;
  border-radius: 18px;
  font-size: 14px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
}

.msg-user .msg-bubble {
  background: var(--msg-user-bg);
  color: #fff;
  border-bottom-right-radius: 6px;
}

.msg-assistant .msg-bubble {
  background: var(--msg-ai-bg);
  backdrop-filter: var(--backdrop);
  color: var(--msg-ai-text);
  border: 1px solid var(--border);
  border-bottom-left-radius: 6px;
}

/* Typing indicator */
.typing-bubble {
  display: flex;
  gap: 5px;
  align-items: center;
  padding: 16px 22px;
}

.typing-bubble .dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text-muted);
  animation: bounce 1.4s ease-in-out infinite;
}

.typing-bubble .dot:nth-child(2) {
  animation-delay: 0.16s;
}
.typing-bubble .dot:nth-child(3) {
  animation-delay: 0.32s;
}

@keyframes bounce {
  0%,
  60%,
  100% {
    transform: translateY(0);
    opacity: 0.4;
  }
  30% {
    transform: translateY(-6px);
    opacity: 1;
  }
}

/* Input area */
.input-area {
  flex-shrink: 0;
  padding-top: 12px;
}

.input-card {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  background: var(--surface);
  backdrop-filter: var(--backdrop);
  border: 1.5px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 10px 10px 10px 18px;
  box-shadow: var(--shadow);
  transition:
    border-color 0.25s,
    box-shadow 0.25s;
}

.input-card:focus-within {
  border-color: var(--primary);
  box-shadow:
    0 0 0 3px var(--ring),
    var(--shadow);
}

.chat-input {
  flex: 1;
  border: none;
  background: none;
  font-family: "DM Sans", system-ui, sans-serif;
  font-size: 15px;
  color: var(--text);
  resize: none;
  outline: none;
  line-height: 1.5;
  max-height: 120px;
  padding: 6px 0;
}

.chat-input::placeholder {
  color: var(--text-muted);
}

.input-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.btn-mic {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius);
  border: none;
  background: var(--surface-hover);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-mic:hover {
  background: var(--primary-soft);
  color: var(--primary);
}

.btn-mic.is-recording {
  background: var(--danger-soft);
  color: var(--danger);
  animation: recPulse 1.5s infinite;
}

@keyframes recPulse {
  0%,
  100% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.3);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(239, 68, 68, 0);
  }
}

.btn-send {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius);
  border: none;
  background: var(--primary);
  color: #fff;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 8px var(--primary-glow);
}

.btn-send:hover:not(:disabled) {
  background: var(--primary-hover);
  transform: scale(1.05);
}

.btn-send:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.voice-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--danger);
  margin-top: 10px;
  padding-left: 4px;
}

.rec-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--danger);
  animation: recPulse 1.5s infinite;
}

@media (max-width: 640px) {
  .interview-view {
    height: calc(100vh - 84px);
  }
  .msg {
    max-width: 92%;
  }
  .iv-header {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
}
</style>
