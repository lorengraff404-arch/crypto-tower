// Web3 Token Balance Functions
import { BLOCKCHAIN_CONFIG } from './blockchain.js';

// ERC-20 ABI (minimal - just balanceOf function)
const ERC20_ABI = [
    {
        "constant": true,
        "inputs": [{ "name": "_owner", "type": "address" }],
        "name": "balanceOf",
        "outputs": [{ "name": "balance", "type": "uint256" }],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [{ "name": "", "type": "uint8" }],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [{ "name": "", "type": "string" }],
        "type": "function"
    }
];

// Initialize Web3 provider
async function getWeb3() {
    if (window.ethereum) {
        // Ensure wallet is connected
        try {
            await window.ethereum.request({ method: 'eth_requestAccounts' });
        } catch (error) {
            console.warn('Wallet connection request failed, using provider anyway:', error);
        }
        return new Web3(window.ethereum);
    } else if (window.web3) {
        return new Web3(window.web3.currentProvider);
    } else {
        // Fallback to BSC RPC (read-only)
        console.warn('No wallet detected, using read-only RPC provider');
        return new Web3(new Web3.providers.HttpProvider(BLOCKCHAIN_CONFIG.rpcUrl));
    }
}

// Get TOWER token balance
export async function getTowerBalance(walletAddress) {
    try {
        const web3 = await getWeb3();
        const contract = new web3.eth.Contract(ERC20_ABI, BLOCKCHAIN_CONFIG.contracts.TOWER);
        const balance = await contract.methods.balanceOf(walletAddress).call();
        const decimals = await contract.methods.decimals().call();

        // Convert from wei to readable format
        return parseFloat(web3.utils.fromWei(balance, 'ether'));
    } catch (error) {
        console.error('Error fetching TOWER balance:', error);
        return 0;
    }
}

// Get GTK token balance
export async function getGTKBalance(walletAddress) {
    try {
        console.log('🔍 Getting GTK balance for:', walletAddress);
        const web3 = await getWeb3();
        console.log('📡 GTK Contract Address:', BLOCKCHAIN_CONFIG.contracts.GTK);
        const contract = new web3.eth.Contract(ERC20_ABI, BLOCKCHAIN_CONFIG.contracts.GTK);
        const balance = await contract.methods.balanceOf(walletAddress).call();
        console.log('💰 Raw GTK balance (wei):', balance);
        const decimals = await contract.methods.decimals().call();
        console.log('🔢 GTK decimals:', decimals);

        // Convert from wei to readable format
        const readable = parseFloat(web3.utils.fromWei(balance, 'ether'));
        console.log('✅ GTK balance (readable):', readable);
        return readable;
    } catch (error) {
        console.error('❌ Error fetching GTK balance:', error);
        return 0;
    }
}

// Get BNB balance
export async function getBNBBalance(walletAddress) {
    try {
        const web3 = await getWeb3();
        const balance = await web3.eth.getBalance(walletAddress);

        // Convert from wei to BNB
        return parseFloat(web3.utils.fromWei(balance, 'ether'));
    } catch (error) {
        console.error('Error fetching BNB balance:', error);
        return 0;
    }
}

// Get all balances at once
export async function getAllBalances(walletAddress) {
    if (!walletAddress) {
        return { tower: 0, gtk: 0, bnb: 0 };
    }

    try {
        const [tower, gtk, bnb] = await Promise.all([
            getTowerBalance(walletAddress),
            getGTKBalance(walletAddress),
            getBNBBalance(walletAddress)
        ]);

        return {
            tower: tower.toFixed(2),
            gtk: gtk.toFixed(2),
            bnb: bnb.toFixed(4)
        };
    } catch (error) {
        console.error('Error fetching balances:', error);
        return { tower: 0, gtk: 0, bnb: 0 };
    }
}

// Watch for balance changes (polls every 5 seconds)
export function watchBalance(walletAddress, callback) {
    let intervalId = null;

    const updateBalances = async () => {
        const balances = await getAllBalances(walletAddress);
        callback(balances);
    };

    // Initial fetch
    updateBalances();

    // Poll every 5 seconds
    intervalId = setInterval(updateBalances, 5000);

    // Return function to stop watching
    return () => {
        if (intervalId) {
            clearInterval(intervalId);
        }
    };
}
