declare function loadMissions(): Promise<void>;
declare function startMission(): Promise<void>;
declare function completeMission(): Promise<void>;
declare global {
    interface Window {
        startMission: () => Promise<void>;
        completeMission: () => Promise<void>;
        confirmStartMission: () => Promise<void>;
        closeDialogue: () => void;
        closeStartModal: () => void;
        closeCompleteModal: () => void;
    }
}
export { loadMissions, startMission, completeMission };
//# sourceMappingURL=quests.d.ts.map