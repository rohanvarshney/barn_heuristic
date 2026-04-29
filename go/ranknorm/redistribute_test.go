package ranknorm

import (
	"math"
	"testing"
)

type sample struct {
	ID    string
	Score float64
}

func TestQuantileMapDefaultSpreadsCluster(t *testing.T) {
	items := []sample{
		{ID: "a", Score: 8.01}, {ID: "b", Score: 8.02}, {ID: "c", Score: 8.01},
		{ID: "d", Score: 8.02}, {ID: "e", Score: 8.01}, {ID: "f", Score: 8.03},
	}
	out, err := Redistribute(items, func(s sample) float64 { return s.Score }, func(s *sample, v float64) {
		s.Score = v
	}, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seen := map[float64]bool{}
	for _, it := range out {
		if it.Score <= 1.0 || it.Score > 10.0 {
			t.Fatalf("score out of range: %f", it.Score)
		}
		if seen[it.Score] {
			t.Fatalf("expected unique quantile mapped score, got duplicate %f", it.Score)
		}
		seen[it.Score] = true
	}
}

func TestAllStrategiesCallable(t *testing.T) {
	items := []sample{{"a", 7.1}, {"b", 7.2}, {"c", 7.3}, {"d", 7.4}, {"e", 7.5}}
	q, _ := Redistribute(items, func(s sample) float64 { return s.Score }, func(s *sample, v float64) { s.Score = v }, StrategyQuantileMap, nil)
	z, _ := Redistribute(items, func(s sample) float64 { return s.Score }, func(s *sample, v float64) { s.Score = v }, StrategyZScoreSigmoid, nil)
	p, _ := Redistribute(items, func(s sample) float64 { return s.Score }, func(s *sample, v float64) { s.Score = v }, StrategyPiecewise, &Options{BucketCount: 5})
	if q[0].Score == z[0].Score && z[0].Score == p[0].Score {
		t.Fatalf("expected different strategy outputs")
	}
}

func TestStableTies(t *testing.T) {
	items := []sample{{"a", 5.5}, {"b", 5.5}, {"c", 9.0}}
	out, err := Redistribute(items, func(s sample) float64 { return s.Score }, func(s *sample, v float64) { s.Score = v }, StrategyQuantileMap, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !(out[0].Score < out[1].Score && out[1].Score < out[2].Score) {
		t.Fatalf("expected stable increasing ties")
	}
}

func TestOutOfRangeRejected(t *testing.T) {
	items := []sample{{"x", 1.0}}
	_, err := Redistribute(items, func(s sample) float64 { return s.Score }, func(s *sample, v float64) { s.Score = v }, StrategyQuantileMap, nil)
	if err == nil {
		t.Fatalf("expected error for out-of-range score")
	}
}

func TestNaNRejected(t *testing.T) {
	items := []sample{{"x", math.NaN()}}
	_, err := Redistribute(items, func(s sample) float64 { return s.Score }, func(s *sample, v float64) { s.Score = v }, StrategyQuantileMap, nil)
	if err == nil {
		t.Fatalf("expected error for NaN score")
	}
}

func TestEmptyAndSingle(t *testing.T) {
	out, err := Redistribute([]sample{}, func(s sample) float64 { return s.Score }, func(s *sample, v float64) { s.Score = v }, StrategyQuantileMap, nil)
	if err != nil || len(out) != 0 {
		t.Fatalf("expected empty output without error")
	}
	one := []sample{{"s", 7.2}}
	single, err := Redistribute(one, func(s sample) float64 { return s.Score }, func(s *sample, v float64) { s.Score = v }, StrategyQuantileMap, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if single[0].Score != 7.2 {
		t.Fatalf("single item should be unchanged")
	}
}
