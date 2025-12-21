// Character and Element Emoji Mapping
const CHARACTER_EMOJIS = {
    Dragon: '🐉',
    Beast: '🦁',
    Insect: '🦋',
    Mineral: '💎',
    Spirit: '👻',
    Avian: '🦅',
    Aqua: '🐠',
    Flora: '🌺'
};

const ELEMENT_EMOJIS = {
    Fire: '🔥',
    Water: '💧',
    Ice: '❄️',
    Thunder: '⚡',
    Dark: '🌑',
    Plant: '🌿',
    Earth: '🪨',
    Wind: '💨'
};

const RARITY_CONFIG = {
    SSS: {
        gradient: 'linear-gradient(135deg, #FFD700, #FFA500, #FF6347)',
        stars: '⭐⭐⭐⭐⭐',
        glow: '0 0 20px rgba(255, 215, 0, 0.6)'
    },
    SS: {
        gradient: 'linear-gradient(135deg, #C0C0C0, #E8E8E8, #B0B0B0)',
        stars: '⭐⭐⭐⭐',
        glow: '0 0 15px rgba(192, 192, 192, 0.5)'
    },
    S: {
        gradient: 'linear-gradient(135deg, #CD7F32, #A0522D)',
        stars: '⭐⭐⭐',
        glow: '0 0 10px rgba(205, 127, 50, 0.4)'
    },
    A: {
        gradient: 'linear-gradient(135deg, #4CAF50, #45a049)',
        stars: '⭐⭐',
        glow: 'none'
    },
    B: {
        gradient: 'linear-gradient(135deg, #2196F3, #1976D2)',
        stars: '⭐',
        glow: 'none'
    },
    C: {
        gradient: 'linear-gradient(135deg, #9E9E9E, #757575)',
        stars: '★',
        glow: 'none'
    }
};

const BATTLE_EMOJIS = {
    attack: '⚔️',
    defend: '🛡️',
    heal: '💚',
    buff: '✨',
    debuff: '💢',
    critical: '💥',
    miss: '💨',
    victory: '🏆',
    defeat: '💀'
};

// Helper functions
function getCharacterEmoji(type) {
    return CHARACTER_EMOJIS[type] || '❓';
}

function getElementEmoji(element) {
    return ELEMENT_EMOJIS[element] || '❓';
}

function getRarityConfig(rarity) {
    return RARITY_CONFIG[rarity] || RARITY_CONFIG.C;
}

function getBattleEmoji(action) {
    return BATTLE_EMOJIS[action] || '❓';
}
