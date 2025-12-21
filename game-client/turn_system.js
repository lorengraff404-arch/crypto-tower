// Turn-based combat helper functions (Phase 13)

// Update turn indicator displaying whose turn it is
function updateTurnDisplay(turnQueue, currentIndex) {
    const turn = turnQueue[currentIndex % turnQueue.length];
    const turnText = document.getElementById('currentTurnText');
    const queuePreview = document.getElementById('turnQueuePreview');
    
    if (!turn) return;
    
    // Display current turn
    if (turn.type === 'player') {
        turnText.textContent = `${turn.name}'s Turn!`;
        turnText.style.color = '#10b981'; // Green for player
    } else {
        turnText.textContent = `Enemy's Turn!`;
        turnText.style.color = '#ef4444'; // Red for enemy
    }
    
    // Show next 3 turns in queue
    const next3 = [];
    for (let i = 1; i <= 3; i++) {
        const nextTurn = turnQueue[(currentIndex + i) % turnQueue.length];
        next3.push(nextTurn.name);
    }
    queuePreview.textContent = `Next: ${next3.join(' → ')}`;
}

// Auto-execute enemy turn (AI)
async function autoExecuteEnemyTurn() {
    if (!currentSession) return;
    
    try {
        const res = await apiClient.request('POST', `/raids/${currentSession.id}/auto-execute`);
        currentSession = res.session;
        
        // Display battle result
        if (res.battle_result) {
            addLog(res.battle_result.message, '#ef4444'); // Red for enemy attacks
        }
        
        // Update UI
        updateBattleUI(currentSession);
        
        // Update turn display
        if (currentSession.turn_queue) {
            const turnQueue = JSON.parse(currentSession.turn_queue);
            updateTurnDisplay(turnQueue, currentSession.current_turn_index);
        }
        
        // Check if battle ended
        if (currentSession.status !== 'IN_PROGRESS') {
            handleBattle End();
        } else {
            // Check if next turn is also enemy (shouldn't happen but just in case)
            const turnQueue = JSON.parse(currentSession.turn_queue);
            const nextTurn = turnQueue[currentSession.current_turn_index % turnQueue.length];
            if (nextTurn.type === 'enemy') {
                setTimeout(() => autoExecuteEnemyTurn(), 2000);
            }
        }
    } catch (e) {
        console.error('Enemy turn failed:', e);
        showNotification('Enemy turn failed: ' + e.message, 'error', 'BATTLE ERROR');
    }
}

// Handle battle end (victory/defeat)
function handleBattleEnd() {
    if (currentSession.status === 'COMPLETED') {
        setTimeout(async () => {
            await showVictoryModal(currentSession);
        }, 1500);
    } else if (currentSession.status === 'FAILED') {
        setTimeout(() => {
            showNotification('DEFEAT! Your team was defeated.', 'error', '💀 DEFEAT');
            location.reload();
        }, 1500);
    }
}

// Execute player turn
async function executePlayerTurn(characterID, moveSlot) {
    if (!currentSession) return;
    
    try {
        const res = await apiClient.request('POST', `/raids/${currentSession.id}/execute-turn`, {
            character_id: characterID,
            move_slot: moveSlot
        });
        
        currentSession = res.session;
        
        // Display battle result
        if (res.battle_result) {
            const result = res.battle_result;
            addLog(result.message, '#fbbf24'); // Gold for player attacks
            
            // Animate damage
            const bossSprite = document.getElementById('bossSprite');
            if (result.damage > 0) {
                const dmgClass = result.effectiveness >= 2.0 ? 'player-dmg-super' :
                    result.effectiveness <= 0.5 ? 'player-dmg-weak' : 'player-dmg';
                spawnDamageNumber(result.damage, dmgClass, bossSprite);
                bossSprite.classList.add('flash-hit');
                setTimeout(() => bossSprite.classList.remove('flash-hit'), 300);
            }
        }
        
        // Update UI
        updateBattleUI(currentSession);
        
        // Update turn display
        if (currentSession.turn_queue) {
            const turnQueue = JSON.parse(currentSession.turn_queue);
            updateTurnDisplay(turnQueue, currentSession.current_turn_index);
        }
        
        // Check if battle ended
        if (currentSession.status !== 'IN_PROGRESS') {
            handleBattleEnd();
        } else {
            // Check if next turn is enemy's
            const turnQueue = JSON.parse(currentSession.turn_queue);
            const nextTurn = turnQueue[currentSession.current_turn_index % turnQueue.length];
            if (nextTurn.type === 'enemy') {
                setTimeout(() => autoExecuteEnemyTurn(), 2000);
            }
        }
        
        // Hide move selection
        document.getElementById('moveSelectionMenu').style.display = 'none';
        document.getElementById('mainActionMenu').style.display = 'grid';
        
    } catch (e) {
        console.error('Player turn failed:', e);
        showNotification('Failed to execute move: ' + e.message, 'error', 'BATTLE ERROR');
    }
}
