export declare const gameAPI: {
    isAuthenticated: () => boolean;
    getNonce: (walletAddress: string) => Promise<import("./types.js").NonceResponse>;
    verifySignature: (walletAddress: string, signature: string) => Promise<import("./types.js").AuthResponse>;
    getProfile: () => Promise<import("./types.js").User>;
    getMyCharacters: () => Promise<{
        characters: any[];
    }>;
    findMatch: (characterIds: number[], betAmount: number) => Promise<any>;
    getBattleStatus: (battleId: string) => Promise<any>;
    logout: () => Promise<void>;
};
declare global {
    interface Window {
        gameAPI: typeof gameAPI;
    }
}
//# sourceMappingURL=game-api.d.ts.map