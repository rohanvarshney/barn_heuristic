import { describe, expect, it } from 'vitest';

import { interpolateScores, normalizeScores } from './normalize';

const sampleScores = [8.2, 8.2, 8.4, 7.3, 9.1, 6.2, 7.9, 8.0, 8.1, 8.7];

describe('normalizeScores', () => {
  it('keeps all strategies within bounds', () => {
    const strategies = ['quantile_map', 'zscore_sigmoid', 'piecewise_bucket'] as const;
    for (const strategy of strategies) {
      const out = normalizeScores(sampleScores, strategy, { bucketCount: 5 });
      expect(out).toHaveLength(sampleScores.length);
      expect(Math.min(...out)).toBeGreaterThan(1.0);
      expect(Math.max(...out)).toBeLessThanOrEqual(10.0);
    }
  });

  it('is deterministic', () => {
    const a = normalizeScores(sampleScores, 'piecewise_bucket', { bucketCount: 6 });
    const b = normalizeScores(sampleScores, 'piecewise_bucket', { bucketCount: 6 });
    expect(a).toEqual(b);
  });

  it('preserves monotonic order across strategies', () => {
    const indices = [...sampleScores.keys()].sort((i, j) =>
      sampleScores[i] === sampleScores[j] ? i - j : sampleScores[i] - sampleScores[j],
    );

    const strategies = ['quantile_map', 'zscore_sigmoid', 'piecewise_bucket'] as const;
    for (const strategy of strategies) {
      const out = normalizeScores(sampleScores, strategy, { bucketCount: 6 });
      for (let k = 0; k < indices.length - 1; k += 1) {
        expect(out[indices[k]]).toBeLessThanOrEqual(out[indices[k + 1]]);
      }
    }
  });
});

describe('interpolateScores', () => {
  it('interpolates correctly at t=0', () => {
    const before = [1, 2, 3];
    const after = [10, 20, 30];
    expect(interpolateScores(before, after, 0)).toEqual(before);
  });

  it('interpolates correctly at t=1', () => {
    const before = [1, 2, 3];
    const after = [10, 20, 30];
    expect(interpolateScores(before, after, 1)).toEqual(after);
  });

  it('interpolates correctly at midpoint t=0.5', () => {
    const before = [1, 2, 3];
    const after = [11, 22, 33];
    expect(interpolateScores(before, after, 0.5)).toEqual([6, 12, 18]);
  });

  it('clamps t to [0, 1]', () => {
    const before = [1, 2, 3];
    const after = [10, 20, 30];
    expect(interpolateScores(before, after, -1)).toEqual(before);
    expect(interpolateScores(before, after, 2)).toEqual(after);
  });

  it('handles empty arrays', () => {
    expect(interpolateScores([], [], 0.5)).toEqual([]);
  });
});
