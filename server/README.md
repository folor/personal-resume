# 90 天每日打卡 · 后端服务（Go + MySQL）· 多路线图版

支持 4 条并行路线图，每条独立打卡、独立进度统计。

## 路线图

| ID | 名称 | 天数 | 任务数 | 说明 |
|----|------|------|--------|------|
| 1 | 🤖 AI 全栈学习 | 91 | 428 | LLM + RAG + Agent + 工作流画布 |
| 2 | 🧠 AI 应用工程师转行 | 90 | 360 | 前端转 AI 应用：Python + LLM API + RAG + Agent + 求职 |
| 3 | 📱 小红书 + 视频号自媒体 | 90 | 360 | 图文起步 → 短视频练表达 → 变现 |
| 4 | 🌍 Alibaba.com 跨境电商 | 90 | 360 | 建材货源（吊顶/线条/角花）B2B 出海 |

路线图 2/3/4 每天含「extra」时段的表达能力练习（录音/口述/公开发布）。

## 目录

- `main.go` — HTTP 服务（net/http + database/sql + go-sql-driver/mysql）
- `go.mod` — Go 模块声明
- `seed/gen_seed.py` — 路线图 1 种子生成器
- `seed/gen_new_seeds.py` — 路线图 2/3/4 种子生成器
- `seed/schema.sql` — 数据库表结构
- `seed/seed_days.sql` — 路线图 1（AI 全栈学习，INSERT 不带 roadmap_id，默认 1）
- `seed/seed_ai_eng.sql` — 路线图 2（AI 应用工程师转行）
- `seed/seed_selfmedia.sql` — 路线图 3（小红书 + 视频号自媒体）
- `seed/seed_crossborder.sql` — 路线图 4（Alibaba.com 跨境电商）
- `start-server.bat` — 一键启动脚本

## 启动

### 1. MySQL 服务
```cmd
net start MySQL80
```
（首次安装见 README 末）

### 2. 启动 Go 服务
```cmd
start-server.bat
```
或手动：
```cmd
set GOPROXY=https://goproxy.cn,direct
go run main.go
```

注意：Go 装在 `D:\go\bin`，若提示找不到命令先 `set PATH=D:\go\bin;%PATH%`。
首次启动时自动建表并导入全部种子数据（`roadmaps` 表为空时触发）。

### 3. 接口

所有查询接口支持 `?roadmap=N` 参数（默认 1）。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET  | `/api/health` | 健康检查 |
| GET  | `/api/roadmaps` | 所有路线图元数据（名称/颜色/阶段名） |
| GET  | `/api/days?roadmap=&week=&stage=&day=` | 查询天（可按 week/stage/day 过滤） |
| GET  | `/api/stats?roadmap=` | 该路线图总进度 + 三阶段进度 |
| PATCH | `/api/todo/{id}` | 勾选/取消单条（body: `{"done": true/false}`） |
| POST | `/api/reset?roadmap=` | 重置指定路线图任务为未完成 |

默认端口 8080，前端通过 Vite 代理 `/api → http://127.0.0.1:8080`。

前端打卡页面路由：`http://localhost:5173/daily`（顶部 Tab 切换路线图）。

## 数据库

- 数据库：`roadmap`
- 用户：`root` / `roadmap123`（如需修改：设环境变量 `ROADMAP_DSN`）
- 表：
  - `roadmaps(id, name, short, color, icon, desc, stage1-3_name)` — 4 条路线图元数据
  - `plan_days(roadmap_id, day_no, week_no, stage, theme, output)` — 主键 (roadmap_id, day_no)，共 361 天
  - `todos(id, roadmap_id, day_no, slot, slot_name, content, done, done_at)` — 共 1508 条任务，外键关联 plan_days

### 重新导入种子数据

```cmd
D:\mysql-8.0.46-winx64\bin\mysql.exe -u root -proadmap123 -e "SET FOREIGN_KEY_CHECKS=0; DROP TABLE IF EXISTS todos; DROP TABLE IF EXISTS plan_days; DROP TABLE IF EXISTS roadmaps; SET FOREIGN_KEY_CHECKS=1;" roadmap
```

然后重启服务，自动重建并导入。

## MySQL 安装到 D 盘（首次）

```cmd
cd D:\mysql-setup
unzip mysql-8.0.46-winx64.zip -d D:\
```

写入 `D:\mysql-8.0.46-winx64\my.ini`：
```ini
[mysqld]
basedir=D:/mysql-8.0.46-winx64
datadir=D:/mysql-8.0.46-winx64/data
port=3306
character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci
[client]
port=3306
default-character-set=utf8mb4
```

```cmd
D:\mysql-8.0.46-winx64\bin\mysqld.exe --defaults-file=D:/mysql-8.0.46-winx64/my.ini --initialize-insecure
D:\mysql-8.0.46-winx64\bin\mysqld.exe --install MySQL80 --defaults-file=D:/mysql-8.0.46-winx64/my.ini
net start MySQL80
D:\mysql-8.0.46-winx64\bin\mysql.exe -uroot --skip-password -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'roadmap123'; CREATE DATABASE roadmap CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```
