#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
分析 enemy_data 目录下所有敌人的 faction 类型，输出JSON文件
"""
import json
import os
from collections import defaultdict
from concurrent.futures import ThreadPoolExecutor, as_completed

def process_file(file_path, filename):
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            data = json.load(f)
        if 'faction' in data:
            faction = data['faction']
            enemy_name = data.get('enemyName', filename.replace('.json', ''))
            return (faction, enemy_name)
        else:
            print(f"警告：文件 {filename} 中没有找到faction字段")
            return None
    except json.JSONDecodeError as e:
        print(f"错误：无法解析文件 {filename}: {e}")
        return None
    except Exception as e:
        print(f"错误：处理文件 {filename} 时出错: {e}")
        return None

def analyze_factions():
    enemy_data_dir = "data/enemy_data"
    factions = defaultdict(list)  # faction -> [enemy_names]

    if not os.path.exists(enemy_data_dir):
        print(f"错误：目录 {enemy_data_dir} 不存在")
        return

    json_files = [f for f in os.listdir(enemy_data_dir) if f.endswith('.json')]
    if not json_files:
        print(f"在 {enemy_data_dir} 目录中没有找到JSON文件")
        return

    print(f"正在分析 {len(json_files)} 个敌人文件（多线程加速）...")

    with ThreadPoolExecutor(max_workers=8) as executor:
        future_to_file = {
            executor.submit(process_file, os.path.join(enemy_data_dir, filename), filename): filename
            for filename in json_files
        }
        for future in as_completed(future_to_file):
            result = future.result()
            if result:
                faction, enemy_name = result
                factions[faction].append(enemy_name)

    # 生成JSON文件
    factions_dict = dict(factions)
    
    # 确保data目录存在
    os.makedirs('data', exist_ok=True)
    
    # 保存派系数据
    with open('data/factions_data.json', 'w', encoding='utf-8') as f:
        json.dump(factions_dict, f, ensure_ascii=False, indent=2)
    
    print(f"\n✅ 派系数据已保存到 data/factions_data.json")
    
    # 输出统计信息
    print(f"\n📊 统计信息:")
    print(f"总共发现 {len(factions)} 种不同的faction")
    total_enemies = sum(len(enemies) for enemies in factions.values())
    print(f"总计：{len(factions)} 种faction，{total_enemies} 个敌人")
    
    print(f"\n📈 按敌人数量排序：")
    sorted_factions = sorted(factions.items(), key=lambda x: len(x[1]), reverse=True)
    for faction, enemies in sorted_factions:
        print(f"  {faction}: {len(enemies)} 个敌人")

if __name__ == "__main__":
    analyze_factions() 