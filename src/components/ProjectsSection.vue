<template>
  <section id="projects" class="section section-compact">
    <div class="container">
      <div class="reveal section-head">
        <h2 class="section-title">项目<span class="gradient-text">经验</span></h2>
        <p class="section-subtitle">9 个核心项目 · 覆盖 B 端中台 / 3D / 跨端 / AI</p>
      </div>

      <div class="projects-list">
        <div
          v-for="project in projects"
          :key="project.title"
          class="project-row glass-card reveal"
        >
          <div class="project-icon" :style="{ background: project.gradient }">{{ project.icon }}</div>
          <div class="project-main">
            <div class="project-title-line">
              <h3 class="project-title">{{ project.title }}</h3>
              <span v-if="project.role" class="project-role" :class="project.roleClass">{{ project.role }}</span>
            </div>
            <p class="project-desc">{{ project.desc }}</p>
            <div class="project-techs">
              <span v-for="tech in project.techs" :key="tech" class="tech-chip">{{ tech }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const projects = [
  {
    icon: '🤖',
    title: 'AI 集成平台（量子系统）',
    role: '独立开发',
    roleClass: 'role-star',
    desc: 'Vue3 全量前端架构 + RBAC 三级权限；Three.js 3D 交互（GLTF 骨骼动画、Raycaster 拾取拖拽、显存优化）；封装大模型 WebSocket 流式通信，AI 文本秒级实时渲染。',
    techs: ['Vue3', 'Three.js', 'WebSocket', 'ECharts', 'JSZip'],
    gradient: 'linear-gradient(135deg, #4da6ff, #a855f7)'
  },
  {
    icon: '🎮',
    title: 'MC 管理系统',
    role: '独立开发',
    roleClass: 'role-star',
    desc: '游戏运维后台架构与 10+ 核心模块；大文件分片上传、MD5 秒传与断点续传；WebSocket 秒级推送；深度二开 WangEditor5，封装高频组件减少 30% 重复代码。',
    techs: ['Vue3', 'WebSocket', 'WangEditor5', '分片上传'],
    gradient: 'linear-gradient(135deg, #22d3ee, #0ea5e9)'
  },
  {
    icon: '🕹️',
    title: '三款游戏 GM 运营中台',
    role: '前端组长',
    roleClass: 'role-lead',
    desc: '主导中台架构设计，抽象三款游戏共性业务搭建高复用标准化中台；支撑活动配置、邮件批量推送、玩家数据筛查；大数据列表优化，单页稳定 500+ 条无卡顿。',
    techs: ['Vue2', 'Vue-Router', 'Vuex', '列表优化'],
    gradient: 'linear-gradient(135deg, #42b883, #35495e)'
  },
  {
    icon: '🏛️',
    title: '广安安心办政务小程序及自助终端',
    role: '',
    roleClass: '',
    desc: 'uni-app 独立完成政务小程序、自助大屏终端、PC 管理后台三套前端；重构纯触控大屏交互逻辑，解决触屏偏移与手势冲突；封装政务通用组件库。',
    techs: ['Vue', 'uni-app', '多端适配', '触控优化'],
    gradient: 'linear-gradient(135deg, #ec4899, #8b5cf6)'
  },
  {
    icon: '📺',
    title: 'TV 大屏安卓应用',
    role: '全栈开发',
    roleClass: 'role-lead',
    desc: 'Cordova 完成 Web 项目原生编译与 APK 打包发布；专项适配纯遥控器操作（焦点流转、大尺寸热区）；独立搭建 Node.js 中间层服务支撑 TV 端动态渲染。',
    techs: ['Vue', 'Cordova', 'Node.js', 'TV 适配'],
    gradient: 'linear-gradient(135deg, #f59e0b, #ef4444)'
  },
  {
    icon: '🎬',
    title: 'uni-app 短视频小程序',
    role: '',
    roleClass: '',
    desc: '付费视频小程序全量前端：视频上传、广告位配置、在线支付；优化预加载与懒加载，首屏加载提升 22%；全机型多尺寸适配。',
    techs: ['uni-app', '多端适配', '懒加载', '支付对接'],
    gradient: 'linear-gradient(135deg, #10b981, #06b6d4)'
  },
  {
    icon: '📦',
    title: '智企进销存系统',
    role: '',
    roleClass: '',
    desc: '递归算法解析扁平 BOM 数据，构建高性能多级树形物料组件；RBAC 页面级/按钮级管控；封装批量导入导出通用单据组件 + ECharts 库存可视化。',
    techs: ['React', 'Axios', 'RBAC', 'ECharts', '树形递归'],
    gradient: 'linear-gradient(135deg, #6366f1, #8b5cf6)'
  },
  {
    icon: '📊',
    title: '金地数字可视化大屏',
    role: '',
    roleClass: '',
    desc: '地产投拓数据大屏：封装自适应图表组件库，多分辨率适配与明暗主题切换；Suspense 异步懒加载；重写全局 errorHandler 统一异常捕获上报。',
    techs: ['Vue3', 'ECharts', 'Suspense', '异常监控'],
    gradient: 'linear-gradient(135deg, #0ea5e9, #6366f1)'
  },
  {
    icon: '🔐',
    title: 'React 后台权限管理系统',
    role: '',
    roleClass: '',
    desc: 'React + Redux 搭建后台骨架，双重 Token 校验与无感续登；递归算法转树形部门结构，路由拦截器实现动态菜单渲染。',
    techs: ['React', 'Redux', 'RBAC', 'Token 无感续登'],
    gradient: 'linear-gradient(135deg, #14b8a6, #0ea5e9)'
  }
]
</script>

<style scoped>
.projects-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
  max-width: 900px;
  margin: 0 auto;
}

.project-row {
  display: flex;
  gap: 18px;
  padding: 18px 20px;
  align-items: flex-start;
}

.project-icon {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.6rem;
  flex-shrink: 0;
}

.project-main {
  flex: 1;
  min-width: 0;
}

.project-title-line {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 6px;
}

.project-title {
  font-size: 1.05rem;
  font-weight: 700;
}

.project-role {
  font-size: 0.7rem;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: 100px;
  white-space: nowrap;
}

.role-star {
  background: rgba(255, 215, 0, 0.12);
  color: #ffd700;
  border: 1px solid rgba(255, 215, 0, 0.3);
}

.role-lead {
  background: rgba(77, 166, 255, 0.12);
  color: var(--accent-blue);
  border: 1px solid rgba(77, 166, 255, 0.3);
}

.project-desc {
  color: var(--text-secondary);
  font-size: 0.88rem;
  line-height: 1.65;
  margin-bottom: 10px;
}

.project-techs {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tech-chip {
  padding: 2px 9px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .project-row {
    padding: 14px;
    gap: 12px;
  }

  .project-icon {
    width: 42px;
    height: 42px;
    border-radius: 11px;
    font-size: 1.3rem;
  }

  .project-title {
    font-size: 0.98rem;
  }

  .project-desc {
    font-size: 0.84rem;
    line-height: 1.6;
  }
}

@media (max-width: 480px) {
  .project-row {
    flex-direction: column;
    gap: 10px;
  }
}
</style>
