<template>
  <div class="page-content">
    <h1 class="page-title">内容管理</h1>
    <el-card class="content-card" shadow="hover">
      <div class="filter-row">
        <span class="filter-label">状态筛选：</span>
        <el-radio-group v-model="filterStatus" size="default" @change="onFilterChange" class="filter-group">
          <el-radio-button value="all">全部</el-radio-button>
          <el-radio-button value="normal">正常</el-radio-button>
          <el-radio-button value="hidden">隐藏</el-radio-button>
          <el-radio-button value="flagged">标记</el-radio-button>
        </el-radio-group>
      </div>
      <el-table :data="list" v-loading="loading" stripe class="content-table">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip min-width="160" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small" class="status-tag">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发布时间" width="180" />
        <el-table-column label="操作" width="340" fixed="right">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="goDetail(row.id)">详情</el-button>
            <el-button v-if="row.status !== 'normal'" size="small" type="success" @click="setStatus(row.id, 'normal')">恢复</el-button>
            <el-button v-if="row.status === 'normal'" size="small" @click="setStatus(row.id, 'hidden')">隐藏</el-button>
            <el-button v-if="row.status === 'normal'" size="small" @click="setStatus(row.id, 'flagged')">标记</el-button>
            <el-button size="small" type="danger" @click="del(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="onPageChange"
        class="content-pagination"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { posts as postsApi } from '../api'
import { ElMessage } from 'element-plus'

const router = useRouter()
function goDetail(id) {
  router.push({ name: 'ContentDetail', params: { id } })
}

const list = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filterStatus = ref('all')

function statusLabel(s) {
  const map = { normal: '正常', hidden: '隐藏', flagged: '标记' }
  return map[s] || s
}
function statusTagType(s) {
  const map = { normal: 'success', hidden: 'warning', flagged: 'danger' }
  return map[s] || 'info'
}

function onFilterChange() {
  page.value = 1
  load()
}

function onPageChange(p) {
  page.value = p
  load()
}

async function load() {
  loading.value = true
  try {
    const params = {
      offset: (page.value - 1) * pageSize.value,
      limit: pageSize.value,
    }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    const { data } = await postsApi.list(params)
    list.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    list.value = []
  } finally {
    loading.value = false
  }
}

async function setStatus(id, status) {
  try {
    await postsApi.updateStatus(id, status)
    ElMessage.success('已更新')
    load()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '失败')
  }
}

async function del(id) {
  if (!confirm('确定删除？')) return
  try {
    await postsApi.delete(id)
    ElMessage.success('已删除')
    load()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '失败')
  }
}

onMounted(load)
</script>

<style scoped>
.page-content { max-width: 1200px; }
.page-title {
  margin: 0 0 24px;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--admin-text);
  letter-spacing: -0.02em;
}
.content-card {
  border-radius: var(--admin-card-radius);
  border: 1px solid var(--admin-border);
  box-shadow: var(--admin-card-shadow);
  overflow: hidden;
}
.content-card :deep(.el-card__body) {
  padding: 20px 24px 24px;
}
.filter-row {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}
.filter-label {
  font-size: 0.9375rem;
  font-weight: 500;
  color: var(--admin-text-secondary);
}
.filter-group :deep(.el-radio-button__inner) {
  border-color: var(--admin-border);
  color: var(--admin-text-secondary);
}
.content-table :deep(.el-table__header th) {
  background: var(--admin-primary-light);
  color: var(--admin-text-secondary);
  font-weight: 600;
  font-size: 0.8125rem;
}
.content-pagination {
  margin-top: 20px;
  justify-content: flex-start;
}
.content-pagination :deep(.el-pagination__total) {
  color: var(--admin-text-secondary);
}
</style>
