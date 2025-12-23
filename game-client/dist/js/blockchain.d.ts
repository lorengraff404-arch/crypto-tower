export function switchToOpBNB(): Promise<void>;
export function checkNetwork(): Promise<boolean>;
export namespace BLOCKCHAIN_CONFIG {
    let rpcUrl: string;
    let chainId: number;
    let chainName: string;
    namespace nativeCurrency {
        let name: string;
        let symbol: string;
        let decimals: number;
    }
    let blockExplorerUrls: string[];
    namespace contracts {
        let TOWER: string;
        let GTK: string;
    }
}
//# sourceMappingURL=blockchain.d.ts.map