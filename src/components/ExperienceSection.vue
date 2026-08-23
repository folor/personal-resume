<template>
  <section id="experience" class="section">
    <div class="container">
      <div class="reveal section-head">
        <h2 class="section-title">工作<span class="gradient-text">经历</span></h2>
        <p class="section-subtitle">5 年前端开发 · 从工程师到前端组长</p>
      </div>

      <div class="timeline">
        <div
          v-for="(item, i) in experiences"
          :key="i"
          class="timeline-item reveal"
          :class="{ edu: item.edu }"
        >
          <div class="timeline-dot"></div>
          <div class="timeline-content glass-card">
            <div class="timeline-top">
              <h3 class="timeline-role">{{ item.role }}</h3>
              <span class="timeline-date">{{ item.date }}</span>
            </div>
            <span class="timeline-company">{{ item.company }}</span>
            <p class="timeline-desc">{{ item.desc }}</p>
            <div class="timeline-techs">
              <span v-for="tech in item.techs" :key="tech" class="tech-chip">{{ tech }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const experiences = [
  {
    date: '2023.05 — 2026.05',
    role: '前端开发工程师 · 前端组长',
    company: '成都星柒玩科技有限公司',
    desc: '独立负责 AI 集成平台（量子系统）全量前端开发与架构落地，基于 Vue3 + Three.js 实现 3D 场景交互与 AI 大模型流式通信；独立搭建 MC 管理系统后台架构，攻克大文件分片上传、断点续传难题；作为前端组长主导三款游戏 GM 运营中台建设，统筹架构设计与团队分工，保障多个复杂项目并行高质量交付。',
    techs: ['Vue3', 'Three.js', 'WebSocket', 'WangEditor5', '大文件分片上传', 'Node.js'],
    edu: false
  },
  {
    date: '2021.07 — 2023.04',
    role: '前端开发工程师',
    company: '成都数拓科技有限公司',
    desc: '参与智企进销存系统开发，基于 React 实现多级 BOM 树形物料管理与 RBAC 权限体系落地；负责金地数字可视化大屏，封装自适应图表组件库并搭建全局异常监控；独立完成 React 后台权限管理系统，设计双重 Token 校验与无感续登机制。',
    techs: ['React', 'Vue3', 'ECharts', 'RBAC', 'Vite', 'Redux'],
    edu: false
  },
  {
    date: '2017.09 — 2021.06',
    role: '通信工程 · 本科',
    company: '西南民族大学',
    desc: '系统学习通信工程专业知识，打下扎实的数理与编程基础，在校期间自学前端技术并完成多个实践项目。',
    techs: [],
    edu: true
  }
]
</script>

<style scoped>
.timeline {
  position: relative;
  max-width: 760px;
  margin: 0 auto;
  padding-left: 32px;
}

.timeline::before {
  content: '';
  position: absolute;
  top: 8px;
  bottom: 8px;
  left: 7px;
  width: 2px;
  background: linear-gradient(var(--accent-blue), var(--accent-purple), transparent);
}

.timeline-item {
  position: relative;
  padding-bottom: 28px;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-dot {
  position: absolute;
  top: 26px;
  left: -31px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--bg-primary);
  border: 3px solid var(--accent-blue);
  box-shadow: 0 0 12px rgba(77, 166, 255, 0.5);
  z-index: 2;
}

.timeline-item.edu .timeline-dot {
  border-color: var(--accent-purple);
  box-shadow: 0 0 12px rgba(168, 85, 247, 0.5);
}

.timeline-content {
  padding: 20px 22px;
}

.timeline-item.edu .timeline-content {
  padding: 16px 22px;
  background: rgba(168, 85, 247, 0.04);
}

.timeline-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 2px;
}

.timeline-role {
  font-size: 1.1rem;
  font-weight: 700;
}

.timeline-date {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--accent-blue);
  padding: 3px 12px;
  border-radius: 6px;
  background: rgba(77, 166, 255, 0.1);
  white-space: nowrap;
}

.timeline-item.edu .timeline-date {
  color: var(--accent-purple);
  background: rgba(168, 85, 247, 0.1);
}

.timeline-company {
  display: block;
  font-size: 0.9rem;
  color: var(--text-muted);
  margin-bottom: 10px;
}

.timeline-desc {
  color: var(--text-secondary);
  font-size: 0.92rem;
  line-height: 1.75;
  margin-bottom: 12px;
}

.timeline-techs {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tech-chip {
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .timeline {
    padding-left: 26px;
  }

  .timeline-dot {
    left: -25px;
    top: 24px;
  }

  .timeline-content {
    padding: 16px;
  }

  .timeline-role {
    font-size: 1rem;
  }

  .timeline-desc {
    font-size: 0.88rem;
    line-height: 1.7;
  }

  .timeline-top {
    flex-direction: column;
    align-items: flex-start;
    gap: 6px;
  }
}
</style>
