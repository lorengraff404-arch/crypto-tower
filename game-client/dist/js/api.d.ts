import type { User } from "./types";
import type { AuthResponse, NonceResponse, HealthResponse } from './types';
export declare class APIClient {
    private baseURL;
    constructor(baseURL?: string);
    setToken(token: string): void;
    clearToken(): void;
    private getHeaders;
    private fetch;
    get<T>(endpoint: string, requiresAuth?: boolean): Promise<T>;
    post<T>(endpoint: string, data?: any, requiresAuth?: boolean): Promise<T>;
    put<T>(endpoint: string, data?: any, requiresAuth?: boolean): Promise<T>;
    delete<T>(endpoint: string, requiresAuth?: boolean): Promise<T>;
    health(): Promise<HealthResponse>;
    getNonce(walletAddress: string): Promise<NonceResponse>;
    verifySignature(walletAddress: string, signature: string): Promise<AuthResponse>;
    getProfile(): Promise<User>;
    getMyCharacters(): Promise<{
        characters: any[];
    }>;
    getCharacter(id: number): Promise<any>;
    findMatch(characterIds: number[], betAmount: number): Promise<any>;
    getBattleStatus(battleId: string): Promise<any>;
    logout(): Promise<void>;
}
export declare const apiClient: APIClient;
export default apiClient;
//# sourceMappingURL=api.d.ts.map