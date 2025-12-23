import { BLOCKCHAIN_CONFIG } from './blockchain.js';

// Treasury Address (Deployer from .env)
const TREASURY_ADDRESS = "0xdCb8ca66Ae0809Eed5dB73E8e1c3787c8178327e"; 

// Minimal ERC20 ABI
const ERC20_ABI = [
    {
        "constant": false,
        "inputs": [
            { "name": "_to", "type": "address" },
            { "name": "_value", "type": "uint256" }
        ],
        "name": "transfer",
        "outputs": [{ "name": "", "type": "bool" }],
        "type": "function"
    }
];

export async function transferGTK(amount) {
    return transferToken(BLOCKCHAIN_CONFIG.contracts.GTK, amount);
}

export async function transferTOWER(amount) {
    return transferToken(BLOCKCHAIN_CONFIG.contracts.TOWER, amount);
}

async function transferToken(contractAddress, amount) {
    if (!window.ethereum) throw new Error("MetaMask not installed");
    
    // Check if Web3 is loaded
    if (typeof Web3 === 'undefined') throw new Error("Web3.js not loaded. Please ensure web3.min.js is included.");

    const web3 = new Web3(window.ethereum);
    const accounts = await web3.eth.requestAccounts();
    const sender = accounts[0];
    
    const contract = new web3.eth.Contract(ERC20_ABI, contractAddress);
    
    // Convert amount to wei (18 decimals)
    // Note: 'amount' is passed as a number/string representing tokens (e.g. 50 GTK)
    const amountWei = web3.utils.toWei(amount.toString(), 'ether');
    
    console.log(`Transferring ${amount} (${amountWei} wei) of ${contractAddress} to ${TREASURY_ADDRESS}`);

    try {
        const tx = await contract.methods.transfer(TREASURY_ADDRESS, amountWei).send({ from: sender });
        console.log("Transaction successful, hash:", tx.transactionHash);
        // Return hash whether it's in a receipt object or direct
        return tx.transactionHash || tx;
    } catch (err) {
        console.error("Transfer failed", err);
        throw err;
    }
}

export const waitForTransaction = async (txHash) => {
    if (!window.ethereum) throw new Error("MetaMask not installed");
    const web3 = new Web3(window.ethereum);
    
    console.log("Waiting for transaction:", txHash);
    let receipt = null;
    while (receipt === null) {
        try {
            receipt = await web3.eth.getTransactionReceipt(txHash);
            if (receipt === null) {
                await new Promise(resolve => setTimeout(resolve, 2000));
            }
        } catch (e) {
             console.warn("Error fetching receipt, retrying...", e);
             await new Promise(resolve => setTimeout(resolve, 2000));
        }
    }
    return receipt;
};
