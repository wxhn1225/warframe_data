let currentLanguage = 'zh';
let currentEnemy = null;
let currentEnemyName = '';

const API_BASE = 'http://localhost:5000/api';

async function loadFactionsData() {
    try {
        const response = await fetch(`${API_BASE}/factions?lang=${currentLanguage}`);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const enemyData = await response.json();
        showFactions(enemyData);
    } catch (error) {
        document.getElementById('factionsContainer').innerHTML = 
            '<div class="error">无法连接到服务器，请确保服务器正在运行</div>';
    }
}

function showFactions(enemyData) {
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
    backCard.onclick = () => loadFactionsData();
    backCard.innerHTML = `
        <div class="faction-name">← 返回派系列表</div>
        <div class="enemy-count">${faction}</div>
    `;
    container.appendChild(backCard);
    enemies.forEach(enemyName => {
        const card = document.createElement('div');
        card.className = 'faction-card';
        card.onclick = () => showEnemyDetail(enemyName);
        card.innerHTML = `
            <div class="faction-name">${enemyName}</div>
        `;
        container.appendChild(card);
    });
}

async function showEnemyDetail(enemyName) {
    document.getElementById('factionsContainer').style.display = 'none';
    document.getElementById('enemyDetail').classList.add('show');
    currentEnemyName = enemyName;
    document.getElementById('enemyTitle').textContent = enemyName;
    showManualLevelInput(enemyName);
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
        fetchAndShowLevelData(enemyName, manualLevel.value);
    };
    document.getElementById('levelPlus').onclick = function() {
        let v = parseInt(manualLevel.value) || 1;
        manualLevel.value = v + 1;
        fetchAndShowLevelData(enemyName, manualLevel.value);
    };
    manualLevel.addEventListener('input', function() {
        fetchAndShowLevelData(enemyName, manualLevel.value);
    });
    fetchAndShowLevelData(enemyName, manualLevel.value);
}

function fetchAndShowLevelData(enemyName, level) {
    if (!level || isNaN(level) || parseInt(level) < 1) return;
    fetch(`${API_BASE}/enemy_data?name=${encodeURIComponent(enemyName)}&level=${level}&lang=${currentLanguage}`)
      .then(res => res.json())
      .then(data => {
          if (data.error) {
              document.getElementById('statsContainer').innerHTML = `<div class="error">${data.error}</div>`;
              document.getElementById('enemyTitle').textContent = enemyName;
              return;
          }
          document.getElementById('enemyTitle').textContent = data.enemyName;
          const d = data.data;
          document.getElementById('statsContainer').innerHTML = `
            <div class="stat-card">
                <div class="stat-label">生命值</div>
                <div class="stat-value">${d.maxHealth ?? '-'}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">护甲</div>
                <div class="stat-value">${d.armour ?? '-'}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">护盾</div>
                <div class="stat-value">${d.maxShield ?? '-'}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">超宏</div>
                <div class="stat-value">${d.maxOverguard ?? '-'}</div>
            </div>
          `;
      });
}

document.addEventListener('DOMContentLoaded', function() {
    loadFactionsData();
    document.getElementById('languageSelect').addEventListener('change', function() {
        currentLanguage = this.value;
        loadFactionsData();
    });
}); 