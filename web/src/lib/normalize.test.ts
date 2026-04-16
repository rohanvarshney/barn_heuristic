import { describe, expect, it } from 'vitest';

import { normalizeScores } from './normalize';

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
