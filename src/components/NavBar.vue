<template>
  <nav class="navbar" :class="{ scrolled: isScrolled }">
    <div class="nav-container">
      <a href="#hero" class="nav-logo">
        <span class="logo-icon">YS</span>
        <span class="logo-text">杨素</span>
      </a>

      <ul class="nav-links" :class="{ open: menuOpen }">
        <li v-for="link in links" :key="link.href">
          <a :href="link.href" @click="menuOpen = false">{{ link.label }}</a>
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
import { ref, onMounted, onUnmounted } from 'vue'

const isScrolled = ref(false)
const menuOpen = ref(false)

const links = [
  { label: '首页', href: '#hero' },
  { label: '技术栈', href: '#skills' },
  { label: '经历', href: '#experience' },
  { label: '项目', href: '#projects' },
  { label: '3D Demo', href: '#demo' },
  { label: '联系', href: '#contact' }
]

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
  gap: 8px;
  list-style: none;
}

.nav-links a {
  display: block;
  padding: 8px 16px;
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.95rem;
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

.nav-links a:hover::after {
  width: 50%;
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

.nav-toggle span:nth-child(1) {
  transform: rotate(0deg);
}

.nav-toggle span:nth-child(3) {
  transform: rotate(0deg);
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
    gap: 0;
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

  .nav-links a:hover {
    background: rgba(255, 255, 255, 0.05);
  }
}
</style>
