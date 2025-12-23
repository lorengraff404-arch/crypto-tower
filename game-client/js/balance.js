// Web3 Token Balance Functions - AUTO NETWORK SWITCH
import { BLOCKCHAIN_CONFIG, switchToOpBNB } from './blockchain.js';

const ERC20_ABI = [{
    "constant": true,
    "inputs": [{"name": "_owner", "type": "address"}],
    "name": "balanceOf",
    "outputs": [{"name": "balance", "type": "uint256"}],
    "type": "function"
}];

async function getWeb3() {
    // ALWAYS use opBNB RPC (no wallet needed for read-only)
    console.log('📡 Using opBNB Testnet RPC:', BLOCKCHAIN_CONFIG.rpcUrl);
    return new Web3(new Web3.providers.HttpProvider(BLOCKCHAIN_CONFIG.rpcUrl));
}

export async function getTowerBalance(walletAddress) {
    try {
        const web3 = await getWeb3();
        const contract = new web3.eth.Contract(ERC20_ABI, BLOCKCHAIN_CONFIG.contracts.TOWER);
        const balance = await contract.methods.balanceOf(walletAddress).call();
        return parseFloat(web3.utils.fromWei(balance, 'ether'));
    } catch (error) {
        console.error('TOWER balance error:', error.message);
        return 0;
    }
}

export async function getGTKBalance(walletAddress) {
    try {
        const web3 = await getWeb3();
        const contract = new web3.eth.Contract(ERC20_ABI, BLOCKCHAIN_CONFIG.contracts.GTK);
        const balance = await contract.methods.balanceOf(walletAddress).call();
        return parseFloat(web3.utils.fromWei(balance, 'ether'));
    } catch (error) {
        console.error('GTK balance error:', error.message);
        return 0;
    }
}

export async function getBNBBalance(walletAddress) {
    try {
        const web3 = await getWeb3();
        const balance = await web3.eth.getBalance(walletAddress);
        return parseFloat(web3.utils.fromWei(balance, 'ether'));
    } catch (error) {
        console.error('BNB balance error:', error.message);
        return 0;
    }
}

export async function getAllBalances(walletAddress) {
    if (!walletAddress) return { tower: 0, gtk: 0, bnb: 0 };
    
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
}

export function watchBalance(walletAddress, callback) {
    let intervalId = null;
    const updateBalances = async () => {
        const balances = await getAllBalances(walletAddress);
        callback(balances);
    };
    updateBalances();
    intervalId = setInterval(updateBalances, 5000);
    return () => { if (intervalId) clearInterval(intervalId); };
}
