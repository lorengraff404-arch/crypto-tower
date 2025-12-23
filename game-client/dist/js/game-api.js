// game-api.ts – wrapper exposing API client as global `gameAPI`
import { APIClient } from "./api.js";
// Initialize the typed API client
const api = new APIClient();
// Define the public interface expected by HTML pages
export const gameAPI = {
    isAuthenticated: () => {
        return !!localStorage.getItem("wallet_address");
    },
    getNonce: async (walletAddress) => {
        return await api.getNonce(walletAddress);
    },
    verifySignature: async (walletAddress, signature) => {
        return await api.verifySignature(walletAddress, signature);
    },
    getProfile: async () => {
        return await api.getProfile();
    },
    getMyCharacters: async () => {
        return await api.getMyCharacters();
    },
    findMatch: async (characterIds, betAmount) => {
        return await api.findMatch(characterIds, betAmount);
    },
    getBattleStatus: async (battleId) => {
        return await api.getBattleStatus(battleId);
    },
    logout: async () => {
        await api.logout();
    }
};
window.gameAPI = gameAPI;
//# sourceMappingURL=game-api.js.map