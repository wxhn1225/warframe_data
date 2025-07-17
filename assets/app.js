let currentLanguage = 'zh';
let dictZh = null;
let dictTarget = null;
let currentEnemyData = null;
let currentEnemyName = '';
let enemyDataCache = {}; // 添加缓存对象
let factionsDataCache = null; // 缓存派系数据
let currentMode = 'normal'; // 添加游戏模式变量，默认为普通模式

// 添加全局变量存储所有敌人数据
let allEnemiesData = {};

const API_BASE = window.location.protocol + '//' + window.location.host + '/api';

// 防抖函数，避免频繁更新
function debounce(func, wait) {
    let timeout;
    return function(...args) {
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(this, args), wait);
    };
}

// 缓存语言字典
const languageDicts = {};

async function loadLanguageDict(lang) {
    // 如果已经缓存过该语言，直接返回
    if (languageDicts[lang]) {
        return languageDicts[lang];
    }
    
    const response = await fetch(`languages/dict.${lang}.json`);
    if (!response.ok) throw new Error('语言文件加载失败');
    const data = await response.json();
    
    // 缓存结果
    languageDicts[lang] = data;
    return data;
}

async function loadFactionsData() {
    try {
        // 显示加载中状态
        document.getElementById('factionsContainer').innerHTML = 
            '<div class="loading">正在加载数据...</div>';
            
        // 加载语言字典
        dictZh = await loadLanguageDict('zh');
        dictTarget = (currentLanguage === 'zh') ? dictZh : await loadLanguageDict(currentLanguage);
        
        // 清除敌人数据缓存，确保在切换模式后重新加载
        enemyDataCache = {};
        
        // 如果有缓存且模式没变，使用缓存
        if (factionsDataCache) {
            showFactions(factionsDataCache);
            return;
        }
        
        // 否则从服务器获取
        const response = await fetch(`${API_BASE}/factions`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const enemyData = await response.json();
        factionsDataCache = enemyData; // 缓存结果
        
        // 存储所有敌人数据，用于搜索
        allEnemiesData = {};
        Object.entries(enemyData).forEach(([faction, enemies]) => {
            enemies.forEach(enemy => {
                allEnemiesData[enemy] = faction;
            });
        });
        
        showFactions(enemyData);
    } catch (error) {
        document.getElementById('factionsContainer').innerHTML = 
            '<div class="error">无法连接到服务器，请确保服务器正在运行</div>';
        console.error('加载派系数据失败:', error);
    }
}

function getTranslatedName(zhName) {
    if (!dictZh || !dictTarget) return zhName;
    let key = null;
    for (const [k, v] of Object.entries(dictZh)) {
        if (v === zhName) {
            key = k;
            break;
        }
    }
    if (!key) return zhName;
    return dictTarget[key] || zhName;
}

function showFactions(enemyData) {
    document.getElementById('enemyDetail').classList.remove('show');
    document.getElementById('factionsContainer').style.display = 'grid';
    const container = document.getElementById('factionsContainer');
    container.innerHTML = '';
    
    // 如果没有传入enemyData，则重新加载数据
    if (!enemyData) {
        loadFactionsData();
        return;
    }
    
    // 使用文档片段减少DOM操作
    const fragment = document.createDocumentFragment();
    
    Object.entries(enemyData).forEach(([faction, enemies]) => {
        if (enemies.length > 0) {
            const card = document.createElement('div');
            card.className = 'faction-card';
            card.onclick = () => showEnemiesInFaction(faction, enemies);
            card.innerHTML = `
                <div class="faction-name">${faction}</div>
                <div class="enemy-count">${enemies.length} 个敌人</div>
                <div class="enemy-list">
                    <div class="enemy-item">点击查看敌人列表</div>
                </div>
            `;
            fragment.appendChild(card);
        }
    });
    
    container.appendChild(fragment);
}

function showEnemiesInFaction(faction, enemies) {
    const container = document.getElementById('factionsContainer');
    container.innerHTML = '';
    
    // 使用文档片段减少DOM操作
    const fragment = document.createDocumentFragment();
    
    const backCard = document.createElement('div');
    backCard.className = 'faction-card';
    backCard.onclick = () => loadFactionsData();
    backCard.innerHTML = `
        <div class="faction-name">← 返回派系列表</div>
        <div class="enemy-count">${faction}</div>
    `;
    fragment.appendChild(backCard);
    
    enemies.forEach(enemyName => {
        const card = document.createElement('div');
        card.className = 'faction-card';
        card.onclick = () => showEnemyDetail(enemyName);
        card.innerHTML = `
            <div class="faction-name">${getTranslatedName(enemyName)}</div>
            <div class="enemy-count">${enemyName}</div>
        `;
        fragment.appendChild(card);
    });
    
    container.appendChild(fragment);
}

async function showEnemyDetail(enemyName) {
    document.getElementById('factionsContainer').style.display = 'none';
    document.getElementById('searchResults').classList.remove('show');
    document.getElementById('enemyDetail').classList.add('show');
    currentEnemyName = enemyName;
    document.getElementById('enemyTitle').textContent = getTranslatedName(enemyName);
    
    try {
        // 显示加载状态
        document.getElementById('statsContainer').innerHTML = '<div class="loading">加载敌人数据中...</div>';
        
        // 检查缓存
        const cacheKey = `${currentMode}_${enemyName}`;
        if (enemyDataCache[cacheKey]) {
            currentEnemyData = enemyDataCache[cacheKey];
            showManualLevelInput(enemyName);
            return;
        }
        
        // 根据当前模式选择数据文件夹
        const dataFolder = currentMode === 'normal' ? 'enemy_data' : 'enemy_data_steel';
        
        // 加载敌人数据
        const response = await fetch(`data/${dataFolder}/${encodeURIComponent(enemyName)}.json`);
        if (!response.ok) {
            document.getElementById('statsContainer').innerHTML = `<div class="error">敌人数据文件不存在</div>`;
            return;
        }
        
        currentEnemyData = await response.json();
        enemyDataCache[cacheKey] = currentEnemyData; // 缓存结果
        showManualLevelInput(enemyName);
    } catch (error) {
        document.getElementById('statsContainer').innerHTML = `<div class="error">加载敌人数据失败</div>`;
        console.error('加载敌人数据失败:', error);
    }
}

function showManualLevelInput(enemyName) {
    const levelSelector = document.querySelector('.level-selector');
    levelSelector.innerHTML = `
        <label for="manualLevel">手动输入等级:</label>
        <button id="levelMinus">-</button>
        <input type="number" id="manualLevel" min="1" max="9999" value="1">
        <button id="levelPlus">+</button>
    `;
    const manualLevel = document.getElementById('manualLevel');
    
    document.getElementById('levelMinus').onclick = function() {
        let v = parseInt(manualLevel.value) || 1;
        if (v > 1) manualLevel.value = v - 1;
        showLevelData(manualLevel.value);
    };
    
    document.getElementById('levelPlus').onclick = function() {
        let v = parseInt(manualLevel.value) || 1;
        manualLevel.value = v + 1;
        showLevelData(manualLevel.value);
    };
    
    // 使用防抖函数处理输入事件
    manualLevel.addEventListener('input', debounce(function() {
        showLevelData(manualLevel.value);
    }, 200));
    
    showLevelData(manualLevel.value);
}

// 格式化数字，添加千位分隔符
function formatNumber(num) {
    if (num === undefined || num === null) return '-';
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function showLevelData(level) {
    if (!currentEnemyData || !currentEnemyData.levelData) return;
    const d = currentEnemyData.levelData[level];
    if (!d) {
        document.getElementById('statsContainer').innerHTML = `<div class="error">该等级数据不存在</div>`;
        return;
    }
    document.getElementById('statsContainer').innerHTML = `
        <div class="stat-card">
            <div class="stat-label">生命值</div>
            <div class="stat-value">${formatNumber(d.maxHealth)}</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">护甲</div>
            <div class="stat-value">${formatNumber(d.armour)}</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">护盾</div>
            <div class="stat-value">${formatNumber(d.maxShield)}</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">超宏</div>
            <div class="stat-value">${formatNumber(d.maxOverguard)}</div>
        </div>
    `;
}

// 搜索敌人
function searchEnemies(query) {
    if (!query || query.trim() === '') return [];
    
    query = query.toLowerCase();
    const results = [];
    
    // 搜索原始名称
    Object.keys(allEnemiesData).forEach(enemy => {
        if (enemy.toLowerCase().includes(query)) {
            results.push({
                name: enemy,
                faction: allEnemiesData[enemy],
                originalName: enemy
            });
        }
    });
    
    // 搜索翻译后的名称
    if (dictZh && dictTarget) {
        Object.entries(dictZh).forEach(([key, zhName]) => {
            const targetName = dictTarget[key];
            if (!targetName) return;
            
            // 检查翻译后的名称是否匹配查询
            if (targetName.toLowerCase().includes(query)) {
                // 查找对应的敌人原名
                Object.keys(allEnemiesData).forEach(enemy => {
                    if (zhName === enemy) {
                        // 避免重复添加
                        if (!results.some(r => r.name === enemy)) {
                            results.push({
                                name: enemy,
                                faction: allEnemiesData[enemy],
                                originalName: enemy
                            });
                        }
                    }
                });
            }
        });
    }
    
    return results;
}

// 显示搜索结果
function displaySearchResults(results) {
    const searchResults = document.getElementById('searchResults');
    
    if (results.length === 0) {
        searchResults.innerHTML = `
            <div class="search-title">搜索结果 <span class="search-clear" onclick="clearSearch()">清除搜索</span></div>
            <div class="no-results">未找到匹配的敌人</div>
        `;
        searchResults.classList.add('show');
        return;
    }
    
    searchResults.innerHTML = `
        <div class="search-title">搜索结果 (${results.length}) <span class="search-clear" onclick="clearSearch()">清除搜索</span></div>
    `;
    
    // 使用文档片段减少DOM操作
    const fragment = document.createDocumentFragment();
    
    results.forEach(result => {
        const card = document.createElement('div');
        card.className = 'faction-card';
        card.onclick = () => showEnemyDetail(result.name);
        card.innerHTML = `
            <div class="faction-name">${getTranslatedName(result.name)}</div>
            <div class="enemy-count">${result.faction}</div>
        `;
        fragment.appendChild(card);
    });
    
    searchResults.appendChild(fragment);
    searchResults.classList.add('show');
    
    // 隐藏派系列表
    document.getElementById('factionsContainer').style.display = 'none';
}

// 清除搜索
function clearSearch() {
    document.getElementById('searchInput').value = '';
    document.getElementById('searchResults').innerHTML = '';
    document.getElementById('searchResults').classList.remove('show');
    document.getElementById('factionsContainer').style.display = 'grid';
}

document.addEventListener('DOMContentLoaded', function() {
    loadFactionsData();
    
    document.getElementById('languageSelect').addEventListener('change', function() {
        currentLanguage = this.value;
        loadFactionsData();
    });
    
    // 添加模式选择事件监听
    document.getElementById('modeSelect').addEventListener('change', function() {
        currentMode = this.value;
        // 清除缓存，重新加载数据
        factionsDataCache = null;
        enemyDataCache = {};
        loadFactionsData();
    });
    
    // 添加搜索事件监听
    document.getElementById('searchButton').addEventListener('click', function() {
        const query = document.getElementById('searchInput').value;
        const results = searchEnemies(query);
        displaySearchResults(results);
    });
    
    // 添加回车键搜索
    document.getElementById('searchInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            const query = this.value;
            const results = searchEnemies(query);
            displaySearchResults(results);
        }
    });
});