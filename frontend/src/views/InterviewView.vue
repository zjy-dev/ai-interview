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
    <div class="interview-header">
      <div>
        <h2>{{ store.current?.title || "面试" }}</h2>
        <span class="position">{{ store.current?.position }}</span>
      </div>
      <div class="header-actions">
        <button
          class="btn btn-secondary"
          @click="audio.enabled.value = !audio.enabled.value"
        >
          {{ audio.enabled.value ? "🔊 语音开" : "🔇 语音关" }}
        </button>
        <button
          v-if="store.current?.status === 'active'"
          class="btn btn-danger"
          @click="endInterview"
        >
          结束面试
        </button>
        <button
          v-else
          class="btn btn-secondary"
          @click="router.push(`/interviews/${interviewId}/report`)"
        >
          查看报告
        </button>
      </div>
    </div>

    <div ref="chatContainer" class="chat-container card">
      <div
        v-for="msg in store.messages"
        :key="msg.id || msg.content"
        :class="['message', `message-${msg.role}`]"
      >
        <div class="message-avatar">
          {{ msg.role === "user" ? "👤" : "🤖" }}
        </div>
        <div class="message-content">
          <div class="message-role">
            {{ msg.role === "user" ? "你" : "面试官" }}
          </div>
          <div class="message-text">{{ msg.content }}</div>
        </div>
      </div>

      <div v-if="sending" class="message message-assistant">
        <div class="message-avatar">🤖</div>
        <div class="message-content">
          <div class="message-role">面试官</div>
          <div class="message-text typing">正在思考中...</div>
        </div>
      </div>
    </div>

    <div v-if="store.current?.status === 'active'" class="input-area">
      <div class="input-row">
        <textarea
          v-model="inputText"
          class="form-control input-textarea"
          :placeholder="
            inputMode === 'voice'
              ? '点击麦克风开始语音输入...'
              : '输入你的回答...'
          "
          rows="2"
          @keydown="handleKeydown"
        />
        <div class="input-actions">
          <button
            v-if="speech.isSupported.value"
            :class="[
              'btn',
              'btn-icon',
              { recording: speech.isListening.value },
            ]"
            @click="toggleVoice"
          >
            {{ speech.isListening.value ? "⏹️" : "🎤" }}
          </button>
          <button
            class="btn btn-primary"
            :disabled="!inputText.trim() || sending"
            @click="sendMessage"
          >
            发送
          </button>
        </div>
      </div>
      <p v-if="speech.isListening.value" class="voice-hint">
        正在录音... {{ speech.transcript.value }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.interview-view {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 104px);
}

.interview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.interview-header h2 {
  font-size: 20px;
}

.position {
  color: var(--text-secondary);
  font-size: 13px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.chat-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message {
  display: flex;
  gap: 12px;
  max-width: 80%;
}

.message-user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-assistant {
  align-self: flex-start;
}

.message-avatar {
  font-size: 24px;
  flex-shrink: 0;
}

.message-content {
  background: #f1f5f9;
  padding: 10px 14px;
  border-radius: 12px;
}

.message-user .message-content {
  background: var(--primary);
  color: white;
}

.message-role {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 4px;
  opacity: 0.7;
}

.message-text {
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.typing {
  animation: blink 1.4s infinite;
}

@keyframes blink {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

.input-area {
  margin-top: 16px;
}

.input-row {
  display: flex;
  gap: 8px;
  align-items: flex-end;
}

.input-textarea {
  flex: 1;
  resize: none;
}

.input-actions {
  display: flex;
  gap: 8px;
}

.btn-icon {
  width: 40px;
  height: 40px;
  padding: 0;
  font-size: 18px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  cursor: pointer;
}

.recording {
  background: #fef2f2;
  border-color: #ef4444;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%,
  100% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.4);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(239, 68, 68, 0);
  }
}

.voice-hint {
  font-size: 13px;
  color: #ef4444;
  margin-top: 8px;
}
</style>
