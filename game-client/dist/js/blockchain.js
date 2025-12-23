// Blockchain Configuration & Utilities
export const BLOCKCHAIN_CONFIG = {
    // Default to opBNB Testnet or relevant chain
    rpcUrl: 'https://opbnb-testnet-rpc.bnbchain.org',
    chainId: 5611, // opBNB Testnet
    chainName: 'opBNB Testnet',
    nativeCurrency: {
        name: 'tBNB',
        symbol: 'tBNB',
        decimals: 18
    },
    blockExplorerUrls: ['https://opbnb-testnet.bscscan.com/'],
    contracts: {
        // Placeholder addresses - TO BE UPDATED WITH REAL DEPLOYMENT
        TOWER: '0x66c83Ff8839D4DAD387AD22D48d48605D5C1aDa3',
        GTK: '0xDc4Bf3BbADC42Cc9feE3790B84ca1D532C036588'
    }
};
export async function switchToOpBNB() {
    if (!window.ethereum)
        return;
    try {
        await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: `0x${BLOCKCHAIN_CONFIG.chainId.toString(16)}` }],
        });
    }
    catch (switchError) {
        // This error code indicates that the chain has not been added to MetaMask.
        if (switchError.code === 4902) {
            try {
                await window.ethereum.request({
                    method: 'wallet_addEthereumChain',
                    params: [
                        {
                            chainId: `0x${BLOCKCHAIN_CONFIG.chainId.toString(16)}`,
                            chainName: BLOCKCHAIN_CONFIG.chainName,
                            nativeCurrency: BLOCKCHAIN_CONFIG.nativeCurrency,
                            rpcUrls: [BLOCKCHAIN_CONFIG.rpcUrl],
                            blockExplorerUrls: BLOCKCHAIN_CONFIG.blockExplorerUrls,
                        },
                    ],
                });
            }
            catch (addError) {
                console.error(addError);
            }
        }
    }
}
export async function checkNetwork() {
    if (!window.ethereum)
        return false;
    const chainId = await window.ethereum.request({ method: 'eth_chainId' });
    return parseInt(chainId, 16) === BLOCKCHAIN_CONFIG.chainId;
}
//# sourceMappingURL=blockchain.js.map