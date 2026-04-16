<script lang="ts">
  type Point = { x: number; y: number; label: string };

  type Props = {
    points: Point[];
    title?: string;
  };

  let { points, title = 'Score Scatter' }: Props = $props();

  const width = 980;
  const height = 290;
  const pad = 28;

  const x = (p: number) => pad + p * (width - pad * 2);
  const y = (score: number) => pad + (1 - (score - 1) / 9) * (height - pad * 2);
</script>

<section class="viz-card">
  <h3>{title}</h3>
  <svg viewBox={`0 0 ${width} ${height}`} role="img" aria-label="scatter chart">
    <line x1={pad} y1={height - pad} x2={width - pad} y2={height - pad} class="axis" />
    <line x1={pad} y1={pad} x2={pad} y2={height - pad} class="axis" />

    {#each [2, 4, 6, 8, 10] as tick}
      <line x1={pad} x2={width - pad} y1={y(tick)} y2={y(tick)} class="grid" />
      <text x="6" y={y(tick) + 4} class="tick">{tick.toFixed(1)}</text>
    {/each}

    {#each points as point}
      <circle cx={x(point.x)} cy={y(point.y)} r="2.2">
        <title>{point.label}: {point.y.toFixed(3)}</title>
      </circle>
    {/each}
  </svg>
</section>
