declare class WalletConnector {
    provider: any;
    address: any;
    chainId: any;
    isMetaMaskInstalled(): boolean;
    connect(): Promise<any>;
    setupEventListeners(): void;
    signMessage(message: any): Promise<any>;
    getBalance(): Promise<number>;
    switchToBSC(testnet?: boolean): Promise<void>;
    disconnect(): void;
    getAddress(): any;
    formatAddress(address?: any): string;
}
declare const wallet: WalletConnector;
//# sourceMappingURL=wallet.d.ts.map