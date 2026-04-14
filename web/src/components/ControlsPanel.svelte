<script lang="ts">
  import type { Strategy } from '../lib/types';

  type Props = {
    strategy: Strategy;
    bucketCount: number;
    speedMs: number;
    seed: number;
    playing: boolean;
    onPlayPause: () => void;
    onReset: () => void;
    onRegenerate: () => void;
  };

  let {
    strategy = $bindable('quantile_map'),
    bucketCount = $bindable(4),
    speedMs = $bindable(2200),
    seed = $bindable(1337),
    playing,
    onPlayPause,
    onReset,
    onRegenerate,
  }: Props = $props();
</script>

<section class="panel">
  <h2>Controls</h2>
  <label>
    Strategy
    <select bind:value={strategy}>
      <option value="quantile_map">quantile_map</option>
      <option value="zscore_sigmoid">zscore_sigmoid</option>
      <option value="piecewise_bucket">piecewise_bucket</option>
    </select>
  </label>

  <label>
    Bucket count ({bucketCount})
    <input type="range" min="2" max="12" step="1" bind:value={bucketCount} disabled={strategy !== 'piecewise_bucket'} />
  </label>

  <label>
    Animation speed: {speedMs}ms
    <input type="range" min="600" max="5000" step="100" bind:value={speedMs} />
  </label>

  <label>
    Seed
    <input type="number" bind:value={seed} min="1" />
  </label>

  <div class="buttons">
    <button type="button" onclick={onPlayPause}>{playing ? 'Pause' : 'Play'}</button>
    <button type="button" onclick={onReset}>Reset</button>
    <button type="button" onclick={onRegenerate}>Regenerate Data</button>
  </div>

  <p class="hint">
    Watch the same restaurants flow from original to normalized scores. The animation keeps item identity stable across views.
  </p>
</section>

