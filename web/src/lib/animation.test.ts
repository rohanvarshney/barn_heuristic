import { describe, expect, it } from 'vitest';

import { lerp } from './animation';

describe('animation utilities', () => {
  it('lerp interpolates and clamps', () => {
    expect(lerp(2, 10, 0)).toBe(2);
    expect(lerp(2, 10, 0.5)).toBe(6);
    expect(lerp(2, 10, 1)).toBe(10);
    expect(lerp(2, 10, -10)).toBe(2);
    expect(lerp(2, 10, 3)).toBe(10);
  });
});
