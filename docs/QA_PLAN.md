# CyclingStream QA & Test Plan

This checklist covers the Motion + Sound + Chat Enhancements scope (Epics 1‑5) plus the shared foundations.

## Automated Verification

| Command | Purpose | Result / Notes |
| --- | --- | --- |
| `npm run lint` (frontend) | ESLint across entire Next.js workspace | **Fails upstream**: existing lint debts (unused vars, `any`, etc.) in unrelated routes/components. No errors triggered by the new motion/sound/chat code. |
| `npx eslint components/video/VideoOverlay.tsx components/video/VideoInfoOverlay.tsx motion/index.ts` | Scoped lint for new brand-motion code | ✅ clean |
| `go test ./...` (backend) | Backend unit tests | **Fails upstream**: handler tests cannot boot Postgres (`pq: role "cyclingstream" does not exist`). Repository/service packages build. |

> When running the global linters/tests, capture the known failures above so regressions are easier to spot. If you have the Postgres role locally you should expect `go test ./...` to pass.

## Manual Test Matrix

| Area | Scenario | Steps | Expected |
| --- | --- | --- | --- |
| **Motion tokens** | Reduced motion override | Settings → disable *Reduced Motion* toggle (and OS `prefers-reduced-motion`). Reload chat/video pages. | All `.motion-*` utilities stop animating; components stay functional. |
| **Chat animations** | Entry + role pulses | Join live race chat (with admin/mod/vip accounts). Send messages with/without special emotes. | Messages slide/fade in; role badges pulse/shimmer; emotes bounce only when `special_emote` true and animations enabled. |
| | Poll announcement | Launch chat poll via admin panel. Vote with viewer accounts. | Poll card slide-in from right; bar widths animate; notification sound fired once per poll create/close. |
| **UI buttons** | Hover + ripple | Navigate to `Settings`, `Chat`, `Watch` screens. Hover/click primary, secondary, icon buttons. Toggle *Button Pulse* off and repeat. | Hover scale + glow when enabled; ripple effect plays on click and is suppressed when disabled or reduced motion on. |
| **Sound cues** | Button click | Toggle preferences ON; click CTA buttons; mute via settings. | `button-click` sound only when enabled; respects master volume. |
| | Notification & mention | Trigger chat poll, level-up, and mention events. | Poll create/close + level-up chime uses `notification` sound (rate-limited); mention ping only when `@username` appears and mention pings enabled. |
| **Achievements** | Unlock journey | 1) Send first chat, 2) accrue 30+ minutes watch time, 3) level up via XP award. | Backend creates `user_achievements` rows; toast queue shows pop-ups sequentially; overlay burst plays once per level-up. Achievements endpoint returns unlocked list. |
| **Video overlays (Brand Motion)** | Live badge & gradients | Open race watch page, toggle controls (mouse move / hide). | Gradient canopy fades with sweep line when controls visible; live badge pops in with glow pulse. |
| | Info ribbon | Ensure race has stage stats; toggle controls. | Card slides up with stat chips popping individually. Reduced-motion disables transitions without layout jumps. |
| **Onboarding** | Completion flow | Finish wizard with valid inputs; ensure `completeOnboarding` succeeds. | Preferences saved, `refreshPreferences()` runs, guard stops redirecting back to `/onboarding`. Reduced motion defaults obey view preference. |
| **Settings persistence** | Toggle UI/audio prefs | Flip each toggle/slider, refresh page. | `ExperienceContext` refetch shares latest state; backend stores JSONB preferences. |

## Regression Checklist

1. **Auth flows** – Login/logout still function since providers wrap layout (`ExperienceProvider`, `SoundProvider`, `AchievementProvider`, `OnboardingGuard` nesting unbroken).
2. **Chat history** – Loading older messages still works; role/badge metadata gracefully handles messages lacking new fields.
3. **Sound fallback** – When Howler fails to preload, logs warning but UI remains silent (no crashes).
4. **WebSocket** – New message types (`poll_created`, `poll_update`, `poll_closed`) ignored safely by older clients.
5. **Backend migrations** – Apply `033/034/035` sequentially on a clean database; verify `achievements` + `user_achievements` tables exist with constraints.

## Follow-ups / Open Risks

- Global ESLint debt should be tackled to allow repo-wide linting to pass.
- Backend tests need a local Postgres role `cyclingstream`. Add instructions to README for running tests inside Docker if preferred.
- Consider Storybook or motion playground for faster visual QA of the motion primitives.

Use this document as the living checklist for regression runs before releases. Update as new epics land.

