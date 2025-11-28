'use client';

export type SoundId = 'button-click' | 'notification' | 'chat-mention' | 'level-up';

export interface PlaySoundOptions {
  volume?: number;
  masterVolume?: number;
  rateLimit?: {
    windowMs: number;
    max: number;
  };
}

interface SoundDefinition {
  src: string;
  volume: number;
  rateLimit?: {
    windowMs: number;
    max: number;
  };
}

type BrowserAudioWindow = Window & {
  webkitAudioContext?: typeof AudioContext;
};

const SOUND_MANIFEST: Record<SoundId, SoundDefinition> = {
  'button-click': { src: '/sounds/click-soft.wav', volume: 0.25 },
  notification: { src: '/sounds/notify-chime.wav', volume: 0.35, rateLimit: { windowMs: 60_000, max: 5 } },
  'chat-mention': { src: '/sounds/chat-mention.wav', volume: 0.2, rateLimit: { windowMs: 30_000, max: 3 } },
  'level-up': { src: '/sounds/level-up.wav', volume: 0.4, rateLimit: { windowMs: 120_000, max: 1 } },
};

export class SoundManager {
  private manifest: Record<SoundId, SoundDefinition>;
  private buffers = new Map<SoundId, AudioBuffer>();
  private loading = new Map<SoundId, Promise<AudioBuffer | null>>();
  private plays = new Map<SoundId, number[]>();
  private audioContext?: AudioContext;

  constructor(manifest: Record<SoundId, SoundDefinition>) {
    this.manifest = manifest;
  }

  private isBrowser() {
    return typeof window !== 'undefined';
  }

  private getContext() {
    if (!this.isBrowser()) {
      return undefined;
    }

    if (!this.audioContext) {
      const audioWindow = window as BrowserAudioWindow;
      const AudioCtx = audioWindow.AudioContext || audioWindow.webkitAudioContext;
      if (AudioCtx) {
        this.audioContext = new AudioCtx();
      }
    }

    return this.audioContext;
  }

  private async loadBuffer(id: SoundId) {
    if (this.buffers.has(id)) {
      return this.buffers.get(id) ?? null;
    }

    if (this.loading.has(id)) {
      return this.loading.get(id) ?? null;
    }

    if (!this.isBrowser()) {
      return null;
    }

    const definition = this.manifest[id];
    if (!definition) {
      return null;
    }

    const loadPromise = (async () => {
      try {
        const response = await fetch(definition.src);
        const arrayBuffer = await response.arrayBuffer();
        const context = this.getContext();
        if (!context) {
          return null;
        }
        const audioBuffer = await context.decodeAudioData(arrayBuffer.slice(0));
        this.buffers.set(id, audioBuffer);
        return audioBuffer;
      } catch (error) {
        console.error(`Failed to load sound "${id}"`, error);
        return null;
      } finally {
        this.loading.delete(id);
      }
    })();

    this.loading.set(id, loadPromise);
    return loadPromise;
  }

  async preload(ids: SoundId[]) {
    await Promise.all(ids.map((id) => this.loadBuffer(id)));
  }

  private canPlay(soundId: SoundId, limiter?: { windowMs: number; max: number }) {
    if (!limiter) {
      return true;
    }

    const now = Date.now();
    const windowStart = now - limiter.windowMs;
    const history = (this.plays.get(soundId) ?? []).filter((timestamp) => timestamp >= windowStart);

    if (history.length >= limiter.max) {
      this.plays.set(soundId, history);
      return false;
    }

    history.push(now);
    this.plays.set(soundId, history);
    return true;
  }

  async play(soundId: SoundId, options?: PlaySoundOptions) {
    if (!this.isBrowser()) {
      return;
    }

    const definition = this.manifest[soundId];
    if (!definition) {
      return;
    }

    const limiter = options?.rateLimit ?? definition.rateLimit;
    if (!this.canPlay(soundId, limiter)) {
      return;
    }

    const context = this.getContext();
    if (!context) {
      return;
    }

    const buffer = await this.loadBuffer(soundId);
    if (!buffer) {
      return;
    }

    if (context.state === 'suspended') {
      await context.resume();
    }

    const source = context.createBufferSource();
    source.buffer = buffer;

    const gainNode = context.createGain();
    const baseVolume = options?.volume ?? definition.volume;
    const masterVolume = options?.masterVolume ?? 1;
    const gainValue = Math.min(1, Math.max(0, baseVolume * masterVolume));
    gainNode.gain.value = gainValue;

    source.connect(gainNode);
    gainNode.connect(context.destination);
    source.start(0);
  }
}

export const soundManager = new SoundManager(SOUND_MANIFEST);

