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
  <div class="new-interview">
    <h1>新建面试</h1>
    <p class="subtitle">配置面试参数，开始一次模拟面试</p>

    <form class="card" @submit.prevent="handleCreate">
      <div class="form-group">
        <label>面试标题 *</label>
        <input
          v-model="form.title"
          class="form-control"
          placeholder="例如：阿里巴巴前端一面"
          required
        />
      </div>

      <div class="form-group">
        <label>目标职位 *</label>
        <select v-model="form.position" class="form-control" required>
          <option value="" disabled>选择职位</option>
          <option v-for="p in positions" :key="p" :value="p">{{ p }}</option>
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

      <div class="form-group">
        <label>简历内容 (可选)</label>
        <textarea
          v-model="form.resume"
          class="form-control"
          rows="5"
          placeholder="粘贴你的简历内容，面试官会据此提问..."
        />
      </div>

      <p v-if="error" class="error">{{ error }}</p>

      <div class="form-actions">
        <button type="button" class="btn btn-secondary" @click="router.back()">
          取消
        </button>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? "创建中..." : "开始面试" }}
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.new-interview {
  max-width: 600px;
}

.new-interview h1 {
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

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
}
</style>
