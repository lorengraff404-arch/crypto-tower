// Debug script to check wallet and balance state
// Run this in browser console when on game.html

console.log('=== WALLET DEBUG INFO ===');

// 1. Check localStorage
const walletAddress = localStorage.getItem('wallet_address');
console.log('1. Wallet from localStorage:', walletAddress);

// 2. Check if MetaMask is available
console.log('2. MetaMask available:', typeof window.ethereum !== 'undefined');

// 3. Try to get accounts from MetaMask
if (window.ethereum) {
    window.ethereum.request({ method: 'eth_accounts' })
        .then(accounts => {
            console.log('3. Accounts from eth_accounts:', accounts);
            if (accounts.length === 0) {
                console.warn('   ⚠️ No accounts returned - wallet not connected!');
            }
        })
        .catch(err => console.error('   Error getting accounts:', err));
}

// 4. Check if Web3 is loaded
console.log('4. Web3 loaded:', typeof Web3 !== 'undefined');

// 5. Check balance display elements
const gtkBalanceHeader = document.getElementById('gtkBalance');
const shopGTKBalance = document.getElementById('shopGTKBalance');
console.log('5. GTK Balance (header):', gtkBalanceHeader?.textContent);
console.log('   GTK Balance (shop):', shopGTKBalance?.textContent);

// 6. Try to manually get GTK balance
if (window.ethereum && walletAddress) {
    import('./dist/js/balance.js').then(async ({ getGTKBalance }) => {
        try {
            const balance = await getGTKBalance(walletAddress);
            console.log('6. Manual GTK balance fetch:', balance);
        } catch (error) {
            console.error('   Error fetching balance:', error);
        }
    });
}

console.log('=== END DEBUG INFO ===');
