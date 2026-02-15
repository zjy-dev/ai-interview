<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();
const router = useRouter();

const email = ref("");
const password = ref("");
const nickname = ref("");
const error = ref("");
const loading = ref(false);

async function handleRegister() {
  error.value = "";
  loading.value = true;
  try {
    await auth.register(email.value, password.value, nickname.value);
    router.push("/");
  } catch (e: any) {
    error.value = e.response?.data?.error || "注册失败";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card card">
      <h2>注册账号</h2>
      <p class="subtitle">创建一个新的 AI Interview 账号</p>

      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label>昵称</label>
          <input
            v-model="nickname"
            type="text"
            class="form-control"
            placeholder="你的昵称"
            required
          />
        </div>
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
            placeholder="至少 8 位"
            required
            minlength="8"
          />
        </div>
        <p v-if="error" class="error">{{ error }}</p>
        <button
          type="submit"
          class="btn btn-primary btn-block"
          :disabled="loading"
        >
          {{ loading ? "注册中..." : "注册" }}
        </button>
      </form>

      <p class="footer-text">
        已有账号？<RouterLink to="/login">立即登录</RouterLink>
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
