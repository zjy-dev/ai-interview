<script setup lang="ts">
import { RouterLink, useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useTheme } from "@/composables/useTheme";

const auth = useAuthStore();
const router = useRouter();
const { theme, toggle } = useTheme();

function logout() {
  auth.logout();
  router.push("/login");
}
</script>

<template>
  <nav class="navbar">
    <div class="navbar-inner">
      <RouterLink to="/" class="brand">
        <span class="brand-icon">◆</span>
        <span class="brand-text">AI Interview</span>
      </RouterLink>

      <div class="nav-center">
        <RouterLink to="/interviews" class="nav-link">
          <span class="nav-icon">⊞</span>
          面试记录
        </RouterLink>
        <RouterLink to="/interviews/new" class="nav-link">
          <span class="nav-icon">＋</span>
          新建面试
        </RouterLink>
        <RouterLink to="/settings" class="nav-link">
          <span class="nav-icon">⚙</span>
          设置
        </RouterLink>
      </div>

      <div class="nav-end">
        <button class="theme-toggle" @click="toggle" :title="theme === 'light' ? '切换夜间模式' : '切换日间模式'">
          <span class="theme-icon" :class="{ 'is-dark': theme === 'dark' }">
            <svg v-if="theme === 'light'" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
            </svg>
            <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            </svg>
          </span>
        </button>
        <button class="btn-logout" @click="logout">退出</button>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--nav-bg);
  backdrop-filter: var(--backdrop);
  -webkit-backdrop-filter: var(--backdrop);
  border-bottom: 1px solid var(--border);
  transition: background 0.3s, border-color 0.3s;
}

.navbar-inner {
  max-width: 1100px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 24px;
}

/* Brand */
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  flex-shrink: 0;
}

.brand-icon {
  font-size: 22px;
  color: var(--primary);
  filter: drop-shadow(0 0 6px var(--primary-glow));
  transition: filter 0.3s;
}

.brand:hover .brand-icon {
  filter: drop-shadow(0 0 12px var(--primary-glow));
}

.brand-text {
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--text);
  transition: color 0.3s;
}

/* Nav links */
.nav-center {
  display: flex;
  align-items: center;
  gap: 4px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: var(--radius);
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.2s;
}

.nav-link:hover {
  color: var(--primary);
  background: var(--primary-soft);
}

.nav-link.router-link-active {
  color: var(--primary);
  background: var(--primary-soft);
  font-weight: 600;
}

.nav-icon {
  font-size: 15px;
  opacity: 0.7;
}

/* Right section */
.nav-end {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

/* Theme toggle */
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
  background: var(--surface);
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.3s;
}

.theme-toggle:hover {
  border-color: var(--primary);
  color: var(--primary);
  background: var(--primary-soft);
  box-shadow: 0 0 12px var(--primary-glow);
}

.theme-icon {
  display: flex;
  transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.theme-icon.is-dark {
  transform: rotate(180deg);
}

/* Logout button */
.btn-logout {
  padding: 8px 16px;
  border-radius: var(--radius);
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-muted);
  background: transparent;
  border: 1px solid var(--border);
  cursor: pointer;
  transition: all 0.2s;
}

.btn-logout:hover {
  color: var(--danger);
  border-color: var(--danger);
  background: var(--danger-soft);
}

/* Mobile */
@media (max-width: 640px) {
  .navbar-inner { padding: 0 16px; }
  .nav-center { gap: 0; }
  .nav-link { padding: 8px 10px; font-size: 13px; }
  .nav-icon { display: none; }
}
</style>
