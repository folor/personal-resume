<template>
  <section id="demo" class="section section-compact">
    <div class="container">
      <div class="reveal section-head">
        <h2 class="section-title">3D <span class="gradient-text">Demo</span></h2>
        <p class="section-subtitle">
          Three.js 交互演示 · GLTF 骨骼动画 / Raycaster 拾取 / 粒子特效 / 后期处理
        </p>
      </div>

      <div class="demo-wrapper glass-card reveal">
        <div class="demo-header">
          <span class="demo-dots">
            <span></span><span></span><span></span>
          </span>
          <span class="demo-title">Three.js 3D 交互演示</span>
          <span class="demo-badge">LIVE</span>
        </div>
        <div class="demo-container">
          <ThreeDemo />
        </div>
        <div class="demo-hints">
          <span>🖱️ 点击地面移动</span>
          <span>✋ 长按拖拽角色</span>
          <span>🔮 滚轮缩放视角</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import ThreeDemo from './ThreeDemo.vue'
</script>

<style scoped>
.demo-wrapper {
  overflow: hidden;
  border-color: rgba(77, 166, 255, 0.15);
  max-width: 960px;
  margin: 0 auto;
}

.demo-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.demo-dots {
  display: flex;
  gap: 6px;
}

.demo-dots span {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.demo-dots span:nth-child(1) { background: #ff5f57; }
.demo-dots span:nth-child(2) { background: #febc2e; }
.demo-dots span:nth-child(3) { background: #28c840; }

.demo-title {
  font-size: 0.85rem;
  color: var(--text-muted);
  flex: 1;
  text-align: center;
}

.demo-badge {
  font-size: 0.7rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(236, 72, 153, 0.15);
  color: var(--accent-pink);
  letter-spacing: 0.05em;
}

.demo-container {
  border: none !important;
}

.demo-container :deep(> div) {
  border: none !important;
  border-radius: 0 !important;
  height: 520px !important;
}

.demo-hints {
  display: flex;
  gap: 20px;
  padding: 10px 16px;
  background: rgba(0, 0, 0, 0.2);
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 0.8rem;
  color: var(--text-muted);
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .demo-container :deep(> div) {
    height: 380px !important;
  }

  .demo-hints {
    gap: 12px;
    font-size: 0.72rem;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .demo-container :deep(> div) {
    height: 320px !important;
  }
}
</style>
