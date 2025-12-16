<template>
  <el-card class="stats-card" shadow="never">
    <div class="card-header">
      <span class="title">敏感词库统计</span>
    </div>

    <div class="stats-list">
      <div class="stats-row" v-for="item in items" :key="item.label">
        <span class="label">{{ item.label }}</span>
        <span class="value">{{ format(item.value) }}</span>
      </div>
    </div>
  </el-card>
</template>

<script setup>


const items = [
  { label: "总计", value: 833926 },
  { label: "涉黄类敏感词", value: 301770 },
  { label: "广告敏感词", value: 234175 },
  { label: "辱骂类敏感词", value: 224009 },
  { label: "暴恐类敏感词", value: 60025 },
  { label: "政治敏感词", value: 13947 },
]

const format = (n) => n.toLocaleString()
</script>

<style scoped>
/* ================================
   📌 卡片基础结构（深浅主题通用）
================================ */
.stats-card {
  padding: 20px;
  border-radius: 14px;
  position: relative;
  /* z-index: 1; */

  background: var(--stats-card-bg);       /* 主题变量 */
  border: 1px solid var(--stats-card-border);
  box-shadow: var(--stats-card-shadow);

  color: var(--text-main);               /* 字体走主题变量 */
  
  overflow: hidden;
}

/* ================================
   标题
================================ */
.card-header {
  margin-bottom: 12px;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-main);
  /* color: black; */
}

/* ================================
   行容器
================================ */
.stats-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: relative;
  z-index: 2;
}

.stats-row {
  display: flex;
  justify-content: space-between;
  padding: 10px 4px;
  border-bottom: 1px solid var(--border-subtle);
}

.label {
  font-size: 14px;
  color: var(--text-sub);
}

.value {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-main);
}

/* ================================
   🌞 浅色主题：启用玻璃效果（高亮）
================================ */
.app-shell:not(.app-shell--dark) .stats-card {
  backdrop-filter: blur(22px);
  -webkit-backdrop-filter: blur(22px);

  background: rgba(255, 255, 255, 0.25);    /* 半透明水晶白 */
  border: 1px solid rgba(148, 163, 184, 0.35);

  box-shadow:
    0 25px 45px rgba(0, 0, 0, 0.09),
    inset 0 0 25px rgba(255,255,255,0.25);

  /* 玻璃高光层（非常关键） */
}

.app-shell:not(.app-shell--dark) .stats-card::before {
  content: "";
  position: absolute;
  inset: 0;
  background: radial-gradient(
      circle at top left,
      rgba(255,255,255,0.1),
      rgba(255,255,255,0.05)
  );
  pointer-events: none;
}

/* 外层细白高光框 */
.app-shell:not(.app-shell--dark) .stats-card::after {
  content: "";
  position: absolute;
  inset: 0;
  border-radius: inherit;
  border: 1px solid rgba(255,255,255,0.55);
  pointer-events: none;
}

/* ================================
   🌙 深色主题：无玻璃，稳重实卡片
================================ */
.app-shell--dark .stats-card {
  background: #1e293b;
  border: 1px solid rgba(255,255,255,0.06);
  box-shadow: 0 8px 22px rgba(0,0,0,0.35);
}

.app-shell--dark .stats-card::before,
.app-shell--dark .stats-card::after {
  display: none;   /* 移除玻璃光晕层 */
}
</style>
