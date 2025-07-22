package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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

var (
	factionsData FactionsData
	// 添加文件缓存
	fileCache     = make(map[string][]byte)
	fileCacheLock sync.RWMutex
)

// 添加一个新的处理函数，返回原始英文派系名称
func handleOriginalFactions(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	// 添加缓存控制
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// 读取原始派系数据
	data, err := os.ReadFile("data/factions_data.json")
	if err != nil {
		http.Error(w, "无法读取派系数据", http.StatusInternalServerError)
		return
	}

	// 解析JSON
	var factionsData map[string][]string
	err = json.Unmarshal(data, &factionsData)
	if err != nil {
		http.Error(w, "解析派系数据失败", http.StatusInternalServerError)
		return
	}

	// 创建新的数据结构，保持派系名称为英文
	originalFactionsData := make(map[string][]string)
	for faction, enemies := range factionsData {
		// 使用原始英文派系名称作为键
		originalFactionsData[faction] = enemies
	}

	// 返回原始派系数据
	json.NewEncoder(w).Encode(originalFactionsData)
}

func main() {
	fmt.Println("🚀 启动 Warframe 敌人数据查看器服务器...")
	fmt.Println("📁 请确保以下文件存在：")
	fmt.Println("   - data/factions_data.json")
	fmt.Println("   - data/enemy_data/ 目录 (包含所有敌人JSON文件)")
	fmt.Println("   - index.html (前端页面)")
	fmt.Println("\n🌐 服务器将在 http://localhost:5000 启动")

	// 加载派系数据
	loadFactionsData()

	// 预加载常用静态文件到缓存
	preloadStaticFiles()

	// 设置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/api/factions", handleOriginalFactions) // 使用新的处理函数
	mux.HandleFunc("/api/enemy/", handleEnemy)
	mux.HandleFunc("/api/status", handleStatus)
	mux.HandleFunc("/", handleStatic)

	// 创建带超时的服务器
	server := &http.Server{
		Addr:         ":5000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器
	log.Fatal(server.ListenAndServe())
}

func preloadStaticFiles() {
	// 预加载关键静态文件
	filesToPreload := []string{
		"index.html",
		"assets/app.js",
		"assets/styles.css",
	}

	for _, file := range filesToPreload {
		data, err := os.ReadFile(file)
		if err == nil {
			fileCacheLock.Lock()
			fileCache[file] = data
			fileCacheLock.Unlock()
			fmt.Printf("✅ 预加载文件: %s\n", file)
		} else {
			fmt.Printf("❌ 无法预加载文件: %s\n", file)
		}
	}

	// 预加载语言文件
	languageFiles, err := filepath.Glob("languages/dict.*.json")
	if err == nil {
		for _, file := range languageFiles {
			data, err := os.ReadFile(file)
			if err == nil {
				fileCacheLock.Lock()
				fileCache[file] = data
				fileCacheLock.Unlock()
				fmt.Printf("✅ 预加载语言文件: %s\n", file)
			}
		}
	}
}

func loadFactionsData() {
	data, err := os.ReadFile("data/factions_data.json")
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
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	// 添加缓存控制
	w.Header().Set("Cache-Control", "public, max-age=3600")

	json.NewEncoder(w).Encode(factionsData)
}

func handleEnemy(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/enemy/")
	if path == "" {
		http.Error(w, "敌人名称不能为空", http.StatusBadRequest)
		return
	}

	// 获取查询参数中的模式
	mode := r.URL.Query().Get("mode")
	dataDir := "enemy_data"

	// 如果指定为钢铁模式，使用钢铁数据目录
	if mode == "steel" {
		dataDir = "enemy_data_steel"
	}

	filePath := filepath.Join("data", dataDir, path+".json")

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		errorResp := ErrorResponse{Error: "敌人数据文件不存在"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// 添加缓存控制
	w.Header().Set("Cache-Control", "public, max-age=86400") // 24小时缓存

	data, err := os.ReadFile(filePath)
	if err != nil {
		errorResp := ErrorResponse{Error: "读取敌人数据失败"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Write(data)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

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
	if path == "/" {
		path = "/index.html"
	}

	filePath := strings.TrimPrefix(path, "/")

	// 设置适当的内容类型
	setContentType(w, filePath)

	// 添加缓存控制
	if !strings.Contains(filePath, "data/") {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 静态资源缓存1小时
	}

	// 检查缓存中是否有文件
	fileCacheLock.RLock()
	cachedData, found := fileCache[filePath]
	fileCacheLock.RUnlock()

	if found {
		w.Write(cachedData)
		return
	}

	// 如果缓存中没有，则从磁盘读取
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}

	// 缓存较小的文件（小于1MB）
	if len(data) < 1024*1024 && !strings.Contains(filePath, "data/enemy_data/") {
		fileCacheLock.Lock()
		fileCache[filePath] = data
		fileCacheLock.Unlock()
	}

	w.Write(data)
}

func setContentType(w http.ResponseWriter, filePath string) {
	switch {
	case strings.HasSuffix(filePath, ".html"):
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case strings.HasSuffix(filePath, ".css"):
		w.Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(filePath, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
	case strings.HasSuffix(filePath, ".json"):
		w.Header().Set("Content-Type", "application/json")
	case strings.HasSuffix(filePath, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(filePath, ".jpg"), strings.HasSuffix(filePath, ".jpeg"):
		w.Header().Set("Content-Type", "image/jpeg")
	case strings.HasSuffix(filePath, ".svg"):
		w.Header().Set("Content-Type", "image/svg+xml")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
