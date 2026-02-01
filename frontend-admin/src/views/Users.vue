<template>
  <div class="page-users">
    <h1 class="page-title">用户管理（C 端用户列表）</h1>
    <el-card class="users-card" shadow="hover">
      <el-table :data="list" v-loading="loading" stripe class="users-table">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column prop="created_at" label="注册时间" width="180" />
      </el-table>
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="onPageChange"
        class="users-pagination"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { users as usersApi } from '../api'

const list = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

function onPageChange(p) {
  page.value = p
  load()
}

async function load() {
  loading.value = true
  try {
    const { data } = await usersApi.list({ offset: (page.value - 1) * pageSize.value, limit: pageSize.value })
    list.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    list.value = []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.page-users { max-width: 1200px; }
.page-title {
  margin: 0 0 24px;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--admin-text);
  letter-spacing: -0.02em;
}
.users-card {
  border-radius: var(--admin-card-radius);
  border: 1px solid var(--admin-border);
  box-shadow: var(--admin-card-shadow);
  overflow: hidden;
}
.users-card :deep(.el-card__body) {
  padding: 20px 24px 24px;
}
.users-table :deep(.el-table__header th) {
  background: var(--admin-primary-light);
  color: var(--admin-text-secondary);
  font-weight: 600;
  font-size: 0.8125rem;
}
.users-pagination {
  margin-top: 20px;
  justify-content: flex-start;
}
.users-pagination :deep(.el-pagination__total) {
  color: var(--admin-text-secondary);
}
</style>
