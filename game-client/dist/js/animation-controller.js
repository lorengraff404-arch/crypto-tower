"use strict";
/**
 * Animation Controller - Manages character animation states and playback
 */
class AnimationController {
    constructor(canvas, ctx) {
        this.canvas = canvas;
        this.ctx = ctx;
        this.animations = new Map(); // character_id -> animation state
        this.frameTime = 1000 / 12; // 12 FPS for pixel art
        this.lastFrameTime = 0;
    }
    /**
     * Register a character for animation
     */
    registerCharacter(characterId, sprites, position) {
        this.animations.set(characterId, {
            sprites: sprites,
            currentAnimation: 'idle',
            currentFrame: 0,
            position: position,
            scale: 2, // Scale factor for rendering
            flipX: false,
            loop: true,
            onComplete: null,
            playing: true
        });
    }
    /**
     * Set animation for a character
     */
    setAnimation(characterId, animationType, options = {}) {
        const state = this.animations.get(characterId);
        if (!state)
            return;
        // Don't restart if already playing this animation
        if (state.currentAnimation === animationType && state.playing) {
            return;
        }
        state.currentAnimation = animationType;
        state.currentFrame = 0;
        state.loop = options.loop !== undefined ? options.loop : this._isLoopAnimation(animationType);
        state.onComplete = options.onComplete || null;
        state.playing = true;
        state.flipX = options.flipX || false;
    }
    /**
     * Check if animation should loop
     */
    _isLoopAnimation(type) {
        const loopAnimations = ['idle', 'walk', 'run', 'block', 'victory'];
        return loopAnimations.includes(type);
    }
    /**
     * Update all animations
     */
    update(timestamp) {
        if (timestamp - this.lastFrameTime < this.frameTime) {
            return; // Not time for next frame yet
        }
        this.lastFrameTime = timestamp;
        for (const [characterId, state] of this.animations.entries()) {
            if (!state.playing)
                continue;
            const sprite = state.sprites[state.currentAnimation];
            if (!sprite)
                continue;
            // Advance frame
            state.currentFrame++;
            // Check if animation complete
            if (state.currentFrame >= sprite.frames.length) {
                if (state.loop) {
                    state.currentFrame = 0; // Loop
                }
                else {
                    state.currentFrame = sprite.frames.length - 1; // Hold last frame
                    state.playing = false;
                    // Call completion callback
                    if (state.onComplete) {
                        state.onComplete();
                        state.onComplete = null;
                    }
                    // Return to idle after non-loop animation
                    setTimeout(() => {
                        this.setAnimation(characterId, 'idle');
                    }, 100);
                }
            }
        }
    }
    /**
     * Render all character animations
     */
    render() {
        for (const [characterId, state] of this.animations.entries()) {
            this.renderCharacter(characterId, state);
        }
    }
    /**
     * Render a single character
     */
    renderCharacter(characterId, state) {
        const sprite = state.sprites[state.currentAnimation];
        if (!sprite || !sprite.image.complete)
            return;
        const frame = sprite.frames[state.currentFrame];
        const pos = state.position;
        this.ctx.save();
        // Apply transformations
        if (state.flipX) {
            this.ctx.scale(-1, 1);
            this.ctx.translate(-pos.x * 2 - frame.width * state.scale, 0);
        }
        // Draw sprite frame
        this.ctx.drawImage(sprite.image, frame.x, frame.y, frame.width, frame.height, // Source
        pos.x, pos.y, frame.width * state.scale, frame.height * state.scale // Destination
        );
        this.ctx.restore();
    }
    /**
     * Play a one-shot animation (attack, skill, etc.)
     */
    playAnimation(characterId, animationType, onComplete) {
        this.setAnimation(characterId, animationType, {
            loop: false,
            onComplete: onComplete
        });
    }
    /**
     * Stop animation for a character
     */
    stopAnimation(characterId) {
        const state = this.animations.get(characterId);
        if (state) {
            state.playing = false;
        }
    }
    /**
     * Resume animation for a character
     */
    resumeAnimation(characterId) {
        const state = this.animations.get(characterId);
        if (state) {
            state.playing = true;
        }
    }
    /**
     * Remove character from animation system
     */
    unregisterCharacter(characterId) {
        this.animations.delete(characterId);
    }
    /**
     * Update character position
     */
    setPosition(characterId, x, y) {
        const state = this.animations.get(characterId);
        if (state) {
            state.position.x = x;
            state.position.y = y;
        }
    }
    /**
     * Set character scale
     */
    setScale(characterId, scale) {
        const state = this.animations.get(characterId);
        if (state) {
            state.scale = scale;
        }
    }
    /**
     * Flip character horizontally
     */
    setFlip(characterId, flipX) {
        const state = this.animations.get(characterId);
        if (state) {
            state.flipX = flipX;
        }
    }
    /**
     * Get current animation state
     */
    getState(characterId) {
        return this.animations.get(characterId);
    }
}
// Export as global
window.AnimationController = AnimationController;
//# sourceMappingURL=animation-controller.js.map