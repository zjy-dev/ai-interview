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
    // 清空 API key 输入 (不再需要重复发送)
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
    <h1>设置</h1>
    <p class="subtitle">配置你的 LLM、TTS 和 STT 服务</p>

    <form @submit.prevent="saveSettings">
      <!-- LLM 设置 -->
      <div class="section card">
        <h3>🧠 LLM 大语言模型</h3>

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
          <label
            >API Key
            {{
              auth.settings?.llm_api_key_set ? "(已设置，留空保持不变)" : ""
            }}</label
          >
          <input
            v-model="form.llm_api_key"
            type="password"
            class="form-control"
            :placeholder="
              auth.settings?.llm_api_key_set ? '••••••••' : '输入 API Key'
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

        <div class="form-group">
          <label>模型名称</label>
          <input
            v-model="form.llm_model"
            class="form-control"
            placeholder="例如：gpt-4o, claude-sonnet-4-20250514"
          />
        </div>
      </div>

      <!-- TTS 设置 -->
      <div class="section card">
        <h3>🔊 TTS 文字转语音</h3>

        <div class="form-group">
          <label class="toggle-label">
            <input type="checkbox" v-model="form.tts_enabled" />
            <span>启用语音播放</span>
          </label>
        </div>

        <template v-if="form.tts_enabled">
          <div class="form-group">
            <label>提供商</label>
            <select v-model="form.tts_provider" class="form-control">
              <option value="" disabled>选择 TTS 提供商</option>
              <option v-for="p in ttsProviders" :key="p.value" :value="p.value">
                {{ p.label }}
              </option>
            </select>
          </div>

          <div v-if="form.tts_provider !== 'edgetts'" class="form-group">
            <label
              >API Key
              {{ auth.settings?.tts_api_key_set ? "(已设置)" : "" }}</label
            >
            <input
              v-model="form.tts_api_key"
              type="password"
              class="form-control"
              :placeholder="
                auth.settings?.tts_api_key_set ? '••••••••' : '输入 API Key'
              "
            />
          </div>

          <div class="form-group">
            <label>音色 / Voice ID</label>
            <input
              v-model="form.tts_voice"
              class="form-control"
              placeholder="例如：alloy, shimmer"
            />
          </div>
        </template>
      </div>

      <!-- STT 设置 -->
      <div class="section card">
        <h3>🎤 STT 语音识别</h3>

        <div class="form-group">
          <label>提供商</label>
          <select v-model="form.stt_provider" class="form-control">
            <option v-for="p in sttProviders" :key="p.value" :value="p.value">
              {{ p.label }}
            </option>
          </select>
        </div>

        <div v-if="form.stt_provider === 'whisper'" class="form-group">
          <label
            >API Key
            {{ auth.settings?.stt_api_key_set ? "(已设置)" : "" }}</label
          >
          <input
            v-model="form.stt_api_key"
            type="password"
            class="form-control"
            :placeholder="
              auth.settings?.stt_api_key_set
                ? '••••••••'
                : '输入 OpenAI API Key'
            "
          />
        </div>
      </div>

      <p
        v-if="message"
        :class="['message', { success: message === '设置已保存' }]"
      >
        {{ message }}
      </p>

      <button type="submit" class="btn btn-primary" :disabled="saving">
        {{ saving ? "保存中..." : "保存设置" }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.settings-view {
  max-width: 600px;
}

.settings-view h1 {
  font-size: 24px;
  margin-bottom: 4px;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 24px;
}

.section {
  margin-bottom: 20px;
}

.section h3 {
  font-size: 16px;
  margin-bottom: 16px;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.toggle-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
}

.message {
  font-size: 14px;
  margin-bottom: 12px;
  color: #ef4444;
}

.message.success {
  color: #15803d;
}
</style>
