<template>
  <div class="user-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="title">用户管理</span>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <div class="search-form">
        <el-form :model="queryForm" inline>
          <el-form-item label="关键词">
            <el-input
              v-model="queryForm.keyword"
              placeholder="用户名/手机号"
              clearable
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="queryForm.status" placeholder="请选择" clearable style="width: 120px">
              <el-option label="正常" :value="1" />
              <el-option label="禁用" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
      
      <!-- 用户表格 -->
      <el-table
        :data="userList"
        v-loading="loading"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              @click="handleView(row)"
            >
              查看
            </el-button>
            <el-button
              size="small"
              :type="row.status === 1 ? 'warning' : 'success'"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button
              size="small"
              type="info"
              @click="handleResetPassword(row)"
            >
              重置密码
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="queryForm.page"
          v-model:page-size="queryForm.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="用户详情"
      width="600px"
    >
      <div v-if="currentUser" class="user-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="用户ID">{{ currentUser.id }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ currentUser.username }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ currentUser.phone }}</el-descriptions-item>
          <el-descriptions-item label="昵称">{{ currentUser.nickname || '-' }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentUser.status === 1 ? 'success' : 'danger'">
              {{ currentUser.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">{{ currentUser.created_at }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
    
    <!-- 重置密码对话框 -->
    <el-dialog
      v-model="passwordVisible"
      title="重置密码"
      width="400px"
    >
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules">
        <el-form-item label="新密码" prop="password">
          <el-input
            v-model="passwordForm.password"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请确认新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="passwordVisible = false">取消</el-button>
          <el-button type="primary" @click="handleConfirmResetPassword">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getUserList, updateUserStatus, resetUserPassword } from '@/api/user'
import type { UserListQuery, UserListItem } from '@/types/user'

const loading = ref(false)
const userList = ref<UserListItem[]>([])
const total = ref(0)
const selectedUsers = ref<UserListItem[]>([])

// 查询表单
const queryForm = reactive<UserListQuery>({
  page: 1,
  page_size: 20,
  keyword: '',
  status: undefined
})

// 用户详情
const detailVisible = ref(false)
const currentUser = ref<UserListItem | null>(null)

// 重置密码
const passwordVisible = ref(false)
const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive({
  password: '',
  confirmPassword: ''
})

const passwordRules: FormRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为6-20个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取用户列表
const fetchUserList = async () => {
  try {
    loading.value = true
    const { data } = await getUserList(queryForm)
    userList.value = data.items || []
    total.value = data.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  queryForm.page = 1
  fetchUserList()
}

// 重置
const handleReset = () => {
  queryForm.keyword = ''
  queryForm.status = undefined
  queryForm.page = 1
  fetchUserList()
}

// 查看用户详情
const handleView = (user: UserListItem) => {
  currentUser.value = user
  detailVisible.value = true
}

// 切换用户状态
const handleToggleStatus = async (user: UserListItem) => {
  const action = user.status === 1 ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const newStatus = user.status === 1 ? 0 : 1
    await updateUserStatus(user.id, newStatus)
    user.status = newStatus
    ElMessage.success(`${action}成功`)
  } catch (error) {
    console.error('更新用户状态失败:', error)
  }
}

// 重置密码
const handleResetPassword = (user: UserListItem) => {
  currentUser.value = user
  passwordForm.password = ''
  passwordForm.confirmPassword = ''
  passwordVisible.value = true
}

// 确认重置密码
const handleConfirmResetPassword = async () => {
  if (!passwordFormRef.value || !currentUser.value) return
  
  try {
    await passwordFormRef.value.validate()
    await resetUserPassword(currentUser.value.id, passwordForm.password)
    ElMessage.success('密码重置成功')
    passwordVisible.value = false
  } catch (error) {
    console.error('重置密码失败:', error)
  }
}

// 选择变化
const handleSelectionChange = (selection: UserListItem[]) => {
  selectedUsers.value = selection
}

// 分页大小变化
const handleSizeChange = (val: number) => {
  queryForm.page_size = val
  queryForm.page = 1
  fetchUserList()
}

// 当前页变化
const handleCurrentChange = (val: number) => {
  queryForm.page = val
  fetchUserList()
}

onMounted(() => {
  fetchUserList()
})
</script>

<style scoped lang="scss">
.user-list {
  .search-form {
    margin-bottom: 20px;
    padding: 20px;
    background: #f8f9fa;
    border-radius: 4px;
  }
  
  .pagination {
    margin-top: 20px;
    text-align: right;
  }
  
  .user-detail {
    padding: 20px 0;
  }
}
</style>