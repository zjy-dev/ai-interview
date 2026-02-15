<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useTheme } from "@/composables/useTheme";

const auth = useAuthStore();
const router = useRouter();
const { theme, toggle } = useTheme();

const email = ref("");
const password = ref("");
const error = ref("");
const loading = ref(false);

async function handleLogin() {
  error.value = "";
  loading.value = true;
  try {
    await auth.login(email.value, password.value);
    router.push("/");
  } catch (e: any) {
    error.value = e.response?.data?.error || "登录失败";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="auth-page">
    <button class="theme-fab" @click="toggle" :title="theme === 'light' ? '夜间模式' : '日间模式'">
      <svg v-if="theme === 'light'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
      </svg>
      <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
      </svg>
    </button>

    <div class="auth-container">
      <div class="auth-brand">
        <span class="brand-diamond">◆</span>
        <h1>AI Interview</h1>
        <p class="tagline">LLM 驱动的模拟面试平台</p>
      </div>

      <div class="auth-card card">
        <div class="card-header">
          <h2>欢迎回来</h2>
          <p>登录以继续你的面试之旅</p>
        </div>

        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label>邮箱地址</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="4" width="20" height="16" rx="2"/><polyline points="22,6 12,13 2,6"/></svg>
              </span>
              <input v-model="email" type="email" class="form-control has-icon" placeholder="your@email.com" required />
            </div>
          </div>

          <div class="form-group">
            <label>密码</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
              </span>
              <input v-model="password" type="password" class="form-control has-icon" placeholder="输入密码" required />
            </div>
          </div>

          <Transition name="shake">
            <p v-if="error" class="error-msg">{{ error }}</p>
          </Transition>

          <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            {{ loading ? "登录中..." : "登 录" }}
          </button>
        </form>

        <div class="divider"><span>或</span></div>

        <p class="switch-text">
          还没有账号？
          <RouterLink to="/register">立即注册</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 24px;
  position: relative;
}

.theme-fab {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 50;
  width: 42px;
  height: 42px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: var(--surface);
  backdrop-filter: var(--backdrop);
  cursor: pointer;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}
.theme-fab:hover {
  border-color: var(--primary);
  color: var(--primary);
  box-shadow: 0 0 16px var(--primary-glow);
}

.auth-container {
  width: 100%;
  max-width: 420px;
  animation: slideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(24px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.auth-brand {
  text-align: center;
  margin-bottom: 32px;
}

.brand-diamond {
  display: inline-block;
  font-size: 36px;
  color: var(--primary);
  filter: drop-shadow(0 0 12px var(--primary-glow));
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

.auth-brand h1 {
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 32px;
  font-weight: 800;
  letter-spacing: -0.04em;
  margin-top: 8px;
}

.tagline {
  color: var(--text-muted);
  font-size: 14px;
  margin-top: 4px;
}

.card-header {
  margin-bottom: 28px;
}

.card-header h2 {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 4px;
}

.card-header p {
  color: var(--text-secondary);
  font-size: 14px;
}

.input-wrapper {
  position: relative;
}

.input-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  display: flex;
  pointer-events: none;
  transition: color 0.2s;
}

.input-wrapper:focus-within .input-icon {
  color: var(--primary);
}

.form-control.has-icon {
  padding-left: 42px;
}

.error-msg {
  color: var(--danger);
  font-size: 13px;
  margin-bottom: 16px;
  padding: 10px 14px;
  background: var(--danger-soft);
  border-radius: var(--radius);
  border: 1px solid rgba(239, 68, 68, 0.15);
}

.btn-block {
  width: 100%;
  height: 48px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.05em;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.divider {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 24px 0;
  color: var(--text-muted);
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border);
}

.switch-text {
  text-align: center;
  font-size: 14px;
  color: var(--text-secondary);
}

.switch-text a {
  font-weight: 600;
}

/* Shake animation for error */
.shake-enter-active { animation: shake 0.4s; }
@keyframes shake {
  10%, 90% { transform: translateX(-1px); }
  20%, 80% { transform: translateX(2px); }
  30%, 50%, 70% { transform: translateX(-3px); }
  40%, 60% { transform: translateX(3px); }
}
</style>
