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

type TranslationResponse struct {
	OriginalName   string `json:"originalName"`
	TranslatedName string `json:"translatedName"`
	Language       string `json:"language"`
}

var factionsData FactionsData
var languageDicts map[string]map[string]string

func main() {
	fmt.Println("🚀 启动 Warframe 敌人数据查看器服务器...")
	fmt.Println("📁 请确保以下文件存在：")
	fmt.Println("   - data/factions_data.json")
	fmt.Println("   - data/enemy_data/ 目录 (包含所有敌人JSON文件)")
	fmt.Println("   - languages/ 目录 (包含语言字典文件)")
	fmt.Println("   - index.html (前端页面)")
	fmt.Println("\n🌐 服务器将在 http://localhost:5000 启动")

	// 初始化语言字典
	languageDicts = make(map[string]map[string]string)
	// 加载派系数据
	loadFactionsData()
	// 预加载常用语言字典
	loadLanguageDict("zh")
	loadLanguageDict("en")

	// 设置路由
	http.HandleFunc("/api/factions", handleFactions)
	http.HandleFunc("/api/enemy/", handleEnemy)
	http.HandleFunc("/api/status", handleStatus)
	http.HandleFunc("/api/translate/", handleTranslate)
	http.HandleFunc("/api/enemy_data", handleEnemyData)
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

func loadLanguageDict(lang string) {
	if languageDicts[lang] != nil {
		return // 已经加载过了
	}
	filePath := filepath.Join("languages", fmt.Sprintf("dict.%s.json", lang))
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("⚠️  未找到语言文件: %s\n", filePath)
		return
	}
	var langDict map[string]string
	err = json.Unmarshal(data, &langDict)
	if err != nil {
		fmt.Printf("❌ 解析语言文件失败 %s: %v\n", lang, err)
		return
	}
	languageDicts[lang] = langDict
	fmt.Printf("✅ 语言字典加载成功: %s\n", lang)
}

func handleFactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	lang := r.URL.Query().Get("lang")
	if lang == "" || lang == "zh" {
		json.NewEncoder(w).Encode(factionsData)
		return
	}
	// 翻译派系和敌人名
	if languageDicts[lang] == nil {
		loadLanguageDict(lang)
	}
	if languageDicts[lang] == nil || languageDicts["zh"] == nil {
		json.NewEncoder(w).Encode(factionsData)
		return
	}
	zhDict := languageDicts["zh"]
	targetDict := languageDicts[lang]
	result := make(map[string][]string)
	for faction, enemies := range factionsData {
		factionKey := ""
		for k, v := range zhDict {
			if v == faction {
				factionKey = k
				break
			}
		}
		translatedFaction := faction
		if factionKey != "" {
			if val, ok := targetDict[factionKey]; ok {
				translatedFaction = val
			}
		}
		var translatedEnemies []string
		for _, enemy := range enemies {
			enemyKey := ""
			for k, v := range zhDict {
				if v == enemy {
					enemyKey = k
					break
				}
			}
			translatedEnemy := enemy
			if enemyKey != "" {
				if val, ok := targetDict[enemyKey]; ok {
					translatedEnemy = val
				}
			}
			translatedEnemies = append(translatedEnemies, translatedEnemy)
		}
		result[translatedFaction] = translatedEnemies
	}
	json.NewEncoder(w).Encode(result)
}

func handleEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	path := strings.TrimPrefix(r.URL.Path, "/api/enemy/")
	if path == "" {
		http.Error(w, "敌人名称不能为空", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join("data", "enemy_data", path+".json")
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

func handleTranslate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	path := strings.TrimPrefix(r.URL.Path, "/api/translate/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.Error(w, "参数不足", http.StatusBadRequest)
		return
	}
	lang := parts[0]
	enemyName := strings.Join(parts[1:], "/")
	if languageDicts[lang] == nil {
		loadLanguageDict(lang)
		if languageDicts[lang] == nil {
			errorResp := ErrorResponse{Error: "语言不支持"}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResp)
			return
		}
	}
	translatedName := translateEnemyName(enemyName, lang)
	response := TranslationResponse{
		OriginalName:   enemyName,
		TranslatedName: translatedName,
		Language:       lang,
	}
	json.NewEncoder(w).Encode(response)
}

func translateEnemyName(enemyName, lang string) string {
	langDict := languageDicts[lang]
	if langDict == nil {
		return enemyName
	}
	// 尝试直接匹配
	if translated, exists := langDict[enemyName]; exists {
		return translated
	}
	// 尝试通过路径查找
	for path, translation := range langDict {
		if translation == enemyName {
			if targetDict := languageDicts[lang]; targetDict != nil {
				if targetTranslation, exists := targetDict[path]; exists {
					return targetTranslation
				}
			}
		}
	}
	return enemyName
}

// 新增：获取指定敌人指定等级数据（已翻译）
func handleEnemyData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	name := r.URL.Query().Get("name")
	level := r.URL.Query().Get("level")
	lang := r.URL.Query().Get("lang")
	if name == "" || level == "" || lang == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "参数不完整"})
		return
	}
	filePath := filepath.Join("data", "enemy_data", name+".json")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "敌人数据文件不存在"})
		return
	}
	var enemy EnemyData
	err = json.Unmarshal(data, &enemy)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "敌人数据解析失败"})
		return
	}
	if languageDicts[lang] == nil {
		loadLanguageDict(lang)
	}
	if languageDicts[lang] == nil || languageDicts["zh"] == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "语言包加载失败"})
		return
	}
	zhDict := languageDicts["zh"]
	targetDict := languageDicts[lang]
	translatedName := enemy.EnemyName
	for k, v := range zhDict {
		if v == enemy.EnemyName {
			if val, ok := targetDict[k]; ok {
				translatedName = val
			}
			break
		}
	}
	levelData, ok := enemy.LevelData[level]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "该等级数据不存在"})
		return
	}
	resp := map[string]interface{}{
		"enemyName": translatedName,
		"level": level,
		"data": levelData,
	}
	json.NewEncoder(w).Encode(resp)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	filePath := strings.TrimPrefix(path, "/")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}
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
