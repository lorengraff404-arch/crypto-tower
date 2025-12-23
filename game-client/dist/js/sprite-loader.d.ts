/**
 * Sprite Loader - Handles loading and caching of character sprite sheets
 */
declare class SpriteLoader {
    cache: Map<any, any>;
    loading: Map<any, any>;
    /**
     * Load all sprite sheets for a character
     * @param {Object} character - Character object with sprite URLs
     * @returns {Promise<Object>} - Object with loaded sprite images
     */
    loadCharacterSprites(character: Object): Promise<Object>;
    /**
     * Internal method to load all sprite sheets
     */
    _loadSprites(character: any): Promise<{}>;
    /**
     * Load a single sprite sheet
     */
    _loadSpriteSheet(url: any, type: any): Promise<any>;
    /**
     * Parse sprite sheet into individual frames
     */
    _parseSpriteSheet(image: any, type: any): {
        x: number;
        y: number;
        width: number;
        height: any;
    }[];
    /**
     * Get frame count for animation type
     */
    _getFrameCount(type: any): any;
    /**
     * Get placeholder sprite for missing animations
     */
    _getPlaceholderSprite(type: any): {
        image: HTMLImageElement;
        frames: {
            x: number;
            y: number;
            width: number;
            height: any;
        }[];
        type: any;
        isPlaceholder: boolean;
    };
    /**
     * Clear cache for a character
     */
    clearCache(characterId: any): void;
    /**
     * Clear all cache
     */
    clearAllCache(): void;
}
//# sourceMappingURL=sprite-loader.d.ts.map