import { expect, test } from 'vitest';
import { interpolateScores, normalizeScores } from './src/lib/normalize';
import { flattenUsers, generateMockUsers } from './src/lib/mockData';

test('benchmark scatter points allocation strategy', () => {
  const items = flattenUsers(generateMockUsers(1337));
  const beforeScores = items.map((x) => x.score);
  const afterScores = normalizeScores(beforeScores, 'quantile_map', { bucketCount: 4 });

  function runBaseline() {
    let currentScores = [];
    let scatterPoints = [];
    for (let progress = 0; progress <= 1; progress += 0.001) {
      currentScores = interpolateScores(beforeScores, afterScores, progress);
      scatterPoints = items.map((item, i) => ({
        x: i / Math.max(1, items.length - 1),
        y: currentScores[i],
        label: item.restaurantName,
      }));
    }
    return scatterPoints;
  }

  function runOptimized() {
    let currentScores = [];
    let baseScatterPoints = items.map((item, i) => ({
      x: i / Math.max(1, items.length - 1),
      label: item.restaurantName,
    }));
    // We only simulate the creation and passing of base points and current scores, the DOM logic avoids object creation per tick.
    for (let progress = 0; progress <= 1; progress += 0.001) {
      currentScores = interpolateScores(beforeScores, afterScores, progress);
      // in Svelte 5, updating currentScores array is O(N) primitive operations, instead of N object allocations.
    }
    return currentScores;
  }

  const N = 10;
  let start = performance.now();
  for (let i = 0; i < N; i++) runBaseline();
  const baselineTime = performance.now() - start;

  start = performance.now();
  for (let i = 0; i < N; i++) runOptimized();
  const optimizedTime = performance.now() - start;

  console.log(`Baseline: ${baselineTime.toFixed(2)}ms`);
  console.log(`Optimized: ${optimizedTime.toFixed(2)}ms`);
  console.log(`Improvement: ${((baselineTime - optimizedTime) / baselineTime * 100).toFixed(2)}%`);
  expect(optimizedTime).toBeLessThan(baselineTime);
});
