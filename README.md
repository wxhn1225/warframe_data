# Warframe 敌人数据查看器

一个用于查看 Warframe 游戏中敌人数据的独立可执行应用程序。

## 功能特点

- 🎯 按派系分类查看敌人
- 📊 显示敌人详细属性数据
- 🌍 支持多语言界面
- 🔄 支持普通模式和钢铁模式切换
- ⚡ 快速加载和响应
- 📱 响应式设计，支持移动设备
- 🚀 独立可执行文件，无需安装任何依赖

## 文件结构

```
warframe_data/
├── See.exe            # 可执行文件
├── index.html         # 主页面
├── server.go          # 服务器源码
├── README.md          # 说明文档
├── src/               # 源代码
│   └── analyze_factions.py
├── data/              # 数据文件
│   ├── enemy_data/    # 普通模式敌人数据
│   ├── enemy_data_steel/ # 钢铁模式敌人数据
│   └── factions_data.json
├── languages/         # 多语言文件
│   ├── dict.zh.json
│   ├── dict.en.json
│   └── ...
└── assets/            # 资源文件
    ├── styles.css
    └── app.js
```

## 快速开始

1. 双击运行 `See.exe`
2. 在浏览器中访问 `http://localhost:5000`

## 使用方法

1. **启动应用**：双击 `See.exe`
2. **打开浏览器**：访问 `http://localhost:5000`
3. **选择模式**：在顶部选择普通模式或钢铁模式
4. **浏览派系**：点击任意派系查看敌人列表
5. **查看详情**：点击敌人名称查看详细数据
6. **选择等级**：在详情页面选择不同等级查看属性
7. **搜索敌人**：使用顶部搜索框搜索敌人

## 模式切换

- **普通模式**：显示常规游戏中的敌人数据
- **钢铁模式**：显示钢铁之路中的敌人数据（生命值、护甲、护盾和超宏数值不同）

## API 端点

- `GET /api/factions` - 获取所有派系数据
- `GET /api/enemy/<n>` - 获取敌人详细数据
- `GET /api/status` - 服务器状态

## 故障排除

### 常见问题

1. **端口被占用**
   - 关闭占用端口的程序
   - 或修改 `server.go` 中的端口号

2. **无法加载数据**
   - 确保 `data/enemy_data/` 目录存在
   - 确保 `data/enemy_data_steel/` 目录存在（钢铁模式）

3. **页面加载缓慢**
   - 检查网络连接
   - 确保数据文件完整

## 开发说明

### 修改代码

1. 修改 `server.go` 添加新功能
2. 修改 `assets/app.js` 更新前端逻辑
3. 修改 `assets/styles.css` 调整样式
4. 重新编译服务器

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