package ranknorm

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

type restaurantItem struct {
	Name  string
	Score float64
}

func makeRestaurantItems(seed int64, n int) []restaurantItem {
	rng := rand.New(rand.NewSource(seed))
	items := make([]restaurantItem, 0, n)
	for i := 0; i < n; i++ {
		var score float64
		// Mix clustered + spread values, with occasional exact ties.
		if i%11 == 0 {
			score = 8.25
		} else if i%17 == 0 {
			score = 6.75
		} else {
			score = minScore + (maxScore-minScore)*rng.Float64()
		}
		if score <= 1.0 {
			score = minScore
		}
		if score > 10.0 {
			score = 10.0
		}
		items = append(items, restaurantItem{
			Name:  fmt.Sprintf("Restaurant #%03d", i),
			Score: score,
		})
	}
	return items
}

func assertOrderPreserved(t *testing.T, original []restaurantItem, redistributed []restaurantItem) {
	t.Helper()

	indices := make([]int, len(original))
	for i := range indices {
		indices[i] = i
	}
	sort.SliceStable(indices, func(i, j int) bool {
		ii := indices[i]
		jj := indices[j]
		if original[ii].Score == original[jj].Score {
			return ii < jj
		}
		return original[ii].Score < original[jj].Score
	})

	for k := 0; k < len(indices)-1; k++ {
		a := indices[k]
		b := indices[k+1]
		if redistributed[a].Score > redistributed[b].Score {
			t.Fatalf("order not preserved: original[%d]=%f <= original[%d]=%f but new[%d]=%f > new[%d]=%f",
				a, original[a].Score, b, original[b].Score, a, redistributed[a].Score, b, redistributed[b].Score)
		}
	}
}

func TestOrderPreservationQuantileMap(t *testing.T) {
	items := makeRestaurantItems(2026, 300)
	out, err := Redistribute(items, func(r restaurantItem) float64 { return r.Score }, func(r *restaurantItem, v float64) { r.Score = v }, StrategyQuantileMap, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertOrderPreserved(t, items, out)
}

func TestOrderPreservationZScoreSigmoid(t *testing.T) {
	items := makeRestaurantItems(2026, 300)
	out, err := Redistribute(items, func(r restaurantItem) float64 { return r.Score }, func(r *restaurantItem, v float64) { r.Score = v }, StrategyZScoreSigmoid, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertOrderPreserved(t, items, out)
}

func TestOrderPreservationPiecewiseBucket(t *testing.T) {
	items := makeRestaurantItems(2026, 300)
	out, err := Redistribute(items, func(r restaurantItem) float64 { return r.Score }, func(r *restaurantItem, v float64) { r.Score = v }, StrategyPiecewise, &Options{BucketCount: 6})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertOrderPreserved(t, items, out)
}

