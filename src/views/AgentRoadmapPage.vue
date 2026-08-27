<template>
  <div class="tp-page">
    <!-- 页头 -->
    <header class="tp-head">
      <div class="tp-badge">🤖 AI Agent 工作流 · 3 个月自学冲刺版</div>
      <h1 class="tp-title">从 Vue 前端到 AI 应用工程师</h1>
      <p class="tp-desc">
        脱产 3 个月、零成本、纯本地运行。基于 5 年 Vue 经验补齐 Agent / RAG / 工作流画布短板，产出可手机远程演示的完整作品集，直接面试工业云、AI 售前、AI 前端岗位。
      </p>
      <div class="tp-stats">
        <span class="tp-stat"><b>90</b> 天</span>
        <span class="tp-stat-sep"></span>
        <span class="tp-stat"><b>3</b> 阶段</span>
        <span class="tp-stat-sep"></span>
        <span class="tp-stat"><b>3</b> 完整作品</span>
        <span class="tp-stat-sep"></span>
        <span class="tp-stat"><b>5</b> 面试话术</span>
        <span class="tp-stat-sep"></span>
        <span class="tp-stat"><b>0</b> 元成本</span>
      </div>

      <!-- 场景定位标签 -->
      <div class="tp-tags-row">
        <span class="tp-tag tag-person">👤 5 年 Vue 前端</span>
        <span class="tp-tag tag-state">💼 脱产 3 个月</span>
        <span class="tp-tag tag-hw">💻 32G 内存笔记本</span>
        <span class="tp-tag tag-cost">🆓 零成本 / 本地运行</span>
        <span class="tp-tag tag-target">🎯 工业云 · AI 售前 · AI 前端</span>
      </div>
    </header>

    <!-- 0 天准备 -->
    <section class="tp-section">
      <h2 class="tp-h2">⏱ 第 0 天：环境与资料准备（半天搞定）</h2>
      <p class="tp-sub">所有依赖一次性配齐，后续 90 天零障碍推进</p>

      <div class="tp-prep">
        <div class="prep-card">
          <div class="prep-icon">🛠</div>
          <h3>环境搭建</h3>
          <ul>
            <li>VS Code + Vue 3 + Vite 已就绪</li>
            <li>安装 <a href="https://ollama.com/download" target="_blank" rel="noopener">Ollama</a>（本地大模型运行时）</li>
            <li>下载轻量模型：<code>ollama pull qwen2.5:7b</code>、<code>ollama pull deepseek-r1:7b</code></li>
          </ul>
        </div>

        <div class="prep-card">
          <div class="prep-icon">📡</div>
          <h3>内网穿透</h3>
          <ul>
            <li>安装 <a href="https://tailscale.com/download" target="_blank" rel="noopener">Tailscale</a>（免费版够用）</li>
            <li>手机 / 电脑同账号登录</li>
            <li>面试时直接用 Tailscale IP 远程演示</li>
          </ul>
        </div>

        <div class="prep-card">
          <div class="prep-icon">📚</div>
          <h3>资料来源（全免费）</h3>
          <ul>
            <li>B 站搜索"Vue3 Ollama / LangChain.js / Dify 部署"</li>
            <li><a href="https://js.langchain.com/docs/" target="_blank" rel="noopener">LangChain.js 官方文档</a> + <a href="https://docs.dify.ai/" target="_blank" rel="noopener">Dify 文档</a></li>
            <li>掘金 / 思否搜"AI Agent 工作流 + Vue 实战"</li>
          </ul>
        </div>

        <div class="prep-card">
          <div class="prep-icon">📦</div>
          <h3>产出标准</h3>
          <ul>
            <li>每个阶段都写代码、存项目</li>
            <li>每月整理一次作品集</li>
            <li>所有 Demo 都接 Tailscale 远程可访问</li>
          </ul>
        </div>
      </div>
    </section>

    <!-- 三阶段路线 -->
    <section class="tp-section">
      <h2 class="tp-h2">🗺 三阶段学习路线</h2>
      <p class="tp-sub">每个阶段有明确目标、核心知识点、可交付 Demo 与验收标准</p>

      <div class="tp-stages">
        <div class="tp-stage" style="--sc: #16a34a">
          <div class="stage-left">
            <div class="stage-num">1</div>
            <div class="stage-dur">Day 1-30</div>
          </div>
          <div class="stage-body">
            <h3 class="stage-title">AI 基础 + 流式对话底座</h3>
            <p class="stage-goal"><span class="goal-tag">目标</span>补齐 LLM 调用、流式输出、Vue3 对话页面的基础能力。</p>
            <div class="stage-tags">
              <span class="topic-tag">Ollama 本地调用</span>
              <span class="topic-tag">SSE 流式输出</span>
              <span class="topic-tag">WebSocket 双向</span>
              <span class="topic-tag">Function Calling</span>
              <span class="topic-tag">上下文管理</span>
              <span class="topic-tag">移动端适配</span>
            </div>
            <div class="stage-output">
              <span class="output-label">作品 1</span>
              <span class="output-text">本地 AI 对话平台：流式打字机 + 工具调用 + 手机可访问</span>
            </div>
            <div class="stage-check">
              <span class="check-label">验收</span>
              <span class="check-text">电脑本地跑，手机连 WiFi 打开页面能正常对话 + 调用工具</span>
            </div>
          </div>
        </div>

        <div class="tp-stage" style="--sc: #0891b2">
          <div class="stage-left">
            <div class="stage-num">2</div>
            <div class="stage-dur">Day 31-60</div>
          </div>
          <div class="stage-body">
            <h3 class="stage-title">RAG 知识库 + Agent 基础</h3>
            <p class="stage-goal"><span class="goal-tag">目标</span>掌握 RAG 全链路 + Agent 四大模块，能搭简单智能体。</p>
            <div class="stage-tags">
              <span class="topic-tag">文档切片</span>
              <span class="topic-tag">向量嵌入</span>
              <span class="topic-tag">Chroma 向量库</span>
              <span class="topic-tag">重排 / 引用溯源</span>
              <span class="topic-tag">Agent 记忆·工具·反思</span>
              <span class="topic-tag">LangChain.js</span>
              <span class="topic-tag">Dify 本地部署</span>
            </div>
            <div class="stage-output">
              <span class="output-label">作品 2</span>
              <span class="output-text">工业文档问答系统：上传 PDF → AI 检索回答 + 引用来源</span>
            </div>
            <div class="stage-check">
              <span class="check-label">验收</span>
              <span class="check-text">上传手册可答专业问题，无文档时不会乱编</span>
            </div>
          </div>
        </div>

        <div class="tp-stage" style="--sc: #7c3aed">
          <div class="stage-left">
            <div class="stage-num">3</div>
            <div class="stage-dur">Day 61-90</div>
          </div>
          <div class="stage-body">
            <h3 class="stage-title">Agent 工作流可视化画布</h3>
            <p class="stage-goal"><span class="goal-tag">目标</span>集成前三阶段所有能力，做出可面试演示的完整作品集。</p>
            <div class="stage-tags">
              <span class="topic-tag">DAG / 分支 / 循环</span>
              <span class="topic-tag">LangGraph</span>
              <span class="topic-tag">Vue3 拖拽画布</span>
              <span class="topic-tag">参数配置面板</span>
              <span class="topic-tag">执行监控大屏</span>
              <span class="topic-tag">Tailscale 远程演示</span>
            </div>
            <div class="stage-output">
              <span class="output-label">作品 3</span>
              <span class="output-text">工业场景 AI Agent 工作流平台（简历核心项目）</span>
            </div>
            <div class="stage-check">
              <span class="check-label">验收</span>
              <span class="check-text">面试时电脑 + 手机远程完整演示整套工作流执行</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 12 周时间安排 -->
    <section class="tp-section tp-weekly-bg">
      <h2 class="tp-h2">📅 12 周时间安排</h2>
      <p class="tp-sub">按周推进，每周末做一次回顾与产出整理</p>

      <div class="tp-weekly">
        <div class="week-cell" style="--wc: #94a3b8">
          <div class="week-head"><span class="week-num">W0</span><span class="week-stage">第 0 天</span></div>
          <p class="week-task">环境 / Ollama / Tailscale / 资料准备</p>
        </div>
        <div class="week-cell" style="--wc: #16a34a">
          <div class="week-head"><span class="week-num">W1</span><span class="week-stage">阶段一</span></div>
          <p class="week-task">LLM 原理 + Ollama 调用 + Prompt 基础</p>
        </div>
        <div class="week-cell" style="--wc: #16a34a">
          <div class="week-head"><span class="week-num">W2</span><span class="week-stage">阶段一</span></div>
          <p class="week-task">SSE 流式输出 + WebSocket 双向通信</p>
        </div>
        <div class="week-cell" style="--wc: #16a34a">
          <div class="week-head"><span class="week-num">W3</span><span class="week-stage">阶段一</span></div>
          <p class="week-task">Function Calling + Vue3 对话页面开发</p>
        </div>
        <div class="week-cell" style="--wc: #16a34a">
          <div class="week-head"><span class="week-num">W4</span><span class="week-stage">阶段一</span></div>
          <p class="week-task">移动端适配 + 作品 1 整合上线</p>
        </div>
        <div class="week-cell" style="--wc: #0891b2">
          <div class="week-head"><span class="week-num">W5</span><span class="week-stage">阶段二</span></div>
          <p class="week-task">RAG 全流程 + Chroma 向量库</p>
        </div>
        <div class="week-cell" style="--wc: #0891b2">
          <div class="week-head"><span class="week-num">W6</span><span class="week-stage">阶段二</span></div>
          <p class="week-task">重排 / HyDE / 引用溯源进阶</p>
        </div>
        <div class="week-cell" style="--wc: #0891b2">
          <div class="week-head"><span class="week-num">W7</span><span class="week-stage">阶段二</span></div>
          <p class="week-task">Agent 四大模块 + LangChain.js</p>
        </div>
        <div class="week-cell" style="--wc: #0891b2">
          <div class="week-head"><span class="week-num">W8</span><span class="week-stage">阶段二</span></div>
          <p class="week-task">Dify 本地部署 + 作品 2 完成</p>
        </div>
        <div class="week-cell" style="--wc: #7c3aed">
          <div class="week-head"><span class="week-num">W9</span><span class="week-stage">阶段三</span></div>
          <p class="week-task">DAG / 分支 / 循环 / 失败重试</p>
        </div>
        <div class="week-cell" style="--wc: #7c3aed">
          <div class="week-head"><span class="week-num">W10</span><span class="week-stage">阶段三</span></div>
          <p class="week-task">LangGraph 复杂流 + 拖拽画布</p>
        </div>
        <div class="week-cell" style="--wc: #7c3aed">
          <div class="week-head"><span class="week-num">W11</span><span class="week-stage">阶段三</span></div>
          <p class="week-task">执行监控大屏 + Tailscale 远程</p>
        </div>
        <div class="week-cell" style="--wc: #7c3aed">
          <div class="week-head"><span class="week-num">W12</span><span class="week-stage">阶段三</span></div>
          <p class="week-task">终版作品集 + 简历 + 面试话术</p>
        </div>
      </div>
    </section>

    <!-- 每日作息 -->
    <section class="tp-section">
      <h2 class="tp-h2">🌅 每日固定作息（脱产版）</h2>
      <p class="tp-sub">避免摆烂：上午理论 / 下午实战 / 晚上复盘，3 个月完整闭环</p>

      <div class="tp-schedule">
        <div class="sched-row">
          <span class="sched-time">09:00 - 11:30</span>
          <span class="sched-label">📖 理论学习</span>
          <span class="sched-text">官方文档 / 实操视频 / 源码阅读</span>
        </div>
        <div class="sched-row">
          <span class="sched-time">14:00 - 17:30</span>
          <span class="sched-label">💻 手写代码</span>
          <span class="sched-text">当日 Demo / 作品 / Git 提交</span>
        </div>
        <div class="sched-row">
          <span class="sched-time">19:00 - 21:00</span>
          <span class="sched-label">🔄 复盘口述</span>
          <span class="sched-text">整理代码、写面试话术、录讲解 demo</span>
        </div>
        <div class="sched-row">
          <span class="sched-time">21:00 - 21:30</span>
          <span class="sched-label">🎯 简历素材</span>
          <span class="sched-text">记录亮点、量化产出、维护作品集 README</span>
        </div>
      </div>
    </section>

    <!-- 求职同步时间线 -->
    <section class="tp-section">
      <h2 class="tp-h2">💼 求职同步动作（边学边投）</h2>
      <p class="tp-sub">每月结束更新一次简历 + 投递重点调整，不浪费空档期</p>

      <div class="tp-hunt">
        <div class="hunt-card" style="--hc: #16a34a">
          <div class="hunt-head">
            <span class="hunt-month">第 1 个月末</span>
            <span class="hunt-target">普通 AI 前端岗</span>
          </div>
          <ul class="hunt-list">
            <li>简历补充：AI 对话 / SSE 流式开发经验</li>
            <li>作品 1：本地 AI 对话平台（流式 + 工具调用）</li>
            <li>投递：Boss 直聘搜"AI 前端 / 智能体前端"</li>
          </ul>
        </div>

        <div class="hunt-card" style="--hc: #0891b2">
          <div class="hunt-head">
            <span class="hunt-month">第 2 个月末</span>
            <span class="hunt-target">AI 应用 / SaaS 厂商</span>
          </div>
          <ul class="hunt-list">
            <li>简历新增：RAG 知识库项目</li>
            <li>作品 2：工业文档问答系统（Chroma + 引用）</li>
            <li>投递：AI 应用 / SaaS 厂商 / 文档智能</li>
          </ul>
        </div>

        <div class="hunt-card" style="--hc: #7c3aed">
          <div class="hunt-head">
            <span class="hunt-month">第 3 个月末</span>
            <span class="hunt-target">工业云 / 智慧制造 / Agent 平台</span>
          </div>
          <ul class="hunt-list">
            <li>简历主推：完整工作流画布 + 远程演示能力</li>
            <li>作品 3：工业 Agent 工作流平台（核心项目）</li>
            <li>投递：工业云 / 智慧制造 / 工业互联网 / Agent 平台</li>
          </ul>
        </div>
      </div>
    </section>

    <!-- 面试话术 -->
    <section class="tp-section">
      <h2 class="tp-h2">🎤 高频面试话术（每天 10 分钟背诵）</h2>
      <p class="tp-sub">点击展开标准回答要点，结合你已有的 5 年 Vue 经验组织口述</p>

      <div class="tp-qa">
        <details class="qa-item">
          <summary><span class="qa-q">Q1：SSE 和 WebSocket 的区别？业务落地场景？</span><span class="qa-arr">+</span></summary>
          <div class="qa-body">
            <p><strong>SSE（Server-Sent Events）</strong>：单向服务器推送，基于 HTTP 长连接，断点自动重连，浏览器原生 <code>EventSource</code> 支持。</p>
            <p><strong>WebSocket</strong>：全双工双向通信，独立协议（ws://），握手升级，适合频繁双向交互。</p>
            <p><strong>落地场景</strong>：</p>
            <ul>
              <li>✅ AI 流式对话 / 实时通知 / 日志推送 → 选 SSE（轻量、自动重连、断线恢复简单）</li>
              <li>✅ 多人协作画布 / 实时游戏 / IM 聊天 → 选 WebSocket（双向、需要双向同步状态）</li>
              <li>我的 Demo 中 AI 对话用 SSE，节点执行状态推送也用 SSE，拖拽画布多人协同才用 WebSocket。</li>
            </ul>
          </div>
        </details>

        <details class="qa-item">
          <summary><span class="qa-q">Q2：RAG 完整实现流程？如何解决幻觉问题？</span><span class="qa-arr">+</span></summary>
          <div class="qa-body">
            <p><strong>完整流程</strong>：文档上传 → 文本切片（Chunking）→ Embedding 向量化 → 入向量库 → 用户提问向量化 → 相似度检索 → 重排（Rerank）→ Top-K 拼进 Prompt → LLM 生成带引用的答案。</p>
            <p><strong>解决幻觉的 4 个手段</strong>：</p>
            <ul>
              <li>① <strong>限定回答范围</strong>：Prompt 里强制"只能基于以下文档回答，不知道就说不知道"</li>
              <li>② <strong>引用溯源</strong>：返回每条答案对应的文档片段 + 相似度分数</li>
              <li>③ <strong>重排 Rerank</strong>：用 BGE-reranker 等模型对初检结果二次排序</li>
              <li>④ <strong>事实校验</strong>：生成后让另一个 LLM 判断答案是否被上下文支持</li>
            </ul>
          </div>
        </details>

        <details class="qa-item">
          <summary><span class="qa-q">Q3：AI Agent 工作流是什么？LangGraph 怎么编排？</span><span class="qa-arr">+</span></summary>
          <div class="qa-body">
            <p><strong>Agent 工作流</strong>：把一个复杂任务拆成多个有依赖关系的节点（LLM 调用 / 工具调用 / 条件判断 / 人工确认），按 DAG（有向无环图）执行，支持分支、循环、失败重试。</p>
            <p><strong>LangGraph 核心</strong>：</p>
            <ul>
              <li>把工作流定义为状态机：<code>State</code>（共享数据）+ <code>Node</code>（处理函数）+ <code>Edge</code>（流转关系）</li>
              <li><code>add_conditional_edges</code> 实现条件分支，<code>add_node</code> 实现节点</li>
              <li>每个节点可独立重试、可单独接 LangSmith 观测</li>
              <li>工业巡检案例：知识库查询 → 报表生成 → 异常汇总 → 告警通知 → 人工介入</li>
            </ul>
          </div>
        </details>

        <details class="qa-item">
          <summary><span class="qa-q">Q4：本地部署项目如何实现手机跨局域网访问？</span><span class="qa-arr">+</span></summary>
          <div class="qa-body">
            <p><strong>方案对比</strong>：</p>
            <ul>
              <li>❌ 花生壳 / ngrok 付费版：免费版带宽极小、域名强制</li>
              <li>✅ <strong>Tailscale</strong>（推荐）：基于 WireGuard，点对点加密，免费 100 台设备，零配置</li>
              <li>✅ Cloudflare Tunnel：免费但需要绑域名，配置稍复杂</li>
              <li>✅ frp 自建：完全可控，但需要一台公网服务器</li>
            </ul>
            <p><strong>Tailscale 实操</strong>：</p>
            <ul>
              <li>电脑 + 手机同账号登录 → 分配 100.x.x.x 虚拟 IP</li>
              <li>Vite 起 dev server 时加 <code>--host</code>，监听 0.0.0.0</li>
              <li>手机浏览器输入 <code>http://电脑-tailscale-ip:端口</code> 即可访问</li>
              <li>面试场景：现场电脑跑 Demo，面试官手机扫码直接演示</li>
            </ul>
          </div>
        </details>

        <details class="qa-item">
          <summary><span class="qa-q">Q5：拖拽工作流画布前端实现方案？</span><span class="qa-arr">+</span></summary>
          <div class="qa-body">
            <p><strong>核心模块</strong>：</p>
            <ul>
              <li>① <strong>画布渲染</strong>：SVG 或 Canvas，我选 SVG（节点少时 DOM 可控，事件好处理）</li>
              <li>② <strong>节点拖拽</strong>：Pointer Events 监听 <code>pointerdown/move/up</code>，实时更新 <code>transform: translate3d</code></li>
              <li>③ <strong>连线绘制</strong>：贝塞尔曲线 <code>M x1 y1 C cx1 cy1, cx2 cy2, x2 y2</code>，箭头用 marker</li>
              <li>④ <strong>状态管理</strong>：Pinia 存 <code>{ nodes, edges }</code>，支持 undo/redo（命令模式）</li>
              <li>⑤ <strong>参数面板</strong>：节点选中后右侧抽屉，动态表单用 JSON Schema 渲染</li>
              <li>⑥ <strong>导入导出</strong>：JSON 序列化，支持存为模板</li>
            </ul>
            <p><strong>可选库</strong>：vue-flow（推荐）/ LogicFlow / 自研。优先自研，面试能讲清楚原理。</p>
          </div>
        </details>
      </div>
    </section>

    <!-- 避坑要点 -->
    <section class="tp-section">
      <h2 class="tp-h2">⚠️ 低成本避坑要点</h2>
      <p class="tp-sub">这 90 天最容易踩的坑，提前规避</p>

      <div class="tp-pitfalls">
        <div class="pitfall">
          <span class="pitfall-icon">💸</span>
          <div>
            <h4>不充值任何云大模型 API</h4>
            <p>全程本地 Ollama 跑 7B 模型，无月费、无 token 成本；面试时强调"零云依赖、可私有化部署"是加分项。</p>
          </div>
        </div>
        <div class="pitfall">
          <span class="pitfall-icon">🔗</span>
          <div>
            <h4>内网穿透只用 Tailscale 免费版</h4>
            <p>不要买付费隧道（frp / ngrok 付费版），Tailscale 免费版足够面试用，零运维。</p>
          </div>
        </div>
        <div class="pitfall">
          <span class="pitfall-icon">📚</span>
          <div>
            <h4>不买付费培训课</h4>
            <p>官方文档 + B 站免费教程 + 论文原文完全够；警惕"DLAI 课程打折"等套路。</p>
          </div>
        </div>
        <div class="pitfall">
          <span class="pitfall-icon">💻</span>
          <div>
            <h4>32G 内存电脑完全够用</h4>
            <p>qwen2.5:7b 量化版只需 4-6G 显存（CPU 也能跑，10-15 tokens/s），无需额外买显卡。</p>
          </div>
        </div>
        <div class="pitfall">
          <span class="pitfall-icon">📦</span>
          <div>
            <h4>作品必须可演示，不是堆代码</h4>
            <p>每个作品都要能跑、能截图、能手机访问、能讲 3 分钟；面试官看 demo 不看 GitHub star。</p>
          </div>
        </div>
        <div class="pitfall">
          <span class="pitfall-icon">🎯</span>
          <div>
            <h4>不要追新模型 / 新框架</h4>
            <p>3 个月只够深耕一套技术栈（Vue3 + LangChain.js + LangGraph），切来切去只会样样松。</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 教程资源 -->
    <section class="tp-section">
      <h2 class="tp-h2">📚 精选教程资源</h2>
      <p class="tp-sub">按阶段分类，全免费可直接访问</p>

      <div class="tp-resources">
        <div class="res-cat">
          <h3 class="cat-title">第 0-30 天：LLM 基础 + 流式对话</h3>
          <div class="res-list">
            <a href="https://ollama.com/" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">Ollama 官方文档</span><span class="res-arrow">↗</span>
              <p class="res-desc">本地大模型运行时，pull / run / API 全覆盖</p>
            </a>
            <a href="https://developer.mozilla.org/en-US/docs/Web/API/EventSource" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">MDN - EventSource</span><span class="res-arrow">↗</span>
              <p class="res-desc">浏览器原生 SSE 客户端 API 权威文档</p>
            </a>
            <a href="https://platform.openai.com/docs/guides/function-calling" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">OpenAI Function Calling</span><span class="res-arrow">↗</span>
              <p class="res-desc">工具调用的标准结构与 JSON Schema 写法</p>
            </a>
            <a href="https://www.bilibili.com/video/BV1Bo4y1A7FU" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">B 站：Vue3 + Ollama 流式对话实战</span><span class="res-arrow">↗</span>
              <p class="res-desc">国内最系统的 Vue3 + AI 集成实战视频</p>
            </a>
          </div>
        </div>

        <div class="res-cat">
          <h3 class="cat-title">第 31-60 天：RAG + Agent 基础</h3>
          <div class="res-list">
            <a href="https://js.langchain.com/docs/" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">LangChain.js 官方文档</span><span class="res-arrow">↗</span>
              <p class="res-desc">前端友好的 JS 版 LangChain，含 RAG / Agent / Memory 模块</p>
            </a>
            <a href="https://docs.trychroma.com/" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">Chroma 官方文档</span><span class="res-arrow">↗</span>
              <p class="res-desc">轻量级本地向量库，pip 一行安装，无数据库依赖</p>
            </a>
            <a href="https://docs.dify.ai/" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">Dify 官方文档</span><span class="res-arrow">↗</span>
              <p class="res-desc">可视化 AI 应用平台，本地 Docker 部署，含知识库 / Agent 编排</p>
            </a>
            <a href="https://lilianweng.github.io/posts/2023-06-23-agent/" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">Lil'Log - LLM Powered Autonomous Agents</span><span class="res-arrow">↗</span>
              <p class="res-desc">Agent 四大模块综述，记忆 / 规划 / 工具 / 反思权威解读</p>
            </a>
          </div>
        </div>

        <div class="res-cat">
          <h3 class="cat-title">第 61-90 天：工作流画布 + LangGraph</h3>
          <div class="res-list">
            <a href="https://langchain-ai.github.io/langgraphjs/" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">LangGraph.js 官方文档</span><span class="res-arrow">↗</span>
              <p class="res-desc">状态机式工作流编排，前端 JS 版本</p>
            </a>
            <a href="https://vueflow.dev/" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">Vue Flow</span><span class="res-arrow">↗</span>
              <p class="res-desc">Vue3 节点编辑器库，可二次封装为拖拽画布</p>
            </a>
            <a href="https://www.anthropic.com/engineering/building-effective-agents" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">Anthropic - Building Effective Agents</span><span class="res-arrow">↗</span>
              <p class="res-desc">从简单 Prompt 到多 Agent 的工程化指南</p>
            </a>
            <a href="https://tailscale.com/kb/" target="_blank" rel="noopener" class="res-item">
              <span class="res-name">Tailscale 文档</span><span class="res-arrow">↗</span>
              <p class="res-desc">零配置内网穿透，手机远程访问 dev server</p>
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
  --bg: #f3f8f5;
  max-width: 920px;
}

/* ============ 页头 ============ */
.tp-head {
  margin-bottom: 36px;
  padding-bottom: 24px;
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
  margin-bottom: 12px;
  font-weight: 600;
}

.tp-title {
  font-size: 1.85rem;
  font-weight: 800;
  color: var(--text);
  margin-bottom: 10px;
  line-height: 1.25;
  letter-spacing: -0.01em;
}

.tp-desc {
  color: var(--text-2);
  font-size: 0.9rem;
  line-height: 1.7;
  max-width: 720px;
  margin-bottom: 16px;
}

.tp-stats {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
  margin-bottom: 16px;
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

/* 场景定位标签 */
.tp-tags-row {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tp-tag {
  padding: 4px 10px;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 6px;
  font-size: 0.74rem;
  color: var(--text-2);
  font-weight: 500;
}

.tp-tag.tag-person { background: #eff6ff; border-color: #bfdbfe; color: #1d4ed8; }
.tp-tag.tag-state  { background: #fef3c7; border-color: #fde68a; color: #b45309; }
.tp-tag.tag-hw     { background: #f3e8ff; border-color: #ddd6fe; color: #6d28d9; }
.tp-tag.tag-cost   { background: var(--green-bg); border-color: var(--green-l); color: var(--green-d); }
.tp-tag.tag-target { background: #fce7f3; border-color: #fbcfe8; color: #be185d; }

/* ============ 区块 ============ */
.tp-section {
  margin-bottom: 40px;
}

.tp-h2 {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 6px;
}

.tp-sub {
  font-size: 0.83rem;
  color: var(--text-3);
  margin-bottom: 18px;
}

/* ============ 0 天准备 ============ */
.tp-prep {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
}

.prep-card {
  padding: 16px;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 10px;
}

.prep-icon {
  font-size: 1.6rem;
  margin-bottom: 6px;
}

.prep-card h3 {
  font-size: 0.92rem;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 10px;
}

.prep-card ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.prep-card li {
  font-size: 0.78rem;
  color: var(--text-2);
  line-height: 1.7;
  padding-left: 14px;
  position: relative;
}

.prep-card li::before {
  content: '·';
  position: absolute;
  left: 4px;
  color: var(--green);
  font-weight: 700;
}

.prep-card code {
  background: #f3f4f6;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 0.72rem;
  color: #1f2937;
  font-family: 'SF Mono', Consolas, monospace;
}

.prep-card a {
  color: var(--green-d);
  text-decoration: none;
  border-bottom: 1px solid transparent;
}

.prep-card a:hover {
  border-bottom-color: var(--green-d);
}

/* ============ 阶段卡片 ============ */
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
  min-width: 64px;
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

.stage-body {
  flex: 1;
}

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
  margin-bottom: 8px;
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

.stage-check {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding: 8px 10px;
  background: var(--green-bg);
  border-radius: 6px;
}

.check-label {
  font-size: 0.66rem;
  font-weight: 700;
  color: var(--green-d);
}

.check-text {
  font-size: 0.78rem;
  color: var(--text-2);
  line-height: 1.5;
}

/* ============ 周计划 ============ */
.tp-weekly-bg {
  background: var(--green-bg);
  margin: 0 -40px -40px;
  padding: 28px 40px 40px;
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

/* ============ 每日作息 ============ */
.tp-schedule {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 10px;
  overflow: hidden;
}

.sched-row {
  display: grid;
  grid-template-columns: 140px 140px 1fr;
  gap: 16px;
  padding: 14px 18px;
  border-bottom: 1px solid var(--border);
  align-items: center;
}

.sched-row:last-child {
  border-bottom: none;
}

.sched-time {
  font-family: 'SF Mono', Consolas, monospace;
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--green-d);
}

.sched-label {
  font-size: 0.84rem;
  font-weight: 700;
  color: var(--text);
}

.sched-text {
  font-size: 0.82rem;
  color: var(--text-2);
  line-height: 1.5;
}

/* ============ 求职同步时间线 ============ */
.tp-hunt {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 12px;
}

.hunt-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-top: 3px solid var(--hc);
  border-radius: 10px;
  padding: 16px;
}

.hunt-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px dashed var(--border);
}

.hunt-month {
  font-size: 0.84rem;
  font-weight: 800;
  color: var(--hc);
}

.hunt-target {
  font-size: 0.72rem;
  color: var(--text-3);
  background: #f7faf8;
  padding: 2px 8px;
  border-radius: 4px;
}

.hunt-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.hunt-list li {
  font-size: 0.78rem;
  color: var(--text-2);
  line-height: 1.7;
  padding-left: 14px;
  position: relative;
}

.hunt-list li::before {
  content: '✓';
  position: absolute;
  left: 0;
  color: var(--hc);
  font-weight: 700;
}

/* ============ 面试话术 ============ */
.tp-qa {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.qa-item {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 10px;
  transition: box-shadow 0.2s ease;
}

.qa-item[open] {
  border-color: var(--green);
  box-shadow: 0 2px 10px rgba(22, 163, 74, 0.06);
}

.qa-item summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 18px;
  cursor: pointer;
  list-style: none;
  user-select: none;
}

.qa-item summary::-webkit-details-marker {
  display: none;
}

.qa-q {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text);
}

.qa-arr {
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--green);
  transition: transform 0.2s ease;
}

.qa-item[open] .qa-arr {
  transform: rotate(45deg);
}

.qa-body {
  padding: 0 18px 16px;
  border-top: 1px solid var(--border);
}

.qa-body p {
  font-size: 0.82rem;
  color: var(--text-2);
  line-height: 1.7;
  margin: 10px 0;
}

.qa-body ul {
  list-style: none;
  padding: 0;
  margin: 8px 0;
}

.qa-body li {
  font-size: 0.82rem;
  color: var(--text-2);
  line-height: 1.7;
  padding-left: 18px;
  position: relative;
  margin-bottom: 4px;
}

.qa-body li::before {
  content: '•';
  position: absolute;
  left: 6px;
  color: var(--green);
  font-weight: 700;
}

.qa-body strong {
  color: var(--text);
  font-weight: 700;
}

.qa-body code {
  background: #f3f4f6;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 0.74rem;
  color: #1f2937;
  font-family: 'SF Mono', Consolas, monospace;
}

/* ============ 避坑要点 ============ */
.tp-pitfalls {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 10px;
}

.pitfall {
  display: flex;
  gap: 12px;
  padding: 14px 16px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 8px;
}

.pitfall-icon {
  font-size: 1.4rem;
  flex-shrink: 0;
}

.pitfall h4 {
  font-size: 0.85rem;
  font-weight: 700;
  color: #92400e;
  margin-bottom: 4px;
}

.pitfall p {
  font-size: 0.76rem;
  color: #78350f;
  line-height: 1.6;
}

/* ============ 教程资源 ============ */
.tp-resources {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 14px;
}

.res-cat {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 14px;
}

.cat-title {
  font-size: 0.85rem;
  font-weight: 700;
  margin-bottom: 10px;
  color: var(--green-d);
  padding-bottom: 6px;
  border-bottom: 1px solid var(--border);
}

.res-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.res-item {
  display: block;
  padding: 8px 10px;
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
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text);
}

.res-arrow {
  color: var(--green);
  font-size: 0.78rem;
  float: right;
}

.res-desc {
  font-size: 0.7rem;
  color: var(--text-3);
  line-height: 1.5;
  margin-top: 3px;
  margin-right: 16px;
}

/* ============ 响应式 ============ */
@media (max-width: 768px) {
  .tp-title { font-size: 1.35rem; }
  .tp-stats { gap: 10px; }
  .tp-stage { flex-direction: column; gap: 10px; padding: 16px; }
  .stage-left { display: flex; align-items: center; gap: 10px; text-align: left; min-width: 0; }
  .stage-num { margin-bottom: 0; width: 36px; height: 36px; font-size: 1rem; }
  .tp-weekly-bg { margin: 0 -16px -40px; padding: 20px 16px 40px; }
  .tp-weekly { grid-template-columns: 1fr; }
  .tp-resources { grid-template-columns: 1fr; }
  .sched-row { grid-template-columns: 1fr; gap: 4px; padding: 12px 14px; }
  .tp-prep { grid-template-columns: 1fr; }
  .tp-hunt { grid-template-columns: 1fr; }
  .tp-pitfalls { grid-template-columns: 1fr; }
}
</style>