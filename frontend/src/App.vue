<script setup lang="ts">
import { RouterView } from "vue-router";
import NavBar from "./components/NavBar.vue";
import { useAuthStore } from "./stores/auth";

const auth = useAuthStore();
</script>

<template>
  <div class="app">
    <NavBar v-if="auth.isAuthenticated" />
    <main class="main-content" :class="{ 'has-nav': auth.isAuthenticated }">
      <RouterView />
    </main>
  </div>
</template>

<style>
/* ========== THEME VARIABLES ========== */
:root,
[data-theme="light"] {
  --primary: #0EA5E9;
  --primary-hover: #0284C7;
  --primary-soft: #E0F2FE;
  --primary-glow: rgba(14, 165, 233, 0.25);

  --bg: #F0F9FF;
  --bg-gradient: linear-gradient(135deg, #F0F9FF 0%, #E0F2FE 40%, #F0F9FF 70%, #BAE6FD 100%);
  --surface: rgba(255, 255, 255, 0.72);
  --surface-solid: #ffffff;
  --surface-hover: rgba(255, 255, 255, 0.9);

  --text: #0F172A;
  --text-secondary: #475569;
  --text-muted: #94A3B8;

  --border: rgba(14, 165, 233, 0.12);
  --border-strong: rgba(14, 165, 233, 0.25);
  --ring: rgba(14, 165, 233, 0.3);

  --success: #10B981;
  --success-soft: #D1FAE5;
  --warning: #F59E0B;
  --warning-soft: #FEF3C7;
  --danger: #EF4444;
  --danger-soft: #FEE2E2;

  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.03);
  --shadow: 0 4px 16px rgba(14, 165, 233, 0.08), 0 1px 3px rgba(0, 0, 0, 0.04);
  --shadow-lg: 0 12px 40px rgba(14, 165, 233, 0.12), 0 4px 12px rgba(0, 0, 0, 0.04);

  --radius: 12px;
  --radius-lg: 20px;
  --radius-full: 9999px;

  --backdrop: blur(16px) saturate(180%);

  --nav-bg: rgba(255, 255, 255, 0.8);
  --input-bg: rgba(255, 255, 255, 0.6);
  --msg-user-bg: linear-gradient(135deg, #0EA5E9, #0284C7);
  --msg-ai-bg: rgba(241, 245, 249, 0.8);
  --msg-ai-text: #0F172A;
}

[data-theme="dark"] {
  --primary: #38BDF8;
  --primary-hover: #7DD3FC;
  --primary-soft: rgba(56, 189, 248, 0.12);
  --primary-glow: rgba(56, 189, 248, 0.2);

  --bg: #0B1120;
  --bg-gradient: linear-gradient(135deg, #0B1120 0%, #0F1D32 40%, #0B1120 70%, #162033 100%);
  --surface: rgba(15, 23, 42, 0.72);
  --surface-solid: #0F172A;
  --surface-hover: rgba(30, 41, 59, 0.8);

  --text: #E2E8F0;
  --text-secondary: #94A3B8;
  --text-muted: #64748B;

  --border: rgba(56, 189, 248, 0.1);
  --border-strong: rgba(56, 189, 248, 0.2);
  --ring: rgba(56, 189, 248, 0.25);

  --success: #34D399;
  --success-soft: rgba(52, 211, 153, 0.12);
  --warning: #FBBF24;
  --warning-soft: rgba(251, 191, 36, 0.12);
  --danger: #F87171;
  --danger-soft: rgba(248, 113, 113, 0.12);

  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.2);
  --shadow: 0 4px 16px rgba(0, 0, 0, 0.3), 0 0 1px rgba(56, 189, 248, 0.05);
  --shadow-lg: 0 12px 40px rgba(0, 0, 0, 0.4), 0 0 1px rgba(56, 189, 248, 0.08);

  --nav-bg: rgba(11, 17, 32, 0.85);
  --input-bg: rgba(15, 23, 42, 0.6);
  --msg-user-bg: linear-gradient(135deg, #0369A1, #0EA5E9);
  --msg-ai-bg: rgba(30, 41, 59, 0.7);
  --msg-ai-text: #E2E8F0;

  --backdrop: blur(16px) saturate(150%);
}

/* ========== RESETS & BASE ========== */
*,
*::before,
*::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html {
  scroll-behavior: smooth;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

body {
  font-family: 'DM Sans', system-ui, sans-serif;
  background: var(--bg-gradient);
  background-attachment: fixed;
  color: var(--text);
  line-height: 1.6;
  transition: background 0.4s ease, color 0.3s ease;
  min-height: 100vh;
}

/* ========== LAYOUT ========== */
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  max-width: 1100px;
  margin: 0 auto;
  padding: 32px 24px;
  width: 100%;
  animation: fadeInUp 0.5s ease both;
}

.main-content.has-nav {
  padding-top: 32px;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

/* ========== TYPOGRAPHY ========== */
h1, h2, h3, h4, h5, h6 {
  font-family: 'Outfit', system-ui, sans-serif;
  font-weight: 600;
  letter-spacing: -0.02em;
  line-height: 1.3;
  color: var(--text);
}

a {
  color: var(--primary);
  text-decoration: none;
  transition: color 0.2s;
}
a:hover { color: var(--primary-hover); }

/* ========== GLASS CARD ========== */
.card {
  background: var(--surface);
  backdrop-filter: var(--backdrop);
  -webkit-backdrop-filter: var(--backdrop);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow);
  padding: 28px;
  transition: background 0.3s, border-color 0.3s, box-shadow 0.3s, transform 0.2s;
}

.card:hover {
  box-shadow: var(--shadow-lg);
}

/* ========== BUTTONS ========== */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 11px 22px;
  border-radius: var(--radius);
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.01em;
  cursor: pointer;
  border: none;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.btn::after {
  content: '';
  position: absolute;
  inset: 0;
  background: currentColor;
  opacity: 0;
  transition: opacity 0.25s;
}

.btn:active::after { opacity: 0.08; }

.btn-primary {
  background: var(--primary);
  color: #fff;
  box-shadow: 0 2px 8px var(--primary-glow);
}
.btn-primary:hover {
  background: var(--primary-hover);
  box-shadow: 0 4px 16px var(--primary-glow);
  transform: translateY(-1px);
}
.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.btn-secondary {
  background: var(--surface);
  backdrop-filter: var(--backdrop);
  color: var(--text);
  border: 1px solid var(--border-strong);
}
.btn-secondary:hover {
  background: var(--surface-hover);
  border-color: var(--primary);
}

.btn-danger {
  background: var(--danger);
  color: #fff;
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.25);
}
.btn-danger:hover {
  background: #DC2626;
  box-shadow: 0 4px 16px rgba(239, 68, 68, 0.3);
  transform: translateY(-1px);
}

.btn-ghost {
  background: transparent;
  color: var(--text-secondary);
  padding: 8px 14px;
}
.btn-ghost:hover {
  background: var(--primary-soft);
  color: var(--primary);
}

/* ========== FORM CONTROLS ========== */
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 13px;
  font-weight: 500;
  letter-spacing: 0.03em;
  text-transform: uppercase;
  margin-bottom: 8px;
  color: var(--text-secondary);
}

.form-control {
  width: 100%;
  padding: 12px 16px;
  background: var(--input-bg);
  backdrop-filter: var(--backdrop);
  border: 1.5px solid var(--border);
  border-radius: var(--radius);
  font-family: 'DM Sans', system-ui, sans-serif;
  font-size: 15px;
  color: var(--text);
  outline: none;
  transition: border-color 0.25s, box-shadow 0.25s, background 0.25s;
}

.form-control::placeholder {
  color: var(--text-muted);
}

.form-control:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--ring);
  background: var(--surface);
}

select.form-control {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' fill='%2394A3B8' viewBox='0 0 16 16'%3E%3Cpath d='M8 11L3 6h10l-5 5z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  padding-right: 36px;
}

textarea.form-control {
  resize: vertical;
  min-height: 48px;
}

/* ========== SCROLLBAR ========== */
::-webkit-scrollbar { width: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb {
  background: var(--border-strong);
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover { background: var(--primary); }

/* ========== RESPONSIVE ========== */
@media (max-width: 640px) {
  .main-content { padding: 20px 16px; }
  .card { padding: 20px; border-radius: var(--radius); }
}
</style>
