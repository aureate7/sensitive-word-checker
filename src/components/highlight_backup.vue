
<template>
  <!-- 用 v-html 渲染高亮后的文本 -->
  <div class="highlight-box" v-html="highlightedHtml"></div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  text: {
    type: String,
    default: '',
  },
  // 后端的 detected_words / 每项里至少有 word 字段
  words: {
    type: Array,
    default: () => [],
  },
})

// 转义正则特殊字符
const escapeRegExp = (s) =>
  s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')

const highlightedHtml = computed(() => {
  if (!props.text) return ''

  // 提取唯一敏感词字符串
  const uniq = Array.from(
    new Set(
      props.words
        .map((w) => (typeof w === 'string' ? w : w.word))
        .filter(Boolean),
    ),
  )
  if (!uniq.length) return props.text

  // 长的词优先匹配，避免子串互相抢
  uniq.sort((a, b) => b.length - a.length)

  const pattern = uniq.map(escapeRegExp).join('|')
  const re = new RegExp(pattern, 'g')

  return props.text.replace(
    re,
    (m) => `<span class="highlight">${m}</span>`,
  )
})
</script>

<style scoped>
.highlight-box {
  white-space: pre-wrap;
  line-height: 1.6;
  font-size: 13px;
}

/* 高亮样式 */
.highlight {
  background: #ffe58f;
  padding: 0 3px;
  border-radius: 3px;
}
</style>