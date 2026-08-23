<template>
  <section id="contact" class="section section-compact">
    <div class="container">
      <div class="reveal section-head">
        <h2 class="section-title">联系<span class="gradient-text">我</span></h2>
        <p class="section-subtitle">有合作意向或技术交流？随时联系我</p>
      </div>

      <div class="contact-grid">
        <a
          v-for="item in contacts"
          :key="item.label"
          :href="item.href"
          class="contact-card glass-card reveal"
        >
          <div class="contact-icon" :style="{ background: item.bg }">{{ item.icon }}</div>
          <div class="contact-info">
            <span class="contact-label">{{ item.label }}</span>
            <span class="contact-value">{{ item.value }}</span>
          </div>
          <svg class="contact-arrow" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14M12 5l7 7-7 7"/>
          </svg>
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const contacts = [
  {
    icon: '✉️',
    label: '邮箱',
    value: 'yangsugogogo@qq.com',
    href: 'mailto:yangsugogogo@qq.com',
    bg: 'linear-gradient(135deg, #4da6ff20, #3178c620)'
  },
  {
    icon: '📱',
    label: '电话',
    value: '184-8231-2017',
    href: 'tel:18482312017',
    bg: 'linear-gradient(135deg, #42b88320, #22d3ee20)'
  },
  {
    icon: '🎓',
    label: '教育背景',
    value: '西南民族大学 · 通信工程（本科）',
    href: '#experience',
    bg: 'linear-gradient(135deg, #22d3ee20, #0ea5e920)'
  },
  {
    icon: '📍',
    label: '所在地',
    value: '中国 · 成都',
    href: '#',
    bg: 'linear-gradient(135deg, #ec489920, #a855f720)'
  }
]
</script>

<style scoped>
.contact-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  max-width: 760px;
  margin: 0 auto;
}

.contact-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px 22px;
  text-decoration: none;
  color: inherit;
}

.contact-icon {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.35rem;
  flex-shrink: 0;
}

.contact-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.contact-label {
  font-size: 0.76rem;
  color: var(--text-muted);
  font-weight: 500;
}

.contact-value {
  font-size: 0.98rem;
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.contact-arrow {
  color: var(--text-muted);
  transition: var(--transition);
  flex-shrink: 0;
}

.contact-card:hover .contact-arrow {
  color: var(--accent-blue);
  transform: translateX(4px);
}

@media (max-width: 768px) {
  .contact-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .contact-card {
    padding: 14px 16px;
  }

  .contact-value {
    font-size: 0.9rem;
  }
}
</style>
