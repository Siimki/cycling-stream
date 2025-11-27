/**
 * Color palettes and user color definitions
 * 
 * Note: USER_COLORS uses hex format instead of OKLCH for the following reasons:
 * 1. User-generated content: These colors are used for chat usernames and may need
 *    to be displayed in external systems or exported
 * 2. Compatibility: Hex colors are universally supported across all browsers and
 *    external integrations
 * 3. Simplicity: These are static color values that don't need the advanced features
 *    of OKLCH (perceptual uniformity, color manipulation)
 * 4. Performance: Hex colors are slightly more performant for inline styles
 * 
 * The core design system uses OKLCH (see globals.css), but user-facing colors
 * can remain in hex format for practical reasons.
 */

export const USER_COLORS = [
  '#22d3ee', // Cyan
  '#a78bfa', // Purple
  '#fb7185', // Pink
  '#4ade80', // Green
  '#fbbf24', // Yellow
  '#f472b6', // Rose
  '#60a5fa', // Blue
  '#34d399', // Emerald
];

