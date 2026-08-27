<template>
  <nav class="navbar" :class="{ scrolled: isScrolled }">
    <div class="nav-container">
      <a href="/" class="nav-logo" @click.prevent="goHome">
        <span class="logo-icon">YS</span>
        <span class="logo-text">杨素</span>
      </a>

      <ul class="nav-links" :class="{ open: menuOpen }">
        <!-- 首页内锚点（仅在首页路由显示） -->
        <template v-if="isHome">
          <li v-for="link in anchorLinks" :key="link.href">
            <a :href="link.href" @click.prevent="scrollTo(link.href)">{{ link.label }}</a>
          </li>
        </template>

        <!-- 路由菜单：工具页 -->
        <li v-for="r in routeLinks" :key="r.to">
          <router-link :to="r.to" class="route-link" @click="menuOpen = false">
            <span class="route-icon">{{ r.icon }}</span>
            {{ r.label }}
          </router-link>
        </li>

        <!-- 工具页显示"返回首页" -->
        <li v-if="!isHome">
          <a href="/" class="route-link home-back" @click.prevent="goHome">← 返回首页</a>
        </li>
      </ul>

      <button class="nav-toggle" @click="menuOpen = !menuOpen" aria-label="菜单">
        <span></span>
        <span></span>
        <span></span>
      </button>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isScrolled = ref(false)
const menuOpen = ref(false)

const isHome = computed(() => route.name === 'home')

const anchorLinks = [
  { label: '技术栈', href: '#skills' },
  { label: '经历', href: '#experience' },
  { label: '项目', href: '#projects' },
  { label: '3D Demo', href: '#demo' },
  { label: '联系', href: '#contact' },
]

const routeLinks = [
  { to: '/agent', label: 'Agent 路线', icon: '🤖' },
  { to: '/daily', label: '每日打卡', icon: '✅' },
  { to: '/geo', label: 'GEO 路线', icon: '🔍' },
]

function scrollTo(href: string) {
  menuOpen.value = false
  const el = document.querySelector(href)
  if (el) el.scrollIntoView({ behavior: 'smooth' })
}

async function goHome() {
  menuOpen.value = false
  if (!isHome.value) {
    await router.push('/')
    await nextTick()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } else {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function onScroll() {
  isScrolled.value = window.scrollY > 40
}

onMounted(() => window.addEventListener('scroll', onScroll))
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: 16px 0;
  transition: all 0.3s ease;
}

.navbar.scrolled {
  padding: 10px 0;
  background: rgba(10, 14, 26, 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  font-weight: 800;
  font-size: 1.25rem;
  cursor: pointer;
}

.logo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: var(--gradient-primary);
  color: #fff;
  font-size: 0.85rem;
  font-weight: 800;
  box-shadow: 0 4px 16px rgba(77, 166, 255, 0.4);
}

.logo-text {
  background: var(--gradient-text);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 4px;
  list-style: none;
}

.nav-links a {
  display: block;
  padding: 8px 14px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.92rem;
  font-weight: 500;
  border-radius: 8px;
  transition: var(--transition);
  position: relative;
}

.nav-links a::after {
  content: '';
  position: absolute;
  bottom: 2px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 2px;
  border-radius: 2px;
  background: var(--gradient-primary);
  transition: width 0.3s ease;
}

.nav-links a:hover {
  color: var(--text-primary);
}

.nav-links a:hover::after,
.nav-links a.router-link-active::after {
  width: 50%;
}

/* 路由链接（工具页入口）特殊样式 */
.route-link {
  display: inline-flex !important;
  align-items: center;
  gap: 5px;
  background: rgba(77, 166, 255, 0.08);
  border: 1px solid rgba(77, 166, 255, 0.25);
}

.route-link:hover {
  background: rgba(77, 166, 255, 0.18);
  border-color: var(--accent-blue);
}

.route-link.router-link-active {
  background: linear-gradient(135deg, rgba(77, 166, 255, 0.3), rgba(168, 85, 247, 0.25));
  border-color: var(--accent-blue);
  color: #fff;
}

.route-link.router-link-active::after {
  width: 0;
}

.route-icon {
  font-size: 0.85rem;
}

.home-back {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.nav-toggle {
  display: none;
  flex-direction: column;
  gap: 5px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px;
}

.nav-toggle span {
  display: block;
  width: 24px;
  height: 2px;
  border-radius: 2px;
  background: var(--text-primary);
  transition: var(--transition);
}

@media (max-width: 768px) {
  .nav-toggle {
    display: flex;
  }

  .nav-links {
    position: fixed;
    top: 64px;
    left: 0;
    right: 0;
    flex-direction: column;
    gap: 4px;
    background: rgba(10, 14, 26, 0.95);
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    padding: 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
    transform: translateY(-120%);
    opacity: 0;
    transition: all 0.3s ease;
    pointer-events: none;
  }

  .nav-links.open {
    transform: translateY(0);
    opacity: 1;
    pointer-events: auto;
  }

  .nav-links li {
    width: 100%;
  }

  .nav-links a {
    display: block;
    padding: 14px 16px;
    border-radius: 10px;
  }
}
</style>
