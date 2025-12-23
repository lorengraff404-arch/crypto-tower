export async function showWagerModal() {
    const modal = document.createElement('div');
    modal.className = 'modal-overlay';
    modal.id = 'wagerModal';
    modal.innerHTML = `
        <div class="modal" style="background: linear-gradient(135deg, #1e1e24 0%, #2a2a35 100%); border: 1px solid #4a4a5e; box-shadow: 0 0 20px rgba(0,0,0,0.8);">
            <div class="modal-header" style="border-bottom: 1px solid #4a4a5e;">
                <h2 style="color: #ffcc00; text-shadow: 0 0 10px rgba(255, 204, 0, 0.5);">⚔️ High Stakes WAGER</h2>
                <button class="close-btn" onclick="document.getElementById('wagerModal').remove()">×</button>
            </div>
            <div class="modal-body" style="padding: 20px; color: #fff;">
                <p>Enter the arena and bet your GTK. Winner takes all (dynamic multiplier).</p>
                
                <div style="margin: 20px 0; text-align: center;">
                    <label style="display: block; margin-bottom: 10px; font-size: 1.2em;">Wager Amount (GTK)</label>
                    <input type="number" id="wagerAmount" value="100" min="100" step="10" style="
                        padding: 10px; 
                        font-size: 1.5em; 
                        width: 150px; 
                        text-align: center; 
                        background: #111; 
                        border: 2px solid #ffcc00; 
                        color: #fff; 
                        border-radius: 8px;
                    ">
                    <p style="font-size: 0.8em; color: #aaa; margin-top: 5px;">Min: 100 GTK</p>
                </div>

                <div id="wagerStatus" style="display: none; margin: 20px 0; padding: 15px; background: rgba(0,0,0,0.3); border-radius: 8px;">
                    <span id="statusText">Searching for opponent...</span>
                    <div class="spinner" style="display: inline-block; width: 20px; height: 20px; border: 3px solid #ffcc00; border-top: 3px solid transparent; border-radius: 50%; animation: spin 1s linear infinite; vertical-align: middle; margin-left: 10px;"></div>
                </div>

                <div style="display: flex; gap: 10px; justify-content: center; margin-top: 20px;">
                    <button id="startWagerBtn" class="premium-btn" style="
                        background: linear-gradient(45deg, #ffcc00, #ffaa00); 
                        color: #000; 
                        font-weight: bold; 
                        padding: 12px 30px; 
                        border: none; 
                        border-radius: 25px; 
                        cursor: pointer; 
                        font-size: 1.1em;
                        transition: transform 0.2s;
                    ">Find Match</button>
                    
                    <button id="cancelWagerBtn" class="premium-btn" style="
                        display: none;
                        background: linear-gradient(45deg, #ff4444, #cc0000); 
                        color: #fff; 
                        font-weight: bold; 
                        padding: 12px 30px; 
                        border: none; 
                        border-radius: 25px; 
                        cursor: pointer; 
                        font-size: 1.1em;
                    ">Cancel Search</button>
                </div>
            </div>
        </div>
    `;

    document.body.appendChild(modal);

    const amountInput = document.getElementById('wagerAmount');
    const startBtn = document.getElementById('startWagerBtn');
    const cancelBtn = document.getElementById('cancelWagerBtn');
    const statusDiv = document.getElementById('wagerStatus');
    const statusText = document.getElementById('statusText');

    let pollInterval;

    startBtn.onclick = async () => {
        const amount = parseInt(amountInput.value);
        if (amount < 100) {
            alert('Minimum wager is 100 GTK');
            return;
        }

        // 1. Start Wager (Lock Funds)
        startBtn.disabled = true;
        startBtn.textContent = 'Locking Funds...';

        try {
            // Get selected team (Assuming current active team for now, ideally user selects)
            const teamRes = await apiClient.get('/teams/active');
            if (!teamRes.team) {
                alert('Please select an active team first!');
                startBtn.disabled = false;
                startBtn.textContent = 'Find Match';
                return;
            }

            const res = await apiClient.post('/battle/wager', {
                team_id: teamRes.team.id, // Assuming team object has id
                wager_amount: amount
            });

            // 2. Update UI to Searching State
            startBtn.style.display = 'none';
            cancelBtn.style.display = 'inline-block';
            statusDiv.style.display = 'block';
            amountInput.disabled = true;

            if (res.status === 'active') {
                // INSTANT MATCH!
                statusText.textContent = 'Match Found! Starting...';
                setTimeout(() => {
                    modal.remove();
                    // Redirect to battle or show battle modal
                    // window.startBattle(res.battle_id); // Assuming this exists global or we emit event
                    // For now, reload or trigger existing battle logic
                    window.location.reload();
                }, 1000);
            } else {
                // QUEUE
                statusText.textContent = 'Funds Locked. Waiting for challenger...';
                // Poll for status
                pollInterval = setInterval(async () => {
                    // We need an endpoint to check battle status specifically?
                    // Or battle/current?
                    /* 
                       Ideally: GET /battle/:id 
                       But we just have res.battle_id
                    */
                    // Simplified polling for MVP: Check active battle via existing logic or new endpoint
                    // Let's assume we can query battle status
                    try {
                        // Re-use battle check logic if available
                        // Logic gap: We need a way to check specific battle ID status
                    } catch (e) { }
                }, 3000);
            }

        } catch (error) {
            console.error(error);
            alert(error.response?.data?.error || 'Failed to start wager');
            startBtn.disabled = false;
            startBtn.textContent = 'Find Match';
        }
    };

    cancelBtn.onclick = async () => {
        if (!confirm('Cancel search and refund funds?')) return;

        cancelBtn.disabled = true;
        cancelBtn.textContent = 'Refunding...';

        try {
            await apiClient.post('/battle/wager/cancel', {});

            clearInterval(pollInterval);
            modal.remove();
            alert('Wager cancelled. Funds refunded.');
            // Refresh balance
            if (window.loadUserBalance) window.loadUserBalance();

        } catch (error) {
            alert('Failed to cancel: ' + error.response?.data?.error);
            cancelBtn.disabled = false;
            cancelBtn.textContent = 'Cancel Search';
        }
    };
}

// Add styles/animations if needed
const style = document.createElement('style');
style.textContent = `
    @keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
`;
document.head.appendChild(style);
