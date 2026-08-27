package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Todo 单条任务
type Todo struct {
	ID       int64  `json:"id"`
	DayNo    int    `json:"day_no"`
	Slot     string `json:"slot"`
	SlotName string `json:"slot_name"`
	Content  string `json:"content"`
	Done     bool   `json:"done"`
}

// Day 一天的计划（含任务列表）
type Day struct {
	DayNo    int    `json:"day_no"`
	WeekNo   int    `json:"week_no"`
	Stage    int    `json:"stage"`
	Theme    string `json:"theme"`
	Output   string `json:"output"`
	Todos    []Todo `json:"todos"`
	DoneCnt  int    `json:"done_cnt"`
	Total    int    `json:"total"`
}

// Roadmap 一个规划
type Roadmap struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Short      string `json:"short"`
	Color      string `json:"color"`
	Icon       string `json:"icon"`
	Desc       string `json:"desc"`
	StageNames []string `json:"stage_names"`
}

// Stats 总体进度
type Stats struct {
	TotalTodos   int `json:"total_todos"`
	DoneTodos    int `json:"done_todos"`
	TotalDays    int `json:"total_days"`
	FinishedDays int `json:"finished_days"`
	StageStats   []struct {
		Stage int    `json:"stage"`
		Name  string `json:"name"`
		Done  int    `json:"done"`
		Total int    `json:"total"`
	} `json:"stage_stats"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, PUT, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// roadmapIDFromQuery 解析 ?roadmap= 参数，默认 1
func roadmapIDFromQuery(r *http.Request) int {
	v := r.URL.Query().Get("roadmap")
	if v == "" {
		return 1
	}
	id, err := strconv.Atoi(v)
	if err != nil || id < 1 {
		return 1
	}
	return id
}

// GET /api/roadmaps  返回所有规划
func handleRoadmaps(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, short, color, icon, `desc`, stage1_name, stage2_name, stage3_name FROM roadmaps ORDER BY id")
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	defer rows.Close()
	out := []Roadmap{}
	for rows.Next() {
		var rm Roadmap
		var s1, s2, s3 string
		if err := rows.Scan(&rm.ID, &rm.Name, &rm.Short, &rm.Color, &rm.Icon, &rm.Desc, &s1, &s2, &s3); err != nil {
			continue
		}
		rm.StageNames = []string{s1, s2, s3}
		out = append(out, rm)
	}
	writeJSON(w, out)
}

// GET /api/days?roadmap=&week=&stage=&day=
func handleDays(w http.ResponseWriter, r *http.Request) {
	rid := roadmapIDFromQuery(r)
	q := r.URL.Query()
	where := []string{"roadmap_id = ?"}
	args := []any{rid}
	if v := q.Get("week"); v != "" {
		where = append(where, "week_no = ?")
		args = append(args, v)
	}
	if v := q.Get("stage"); v != "" {
		where = append(where, "stage = ?")
		args = append(args, v)
	}
	if v := q.Get("day"); v != "" {
		where = append(where, "day_no = ?")
		args = append(args, v)
	}
	cond := strings.Join(where, " AND ")

	rows, err := db.Query("SELECT day_no, week_no, stage, theme, output FROM plan_days WHERE "+cond+" ORDER BY day_no", args...)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	defer rows.Close()

	days := []Day{}
	for rows.Next() {
		var d Day
		if err := rows.Scan(&d.DayNo, &d.WeekNo, &d.Stage, &d.Theme, &d.Output); err != nil {
			writeErr(w, 500, err.Error())
			return
		}
		days = append(days, d)
	}
	rows.Close()

	for i := range days {
		loadTodos(rid, &days[i])
	}
	writeJSON(w, days)
}

func loadTodos(rid int, d *Day) {
	rows, err := db.Query("SELECT id, slot, slot_name, content, done FROM todos WHERE roadmap_id = ? AND day_no = ? ORDER BY FIELD(slot,'am','pm','ev','extra'), id", rid, d.DayNo)
	if err != nil {
		return
	}
	defer rows.Close()
	d.Total = 0
	d.DoneCnt = 0
	for rows.Next() {
		var t Todo
		var done int
		if err := rows.Scan(&t.ID, &t.Slot, &t.SlotName, &t.Content, &done); err != nil {
			continue
		}
		t.Done = done == 1
		t.DayNo = d.DayNo
		d.Todos = append(d.Todos, t)
		if t.Done {
			d.DoneCnt++
		}
		d.Total++
	}
}

// GET /api/stats?roadmap=
func handleStats(w http.ResponseWriter, r *http.Request) {
	rid := roadmapIDFromQuery(r)
	var s Stats

	// 取该 roadmap 的 stage 名
	stageNames := map[int]string{}
	rnameRows, _ := db.Query("SELECT stage1_name, stage2_name, stage3_name FROM roadmaps WHERE id = ?", rid)
	if rnameRows != nil {
		if rnameRows.Next() {
			var s1, s2, s3 string
			rnameRows.Scan(&s1, &s2, &s3)
			stageNames[1] = s1
			stageNames[2] = s2
			stageNames[3] = s3
		}
		rnameRows.Close()
	}

	db.QueryRow("SELECT COUNT(*), COALESCE(SUM(done),0) FROM todos WHERE roadmap_id = ?", rid).Scan(&s.TotalTodos, &s.DoneTodos)
	db.QueryRow("SELECT COUNT(*), COALESCE(SUM(done_cnt = total AND total > 0),0) FROM (SELECT day_no, SUM(done) done_cnt, COUNT(*) total FROM todos WHERE roadmap_id = ? GROUP BY day_no) t", rid).Scan(&s.TotalDays, &s.FinishedDays)

	rows, err := db.Query(`SELECT stage, SUM(t.done), COUNT(*) FROM todos t JOIN plan_days d ON t.roadmap_id = d.roadmap_id AND t.day_no = d.day_no WHERE t.roadmap_id = ? GROUP BY stage ORDER BY stage`, rid)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item struct {
				Stage int    `json:"stage"`
				Name  string `json:"name"`
				Done  int    `json:"done"`
				Total int    `json:"total"`
			}
			rows.Scan(&item.Stage, &item.Done, &item.Total)
			item.Name = stageNames[item.Stage]
			s.StageStats = append(s.StageStats, item)
		}
	}
	writeJSON(w, s)
}

// PATCH /api/todo/{id}  body: {"done": true}
func handleToggleTodo(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		writeErr(w, 400, "bad path")
		return
	}
	id, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil {
		writeErr(w, 400, "bad id")
		return
	}
	var body struct {
		Done *bool `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Done == nil {
		writeErr(w, 400, "body must be {\"done\": bool}")
		return
	}
	doneInt := 0
	if *body.Done {
		doneInt = 1
	}
	res, err := db.Exec("UPDATE todos SET done = ?, done_at = "+nullIf(*body.Done)+" WHERE id = ?", doneInt, id)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	n, _ := res.RowsAffected()
	writeJSON(w, map[string]any{"id": id, "done": *body.Done, "affected": n})
}

func nullIf(done bool) string {
	if done {
		return "NOW()"
	}
	return "NULL"
}

// POST /api/reset?roadmap=  清空指定规划的完成状态
func handleReset(w http.ResponseWriter, r *http.Request) {
	rid := roadmapIDFromQuery(r)
	db.Exec("UPDATE todos SET done = 0, done_at = NULL WHERE roadmap_id = ?", rid)
	writeJSON(w, map[string]string{"msg": "ok", "roadmap": strconv.Itoa(rid)})
}

func main() {
	dsn := os.Getenv("ROADMAP_DSN")
	if dsn == "" {
		dsn = "root:roadmap123@tcp(127.0.0.1:3306)/roadmap?charset=utf8mb4&parseTime=true&loc=Local"
	}
	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("mysql connect failed: ", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(10)

	initSchema()
	seedIfEmpty()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/roadmaps", handleRoadmaps)
	mux.HandleFunc("/api/days", handleDays)
	mux.HandleFunc("/api/stats", handleStats)
	mux.HandleFunc("/api/todo/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch || r.Method == http.MethodPut {
			handleToggleTodo(w, r)
			return
		}
		writeErr(w, 405, "method not allowed")
	})
	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeErr(w, 405, "method not allowed")
			return
		}
		handleReset(w, r)
	})
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})

	fmt.Println("roadmap-server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}

func initSchema() {
	// 检测旧表是否缺少 roadmap_id 列，缺则重建（个人项目，重建可接受）
	var hasCol int
	db.QueryRow("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='plan_days' AND COLUMN_NAME='roadmap_id'").Scan(&hasCol)
	if hasCol == 0 {
		log.Println("detecting old schema, rebuilding plan_days/todos for multi-roadmap...")
		db.Exec("DROP TABLE IF EXISTS todos")
		db.Exec("DROP TABLE IF EXISTS plan_days")
		db.Exec("DROP TABLE IF EXISTS roadmaps")
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS roadmaps (
			id INT PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			short VARCHAR(20) NOT NULL,
			color VARCHAR(20) NOT NULL,
			icon VARCHAR(10) NOT NULL,
			` + "`desc`" + ` VARCHAR(200) NOT NULL DEFAULT '',
			stage1_name VARCHAR(100) NOT NULL DEFAULT '',
			stage2_name VARCHAR(100) NOT NULL DEFAULT '',
			stage3_name VARCHAR(100) NOT NULL DEFAULT ''
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	`CREATE TABLE IF NOT EXISTS plan_days (
		roadmap_id INT NOT NULL DEFAULT 1,
		day_no INT NOT NULL,
		week_no INT NOT NULL,
		stage TINYINT NOT NULL,
		theme VARCHAR(200) NOT NULL,
		output VARCHAR(200) NOT NULL,
		PRIMARY KEY (roadmap_id, day_no)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	`CREATE TABLE IF NOT EXISTS todos (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		roadmap_id INT NOT NULL DEFAULT 1,
		day_no INT NOT NULL,
		slot VARCHAR(10) NOT NULL,
		slot_name VARCHAR(20) NOT NULL,
		content VARCHAR(500) NOT NULL,
		done TINYINT DEFAULT 0,
		done_at DATETIME NULL,
		KEY idx_roadmap_day (roadmap_id, day_no),
		CONSTRAINT fk_todos_roadmap_day FOREIGN KEY (roadmap_id, day_no) REFERENCES plan_days(roadmap_id, day_no) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			log.Fatal("init schema: ", err)
		}
	}
}

// seedRoadmaps 写入规划元数据
var seedRoadmaps = []Roadmap{
	{ID: 1, Name: "AI 全栈学习", Short: "AI学习", Color: "#16a34a", Icon: "🤖", Desc: "LLM + RAG + Agent + 工作流，90 天打卡", StageNames: []string{"阶段一 · AI 基础 + 流式对话", "阶段二 · RAG + Agent", "阶段三 · 工作流画布"}},
	{ID: 2, Name: "AI 应用工程师转行", Short: "AI工程师", Color: "#7c3aed", Icon: "🧠", Desc: "前端转 AI 应用方向：Python + LLM API + RAG + Agent + 前端集成 + 求职", StageNames: []string{"阶段一 · Python + LLM 基础", "阶段二 · RAG + Agent 实战", "阶段三 · 项目落地 + 求职"}},
	{ID: 3, Name: "小红书 + 视频号自媒体", Short: "自媒体", Color: "#db2777", Icon: "📱", Desc: "图文起步 → 短视频练表达 → 变现，含表达能力刻意练习", StageNames: []string{"阶段一 · 账号基建 + 内容打底", "阶段二 · 差异化 + 流量破圈", "阶段三 · 变现验证"}},
	{ID: 4, Name: "Alibaba.com 跨境电商", Short: "跨境电商", Color: "#0891b2", Icon: "🌍", Desc: "建材货源（吊顶/线条/角花）B2B 出海：开店 → 上架 → 询盘 → 信保订单", StageNames: []string{"阶段一 · 开店基建", "阶段二 · 上架引流", "阶段三 · 询盘转化"}},
}

func seedIfEmpty() {
	var cnt int
	db.QueryRow("SELECT COUNT(*) FROM roadmaps").Scan(&cnt)
	if cnt > 0 {
		fmt.Println("seed skipped, roadmaps =", cnt)
		return
	}

	// 1. 写入规划元数据
	for _, rm := range seedRoadmaps {
		db.Exec("INSERT INTO roadmaps (id, name, short, color, icon, `desc`, stage1_name, stage2_name, stage3_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			rm.ID, rm.Name, rm.Short, rm.Color, rm.Icon, rm.Desc, rm.StageNames[0], rm.StageNames[1], rm.StageNames[2])
	}

	// 2. 依次执行种子 SQL 文件（位于 seed/ 目录）
	seedDir := "seed"
	seedFiles := []string{
		"seed_days.sql",      // roadmap_id=1（原 AI 学习规划，INSERT 不带 roadmap_id，默认 1）
		"seed_ai_eng.sql",    // roadmap_id=2
		"seed_selfmedia.sql", // roadmap_id=3
		"seed_crossborder.sql", // roadmap_id=4
	}
	for _, f := range seedFiles {
		path := filepath.Join(seedDir, f)
		execSeedFile(path)
	}

	db.QueryRow("SELECT COUNT(*) FROM plan_days").Scan(&cnt)
	fmt.Println("seeded, plan_days =", cnt)
}

func execSeedFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("WARN: seed file missing, skip %s: %v", path, err)
		return
	}
	// 规范化换行符：Windows \r\n → \n，保证按 ";\n" 正确分割
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	// 关闭外键检查，避免 TRUNCATE/DELETE 被 FK 约束阻止
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	for _, stmt := range strings.Split(content, ";\n") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "USE ") || strings.HasPrefix(stmt, "SET ") {
			continue
		}
		// TRUNCATE 在有 FK 约束时可能仍然失败，替换为 DELETE FROM
		stmt = strings.ReplaceAll(stmt, "TRUNCATE TABLE", "DELETE FROM")
		if _, err := db.Exec(stmt); err != nil {
			preview := stmt
			if len(preview) > 80 {
				preview = preview[:80]
			}
			log.Printf("seed stmt failed in %s: %v | %s...", path, err, preview)
		}
	}
}
