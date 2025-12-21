// Quests.ts - Type-safe Mission/Quest UI Logic

import { apiClient } from './api.js';
import type {
    Mission,
    MissionListResponse,
    MissionCompleteResponse,
    DialogueResponse,
    StoryProgress,
    User
} from './types';

// State management with types
let userToken: string | null = null;
let currentMission: Mission | null = null;
let userStoryProgress: StoryProgress | null = null;

// Initialize
document.addEventListener('DOMContentLoaded', async (): Promise<void> => {
    userToken = localStorage.getItem('token');

    if (!userToken) {
        alert('Please connect your wallet first');
        window.location.href = 'index.html';
        return;
    }

    await loadUserStats();
    await loadStoryProgress();
    await loadMissions();
});

// Load user stats
async function loadUserStats(): Promise<void> {
    try {
        const data = await apiClient.get<User>('/auth/profile');

        const levelElement = document.getElementById('user-level');
        const xpElement = document.getElementById('user-xp');

        if (levelElement) levelElement.textContent = `Level: ${data.level || 1}`;
        if (xpElement) xpElement.textContent = `XP: ${data.experience || 0}`;
    } catch (error) {
        console.error('Error loading user stats:', error);
    }
}

// Load story progress
async function loadStoryProgress(): Promise<void> {
    try {
        const data = await apiClient.get<{ progress: StoryProgress }>('/story/progress');
        userStoryProgress = data.progress;

        // Update Aria corruption bar
        const corruptionPercent = userStoryProgress.aria_corruption || 0;
        const corruptionFill = document.querySelector('.corruption-fill') as HTMLElement;
        const corruptionText = document.querySelector('.corruption-text');

        if (corruptionFill) corruptionFill.style.width = `${corruptionPercent}%`;
        if (corruptionText) corruptionText.textContent = `Corruption: ${corruptionPercent}%`;
    } catch (error) {
        console.error('Error loading story progress:', error);
    }
}

// Load missions
async function loadMissions(): Promise<void> {
    try {
        const data = await apiClient.get<MissionListResponse>('/missions');
        displayMissions(data.missions || []);

        // Update mission count
        const completedCount = (data.missions || []).filter(m => m.status === 'completed').length;
        const missionsComplete = document.getElementById('missions-complete');
        if (missionsComplete) {
            missionsComplete.textContent = `Missions: ${completedCount}/30`;
        }
    } catch (error) {
        console.error('Error loading missions:', error);
    }
}

// Display missions
function displayMissions(missions: Mission[]): void {
    const container = document.getElementById('mission-list');
    if (!container) return;

    container.innerHTML = '';

    missions.forEach((mission: Mission) => {
        const card = document.createElement('div');
        card.className = `mission-card ${mission.status}`;
        if (mission.status === 'locked') card.classList.add('locked');

        const statusClass = `status-${mission.status.replace('_', '-')}`;

        card.innerHTML = `
      <div class="mission-card-header">
        <span class="mission-level">L${mission.level}</span>
        <span class="mission-status ${statusClass}">${mission.status}</span>
      </div>
      <div class="mission-name">${mission.name}</div>
      <div class="mission-description">${mission.description?.substring(0, 60)}...</div>
    `;

        if (mission.status !== 'locked') {
            card.onclick = () => void showMissionDetail(mission);
        }

        container.appendChild(card);
    });
}

// Show mission detail
async function showMissionDetail(mission: Mission): Promise<void> {
    currentMission = mission;

    // Remove active class from all cards
    document.querySelectorAll('.mission-card').forEach(c => c.classList.remove('active'));

    // Add active to selected
    const cards = Array.from(document.querySelectorAll('.mission-card'));
    const clickedCard = cards.find(c => c.textContent?.includes(mission.name));
    if (clickedCard) clickedCard.classList.add('active');

    // Load dialogue if exists
    await loadMissionDialogue(mission.level);

    // Display mission detail
    const detailContainer = document.getElementById('mission-detail');
    if (!detailContainer) return;

    detailContainer.innerHTML = `
    <div class="mission-detail-content">
      <div class="mission-detail-header">
        <h2>${mission.name}</h2>
        <div class="mission-detail-meta">
          <span class="mission-level">Level ${mission.level}</span>
          <span class="mission-status status-${mission.status.replace('_', '-')}">${mission.status}</span>
        </div>
      </div>

      <div class="mission-objectives">
        <h3>Objectives</h3>
        <div class="objective-list" id="objective-list"></div>
      </div>

      <div class="mission-rewards">
        <h3>Rewards</h3>
        <div class="rewards-grid" id="rewards-grid"></div>
      </div>

      ${mission.unlock_feature ? `
        <div class="unlock-alert">
          <strong>Unlocks:</strong> ${formatFeatureName(mission.unlock_feature)}
        </div>
      ` : ''}

      <div class="mission-actions">
        ${getMissionActions(mission)}
      </div>
    </div>
  `;

    displayObjectives(mission);
    displayRewards(mission);
}

// Load mission dialogue
async function loadMissionDialogue(level: number): Promise<void> {
    try {
        const data = await apiClient.get<DialogueResponse>(`/story/missions/${level}/dialogues`);

        if (currentMission?.status === 'available' && data.briefings && data.briefings.length > 0 && data.briefings[0]) {
            showDialogue(data.briefings[0].dialogue_text);
        }
    } catch (error) {
        console.error('Error loading dialogue:', error);
    }
}

// Show Aria dialogue
function showDialogue(text: string): void {
    const dialogueContent = document.getElementById('dialogue-content');
    const ariaDialogue = document.getElementById('aria-dialogue');

    if (dialogueContent) dialogueContent.textContent = text;
    if (ariaDialogue) ariaDialogue.classList.remove('hidden');
}

// Close dialogue
function closeDialogue(): void {
    const ariaDialogue = document.getElementById('aria-dialogue');
    if (ariaDialogue) ariaDialogue.classList.add('hidden');
}

// Display objectives
function displayObjectives(mission: Mission): void {
    const container = document.getElementById('objective-list');
    if (!container) return;

    const objectivesProgress = mission.objectives_progress || [];

    mission.objectives.forEach((obj, index) => {
        const current = objectivesProgress[index]?.current || 0;
        const target = obj.target || 1;
        const completed = current >= target;
        const progressPercent = Math.min((current / target) * 100, 100);

        const objDiv = document.createElement('div');
        objDiv.className = 'objective-item';
        objDiv.innerHTML = `
      <span class="objective-icon">${completed ? '✅' : '⚪'}</span>
      <div class="objective-text">
        <div>${obj.description || obj.type}</div>
        <div class="objective-progress">${current} / ${target}</div>
      </div>
      <div class="progress-bar">
        <div class="progress-fill" style="width: ${progressPercent}%"></div>
      </div>
    `;
        container.appendChild(objDiv);
    });
}

// Display rewards
function displayRewards(mission: Mission): void {
    const container = document.getElementById('rewards-grid');
    if (!container) return;

    const rewardIcons: Record<string, string> = {
        'gtk': '💰',
        'experience': '⭐',
        'materials': '⚙️',
        'items': '📦'
    };

    for (const [type, value] of Object.entries(mission.rewards)) {
        if (Array.isArray(value)) {
            value.forEach(item => {
                const rewardDiv = createRewardItem(item.type || type, item.quantity || 0, rewardIcons);
                container.appendChild(rewardDiv);
            });
        } else if (typeof value === 'number') {
            const rewardDiv = createRewardItem(type, value, rewardIcons);
            container.appendChild(rewardDiv);
        }
    }
}

// Create reward item
function createRewardItem(type: string, amount: number, icons: Record<string, string>): HTMLDivElement {
    const div = document.createElement('div');
    div.className = 'reward-item';
    div.innerHTML = `
    <div class="reward-icon">${icons[type] || '🎁'}</div>
    <div class="reward-name">${type.toUpperCase()}</div>
    <div class="reward-amount">${amount}</div>
  `;
    return div;
}

// Get mission action buttons
function getMissionActions(mission: Mission): string {
    switch (mission.status) {
        case 'available':
            return '<button onclick="window.startMission()" class="btn-primary">Start Mission</button>';
        case 'in_progress':
            return '<button onclick="window.completeMission()" class="btn-primary">Complete Mission</button>';
        case 'completed':
            return '<button disabled class="btn-secondary">Completed</button>';
        default:
            return '';
    }
}

// Start mission
async function startMission(): Promise<void> {
    if (!currentMission) return;

    const missionNameElement = document.getElementById('start-mission-name');
    const startModal = document.getElementById('start-modal');

    if (missionNameElement) missionNameElement.textContent = currentMission.name;
    if (startModal) {
        startModal.classList.remove('hidden');
        (startModal as HTMLElement).style.display = 'flex';
    }
}

// Confirm start mission
async function confirmStartMission(): Promise<void> {
    if (!currentMission) return;

    try {
        await apiClient.post(`/missions/${currentMission.id}/start`);
        closeStartModal();
        await loadMissions();
        await showMissionDetail(currentMission);
    } catch (error) {
        if (error instanceof Error) {
            alert(`Error: ${error.message}`);
        }
    }
}

// Complete mission
async function completeMission(): Promise<void> {
    if (!currentMission) return;

    try {
        const data = await apiClient.post<MissionCompleteResponse>(
            `/missions/${currentMission.id}/complete`
        );
        showCompletionModal(data);
    } catch (error) {
        if (error instanceof Error) {
            alert(`Error: ${error.message}`);
        }
    }
}

// Show completion modal
function showCompletionModal(data: MissionCompleteResponse): void {
    const container = document.getElementById('rewards-container');
    if (!container) return;

    container.innerHTML = '';

    const rewards = data.rewards || {};
    const rewardIcons: Record<string, string> = { 'gtk': '💰', 'experience': '⭐', 'materials': '⚙️' };

    for (const [type, value] of Object.entries(rewards)) {
        if (typeof value === 'number') {
            const div = createRewardItem(type, value, rewardIcons);
            container.appendChild(div);
        }
    }

    // Show unlock if exists
    const unlockNotification = document.getElementById('unlock-notification');
    const unlockFeature = document.getElementById('unlock-feature');

    if (data.unlock_feature) {
        if (unlockNotification) unlockNotification.classList.remove('hidden');
        if (unlockFeature) unlockFeature.textContent = formatFeatureName(data.unlock_feature);
    } else {
        if (unlockNotification) unlockNotification.classList.add('hidden');
    }

    const completeModal = document.getElementById('complete-modal');
    if (completeModal) {
        completeModal.classList.remove('hidden');
        (completeModal as HTMLElement).style.display = 'flex';
    }
}

// Close modals
function closeStartModal(): void {
    const modal = document.getElementById('start-modal');
    if (modal) {
        modal.classList.add('hidden');
        (modal as HTMLElement).style.display = 'none';
    }
}

function closeCompleteModal(): void {
    const modal = document.getElementById('complete-modal');
    if (modal) {
        modal.classList.add('hidden');
        (modal as HTMLElement).style.display = 'none';
    }
    loadUserStats();
    loadMissions();
}

// Format feature name
function formatFeatureName(feature: string): string {
    return feature.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
}

// Setup window functions for onclick handlers
declare global {
    interface Window {
        startMission: () => Promise<void>;
        completeMission: () => Promise<void>;
        confirmStartMission: () => Promise<void>;
        closeDialogue: () => void;
        closeStartModal: () => void;
        closeCompleteModal: () => void;
    }
}

window.startMission = startMission;
window.completeMission = completeMission;
window.confirmStartMission = confirmStartMission;
window.closeDialogue = closeDialogue;
window.closeStartModal = closeStartModal;
window.closeCompleteModal = closeCompleteModal;

// Close modals on outside click
window.onclick = (event: MouseEvent): void => {
    const startModal = document.getElementById('start-modal');
    const completeModal = document.getElementById('complete-modal');

    if (event.target === startModal) closeStartModal();
    if (event.target === completeModal) closeCompleteModal();
};

// Export for testing
export { loadMissions, startMission, completeMission };
