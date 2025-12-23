export function getTowerBalance(walletAddress: any): Promise<number>;
export function getGTKBalance(walletAddress: any): Promise<number>;
export function getBNBBalance(walletAddress: any): Promise<number>;
export function getAllBalances(walletAddress: any): Promise<{
    tower: number;
    gtk: number;
    bnb: number;
} | {
    tower: string;
    gtk: string;
    bnb: string;
}>;
export function watchBalance(walletAddress: any, callback: any): () => void;
//# sourceMappingURL=balance.d.ts.map