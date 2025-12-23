
// --- Battle Management Logic ---

function showSection(sectionId) {
    // Hide all main sections (simple toggle for now)
    // In a full app, we'd have IDs for the dashboard container too.
    // For now, let's just toggle 'battles-section' and the maindashboard content.
    // Assumption: The main content is inside <div class="col-md-10 p-4"> which lacks an ID.
    // I will add ID to main container via script or assume 'battles-section' is overlay/sibling.
    // To keep it simple:
    const dashboard = document.getElementById('dashboard-content');
    const battles = document.getElementById('battles-section');

    if (sectionId === 'battles') {
        if (dashboard) dashboard.classList.add('d-none');
        battles.classList.remove('d-none');
        loadActiveBattles();
        loadBattleHistory();
    } else {
        if (dashboard) dashboard.classList.remove('d-none');
        battles.classList.add('d-none');
    }
}

async function loadActiveBattles() {
    try {
        const res = await fetch(`${API_URL}/admin/battles/active`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await res.json();
        const tbody = document.getElementById('activeBattlesTable');
        tbody.innerHTML = '';

        if (data.battles && data.battles.length > 0) {
            data.battles.forEach(b => {
                const row = `<tr>
                            <td>${b.id}</td>
                            <td>${b.Player1?.username || b.player1_id} <span class="badge bg-success">${b.player1_bet} GTK</span></td>
                            <td>${b.Player2?.username || b.player2_id} <span class="badge bg-danger">${b.player2_bet} GTK</span></td>
                            <td>${b.battle_type}</td>
                            <td>${new Date(b.created_at).toLocaleTimeString()}</td>
                            <td>
                                <button class="btn btn-sm btn-danger" onclick="terminateBattle(${b.id})">Terminate</button>
                            </td>
                        </tr>`;
                tbody.innerHTML += row;
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="6" class="text-center">No active battles</td></tr>';
        }
    } catch (err) {
        console.error(err);
    }
}

async function loadBattleHistory() {
    try {
        const res = await fetch(`${API_URL}/admin/battles/history?limit=10`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await res.json();
        const tbody = document.getElementById('historyBattlesTable');
        tbody.innerHTML = '';

        if (data.battles) {
            data.battles.forEach(b => {
                const winner = b.Winner?.username || (b.winner_id ? `ID ${b.winner_id}` : 'None');
                const row = `<tr>
                            <td>${b.id}</td>
                            <td class="text-success fw-bold">${winner}</td>
                            <td>${b.WinnerID === b.player1_id ? b.Player2?.username : b.Player1?.username}</td>
                            <td>${b.battle_type}</td>
                            <td>${new Date(b.ended_at || b.updated_at).toLocaleString()}</td>
                            <td><small>${b.status}</small></td>
                        </tr>`;
                tbody.innerHTML += row;
            });
        }
    } catch (err) {
        console.error(err);
    }
}

async function terminateBattle(battleId) {
    const reason = prompt("Enter functionality reason for termination:");
    if (!reason) return;

    try {
        const res = await fetch(`${API_URL}/admin/battles/${battleId}/terminate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ reason })
        });

        if (res.ok) {
            alert("Battle terminated.");
            loadActiveBattles();
        } else {
            const data = await res.json();
            alert("Error: " + data.error);
        }
    } catch (err) {
        alert("Network error");
    }
}

