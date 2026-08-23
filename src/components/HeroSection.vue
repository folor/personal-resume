<template>
  <section id="hero" class="hero">
    <div class="container">
      <div class="hero-text reveal">
        <div class="hero-badge">
          <span class="badge-dot"></span>
          5 年前端开发经验 · 求职中 · 成都
        </div>

        <h1 class="hero-title">
          <span class="hero-greeting">你好，我是</span>
          <span class="hero-name">杨素</span>
        </h1>

        <p class="hero-role">
          <span class="gradient-text">前端开发工程师</span>
          <span class="divider">/</span>
          <span>WebGL 3D 交互开发</span>
        </p>

        <p class="hero-desc">
          精通 Vue2/3 与 React 技术栈，擅长 B 端复杂业务架构与 RBAC 权限体系；
          熟练掌握 Three.js 3D 交互开发，具备性能优化与 AI 工程化落地经验。
        </p>

        <!-- 核心技术栈速览 -->
        <div class="hero-techs">
          <span v-for="t in techs" :key="t.name" class="tech-chip" :style="{ borderColor: t.color + '55', color: t.color }">
            {{ t.name }}
          </span>
        </div>

        <!-- 快速数据 -->
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">5+</span>
            <span class="stat-label">年开发经验</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">9</span>
            <span class="stat-label">项目经验</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">10+</span>
            <span class="stat-label">核心业务模块</span>
          </div>
          <div class="stat-item">
            <span class="stat-number">3</span>
            <span class="stat-label">端跨端开发</span>
          </div>
        </div>

        <div class="hero-actions">
          <a href="#experience" class="btn btn-primary">
            查看经历
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </a>
          <a href="#contact" class="btn btn-secondary">联系我</a>
          <a href="mailto:yangsugogogo@qq.com" class="btn btn-secondary btn-icon" title="发邮件">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="4" width="20" height="16" rx="2"/>
              <path d="M22 7l-10 6L2 7"/>
            </svg>
          </a>
          <a href="tel:18482312017" class="btn btn-secondary btn-icon" title="打电话">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
            </svg>
          </a>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const techs = [
  { name: 'Vue2/3', color: '#42b883' },
  { name: 'React', color: '#61dafb' },
  { name: 'TypeScript', color: '#3178c6' },
  { name: 'Three.js', color: '#4da6ff' },
  { name: 'uni-app', color: '#22d3ee' },
  { name: 'Node.js', color: '#f59e0b' },
  { name: 'ECharts', color: '#a855f7' },
  { name: 'WebSocket', color: '#ec4899' }
]
</script>

<style scoped>
.hero {
  position: relative;
  padding: 140px 0 64px;
  overflow: hidden;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  border-radius: 100px;
  background: rgba(77, 166, 255, 0.08);
  border: 1px solid rgba(77, 166, 255, 0.2);
  color: var(--accent-blue);
  font-size: 0.85rem;
  font-weight: 500;
  margin-bottom: 24px;
}

.badge-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent-cyan);
  box-shadow: 0 0 8px var(--accent-cyan);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.hero-title {
  margin-bottom: 14px;
}

.hero-greeting {
  display: block;
  font-size: 1.2rem;
  color: var(--text-secondary);
  font-weight: 500;
  margin-bottom: 4px;
}

.hero-name {
  display: block;
  font-size: 3.2rem;
  font-weight: 900;
  letter-spacing: -0.02em;
  line-height: 1.1;
  background: var(--gradient-text);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-role {
  font-size: 1.2rem;
  font-weight: 600;
  margin-bottom: 18px;
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.gradient-text {
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.divider {
  color: var(--text-muted);
}

.hero-desc {
  color: var(--text-secondary);
  font-size: 1rem;
  line-height: 1.8;
  max-width: 620px;
  margin-bottom: 24px;
}

/* 技术栈速览 */
.hero-techs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 28px;
}

.tech-chip {
  padding: 5px 14px;
  border-radius: 100px;
  font-size: 0.85rem;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

/* 快速数据 */
.hero-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  max-width: 560px;
  margin-bottom: 32px;
}

.stat-item {
  text-align: center;
  padding: 14px 8px;
  border-radius: var(--radius-md);
  background: var(--bg-card);
  border: 1px solid var(--border-glass);
}

.stat-number {
  display: block;
  font-size: 1.5rem;
  font-weight: 900;
  background: var(--gradient-text);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.stat-label {
  display: block;
  font-size: 0.72rem;
  color: var(--text-muted);
  margin-top: 2px;
}

.hero-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 26px;
  border-radius: 12px;
  font-size: 0.95rem;
  font-weight: 600;
  text-decoration: none;
  transition: var(--transition);
  cursor: pointer;
  border: none;
}

.btn-primary {
  background: var(--gradient-primary);
  color: #fff;
  box-shadow: 0 4px 20px rgba(77, 166, 255, 0.3);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(77, 166, 255, 0.45);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-primary);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-2px);
}

.btn-icon {
  padding: 12px 14px;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .hero {
    padding: 110px 0 48px;
  }

  .hero-name {
    font-size: 2.3rem;
  }

  .hero-role {
    font-size: 1.05rem;
  }

  .hero-desc {
    font-size: 0.95rem;
  }

  .hero-stats {
    grid-template-columns: repeat(2, 1fr);
  }

  .hero-actions .btn:not(.btn-icon) {
    flex: 1;
    justify-content: center;
  }
}

@media (max-width: 420px) {
  .hero-badge {
    font-size: 0.75rem;
    padding: 5px 12px;
  }

  .hero-name {
    font-size: 2rem;
  }

  .tech-chip {
    font-size: 0.78rem;
    padding: 4px 11px;
  }
}
</style>
