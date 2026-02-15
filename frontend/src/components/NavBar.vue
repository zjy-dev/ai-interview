<script setup lang="ts">
import { RouterLink, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();
const router = useRouter();

function logout() {
  auth.logout();
  router.push("/login");
}
</script>

<template>
  <nav class="navbar">
    <div class="navbar-inner">
      <RouterLink to="/" class="brand">🎯 AI Interview</RouterLink>
      <div class="nav-links">
        <RouterLink to="/interviews">面试列表</RouterLink>
        <RouterLink to="/interviews/new">新建面试</RouterLink>
        <RouterLink to="/settings">设置</RouterLink>
        <button class="btn btn-secondary btn-sm" @click="logout">退出</button>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  padding: 0 24px;
}

.navbar-inner {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px;
}

.brand {
  font-size: 18px;
  font-weight: 700;
  color: var(--primary);
  text-decoration: none;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-links a {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  color: var(--primary);
}

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
}
</style>
