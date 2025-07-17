# Warframe 敌人数据查看器

一个用于查看 Warframe 游戏中敌人数据的独立可执行应用程序。

## 功能特点

- 🎯 按派系分类查看敌人
- 📊 显示敌人详细属性数据
- 🌍 支持多语言界面
- ⚡ 快速加载和响应
- 📱 响应式设计，支持移动设备
- 🚀 独立可执行文件，无需安装任何依赖

## 文件结构

```
warframe_data/
├── warframe-viewer.exe  # 可执行文件
├── index.html           # 主页面
├── styles.css           # 样式文件
├── app.js              # 脚本文件
├── start.bat           # 启动脚本
├── build.bat           # 构建脚本
├── README.md           # 说明文档
├── src/                # 源代码
│   ├── analyze_factions.py
│   ├── server.go
│   └── go.mod
├── data/               # 数据文件
│   ├── enemy_data/
│   └── factions_data.json
└── assets/             # 资源文件
    ├── styles.css
    └── app.js
```

## 快速开始

### 方法一：直接运行（推荐）

1. 双击运行 `warframe-viewer.exe`
2. 在浏览器中访问 `http://localhost:5000`

### 方法二：使用启动脚本

1. 双击运行 `start.bat`
2. 等待自动生成数据并启动服务器
3. 在浏览器中访问 `http://localhost:5000`

### 方法三：从源码构建

1. 确保已安装 Go 1.21+
2. 运行构建脚本：
   ```bash
   build.bat
   ```
3. 运行生成的可执行文件：
   ```bash
   warframe-viewer.exe
   ```

## 使用方法

1. **启动应用**：双击 `warframe-viewer.exe`
2. **打开浏览器**：访问 `http://localhost:5000`
3. **浏览派系**：点击任意派系查看敌人列表
4. **查看详情**：点击敌人名称查看详细数据
5. **选择等级**：在详情页面选择不同等级查看属性

## API 端点

- `GET /api/factions` - 获取所有派系数据
- `GET /api/enemy/<name>` - 获取敌人详细数据
- `GET /api/status` - 服务器状态

## 故障排除

### 常见问题

1. **端口被占用**
   - 关闭占用端口的程序
   - 或修改 `src/server.go` 中的端口号

2. **无法加载数据**
   - 确保 `data/enemy_data/` 目录存在
   - 运行 `start.bat` 重新生成数据

3. **页面加载缓慢**
   - 检查网络连接
   - 确保数据文件完整

## 开发说明

### 修改代码

1. 修改 `src/server.go` 添加新功能
2. 修改 `app.js` 更新前端逻辑
3. 修改 `styles.css` 调整样式
4. 运行 `build.bat` 重新构建

### 数据格式

敌人JSON文件应包含：
- `enemyName`: 敌人名称
- `faction`: 派系
- `levelData`: 等级数据对象

## 技术栈

- **后端**: Go 1.21+
- **前端**: HTML5 + CSS3 + JavaScript
- **数据格式**: JSON
- **打包**: Go build

## 许可证

本项目仅供学习和研究使用。 