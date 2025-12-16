<template>
  <!-- 这里通过 app-shell / app-shell--dark 挂载全局 CSS 变量 -->
  <div :class="['app-shell', { 'app-shell--dark': isDark }]">
    <el-container class="layout">
      <!-- 顶部导航栏（ChatGPT 风格） -->
      <el-header class="nav">
        <!-- 左侧：品牌区 -->
        <div class="nav-left">
          <div class="brand">
            <el-icon class="brand-icon"><Lock /></el-icon>
            <div class="brand-text">
              <div class="brand-title">敏感词检测系统</div>
              <div class="brand-subtitle">Golang · Gin · Vue3</div>
            </div>
          </div>
        </div>

        <!-- 右侧：菜单 -->
        <div class="nav-right">
          <button class="nav-btn" @click="goHome">首页</button>

          <button class="nav-btn ghost" @click="toggleDark">
                <el-icon>
                    <!-- <Moon v-if="isDark" />
                    <Sunny v-else /> -->
                    <component :is="isDark ? Moon : Sunny"></component>
                </el-icon>
            <span>{{ isDark ? '深色模式' : '亮色模式' }}</span>
          </button>

          <a
            class="nav-btn primary"
            href="https://github.com"
            target="_blank"
          >
            GitHub
          </a>
        </div>
      </el-header>

      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import router from './router'
import {Moon, Sunny, Lock} from '@element-plus/icons-vue'

const isDark = ref(false)
const toggleDark = () => (isDark.value = !isDark.value)
const goHome = () => router.push('/')
</script>

<!-- /* 不要 scoped */ -->
<style>
/* 🌞 默认：浅色主题（挂在 .app-shell 上） */
.app-shell {
  /* 页面背景 & 文字 */
  --bg-page: #f3f4f6;
  --text-main: #111827;
  --text-sub: #4b5563;

  /* 导航栏 */
  --nav-bg: rgba(255, 255, 255, 0.96);
  --nav-border: rgba(15, 23, 42, 0.06);
  --nav-shadow: 0 8px 24px rgba(15, 23, 42, 0.08);

  --brand-title: #0f172a;
  --brand-sub: #6b7280;

  --nav-btn-bg: rgba(15, 23, 42, 0.04);
  --nav-btn-border: rgba(15, 23, 42, 0.12);
  --nav-btn-text: #111827;

  --nav-btn-primary-bg: rgba(22, 163, 74, 0.12);
  --nav-btn-primary-border: rgba(22, 163, 74, 0.35);
  --nav-btn-primary-text: #15803d;

  /* 主体卡片 */
  --bg-main: #eef2ff;
  --bg-shell: #ffffff;
  --bg-card-soft: #ffffff;
  --bg-card-alt: #f9fafb;
  --border-subtle: rgba(148, 163, 184, 0.25);
  --shadow-subtle: 0 24px 50px rgba(15, 23, 42, 0.08);

  /* Home 大面板背景 */
  --panel-shell-bg:
    radial-gradient(circle at top left, rgba(129, 140, 248, 0.18), transparent 55%),
    radial-gradient(circle at bottom right, rgba(56, 189, 248, 0.2), transparent 55%),
    #ffffff;

  /* 输入 / 预览背景 */
  --input-bg: #f9fafb;
  --input-border: rgba(148, 163, 184, 0.5);
  --preview-bg: #f9fafb;

  /* 外壳颜色（你刚才做的外黑内白那一层） */
  --result-shell-bg: #f3f4f6;

  transition: background-color 0.3s ease, color 0.3s ease, box-shadow 0.3s ease;
}

/* 🌙 深色主题：在 .app-shell--dark 上覆盖 */
.app-shell.app-shell--dark {
  --bg-page: #020617;
  --text-main: #e2e8f0;
  --text-sub: #94a3b8;

  --nav-bg: rgba(17, 25, 40, 0.78);
  --nav-border: rgba(148, 163, 184, 0.28);
  --nav-shadow:
    0 4px 20px rgba(0, 0, 0, 0.5),
    0 0 32px rgba(99, 102, 241, 0.2);

  --brand-title: #f1f5f9;
  --brand-sub: rgba(148, 163, 184, 0.85);

  --nav-btn-bg: rgba(255, 255, 255, 0.06);
  --nav-btn-border: rgba(255, 255, 255, 0.16);
  --nav-btn-text: #e5e7eb;

  --nav-btn-primary-bg: rgba(16, 185, 129, 0.22);
  --nav-btn-primary-border: rgba(16, 185, 129, 0.5);
  --nav-btn-primary-text: #6ee7b7;

  --bg-main: #020617;
  --bg-shell: rgba(15, 23, 42, 0.98);
  --bg-card-soft: rgba(15, 23, 42, 0.96);
  --bg-card-alt: rgba(15, 23, 42, 0.9);
  --border-subtle: rgba(148, 163, 184, 0.35);
  --shadow-subtle: 0 24px 60px rgba(15, 23, 42, 0.9);

  --panel-shell-bg:
    radial-gradient(circle at top left, rgba(148, 163, 184, 0.22), transparent 55%),
    radial-gradient(circle at bottom right, rgba(56, 189, 248, 0.22), transparent 55%),
    rgba(15, 23, 42, 0.96);

  --input-bg: rgba(15, 23, 42, 0.96);
  --input-border: rgba(148, 163, 184, 0.6);
  --preview-bg: rgba(15, 23, 42, 0.96);

  --result-shell-bg: #020617;
}

/* 让 body 也吃到变量（防止有旧样式残留） */
body {
  margin: 0;
  background: var(--bg-page);
  color: var(--text-main);
}
</style>


<!-- ========= 导航 / 布局具体样式（scoped） ========= -->
<style scoped>
.app-shell {
  min-height: 100vh;
  background: var(--bg-page);
  color: var(--text-main);
}

/* 让主内容居中一点，看起来更像控制台 */
.layout {
  max-width: 1440px;
  margin: 0 auto;
}

/* ===== 导航栏（使用主题变量） ===== */
.nav {
  height: 70px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 32px;
  margin-bottom: 12px;

  background: var(--nav-bg);
  border-bottom: 1px solid var(--nav-border);
  box-shadow: var(--nav-shadow);
}

/* 左侧品牌 */
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-icon {
  font-size: 28px;
  color: #93c5fd;
  filter: drop-shadow(0 0 8px rgba(96,165,250,0.4));
}

.brand-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--brand-title);
}

.brand-subtitle {
  font-size: 12px;
  color: var(--brand-sub);
}

/* ===== 导航按钮 ===== */
.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-btn {
  padding: 8px 18px;
  border-radius: 999px;

  background: var(--nav-btn-bg);
  border: 1px solid var(--nav-btn-border);

  color: var(--nav-btn-text);
  font-weight: 500;
  letter-spacing: 0.5px;

  display: inline-flex;
  align-items: center;
  gap: 6px;

  cursor: pointer;
  transition: 0.22s ease;
  text-decoration: none;
}

/* Hover 发光效果：这里用透明度，不再写死白色，深浅模式都好看 */
.nav-btn:hover {
  background: rgba(255, 255, 255, 0.10);
  border-color: rgba(255, 255, 255, 0.22);
  transform: translateY(-2px);
  box-shadow: 0 0 12px rgba(148,163,255,0.32);
}

/* GitHub 主按钮 */
.nav-btn.primary {
  background: var(--nav-btn-primary-bg);
  border-color: var(--nav-btn-primary-border);
  color: var(--nav-btn-primary-text);
}

.nav-btn.primary:hover {
  background: rgba(16, 185, 129, 0.3);
  box-shadow: 0 0 16px rgba(52,211,153,0.55);
}

/* 主内容区 */
.main {
  padding: 12px 24px 24px;
}
</style>
