<template>
  <el-card class="category-card" shadow="hover">
    <template #header>
      <div class="card-header">
        <span>文本检测</span>
      </div>
    </template>

    <el-form label-position="top">
      <!-- 文本输入 -->
      <el-form-item label="输入要检测的文本" :for="inputId">
        <el-input
          v-model="text"
          type="textarea"
          :id="inputId"
          :rows="8"
          class="dark-textarea"
          placeholder="请输入要检测的文本内容..."
        />
      </el-form-item>

      <!-- 选择检测类别 + 按钮 -->
      <el-form-item>
        
        <div class="category-label-row" style="margin-bottom: 8px;">
          <span class="category-label-main">选择检测类别</span>

          <div class="category-actions">
            <el-button
              class="btn-clear"
              size="small"
              @click="clearCategories"
            >
              清空类别
            </el-button>

            <el-button
              class="btn-select-all"
              size="small"
              @click="selectAll"
            >
              全选
            </el-button>
          </div>
        </div>

        <el-checkbox-group v-model="selectedKeys">
          <el-row :gutter="12">
            <el-col
              v-for="(label, key) in categories"
              :key="key"
              :xs="12"
              :sm="8"
            >
              <el-checkbox :label="key">
                {{ label }}
              </el-checkbox>
            </el-col>
          </el-row>
        </el-checkbox-group>
      </el-form-item>
      
      <!-- 底部按钮 -->
      <el-form-item>
        <el-space>
          <el-button @click="handleClear">清空文本</el-button>
          <el-button
            type="primary"
            :loading="props.loading"
            :disabled="!text.trim() || selectedKeys.length === 0"
            @click="handleSubmit"
          >
            {{ props.loading ? '检测中...' : '开始检测' }}
          </el-button>
        </el-space>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['detect'])

const text = ref('')
const categories = ref({})
const selectedKeys = ref([])


// 动态生成 inputId，确保唯一性
const inputId = `el-id-${Math.random().toString(36).substr(2, 9)}`

// 加载类别（初始状态：不选任何一个）
onMounted(async () => {
  try {
    const resp = await axios.get('/api/categories')
    categories.value = resp.data || {}
    selectedKeys.value = [] // 初始不选
    console.log('CategoryPanel 加载类别:', categories.value)
  } catch (e) {
    console.error('加载类别失败', e)
  }
})

// 点击“开始检测”
const handleSubmit = () => {
  if (!text.value.trim()) {
    console.warn('文本为空，不发请求')
    return
  }
  if (selectedKeys.value.length === 0) {
    console.warn('未选择任何检测类别，不发请求')
    return
  }

  console.log('CategoryPanel emit detect:', {
    textPreview: text.value.slice(0, 30) + '...',
    categories: selectedKeys.value,
  })

  emit('detect', {
    text: text.value,
    categories: selectedKeys.value,
  })
}

// 清空文本
const handleClear = () => {
  text.value = ''
}

// 全选所有类别
const selectAll = () => {
  selectedKeys.value = Object.keys(categories.value)
}

// 清空类别
const clearCategories = () => {
  selectedKeys.value = []
}
</script>

<style scoped>
.category-card {
  background: var(--bg-card-soft);
  border: 1px solid var(--border-subtle);
  color: var(--text-main);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  color: var(--text-sub);
}

/* 文本输入区域 */
.dark-textarea .el-textarea__inner {
  background: var(--input-bg);
  border: 1px solid var(--input-border);
  color: var(--text-main);
  font-size: 14px;
  border-radius: 12px;
  padding: 12px;
  resize: none;
  min-height: 150px;
}

.dark-textarea .el-textarea__inner::placeholder {
  color: rgba(148, 163, 184, 0.7);
}

.app-shell--dark .dark-textarea .el-textarea__inner::placeholder {
  color: rgba(148, 163, 184, 0.9);
}

/* 顶部 label + 按钮在同一行 */
.category-label-row {
  display: flex;
  align-items: center;
}

.category-label-main {
  font-weight: 500;
  color: var(--text-sub);
}

/* 右侧按钮组：贴右侧，小胶囊排布 */
.category-actions {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

/* 统一压一压 Element 默认按钮高度 */
.category-actions .el-button {
  padding: 4px 14px;
  font-size: 12px;
  border-radius: 999px;
  line-height: 1.2;
}

/* —— 全选：主色渐变小胶囊 —— */
.btn-select-all {
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  color: #fff;

  border: 1px solid rgba(255, 255, 255, 0.18);
  box-shadow: 0 2px 5px rgba(79, 70, 229, 0.28);

  cursor: pointer;
  transition: 0.18s ease;
}

.btn-select-all:hover {
  opacity: 0.97;
  transform: translateY(-1px);
  box-shadow: 0 4px 10px rgba(79, 70, 229, 0.38);
}

/* —— 清空类别：浅色线框按钮 —— */
.btn-clear {
  background: transparent;
  border: 1px solid var(--border-subtle);
  color: var(--text-sub);
  cursor: pointer;
  transition: 0.18s ease;
}

/* 浅色主题 hover */
.app-shell:not(.app-shell--dark) .btn-clear:hover {
  background: rgba(0, 0, 0, 0.04);
  color: var(--text-main);
}

/* 深色主题 hover */
.app-shell--dark .btn-clear:hover {
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-main);
}

/* label 颜色走主题变量 */
:deep(.el-form-item__label) {
  color: var(--text-sub);
}

</style>
