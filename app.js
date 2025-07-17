let currentLanguage = 'zh';
let languageDict = {};
let enemyData = {};
let currentEnemy = null;

const API_BASE = 'http://localhost:5000/api';

async function loadLanguageData() {
    try {
        languageDict = {
            "zh": {
                "/Lotus/Language/EntratiLab/EntratiGeneral/AlchemyNecramechMelee": "元素弧犬"
            },
            "en": {
                "/Lotus/Language/EntratiLab/EntratiGeneral/AlchemyNecramechMelee": "Elementa Arcocanid"
            }
        };
        
        await loadFactionsData();
    } catch (error) {
        console.error('加载数据失败:', error);
        document.getElementById('factionsContainer').innerHTML = 
            '<div class="error">加载数据失败，请检查网络连接</div>';
    }
}

async function loadFactionsData() {
    try {
        const response = await fetch(`${API_BASE}/factions`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        enemyData = await response.json();
        showFactions();
    } catch (error) {
        console.error('加载派系数据失败:', error);
        document.getElementById('factionsContainer').innerHTML = 
            '<div class="error">无法连接到服务器，请确保服务器正在运行</div>';
    }
}

function showFactions() {
    document.getElementById('enemyDetail').classList.remove('show');
    document.getElementById('factionsContainer').style.display = 'grid';
    
    const container = document.getElementById('factionsContainer');
    container.innerHTML = '';

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
            
            container.appendChild(card);
        }
    });
}

function showEnemiesInFaction(faction, enemies) {
    const container = document.getElementById('factionsContainer');
    container.innerHTML = '';

    const backCard = document.createElement('div');
    backCard.className = 'faction-card';
    backCard.onclick = showFactions;
    backCard.innerHTML = `
        <div class="faction-name">← 返回派系列表</div>
        <div class="enemy-count">${faction}</div>
    `;
    container.appendChild(backCard);

    enemies.forEach(enemyName => {
        const card = document.createElement('div');
        card.className = 'faction-card';
        card.onclick = () => showEnemyDetail(enemyName);
        
        const translatedName = translateEnemyName(enemyName);
        
        card.innerHTML = `
            <div class="faction-name">${translatedName}</div>
            <div class="enemy-count">${enemyName}</div>
        `;
        
        container.appendChild(card);
    });
}

function translateEnemyName(enemyName) {
    const langDict = languageDict[currentLanguage] || {};
    for (const [key, value] of Object.entries(langDict)) {
        if (value === enemyName) {
            return value;
        }
    }
    return enemyName;
}

async function showEnemyDetail(enemyName) {
    document.getElementById('factionsContainer').style.display = 'none';
    document.getElementById('enemyDetail').classList.add('show');
    
    const translatedName = translateEnemyName(enemyName);
    document.getElementById('enemyTitle').textContent = translatedName;
    
    try {
        const response = await fetch(`${API_BASE}/enemy/${encodeURIComponent(enemyName)}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        
        if (data) {
            currentEnemy = data;
            
            const levelSelect = document.getElementById('levelSelect');
            levelSelect.innerHTML = '';
            
            const levels = Object.keys(data.levelData).sort((a, b) => parseInt(a) - parseInt(b));
            levels.forEach(level => {
                const option = document.createElement('option');
                option.value = level;
                option.textContent = `等级 ${level}`;
                levelSelect.appendChild(option);
            });
            
            if (levels.length > 0) {
                showLevelData(data, levels[0]);
            }
            
            levelSelect.addEventListener('change', function() {
                showLevelData(data, this.value);
            });
        } else {
            console.error(`无法加载敌人数据: ${enemyName}`);
        }
    } catch (error) {
        console.error(`加载敌人数据失败: ${enemyName}`, error);
        document.getElementById('statsContainer').innerHTML = 
            '<div class="error">加载敌人数据失败，请检查网络连接</div>';
    }
}

function showLevelData(enemyData, level) {
    if (!enemyData || !enemyData.levelData[level]) return;
    
    const levelData = enemyData.levelData[level];
    const container = document.getElementById('statsContainer');
    
    container.innerHTML = `
        <div class="stat-card">
            <div class="stat-label">生命值</div>
            <div class="stat-value">${levelData.maxHealth}</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">护甲</div>
            <div class="stat-value">${levelData.armour}</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">护盾</div>
            <div class="stat-value">${levelData.maxShield}</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">超宏</div>
            <div class="stat-value">${levelData.maxOverguard}</div>
        </div>
    `;
}

document.addEventListener('DOMContentLoaded', function() {
    loadLanguageData();
    document.getElementById('languageSelect').addEventListener('change', function() {
        currentLanguage = this.value;
        loadLanguageData();
    });
}); 