import type { Strategy } from './types';

export const EPSILON = 1e-9;
export const MIN_SCORE = 1.0 + EPSILON;
export const MAX_SCORE = 10.0;

export type NormalizeOptions = {
  bucketCount?: number;
};

const clamp = (v: number): number => {
  if (v <= 1.0) return MIN_SCORE;
  if (v > MAX_SCORE) return MAX_SCORE;
  return v;
};

const validateInput = (scores: number[]): number[] => {
  return scores.map((s) => {
    if (s <= 1.0 || s > 10.0) {
      throw new Error(`score ${s} outside supported range (1.0, 10.0]`);
    }
    return s;
  });
};

const quantileMap = (values: number[]): number[] => {
  const n = values.length;
  if (n <= 1) return [...values];

  const indices = Array.from({ length: n }, (_, i) => i);
  indices.sort((a, b) => (values[a] === values[b] ? a - b : values[a] - values[b]));

  const out = new Array<number>(n);
  indices.forEach((index, rank) => {
    const percentile = rank / (n - 1);
    out[index] = MIN_SCORE + percentile * (MAX_SCORE - MIN_SCORE);
  });
  return out;
};

const zscoreSigmoid = (values: number[]): number[] => {
  const n = values.length;
  if (n <= 1) return [...values];

  let sum = 0;
  for (let i = 0; i < n; i++) {
    sum += values[i];
  }
  const mean = sum / n;

  let varianceSum = 0;
  for (let i = 0; i < n; i++) {
    const diff = values[i] - mean;
    varianceSum += diff * diff;
  }
  const variance = varianceSum / n;
  const std = Math.sqrt(variance);
  if (std < EPSILON) return quantileMap(values);

  const out = new Array<number>(n);
  for (let i = 0; i < n; i++) {
    const z = (values[i] - mean) / std;
    const logistic = 1 / (1 + Math.exp(-z));
    out[i] = clamp(MIN_SCORE + logistic * (MAX_SCORE - MIN_SCORE));
  }
  return out;
};

const piecewiseBucket = (values: number[], bucketCount = 4): number[] => {
  const n = values.length;
  if (n <= 1) return [...values];

  const buckets = Math.max(2, bucketCount);
  const width = (MAX_SCORE - MIN_SCORE) / buckets;
  const grouped: number[][] = Array.from({ length: buckets }, () => []);

  values.forEach((score, index) => {
    const raw = width > 0 ? Math.floor((score - MIN_SCORE) / width) : 0;
    const bucketIdx = Math.max(0, Math.min(buckets - 1, raw));
    grouped[bucketIdx].push(index);
  });

  const out = new Array<number>(n);
  let writeStart = MIN_SCORE;
  const total = n;

  grouped.forEach((items) => {
    if (items.length === 0) return;
    const span = (items.length / total) * (MAX_SCORE - MIN_SCORE);
    const writeEnd = Math.min(MAX_SCORE, writeStart + span);
    items.sort((a, b) => (values[a] === values[b] ? a - b : values[a] - values[b]));

    if (items.length === 1) {
      out[items[0]] = clamp((writeStart + writeEnd) / 2);
    } else {
      items.forEach((index, pos) => {
        const p = pos / (items.length - 1);
        out[index] = clamp(writeStart + p * (writeEnd - writeStart));
      });
    }

    writeStart = writeEnd;
  });

  return out.map(clamp);
};

export const normalizeScores = (
  scores: number[],
  strategy: Strategy = 'quantile_map',
  options: NormalizeOptions = {},
): number[] => {
  const valid = validateInput(scores);
  if (strategy === 'quantile_map') return quantileMap(valid);
  if (strategy === 'zscore_sigmoid') return zscoreSigmoid(valid);
  return piecewiseBucket(valid, options.bucketCount ?? 4);
};

export const interpolateScores = (before: number[], after: number[], t: number): number[] => {
  const clampedT = Math.max(0, Math.min(1, t));
  return before.map((v, i) => v + (after[i] - v) * clampedT);
};
