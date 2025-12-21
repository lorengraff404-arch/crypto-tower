/**
 * Sprite Loader - Handles loading and caching of character sprite sheets
 */

class SpriteLoader {
    constructor() {
        this.cache = new Map(); // Cache loaded sprites
        this.loading = new Map(); // Track loading promises
    }

    /**
     * Load all sprite sheets for a character
     * @param {Object} character - Character object with sprite URLs
     * @returns {Promise<Object>} - Object with loaded sprite images
     */
    async loadCharacterSprites(character) {
        const cacheKey = `character_${character.id}`;

        // Return from cache if available
        if (this.cache.has(cacheKey)) {
            return this.cache.get(cacheKey);
        }

        // Wait if already loading
        if (this.loading.has(cacheKey)) {
            return await this.loading.get(cacheKey);
        }

        // Load sprites
        const loadPromise = this._loadSprites(character);
        this.loading.set(cacheKey, loadPromise);

        try {
            const sprites = await loadPromise;
            this.cache.set(cacheKey, sprites);
            return sprites;
        } finally {
            this.loading.delete(cacheKey);
        }
    }

    /**
     * Internal method to load all sprite sheets
     */
    async _loadSprites(character) {
        const animationTypes = [
            'idle', 'walk', 'run', 'attack', 'skill',
            'hit', 'block', 'dodge', 'death', 'victory'
        ];

        const sprites = {};
        const loadPromises = [];

        for (const type of animationTypes) {
            const spriteKey = `sprite_${type}`;
            const spriteURL = character[spriteKey];

            if (spriteURL) {
                loadPromises.push(
                    this._loadSpriteSheet(spriteURL, type)
                        .then(spriteData => {
                            sprites[type] = spriteData;
                        })
                        .catch(err => {
                            console.warn(`Failed to load ${type} sprite:`, err);
                            sprites[type] = this._getPlaceholderSprite(type);
                        })
                );
            } else {
                // Use placeholder if sprite not available
                sprites[type] = this._getPlaceholderSprite(type);
            }
        }

        await Promise.all(loadPromises);
        return sprites;
    }

    /**
     * Load a single sprite sheet
     */
    async _loadSpriteSheet(url, type) {
        return new Promise((resolve, reject) => {
            const img = new Image();
            img.crossOrigin = 'anonymous';

            img.onload = () => {
                // Parse sprite sheet into frames
                const frames = this._parseSpriteSheet(img, type);
                resolve({
                    image: img,
                    frames: frames,
                    type: type
                });
            };

            img.onerror = () => {
                reject(new Error(`Failed to load sprite: ${url}`));
            };

            img.src = url;
        });
    }

    /**
     * Parse sprite sheet into individual frames
     */
    _parseSpriteSheet(image, type) {
        // Get frame count for this animation type
        const frameCount = this._getFrameCount(type);
        const frames = [];

        // Assuming horizontal sprite sheet layout
        const frameWidth = image.width / frameCount;
        const frameHeight = image.height;

        for (let i = 0; i < frameCount; i++) {
            frames.push({
                x: i * frameWidth,
                y: 0,
                width: frameWidth,
                height: frameHeight
            });
        }

        return frames;
    }

    /**
     * Get frame count for animation type
     */
    _getFrameCount(type) {
        const frameCounts = {
            'idle': 8,
            'walk': 8,
            'run': 8,
            'attack': 10,
            'skill': 10,
            'hit': 6,
            'block': 6,
            'dodge': 8,
            'death': 10,
            'victory': 10
        };
        return frameCounts[type] || 8;
    }

    /**
     * Get placeholder sprite for missing animations
     */
    _getPlaceholderSprite(type) {
        // Create a simple colored rectangle as placeholder
        const canvas = document.createElement('canvas');
        canvas.width = 512;
        canvas.height = 64;
        const ctx = canvas.getContext('2d');

        // Draw placeholder
        ctx.fillStyle = '#4A90E2';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = '#FFFFFF';
        ctx.font = '16px Arial';
        ctx.textAlign = 'center';
        ctx.fillText(`${type.toUpperCase()} (Loading...)`, canvas.width / 2, canvas.height / 2);

        const img = new Image();
        img.src = canvas.toDataURL();

        return {
            image: img,
            frames: this._parseSpriteSheet(img, type),
            type: type,
            isPlaceholder: true
        };
    }

    /**
     * Clear cache for a character
     */
    clearCache(characterId) {
        const cacheKey = `character_${characterId}`;
        this.cache.delete(cacheKey);
    }

    /**
     * Clear all cache
     */
    clearAllCache() {
        this.cache.clear();
    }
}

// Export as global
window.SpriteLoader = SpriteLoader;
