@echo off
echo ========================================
echo Warframe 敌人数据查看器
echo ========================================
echo.

echo 正在生成派系数据...
python src\analyze_factions.py

echo 启动服务器...
echo 请在浏览器中访问: http://localhost:5000
echo 按 Ctrl+C 停止服务器
echo.
warframe-viewer.exe

pause 