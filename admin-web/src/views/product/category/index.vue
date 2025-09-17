<template>
  <div class="category-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>商品分类管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增分类
          </el-button>
        </div>
      </template>

      <div class="search-bar">
        <el-form :model="searchForm" inline>
          <el-form-item label="父级分类">
            <el-select v-model="searchForm.parent_id" placeholder="请选择父级分类" clearable style="width: 200px">
              <el-option label="顶级分类" :value="0" />
              <el-option
                v-for="category in flatCategories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 120px">
              <el-option label="启用" :value="1" />
              <el-option label="禁用" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
        v-loading="loading"
        :data="categoryTree"
        style="width: 100%"
        row-key="id"
        :tree-props="{ children: 'children' }"
        default-expand-all
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分类名称" min-width="200">
          <template #default="{ row }">
            <div class="category-name">
              <el-image
                v-if="row.icon"
                :src="row.icon"
                style="width: 24px; height: 24px; margin-right: 8px"
                fit="cover"
              />
              {{ row.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="层级" width="80">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)">{{ row.level }}级</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleAddChild(row)">添加子分类</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 分类表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="父级分类" prop="parent_id">
          <el-select v-model="form.parent_id" placeholder="请选择父级分类" style="width: 100%">
            <el-option label="顶级分类" :value="0" />
            <el-option
              v-for="category in flatCategories.filter(c => c.id !== form.id)"
              :key="category.id"
              :label="`${'　'.repeat(category.level - 1)}${category.name}`"
              :value="category.id"
              :disabled="category.level >= 3"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="分类图标" prop="icon">
          <div class="icon-upload">
            <el-upload
              :show-file-list="false"
              :on-success="handleIconSuccess"
              :before-upload="beforeIconUpload"
              action="/api/admin/upload/image"
              :headers="{ Authorization: `Bearer ${getToken()}` }"
            >
              <el-image
                v-if="form.icon"
                :src="form.icon"
                style="width: 80px; height: 80px; border: 1px dashed #d9d9d9; border-radius: 6px;"
                fit="cover"
              />
              <el-icon v-else class="icon-uploader"><Plus /></el-icon>
            </el-upload>
            <div class="icon-tips">建议尺寸：80x80px</div>
          </div>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入分类描述"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getCategoryTree,
  createCategory,
  updateCategory,
  deleteCategory,
  updateCategoryStatus,
  uploadImage,
  type Category
} from '@/api/product'
import { getToken } from '@/utils/auth'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()

const searchForm = reactive({
  parent_id: undefined as number | undefined,
  status: undefined as number | undefined
})

const form = reactive({
  id: undefined as number | undefined,
  name: '',
  parent_id: 0,
  icon: '',
  sort: 0,
  description: '',
  status: 1
})

const formRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 1, max: 50, message: '分类名称长度在 1 到 50 个字符', trigger: 'blur' }
  ],
  parent_id: [
    { required: true, message: '请选择父级分类', trigger: 'change' }
  ],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' }
  ]
}

const categoryTree = ref<Category[]>([])
const originalCategories = ref<Category[]>([])

// 扁平化分类列表，用于父级分类选择
const flatCategories = computed(() => {
  const flatten = (categories: Category[], level = 1): Category[] => {
    const result: Category[] = []
    categories.forEach(category => {
      result.push({ ...category, level })
      if (category.children && category.children.length > 0) {
        result.push(...flatten(category.children, level + 1))
      }
    })
    return result
  }
  return flatten(originalCategories.value)
})

const dialogTitle = computed(() => {
  return form.id ? '编辑分类' : '新增分类'
})

const getLevelType = (level: number) => {
  const types = ['', 'primary', 'success', 'warning']
  return types[level] || 'info'
}

const loadCategories = async () => {
  try {
    loading.value = true
    const { data } = await getCategoryTree()
    originalCategories.value = data
    filterCategories()
  } catch (error) {
    ElMessage.error('加载分类列表失败')
  } finally {
    loading.value = false
  }
}

const filterCategories = () => {
  let filtered = [...originalCategories.value]
  
  // 根据搜索条件过滤
  if (searchForm.parent_id !== undefined) {
    if (searchForm.parent_id === 0) {
      // 只显示顶级分类
      filtered = filtered.filter(cat => cat.parent_id === 0)
    } else {
      // 显示指定父级分类的子分类
      const parent = flatCategories.value.find(cat => cat.id === searchForm.parent_id)
      if (parent) {
        filtered = parent.children || []
      }
    }
  }
  
  if (searchForm.status !== undefined) {
    const filterByStatus = (categories: Category[]): Category[] => {
      return categories
        .filter(cat => cat.status === searchForm.status)
        .map(cat => ({
          ...cat,
          children: cat.children ? filterByStatus(cat.children) : []
        }))
    }
    filtered = filterByStatus(filtered)
  }
  
  categoryTree.value = filtered
}

const handleSearch = () => {
  filterCategories()
}

const handleReset = () => {
  searchForm.parent_id = undefined
  searchForm.status = undefined
  filterCategories()
}

const handleAdd = () => {
  resetForm()
  dialogVisible.value = true
}

const handleAddChild = (row: Category) => {
  resetForm()
  form.parent_id = row.id
  dialogVisible.value = true
}

const handleEdit = (row: Category) => {
  form.id = row.id
  form.name = row.name
  form.parent_id = row.parent_id
  form.icon = row.icon || ''
  form.sort = row.sort
  form.description = row.description || ''
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = async (row: Category) => {
  if (row.children && row.children.length > 0) {
    ElMessage.warning('该分类下存在子分类，无法删除')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确认删除分类"${row.name}"吗？`,
      '删除确认',
      { type: 'warning' }
    )
    
    await deleteCategory(row.id)
    ElMessage.success('删除成功')
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleStatusChange = async (row: Category) => {
  try {
    await updateCategoryStatus(row.id, row.status)
    ElMessage.success('状态更新成功')
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1 // 回滚状态
    ElMessage.error('状态更新失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    if (form.id) {
      await updateCategory(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await createCategory(form)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadCategories()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

const handleDialogClose = () => {
  formRef.value?.resetFields()
  resetForm()
}

const resetForm = () => {
  form.id = undefined
  form.name = ''
  form.parent_id = 0
  form.icon = ''
  form.sort = 0
  form.description = ''
  form.status = 1
}

const handleIconSuccess = (response: any) => {
  form.icon = response.data.url
}

const beforeIconUpload = (file: File) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('上传头像图片只能是 JPG/PNG 格式!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('上传头像图片大小不能超过 2MB!')
    return false
  }
  return true
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped lang="scss">
.category-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .search-bar {
    margin-bottom: 20px;
  }

  .category-name {
    display: flex;
    align-items: center;
  }

  .icon-upload {
    .icon-uploader {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 80px;
      height: 80px;
      border: 1px dashed #d9d9d9;
      border-radius: 6px;
      cursor: pointer;
      font-size: 28px;
      color: #8c939d;
      
      &:hover {
        border-color: #409EFF;
        color: #409EFF;
      }
    }

    .icon-tips {
      margin-top: 8px;
      font-size: 12px;
      color: #999;
      text-align: center;
    }
  }
}
</style>