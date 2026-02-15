<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { useInterviewStore } from "@/stores/interview";

const store = useInterviewStore();
const router = useRouter();

onMounted(() => {
  store.fetchList();
});

function statusLabel(status: string) {
  const map: Record<string, string> = {
    active: "进行中",
    completed: "已完成",
    cancelled: "已取消",
  };
  return map[status] || status;
}

function statusIcon(status: string) {
  const map: Record<string, string> = {
    active: "⚡",
    completed: "✓",
    cancelled: "—",
  };
  return map[status] || "○";
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString("zh-CN", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}
</script>

<template>
  <div class="list-page">
    <div class="page-header">
      <div>
        <h1>面试记录</h1>
        <p class="page-desc">管理你的所有模拟面试历史</p>
      </div>
      <button class="btn btn-primary" @click="router.push('/interviews/new')">
        <svg
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2.5"
        >
          <line x1="12" y1="5" x2="12" y2="19" />
          <line x1="5" y1="12" x2="19" y2="12" />
        </svg>
        新建面试
      </button>
    </div>

    <!-- Loading skeleton -->
    <div v-if="store.loading" class="skeleton-list">
      <div v-for="i in 3" :key="i" class="skeleton-card card">
        <div class="skel skel-title"></div>
        <div class="skel skel-text"></div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="store.interviews.length === 0" class="empty-state card">
      <div class="empty-icon">◇</div>
      <h3>还没有面试记录</h3>
      <p>开始你的第一次 AI 模拟面试吧</p>
      <button class="btn btn-primary" @click="router.push('/interviews/new')">
        <svg
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2.5"
        >
          <line x1="12" y1="5" x2="12" y2="19" />
          <line x1="5" y1="12" x2="19" y2="12" />
        </svg>
        开始面试
      </button>
    </div>

    <!-- Interview list -->
    <TransitionGroup v-else name="list" tag="div" class="interview-grid">
      <div
        v-for="(item, index) in store.interviews"
        :key="item.id"
        class="interview-card card"
        :style="{ animationDelay: `${index * 60}ms` }"
        @click="router.push(`/interviews/${item.id}`)"
      >
        <div class="card-top">
          <span :class="['status-chip', `status-${item.status}`]">
            <span class="status-dot"></span>
            {{ statusLabel(item.status) }}
          </span>
          <span class="card-date">{{ formatDate(item.created_at) }}</span>
        </div>
        <h3 class="card-title">{{ item.title }}</h3>
        <p class="card-position">{{ item.position }}</p>
        <div class="card-footer">
          <span class="card-arrow">→</span>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.list-page {
  animation: fadeInUp 0.5s ease both;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 800;
  letter-spacing: -0.03em;
}

.page-desc {
  color: var(--text-secondary);
  font-size: 14px;
  margin-top: 4px;
}

/* Skeleton loading */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.skeleton-card {
  padding: 28px;
}

.skel {
  border-radius: 8px;
  background: var(--border);
  animation: shimmer 1.5s infinite;
}

.skel-title {
  width: 60%;
  height: 20px;
  margin-bottom: 12px;
}
.skel-text {
  width: 40%;
  height: 14px;
}

@keyframes shimmer {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
  100% {
    opacity: 1;
  }
}

/* Empty state */
.empty-state {
  text-align: center;
  padding: 64px 32px;
}

.empty-icon {
  font-size: 48px;
  color: var(--primary);
  filter: drop-shadow(0 0 16px var(--primary-glow));
  margin-bottom: 16px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-8px);
  }
}

.empty-state h3 {
  font-size: 20px;
  margin-bottom: 6px;
}

.empty-state p {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 24px;
}

/* Interview grid */
.interview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}

.interview-card {
  cursor: pointer;
  padding: 24px;
  position: relative;
  animation: cardIn 0.5s cubic-bezier(0.16, 1, 0.3, 1) both;
  transition:
    transform 0.25s,
    box-shadow 0.25s,
    border-color 0.25s;
}

@keyframes cardIn {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.interview-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-lg);
  border-color: var(--primary);
}

.interview-card:hover .card-arrow {
  transform: translateX(4px);
  color: var(--primary);
}

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.status-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-family: "Outfit", system-ui, sans-serif;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  padding: 4px 12px;
  border-radius: var(--radius-full);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-active {
  color: var(--primary);
  background: var(--primary-soft);
}
.status-active .status-dot {
  background: var(--primary);
  box-shadow: 0 0 6px var(--primary-glow);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

.status-completed {
  color: var(--success);
  background: var(--success-soft);
}
.status-completed .status-dot {
  background: var(--success);
}

.status-cancelled {
  color: var(--text-muted);
  background: var(--border);
}
.status-cancelled .status-dot {
  background: var(--text-muted);
}

.card-date {
  font-size: 12px;
  color: var(--text-muted);
}

.card-title {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 6px;
  letter-spacing: -0.01em;
}

.card-position {
  color: var(--text-secondary);
  font-size: 14px;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.card-arrow {
  font-size: 18px;
  color: var(--text-muted);
  transition:
    transform 0.25s,
    color 0.25s;
}

/* List transition */
.list-enter-active {
  transition: all 0.4s ease;
}
.list-leave-active {
  transition: all 0.3s ease;
}
.list-enter-from {
  opacity: 0;
  transform: translateY(20px);
}
.list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

@media (max-width: 640px) {
  .interview-grid {
    grid-template-columns: 1fr;
  }
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
