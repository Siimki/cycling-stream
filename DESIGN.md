# CyclingStream Design System

> **Version 1.0** | Last Updated: 2025  
> This document serves as the single source of truth for design decisions and ensures consistency across the CyclingStream platform.

---

## Table of Contents

1. [Design Philosophy](#design-philosophy)
2. [Color System](#color-system)
3. [Typography](#typography)
4. [Spacing & Layout](#spacing--layout)
5. [Components](#components)
6. [Effects & Animations](#effects--animations)
7. [Responsive Design](#responsive-design)
8. [Accessibility](#accessibility)
9. [Best Practices](#best-practices)
10. [Inconsistencies & Future Improvements](#inconsistencies--future-improvements)

---

## Design Philosophy

### Vision

CyclingStream is a modern, professional streaming platform designed for cycling enthusiasts. Our design philosophy centers on:

- **Performance First**: Fast, responsive, and optimized for live streaming
- **Dark Theme Excellence**: Inspired by modern streaming platforms (Kick/Twitch) with a cycling-focused aesthetic
- **Accessibility**: WCAG 2.1 AA compliant, ensuring the platform is usable by everyone
- **Consistency**: Unified design language across all components and pages
- **Mobile-First**: Responsive design that works beautifully on all devices

### Design Principles

1. **Clarity Over Cleverness**: Simple, intuitive interfaces that don't require explanation
2. **Energy & Motion**: Reflects the dynamic nature of cycling with subtle animations and transitions
3. **Content First**: Design supports the content (streams, races, chat) without competing for attention
4. **Progressive Enhancement**: Core functionality works everywhere, enhanced experiences where supported

### Inspiration

- **Streaming Platforms**: Kick, Twitch (dark themes, chat integration, live indicators)
- **Cycling Culture**: Green primary color represents energy, nature, and cycling
- **Modern Web**: OKLCH color space, CSS custom properties, Tailwind CSS

---

## Color System

### Color Space: OKLCH

We use **OKLCH** (OK Lightness Chroma Hue) color space for all colors. OKLCH provides:
- Perceptually uniform color gradients
- Better color manipulation (lightness, saturation, hue)
- Future-proof color system
- Consistent appearance across devices

### Color Palette

#### Primary Colors

```css
--primary: oklch(0.72 0.19 155);              /* Green - Main brand color */
--primary-foreground: oklch(0.1 0.02 155);    /* Dark green for text on primary */
```

**Usage**: Primary actions, links, highlights, brand elements, live indicators

#### Background Colors

```css
--background: oklch(0.11 0.005 260);          /* Main background - very dark */
--card: oklch(0.14 0.005 260);                /* Card/surface background */
--popover: oklch(0.14 0.005 260);             /* Popover/dropdown background */
--sidebar: oklch(0.12 0.005 260);             /* Sidebar background */
```

**Usage**:
- `background`: Main page background
- `card`: Cards, modals, elevated surfaces
- `popover`: Dropdowns, tooltips, popovers
- `sidebar`: Sidebar/navigation panels

#### Foreground Colors

```css
--foreground: oklch(0.95 0 0);                /* Primary text color */
--card-foreground: oklch(0.95 0 0);           /* Text on cards */
--popover-foreground: oklch(0.95 0 0);        /* Text on popovers */
--sidebar-foreground: oklch(0.95 0 0);        /* Text on sidebar */
```

**Usage**: All text content, icons, borders (when used as foreground)

#### Secondary Colors

```css
--secondary: oklch(0.18 0.005 260);           /* Secondary background */
--secondary-foreground: oklch(0.95 0 0);      /* Text on secondary */
```

**Usage**: Secondary buttons, less prominent UI elements

#### Muted Colors

```css
--muted: oklch(0.22 0.005 260);               /* Muted background */
--muted-foreground: oklch(0.55 0 0);          /* Muted text */
```

**Usage**: 
- `muted`: Subtle backgrounds, disabled states
- `muted-foreground`: Secondary text, placeholders, labels

#### Accent Colors

```css
--accent: oklch(0.72 0.19 155);               /* Accent color (same as primary) */
--accent-foreground: oklch(0.1 0.02 155);     /* Text on accent */
```

**Usage**: Hover states, focus indicators, subtle highlights

#### Destructive Colors

```css
--destructive: oklch(0.577 0.245 27.325);     /* Error/danger color */
--destructive-foreground: oklch(0.98 0 0);    /* Text on destructive */
```

**Usage**: Error messages, delete actions, destructive buttons

#### Status Colors

```css
--success: oklch(0.65 0.2 150);              /* Success/green */
--success-foreground: oklch(0.1 0.02 150);    /* Text on success */
--error: oklch(0.577 0.245 27.325);          /* Error/red (same as destructive in :root) */
--error-foreground: oklch(0.98 0 0);         /* Text on error */
--warning: oklch(0.75 0.15 85);              /* Warning/yellow */
--warning-foreground: oklch(0.1 0.02 85);    /* Text on warning */
--info: oklch(0.65 0.15 220);                /* Info/blue */
--info-foreground: oklch(0.1 0.02 220);      /* Text on info */
```

**Usage**: 
- `success`: Success messages, positive feedback, completed states
- `error`: Error states (can use `destructive` as alternative)
- `warning`: Warning messages, caution states
- `info`: Informational messages, neutral states

#### Connection & Live Status Colors

```css
--live: oklch(0.577 0.245 27.325);           /* Live indicator red */
--connected: oklch(0.65 0.2 150);            /* Connected green */
--disconnected: oklch(0.577 0.245 27.325);   /* Disconnected red */
```

**Usage**:
- `live`: Live stream indicators, "ON AIR" badges
- `connected`: Connection status (connected state)
- `disconnected`: Connection status (disconnected state)

#### Border & Input Colors

```css
--border: oklch(0.24 0.005 260);              /* Border color */
--input: oklch(0.18 0.005 260);               /* Input background */
--ring: oklch(0.72 0.19 155);                 /* Focus ring color */
```

**Usage**:
- `border`: All borders, dividers
- `input`: Input field backgrounds
- `ring`: Focus indicators, active states

#### Chart Colors

```css
--chart-1: oklch(0.72 0.19 155);              /* Primary chart color */
--chart-2: oklch(0.7 0.15 45);                /* Secondary chart color */
--chart-3: oklch(0.6 0.15 220);               /* Tertiary chart color */
--chart-4: oklch(0.7 0.18 280);               /* Quaternary chart color */
--chart-5: oklch(0.65 0.2 330);               /* Quinary chart color */
```

**Usage**: Data visualization, charts, graphs

### Color Usage Guidelines

#### Opacity Modifiers

Use Tailwind opacity modifiers consistently:

- `/10` - Very subtle backgrounds (error backgrounds, hover states)
- `/20` - Subtle backgrounds (badges, tags)
- `/30` - Light backgrounds (hover states, subtle overlays)
- `/40` - Medium backgrounds (borders with opacity)
- `/50` - Semi-transparent (backdrop overlays, dividers)
- `/70` - Mostly opaque (gradient stops)
- `/80` - Nearly opaque (card backgrounds with blur)
- `/90` - Almost fully opaque (hover states)

**Examples**:
```tsx
// ✅ Correct
<div className="bg-card/80 backdrop-blur-sm">
<div className="bg-primary/20 text-primary">
<div className="border-primary/40">

// ❌ Avoid
<div className="bg-card opacity-80">  // Use /80 instead
```

#### Semantic Color Usage

- **Primary**: Actions, links, active states, brand elements
- **Destructive**: Errors, delete actions, warnings
- **Success**: Success messages, positive feedback
- **Error**: Error states (alternative to destructive)
- **Warning**: Warning messages, caution states
- **Info**: Informational messages
- **Live/Connected/Disconnected**: Status indicators
- **Muted**: Secondary information, disabled states
- **Foreground**: All text content
- **Background/Card**: Surfaces and containers

#### Color Accessibility

All color combinations meet WCAG 2.1 AA contrast requirements:
- Text on background: 4.5:1 minimum
- Large text (18px+): 3:1 minimum
- Interactive elements: Clear focus states

---

## Typography

### Font Families

#### Primary Font: Geist Sans

```css
--font-sans: var(--font-geist-sans), "Geist", sans-serif;
```

**Usage**: All body text, headings, UI elements, navigation

**Characteristics**:
- Modern, clean, highly readable
- Excellent for screens
- Optimized for performance (Next.js font optimization)

#### Monospace Font: Geist Mono

```css
--font-mono: "JetBrains Mono", "JetBrains Mono Fallback", monospace;
```

**Usage**: Code snippets, technical data, timestamps, IDs

### Type Scale

#### Headings

| Element | Size | Weight | Line Height | Usage |
|---------|------|--------|-------------|-------|
| `h1` | `text-4xl` (2.25rem / 36px) | `font-bold` (700) | `leading-tight` | Page titles |
| `h2` | `text-3xl` (1.875rem / 30px) | `font-bold` (700) | `leading-tight` | Section titles |
| `h3` | `text-2xl` (1.5rem / 24px) | `font-semibold` (600) | `leading-snug` | Subsection titles |
| `h4` | `text-xl` (1.25rem / 20px) | `font-semibold` (600) | `leading-snug` | Card titles |
| `h5` | `text-lg` (1.125rem / 18px) | `font-semibold` (600) | `leading-normal` | Small headings |
| `h6` | `text-base` (1rem / 16px) | `font-semibold` (600) | `leading-normal` | Labels |

#### Body Text

| Element | Size | Weight | Line Height | Usage |
|---------|------|--------|-------------|-------|
| Body | `text-base` (1rem / 16px) | `font-normal` (400) | `leading-relaxed` | Main content |
| Small | `text-sm` (0.875rem / 14px) | `font-normal` (400) | `leading-normal` | Secondary text |
| Extra Small | `text-xs` (0.75rem / 12px) | `font-medium` (500) | `leading-normal` | Labels, captions |

#### Special Text

| Element | Size | Weight | Usage |
|---------|------|--------|-------|
| Large | `text-lg` (1.125rem / 18px) | `font-normal` | Emphasized body text |
| Display | `text-6xl` (3.75rem / 60px) | `font-bold` | Hero text, large displays |

### Typography Guidelines

#### Text Colors

```tsx
// Primary text
<p className="text-foreground">Main content</p>

// Secondary text
<p className="text-muted-foreground">Secondary information</p>

// Muted text (lighter)
<p className="text-muted-foreground/70">Tertiary information</p>

// Primary colored text
<span className="text-primary">Highlighted text</span>
```

#### Text Opacity

Use opacity modifiers for text hierarchy:

- `text-foreground/95` - Primary text (slight transparency for depth)
- `text-foreground/90` - Secondary text
- `text-muted-foreground` - Tertiary text (default: 55% opacity)
- `text-muted-foreground/60` - Very subtle text

#### Font Weights

- `font-normal` (400) - Body text, default
- `font-medium` (500) - Emphasized text, labels
- `font-semibold` (600) - Headings, important text
- `font-bold` (700) - Strong emphasis, page titles

#### Text Alignment

- **Left**: Default for most content
- **Center**: Hero sections, empty states, error pages
- **Right**: Numbers, timestamps (rare)

#### Text Truncation

```tsx
// Single line truncation
<p className="truncate">Long text that will be truncated</p>

// Multi-line truncation (2 lines)
<p className="line-clamp-2">Long text that will be truncated after 2 lines</p>
```

---

## Spacing & Layout

### Spacing Scale

We use Tailwind's default spacing scale (4px base unit):

| Token | Value | Usage |
|-------|-------|-------|
| `0` | 0px | No spacing |
| `0.5` | 2px | Tight spacing, icons |
| `1` | 4px | Very tight spacing |
| `1.5` | 6px | Tight spacing |
| `2` | 8px | Small spacing |
| `2.5` | 10px | Small-medium spacing |
| `3` | 12px | Medium spacing |
| `4` | 16px | Default spacing |
| `5` | 20px | Medium-large spacing |
| `6` | 24px | Large spacing |
| `8` | 32px | Extra large spacing |
| `10` | 40px | Section spacing |
| `12` | 48px | Large section spacing |
| `16` | 64px | Extra large section spacing |

### Spacing Guidelines

#### Padding

- **Cards**: `p-4 sm:p-6` (16px mobile, 24px desktop)
- **Sections**: `py-6 sm:py-8` (24px/32px vertical)
- **Buttons**: `px-4 py-2` (16px horizontal, 8px vertical)
- **Inputs**: `px-3 py-2` (12px horizontal, 8px vertical)
- **Icons**: `p-2` or `p-1.5` (8px or 6px)

#### Margins

- **Between elements**: `mb-4` or `gap-4` (16px)
- **Section spacing**: `mt-8` or `mb-8` (32px)
- **Component spacing**: `gap-2` to `gap-6` (8px to 24px)

#### Gaps (Flexbox/Grid)

- **Tight**: `gap-1` or `gap-2` (4px-8px)
- **Default**: `gap-4` (16px)
- **Loose**: `gap-6` or `gap-8` (24px-32px)

### Container Widths

| Container | Max Width | Usage |
|-----------|-----------|-------|
| Full | `w-full` | Full width |
| Small | `max-w-md` (28rem / 448px) | Forms, modals |
| Medium | `max-w-4xl` (56rem / 896px) | Content pages |
| Large | `max-w-7xl` (80rem / 1280px) | Homepage, race listings |

### Layout Patterns

#### Page Layout

```tsx
<div className="min-h-screen bg-background flex flex-col">
  <Navigation variant="full" />
  <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 w-full">
    {/* Content */}
  </main>
  <Footer />
</div>
```

#### Card Layout

```tsx
<div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6">
  {/* Card content */}
</div>
```

#### Grid Layouts

```tsx
// Responsive grid
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
  {/* Grid items */}
</div>
```

#### Flex Layouts

```tsx
// Horizontal flex
<div className="flex items-center justify-between gap-4">
  {/* Flex items */}
</div>

// Vertical flex
<div className="flex flex-col gap-4">
  {/* Flex items */}
</div>
```

---

## Components

### Button

#### Variants

```tsx
// Primary (default)
<Button>Primary Action</Button>

// Destructive
<Button variant="destructive">Delete</Button>

// Outline
<Button variant="outline">Secondary Action</Button>

// Secondary
<Button variant="secondary">Tertiary Action</Button>

// Ghost
<Button variant="ghost">Subtle Action</Button>

// Link
<Button variant="link">Link Style</Button>
```

#### Sizes

```tsx
<Button size="sm">Small</Button>
<Button size="default">Default</Button>
<Button size="lg">Large</Button>
<Button size="icon">Icon Only</Button>
```

#### Primary Button Pattern

For important actions (sign up, watch race, etc.):

```tsx
<Button className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-semibold">
  Watch Race
</Button>
```

**Usage**: CTA buttons, primary actions, important links

### Input

```tsx
<Input
  type="email"
  placeholder="Enter email"
  className="bg-background"
/>
```

**Styling**:
- Height: `h-10` (40px)
- Padding: `px-3 py-2`
- Border: `border border-input`
- Focus: `focus-visible:ring-2 focus-visible:ring-ring`

### Card

#### Standard Card

```tsx
<div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6">
  <h3 className="text-xl font-semibold text-foreground/95 mb-2">Card Title</h3>
  <p className="text-muted-foreground">Card content</p>
</div>
```

#### Interactive Card (RaceCard pattern)

```tsx
<Link href="/races/123">
  <div className="bg-card/80 backdrop-blur-sm border-2 border-primary/40 rounded-lg p-4 sm:p-6 hover:border-primary hover:border-[3px] hover:bg-card/90 hover:shadow-[0_4px_20px_rgba(0,0,0,0.3)] hover:shadow-primary/20 hover:-translate-y-1 transition-all duration-200 cursor-pointer h-full flex flex-col group">
    <h3 className="text-lg sm:text-xl font-semibold text-foreground/95 mb-2 group-hover:text-primary transition-colors">
      Card Title
    </h3>
    {/* Content */}
  </div>
</Link>
```

### Navigation

#### Full Navigation

```tsx
<Navigation variant="full" />
```

**Features**:
- Sticky header: `sticky top-0 z-50`
- Background: `bg-card/95 backdrop-blur-xl`
- Border: `border-b border-border/30`
- Height: `h-14 sm:h-16` (56px mobile, 64px desktop)

#### Minimal Navigation

```tsx
<Navigation variant="minimal" />
```

**Usage**: Auth pages, minimal layouts

### Chat Component

```tsx
<div className="flex flex-col h-full bg-card/50 border-l border-border min-h-0">
  {/* Header */}
  <div className="px-4 py-3 border-b border-border/50 flex items-center justify-between shrink-0 h-12">
    {/* Chat header content */}
  </div>
  
  {/* Messages */}
  <div className="flex-1 overflow-y-auto chat-scroll px-4 py-3 min-h-0">
    {/* Messages */}
  </div>
  
  {/* Input */}
  <div className="p-3 border-t border-border/50 shrink-0 bg-card/30">
    {/* Input field */}
  </div>
</div>
```

**Features**:
- Custom scrollbar (`.chat-scroll`)
- Fixed header and footer
- Scrollable message area

### Video Player

```tsx
<div className="aspect-video w-full h-full bg-black rounded-lg overflow-hidden border border-border">
  {/* Video element */}
</div>
```

**Features**:
- 16:9 aspect ratio (`aspect-video`)
- Black background
- Rounded corners
- Border for definition

### Error Message

```tsx
// Default variant
<ErrorMessage message="Error message" />

// Inline variant
<ErrorMessage variant="inline" message="Error message" />

// Full variant
<ErrorMessage variant="full" message="Error message" onRetry={handleRetry} />
```

**Styling**:
- Background: `bg-destructive/10`
- Border: `border border-destructive/20`
- Text: `text-destructive`

### Loading States

#### Loading Spinner

```tsx
// Standard spinner
<div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>

// Small spinner
<div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
```

**Note**: Current `LoadingSpinner.tsx` uses hardcoded colors - should be updated to use design tokens.

#### Skeleton Loader

```tsx
<div className="animate-pulse bg-muted rounded">
  {/* Skeleton content */}
</div>
```

**Note**: Current `SkeletonLoader.tsx` uses hardcoded colors - should be updated to use design tokens.

### Mission Card

#### Standard Mission Card

```tsx
<div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6">
  <div className="grid grid-cols-1 sm:grid-cols-[1fr_auto] gap-3 sm:gap-4">
    {/* Left: Title + Description */}
    <div className="flex-1 min-w-0">
      <h3 className="text-base sm:text-lg font-semibold text-foreground/95 mb-1 truncate">
        Mission Title
      </h3>
      {description && (
        <p className="text-sm text-muted-foreground truncate">{description}</p>
      )}
    </div>

    {/* Right: Reward */}
    <div className="text-right sm:flex sm:items-start sm:justify-end">
      <div>
        <div className="text-lg sm:text-xl font-bold text-primary">
          +{points_reward}
        </div>
        <div className="text-xs text-muted-foreground">points</div>
      </div>
    </div>
  </div>

  {/* Middle: Progress Bar */}
  <div className="mt-3">
    <MissionProgress progress={progress} target={target_value} />
  </div>

  {/* Status/Claim Button */}
  <div className="flex items-center justify-between mt-3">
    {/* Status or claim button */}
  </div>
</div>
```

**Design Rules**:
- **Layout**: Grid-based with left (title+description), middle (progress), right (reward)
- **Padding**: `p-4 sm:p-6` (consistent across all cards)
- **Border**: `border border-border/50` with `rounded-lg`
- **Title**: `text-base sm:text-lg font-semibold text-foreground/95` with `truncate`
- **Description**: `text-sm text-muted-foreground` with `truncate` (one line only)
- **Reward**: Right-aligned, `text-lg sm:text-xl font-bold text-primary`
- **Status Labels**: `text-sm text-muted-foreground` or `text-sm text-primary font-medium` (for completed)
- **Spacing**: `gap-3 sm:gap-4` between grid items, `mt-3` between sections

### Mission Progress Bar

#### Standard Progress Bar

```tsx
<div className="w-full">
  <div className="flex items-center justify-between mb-1">
    <span className="text-xs text-muted-foreground">
      {progress} / {target}
    </span>
    <span className="text-xs text-muted-foreground">{Math.round(percentage)}%</span>
  </div>
  <div className="w-full bg-muted rounded-full h-2 overflow-hidden">
    <div
      className="h-full bg-primary"
      style={{ width: `${percentage}%` }}
    />
  </div>
</div>
```

**Design Rules**:
- **Height**: `h-2` (8px) - consistent across all progress bars
- **Border Radius**: `rounded-full`
- **Fill Color**: `bg-primary` (brand-green)
- **Background**: `bg-muted` (neutral dark)
- **Percentage Text**: Right-aligned, `text-xs text-muted-foreground`
- **Progress Text**: Left-aligned, `text-xs text-muted-foreground`
- **No transitions or animations** for this release (static UI only)

### Weekly Overview Card

#### Combined Weekly Overview

```tsx
<div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4 sm:p-6">
  {/* Level & XP Progress - Main progress bar */}
  <div className="mb-4">
    <div className="flex items-center justify-between mb-2">
      <h3 className="text-lg font-semibold text-foreground/95">
        Level {level}
      </h3>
      <div className="text-sm text-muted-foreground">
        {xp_total} XP
      </div>
    </div>
    <div className="space-y-1">
      <div className="flex justify-between text-sm text-muted-foreground">
        <span>Progress to Level {level + 1}</span>
        <span>{progress} / {target} XP</span>
      </div>
      <div className="h-3 bg-muted rounded-full overflow-hidden">
        <div
          className="h-full bg-primary transition-all duration-300"
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  </div>

  {/* Weekly Goals - Simple 2-row blocks */}
  <div className="space-y-3 pt-3 border-t border-border/50">
    <div>
      <div className="flex justify-between text-sm mb-1">
        <span className="text-muted-foreground">Watch Time</span>
        <span className="font-medium text-foreground">{minutes} / 30 min</span>
      </div>
    </div>
    <div>
      <div className="flex justify-between text-sm mb-1">
        <span className="text-muted-foreground">Chat Messages</span>
        <span className="font-medium text-foreground">{messages} / 3</span>
      </div>
    </div>
  </div>
</div>
```

**Design Rules**:
- **Single cohesive card** (not multiple cards)
- **Main progress bar**: XP progress with `h-3` height, `bg-primary` fill
- **Weekly goals**: Simple 2-row blocks (no progress bars, just numbers)
- **Numbers**: Clearly visible with `font-medium text-foreground`
- **Labels**: Muted with `text-muted-foreground`
- **Spacing**: `mb-4` for main section, `pt-3` with `border-t border-border/50` for weekly goals
- **Typography**: Follow design system scale

### Mission Grouping

#### Grouped Missions Display

```tsx
<div className="space-y-6">
  {groupedMissions.map((group) => (
    <div key={group.title} className="space-y-4">
      <h2 className="text-sm font-medium text-muted-foreground uppercase tracking-wider">
        {group.title}
      </h2>
      <div className="space-y-4">
        {group.missions.map((mission) => (
          <MissionCard key={mission.id} userMission={mission} />
        ))}
      </div>
    </div>
  ))}
</div>
```

**Grouping Rules**:
- **Weekly Missions**: `watch_time` mission types
- **Chat Missions**: `chat_message` mission types
- **Prediction Missions**: `predict_winner` mission types
- **Special / Limited-time**: `watch_race`, `follow_series`, `streak` mission types
- **Section Headers**: `text-sm font-medium text-muted-foreground uppercase tracking-wider`
- **Spacing**: `space-y-6` between groups, `space-y-4` between cards within a group
- **Do NOT intermix** mission types - each group contains only its designated types

### Streak Display

#### Streak Card

```tsx
<div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-4">
  <h3 className="text-lg font-semibold text-foreground/95 mb-2">Streak</h3>
  <div className="flex items-center gap-4">
    <div>
      <div className="flex items-center gap-2">
        <Link className={`w-5 h-5 ${active ? 'text-primary' : 'text-muted-foreground'}`} />
        <span className={`text-2xl font-bold ${active ? 'text-primary' : 'text-foreground'}`}>
          {current_streak}
        </span>
      </div>
      <div className="text-sm text-muted-foreground">Current Streak</div>
    </div>
    <div>
      <div className="text-2xl font-bold text-foreground">
        {best_streak}
      </div>
      <div className="text-sm text-muted-foreground">Best Streak</div>
    </div>
  </div>
</div>
```

**Design Rules**:
- **Icon**: Use `Link` icon from lucide-react (no emojis)
- **Icon Size**: `w-5 h-5` (consistent with design system)
- **Active Color**: `text-primary` (brand-green only when streak > 0)
- **Inactive Color**: `text-muted-foreground` or `text-foreground`
- **Compact Layout**: Minimal spacing, clear hierarchy
- **No emojis**: Ever. Use icons from lucide-react only.

---

## Effects & Animations

### Border Radius

#### Standard Values

| Token | Value | Usage |
|-------|-------|-------|
| `rounded-sm` | `calc(var(--radius) - 4px)` | Small elements |
| `rounded-md` | `calc(var(--radius) - 2px)` | Default (buttons, inputs) |
| `rounded-lg` | `var(--radius)` (10px) | Cards, containers |
| `rounded-xl` | `calc(var(--radius) + 4px)` | Large cards |
| `rounded-full` | `9999px` | Pills, avatars, icons |

**Base Radius**: `--radius: 0.625rem` (10px)

#### Usage Guidelines

- **Buttons**: `rounded-md` or `rounded-lg`
- **Inputs**: `rounded-md`
- **Cards**: `rounded-lg`
- **Badges/Tags**: `rounded-md` or `rounded-full`
- **Avatars**: `rounded-full`
- **Icons in containers**: `rounded-lg` or `rounded-full`

### Shadows

#### Standard Shadows

```tsx
// Subtle shadow
className="shadow-sm"

// Default shadow
className="shadow-md"

// Large shadow
className="shadow-lg"

// Custom shadow (for cards)
className="shadow-[0_4px_20px_rgba(0,0,0,0.3)]"

// Primary glow
className="shadow-lg shadow-primary/20"
```

#### Glow Effect

Custom glow utility class:

```css
.glow-primary {
  box-shadow: 0 0 20px oklch(0.72 0.19 155 / 0.15);
}
```

**Usage**: Important elements, active states, highlights

### Borders

#### Border Styles

```tsx
// Standard border
className="border border-border"

// Border with opacity
className="border border-border/50"

// Primary border
className="border-2 border-primary/40"

// Thick border on hover
className="hover:border-[3px] hover:border-primary"
```

#### Border Guidelines

- **Default**: `border border-border` (1px)
- **Emphasized**: `border-2 border-primary/40` (2px)
- **Hover**: `hover:border-[3px]` (3px)
- **Opacity**: Use `/50` or `/40` for subtle borders

### Backdrop Blur

#### Standard Values

```tsx
// Small blur (cards)
className="backdrop-blur-sm"

// Large blur (navigation)
className="backdrop-blur-xl"
```

#### Usage Guidelines

- **Cards**: `backdrop-blur-sm` (subtle blur)
- **Navigation**: `backdrop-blur-xl` (strong blur for glass effect)
- **Modals**: `backdrop-blur-md` or `backdrop-blur-lg`

**Note**: Always combine with semi-transparent background:
```tsx
className="bg-card/80 backdrop-blur-sm"
```

### Transitions

#### Standard Transitions

```tsx
// Color transitions (most common)
className="transition-colors"

// All properties
className="transition-all duration-200"

// Specific duration
className="transition-all duration-300 ease-in-out"
```

#### Transition Durations

- **Fast**: `duration-150` (150ms) - Hover states
- **Default**: `duration-200` (200ms) - Most interactions
- **Slow**: `duration-300` (300ms) - Complex animations

#### Common Transition Patterns

```tsx
// Hover color change
className="hover:text-primary transition-colors"

// Hover background
className="hover:bg-muted/50 transition-colors"

// Transform + shadow
className="hover:-translate-y-1 hover:shadow-lg transition-all duration-200"

// Opacity
className="hover:opacity-90 transition-opacity"
```

### Animations

#### Built-in Animations

```tsx
// Spin (loading)
className="animate-spin"

// Pulse (loading, live indicators)
className="animate-pulse"

// Custom live pulse
className="animate-live-pulse"
```

#### Custom Animations

**Live Pulse** (for live indicators):

```css
@keyframes live-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.animate-live-pulse {
  animation: live-pulse 2s ease-in-out infinite;
}
```

**Usage**:
```tsx
<span className="w-2 h-2 rounded-full bg-red-500 animate-live-pulse" />
```

---

## Responsive Design

### Breakpoints

We use Tailwind's default breakpoints:

| Breakpoint | Min Width | Usage |
|------------|-----------|-------|
| `sm` | 640px | Small tablets, large phones |
| `md` | 768px | Tablets |
| `lg` | 1024px | Small desktops |
| `xl` | 1280px | Desktops |
| `2xl` | 1536px | Large desktops |

### Mobile-First Approach

Always design for mobile first, then enhance for larger screens:

```tsx
// Mobile-first pattern
<div className="
  text-base          // Mobile: 16px
  sm:text-lg         // Small+: 18px
  lg:text-xl         // Large+: 20px
">
  Responsive text
</div>
```

### Common Responsive Patterns

#### Padding

```tsx
className="px-4 sm:px-6 lg:px-8"  // 16px → 24px → 32px
className="py-6 sm:py-8"          // 24px → 32px
```

#### Grid Layouts

```tsx
// 1 column mobile, 2 tablet, 3 desktop
className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6"
```

#### Typography

```tsx
className="text-lg sm:text-xl"           // 18px → 20px
className="text-2xl sm:text-3xl"         // 24px → 30px
```

#### Spacing

```tsx
className="gap-2 sm:gap-4"               // 8px → 16px
className="mb-4 sm:mb-6"                 // 16px → 24px
```

#### Flex Direction

```tsx
className="flex flex-col sm:flex-row"    // Column → Row
```

#### Visibility

```tsx
className="hidden sm:block"              // Hidden mobile, visible desktop
className="block sm:hidden"              // Visible mobile, hidden desktop
```

### Container Patterns

```tsx
// Standard page container
<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
  {/* Content */}
</div>

// Content container
<div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
  {/* Content */}
</div>
```

---

## Accessibility

### WCAG 2.1 AA Compliance

All design decisions prioritize accessibility:

- **Color Contrast**: Minimum 4.5:1 for text, 3:1 for large text
- **Focus States**: Clear, visible focus indicators
- **Keyboard Navigation**: All interactive elements are keyboard accessible
- **Screen Readers**: Proper semantic HTML and ARIA labels

### Focus States

#### Standard Focus Ring

```tsx
className="focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
```

**Usage**: All interactive elements (buttons, inputs, links)

#### Focus Colors

- **Ring Color**: `--ring: oklch(0.72 0.19 155)` (primary green)
- **Offset**: `ring-offset-2` (8px) for visibility

### Color Contrast

All text/background combinations meet WCAG AA:

- **Foreground on Background**: 4.5:1 ✅
- **Muted Foreground on Background**: 4.5:1 ✅
- **Primary on Primary Foreground**: 4.5:1 ✅
- **Destructive on Destructive Foreground**: 4.5:1 ✅

### Semantic HTML

```tsx
// ✅ Correct
<button>Click me</button>
<a href="/page">Link</a>
<nav>Navigation</nav>

// ❌ Avoid
<div onClick={handleClick}>Click me</div>
<div role="link">Link</div>
```

### ARIA Labels

```tsx
// Icon buttons
<button aria-label="Close menu">
  <X className="w-5 h-5" />
</button>

// Loading states
<div aria-live="polite" aria-busy="true">
  Loading...
</div>
```

### Keyboard Navigation

- **Tab Order**: Logical, sequential
- **Enter/Space**: Activate buttons and links
- **Escape**: Close modals, dropdowns
- **Arrow Keys**: Navigate lists, menus

---

## Best Practices

### Do's ✅

1. **Use Design Tokens**: Always use CSS variables and Tailwind classes
   ```tsx
   // ✅ Correct
   <div className="bg-card text-foreground">
   
   // ❌ Avoid
   <div className="bg-[#1a1a1a] text-white">
   ```

2. **Consistent Spacing**: Use the spacing scale
   ```tsx
   // ✅ Correct
   <div className="p-4 gap-4">
   
   // ❌ Avoid
   <div className="p-[17px] gap-[18px]">
   ```

3. **Responsive Design**: Mobile-first approach
   ```tsx
   // ✅ Correct
   <div className="text-base sm:text-lg">
   
   // ❌ Avoid
   <div className="text-lg sm:text-base">  // Desktop-first
   ```

4. **Semantic Colors**: Use semantic color names
   ```tsx
   // ✅ Correct
   <div className="bg-destructive text-destructive-foreground">
   
   // ❌ Avoid
   <div className="bg-red-500 text-white">
   ```

5. **Consistent Transitions**: Use standard transition patterns
   ```tsx
   // ✅ Correct
   <button className="hover:bg-primary transition-colors">
   
   // ❌ Avoid
   <button className="hover:bg-primary transition-all duration-500">
   ```

### Don'ts ❌

1. **Hardcoded Colors**: Don't use hex/rgb directly
   ```tsx
   // ❌ Avoid
   <div className="bg-[#1a1a1a]">
   <div style={{ backgroundColor: '#1a1a1a' }}>
   ```

2. **Inconsistent Opacity**: Don't mix opacity methods
   ```tsx
   // ❌ Avoid
   <div className="bg-card opacity-80">  // Use /80 instead
   ```

3. **Magic Numbers**: Don't use arbitrary values
   ```tsx
   // ❌ Avoid
   <div className="p-[13px] gap-[7px]">
   ```

4. **Inline Styles**: Avoid inline styles (except dynamic values)
   ```tsx
   // ❌ Avoid
   <div style={{ padding: '16px', backgroundColor: '#1a1a1a' }}>
   
   // ✅ Acceptable (dynamic values)
   <div style={{ width: `${percentage}%` }}>
   ```

5. **Breaking Responsive**: Don't use fixed widths without breakpoints
   ```tsx
   // ❌ Avoid
   <div className="w-[500px]">
   
   // ✅ Correct
   <div className="w-full max-w-md">
   ```

### Common Patterns

#### Card with Hover Effect

```tsx
<div className="
  bg-card/80 backdrop-blur-sm
  border-2 border-primary/40
  rounded-lg
  p-4 sm:p-6
  hover:border-primary
  hover:border-[3px]
  hover:bg-card/90
  hover:shadow-[0_4px_20px_rgba(0,0,0,0.3)]
  hover:shadow-primary/20
  hover:-translate-y-1
  transition-all duration-200
  cursor-pointer
">
  {/* Content */}
</div>
```

#### Primary Button

```tsx
<Button className="
  bg-gradient-to-r from-primary to-primary/80
  hover:from-primary/90 hover:to-primary/70
  text-primary-foreground
  font-semibold
">
  Action
</Button>
```

#### Section Container

```tsx
<section className="
  max-w-7xl mx-auto
  px-4 sm:px-6 lg:px-8
  py-6 sm:py-8
  w-full
">
  {/* Content */}
</section>
```

#### Responsive Grid

```tsx
<div className="
  grid
  grid-cols-1
  sm:grid-cols-2
  lg:grid-cols-3
  gap-4 sm:gap-6
">
  {/* Grid items */}
</div>
```

---

## Inconsistencies & Future Improvements

### Known Inconsistencies

#### 1. LoadingSpinner Component

**Issue**: Uses hardcoded colors instead of design tokens

**Current**:
```tsx
<div className="border-b-2 border-blue-600" />
<p className="text-gray-600">Loading...</p>
```

**Recommended Fix**:
```tsx
<div className="border-b-2 border-primary" />
<p className="text-muted-foreground">Loading...</p>
```

**File**: `frontend/components/LoadingSpinner.tsx`

#### 2. SkeletonLoader Component

**Issue**: Uses hardcoded colors (`bg-gray-200`, `bg-white`) instead of design tokens

**Current**:
```tsx
<div className="bg-gray-200 rounded" />
<div className="bg-white rounded-lg shadow-md" />
```

**Recommended Fix**:
```tsx
<div className="bg-muted rounded animate-pulse" />
<div className="bg-card rounded-lg shadow-md" />
```

**File**: `frontend/components/SkeletonLoader.tsx`

#### 3. USER_COLORS Constant

**Issue**: Uses hex colors instead of OKLCH format

**Current**:
```tsx
export const USER_COLORS = [
  '#22d3ee',
  '#a78bfa',
  // ...
];
```

**Decision**: Keep as hex format (documented in `frontend/constants/colors.ts`)

**Rationale**:
- User-generated content: These colors are used for chat usernames and may need to be displayed in external systems or exported
- Compatibility: Hex colors are universally supported across all browsers and external integrations
- Simplicity: These are static color values that don't need the advanced features of OKLCH
- Performance: Hex colors are slightly more performant for inline styles

**File**: `frontend/constants/colors.ts` - Now includes documentation explaining the decision

#### 4. Border Radius Inconsistencies

**Issue**: Mixed usage of `rounded-md`, `rounded-lg`, `rounded-xl` without clear guidelines

**Recommendation**: 
- Standardize on `rounded-lg` for cards
- Use `rounded-md` for buttons/inputs
- Document usage guidelines (see [Border Radius](#border-radius) section)

#### 5. Opacity Value Inconsistencies

**Issue**: Mixed usage of opacity values (`/50`, `/80`, `/90`) without clear patterns

**Recommendation**: 
- Document standard opacity values (see [Color Usage Guidelines](#opacity-modifiers))
- Audit and standardize existing components

#### 6. Backdrop Blur Inconsistencies

**Issue**: Mixed usage of `backdrop-blur-sm` and `backdrop-blur-xl`

**Recommendation**: 
- Standardize: `backdrop-blur-sm` for cards, `backdrop-blur-xl` for navigation
- Document usage guidelines (see [Backdrop Blur](#backdrop-blur) section)

#### 7. Shadow Inconsistencies

**Issue**: Some hardcoded shadow values instead of design tokens

**Current**:
```tsx
className="shadow-[0_4px_20px_rgba(0,0,0,0.3)]"
```

**Recommendation**: 
- Create shadow utility classes or document standard shadow patterns
- Consider adding to design tokens if used frequently

### Future Improvements

#### 1. Design Token System

- **Goal**: Centralize all design tokens in a single source of truth
- **Approach**: Consider using a design token file (JSON/TypeScript) that generates CSS variables
- **Benefits**: Better tooling, easier theming, type safety

#### 2. Component Documentation

- **Goal**: Document all reusable components with Storybook or similar
- **Benefits**: Visual component library, easier onboarding, design consistency

#### 3. Dark/Light Theme Support

- **Goal**: Support both dark and light themes
- **Approach**: Extend CSS variables to support both themes
- **Current**: Dark theme only

#### 4. Animation System

- **Goal**: Standardize animation timings and easing functions
- **Approach**: Create animation utility classes
- **Benefits**: Consistent feel across interactions

#### 5. Spacing System Refinement

- **Goal**: Audit and standardize all spacing usage
- **Approach**: Create spacing utility classes for common patterns
- **Benefits**: More consistent layouts

#### 6. Typography Scale Refinement

- **Goal**: Ensure typography scale is used consistently
- **Approach**: Create typography utility classes
- **Benefits**: Better visual hierarchy

### Migration Checklist

When fixing inconsistencies:

- [x] Update `LoadingSpinner.tsx` to use design tokens ✅
- [x] Update `SkeletonLoader.tsx` to use design tokens ✅
- [x] Document `USER_COLORS` decision (keeping hex format) ✅
- [ ] Audit and standardize border radius usage
- [ ] Audit and standardize opacity values
- [ ] Audit and standardize backdrop blur values
- [ ] Audit and standardize shadow usage
- [ ] Update this document with any new patterns

---

## Summary

### Design System Overview

CyclingStream uses a **modern, dark-themed design system** built on:

- **OKLCH color space** for perceptually uniform colors
- **Tailwind CSS** for utility-first styling
- **shadcn/ui** style components for consistency
- **Mobile-first responsive design** for all devices
- **WCAG 2.1 AA** accessibility standards

### Key Design Tokens

- **Primary Color**: Green `oklch(0.72 0.19 155)` - Energy, cycling, nature
- **Background**: Very dark `oklch(0.11 0.005 260)` - Streaming platform aesthetic
- **Typography**: Geist Sans (primary), Geist Mono (code)
- **Border Radius**: 10px base (`0.625rem`)
- **Spacing**: 4px base unit (Tailwind scale)

### Consistency Principles

1. **Always use design tokens** - Never hardcode colors, spacing, or other values
2. **Mobile-first** - Design for mobile, enhance for desktop
3. **Semantic naming** - Use semantic color names (`primary`, `destructive`, etc.)
4. **Accessibility first** - All designs must meet WCAG AA standards
5. **Progressive enhancement** - Core functionality works everywhere

### Quick Reference

- **Colors**: See [Color System](#color-system)
- **Typography**: See [Typography](#typography)
- **Spacing**: See [Spacing & Layout](#spacing--layout)
- **Components**: See [Components](#components)
- **Effects**: See [Effects & Animations](#effects--animations)
- **Responsive**: See [Responsive Design](#responsive-design)

---

**Last Updated**: 2025  
**Maintained By**: CyclingStream Team  
**Questions?**: Refer to this document or update it with new patterns

