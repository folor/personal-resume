<template>
  <!-- ===== 工具页：浅绿简约后台管理布局 ===== -->
  <div v-if="isToolPage" class="tool-app">
    <!-- 左侧固定侧边栏 -->
    <aside class="tool-sidebar">
      <div class="sidebar-brand" @click="goHome">
        <span class="brand-icon">YS</span>
        <div class="brand-text">
          <span class="brand-name">杨素</span>
          <span class="brand-sub">学习路线工具</span>
        </div>
      </div>

      <nav class="sidebar-nav">
        <div class="nav-group-label">学习路线</div>
        <router-link
          v-for="r in toolRoutes"
          :key="r.to"
          :to="r.to"
          class="sidebar-item"
        >
          <span class="item-icon">{{ r.icon }}</span>
          <span class="item-label">{{ r.label }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <a href="/" class="sidebar-item home-back" @click.prevent="goHome">
          <span class="item-icon">←</span>
          <span class="item-label">返回首页</span>
        </a>
      </div>
    </aside>

    <!-- 右侧内容区 -->
    <main class="tool-content">
      <router-view v-slot="{ Component }">
        <transition name="tool-page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- 移动端侧边栏开关 -->
    <button
      class="tool-mobile-toggle"
      :class="{ open: mobileSidebarOpen }"
      @click="mobileSidebarOpen = !mobileSidebarOpen"
      aria-label="菜单"
    >
      <span></span><span></span><span></span>
    </button>

    <!-- 移动端遮罩 -->
    <div
      v-if="mobileSidebarOpen"
      class="tool-mobile-mask"
      @click="mobileSidebarOpen = false"
    ></div>
    <aside
      v-if="mobileSidebarOpen"
      class="tool-sidebar mobile"
    >
      <div class="sidebar-brand" @click="goHome">
        <span class="brand-icon">YS</span>
        <div class="brand-text">
          <span class="brand-name">杨素</span>
          <span class="brand-sub">学习路线工具</span>
        </div>
      </div>
      <nav class="sidebar-nav">
        <div class="nav-group-label">学习路线</div>
        <router-link
          v-for="r in toolRoutes"
          :key="r.to"
          :to="r.to"
          class="sidebar-item"
          @click="mobileSidebarOpen = false"
        >
          <span class="item-icon">{{ r.icon }}</span>
          <span class="item-label">{{ r.label }}</span>
        </router-link>
      </nav>
      <div class="sidebar-footer">
        <a href="/" class="sidebar-item home-back" @click.prevent="goHome">
          <span class="item-icon">←</span>
          <span class="item-label">返回首页</span>
        </a>
      </div>
    </aside>
  </div>

  <!-- ===== 首页：深色玻璃风 ===== -->
  <template v-else>
    <!-- 动态背景光斑 -->
    <div class="bg-orbs">
      <div class="bg-orb"></div>
      <div class="bg-orb"></div>
      <div class="bg-orb"></div>
    </div>

    <!-- 导航栏 -->
    <NavBar />

    <!-- 路由出口 -->
    <router-view v-slot="{ Component }">
      <transition name="page" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>

    <!-- 页脚 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-inner">
          <div class="footer-brand">
            <span class="logo-icon">YS</span>
            <span>杨素 · 前端开发工程师</span>
          </div>
          <p class="footer-copy">
            &copy; {{ year }} 杨素 · yangsugogogo@qq.com · 184-8231-2017 · 成都
          </p>
        </div>
      </div>
    </footer>
  </template>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NavBar from './components/NavBar.vue'

const route = useRoute()
const router = useRouter()
const year = new Date().getFullYear()
const mobileSidebarOpen = ref(false)

const isToolPage = computed(() => !!route.meta.tool)

const toolRoutes = [
  { to: '/agent', label: 'Agent 开发路线', icon: '🤖' },
  { to: '/daily', label: '90 天每日打卡', icon: '✅' },
  { to: '/geo', label: 'GEO 优化路线', icon: '🔍' },
]

async function goHome() {
  mobileSidebarOpen.value = false
  await router.push('/')
}
</script>

<style scoped>
/* ===== 工具页布局 ===== */
.tool-app {
  display: flex;
  min-height: 100vh;
  background: #f3f8f5;
  --tool-green: #16a34a;
  --tool-green-dark: #15803d;
  --tool-green-light: #dcfce7;
  --tool-green-bg: #f0fdf4;
  --tool-text: #1a2e23;
  --tool-text-2: #4a5d52;
  --tool-text-3: #8a9b92;
  --tool-border: #d1e7da;
  --tool-card: #ffffff;
  --tool-cyan: #0891b2;
  --tool-amber: #d97706;
  --tool-pink: #db2777;
  --tool-purple: #7c3aed;
}

/* 侧边栏 */
.tool-sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 240px;
  background: #fff;
  border-right: 1px solid var(--tool-border);
  display: flex;
  flex-direction: column;
  z-index: 100;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 18px;
  cursor: pointer;
  border-bottom: 1px solid var(--tool-border);
}

.brand-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: var(--tool-green);
  color: #fff;
  font-size: 0.78rem;
  font-weight: 800;
  flex-shrink: 0;
}

.brand-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.brand-name {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--tool-text);
}

.brand-sub {
  font-size: 0.7rem;
  color: var(--tool-text-3);
}

.sidebar-nav {
  flex: 1;
  padding: 16px 10px;
  overflow-y: auto;
}

.nav-group-label {
  font-size: 0.68rem;
  font-weight: 600;
  color: var(--tool-text-3);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  padding: 0 10px;
  margin-bottom: 8px;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  text-decoration: none;
  color: var(--tool-text-2);
  font-size: 0.88rem;
  font-weight: 500;
  transition: all 0.2s ease;
  margin-bottom: 2px;
}

.sidebar-item:hover {
  background: var(--tool-green-bg);
  color: var(--tool-green-dark);
}

.sidebar-item.router-link-active {
  background: var(--tool-green-light);
  color: var(--tool-green-dark);
  font-weight: 600;
  border-left: 3px solid var(--tool-green);
  padding-left: 9px;
}

.item-icon {
  font-size: 1rem;
  flex-shrink: 0;
}

.sidebar-footer {
  padding: 10px;
  border-top: 1px solid var(--tool-border);
}

.home-back {
  color: var(--tool-text-3);
}

.home-back:hover {
  color: var(--tool-green-dark);
}

/* 内容区 */
.tool-content {
  margin-left: 240px;
  flex: 1;
  min-height: 100vh;
  padding: 32px 40px;
  max-width: 1100px;
}

/* 工具页过渡 */
.tool-page-enter-active,
.tool-page-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.tool-page-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.tool-page-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* 移动端侧边栏开关 */
.tool-mobile-toggle {
  display: none;
  position: fixed;
  top: 16px;
  left: 16px;
  z-index: 200;
  flex-direction: column;
  gap: 5px;
  background: #fff;
  border: 1px solid var(--tool-border);
  border-radius: 10px;
  padding: 10px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.tool-mobile-toggle span {
  display: block;
  width: 20px;
  height: 2px;
  border-radius: 2px;
  background: var(--tool-text);
  transition: 0.3s;
}

.tool-mobile-toggle.open span:nth-child(1) {
  transform: translateY(7px) rotate(45deg);
}
.tool-mobile-toggle.open span:nth-child(2) {
  opacity: 0;
}
.tool-mobile-toggle.open span:nth-child(3) {
  transform: translateY(-7px) rotate(-45deg);
}

.tool-sidebar.mobile {
  display: none;
}

.tool-mobile-mask {
  display: none;
}

/* ===== 首页布局（保留原样） ===== */
.footer {
  position: relative;
  z-index: 1;
  padding: 28px 0;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(0, 0, 0, 0.2);
}

.footer-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  text-align: center;
}

.footer-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 800;
  font-size: 1.05rem;
  color: var(--text-primary);
}

.logo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: var(--gradient-primary);
  color: #fff;
  font-size: 0.72rem;
  font-weight: 800;
}

.footer-copy {
  color: var(--text-muted);
  font-size: 0.82rem;
}

.page-enter-active,
.page-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .tool-sidebar {
    display: none;
  }

  .tool-mobile-toggle {
    display: flex;
  }

  .tool-sidebar.mobile {
    display: flex;
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 240px;
    z-index: 201;
    animation: slideIn 0.3s ease;
  }

  .tool-mobile-mask {
    display: block;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.3);
    z-index: 150;
  }

  .tool-content {
    margin-left: 0;
    padding: 60px 16px 32px;
  }

  .footer-copy {
    font-size: 0.72rem;
    line-height: 1.7;
  }
}

@keyframes slideIn {
  from { transform: translateX(-100%); }
  to { transform: translateX(0); }
}
</style>
