<template>
  <div class="highlight-box">
    <template
      v-for="segment in highlightedSegments"
      :key="segment.key || segment.text"
    >
      <span
        v-if="segment.highlight"
        :class="['highlight', { 'highlight--focused': segment.focused }]"
        :data-range-key="segment.rangeKey || ''"
      >
        {{ segment.text }}
      </span>
      <span v-else>{{ segment.text }}</span>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  text: {
    type: String,
    default: '',
  },
  // 优先使用区间高亮（start/end 为字符下标，end 为开区间）
  ranges: {
    type: Array,
    default: () => [],
  },
  activeRangeKey: {
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

const normalizeRanges = (ranges, textLen) => {
  const uniq = new Map()
  ;(Array.isArray(ranges) ? ranges : []).forEach((r) => {
    const start = Number(r?.start)
    const end = Number(r?.end)
    if (!Number.isFinite(start) || !Number.isFinite(end) || end <= start) return
    const safeStart = Math.max(0, Math.min(textLen, start))
    const safeEnd = Math.max(0, Math.min(textLen, end))
    if (safeEnd <= safeStart) return
    const key = String(r?.key || `${safeStart}:${safeEnd}`)
    uniq.set(`${safeStart}:${safeEnd}`, { start: safeStart, end: safeEnd, key })
  })

  return Array.from(uniq.values()).sort((a, b) => {
    if (a.start !== b.start) return a.start - b.start
    return a.end - b.end
  })
}

const highlightedSegments = computed(() => {
  const text = String(props.text || '')
  if (!text) return []
  const runes = Array.from(text)

  const rangeList = normalizeRanges(props.ranges, runes.length)
  if (rangeList.length) {
    const segments = []
    let cursor = 0

    rangeList.forEach((range) => {
      if (range.start > cursor) {
        segments.push({
          key: `plain-${cursor}-${range.start}`,
          text: runes.slice(cursor, range.start).join(''),
          highlight: false,
          rangeKey: '',
          focused: false,
        })
      }
      segments.push({
        key: `hit-${range.key}`,
        text: runes.slice(range.start, range.end).join(''),
        highlight: true,
        rangeKey: range.key,
        focused: range.key === props.activeRangeKey,
      })
      cursor = range.end
    })

    if (cursor < runes.length) {
      segments.push({
        key: `plain-${cursor}-end`,
        text: runes.slice(cursor).join(''),
        highlight: false,
        rangeKey: '',
        focused: false,
      })
    }
    return segments
  }

  // 提取唯一敏感词字符串
  const uniq = Array.from(
    new Set(
      props.words
        .map((w) => (typeof w === 'string' ? w : w.word))
        .filter(Boolean),
    ),
  )
  if (!uniq.length) return [{ key: 'plain-all', text, highlight: false, rangeKey: '', focused: false }]

  // 长的词优先匹配，避免子串互相抢
  uniq.sort((a, b) => {
    const lenDiff = Array.from(b).length - Array.from(a).length
    if (lenDiff !== 0) return lenDiff
    return String(a).localeCompare(String(b), 'zh-Hans-CN')
  })

  const pattern = uniq.map(escapeRegExp).join('|')
  if (!pattern) return [{ key: 'plain-all', text, highlight: false, rangeKey: '', focused: false }]

  const re = new RegExp(pattern, 'g')
  const segments = []

  let lastIndex = 0
  for (const match of text.matchAll(re)) {
    const start = match.index ?? 0
    const hit = match[0]
    if (!hit) continue

    if (start > lastIndex) {
      segments.push({
        key: `plain-${lastIndex}-${start}`,
        text: text.slice(lastIndex, start),
        highlight: false,
        rangeKey: '',
        focused: false,
      })
    }

    segments.push({
      key: `hit-fallback-${start}`,
      text: hit,
      highlight: true,
      rangeKey: '',
      focused: false,
    })
    lastIndex = start + hit.length
  }

  if (!segments.length) {
    return [{ key: 'plain-all', text, highlight: false, rangeKey: '', focused: false }]
  }
  if (lastIndex < text.length) {
    segments.push({
      key: `plain-${lastIndex}-tail`,
      text: text.slice(lastIndex),
      highlight: false,
      rangeKey: '',
      focused: false,
    })
  }
  return segments
})

</script>

<style>
.highlight-box {
  white-space: pre-wrap;
  line-height: 1.6;
  font-size: 13px;
}

/* 高亮样式 */
.highlight {
  background: #ffe58f !important; /* 强制应用背景色 */
  padding: 0 3px !important;
  border-radius: 3px !important;
  color: #d44f4f !important; /* 强制应用字体颜色 */
  font-weight: bold !important;
}

.highlight--focused {
  position: relative;
  background: linear-gradient(90deg, #fb923c 0%, #ef4444 100%) !important;
  color: #fff !important;
  border-radius: 4px !important;
  box-shadow:
    0 0 0 2px #2563eb,
    0 0 0 6px rgba(59, 130, 246, 0.22),
    0 8px 18px rgba(239, 68, 68, 0.34);
  animation: focusedPulse 1.15s ease-in-out infinite;
}

@keyframes focusedPulse {
  0%,
  100% {
    box-shadow:
      0 0 0 2px #2563eb,
      0 0 0 6px rgba(59, 130, 246, 0.22),
      0 8px 18px rgba(239, 68, 68, 0.34);
  }
  50% {
    box-shadow:
      0 0 0 2px #2563eb,
      0 0 0 9px rgba(59, 130, 246, 0.3),
      0 10px 22px rgba(239, 68, 68, 0.4);
  }
}

</style>
