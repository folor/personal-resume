<template>
  <div class="tp-page">
    <!-- 页头 -->
    <header class="tp-head">
      <div class="tp-badge">🔍 GEO · 生成式引擎优化路线</div>
      <h1 class="tp-title">让 AI 搜索主动引用你</h1>
      <p class="tp-desc">
        GEO（Generative Engine Optimization）是针对 ChatGPT、Perplexity、Google AI Overviews、Gemini 等生成式搜索引擎的内容优化方法。区别于传统 SEO 排名争夺，GEO 的核心是让 AI 在合成答案时<strong>引用你的内容</strong>。本路线含认知重构、内容/技术/监测分阶段路线、9 周时间安排与权威教程。
      </p>
      <div class="tp-stats">
        <span class="tp-stat"><b>5</b> 阶段</span>
        <span class="tp-stat-sep"></span>
        <span class="tp-stat"><b>9</b> 周核心</span>
        <span class="tp-stat-sep"></span>
        <span class="tp-stat"><b>20+</b> 教程资源</span>
        <span class="tp-stat-sep"></span>
        <span class="tp-stat"><b>中</b> 难度</span>
      </div>
    </header>

    <!-- SEO vs GEO 对比 -->
    <section class="tp-section">
      <h2 class="tp-h2">先搞懂：SEO vs GEO</h2>
      <p class="tp-sub">用户行为与优化逻辑的根本转变</p>

      <div class="tp-compare">
        <div class="compare-col seo-col">
          <div class="compare-head">
            <span class="compare-icon">🔗</span>
            <div>
              <h3>传统 SEO</h3>
              <p>争夺 10 条蓝链接排名</p>
            </div>
          </div>
          <ul class="compare-list">
            <li><span class="cl-label">目标</span>排名进前 3，抢占点击</li>
            <li><span class="cl-label">机制</span>关键词匹配 + 反向链接权重</li>
            <li><span class="cl-label">产出</span>10 条链接列表，用户逐个点</li>
            <li><span class="cl-label">指标</span>排名位次、自然流量、CTR</li>
            <li><span class="cl-label">内容</span>关键词密度、内链、Title 优化</li>
          </ul>
        </div>

        <div class="compare-col geo-col">
          <div class="compare-head">
            <span class="compare-icon">✨</span>
            <div>
              <h3>GEO 生成式引擎优化</h3>
              <p>让 AI 答案里出现你的引用</p>
            </div>
          </div>
          <ul class="compare-list">
            <li><span class="cl-label">目标</span>被 AI 合成答案时引用/提及</li>
            <li><span class="cl-label">机制</span>RAG 检索 + LLM 合成 + 引用选择</li>
            <li><span class="cl-label">产出</span>一个 AI 答案 + 若干引用来源</li>
            <li><span class="cl-label">指标</span>引用频率、品牌提及、答案准确度</li>
            <li><span class="cl-label">内容</span>事实密度、结构化、可引用性、E-E-A-T</li>
          </ul>
        </div>
      </div>

      <div class="compare-note">
        <strong>关键转变：</strong>用户不再"点击搜索"，而是"提问→得答案"。你的内容要么被 AI 引用（曝光），要么不存在。GEO 不是替代 SEO，而是在 SEO 基础上增加一层"AI 可读性"。
      </div>
    </section>

    <!-- 阶段路线 -->
    <section class="tp-section">
      <h2 class="tp-h2">分阶段学习路线</h2>
      <p class="tp-sub">每个阶段有明确目标与可交付产出</p>

      <div class="tp-stages">
        <div
          v-for="s in stages"
          :key="s.id"
          class="tp-stage"
          :style="{ '--sc': s.color }"
        >
          <div class="stage-left">
            <div class="stage-num">{{ s.id }}</div>
            <div class="stage-dur">{{ s.duration }}</div>
          </div>
          <div class="stage-body">
            <h3 class="stage-title">{{ s.title }}</h3>
            <p class="stage-goal"><span class="goal-tag">目标</span>{{ s.goal }}</p>
            <div class="stage-tags">
              <span v-for="t in s.topics" :key="t" class="topic-tag">{{ t }}</span>
            </div>
            <div class="stage-output">
              <span class="output-label">产出</span>
              <span class="output-text">{{ s.output }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 9 周计划 -->
    <section class="tp-section tp-weekly-bg">
      <h2 class="tp-h2">9 周时间安排</h2>
      <p class="tp-sub">建议配合一个真实站点做实验场</p>

      <div class="tp-weekly">
        <div
          v-for="w in weekly"
          :key="w.week"
          class="week-cell"
          :style="{ '--wc': w.color }"
        >
          <div class="week-head">
            <span class="week-num">W{{ w.week }}</span>
            <span class="week-stage">{{ w.stage }}</span>
          </div>
          <p class="week-task">{{ w.task }}</p>
        </div>
      </div>
    </section>

    <!-- 教程资源 -->
    <section class="tp-section">
      <h2 class="tp-h2">精选教程资源</h2>
      <p class="tp-sub">优先 Aleyda Solis 等权威信源与监测工具官方文档</p>

      <div class="tp-resources">
        <div v-for="cat in resources" :key="cat.category" class="res-cat">
          <h3 class="cat-title">{{ cat.category }}</h3>
          <div class="res-list">
            <a
              v-for="r in cat.items"
              :key="r.url"
              :href="r.url"
              target="_blank"
              rel="noopener noreferrer"
              class="res-item"
            >
              <span class="res-name">{{ r.name }}</span>
              <span class="res-arrow">↗</span>
              <p class="res-desc">{{ r.desc }}</p>
            </a>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

let observer: IntersectionObserver | null = null

onMounted(() => {
  observer = new IntersectionObserver(
    (entries) => entries.forEach((e) => { if (e.isIntersecting) e.target.classList.add('visible') }),
    { threshold: 0.08 }
  )
  document.querySelectorAll('.reveal').forEach((el) => observer!.observe(el))
})

onUnmounted(() => observer?.disconnect())

interface Stage {
  id: string
  title: string
  duration: string
  goal: string
  topics: string[]
  output: string
  color: string
}

const stages: Stage[] = [
  {
    id: '0',
    title: '认知重构',
    duration: '1 周',
    goal: '理解生成式引擎工作原理，从"排名思维"切换到"被引用思维"',
    topics: ['SEO vs GEO 本质', 'RAG 检索→LLM 合成→引用', '用户行为变化', '哪些内容会被 AI 引用', 'AI 爬虫(GPTBot/PerplexityBot)'],
    output: '一份 GEO 认知笔记 + 目标 AI 引擎清单',
    color: '#16a34a',
  },
  {
    id: '1',
    title: '内容优化',
    duration: '2-3 周',
    goal: '让内容具备"可被 AI 提取与引用"的结构与权威性',
    topics: ['信息密度', '清晰 H2/H3 结构', '列表与表格', 'E-E-A-T 权威信号', '长尾问答覆盖', '可引用性(让 AI 好抄)'],
    output: '5 篇 AI 友好的支柱内容(Pillar Content)',
    color: '#0891b2',
  },
  {
    id: '2',
    title: '技术可读性',
    duration: '2 周',
    goal: '让站点对 AI 爬虫与 RAG 检索系统都友好',
    topics: ['Schema.org 结构化数据', 'FAQ/HowTo/Article 标记', '语义 HTML', 'Core Web Vitals', '知识图谱对齐', 'robots.txt 与 AI 爬虫策略'],
    output: '站点结构化数据全覆盖 + AI 爬虫可抓取',
    color: '#7c3aed',
  },
  {
    id: '3',
    title: '监测与迭代',
    duration: '2-3 周',
    goal: '量化你在 AI 答案中的引用情况并持续优化',
    topics: ['GEO 监测工具', '引用频率追踪', '品牌提及监测', '提示词测试矩阵', 'A/B 内容实验'],
    output: '月度 GEO 报告 + 优化清单',
    color: '#db2777',
  },
  {
    id: '4',
    title: '进阶与多模态',
    duration: '持续',
    goal: '扩展到多模态与品牌 Entity 层面的优化',
    topics: ['图片 alt/视频字幕', '品牌 Entity 强化', 'LLM 训练语料覆盖', '本地化 GEO', '反向 GEO(隐私保护)'],
    output: '多模态 GEO 策略 + 品牌 Entity 知识图谱',
    color: '#d97706',
  },
]

interface Week {
  week: number
  stage: string
  task: string
  color: string
}

const weekly: Week[] = [
  { week: 1, stage: '阶段0', task: 'GEO 认知 + 拆解 3 个 AI 引擎答案', color: '#16a34a' },
  { week: 2, stage: '阶段1', task: '信息密度 + 结构化内容改写', color: '#0891b2' },
  { week: 3, stage: '阶段1', task: 'E-E-A-T + 长尾问答覆盖', color: '#0891b2' },
  { week: 4, stage: '阶段1', task: '5 篇支柱内容产出与发布', color: '#0891b2' },
  { week: 5, stage: '阶段2', task: 'Schema.org 结构化数据部署', color: '#7c3aed' },
  { week: 6, stage: '阶段2', task: '语义 HTML + AI 爬虫策略', color: '#7c3aed' },
  { week: 7, stage: '阶段3', task: '监测工具接入 + 基线测量', color: '#db2777' },
  { week: 8, stage: '阶段3', task: '提示词测试矩阵 + A/B 实验', color: '#db2777' },
  { week: 9, stage: '阶段3', task: '月度 GEO 报告 + 迭代清单', color: '#db2777' },
]

interface ResourceCat {
  category: string
  items: { name: string; url: string; desc: string }[]
}

const resources: ResourceCat[] = [
  {
    category: '权威博客',
    items: [
      { name: 'Aleyda Solis - GEO', url: 'https://www.aleyda.com/blog/', desc: '国际 SEO 权威，GEO 早期布道者，必读' },
      { name: 'Search Engine Land', url: 'https://searchengineland.com/', desc: '搜索行业第一新闻与深度分析' },
      { name: 'Search Engine Roundtable', url: 'https://www.seroundtable.com/', desc: '搜索算法与 AI 搜索变化追踪' },
      { name: 'Otterly.ai Blog', url: 'https://www.otterly.ai/blog/', desc: 'GEO 专项监测工具的实践博客' },
    ],
  },
  {
    category: '监测工具',
    items: [
      { name: 'Profound', url: 'https://www.profound.com/', desc: 'AI 搜索可见性监测，品牌提及追踪' },
      { name: 'Otterly.ai', url: 'https://www.otterly.ai/', desc: 'GEO 专用监测，提示词测试矩阵' },
      { name: 'Goodie', url: 'https://goodie.ai/', desc: 'AI 搜索表现追踪工具' },
      { name: 'Peec.ai', url: 'https://peec.ai/', desc: 'GEO 排名与引用监测平台' },
    ],
  },
  {
    category: '官方指南',
    items: [
      { name: 'Google Search Central', url: 'https://developers.google.com/search/docs', desc: 'Google 搜索官方文档，含 AI Overviews' },
      { name: 'Schema.org', url: 'https://schema.org/', desc: '结构化数据标准，AI 解析依据' },
      { name: 'Bing Webmaster', url: 'https://www.bing.com/webmasters', desc: 'Bing/Copilot 抓取与提交' },
      { name: 'Perplexity Publisher', url: 'https://docs.perplexity.ai/', desc: 'Perplexity 发布者指南' },
    ],
  },
  {
    category: '研究与课程',
    items: [
      { name: 'Originality.ai GEO Study', url: 'https://originality.ai/blog/generative-engine-optimization', desc: 'GEO 概念提出研究，引用率实证' },
      { name: 'Backlinko', url: 'https://backlinko.com/', desc: 'Brian Dean 的 SEO/GEO 实战研究' },
      { name: 'Ahrefs Blog', url: 'https://ahrefs.com/blog/', desc: '数据驱动的搜索营销研究' },
      { name: 'Moz Blog', url: 'https://moz.com/blog', desc: '老牌 SEO 社区，含 AI 搜索章节' },
    ],
  },
]
</script>

<style scoped>
.tp-page {
  --green: #16a34a;
  --green-d: #15803d;
  --green-l: #dcfce7;
  --green-bg: #f0fdf4;
  --text: #1a2e23;
  --text-2: #4a5d52;
  --text-3: #8a9b92;
  --border: #d1e7da;
  --card: #fff;
  max-width: 860px;
}

/* 页头 */
.tp-head {
  margin-bottom: 40px;
  padding-bottom: 28px;
  border-bottom: 2px solid var(--border);
}

.tp-badge {
  display: inline-block;
  padding: 4px 14px;
  background: var(--green-bg);
  border: 1px solid var(--green-l);
  border-radius: 100px;
  font-size: 0.78rem;
  color: var(--green-d);
  margin-bottom: 14px;
  font-weight: 600;
}

.tp-title {
  font-size: 1.9rem;
  font-weight: 800;
  color: var(--text);
  margin-bottom: 12px;
  line-height: 1.25;
  letter-spacing: -0.01em;
}

.tp-desc {
  color: var(--text-2);
  font-size: 0.92rem;
  line-height: 1.7;
  max-width: 700px;
  margin-bottom: 18px;
}

.tp-desc strong { color: var(--text); font-weight: 700; }

.tp-stats {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.tp-stat {
  font-size: 0.85rem;
  color: var(--text-2);
}

.tp-stat b {
  font-size: 1.1rem;
  color: var(--green);
  margin-right: 3px;
}

.tp-stat-sep {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--border);
}

/* 区块 */
.tp-section {
  margin-bottom: 44px;
}

.tp-h2 {
  font-size: 1.3rem;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 6px;
}

.tp-sub {
  font-size: 0.85rem;
  color: var(--text-3);
  margin-bottom: 20px;
}

/* SEO vs GEO 对比 */
.tp-compare {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
}

.compare-col {
  padding: 20px;
  border-radius: 10px;
  border: 1px solid var(--border);
  background: var(--card);
}

.seo-col { background: #f9fafb; }
.geo-col { background: var(--green-bg); border-color: var(--green-l); }

.compare-head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

.compare-icon { font-size: 1.6rem; }
.compare-head h3 { font-size: 1rem; font-weight: 700; color: var(--text); margin-bottom: 2px; }
.compare-head p { font-size: 0.76rem; color: var(--text-3); }

.compare-list { list-style: none; display: flex; flex-direction: column; gap: 9px; }
.compare-list li { font-size: 0.82rem; color: var(--text-2); line-height: 1.5; }
.cl-label {
  display: inline-block;
  min-width: 42px;
  font-size: 0.66rem;
  font-weight: 700;
  color: var(--green-d);
  margin-right: 8px;
}
.seo-col .cl-label { color: var(--text-3); }

.compare-note {
  padding: 14px 18px;
  background: var(--green-bg);
  border: 1px solid var(--green-l);
  border-radius: 8px;
  font-size: 0.82rem;
  color: var(--text-2);
  line-height: 1.7;
}
.compare-note strong { color: var(--green-d); }

/* 阶段卡片 */
.tp-stages {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tp-stage {
  display: flex;
  gap: 18px;
  padding: 20px;
  background: var(--card);
  border: 1px solid var(--border);
  border-left: 3px solid var(--sc);
  border-radius: 10px;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
}

.tp-stage:hover {
  box-shadow: 0 4px 16px rgba(22, 163, 74, 0.08);
  transform: translateX(3px);
}

.stage-left {
  flex-shrink: 0;
  text-align: center;
  min-width: 56px;
}

.stage-num {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: var(--sc);
  color: #fff;
  font-size: 1.25rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 6px;
}

.stage-dur {
  font-size: 0.7rem;
  color: var(--text-3);
  font-weight: 600;
}

.stage-body { flex: 1; }

.stage-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 8px;
}

.stage-goal {
  font-size: 0.86rem;
  color: var(--text-2);
  margin-bottom: 12px;
  line-height: 1.6;
}

.goal-tag {
  display: inline-block;
  padding: 1px 8px;
  background: var(--green-bg);
  color: var(--green-d);
  border-radius: 4px;
  font-size: 0.68rem;
  font-weight: 600;
  margin-right: 8px;
  vertical-align: middle;
}

.stage-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  margin-bottom: 12px;
}

.topic-tag {
  padding: 3px 10px;
  background: #f7faf8;
  border: 1px solid var(--border);
  border-radius: 100px;
  font-size: 0.72rem;
  color: var(--text-2);
}

.stage-output {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding-top: 10px;
  border-top: 1px dashed var(--border);
}

.output-label {
  font-size: 0.66rem;
  font-weight: 700;
  color: var(--sc);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.output-text {
  font-size: 0.82rem;
  color: var(--text);
  font-weight: 500;
}

/* 周计划 */
.tp-weekly-bg {
  background: var(--green-bg);
  margin: 0 -40px -44px;
  padding: 28px 40px 44px;
  border-radius: 12px 0 0 0;
}

.tp-weekly {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 10px;
}

.week-cell {
  padding: 12px;
  background: var(--card);
  border: 1px solid var(--border);
  border-top: 3px solid var(--wc);
  border-radius: 8px;
  transition: box-shadow 0.2s ease;
}

.week-cell:hover {
  box-shadow: 0 3px 12px rgba(0, 0, 0, 0.06);
}

.week-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.week-num {
  font-weight: 800;
  font-size: 0.88rem;
  color: var(--wc);
}

.week-stage {
  font-size: 0.64rem;
  color: var(--text-3);
  padding: 1px 6px;
  background: #f7faf8;
  border-radius: 4px;
}

.week-task {
  font-size: 0.78rem;
  color: var(--text-2);
  line-height: 1.5;
}

/* 教程资源 */
.tp-resources {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.res-cat {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 16px;
}

.cat-title {
  font-size: 0.88rem;
  font-weight: 700;
  margin-bottom: 12px;
  color: var(--green-d);
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
}

.res-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.res-item {
  display: block;
  padding: 10px 12px;
  border: 1px solid transparent;
  border-radius: 8px;
  text-decoration: none;
  transition: all 0.2s ease;
}

.res-item:hover {
  background: var(--green-bg);
  border-color: var(--green-l);
}

.res-name {
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text);
}

.res-arrow {
  color: var(--green);
  font-size: 0.8rem;
  float: right;
}

.res-desc {
  font-size: 0.72rem;
  color: var(--text-3);
  line-height: 1.5;
  margin-top: 3px;
  margin-right: 16px;
}

/* 响应式 */
@media (max-width: 768px) {
  .tp-title { font-size: 1.4rem; }
  .tp-stats { gap: 10px; }
  .tp-compare { grid-template-columns: 1fr; }
  .tp-stage { flex-direction: column; gap: 10px; padding: 16px; }
  .stage-left { display: flex; align-items: center; gap: 10px; text-align: left; min-width: 0; }
  .stage-num { margin-bottom: 0; width: 36px; height: 36px; font-size: 1rem; }
  .tp-weekly-bg { margin: 0 -16px -44px; padding: 20px 16px 44px; }
  .tp-weekly { grid-template-columns: 1fr; }
  .tp-resources { grid-template-columns: 1fr; }
}
</style>
