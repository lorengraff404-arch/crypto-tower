// Logic to connect to Backend API
const API_URL = "http://localhost:8080/api/v1";
// Mock Auth for Demo (Real impl would use JWT from localStorage)
const token = localStorage.getItem('token');

document.getElementById('banForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const target = document.getElementById('banTarget').value;
    const duration = parseInt(document.getElementById('banDuration').value);
    const reason = document.getElementById('banReason').value;

    try {
        const res = await fetch(`${API_URL}/admin/ban`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ target_id: parseInt(target), reason, duration })
        });
        const data = await res.json();
        if (res.ok) {
            alert('User Banned Successfully!');
        } else {
            alert('Error: ' + data.error);
        }
    } catch (err) {
        alert('Network Panic: ' + err.message);
    }
});

// Initial Load
async function loadStats() {
    if (!token) return; // Silent fail if not logged in
    try {
        const res = await fetch(`${API_URL}/revenue/stats`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (res.ok) {
            const data = await res.json();
            document.getElementById('statRevenue').innerText = data.total_revenue || '0';
            // Other stats would be populated here
        }
    } catch (e) { console.error(e); }
}

// Check Auth
if (!token) {
    // Ideally redirect to login, but for dashboard we show message
    console.warn("No token found");
} else {
    loadStats();
}
