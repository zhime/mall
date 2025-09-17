<template>
  <div class="product-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>商品管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增商品
          </el-button>
        </div>
      </template>

      <div class="search-bar">
        <el-form :model="searchForm" inline>
          <el-form-item label="商品名称">
            <el-input
              v-model="searchForm.keyword"
              placeholder="请输入商品名称或SKU"
              style="width: 200px"
              clearable
            />
          </el-form-item>
          <el-form-item label="商品分类">
            <el-select v-model="searchForm.category_id" placeholder="请选择分类" clearable style="width: 150px">
              <el-option
                v-for="category in flatCategories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="商品状态">
            <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 120px">
              <el-option label="上架" :value="1" />
              <el-option label="下架" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <div class="toolbar">
        <el-button
          type="success"
          :disabled="!selectedProducts.length"
          @click="handleBatchStatus(1)"
        >
          批量上架
        </el-button>
        <el-button
          type="warning"
          :disabled="!selectedProducts.length"
          @click="handleBatchStatus(0)"
        >
          批量下架
        </el-button>
        <el-button
          type="danger"
          :disabled="!selectedProducts.length"
          @click="handleBatchDelete"
        >
          批量删除
        </el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="productList"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="商品信息" min-width="300">
          <template #default="{ row }">
            <div class="product-info">
              <el-image
                :src="row.images[0] || '/default-product.png'"
                style="width: 60px; height: 60px; margin-right: 12px"
                fit="cover"
              />
              <div class="product-details">
                <div class="product-name">{{ row.name }}</div>
                <div class="product-sku">SKU: {{ row.sku }}</div>
                <div class="product-category">{{ row.category_name }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="价格" width="120">
          <template #default="{ row }">
            <div>
              <div class="current-price">¥{{ row.price }}</div>
              <div v-if="row.original_price > row.price" class="original-price">
                ¥{{ row.original_price }}
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="80" />
        <el-table-column prop="sales" label="销量" width="80" />
        <el-table-column label="标签" width="120">
          <template #default="{ row }">
            <div class="product-tags">
              <el-tag v-if="row.is_featured" type="success" size="small">精选</el-tag>
              <el-tag v-if="row.is_new" type="warning" size="small">新品</el-tag>
            </div>
          </template>
        </el-table-column>
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
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import {
  getProductList,
  deleteProduct,
  updateProductStatus,
  batchUpdateProductStatus,
  batchDeleteProducts,
  getCategoryTree,
  type Product,
  type Category,
  type ProductQuery
} from '@/api/product'

const router = useRouter()
const loading = ref(false)

const searchForm = reactive({
  keyword: '',
  category_id: undefined as number | undefined,
  status: undefined as number | undefined,
  price_min: undefined as number | undefined,
  price_max: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const productList = ref<Product[]>([])
const selectedProducts = ref<Product[]>([])
const categories = ref<Category[]>([])

// 扁平化分类列表
const flatCategories = computed(() => {
  const flatten = (categories: Category[]): Category[] => {
    const result: Category[] = []
    categories.forEach(category => {
      result.push(category)
      if (category.children && category.children.length > 0) {
        result.push(...flatten(category.children))
      }
    })
    return result
  }
  return flatten(categories.value)
})

const loadProducts = async () => {
  try {
    loading.value = true
    const params: ProductQuery = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm
    }
    
    const { data } = await getProductList(params)
    productList.value = data.list
    pagination.total = data.total
  } catch (error) {
    ElMessage.error('加载商品列表失败')
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const { data } = await getCategoryTree()
    categories.value = data
  } catch (error) {
    ElMessage.error('加载分类列表失败')
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadProducts()
}

const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    category_id: undefined,
    status: undefined,
    price_min: undefined,
    price_max: undefined
  })
  pagination.page = 1
  loadProducts()
}

const handleAdd = () => {
  router.push('/product/add')
}

const handleView = (row: Product) => {
  router.push(`/product/view/${row.id}`)
}

const handleEdit = (row: Product) => {
  router.push(`/product/edit/${row.id}`)
}

const handleDelete = async (row: Product) => {
  try {
    await ElMessageBox.confirm(
      `确认删除商品"${row.name}"吗？`,
      '删除确认',
      { type: 'warning' }
    )
    
    await deleteProduct(row.id)
    ElMessage.success('删除成功')
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleStatusChange = async (row: Product) => {
  try {
    await updateProductStatus(row.id, row.status)
    ElMessage.success('状态更新成功')
  } catch (error) {
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('状态更新失败')
  }
}

const handleSelectionChange = (selection: Product[]) => {
  selectedProducts.value = selection
}

const handleBatchStatus = async (status: number) => {
  try {
    const ids = selectedProducts.value.map(p => p.id)
    await batchUpdateProductStatus(ids, status)
    ElMessage.success('批量操作成功')
    loadProducts()
  } catch (error) {
    ElMessage.error('批量操作失败')
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确认删除选中的 ${selectedProducts.value.length} 个商品吗？`,
      '批量删除确认',
      { type: 'warning' }
    )
    
    const ids = selectedProducts.value.map(p => p.id)
    await batchDeleteProducts(ids)
    ElMessage.success('批量删除成功')
    loadProducts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  loadProducts()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  loadProducts()
}

onMounted(() => {
  loadCategories()
  loadProducts()
})
</script>

<style scoped lang="scss">
.product-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .search-bar {
    margin-bottom: 20px;
  }

  .toolbar {
    margin-bottom: 20px;
  }

  .product-info {
    display: flex;
    align-items: center;

    .product-details {
      .product-name {
        font-weight: 500;
        margin-bottom: 4px;
      }

      .product-sku {
        font-size: 12px;
        color: #666;
        margin-bottom: 2px;
      }

      .product-category {
        font-size: 12px;
        color: #999;
      }
    }
  }

  .current-price {
    color: #f56c6c;
    font-weight: 500;
  }

  .original-price {
    color: #999;
    font-size: 12px;
    text-decoration: line-through;
  }

  .product-tags {
    .el-tag {
      margin-right: 4px;
    }
  }

  .pagination {
    margin-top: 20px;
    text-align: center;
  }
}
</style>