<template>
  <!-- 动态背景光斑 -->
  <div class="bg-orbs">
    <div class="bg-orb"></div>
    <div class="bg-orb"></div>
    <div class="bg-orb"></div>
  </div>

  <!-- 导航栏 -->
  <NavBar />

  <!-- 主内容：经历与技术栈优先，3D Demo 后置辅助 -->
  <main>
    <HeroSection />
    <SkillsSection />
    <ExperienceSection />
    <ProjectsSection />
    <DemoSection />
    <ContactSection />
  </main>

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

<script setup lang="ts">
import { onMounted } from 'vue'
import NavBar from './components/NavBar.vue'
import HeroSection from './components/HeroSection.vue'
import SkillsSection from './components/SkillsSection.vue'
import ExperienceSection from './components/ExperienceSection.vue'
import ProjectsSection from './components/ProjectsSection.vue'
import DemoSection from './components/DemoSection.vue'
import ContactSection from './components/ContactSection.vue'

const year = new Date().getFullYear()

onMounted(() => {
  // 滚动入场动画
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible')
        }
      })
    },
    { threshold: 0.1, rootMargin: '0px 0px -60px 0px' }
  )

  document.querySelectorAll('.reveal').forEach((el) => observer.observe(el))
})
</script>

<style scoped>
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

@media (max-width: 480px) {
  .footer-copy {
    font-size: 0.72rem;
    line-height: 1.7;
  }
}
</style>
