package ranknorm

import (
	"fmt"
	"math"
	"sort"
)

const (
	epsilon  = 1e-9
	minScore = 1.0 + epsilon
	maxScore = 10.0
)

type Strategy string

const (
	StrategyQuantileMap   Strategy = "quantile_map"
	StrategyZScoreSigmoid Strategy = "zscore_sigmoid"
	StrategyPiecewise     Strategy = "piecewise_bucket"
)

type Options struct {
	BucketCount int
}

type scoredItem struct {
	index int
	score float64
}

func clamp(v float64) float64 {
	if v <= 1.0 {
		return minScore
	}
	if v > maxScore {
		return maxScore
	}
	return v
}

func validateScores(scores []float64) error {
	for _, s := range scores {
		if s <= 1.0 || s > 10.0 {
			return fmt.Errorf("score %.6f outside supported range (1.0, 10.0]", s)
		}
	}
	return nil
}

func quantileMap(values []float64) []float64 {
	n := len(values)
	if n <= 1 {
		return values
	}
	ranked := make([]scoredItem, n)
	for i, s := range values {
		ranked[i] = scoredItem{index: i, score: s}
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score == ranked[j].score {
			return ranked[i].index < ranked[j].index
		}
		return ranked[i].score < ranked[j].score
	})
	out := make([]float64, n)
	for rank, item := range ranked {
		pct := float64(rank) / float64(n-1)
		out[item.index] = minScore + pct*(maxScore-minScore)
	}
	return out
}

func zscoreSigmoid(values []float64) []float64 {
	n := len(values)
	if n <= 1 {
		return values
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(n)
	variance := 0.0
	for _, v := range values {
		d := v - mean
		variance += d * d
	}
	variance /= float64(n)
	std := math.Sqrt(variance)
	if std < epsilon {
		return quantileMap(values)
	}
	out := make([]float64, n)
	for i, v := range values {
		z := (v - mean) / std
		logistic := 1.0 / (1.0 + math.Exp(-z))
		out[i] = clamp(minScore + logistic*(maxScore-minScore))
	}
	return out
}

func piecewiseBucket(values []float64, buckets int) []float64 {
	n := len(values)
	if n <= 1 {
		return values
	}
	if buckets < 2 {
		buckets = 4
	}
	width := (maxScore - minScore) / float64(buckets)
	bucketed := make([][]scoredItem, buckets)
	for i, v := range values {
		raw := 0
		if width > 0 {
			raw = int((v - minScore) / width)
		}
		if raw < 0 {
			raw = 0
		}
		if raw >= buckets {
			raw = buckets - 1
		}
		bucketed[raw] = append(bucketed[raw], scoredItem{index: i, score: v})
	}
	out := make([]float64, n)
	writeStart := minScore
	total := float64(n)
	for _, items := range bucketed {
		if len(items) == 0 {
			continue
		}
		fraction := float64(len(items)) / total
		span := fraction * (maxScore - minScore)
		writeEnd := math.Min(maxScore, writeStart+span)
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].score == items[j].score {
				return items[i].index < items[j].index
			}
			return items[i].score < items[j].score
		})
		if len(items) == 1 {
			out[items[0].index] = clamp((writeStart + writeEnd) / 2.0)
		} else {
			for pos, item := range items {
				p := float64(pos) / float64(len(items)-1)
				out[item.index] = clamp(writeStart + p*(writeEnd-writeStart))
			}
		}
		writeStart = writeEnd
	}
	for i, v := range out {
		out[i] = clamp(v)
	}
	return out
}

func Redistribute[T any](
	items []T,
	getScore func(T) float64,
	setScore func(*T, float64),
	strategy Strategy,
	opts *Options,
) ([]T, error) {
	out := make([]T, len(items))
	copy(out, items)

	if strategy == "" {
		strategy = StrategyQuantileMap
	}

	scores := make([]float64, len(items))
	for i := range out {
		scores[i] = getScore(out[i])
	}
	if err := validateScores(scores); err != nil {
		return nil, err
	}

	var redistributed []float64
	switch strategy {
	case StrategyQuantileMap:
		redistributed = quantileMap(scores)
	case StrategyZScoreSigmoid:
		redistributed = zscoreSigmoid(scores)
	case StrategyPiecewise:
		bucketCount := 4
		if opts != nil && opts.BucketCount > 0 {
			bucketCount = opts.BucketCount
		}
		redistributed = piecewiseBucket(scores, bucketCount)
	default:
		return nil, fmt.Errorf("unknown strategy: %s", strategy)
	}

	for i := range out {
		setScore(&out[i], clamp(redistributed[i]))
	}

	return out, nil
}
