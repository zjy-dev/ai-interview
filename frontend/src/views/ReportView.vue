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

function scoreColor(score: number): string {
  if (score >= 80) return "var(--success)";
  if (score >= 60) return "var(--warning)";
  return "var(--danger)";
}

function scoreLabel(score: number): string {
  if (score >= 90) return "优秀";
  if (score >= 80) return "良好";
  if (score >= 60) return "合格";
  return "待提升";
}
</script>

<template>
  <div class="report-view">
    <!-- Header -->
    <div class="rpt-header">
      <div>
        <button class="btn btn-ghost btn-sm" @click="router.push('/interviews')">← 返回列表</button>
        <h1>面试评估报告</h1>
        <p class="rpt-subtitle">{{ store.current?.title }} · {{ store.current?.position }}</p>
      </div>
    </div>

    <!-- Loading skeleton -->
    <div v-if="!store.evaluation" class="skeleton-wrap">
      <div class="skel skel-score"></div>
      <div class="skel skel-bar"></div>
      <div class="skel skel-bar short"></div>
    </div>

    <template v-else>
      <!-- Score Hero -->
      <div class="score-hero card">
        <div class="score-ring" :style="{ '--score-color': scoreColor(store.evaluation.overall_score), '--score-pct': `${store.evaluation.overall_score * 3.6}deg` } as any">
          <div class="score-inner">
            <span class="score-num">{{ store.evaluation.overall_score }}</span>
            <span class="score-max">/ 100</span>
          </div>
        </div>
        <div class="score-info">
          <span class="score-badge" :style="{ background: scoreColor(store.evaluation.overall_score) }">
            {{ scoreLabel(store.evaluation.overall_score) }}
          </span>
          <p class="summary-text">{{ store.evaluation.summary }}</p>
        </div>
      </div>

      <!-- Category breakdown -->
      <div class="card section-card">
        <div class="section-label">分项评分</div>
        <div class="cat-list">
          <div
            v-for="(cat, i) in store.evaluation.categories"
            :key="cat.category"
            class="cat-row"
            :style="{ animationDelay: `${i * 80}ms` }"
          >
            <div class="cat-top">
              <span class="cat-name">{{ cat.category }}</span>
              <span class="cat-score" :style="{ color: scoreColor(cat.score) }">{{ cat.score }}</span>
            </div>
            <div class="progress-track">
              <div
                class="progress-fill"
                :style="{ width: `${cat.score}%`, background: scoreColor(cat.score) }"
              />
            </div>
            <p v-if="cat.comment" class="cat-comment">{{ cat.comment }}</p>
          </div>
        </div>
      </div>

      <!-- Strengths / Weaknesses -->
      <div class="feedback-grid">
        <div class="card feedback-card feedback-good">
          <div class="feedback-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
          </div>
          <h3>优势亮点</h3>
          <p>{{ store.evaluation.strengths }}</p>
        </div>
        <div class="card feedback-card feedback-warn">
          <div class="feedback-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          </div>
          <h3>待改进项</h3>
          <p>{{ store.evaluation.weaknesses }}</p>
        </div>
      </div>

      <!-- Suggestions -->
      <div class="card section-card suggestions-card">
        <div class="section-label">改进建议</div>
        <div class="suggestions-icon">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M12 2a7 7 0 0 0-4 12.7V17h8v-2.3A7 7 0 0 0 12 2z"/></svg>
        </div>
        <p class="suggestions-text">{{ store.evaluation.suggestions }}</p>
      </div>

      <!-- Actions -->
      <div class="rpt-actions">
        <button class="btn btn-secondary" @click="router.push('/interviews')">返回列表</button>
        <button class="btn btn-primary" @click="router.push('/interviews/new')">开始新面试 →</button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.report-view {
  max-width: 780px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
  animation: fadeInUp 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(16px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Header */
.rpt-header h1 {
  font-size: 26px;
  font-weight: 800;
  margin: 8px 0 4px;
}

.rpt-subtitle {
  font-size: 14px;
  color: var(--text-muted);
}

.btn-sm { padding: 8px 14px; font-size: 13px; }

/* Skeleton loading */
.skeleton-wrap {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.skel {
  border-radius: var(--radius-lg);
  background: var(--surface-hover);
  animation: shimmer 1.5s infinite;
}

.skel-score { height: 200px; }
.skel-bar { height: 60px; }
.skel-bar.short { width: 60%; height: 40px; }

@keyframes shimmer {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 0.8; }
}

/* Score Hero */
.score-hero {
  display: flex;
  align-items: center;
  gap: 32px;
  padding: 36px;
}

.score-ring {
  position: relative;
  width: 130px;
  height: 130px;
  border-radius: 50%;
  background: conic-gradient(var(--score-color) var(--score-pct), var(--surface-hover) var(--score-pct));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.score-inner {
  width: 104px;
  height: 104px;
  border-radius: 50%;
  background: var(--surface-solid);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.score-num {
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 40px;
  font-weight: 800;
  line-height: 1;
  color: var(--text);
}

.score-max {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 2px;
}

.score-info {
  flex: 1;
}

.score-badge {
  display: inline-block;
  padding: 4px 14px;
  border-radius: var(--radius-full);
  color: #fff;
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 12px;
}

.summary-text {
  font-size: 15px;
  line-height: 1.7;
  color: var(--text-secondary);
}

/* Section card */
.section-card { padding: 28px 32px; }

.section-label {
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--primary);
  margin-bottom: 20px;
}

/* Category list */
.cat-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.cat-row {
  animation: fadeInUp 0.4s cubic-bezier(0.16, 1, 0.3, 1) both;
}

.cat-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 6px;
}

.cat-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

.cat-score {
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 18px;
  font-weight: 700;
}

.progress-track {
  height: 6px;
  background: var(--surface-hover);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}

.cat-comment {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 6px;
  line-height: 1.5;
}

/* Feedback grid */
.feedback-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.feedback-card {
  padding: 28px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.feedback-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.feedback-good .feedback-icon {
  background: var(--success-soft);
  color: var(--success);
}

.feedback-warn .feedback-icon {
  background: var(--warning-soft);
  color: var(--warning);
}

.feedback-card h3 {
  font-family: 'Outfit', system-ui, sans-serif;
  font-size: 16px;
  font-weight: 700;
}

.feedback-card p {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-secondary);
  white-space: pre-wrap;
}

/* Suggestions */
.suggestions-card {
  text-align: center;
  align-items: center;
  display: flex;
  flex-direction: column;
}

.suggestions-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: var(--primary-soft);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}

.suggestions-text {
  font-size: 15px;
  line-height: 1.7;
  color: var(--text-secondary);
  max-width: 600px;
  white-space: pre-wrap;
}

/* Actions */
.rpt-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding-top: 8px;
}

@media (max-width: 640px) {
  .score-hero {
    flex-direction: column;
    text-align: center;
    padding: 28px 20px;
  }
  .feedback-grid { grid-template-columns: 1fr; }
  .section-card { padding: 20px; }
}
</style>
