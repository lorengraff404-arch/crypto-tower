// game-api.ts – wrapper exposing API client as global `gameAPI`
import { APIClient } from "./api.js";


// Initialize the typed API client
const api = new APIClient();

// Define the public interface expected by HTML pages
export const gameAPI = {
    isAuthenticated: (): boolean => {
        return !!localStorage.getItem("wallet_address");
    },
    getNonce: async (walletAddress: string) => {
        return await api.getNonce(walletAddress);
    },
    verifySignature: async (walletAddress: string, signature: string) => {
        return await api.verifySignature(walletAddress, signature);
    },
    getProfile: async () => {
        return await api.getProfile();
    },
    getMyCharacters: async () => {
        return await api.getMyCharacters();
    },
    findMatch: async (characterIds: number[], betAmount: number) => {
        return await api.findMatch(characterIds, betAmount);
    },
    getBattleStatus: async (battleId: string) => {
        return await api.getBattleStatus(battleId);
    },
    logout: async () => {
        await api.logout();
    }
};


// Expose globally for legacy script usage
declare global {
    interface Window { gameAPI: typeof gameAPI; }
}
window.gameAPI = gameAPI;
