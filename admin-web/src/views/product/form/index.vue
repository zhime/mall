<template>
  <div class="product-form-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑商品' : '新增商品' }}</span>
          <el-button @click="handleBack">返回列表</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="120px"
        size="default"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="商品名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入商品名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="商品SKU" prop="sku">
              <el-input v-model="form.sku" placeholder="请输入商品SKU" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="商品分类" prop="category_id">
              <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 100%">
                <el-option
                  v-for="category in flatCategories"
                  :key="category.id"
                  :label="`${'　'.repeat(category.level - 1)}${category.name}`"
                  :value="category.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="商品重量" prop="weight">
              <el-input-number
                v-model="form.weight"
                :min="0"
                :precision="2"
                style="width: 100%"
                placeholder="商品重量(g)"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="销售价格" prop="price">
              <el-input-number
                v-model="form.price"
                :min="0"
                :precision="2"
                style="width: 100%"
                placeholder="销售价格"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="原价" prop="original_price">
              <el-input-number
                v-model="form.original_price"
                :min="0"
                :precision="2"
                style="width: 100%"
                placeholder="原价"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="库存数量" prop="stock">
              <el-input-number
                v-model="form.stock"
                :min="0"
                style="width: 100%"
                placeholder="库存数量"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="商品图片" prop="images">
          <el-upload
            v-model:file-list="imageFileList"
            action="/api/admin/upload/image"
            :headers="{ Authorization: `Bearer ${getToken()}` }"
            list-type="picture-card"
            :on-success="handleImageSuccess"
            :on-remove="handleImageRemove"
            :before-upload="beforeImageUpload"
            multiple
            :limit="5"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <div class="upload-tips">最多上传5张图片，建议尺寸：800x800px</div>
        </el-form-item>

        <el-form-item label="商品描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请输入商品描述"
          />
        </el-form-item>

        <el-form-item label="商品状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">上架</el-radio>
            <el-radio :label="0">下架</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="商品标签">
          <el-checkbox v-model="form.is_featured">精选商品</el-checkbox>
          <el-checkbox v-model="form.is_new" style="margin-left: 20px">新品</el-checkbox>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ isEdit ? '更新商品' : '创建商品' }}
          </el-button>
          <el-button @click="handleReset">重置</el-button>
          <el-button @click="handleBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, FormInstance, UploadFile } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getProduct,
  createProduct,
  updateProduct,
  getCategoryTree,
  type Product,
  type Category
} from '@/api/product'
import { getToken } from '@/utils/auth'

const router = useRouter()
const route = useRoute()
const formRef = ref<FormInstance>()
const submitting = ref(false)

const productId = computed(() => route.params.id as string)
const isEdit = computed(() => !!productId.value && productId.value !== 'add')

const form = reactive({
  name: '',
  sku: '',
  category_id: undefined as number | undefined,
  price: 0,
  original_price: 0,
  stock: 0,
  weight: 0,
  images: [] as string[],
  description: '',
  status: 1,
  is_featured: false,
  is_new: false
})

const formRules = {
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 1, max: 100, message: '商品名称长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  sku: [
    { required: true, message: '请输入商品SKU', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入销售价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '价格不能小于0', trigger: 'blur' }
  ],
  original_price: [
    { required: true, message: '请输入原价', trigger: 'blur' },
    { type: 'number', min: 0, message: '原价不能小于0', trigger: 'blur' }
  ],
  stock: [
    { required: true, message: '请输入库存数量', trigger: 'blur' },
    { type: 'number', min: 0, message: '库存不能小于0', trigger: 'blur' }
  ],
  weight: [
    { required: true, message: '请输入商品重量', trigger: 'blur' },
    { type: 'number', min: 0, message: '重量不能小于0', trigger: 'blur' }
  ],
  images: [
    { required: true, message: '请上传商品图片', trigger: 'change' }
  ]
}

const categories = ref<Category[]>([])
const imageFileList = ref<UploadFile[]>([])

// 扁平化分类列表
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
  return flatten(categories.value)
})

const loadCategories = async () => {
  try {
    const { data } = await getCategoryTree()
    categories.value = data
  } catch (error) {
    ElMessage.error('加载分类列表失败')
  }
}

const loadProduct = async () => {
  if (!isEdit.value) return
  
  try {
    const { data } = await getProduct(Number(productId.value))
    Object.assign(form, data)
    
    // 设置图片文件列表
    imageFileList.value = data.images.map((url: string, index: number) => ({
      name: `image-${index}`,
      url
    }))
  } catch (error) {
    ElMessage.error('加载商品信息失败')
    handleBack()
  }
}

const handleImageSuccess = (response: any, file: UploadFile) => {
  form.images.push(response.data.url)
}

const handleImageRemove = (file: UploadFile) => {
  const index = form.images.findIndex(url => url === file.url || url === file.response?.data?.url)
  if (index > -1) {
    form.images.splice(index, 1)
  }
}

const beforeImageUpload = (file: File) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('上传图片只能是 JPG/PNG 格式!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('上传图片大小不能超过 2MB!')
    return false
  }
  return true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    if (isEdit.value) {
      await updateProduct(Number(productId.value), form)
      ElMessage.success('更新成功')
    } else {
      await createProduct(form)
      ElMessage.success('创建成功')
    }
    
    handleBack()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  form.images = []
  imageFileList.value = []
}

const handleBack = () => {
  router.push('/product/list')
}

onMounted(() => {
  loadCategories()
  loadProduct()
})
</script>

<style scoped lang="scss">
.product-form-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .upload-tips {
    margin-top: 8px;
    font-size: 12px;
    color: #999;
  }
}
</style>