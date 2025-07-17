package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type FactionsData map[string][]string

type EnemyData struct {
	EnemyName string                 `json:"enemyName"`
	Faction   string                 `json:"faction"`
	LevelData map[string]interface{} `json:"levelData"`
}

type StatusResponse struct {
	Status        string `json:"status"`
	FactionsCount int    `json:"factions_count"`
	TotalEnemies  int    `json:"total_enemies"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var factionsData FactionsData

func main() {
	fmt.Println("🚀 启动 Warframe 敌人数据查看器服务器...")
	fmt.Println("📁 请确保以下文件存在：")
	fmt.Println("   - data/factions_data.json")
	fmt.Println("   - data/enemy_data/ 目录 (包含所有敌人JSON文件)")
	fmt.Println("   - index.html (前端页面)")
	fmt.Println("\n🌐 服务器将在 http://localhost:5000 启动")

	// 加载派系数据
	loadFactionsData()

	// 设置路由
	http.HandleFunc("/api/factions", handleFactions)
	http.HandleFunc("/api/enemy/", handleEnemy)
	http.HandleFunc("/api/status", handleStatus)
	http.HandleFunc("/", handleStatic)

	// 启动服务器
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func loadFactionsData() {
	data, err := ioutil.ReadFile("data/factions_data.json")
	if err != nil {
		fmt.Println("❌ 未找到 data/factions_data.json，请先运行 src/analyze_factions.py")
		return
	}

	err = json.Unmarshal(data, &factionsData)
	if err != nil {
		fmt.Println("❌ 解析派系数据失败:", err)
		return
	}

	fmt.Println("✅ 派系数据加载成功")
}

func handleFactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	json.NewEncoder(w).Encode(factionsData)
}

func handleEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从URL中提取敌人名称
	path := strings.TrimPrefix(r.URL.Path, "/api/enemy/")
	if path == "" {
		http.Error(w, "敌人名称不能为空", http.StatusBadRequest)
		return
	}

	// 构建文件路径
	filePath := filepath.Join("data", "enemy_data", path+".json")
	
	// 读取文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		errorResp := ErrorResponse{Error: "敌人数据文件不存在"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Write(data)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	totalEnemies := 0
	for _, enemies := range factionsData {
		totalEnemies += len(enemies)
	}

	status := StatusResponse{
		Status:        "running",
		FactionsCount: len(factionsData),
		TotalEnemies:  totalEnemies,
	}

	json.NewEncoder(w).Encode(status)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	// 默认页面
	if path == "/" {
		path = "/index.html"
	}

	// 移除开头的斜杠
	filePath := strings.TrimPrefix(path, "/")
	
	// 读取文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}

	// 设置正确的Content-Type
	switch {
	case strings.HasSuffix(filePath, ".html"):
		w.Header().Set("Content-Type", "text/html")
	case strings.HasSuffix(filePath, ".css"):
		w.Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(filePath, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
	}

	w.Write(data)
} 