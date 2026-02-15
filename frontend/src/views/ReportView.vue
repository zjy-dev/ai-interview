<script setup lang="ts">
import { onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useInterviewStore } from "@/stores/interview";

const route = useRoute();
const router = useRouter();
const store = useInterviewStore();

const interviewId = computed(() => Number(route.params.id));

onMounted(async () => {
  await store.fetchEvaluation(interviewId.value);
});

function scoreColor(score: number) {
  if (score >= 80) return "#15803d";
  if (score >= 60) return "#ca8a04";
  return "#dc2626";
}
</script>

<template>
  <div class="report-view">
    <div class="page-header">
      <h1>面试评估报告</h1>
      <button class="btn btn-secondary" @click="router.push('/interviews')">
        返回列表
      </button>
    </div>

    <div v-if="!store.evaluation" class="loading card">加载中...</div>

    <template v-else>
      <div class="score-card card">
        <div
          class="overall-score"
          :style="{ color: scoreColor(store.evaluation.overall_score) }"
        >
          {{ store.evaluation.overall_score }}
        </div>
        <p class="score-label">综合评分</p>
        <p class="summary">{{ store.evaluation.summary }}</p>
      </div>

      <div class="categories card">
        <h3>分项评分</h3>
        <div class="category-grid">
          <div
            v-for="cat in store.evaluation.categories"
            :key="cat.category"
            class="category-item"
          >
            <div class="cat-header">
              <span class="cat-name">{{ cat.category }}</span>
              <span
                class="cat-score"
                :style="{ color: scoreColor(cat.score) }"
                >{{ cat.score }}</span
              >
            </div>
            <div class="progress-bar">
              <div
                class="progress-fill"
                :style="{
                  width: `${cat.score}%`,
                  background: scoreColor(cat.score),
                }"
              />
            </div>
            <p class="cat-comment">{{ cat.comment }}</p>
          </div>
        </div>
      </div>

      <div class="feedback-grid">
        <div class="card">
          <h3>✅ 优势</h3>
          <p>{{ store.evaluation.strengths }}</p>
        </div>
        <div class="card">
          <h3>⚠️ 不足</h3>
          <p>{{ store.evaluation.weaknesses }}</p>
        </div>
      </div>

      <div class="card">
        <h3>💡 改进建议</h3>
        <p>{{ store.evaluation.suggestions }}</p>
      </div>
    </template>
  </div>
</template>

<style scoped>
.report-view {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  font-size: 24px;
}

.loading {
  text-align: center;
  padding: 48px;
  color: var(--text-secondary);
}

.score-card {
  text-align: center;
  padding: 40px 24px;
}

.overall-score {
  font-size: 72px;
  font-weight: 800;
  line-height: 1;
}

.score-label {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 8px 0 16px;
}

.summary {
  font-size: 15px;
  line-height: 1.6;
  max-width: 600px;
  margin: 0 auto;
}

.categories h3 {
  margin-bottom: 16px;
}

.category-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.cat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.cat-name {
  font-weight: 500;
  font-size: 14px;
}

.cat-score {
  font-weight: 700;
  font-size: 16px;
}

.progress-bar {
  height: 6px;
  background: #f1f5f9;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.6s ease;
}

.cat-comment {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 4px;
}

.feedback-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.feedback-grid h3 {
  margin-bottom: 8px;
}

.feedback-grid p,
.card p {
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre-wrap;
}

.card h3 {
  font-size: 16px;
  margin-bottom: 12px;
}
</style>
