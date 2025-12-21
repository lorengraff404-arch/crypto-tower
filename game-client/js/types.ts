// Type definitions for API responses and models

// ============= User & Auth =============
export interface User {
    id: number;
    wallet_address: string;
    level: number;
    experience: number;
    gtk_balance: number;
    tower_balance: number;
    unlocked_features: string[];
    created_at: string;
    updated_at: string;
}

export interface AuthResponse {
    token: string;
    user: User;
}

export interface NonceResponse {
    nonce: string;
}

// ============= Islands =============
export interface Island {
    id: number;
    name: string;
    biome: string;
    difficulty: 'beginner' | 'advanced' | 'expert' | 'legendary';
    entry_fee: number;
    base_reward: number;
    wave_count: number;
    drop_rates: Record<string, number>;
    available_enemies: string;
    created_at: string;
}

export interface Enemy {
    id: number;
    name: string;
    family: string;
    type: string;
    element: string;
    base_stats: {
        attack: number;
        defense: number;
        hp: number;
        speed: number;
    };
    abilities: string[];
    loot: Record<string, number>;
    is_boss: boolean;
}

export interface IslandRaid {
    id: number;
    user_id: number;
    island_id: number;
    status: 'IN_PROGRESS' | 'COMPLETED' | 'FAILED' | 'EXPIRED' | 'ABANDONED';
    current_stage: number;
    total_stages: number;
    current_boss_hp: number;
    current_team_hp: number;
    initial_team_hp?: number;
    turn_count: number;
    character_states?: string; // JSON string in DB, might be parsed in API? No backend sends it as string or interface{}? Backend sends *models.RaidSession. 
    // Backend definition: CharacterStates string `json:"character_states,omitempty"` 
    // Actually, if it's a JSON string in DB, the Go struct has string. 
    // But Gorm might not auto-parse it unless using a hook or generic scanner.
    // Let's assume it comes as string for now, or if I changed it to interface{}? 
    // Looking at backend models/raid.go: 	CharacterStates   string `gorm:"type:text" json:"character_states,omitempty"`
    // So it sends a string. Frontend needs to parse it.
    rewards_claimed: boolean;
    created_at: string;
    updated_at: string;
    // Computed/Legacy support
    waves_cleared?: number; // kept for compatibility if needed, but backend uses current_stage
}
export interface IslandListResponse {
    islands: Island[];
    raids_today: number;
}

export interface IslandDetailResponse {
    island: Island;
    enemies: Enemy[];
}

export interface RaidCompleteResponse {
    raid: IslandRaid;
    rewards: {
        gtk?: number;
        items?: Array<{
            type: string;
            name: string;
            quantity: number;
        }>;
    };
}

// ============= Missions =============
export interface Mission {
    id: number;
    level: number;
    name: string;
    description: string;
    story: string;
    mission_type: 'tutorial' | 'progression' | 'daily' | 'weekly';
    unlock_feature?: string;
    objectives: MissionObjective[];
    rewards: MissionRewards;
    status: 'locked' | 'available' | 'in_progress' | 'completed';
    objectives_progress?: ObjectiveProgress[];
}

export interface MissionObjective {
    type: string;
    description: string;
    target: number;
}

export interface ObjectiveProgress {
    current: number;
    completed: boolean;
}

export interface MissionRewards {
    gtk?: number;
    experience?: number;
    materials?: number;
    items?: Array<{
        type: string;
        quantity: number;
    }>;
}

export interface MissionListResponse {
    missions: Mission[];
    current_mission?: Mission;
}

export interface MissionCompleteResponse {
    mission: Mission;
    rewards: MissionRewards;
    unlock_feature?: string;
    level_up?: boolean;
    new_level?: number;
}

// ============= Story/Narrative =============
export interface StoryDialogue {
    id: number;
    mission_level: number;
    dialogue_type: 'briefing' | 'post_mission' | 'cutscene';
    character: 'aria' | 'kairos' | 'voice' | 'narrator';
    dialogue_text: string;
    sort_order: number;
}

export interface DialogueResponse {
    mission_level: number;
    briefings: StoryDialogue[];
    post_missions: StoryDialogue[];
    cutscenes: StoryDialogue[];
}

export interface StoryProgress {
    user_id: number;
    current_act: number;
    aria_corruption: number;
    kairos_relationship: 'neutral' | 'friend' | 'enemy';
    voice_encounters: number;
    collected_fragments: string[];
    viewed_cutscenes: string[];
}

export interface StoryFragment {
    id: number;
    fragment_id: string;
    title: string;
    content: string;
    fragment_type: string;
    unlock_level: number;
    rarity: string;
}

// ============= API Error Response =============
export interface APIError {
    error: string;
    details?: string;
}

// ============= Generic API Response =============
export interface APIResponse<T> {
    data?: T;
    error?: string;
    message?: string;
}

// ============= Pagination =============
export interface PaginatedResponse<T> {
    data: T[];
    page: number;
    page_size: number;
    total_count: number;
    total_pages: number;
}

// ============= Health Check =============
export interface HealthResponse {
    status: 'healthy' | 'unhealthy';
    service: string;
    version: string;
    database: string;
}
