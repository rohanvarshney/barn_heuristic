# Web UI Visualizer Guide

## Purpose

The web UI is a Svelte-based interactive demo that animates ranking normalization from **before** to **after** for restaurant scores on `(1.0, 10.0]`.

It is designed to make the algorithm behavior visible in real time and explain how each strategy changes score distribution.

## Where the logic lives

- App orchestration: `web/src/App.svelte`
- Strategy engine: `web/src/lib/normalize.ts`
- Deterministic mock data: `web/src/lib/mockData.ts`
- Animation utilities: `web/src/lib/animation.ts`
- Controls: `web/src/components/ControlsPanel.svelte`
- Scatter visualization: `web/src/components/ScatterView.svelte`
- Distribution visualization: `web/src/components/DistributionView.svelte`

## How the UI works

1. The app generates deterministic mock restaurant ratings data.
2. It computes:
   - `beforeScores` (source)
   - `afterScores` (normalized, based on selected strategy)
3. During playback, it interpolates score values for a live in-between state:
   - `currentScores = interpolate(before, after, progress)`
4. Both visual panels consume the same `currentScores` so they stay synchronized.

## Visual panels

### 1) Scatter View

`ScatterView.svelte` plots each restaurant as a point:

- X-axis: stable item index (keeps identity/order)
- Y-axis: score value
- As progress moves 0 -> 1, points transition smoothly from original to normalized scores

### 2) Distribution View (Histogram + CDF)

`DistributionView.svelte` overlays three distributions:

- **Before** distribution
- **Current** interpolated distribution
- **After** normalized distribution

This makes density and shape changes visible while animation runs.

## Strategy controls

From `ControlsPanel.svelte`, users can:

- Switch algorithm:
  - `quantile_map`
  - `zscore_sigmoid`
  - `piecewise_bucket`
- Adjust piecewise bucket count
- Adjust animation speed
- Change data seed
- Play/Pause, Reset, and Regenerate data

## Mock data for 10 users (confirmation)

Yes, the UI supports mock data for 10 users, and this is already wired in code:

- `generateMockUsers(seed, userCount = 10, minRestaurants = 50, maxRestaurants = 500)` in `web/src/lib/mockData.ts`
- `App.svelte` initializes with `generateMockUsers(1337)` (default 10 users)
- The UI includes **Regenerate Data** + **Seed** controls to re-create deterministic datasets
- `App.svelte` also displays metadata text:
  - `Data points: ... restaurants across 10 deterministic users`

Each user gets a random (seeded, deterministic) number of restaurants between 50 and 500, and each item includes:

- `restaurantName`
- `score`

## Determinism model

Determinism is guaranteed by the seeded pseudo-random generator in `mockData.ts`:

- Same seed -> same users, same restaurants, same scores
- Different seed -> different dataset, still reproducible for that seed

## Local run

```bash
cd web
npm install
npm run dev
```

## Validation commands

```bash
cd web
npm run test
npm run check
npm run build
```
