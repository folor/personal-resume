<template>
  <div class="dp-page">
    <!-- 页头 -->
    <header class="dp-head">
      <div class="dp-head-row">
        <div>
          <div class="dp-badge">{{ currentRoadmap.icon }} {{ currentRoadmap.name || '90 天每日打卡' }}</div>
          <h1 class="dp-title">{{ currentRoadmap.desc || '今日任务，一天一天清零' }}</h1>
          <p class="dp-desc">
            数据存于本地 MySQL（roadmap 库），勾选即持久化。上午理论 · 下午实战 · 晚上复盘。
          </p>
        </div>
        <div class="dp-conn" :class="apiOnline ? 'on' : 'off'">
          <span class="conn-dot"></span>
          {{ apiOnline ? '后端已连接' : '后端未连接' }}
        </div>
      </div>

      <!-- 总进度条 -->
      <div class="dp-progress-card">
        <div class="dp-progress-top">
          <span class="dp-progress-label">总进度</span>
          <span class="dp-progress-num">
            <b>{{ stats.done_todos }}</b> / {{ stats.total_todos }} 项 ·
            <b>{{ stats.finished_days }}</b> / {{ stats.total_days }} 天完成 ·
            <b>{{ overallPercent }}%</b>
          </span>
        </div>
        <div class="dp-progress-bar">
          <div class="dp-progress-fill" :style="{ width: overallPercent + '%' }"></div>
        </div>
        <div class="dp-stage-bars">
          <div
            v-for="s in stats.stage_stats"
            :key="s.stage"
            class="stage-bar"
            :style="{ '--sb': stageColor[s.stage] }"
          >
            <span class="sb-name">{{ s.name }}</span>
            <div class="sb-track"><div class="sb-fill" :style="{ width: pct(s.done, s.total) + '%' }"></div></div>
            <span class="sb-num">{{ s.done }}/{{ s.total }}</span>
          </div>
        </div>
      </div>
    </header>

    <!-- 规划切换 -->
    <div class="dp-roadmaps" v-if="roadmaps.length > 1">
      <button
        v-for="rm in roadmaps"
        :key="rm.id"
        class="rm-tab"
        :class="{ active: currentRoadmapId === rm.id }"
        :style="currentRoadmapId === rm.id ? { background: rm.color, borderColor: rm.color } : {}"
        @click="switchRoadmap(rm.id)"
      >
        <span class="rm-icon">{{ rm.icon }}</span>
        <span class="rm-name">{{ rm.short }}</span>
      </button>
    </div>

    <!-- 过滤器 -->
    <div class="dp-filters">
      <div class="filter-tabs">
        <button
          v-for="f in filters"
          :key="f.value"
          class="filter-tab"
          :class="{ active: filter === f.value }"
          :style="filter === f.value && f.color ? { background: f.color, borderColor: f.color } : {}"
          @click="filter = f.value"
        >
          {{ f.label }}
        </button>
      </div>
      <input v-model="keyword" class="dp-search" placeholder="🔍 搜索任务内容 / 主题…" />
    </div>

    <!-- 加载/错误状态 -->
    <div v-if="loading" class="dp-state">加载中…</div>
    <div v-else-if="error" class="dp-state error">
      ⚠️ {{ error }}
      <p class="state-hint">
        请确认 Go 后端已启动：<code>cd server && go run main.go</code>（默认端口 8080）
      </p>
    </div>

    <!-- 天列表 -->
    <template v-else>
      <div v-for="g in groupedDays" :key="g.key" class="dp-group">
        <div class="dp-group-head" :style="{ '--gc': stageColor[g.stage] }">
          <span class="group-title">{{ g.title }}</span>
          <span class="group-meta">{{ g.days.length }} 天 · 完成 {{ g.doneSum }}/{{ g.totalSum }} 项</span>
        </div>

        <div v-for="d in g.days" :key="d.day_no" class="dp-day" :class="{ done: d.done_cnt === d.total && d.total > 0 }">
          <button class="dp-day-head" @click="toggleDay(d.day_no)">
            <span class="day-badge" :style="{ background: stageColor[d.stage] }">
              {{ d.day_no === 0 ? 'D0' : d.day_no }}
            </span>
            <div class="day-info">
              <span class="day-theme">{{ d.theme }}</span>
              <span class="day-output">🎯 产出：{{ d.output }}</span>
            </div>
            <div class="day-right">
              <span class="day-count" :class="{ all: d.done_cnt === d.total }">
                {{ d.done_cnt }}/{{ d.total }}
              </span>
              <span class="day-arrow" :class="{ open: expandedDays.has(d.day_no) }">▾</span>
            </div>
          </button>

          <div v-if="expandedDays.has(d.day_no)" class="dp-day-body">
            <label
              v-for="t in d.todos"
              :key="t.id"
              class="dp-todo"
              :class="{ checked: t.done, [`slot-${t.slot}`]: true }"
            >
              <input
                type="checkbox"
                :checked="t.done"
                :disabled="pendingIds.has(t.id)"
                @change="toggleTodo(t)"
              />
              <span class="todo-check"></span>
              <span class="todo-slot">{{ slotMeta[t.slot]?.icon }} {{ t.slot_name }}</span>
              <span class="todo-content">{{ t.content }}</span>
            </label>
          </div>
        </div>
      </div>
      <div v-if="groupedDays.length === 0" class="dp-state">没有匹配的任务</div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface Todo {
  id: number
  day_no: number
  slot: string
  slot_name: string
  content: string
  done: boolean
}
interface Day {
  day_no: number
  week_no: number
  stage: number
  theme: string
  output: string
  todos: Todo[]
  done_cnt: number
  total: number
}
interface StageStat {
  stage: number
  name: string
  done: number
  total: number
}
interface Stats {
  total_todos: number
  done_todos: number
  total_days: number
  finished_days: number
  stage_stats: StageStat[]
}

// 本地开发走 vite 代理 /api；生产构建时可通过 .env.production 的 VITE_API_BASE 指向后端域名
const API = (import.meta.env.VITE_API_BASE as string | undefined) || '/api'

interface Roadmap {
  id: number
  name: string
  short: string
  color: string
  icon: string
  desc: string
  stage_names: string[]
}

const roadmaps = ref<Roadmap[]>([])
const currentRoadmapId = ref(1)
const currentRoadmap = computed(() =>
  roadmaps.value.find((r) => r.id === currentRoadmapId.value) || {
    id: 1, name: '', short: '', color: '#16a34a', icon: '✅', desc: '', stage_names: [],
  }
)

const days = ref<Day[]>([])
const stats = ref<Stats>({
  total_todos: 0, done_todos: 0, total_days: 0, finished_days: 0, stage_stats: [],
})
const loading = ref(true)
const error = ref('')
const apiOnline = ref(false)
const expandedDays = ref(new Set<number>([1]))
const pendingIds = ref(new Set<number>())
const filter = ref('all')
const keyword = ref('')

const stageColor: Record<number, string> = { 0: '#94a3b8', 1: '#16a34a', 2: '#0891b2', 3: '#7c3aed' }
const stageName = computed<Record<number, string>>(() => {
  const names = currentRoadmap.value.stage_names || []
  return {
    0: '第 0 天准备',
    1: names[0] || '阶段一',
    2: names[1] || '阶段二',
    3: names[2] || '阶段三',
  }
})
const slotMeta: Record<string, { icon: string; time: string }> = {
  am: { icon: '📖', time: '09:00-11:30' },
  pm: { icon: '💻', time: '14:00-17:30' },
  ev: { icon: '🔄', time: '19:00-21:00' },
  extra: { icon: '🎯', time: '21:00-21:30' },
}

const filters = [
  { label: '全部', value: 'all' },
  { label: '阶段一', value: 'stage1', color: stageColor[1] },
  { label: '阶段二', value: 'stage2', color: stageColor[2] },
  { label: '阶段三', value: 'stage3', color: stageColor[3] },
  { label: '未完成', value: 'undone' },
]

const overallPercent = computed(() => pct(stats.value.done_todos, stats.value.total_todos))

function pct(a: number, b: number) {
  return b === 0 ? 0 : Math.round((a / b) * 100)
}

const groupedDays = computed(() => {
  let list = days.value
  if (filter.value.startsWith('stage')) {
    const st = Number(filter.value.replace('stage', ''))
    list = list.filter((d) => d.stage === st)
  } else if (filter.value === 'undone') {
    list = list.filter((d) => d.done_cnt < d.total)
  }
  const kw = keyword.value.trim()
  if (kw) {
    list = list.filter(
      (d) =>
        d.theme.includes(kw) ||
        d.output.includes(kw) ||
        d.todos.some((t) => t.content.includes(kw))
    )
  }
  // 按阶段分组
  const groups: { key: string; stage: number; title: string; days: Day[]; doneSum: number; totalSum: number }[] = []
  for (const d of list) {
    let g = groups.find((x) => x.stage === d.stage)
    if (!g) {
      g = {
        key: 'stage-' + d.stage,
        stage: d.stage,
        title: stageName.value[d.stage] ?? '其他',
        days: [],
        doneSum: 0,
        totalSum: 0,
      }
      groups.push(g)
    }
    g.days.push(d)
    g.doneSum += d.done_cnt
    g.totalSum += d.total
  }
  return groups
})

function toggleDay(dayNo: number) {
  const s = new Set(expandedDays.value)
  if (s.has(dayNo)) s.delete(dayNo)
  else s.add(dayNo)
  expandedDays.value = s
}

async function toggleTodo(todo: Todo) {
  const next = !todo.done
  todo.done = next
  // 本地同步计数
  const day = days.value.find((d) => d.day_no === todo.day_no)
  if (day) day.done_cnt += next ? 1 : -1
  const pend = new Set(pendingIds.value)
  pend.add(todo.id)
  pendingIds.value = pend
  try {
    const res = await fetch(`${API}/todo/${todo.id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ done: next }),
    })
    if (!res.ok) throw new Error('保存失败')
    stats.value.done_todos += next ? 1 : -1
    // 更新阶段统计
    const dayInfo = days.value.find((d) => d.day_no === todo.day_no)
    if (dayInfo) {
      const ss = stats.value.stage_stats.find((x) => x.stage === dayInfo.stage)
      if (ss) ss.done += next ? 1 : -1
    }
    refreshFinishedDays()
  } catch (e) {
    todo.done = !next
    if (day) day.done_cnt += next ? -1 : 1
    error.value = '保存失败，请检查后端服务'
    setTimeout(() => (error.value = ''), 3000)
  } finally {
    const p = new Set(pendingIds.value)
    p.delete(todo.id)
    pendingIds.value = p
  }
}

function refreshFinishedDays() {
  stats.value.finished_days = days.value.filter((d) => d.total > 0 && d.done_cnt === d.total).length
}

async function loadRoadmaps() {
  try {
    const res = await fetch(`${API}/roadmaps`).then((r) => r.json())
    roadmaps.value = res as Roadmap[]
  } catch {
    // 后端未连接时静默
  }
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const rid = currentRoadmapId.value
    const [daysRes, statsRes] = await Promise.all([
      fetch(`${API}/days?roadmap=${rid}`).then((r) => r.json()),
      fetch(`${API}/stats?roadmap=${rid}`).then((r) => r.json()),
    ])
    days.value = daysRes as Day[]
    stats.value = statsRes as Stats
    apiOnline.value = true
    // 默认展开第一个未完成的 天
    const firstUndone = days.value.find((d) => d.done_cnt < d.total)
    expandedDays.value = new Set(firstUndone ? [firstUndone.day_no] : [1])
  } catch {
    apiOnline.value = false
    error.value = '无法连接后端 API（127.0.0.1:8080）'
  } finally {
    loading.value = false
  }
}

async function switchRoadmap(id: number) {
  currentRoadmapId.value = id
  filter.value = 'all'
  keyword.value = ''
  await load()
}

onMounted(async () => {
  await loadRoadmaps()
  await load()
})
</script>

<style scoped>
.dp-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 页头 */
.dp-head-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.dp-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 999px;
  background: #dcfce7;
  color: #15803d;
  font-size: 0.75rem;
  font-weight: 600;
  margin-bottom: 10px;
}

.dp-title {
  font-size: 1.6rem;
  font-weight: 800;
  color: #1a2e23;
  margin: 0 0 8px;
}

.dp-desc {
  color: #4a5d52;
  font-size: 0.9rem;
  margin: 0;
}

.dp-conn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 600;
  flex-shrink: 0;
  border: 1px solid #d1e7da;
  background: #fff;
}

.dp-conn.on { color: #15803d; }
.dp-conn.off { color: #dc2626; }

.conn-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.dp-conn.on .conn-dot { background: #22c55e; box-shadow: 0 0 6px #22c55e; }
.dp-conn.off .conn-dot { background: #ef4444; }

/* 进度卡 */
.dp-progress-card {
  margin-top: 18px;
  background: #fff;
  border: 1px solid #d1e7da;
  border-radius: 14px;
  padding: 18px 20px;
}

.dp-progress-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 10px;
}

.dp-progress-label {
  font-weight: 700;
  color: #1a2e23;
  font-size: 0.95rem;
}

.dp-progress-num {
  font-size: 0.85rem;
  color: #4a5d52;
}

.dp-progress-num b {
  color: #16a34a;
  font-size: 1rem;
}

.dp-progress-bar {
  height: 10px;
  border-radius: 999px;
  background: #e8f2eb;
  overflow: hidden;
}

.dp-progress-fill {
  height: 100%;
  border-radius: 999px;
  background: linear-gradient(90deg, #16a34a, #4ade80);
  transition: width 0.4s ease;
}

.dp-stage-bars {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 14px;
  margin-top: 16px;
}

.stage-bar {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 8px;
}

.sb-name {
  font-size: 0.72rem;
  font-weight: 600;
  color: #4a5d52;
  white-space: nowrap;
}

.sb-track {
  height: 6px;
  border-radius: 999px;
  background: #eef2f0;
  overflow: hidden;
}

.sb-fill {
  height: 100%;
  border-radius: 999px;
  background: var(--sb);
  transition: width 0.4s ease;
}

.sb-num {
  font-size: 0.72rem;
  color: #8a9b92;
  font-variant-numeric: tabular-nums;
}

/* 规划切换 */
.dp-roadmaps {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.rm-tab {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 9px 16px;
  border-radius: 10px;
  border: 1px solid #d1e7da;
  background: #fff;
  color: #1a2e23;
  font-size: 0.88rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.rm-tab:hover {
  border-color: #16a34a;
  transform: translateY(-1px);
}

.rm-tab.active {
  color: #fff;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.1);
}

.rm-icon {
  font-size: 1rem;
}

.rm-name {
  white-space: nowrap;
}

/* 过滤器 */
.dp-filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
}

.filter-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-tab {
  padding: 7px 16px;
  border-radius: 999px;
  border: 1px solid #d1e7da;
  background: #fff;
  color: #4a5d52;
  font-size: 0.82rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-tab:hover { border-color: #16a34a; color: #16a34a; }

.filter-tab.active {
  background: #16a34a;
  border-color: #16a34a;
  color: #fff;
  font-weight: 600;
}

.dp-search {
  flex: 1;
  min-width: 200px;
  max-width: 320px;
  padding: 8px 14px;
  border-radius: 10px;
  border: 1px solid #d1e7da;
  background: #fff;
  font-size: 0.85rem;
  color: #1a2e23;
  outline: none;
  transition: border-color 0.2s;
}

.dp-search:focus { border-color: #16a34a; }

/* 状态 */
.dp-state {
  background: #fff;
  border: 1px solid #d1e7da;
  border-radius: 12px;
  padding: 40px;
  text-align: center;
  color: #4a5d52;
}

.dp-state.error { color: #dc2626; }

.state-hint {
  margin-top: 10px;
  font-size: 0.82rem;
  color: #4a5d52;
}

.state-hint code {
  background: #f0fdf4;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.78rem;
}

/* 分组 */
.dp-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.dp-group-head {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 2px 4px;
}

.dp-group-head::before {
  content: '';
  width: 4px;
  height: 18px;
  border-radius: 2px;
  background: var(--gc);
}

.group-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: #1a2e23;
}

.group-meta {
  font-size: 0.75rem;
  color: #8a9b92;
}

/* 天卡片 */
.dp-day {
  background: #fff;
  border: 1px solid #d1e7da;
  border-radius: 12px;
  overflow: hidden;
  transition: box-shadow 0.2s;
}

.dp-day:hover { box-shadow: 0 4px 14px rgba(22, 163, 74, 0.08); }

.dp-day.done { border-color: #86efac; background: #f8fef9; }

.dp-day-head {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
  padding: 14px 16px;
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
}

.day-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 42px;
  height: 42px;
  border-radius: 10px;
  color: #fff;
  font-weight: 800;
  font-size: 0.9rem;
  flex-shrink: 0;
}

.day-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.day-theme {
  font-size: 0.95rem;
  font-weight: 700;
  color: #1a2e23;
}

.day-output {
  font-size: 0.75rem;
  color: #8a9b92;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.day-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.day-count {
  font-size: 0.8rem;
  font-weight: 600;
  color: #8a9b92;
  font-variant-numeric: tabular-nums;
  padding: 3px 10px;
  border-radius: 999px;
  background: #f0fdf4;
}

.day-count.all {
  background: #dcfce7;
  color: #15803d;
}

.day-arrow {
  color: #8a9b92;
  font-size: 0.8rem;
  transition: transform 0.25s;
}

.day-arrow.open { transform: rotate(180deg); }

/* 任务 */
.dp-day-body {
  border-top: 1px solid #e4efe8;
  padding: 10px 16px 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.dp-todo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.dp-todo:hover { background: #f4faf6; }

.dp-todo input { display: none; }

.todo-check {
  width: 20px;
  height: 20px;
  border-radius: 6px;
  border: 2px solid #cbd9d0;
  flex-shrink: 0;
  position: relative;
  transition: all 0.15s;
}

.dp-todo input:checked + .todo-check {
  background: #16a34a;
  border-color: #16a34a;
}

.dp-todo input:checked + .todo-check::after {
  content: '';
  position: absolute;
  left: 6px;
  top: 2px;
  width: 5px;
  height: 10px;
  border: solid #fff;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.todo-slot {
  font-size: 0.7rem;
  font-weight: 600;
  color: #4a5d52;
  padding: 3px 8px;
  border-radius: 6px;
  background: #f0f4f2;
  white-space: nowrap;
  flex-shrink: 0;
}

.dp-todo.slot-pm .todo-slot { background: #e0f2fe; color: #0369a1; }
.dp-todo.slot-ev .todo-slot { background: #fef3c7; color: #b45309; }
.dp-todo.slot-extra .todo-slot { background: #fce7f3; color: #be185d; }

.todo-content {
  font-size: 0.88rem;
  color: #1a2e23;
  line-height: 1.5;
}

.dp-todo.checked .todo-content {
  color: #a3b5ab;
  text-decoration: line-through;
}

/* 响应式 */
@media (max-width: 768px) {
  .dp-head-row { flex-direction: column; }

  .dp-stage-bars { grid-template-columns: 1fr; }

  .dp-filters { flex-direction: column; align-items: stretch; }

  .dp-search { max-width: none; }

  .day-badge { min-width: 36px; height: 36px; font-size: 0.8rem; }

  .day-output { white-space: normal; }

  .todo-content { font-size: 0.82rem; }
}
</style>
