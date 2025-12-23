/**
 * Animation Controller - Manages character animation states and playback
 */
declare class AnimationController {
    constructor(canvas: any, ctx: any);
    canvas: any;
    ctx: any;
    animations: Map<any, any>;
    frameTime: number;
    lastFrameTime: number;
    /**
     * Register a character for animation
     */
    registerCharacter(characterId: any, sprites: any, position: any): void;
    /**
     * Set animation for a character
     */
    setAnimation(characterId: any, animationType: any, options?: {}): void;
    /**
     * Check if animation should loop
     */
    _isLoopAnimation(type: any): boolean;
    /**
     * Update all animations
     */
    update(timestamp: any): void;
    /**
     * Render all character animations
     */
    render(): void;
    /**
     * Render a single character
     */
    renderCharacter(characterId: any, state: any): void;
    /**
     * Play a one-shot animation (attack, skill, etc.)
     */
    playAnimation(characterId: any, animationType: any, onComplete: any): void;
    /**
     * Stop animation for a character
     */
    stopAnimation(characterId: any): void;
    /**
     * Resume animation for a character
     */
    resumeAnimation(characterId: any): void;
    /**
     * Remove character from animation system
     */
    unregisterCharacter(characterId: any): void;
    /**
     * Update character position
     */
    setPosition(characterId: any, x: any, y: any): void;
    /**
     * Set character scale
     */
    setScale(characterId: any, scale: any): void;
    /**
     * Flip character horizontally
     */
    setFlip(characterId: any, flipX: any): void;
    /**
     * Get current animation state
     */
    getState(characterId: any): any;
}
//# sourceMappingURL=animation-controller.d.ts.map