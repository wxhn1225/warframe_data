@echo off
echo ========================================
echo 构建 Warframe 敌人数据查看器
echo ========================================
echo.

echo 正在检查Go环境...
go version >nul 2>&1
if errorlevel 1 (
    echo 错误：未找到Go，请先安装Go 1.21+
    echo 下载地址：https://golang.org/dl/
    pause
    exit /b 1
)

echo 正在生成派系数据...
python analyze_factions.py

echo 正在编译Go程序...
go build -o warframe-viewer.exe server.go

if errorlevel 1 (
    echo 编译失败！
    pause
    exit /b 1
)

echo 正在创建发布包...
if not exist "release" mkdir release
copy warframe-viewer.exe release\
copy index.html release\
copy styles.css release\
copy app.js release\
xcopy /E /I enemy_data release\enemy_data
copy factions_data.json release\

echo.
echo ✅ 构建完成！
echo 📁 发布文件在 release 目录中
echo 🚀 运行 release\warframe-viewer.exe 启动应用
echo.

pause 