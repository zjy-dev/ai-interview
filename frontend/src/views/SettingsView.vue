<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();

const form = ref({
  llm_provider: "",
  llm_api_key: "",
  llm_base_url: "",
  llm_model: "",
  tts_provider: "",
  tts_api_key: "",
  tts_voice: "",
  tts_enabled: true,
  stt_provider: "browser",
  stt_api_key: "",
});

const saving = ref(false);
const message = ref("");

const llmProviders = [
  { value: "openai", label: "OpenAI" },
  { value: "anthropic", label: "Anthropic (Claude)" },
  { value: "deepseek", label: "DeepSeek" },
  { value: "gemini", label: "Google Gemini" },
  { value: "custom", label: "自定义 (OpenAI 兼容)" },
];

const ttsProviders = [
  { value: "openai", label: "OpenAI TTS" },
  { value: "fishaudio", label: "Fish Audio" },
  { value: "elevenlabs", label: "ElevenLabs" },
  { value: "edgetts", label: "Edge TTS (免费)" },
];

const sttProviders = [
  { value: "browser", label: "浏览器语音识别 (免费)" },
  { value: "whisper", label: "OpenAI Whisper" },
];

onMounted(async () => {
  await auth.fetchSettings();
  if (auth.settings) {
    form.value.llm_provider = auth.settings.llm_provider;
    form.value.llm_base_url = auth.settings.llm_base_url;
    form.value.llm_model = auth.settings.llm_model;
    form.value.tts_provider = auth.settings.tts_provider;
    form.value.tts_voice = auth.settings.tts_voice;
    form.value.tts_enabled = auth.settings.tts_enabled;
    form.value.stt_provider = auth.settings.stt_provider;
  }
});

async function saveSettings() {
  saving.value = true;
  message.value = "";
  try {
    await auth.updateSettings(form.value);
    message.value = "设置已保存";
    form.value.llm_api_key = "";
    form.value.tts_api_key = "";
    form.value.stt_api_key = "";
  } catch (e: any) {
    message.value = "保存失败: " + (e.response?.data?.error || e.message);
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <div class="settings-view">
    <div class="settings-header">
      <h1>设置</h1>
      <p class="settings-subtitle">
        配置 LLM、TTS 和 STT 服务参数，API Key 使用 BYOK 模式安全存储
      </p>
    </div>

    <form @submit.prevent="saveSettings" class="settings-form">
      <!-- LLM -->
      <div class="card section-card">
        <div class="section-top">
          <div class="section-icon icon-llm">
            <svg
              width="20"
              height="20"
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
          <div class="section-label">LLM 大语言模型</div>
        </div>

        <div class="field-grid">
          <div class="form-group">
            <label>提供商</label>
            <select v-model="form.llm_provider" class="form-control">
              <option value="" disabled>选择 LLM 提供商</option>
              <option v-for="p in llmProviders" :key="p.value" :value="p.value">
                {{ p.label }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>模型名称</label>
            <input
              v-model="form.llm_model"
              class="form-control"
              placeholder="例如：gpt-4o, claude-sonnet-4-20250514"
            />
          </div>
        </div>

        <div class="form-group">
          <label>
            API Key
            <span v-if="auth.settings?.llm_api_key_set" class="key-badge"
              >已设置</span
            >
          </label>
          <input
            v-model="form.llm_api_key"
            type="password"
            class="form-control"
            :placeholder="
              auth.settings?.llm_api_key_set ? '留空保持不变' : '输入 API Key'
            "
          />
        </div>

        <div v-if="form.llm_provider === 'custom'" class="form-group">
          <label>Base URL</label>
          <input
            v-model="form.llm_base_url"
            class="form-control"
            placeholder="https://api.example.com/v1"
          />
        </div>
      </div>

      <!-- TTS -->
      <div class="card section-card">
        <div class="section-top">
          <div class="section-icon icon-tts">
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5" />
              <path d="M19.07 4.93a10 10 0 0 1 0 14.14" />
              <path d="M15.54 8.46a5 5 0 0 1 0 7.07" />
            </svg>
          </div>
          <div class="section-label">TTS 文字转语音</div>
        </div>

        <label class="toggle-row">
          <span class="toggle-text">启用语音播放</span>
          <button
            type="button"
            :class="['toggle-switch', { active: form.tts_enabled }]"
            @click="form.tts_enabled = !form.tts_enabled"
          >
            <span class="toggle-knob" />
          </button>
        </label>

        <template v-if="form.tts_enabled">
          <div class="field-grid">
            <div class="form-group">
              <label>提供商</label>
              <select v-model="form.tts_provider" class="form-control">
                <option value="" disabled>选择 TTS 提供商</option>
                <option
                  v-for="p in ttsProviders"
                  :key="p.value"
                  :value="p.value"
                >
                  {{ p.label }}
                </option>
              </select>
            </div>
            <div class="form-group">
              <label>音色 / Voice ID</label>
              <input
                v-model="form.tts_voice"
                class="form-control"
                placeholder="例如：alloy, shimmer"
              />
            </div>
          </div>

          <div v-if="form.tts_provider !== 'edgetts'" class="form-group">
            <label>
              API Key
              <span v-if="auth.settings?.tts_api_key_set" class="key-badge"
                >已设置</span
              >
            </label>
            <input
              v-model="form.tts_api_key"
              type="password"
              class="form-control"
              :placeholder="
                auth.settings?.tts_api_key_set ? '留空保持不变' : '输入 API Key'
              "
            />
          </div>
        </template>
      </div>

      <!-- STT -->
      <div class="card section-card">
        <div class="section-top">
          <div class="section-icon icon-stt">
            <svg
              width="20"
              height="20"
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
          </div>
          <div class="section-label">STT 语音识别</div>
        </div>

        <div class="form-group">
          <label>提供商</label>
          <select v-model="form.stt_provider" class="form-control">
            <option v-for="p in sttProviders" :key="p.value" :value="p.value">
              {{ p.label }}
            </option>
          </select>
        </div>

        <div v-if="form.stt_provider === 'whisper'" class="form-group">
          <label>
            API Key
            <span v-if="auth.settings?.stt_api_key_set" class="key-badge"
              >已设置</span
            >
          </label>
          <input
            v-model="form.stt_api_key"
            type="password"
            class="form-control"
            :placeholder="
              auth.settings?.stt_api_key_set
                ? '留空保持不变'
                : '输入 OpenAI API Key'
            "
          />
        </div>
      </div>

      <!-- Message & Submit -->
      <Transition name="fade">
        <p
          v-if="message"
          :class="['msg-toast', { 'is-success': message === '设置已保存' }]"
        >
          <svg
            v-if="message === '设置已保存'"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
            <polyline points="22 4 12 14.01 9 11.01" />
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
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          {{ message }}
        </p>
      </Transition>

      <button type="submit" class="btn btn-primary btn-save" :disabled="saving">
        <svg
          v-if="saving"
          class="spin"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M21 12a9 9 0 1 1-6.219-8.56" />
        </svg>
        {{ saving ? "保存中..." : "保存设置" }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.settings-view {
  max-width: 660px;
  margin: 0 auto;
  animation: fadeInUp 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.settings-header {
  margin-bottom: 32px;
}

.settings-header h1 {
  font-size: 28px;
  font-weight: 800;
}

.settings-subtitle {
  font-size: 14px;
  color: var(--text-muted);
  margin-top: 6px;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* Section card */
.section-card {
  padding: 28px 32px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-top {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.section-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-llm {
  background: var(--primary-soft);
  color: var(--primary);
}
.icon-tts {
  background: var(--success-soft);
  color: var(--success);
}
.icon-stt {
  background: var(--warning-soft);
  color: var(--warning);
}

.section-label {
  font-family: "Outfit", system-ui, sans-serif;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--primary);
}

/* Field grid */
.field-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

/* Form group */
.form-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.key-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  padding: 1px 8px;
  border-radius: var(--radius-full);
  background: var(--success-soft);
  color: var(--success);
}

/* Toggle */
.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  padding: 12px 16px;
  background: var(--surface-hover);
  border-radius: var(--radius);
}

.toggle-text {
  font-size: 14px;
  font-weight: 500;
  color: var(--text);
}

.toggle-switch {
  position: relative;
  width: 44px;
  height: 24px;
  border-radius: 12px;
  border: none;
  background: var(--border-strong);
  cursor: pointer;
  transition: background 0.25s;
  padding: 0;
}

.toggle-switch.active {
  background: var(--primary);
}

.toggle-knob {
  position: absolute;
  top: 3px;
  left: 3px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.toggle-switch.active .toggle-knob {
  transform: translateX(20px);
}

/* Message toast */
.msg-toast {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: var(--radius);
  font-size: 14px;
  font-weight: 500;
  background: var(--danger-soft);
  color: var(--danger);
}

.msg-toast.is-success {
  background: var(--success-soft);
  color: var(--success);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Save button */
.btn-save {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 14px;
  font-size: 15px;
}

.spin {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 640px) {
  .section-card {
    padding: 20px;
  }
  .field-grid {
    grid-template-columns: 1fr;
  }
}
</style>
