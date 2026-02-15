<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();
const router = useRouter();

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
    <div class="auth-card card">
      <h2>登录 AI Interview</h2>
      <p class="subtitle">使用你的账号登录</p>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label>邮箱</label>
          <input
            v-model="email"
            type="email"
            class="form-control"
            placeholder="your@email.com"
            required
          />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input
            v-model="password"
            type="password"
            class="form-control"
            placeholder="输入密码"
            required
          />
        </div>
        <p v-if="error" class="error">{{ error }}</p>
        <button
          type="submit"
          class="btn btn-primary btn-block"
          :disabled="loading"
        >
          {{ loading ? "登录中..." : "登录" }}
        </button>
      </form>

      <p class="footer-text">
        还没有账号？<RouterLink to="/register">立即注册</RouterLink>
      </p>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.auth-card {
  width: 100%;
  max-width: 400px;
}

.auth-card h2 {
  font-size: 24px;
  margin-bottom: 4px;
}

.subtitle {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 24px;
}

.error {
  color: #ef4444;
  font-size: 13px;
  margin-bottom: 12px;
}

.btn-block {
  width: 100%;
}

.footer-text {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
  color: var(--text-secondary);
}

.footer-text a {
  color: var(--primary);
}
</style>
