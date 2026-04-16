<script lang="ts">
  import ControlsPanel from './components/ControlsPanel.svelte';
  import DistributionView from './components/DistributionView.svelte';
  import ScatterView from './components/ScatterView.svelte';
  import { createProgressStepper } from './lib/animation';
  import { flattenUsers, generateMockUsers } from './lib/mockData';
  import { interpolateScores, normalizeScores } from './lib/normalize';
  import type { Strategy } from './lib/types';

  let strategy = $state<Strategy>('quantile_map');
  let bucketCount = $state(4);
  let speedMs = $state(2200);
  let seed = $state(1337);
  let playing = $state(false);
  let progress = $state(0);

  let items = $state(flattenUsers(generateMockUsers(1337)));
  let beforeScores = $derived(items.map((x) => x.score));
  let afterScores = $derived(normalizeScores(beforeScores, strategy, { bucketCount }));
  let currentScores = $derived(interpolateScores(beforeScores, afterScores, progress));

  let scatterPoints = $derived(
    items.map((item, i) => ({
      x: i / Math.max(1, items.length - 1),
      y: currentScores[i],
      label: item.restaurantName,
    })),
  );

  let stepper = createProgressStepper(
    2200,
    (t) => {
      progress = t;
    },
    () => {
      playing = false;
    },
  );

  $effect(() => {
    stepper.stop();
    stepper = createProgressStepper(
      speedMs,
      (t) => {
        progress = t;
      },
      () => {
        playing = false;
      },
    );
    if (playing) {
      stepper.play();
    }
  });

  const togglePlay = () => {
    if (playing) {
      playing = false;
      stepper.stop();
      return;
    }
    if (progress >= 1) progress = 0;
    playing = true;
    stepper.play();
  };

  const reset = () => {
    playing = false;
    stepper.stop();
    progress = 0;
  };

  const regenerate = () => {
    reset();
    items = flattenUsers(generateMockUsers(seed));
  };

  const strategyDescription = $derived.by(() => {
    if (strategy === 'quantile_map') {
      return 'Quantile mapping repositions each restaurant by rank percentile to spread scores evenly.';
    }
    if (strategy === 'zscore_sigmoid') {
      return 'Z-score + sigmoid smooths score crowding by standardizing values then rescaling.';
    }
    return 'Piecewise bucket allocates output range based on bucket density and smooth interpolation.';
  });
</script>

<main>
  <header>
    <h1>Realtime Ranking Normalization Visualizer</h1>
    <p>
      See how restaurant scores redistribute from crowded ranges into a normalized spread across all three
      algorithms.
    </p>
    <p>
      <a href="https://github.com/rohanvarshney/barn_heuristic" target="_blank" rel="noopener noreferrer">View Repository on GitHub</a>
    </p>
  </header>

  <ControlsPanel
    bind:strategy
    bind:bucketCount
    bind:speedMs
    bind:seed
    {playing}
    onPlayPause={togglePlay}
    onReset={reset}
    onRegenerate={regenerate}
  />

  <section class="meta">
    <p><strong>Progress:</strong> {(progress * 100).toFixed(1)}%</p>
    <p><strong>Data points:</strong> {items.length} restaurants across 10 deterministic users</p>
    <p><strong>Algorithm:</strong> {strategyDescription}</p>
  </section>

  <ScatterView points={scatterPoints} title="Restaurant Score Motion (Before → After)" />
  <DistributionView {beforeScores} {currentScores} {afterScores} bins={20} />
</main>
