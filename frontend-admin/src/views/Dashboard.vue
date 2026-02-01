<template>
  <div class="page-dashboard">
    <h1 class="page-title">数据统计</h1>
    <el-row :gutter="20" class="stat-row">
      <el-col :xs="24" :sm="12" :md="8">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <span class="stat-card-header">用户数</span>
          </template>
          <div class="stat-value">{{ stats.user_count ?? '-' }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="8">
        <el-card shadow="hover" class="stat-card">
          <template #header>
            <span class="stat-card-header">发帖数</span>
          </template>
          <div class="stat-value">{{ stats.post_count ?? '-' }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { stats as statsApi } from '../api'

const stats = ref({})

onMounted(async () => {
  try {
    const { data } = await statsApi.get()
    stats.value = data
  } catch (e) {
    stats.value = {}
  }
})
</script>

<style scoped>
.page-dashboard { max-width: 960px; }
.page-title {
  margin: 0 0 24px;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--admin-text);
  letter-spacing: -0.02em;
}
.stat-row { margin-bottom: 0; }
.stat-card {
  border-radius: var(--admin-card-radius);
  border: 1px solid var(--admin-border);
  box-shadow: var(--admin-card-shadow);
  transition: box-shadow var(--admin-transition);
}
.stat-card:hover {
  box-shadow: var(--admin-card-shadow-hover);
}
.stat-card :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid var(--admin-primary-light);
  font-weight: 600;
  font-size: 0.9375rem;
  color: var(--admin-text-secondary);
}
.stat-card :deep(.el-card__body) {
  padding: 20px;
}
.stat-card-header { color: var(--admin-text-secondary); }
.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--admin-primary);
  letter-spacing: -0.02em;
}
@media (prefers-reduced-motion: reduce) {
  .stat-card { transition: none; }
}
</style>
