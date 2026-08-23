<template>
  <section id="demo" class="section section-compact">
    <div class="container">
      <div class="reveal section-head">
        <h2 class="section-title">3D <span class="gradient-text">Demo</span></h2>
        <p class="section-subtitle">
          Three.js 交互演示 · 智慧城市可视化 / 地形地貌飞行巡航 / GLTF 骨骼动画 / Raycaster 拾取 / 粒子特效 / Bloom 后期处理
        </p>
      </div>

      <div class="demo-wrapper glass-card reveal">
        <div class="demo-header">
          <span class="demo-dots">
            <span></span><span></span><span></span>
          </span>
          <div class="demo-tabs">
            <button
              :class="['demo-tab', { active: activeDemo === 'city' }]"
              @click="activeDemo = 'city'"
            >
              🏙️ 智慧城市
            </button>
            <button
              :class="['demo-tab', { active: activeDemo === 'terrain' }]"
              @click="activeDemo = 'terrain'"
            >
              ✈️ 地形巡航
            </button>
            <button
              :class="['demo-tab', { active: activeDemo === 'bot' }]"
              @click="activeDemo = 'bot'"
            >
              🤖 角色交互
            </button>
          </div>
          <span class="demo-badge">LIVE</span>
        </div>
        <div class="demo-container">
          <SmartCityDemo v-if="activeDemo === 'city'" />
          <TerrainFlightDemo v-else-if="activeDemo === 'terrain'" />
          <ThreeDemo v-else />
        </div>
        <div class="demo-hints" v-if="activeDemo === 'city'">
          <span>🖱️ 悬停查看楼宇数据</span>
          <span>👆 点击楼宇聚焦 / 点击空白复位</span>
          <span>🖱️ 拖拽旋转 · 滚轮缩放</span>
        </div>
        <div class="demo-hints" v-else-if="activeDemo === 'terrain'">
          <span>🖱️ 点击地形设定巡航航点</span>
          <span>🕹️ 手动模式：W/S 油门 · A/D 转向 · ↑↓ 俯仰</span>
          <span>🎥 工具栏切换视角</span>
        </div>
        <div class="demo-hints" v-else>
          <span>🕹️ WASD 移动选中角色</span>
          <span>🖱️ 点击地面移动 · 点击角色切换选中</span>
          <span>✋ 长按拖拽 · 点击能量球拾取计分</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import SmartCityDemo from './SmartCityDemo.vue'
import TerrainFlightDemo from './TerrainFlightDemo.vue'
import ThreeDemo from './ThreeDemo.vue'

const activeDemo = ref<'city' | 'terrain' | 'bot'>('city')
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
  padding: 10px 16px;
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

.demo-tabs {
  flex: 1;
  display: flex;
  justify-content: center;
  gap: 6px;
}

.demo-tab {
  padding: 6px 16px;
  font-size: 0.8rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.04);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.2s;
}

.demo-tab:hover {
  color: var(--text-primary);
  border-color: rgba(255, 255, 255, 0.25);
}

.demo-tab.active {
  background: linear-gradient(135deg, rgba(77, 166, 255, 0.25), rgba(168, 85, 247, 0.25));
  border-color: rgba(77, 166, 255, 0.6);
  color: var(--text-primary);
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

  .demo-tab {
    padding: 5px 10px;
    font-size: 0.72rem;
  }

  .demo-dots {
    display: none;
  }
}

@media (max-width: 480px) {
  .demo-container :deep(> div) {
    height: 320px !important;
  }

  .demo-tab {
    padding: 4px 8px;
    font-size: 0.68rem;
  }
}
</style>
