# Warframe 敌人数据查看器

一个用于查看 Warframe 游戏中敌人数据的独立可执行应用程序。

## 功能特点

- 🎯 按派系分类查看敌人
- 📊 显示敌人详细属性数据
- 🌍 支持多语言界面
- ⚡ 快速加载和响应
- 📱 响应式设计，支持移动设备
- 🚀 独立可执行文件，无需安装任何依赖

## 快速开始

### 方法一：使用预编译版本（推荐）

1. 下载 `release` 目录中的文件
2. 双击运行 `warframe-viewer.exe`
3. 在浏览器中访问 `http://localhost:5000`

### 方法二：从源码构建

1. 确保已安装 Go 1.21+
2. 运行构建脚本：
   ```bash
   build.bat
   ```
3. 运行生成的可执行文件：
   ```bash
   release\warframe-viewer.exe
   ```

### 方法三：开发模式

1. 安装依赖：
   ```bash
   pip install -r requirements.txt
   ```

2. 生成派系数据：
   ```bash
   python analyze_factions.py
   ```

3. 启动Go服务器：
   ```bash
   go run server.go
   ```

4. 在浏览器中访问 `http://localhost:5000`

## 文件结构

```
warframe_data/
├── enemy_data/           # 敌人JSON数据文件
├── analyze_factions.py   # 派系分析脚本
├── server.go            # Go后端服务器
├── server.py            # Python后端服务器（备用）
├── index.html           # 前端页面
├── styles.css           # 样式文件
├── app.js              # 前端逻辑
├── requirements.txt     # Python依赖
├── go.mod              # Go模块文件
├── build.bat           # 构建脚本
├── start.bat           # Python启动脚本
└── README.md           # 说明文档
```

## 打包说明

### Go版本（推荐）

Go版本的优势：
- 编译成单个可执行文件
- 无需安装任何运行时
- 跨平台支持
- 性能优秀

构建命令：
```bash
go build -o warframe-viewer.exe server.go
```

### Python版本（备用）

使用PyInstaller打包：
```bash
pip install pyinstaller
pyinstaller --onefile --add-data "enemy_data;enemy_data" --add-data "index.html;." --add-data "styles.css;." --add-data "app.js;." server.py
```

## API 端点

- `GET /api/factions` - 获取所有派系数据
- `GET /api/enemy/<name>` - 获取敌人详细数据
- `GET /api/status` - 服务器状态

## 故障排除

### 常见问题

1. **端口被占用**
   - 修改 `server.go` 中的端口号
   - 或关闭占用端口的程序

2. **无法加载敌人数据**
   - 确保 `enemy_data/` 目录存在且包含JSON文件
   - 检查文件编码是否为UTF-8

3. **页面加载缓慢**
   - 确保已运行 `analyze_factions.py` 生成派系数据
   - 检查网络连接

### 日志查看

服务器运行时会显示详细的日志信息，包括：
- 数据加载状态
- API请求记录
- 错误信息

## 开发说明

### 添加新功能

1. 修改 `server.go` 添加新的API端点
2. 更新 `app.js` 添加前端逻辑
3. 根据需要修改 `styles.css` 调整样式

### 数据格式

敌人JSON文件应包含以下字段：
- `enemyName`: 敌人名称
- `faction`: 派系
- `levelData`: 等级数据对象

## 技术栈

- **后端**: Go 1.21+ (主要) / Python Flask (备用)
- **前端**: HTML5 + CSS3 + JavaScript
- **数据格式**: JSON
- **打包工具**: Go build / PyInstaller

## 许可证

本项目仅供学习和研究使用。 