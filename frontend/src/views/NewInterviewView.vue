<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useInterviewStore } from "@/stores/interview";

const store = useInterviewStore();
const router = useRouter();

const form = ref({
  title: "",
  position: "",
  language: "zh-CN",
  resume: "",
});
const loading = ref(false);
const error = ref("");

const positions = [
  "前端工程师",
  "后端工程师",
  "全栈工程师",
  "DevOps 工程师",
  "数据工程师",
  "AI/ML 工程师",
  "产品经理",
  "其他",
];

const languages = [
  { value: "zh-CN", label: "中文" },
  { value: "en", label: "English" },
];

async function handleCreate() {
  if (!form.value.title || !form.value.position) {
    error.value = "请填写必填字段";
    return;
  }
  loading.value = true;
  error.value = "";
  try {
    const interview = await store.create(form.value);
    router.push(`/interviews/${interview.id}`);
  } catch (e: any) {
    error.value = e.response?.data?.error || "创建失败";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="new-page">
    <div class="page-header">
      <div>
        <h1>新建面试</h1>
        <p class="page-desc">配置面试参数，开始一次模拟面试</p>
      </div>
    </div>

    <form class="form-card card" @submit.prevent="handleCreate">
      <div class="form-section">
        <div class="section-label">基本信息</div>

        <div class="form-group">
          <label>面试标题 <span class="required">*</span></label>
          <input
            v-model="form.title"
            class="form-control"
            placeholder="例如：阿里巴巴前端一面"
            required
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>目标职位 <span class="required">*</span></label>
            <select v-model="form.position" class="form-control" required>
              <option value="" disabled>选择职位</option>
              <option v-for="p in positions" :key="p" :value="p">
                {{ p }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>面试语言</label>
            <select v-model="form.language" class="form-control">
              <option v-for="l in languages" :key="l.value" :value="l.value">
                {{ l.label }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <div class="form-section">
        <div class="section-label">
          简历内容 <span class="optional">(可选)</span>
        </div>
        <div class="form-group" style="margin-bottom: 0">
          <textarea
            v-model="form.resume"
            class="form-control"
            rows="6"
            placeholder="粘贴你的简历内容，面试官会据此提问..."
          />
        </div>
      </div>

      <Transition name="fade">
        <p v-if="error" class="error-msg">{{ error }}</p>
      </Transition>

      <div class="form-actions">
        <button type="button" class="btn btn-ghost" @click="router.back()">
          ← 返回
        </button>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? "创建中..." : "开始面试 →" }}
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.new-page {
  max-width: 640px;
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
  margin-bottom: 28px;
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

.form-card {
  padding: 32px;
}

.form-section {
  margin-bottom: 28px;
  padding-bottom: 28px;
  border-bottom: 1px solid var(--border);
}

.form-section:last-of-type {
  border-bottom: none;
  padding-bottom: 0;
}

.section-label {
  font-family: "Outfit", system-ui, sans-serif;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--primary);
  margin-bottom: 20px;
}

.optional {
  text-transform: none;
  letter-spacing: 0;
  color: var(--text-muted);
  font-weight: 400;
}

.required {
  color: var(--danger);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
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

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: space-between;
  margin-top: 24px;
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
  to {
    transform: rotate(360deg);
  }
}

.fade-enter-active {
  transition: all 0.3s;
}
.fade-leave-active {
  transition: all 0.2s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

@media (max-width: 640px) {
  .form-card {
    padding: 24px;
  }
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
