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

function statusClass(status: string) {
  return `status-${status}`;
}
</script>

<template>
  <div>
    <div class="page-header">
      <h1>面试记录</h1>
      <button class="btn btn-primary" @click="router.push('/interviews/new')">
        + 新建面试
      </button>
    </div>

    <div v-if="store.loading" class="loading">加载中...</div>

    <div v-else-if="store.interviews.length === 0" class="empty card">
      <p>还没有面试记录</p>
      <button class="btn btn-primary" @click="router.push('/interviews/new')">
        开始你的第一次面试
      </button>
    </div>

    <div v-else class="interview-list">
      <div
        v-for="item in store.interviews"
        :key="item.id"
        class="interview-item card"
        @click="router.push(`/interviews/${item.id}`)"
      >
        <div class="item-header">
          <h3>{{ item.title }}</h3>
          <span :class="['status-badge', statusClass(item.status)]">{{
            statusLabel(item.status)
          }}</span>
        </div>
        <p class="item-meta">
          {{ item.position }} ·
          {{ new Date(item.created_at).toLocaleDateString() }}
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 24px;
}

.loading,
.empty {
  text-align: center;
  padding: 48px;
  color: var(--text-secondary);
}

.empty .btn {
  margin-top: 16px;
}

.interview-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.interview-item {
  cursor: pointer;
  transition: box-shadow 0.2s;
}

.interview-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-header h3 {
  font-size: 16px;
}

.item-meta {
  color: var(--text-secondary);
  font-size: 13px;
  margin-top: 4px;
}

.status-badge {
  font-size: 12px;
  padding: 3px 10px;
  border-radius: 12px;
  font-weight: 500;
}

.status-active {
  background: #dbeafe;
  color: #1d4ed8;
}

.status-completed {
  background: #dcfce7;
  color: #15803d;
}

.status-cancelled {
  background: #f3f4f6;
  color: #6b7280;
}
</style>
