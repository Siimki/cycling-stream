# UI Sound Guide

CyclingStream uses lightweight, tactile sound cues to reinforce key actions without overwhelming the stream audio mix.

## Sound Events

| Event | File | Default Volume | Rate Limit |
|-------|------|----------------|------------|
| Button click | `public/sounds/click-soft.wav` | 25% | none |
| Notification (poll/event/reward) | `public/sounds/notify-chime.wav` | 35% | 5 per minute |
| Chat mention ping | `public/sounds/chat-mention.wav` | 20% | 3 per 30s |
| Level-up | `public/sounds/level-up.wav` | 40% | 1 per 2 min |

All sounds live under `frontend/public/sounds/` and are loaded lazily through the `SoundProvider`.

## Implementation

- `frontend/lib/sound/sound-manager.ts` – Minimal Web Audio wrapper with async loading, global volume, and rate limiting.
- `frontend/components/providers/SoundProvider.tsx` – React context that ties user audio preferences to `play(soundId)`.
- `useSound()` – Hook to request playback anywhere in the client tree.

## Preferences

Audio toggles are stored in `user_preferences.audio_preferences`:

- `button_clicks`
- `notification_sounds`
- `mention_pings`
- `master_volume` (0–1)

Users can edit these settings in `/settings`. Preferences sync with the backend and respect `prefers-reduced-motion` where appropriate.

## Asset Workflow

1. Source sounds (Mixkit, Zapsplat, Freesound, etc.).
2. Trim + normalize (-12dB target) in Audacity.
3. Export short mono WAV (<200ms) to keep payload tiny.
4. Update `SOUND_MANIFEST` if filenames change.

## Testing Checklist

- ✅ Verify sounds mute immediately when toggles are off.
- ✅ Confirm rate limit logic (no more than 5 notifications/minute).
- ✅ Test laptop speakers, TV, and headphones for clarity at 10–20% volume.
- ✅ Measure bundle impact (all files combined <30 KB).

