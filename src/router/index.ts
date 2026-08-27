import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomePage.vue'),
    },
    {
      path: '/agent',
      name: 'agent',
      component: () => import('../views/AgentRoadmapPage.vue'),
      meta: { tool: true },
    },
    {
      path: '/daily',
      name: 'daily',
      component: () => import('../views/DailyPlanPage.vue'),
      meta: { tool: true },
    },
    {
      path: '/geo',
      name: 'geo',
      component: () => import('../views/GeoRoadmapPage.vue'),
      meta: { tool: true },
    },
  ],
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, behavior: 'smooth' }
    return { top: 0 }
  },
})

export default router
