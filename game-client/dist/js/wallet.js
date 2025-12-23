"use strict";
// MetaMask Wallet Integration
class WalletConnector {
    constructor() {
        this.provider = null;
        this.address = null;
        this.chainId = null;
    }
    // Check if MetaMask is installed
    isMetaMaskInstalled() {
        return typeof window.ethereum !== 'undefined';
    }
    // Connect wallet
    async connect() {
        if (!this.isMetaMaskInstalled()) {
            throw new Error('MetaMask is not installed. Please install it from metamask.io');
        }
        try {
            // Request account access
            const accounts = await window.ethereum.request({
                method: 'eth_requestAccounts'
            });
            this.address = accounts[0];
            this.provider = window.ethereum;
            // Get chain ID
            this.chainId = await window.ethereum.request({
                method: 'eth_chainId'
            });
            // Setup event listeners
            this.setupEventListeners();
            return this.address;
        }
        catch (error) {
            console.error('Failed to connect wallet:', error);
            throw error;
        }
    }
    // Setup MetaMask event listeners
    setupEventListeners() {
        window.ethereum.on('accountsChanged', (accounts) => {
            if (accounts.length === 0) {
                this.disconnect();
            }
            else {
                this.address = accounts[0];
                window.location.reload(); // Reload page on account change
            }
        });
        window.ethereum.on('chainChanged', (chainId) => {
            window.location.reload(); // Reload page on network change
        });
    }
    // Sign message (for authentication)
    async signMessage(message) {
        if (!this.address) {
            throw new Error('Wallet not connected');
        }
        try {
            const signature = await window.ethereum.request({
                method: 'personal_sign',
                params: [message, this.address]
            });
            return signature;
        }
        catch (error) {
            console.error('Failed to sign message:', error);
            throw error;
        }
    }
    // Get balance (ETH/BNB)
    async getBalance() {
        if (!this.address) {
            throw new Error('Wallet not connected');
        }
        const balance = await window.ethereum.request({
            method: 'eth_getBalance',
            params: [this.address, 'latest']
        });
        // Convert from wei to ETH
        return parseInt(balance, 16) / 1e18;
    }
    // Switch to BSC network
    async switchToBSC(testnet = false) {
        const chainId = testnet ? '0x61' : '0x38'; // BSC Testnet or Mainnet
        const chainName = testnet ? 'BSC Testnet' : 'BSC Mainnet';
        const rpcUrl = testnet
            ? 'https://data-seed-prebsc-1-s1.binance.org:8545'
            : 'https://bsc-dataseed.binance.org';
        try {
            await window.ethereum.request({
                method: 'wallet_switchEthereumChain',
                params: [{ chainId }]
            });
        }
        catch (switchError) {
            // Network not added, add it
            if (switchError.code === 4902) {
                await window.ethereum.request({
                    method: 'wallet_addEthereumChain',
                    params: [{
                            chainId,
                            chainName,
                            nativeCurrency: {
                                name: 'BNB',
                                symbol: 'BNB',
                                decimals: 18
                            },
                            rpcUrls: [rpcUrl],
                            blockExplorerUrls: [testnet ? 'https://testnet.bscscan.com' : 'https://bscscan.com']
                        }]
                });
            }
            else {
                throw switchError;
            }
        }
    }
    // Disconnect
    disconnect() {
        this.address = null;
        this.provider = null;
        this.chainId = null;
        gameAPI.logout();
    }
    // Get current address
    getAddress() {
        return this.address;
    }
    // Format address for display
    formatAddress(address = this.address) {
        if (!address)
            return '';
        return `${address.substring(0, 6)}...${address.substring(38)}`;
    }
}
// Export for use
const wallet = new WalletConnector();
//# sourceMappingURL=wallet.js.map