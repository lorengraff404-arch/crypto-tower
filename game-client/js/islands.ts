// Islands.ts - Type-safe Island Raids UI Logic

import { apiClient } from './api.js';
import type {
    Island,
    IslandListResponse,
    IslandDetailResponse,
    IslandRaid,
    RaidCompleteResponse,
    Enemy,
    User
} from './types';

// State management
let currentIsland: Island | null = null;
let userToken: string | null = null;
let raidData: IslandRaid | null = null;

// Initialize on page load
document.addEventListener('DOMContentLoaded', async (): Promise<void> => {
    userToken = localStorage.getItem('token');

    if (!userToken) {
        alert('Please connect your wallet first');
        window.location.href = 'index.html';
        return;
    }

    await loadUserStats();
    await loadIslands();
    setupEventListeners();
});

// Load user stats with type safety
async function loadUserStats(): Promise<void> {
    try {
        const data = await apiClient.get<User>('/auth/profile');

        const gtkElement = document.getElementById('user-gtk');
        const levelElement = document.getElementById('user-level');

        if (gtkElement) gtkElement.textContent = `GTK: ${data.gtk_balance || 0}`;
        if (levelElement) levelElement.textContent = `Level: ${data.level || 1}`;

        await loadRaidCount();
    } catch (error) {
        console.error('Error loading user stats:', error);
    }
}

// Load daily raid count
async function loadRaidCount(): Promise<number> {
    try {
        const data = await apiClient.get<IslandListResponse>('/islands');
        const raidsToday = data.raids_today || 0;
        const raidsRemaining = Math.max(0, 5 - raidsToday);

        const raidsElement = document.getElementById('raids-today');
        if (raidsElement) {
            raidsElement.textContent = `Raids: ${raidsToday}/5`;
        }

        return raidsRemaining;
    } catch (error) {
        console.error('Error loading raid count:', error);
        return 5;
    }
}

// Load islands from API
async function loadIslands(): Promise<void> {
    try {
        const data = await apiClient.get<IslandListResponse>('/islands');
        displayIslands(data.islands || []);
    } catch (error) {
        console.error('Error loading islands:', error);
        const container = document.getElementById('island-selection');
        if (container) {
            container.innerHTML = `
        <div style="grid-column: 1/-1; text-align: center; color: #ff6b6b;">
          <p>Failed to load islands. Please try again.</p>
        </div>
      `;
        }
    }
}

// Display islands with type safety
function displayIslands(islands: Island[]): void {
    const container = document.getElementById('island-selection');
    if (!container) return;

    container.innerHTML = '';

    const icons: Record<string, string> = {
        'forest': '🌳',
        'industrial': '⚙️',
        'crystal': '💎',
        'desert': '🏜️',
        'ocean': '🌊',
        'volcanic': '🌋'
    };

    islands.forEach((island: Island) => {
        const card = document.createElement('div');
        card.className = 'island-card';
        card.onclick = () => openIslandModal(island);

        const icon = icons[island.biome] || '🏝️';
        const difficultyClass = `difficulty-${island.difficulty}`;

        card.innerHTML = `
      <div class="island-card-content">
        <div class="island-icon">${icon}</div>
        <h3 class="island-name">${island.name}</h3>
        <div class="island-biome">${island.biome}</div>
        <div class="island-stats">
          <div class="island-stat">
            <span class="island-stat-label">Entry Fee</span>
            <span class="island-stat-value">${island.entry_fee} GTK</span>
          </div>
          <div class="island-stat">
            <span class="island-stat-label">Reward</span>
            <span class="island-stat-value">${island.base_reward} GTK</span>
          </div>
          <div class="island-stat">
            <span class="island-stat-label">Waves</span>
            <span class="island-stat-value">${island.wave_count}</span>
          </div>
        </div>
        <div class="island-difficulty ${difficultyClass}">${island.difficulty}</div>
      </div>
    `;

        container.appendChild(card);
    });
}

// Open island modal with type safety
async function openIslandModal(island: Island): Promise<void> {
    currentIsland = island;

    const elements = {
        name: document.getElementById('modal-island-name'),
        biome: document.getElementById('modal-biome'),
        difficulty: document.getElementById('modal-difficulty'),
        entryFee: document.getElementById('modal-entry-fee'),
        reward: document.getElementById('modal-reward'),
        waves: document.getElementById('modal-waves')
    };

    if (elements.name) elements.name.textContent = island.name;
    if (elements.biome) elements.biome.textContent = island.biome;
    if (elements.difficulty) elements.difficulty.textContent = island.difficulty;
    if (elements.entryFee) elements.entryFee.textContent = String(island.entry_fee);
    if (elements.reward) elements.reward.textContent = String(island.base_reward);
    if (elements.waves) elements.waves.textContent = String(island.wave_count);

    await loadIslandDetails(island.id);

    const raidsRemaining = await loadRaidCount();
    const remainingElement = document.getElementById('modal-raids-remaining');
    if (remainingElement) {
        remainingElement.textContent = String(raidsRemaining);
    }

    const modal = document.getElementById('island-modal');
    if (modal) {
        modal.classList.remove('hidden');
        (modal as HTMLElement).style.display = 'flex';
    }
}

// Load island details
async function loadIslandDetails(islandId: number): Promise<void> {
    try {
        const data = await apiClient.get<IslandDetailResponse>(`/islands/${islandId}`);
        displayEnemies(data.enemies || []);
        displayDrops(data.island.drop_rates);
    } catch (error) {
        console.error('Error loading island details:', error);
    }
}

// Display enemies
function displayEnemies(enemies: Enemy[]): void {
    const container = document.getElementById('modal-enemies');
    if (!container) return;

    container.innerHTML = '<h3>Enemies</h3><div class="enemy-list"></div>';
    const enemyList = container.querySelector('.enemy-list');
    if (!enemyList) return;

    enemies.forEach((enemy: Enemy) => {
        const enemyDiv = document.createElement('div');
        enemyDiv.className = 'enemy-item';

        enemyDiv.innerHTML = `
      <div>
        <span class="enemy-name">${enemy.name}</span>
        <span class="enemy-type"> - ${enemy.type}</span>
      </div>
      ${enemy.is_boss ? '<span class="boss-tag">BOSS</span>' : ''}
    `;

        enemyList.appendChild(enemyDiv);
    });
}

// Display drop rates
function displayDrops(dropRates: Record<string, number>): void {
    const container = document.getElementById('modal-drops');
    if (!container) return;

    container.innerHTML = '<h3>Possible Rewards</h3><div class="drop-list"></div>';
    const dropList = container.querySelector('.drop-list');
    if (!dropList) return;

    for (const [item, probability] of Object.entries(dropRates)) {
        const dropDiv = document.createElement('div');
        dropDiv.className = 'drop-item';

        dropDiv.innerHTML = `
      <span>${item}</span>
      <span class="drop-probability">${(probability * 100).toFixed(0)}%</span>
    `;

        dropList.appendChild(dropDiv);
    }
}

// Enter raid
async function enterRaid(): Promise<void> {
    if (!currentIsland) return;

    const raidsRemaining = await loadRaidCount();
    if (raidsRemaining <= 0) {
        alert('Daily raid limit reached (5/5). Try again tomorrow!');
        return;
    }

    try {
        // Backend expects: POST /raids/start { island_id: 1, team_id: 1 }
        // We hardcode team_id: 1 for now or fetch user's active team
        const data = await apiClient.post<{ session: IslandRaid }>(
            '/raids/start',
            { island_id: currentIsland.id, team_id: 1 }
        );

        raidData = data.session;
        closeIslandModal();
        startRaid();
    } catch (error) {
        if (error instanceof Error) {
            alert(`Error: ${error.message}`);
        }
    }
}

// Start raid
function startRaid(): void {
    if (!currentIsland || !raidData) return;

    const selectionElement = document.getElementById('island-selection');
    const battleElement = document.getElementById('raid-battle');

    if (selectionElement) selectionElement.classList.add('hidden');
    if (battleElement) battleElement.classList.remove('hidden');

    renderBattleUI();
}

// Render Battle UI
function renderBattleUI(): void {
    if (!raidData) return;

    const currentWave = document.getElementById('current-wave');
    const totalWaves = document.getElementById('total-waves');
    const playerHp = document.getElementById('player-hp');

    if (currentWave) currentWave.textContent = String(raidData.current_stage || 1);
    if (totalWaves) totalWaves.textContent = String(raidData.total_stages || currentIsland?.wave_count || 1);
    if (playerHp) playerHp.textContent = String(raidData.current_team_hp);

    const enemiesContainer = document.getElementById('enemies-container');
    if (enemiesContainer) {
        enemiesContainer.innerHTML = `
            <div style="text-align:center; padding: 20px;">
                <h3>Boss HP: ${raidData.current_boss_hp}</h3>
                <div style="background:#333; height:20px; width:100%; border:1px solid white;">
                    <div style="background:red; height:100%; width:${Math.min(100, raidData.current_boss_hp / 10)}%"></div> 
                </div>
                <br>
                <h3>Team HP: ${raidData.current_team_hp}</h3>
                <div style="background:#333; height:20px; width:100%; border:1px solid white;">
                    <div style="background:green; height:100%; width:${Math.min(100, raidData.current_team_hp / 50)}%"></div> 
                </div>
                <hr>
                <div id="battle-actions" style="display:flex; gap:10px; justify-content:center; margin-top:20px;">
                    <button id="btn-attack" class="menu-btn">⚔️ Basic Attack</button>
                    <button id="btn-skill" class="menu-btn">✨ Skill</button>
                    <button id="btn-flee" class="menu-btn" style="background:#ef4444;">🏃 Flee</button>
                </div>
                <div id="battle-log" style="margin-top:20px; color:#ccc; font-size:12px; height:100px; overflow-y:auto; border:1px solid #444; padding:10px;">
                    Battle Started!
                </div>
            </div>
        `;

        // Bind events
        document.getElementById('btn-attack')?.addEventListener('click', () => performAction('attack', 0));
        document.getElementById('btn-skill')?.addEventListener('click', () => performAction('skill', 0)); // Slot 0 for simplicity
        document.getElementById('btn-flee')?.addEventListener('click', () => endRaid());
    }
}

async function performAction(actionType: string, moveSlot: number) {
    if (!raidData) return;

    // Pick first character for testing? Or active character.
    // Backend expects: { character_id: uint, action_type: string, move_slot: int }
    // We need to know which character is active. 
    // For MVP, we assume Character 1.
    // But wait, raidData lacks active_character_id initially? Backend sets it.
    // Let's assume character_id 1 for now if we don't have it, or fetch from character_states.
    // Simplification: Send 0 and let backend handle or user first available.
    // Check backend: RaidTurnHandler parses CharacterID.

    addLog(`Using ${actionType}...`);

    try {
        const result = await apiClient.post<any>(`/raids/${raidData.id}/turn`, {
            character_id: 1, // TODO: Get from UI selection
            action_type: actionType,
            move_slot: moveSlot
        });

        // Update state
        raidData = result.session;
        renderBattleUI();
        addLog(`Turn Results: Damage Dealt: ${result.result?.damageDealt || 0}, Damage Taken: ${result.result?.damageTaken || 0}`);

        if (raidData?.status === 'COMPLETED') {
            alert("Victory!");
            endRaid(true);
        } else if (raidData?.status === 'FAILED') {
            alert("Defeat!");
            closeLootModal(); // effectively close
        }

    } catch (e: any) {
        alert(e.message);
    }
}

function addLog(msg: string) {
    const log = document.getElementById('battle-log');
    if (log) {
        log.innerHTML += `<div>${msg}</div>`;
        log.scrollTop = log.scrollHeight;
    }
}

// End raid (Complete or Flee)
async function endRaid(victory: boolean = false): Promise<void> {
    if (!raidData || !currentIsland) return;

    try {
        if (victory) {
            const data = await apiClient.post<RaidCompleteResponse>(`/raids/${raidData.id}/complete`, {});
            showLootModal(data.raid, data.rewards);
        } else {
            await apiClient.post(`/raids/${raidData.id}/flee`, {});
            closeLootModal(); // Return to map
        }
    } catch (error) {
        if (error instanceof Error) {
            alert(`Error: ${error.message}`);
        }
    }
}

// Show loot modal
function showLootModal(raid: IslandRaid, rewards: RaidCompleteResponse['rewards']): void {
    const wavesElement = document.getElementById('loot-waves');
    const hpElement = document.getElementById('loot-hp');
    const perfectIndicator = document.getElementById('perfect-indicator');

    if (wavesElement) wavesElement.textContent = String(raid.current_stage || raid.waves_cleared || 0);
    if (hpElement) hpElement.textContent = String(raid.current_team_hp);

    // Show perfect clear bonus if applicable
    // Assume perfect if full HP (requires initial_team_hp or knowing max)
    // If initial_team_hp is missing, fallback to > 0 check? Or just omit logic for now.
    const isPerfect = raid.initial_team_hp ? (raid.current_team_hp >= raid.initial_team_hp) : false;

    if (isPerfect && perfectIndicator) {
        perfectIndicator.classList.remove('hidden');
    } else if (perfectIndicator) {
        perfectIndicator.classList.add('hidden');
    }

    // Display rewards
    const container = document.getElementById('loot-items');
    if (!container) return;

    container.innerHTML = '';

    const rewardIcons: Record<string, string> = {
        'gtk': '💰',
        'seeds': '🌱',
        'materials': '⚙️',
        'shards': '💎',
        'runes': '📜',
        'eggs': '🥚'
    };

    // Handle GTK reward
    if (rewards.gtk) {
        const rewardDiv = document.createElement('div');
        rewardDiv.className = 'reward-item';
        rewardDiv.innerHTML = `
      <div class="reward-icon">${rewardIcons['gtk']}</div>
      <div class="reward-name">GTK</div>
      <div class="reward-quantity">${rewards.gtk}</div>
    `;
        container.appendChild(rewardDiv);
    }

    // Handle item rewards
    if (rewards.items && Array.isArray(rewards.items)) {
        rewards.items.forEach(item => {
            const rewardDiv = document.createElement('div');
            rewardDiv.className = 'reward-item';
            rewardDiv.innerHTML = `
        <div class="reward-icon">${rewardIcons[item.type] || '📦'}</div>
        <div class="reward-name">${item.name || item.type}</div>
        <div class="reward-quantity">x${item.quantity}</div>
      `;
            container.appendChild(rewardDiv);
        });
    }

    // Show modal
    const lootModal = document.getElementById('loot-modal');
    if (lootModal) {
        lootModal.classList.remove('hidden');
        (lootModal as HTMLElement).style.display = 'flex';
    }
}

// Close modals
function closeIslandModal(): void {
    const modal = document.getElementById('island-modal');
    if (modal) {
        modal.classList.add('hidden');
        (modal as HTMLElement).style.display = 'none';
    }
}

function closeLootModal(): void {
    const modal = document.getElementById('loot-modal');
    if (modal) {
        modal.classList.add('hidden');
        (modal as HTMLElement).style.display = 'none';
    }

    const battleElement = document.getElementById('raid-battle');
    const selectionElement = document.getElementById('island-selection');

    if (battleElement) battleElement.classList.add('hidden');
    if (selectionElement) selectionElement.classList.remove('hidden');

    loadUserStats();
    loadIslands();
}

// Setup event listeners
function setupEventListeners(): void {
    document.querySelectorAll('.close-modal').forEach((btn: Element) => {
        (btn as HTMLElement).onclick = closeIslandModal;
    });

    window.onclick = (event: MouseEvent) => {
        const islandModal = document.getElementById('island-modal');
        const lootModal = document.getElementById('loot-modal');

        if (event.target === islandModal) closeIslandModal();
        if (event.target === lootModal) closeLootModal();
    };

    const enterButton = document.getElementById('btn-enter-raid');
    const endButton = document.getElementById('btn-end-raid');

    if (enterButton) enterButton.onclick = () => void enterRaid();
    if (endButton) endButton.onclick = () => void endRaid();
}

// Export for potential module usage
export { loadIslands, enterRaid, endRaid };
