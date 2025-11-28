# Motion Semantics

CyclingStream motion rules keep the UI fast, sharp, energetic, and playful without sacrificing performance or accessibility.

## Timing Tokens

| Token | Duration | Usage |
|-------|----------|-------|
| `instant` | 80ms | Snap indicators, tooltips |
| `fast` | 150ms | Chat entries, hover states |
| `base` | 200ms | Buttons, cards, tabs |
| `slow` | 300ms | Overlay transitions |
| `pulse` | 500ms | Glows, breathing loops |

Durations map directly to CSS variables defined in `app/globals.css` and TypeScript tokens in `constants/design-tokens.ts`.

## Easing

- **Sharp (`cubic-bezier(0.32, 0.72, 0, 1)`):** snappy slide/fade for chat + controls.
- **Spring (`cubic-bezier(0.16, 1, 0.3, 1)`):** playful hover/toggle micro-interactions.
- **Bounce (`cubic-bezier(0.34, 1.56, 0.64, 1)`):** emotes, celebratory states.
- **Linear:** progress indicators, shimmer effects.

## Motion Types

| Semantic | Description | Example |
|----------|-------------|---------|
| **Fast** | Immediate feedback (<200ms) | Chat message entry, hover expansion |
| **Sharp** | Decisive ease-out curves | Navigation focus, modal snap |
| **Energetic** | Springy transitions with overshoot | Button click pulse, toggle knob |
| **Playful** | Bounce/glow combos | Emotes, achievements, VIP chat |

## Primitives / Classes

Defined inside `app/globals.css`:

- `.motion-slide-in-up`: GPU-accelerated slide + fade for chat rows and poll banners.
- `.motion-fade-in`: Lightweight opacity transition for counters/text.
- `.motion-bounce`: 4–8px bounce (emotes).
- `.motion-pulse-ring`: Expanding glow for VIP/MOD badges and live states.
- `.motion-shimmer`: Gradient sweep for progress bars or shimmer states.

All utilities automatically disable under `prefers-reduced-motion`.

## Hooks & Utilities

`frontend/motion/index.ts` exports:

- `motionTokens` – shared durations/easings/z-indices.
- `useMotionPref()` – resolved UI preferences (includes reduced motion overrides).
- `useMotionPreset(name)` – supplies className bundles such as `chat-message-entry`, `vip-ring`, `button-hover`.

## Accessibility

- Honor system `prefers-reduced-motion` by forcing animations off when requested.
- Respect per-user toggles (`ui_preferences`) for chat animations, button pulse, poll motion.
- Emote bounce limited to 4–8px travel, 120–180ms duration.
- Heavy animations (level-ups, overlays) run under 600ms total and never block interaction.

## Performance

- Use `transform` + `opacity` only; avoid layout thrash.
- Reuse GPU-friendly primitives (translate3d, scale).
- Keep animation trees shallow to maintain 60 FPS, especially in chat lists.

## Testing Checklist

1. Verify animations and sounds disable when `Reduced Motion` toggle is on.
2. Measure FPS via Chrome Performance while chat spam or overlay transitions run.
3. Test on TV/mobile breakpoints to ensure animations stay within safe bounds.
4. Validate classes via Storybook or dedicated motion sandbox (TODO).

