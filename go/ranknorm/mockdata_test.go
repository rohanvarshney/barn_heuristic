package ranknorm

import "testing"

func TestMockDataDeterministic(t *testing.T) {
	a, err := GenerateMockUserRatings(123, 10, 50, 500)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := GenerateMockUserRatings(123, 10, 50, 500)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a) != len(b) {
		t.Fatalf("expected same user count")
	}
	for i := range a {
		if a[i].UserID != b[i].UserID {
			t.Fatalf("user mismatch at %d", i)
		}
		if len(a[i].Ratings) != len(b[i].Ratings) {
			t.Fatalf("ratings length mismatch at %d", i)
		}
		for j := range a[i].Ratings {
			ar := a[i].Ratings[j]
			br := b[i].Ratings[j]
			if ar.RestaurantName != br.RestaurantName || ar.Rating != br.Rating {
				t.Fatalf("rating mismatch at user %d item %d", i, j)
			}
		}
	}
}

func TestMockDataShapeAndBounds(t *testing.T) {
	users, err := DefaultMockUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 10 {
		t.Fatalf("expected 10 users, got %d", len(users))
	}
	for _, u := range users {
		if len(u.Ratings) < 50 || len(u.Ratings) > 500 {
			t.Fatalf("expected 50..500 ratings, got %d", len(u.Ratings))
		}
		for _, r := range u.Ratings {
			if r.RestaurantName == "" {
				t.Fatalf("expected restaurant name")
			}
			if r.Rating <= 1.0 || r.Rating > 10.0 {
				t.Fatalf("rating out of bounds: %f", r.Rating)
			}
		}
	}
}

func BenchmarkFlattenRatings(b *testing.B) {
	users, err := DefaultMockUsers()
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FlattenRatings(users)
	}
}
