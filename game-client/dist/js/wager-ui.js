import { apiClient } from './api.js';

export async function showBattleSelectionModal() {
    const modal = document.createElement('div');
    modal.className = 'modal-overlay';
    modal.id = 'battleModal';
    modal.innerHTML = `
        <div class="modal" style="background: linear-gradient(135deg, #18181b 0%, #27272a 100%); border: 1px solid #3f3f46; box-shadow: 0 0 40px rgba(0,0,0,0.9); max-width: 800px; width: 90%;">
            <div class="modal-header" style="border-bottom: 1px solid #3f3f46; padding: 20px;">
                <h2 style="color: #fff; font-size: 2em; display: flex; align-items: center; gap: 10px;">
                    ⚔️ Battle Arena
                </h2>
                <button class="close-btn" onclick="document.getElementById('battleModal').remove()" style="font-size: 2em;">×</button>
            </div>
            
            <div class="modal-body" style="padding: 30px; color: #a1a1aa; display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 20px;">
                
                <!-- PVE MODE -->
                <div class="battle-card" onclick="selectMode('pve')" style="
                    background: rgba(255,255,255,0.03); 
                    border: 1px solid #3f3f46; 
                    border-radius: 12px; 
                    padding: 20px; 
                    text-align: center; 
                    cursor: pointer; 
                    transition: all 0.3s ease;
                " onmouseover="this.style.borderColor='#10b981'; this.style.transform='translateY(-5px)'" onmouseout="this.style.borderColor='#3f3f46'; this.style.transform='translateY(0)'">
                    <div style="font-size: 3em; margin-bottom: 10px;">🏝️</div>
                    <h3 style="color: #fff; margin-bottom: 5px;">Island Raid</h3>
                    <p style="font-size: 0.9em;">PvE Campaign</p>
                    <span style="display: inline-block; background: #064e3b; color: #34d399; padding: 4px 8px; border-radius: 4px; font-size: 0.8em; margin-top: 10px;">Low Risk</span>
                </div>

                <!-- RANKED MODE -->
                <div class="battle-card" onclick="selectMode('ranked')" style="
                    background: rgba(255,255,255,0.03); 
                    border: 1px solid #3f3f46; 
                    border-radius: 12px; 
                    padding: 20px; 
                    text-align: center; 
                    cursor: pointer; 
                    transition: all 0.3s ease;
                " onmouseover="this.style.borderColor='#3b82f6'; this.style.transform='translateY(-5px)'" onmouseout="this.style.borderColor='#3f3f46'; this.style.transform='translateY(0)'">
                    <div style="font-size: 3em; margin-bottom: 10px;">🏆</div>
                    <h3 style="color: #fff; margin-bottom: 5px;">Ranked PvP</h3>
                    <p style="font-size: 0.9em;">Climb the Ladder</p>
                    <span style="display: inline-block; background: #1e3a8a; color: #60a5fa; padding: 4px 8px; border-radius: 4px; font-size: 0.8em; margin-top: 10px;">No Cost</span>
                </div>

                <!-- WAGER MODE -->
                <div class="battle-card" onclick="selectMode('wager')" style="
                    background: rgba(255,255,255,0.03); 
                    border: 1px solid #3f3f46; 
                    border-radius: 12px; 
                    padding: 20px; 
                    text-align: center; 
                    cursor: pointer; 
                    transition: all 0.3s ease;
                    position: relative;
                    overflow: hidden;
                " onmouseover="this.style.borderColor='#eab308'; this.style.transform='translateY(-5px)'" onmouseout="this.style.borderColor='#3f3f46'; this.style.transform='translateY(0)'">
                    <div style="position: absolute; top:0; right:0; background: #eab308; color: #000; font-size: 0.7em; font-weight: bold; padding: 2px 10px; transform: rotate(45deg) translate(8px, -8px); box-shadow: 0 2px 5px rgba(0,0,0,0.5);">HIGH STAKES</div>
                    <div style="font-size: 3em; margin-bottom: 10px;">💰</div>
                    <h3 style="color: #fff; margin-bottom: 5px;">Wager Match</h3>
                    <p style="font-size: 0.9em;">Bet GTK & Win Big</p>
                    <span style="display: inline-block; background: #422006; color: #facc15; padding: 4px 8px; border-radius: 4px; font-size: 0.8em; margin-top: 10px;">Dynamic Stake</span>
                </div>
            </div>

            <!-- DYNAMIC CONTENT CONTAINER -->
            <div id="battleContent" style="padding: 0 30px 30px 30px; display: none;">
                <!-- Content injected via JS -->
            </div>
        </div>
    `;
    document.body.appendChild(modal);

    // Attach global handler for this modal instance
    window.selectMode = (mode) => {
        const contentDiv = document.getElementById('battleContent');
        contentDiv.style.display = 'block';

        if (mode === 'pve') {
            window.location.href = 'island-raids.html'; // Direct redirect for now
        } else if (mode === 'ranked') {
            contentDiv.innerHTML = `
                <div style="text-align: center; padding: 20px; background: rgba(0,0,0,0.3); border-radius: 8px;">
                     <h3 style="color: #60a5fa;">Searching for Ranked Object...</h3>
                     <div class="spinner" style="margin: 20px auto;"></div>
                     <p>Finding opponent near your skill level...</p>
                     <button onclick="startRankedSearch()" class="premium-btn">Find Match</button>
                </div>
            `;
            // Simplified for immediate start
            startRankedSearch();
        } else if (mode === 'wager') {
            contentDiv.innerHTML = `
                <div style="text-align: center; padding: 20px; background: rgba(0,0,0,0.3); border-radius: 8px; border: 1px solid #eab308;">
                     <h3 style="color: #facc15;">High Stakes Arena</h3>
                     <p style="margin: 15px 0; color: #ddd;">
                        Stakes are calculated dynamically based on the relative strength (Combat Power) of both teams.
                        <br><br>
                        <strong>Stronger Team risks MORE to win LESS.</strong><br>
                        <strong>Weaker Team risks LESS to win MORE.</strong>
                     </p>
                     <div style="margin: 20px 0;">
                        <span style="color: #a1a1aa;">Required Wallet Balance:</span>
                        <span style="color: #fff; font-weight: bold;">500+ GTK</span>
                     </div>
                     <p id="wagerStatusText" style="color: #eab308; display: none;">Searching for worthy opponent...</p>
                     <button id="startWagerBtn" onclick="startWagerSearch()" class="premium-btn" style="background: linear-gradient(45deg, #eab308, #ca8a04); color: #000; width: 100%; border: none;">
                        🔎 Find High Stakes Match
                     </button>
                </div>
            `;
        }
    };

    window.startRankedSearch = async () => {
        // Placeholder logic
        alert("Ranked matchmaking coming in next update!");
    };

    window.startWagerSearch = async () => {
        const btn = document.getElementById('startWagerBtn');
        const status = document.getElementById('wagerStatusText');

        btn.disabled = true;
        btn.innerHTML = '<div class="spinner"></div> Searching...';
        status.style.display = 'block';

        try {
            // Get Team
            const teamRes = await apiClient.get('/teams/active');
            if (!teamRes.team) {
                alert('No active team selected!');
                return;
            }

            // Start Wager (Dynamic) - No amount sent
            const res = await apiClient.post('/battle/wager', {
                team_id: teamRes.team.id
            });

            if (res.status === 'active') {
                // Match Found!
                status.innerHTML = `<span style="color: #4ade80">Match Found! Stakes: You ${res.your_stake} GTK vs Enemy ${res.enemy_stake} GTK</span>`;
                setTimeout(() => {
                    window.location.href = `battle.html?id=${res.battle_id}`;
                }, 1500);
            } else {
                status.innerHTML = "Funds Verified. In Queue (Auto-match when opponent appears)...";
                // Polling logic would go here
            }

        } catch (error) {
            console.error(error);
            btn.disabled = false;
            btn.innerText = "Try Again";
            status.innerHTML = `<span style="color: #ef4444">${error.response?.data?.error || 'Search failed'}</span>`;
        }
    };
}