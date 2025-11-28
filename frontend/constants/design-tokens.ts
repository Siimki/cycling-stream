/**
 * Design Tokens System
 * 
 * Centralized design tokens for consistent UI/UX across the application.
 * These tokens define colors, typography, spacing, and layout values.
 */

// Color Tokens
// Using OKLCH for perceptual uniformity and better color manipulation
export const colors = {
  // Background colors
  background: 'oklch(0.05 0 0)', // #050505 - Main background
  surface: 'oklch(0.08 0 0)',    // #0A0A0A - Elevated surfaces (cards, panels)
  
  // Text colors
  textPrimary: 'oklch(0.95 0 0)',   // #F5F5F5 - Primary text
  textSecondary: 'oklch(0.65 0 0)', // #A3A3A3 - Secondary text, labels
  
  // Accent (using existing primary green)
  accent: 'oklch(0.72 0.19 155)', // Existing primary green for interactive elements
} as const;

// Typography Scale
export const typography = {
  sizes: {
    xs: '12px',   // Labels, captions
    sm: '14px',   // Small text, secondary info
    base: '16px', // Body text, default
    lg: '20px',   // Subheadings
    xl: '24px',   // Headings
  },
  weights: {
    normal: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
  },
  lineHeights: {
    tight: 1.2,
    normal: 1.5,
    relaxed: 1.6,
  },
} as const;

// Spacing System (8px base unit)
export const spacing = {
  xs: '8px',   // Between related elements
  sm: '16px',  // Between sections
  md: '24px',  // Between major sections
  lg: '32px',  // Page gutters (desktop)
  xl: '24px',  // Page gutters (mobile)
} as const;

// Layout Tokens
export const layout = {
  gridColumns: 12,
  gutterMobile: '24px',
  gutterDesktop: '32px',
  playerColumns: {
    desktop: 8,
    large: 9,
  },
  chatColumns: {
    desktop: 4,
    large: 3,
  },
} as const;

// Border Radius
export const radius = {
  sm: '4px',
  md: '8px',
  lg: '10px',
  xl: '12px',
} as const;

// Motion Tokens
export const motion = {
  durations: {
    instant: 80,
    fast: 150,
    base: 200,
    slow: 300,
    pulse: 500,
  },
  easing: {
    sharp: 'cubic-bezier(0.32, 0.72, 0, 1)',
    spring: 'cubic-bezier(0.16, 1, 0.3, 1)',
    bounce: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
    linear: 'linear',
  },
  zIndex: {
    chatGlow: 15,
    overlay: 40,
    toast: 45,
    levelUp: 60,
  },
} as const;

