export type AnimatorState = {
  progress: number;
  playing: boolean;
};

export const createProgressStepper = (
  durationMs: number,
  onUpdate: (t: number) => void,
  onDone?: () => void,
) => {
  let rafId = 0;
  let start = 0;
  let running = false;

  const tick = (timestamp: number) => {
    if (!running) return;
    if (!start) start = timestamp;
    const elapsed = timestamp - start;
    const t = Math.max(0, Math.min(1, elapsed / durationMs));
    onUpdate(t);
    if (t < 1) {
      rafId = requestAnimationFrame(tick);
    } else {
      running = false;
      onDone?.();
    }
  };

  return {
    play() {
      if (running) return;
      running = true;
      start = 0;
      rafId = requestAnimationFrame(tick);
    },
    stop() {
      running = false;
      if (rafId) cancelAnimationFrame(rafId);
    },
  };
};

export const lerp = (a: number, b: number, t: number): number => {
  const clampedT = Math.max(0, Math.min(1, t));
  return a + (b - a) * clampedT;
};

