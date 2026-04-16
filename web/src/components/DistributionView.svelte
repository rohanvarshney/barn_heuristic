<script lang="ts">
  type Props = {
    beforeScores: number[];
    currentScores: number[];
    afterScores: number[];
    bins?: number;
  };

  let { beforeScores, currentScores, afterScores, bins = 20 }: Props = $props();

  const width = 980;
  const height = 280;
  const pad = 26;

  const buildHistogram = (scores: number[]) => {
    const out = new Array<number>(bins).fill(0);
    for (const score of scores) {
      const normalized = (score - 1.0) / 9.0;
      const idx = Math.max(0, Math.min(bins - 1, Math.floor(normalized * bins)));
      out[idx] += 1;
    }
    return out;
  };

  const buildCdf = (scores: number[]) => {
    const sorted = [...scores].sort((a, b) => a - b);
    return sorted.map((score, i) => ({
      score,
      p: (i + 1) / sorted.length,
    }));
  };

  let histBefore = $derived(buildHistogram(beforeScores));
  let histCurrent = $derived(buildHistogram(currentScores));
  let histAfter = $derived(buildHistogram(afterScores));

  let cdfBefore = $derived(buildCdf(beforeScores));
  let cdfCurrent = $derived(buildCdf(currentScores));
  let cdfAfter = $derived(buildCdf(afterScores));

  let maxBin = $derived(Math.max(1, ...histBefore, ...histCurrent, ...histAfter));

  const x = (score: number) => pad + ((score - 1.0) / 9.0) * (width - pad * 2);
  const yDensity = (count: number) => height - pad - (count / maxBin) * (height - pad * 2);
  const yCdf = (p: number) => height - pad - p * (height - pad * 2);

  let barW = $derived((width - pad * 2) / bins);

  const cdfPath = (cdf: { score: number; p: number }[]) => {
    if (cdf.length === 0) return '';
    return cdf.map((pt, i) => `${i === 0 ? 'M' : 'L'} ${x(pt.score)} ${yCdf(pt.p)}`).join(' ');
  };
</script>

<section class="viz-card">
  <h3>Distribution Morph (Histogram + CDF)</h3>
  <svg viewBox={`0 0 ${width} ${height}`} role="img" aria-label="distribution chart">
    <line x1={pad} y1={height - pad} x2={width - pad} y2={height - pad} class="axis" />
    <line x1={pad} y1={pad} x2={pad} y2={height - pad} class="axis" />

    {#each histAfter as count, i}
      <rect
        x={pad + i * barW + 1}
        y={yDensity(count)}
        width={Math.max(1, barW - 2)}
        height={height - pad - yDensity(count)}
        class="hist-after"
      />
    {/each}

    {#each histBefore as count, i}
      <rect
        x={pad + i * barW + 1}
        y={yDensity(count)}
        width={Math.max(1, barW - 2)}
        height={height - pad - yDensity(count)}
        class="hist-before"
      />
    {/each}

    {#each histCurrent as count, i}
      <rect
        x={pad + i * barW + 1}
        y={yDensity(count)}
        width={Math.max(1, barW - 2)}
        height={height - pad - yDensity(count)}
        class="hist-current"
      />
    {/each}

    <path d={cdfPath(cdfBefore)} class="cdf-before" />
    <path d={cdfPath(cdfAfter)} class="cdf-after" />
    <path d={cdfPath(cdfCurrent)} class="cdf-current" />
  </svg>
  <p class="legend">
    <span class="before">Before</span>
    <span class="current">Current</span>
    <span class="after">After</span>
  </p>
</section>
