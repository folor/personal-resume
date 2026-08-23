<template>
  <section id="skills" class="section section-compact">
    <div class="container">
      <div class="reveal section-head">
        <h2 class="section-title">技术<span class="gradient-text">栈</span></h2>
        <p class="section-subtitle">我擅长并持续精进的技术领域</p>
      </div>

      <div class="skills-grid">
        <div
          v-for="skill in skills"
          :key="skill.title"
          class="skill-card glass-card reveal"
        >
          <div class="skill-head">
            <div class="skill-icon" :style="{ background: skill.bg }">{{ skill.icon }}</div>
            <h3 class="skill-title">{{ skill.title }}</h3>
          </div>
          <p class="skill-desc">{{ skill.desc }}</p>
          <div class="skill-tags">
            <span v-for="tag in skill.tags" :key="tag" class="skill-tag">{{ tag }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const skills = [
  {
    icon: '⚡',
    title: '前端基础与跨端',
    desc: '精通 Vue2/3 及 React 技术栈，深入掌握 Vue 全家桶与 B 端组件库，具备 uni-app 跨端及 Cordova TV 大屏开发经验',
    tags: ['Vue2/3', 'React', 'TypeScript', 'ES6+', 'Pinia', 'Element Plus', 'uni-app'],
    bg: 'linear-gradient(135deg, #42b88320, #35495e20)'
  },
  {
    icon: '🎮',
    title: 'WebGL 与 3D 交互',
    desc: '熟练掌握 Three.js 核心开发：GLTF 模型加载、骨骼动画、Raycaster 射线拾取；理解 WebGL 渲染管线，能规避显存泄漏并优化复杂场景性能',
    tags: ['Three.js', 'WebGL', 'GLTF', '骨骼动画', 'Raycaster', 'WebSocket'],
    bg: 'linear-gradient(135deg, #4da6ff20, #a855f720)'
  },
  {
    icon: '🏗️',
    title: 'B 端复杂业务架构',
    desc: '具备企业级后台前端架构独立搭建能力，熟练实现 RBAC 三级权限体系与动态路由，擅长复杂业务表单、树形物料结构及流程审批模块',
    tags: ['RBAC', '动态路由', '复杂表单', '树形结构', '流程审批'],
    bg: 'linear-gradient(135deg, #3178c620, #235a9720)'
  },
  {
    icon: '🚀',
    title: '前端工程化与构建',
    desc: '精通 Webpack/Vite 构建配置，具备多环境 CI/CD 自动化部署及线上异常监控体系搭建经验，熟悉 Node.js 简易服务开发',
    tags: ['Webpack', 'Vite', 'CI/CD', '异常监控', 'Node.js'],
    bg: 'linear-gradient(135deg, #f59e0b20, #ef444420)'
  },
  {
    icon: '📊',
    title: '性能优化',
    desc: '在大数据表格渲染、大文件分片断点续传、防抖节流及懒加载方面有成熟落地方案，保障复杂业务场景下的流畅体验',
    tags: ['分片上传', '断点续传', '虚拟滚动', '懒加载', '防抖节流', 'MD5 秒传'],
    bg: 'linear-gradient(135deg, #22d3ee20, #0ea5e920)'
  },
  {
    icon: '🤖',
    title: 'AI 工程化应用',
    desc: '善于利用 AI 编程工具辅助需求拆解、代码重构与自动化测试；有大模型流式接口封装、断线重连与多厂商接口适配层搭建经验',
    tags: ['大模型流式', 'WebSocket', '多模态适配', 'AI 编程工具', '自动化测试'],
    bg: 'linear-gradient(135deg, #ec489920, #8b5cf620)'
  }
]
</script>

<style scoped>
.skills-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.skill-card {
  padding: 24px 22px;
  position: relative;
  overflow: hidden;
}

.skill-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--gradient-primary);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.4s ease;
}

.skill-card:hover::before {
  transform: scaleX(1);
}

.skill-head {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 12px;
}

.skill-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.4rem;
  flex-shrink: 0;
}

.skill-title {
  font-size: 1.1rem;
  font-weight: 700;
}

.skill-desc {
  color: var(--text-secondary);
  font-size: 0.9rem;
  line-height: 1.7;
  margin-bottom: 14px;
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.skill-tag {
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .skills-grid {
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .skill-card {
    padding: 18px 16px;
  }

  .skill-desc {
    font-size: 0.86rem;
    margin-bottom: 10px;
  }
}
</style>
