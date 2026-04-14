package ranknorm

import (
	"fmt"
	"math/rand"
)

type RestaurantRating struct {
	RestaurantName string
	Rating         float64
}

type UserRatings struct {
	UserID  string
	Ratings []RestaurantRating
}

var adjectives = []string{
	"Golden", "Rusty", "Hidden", "Spicy", "Cozy", "Lucky", "Midnight", "Sunny",
	"Royal", "Humble", "Crispy", "Saffron", "Electric", "Silver", "Jade", "Maple",
}

var nouns = []string{
	"Noodle", "Bistro", "Tavern", "Kitchen", "Diner", "Sushi", "Bakery", "Grill",
	"Pizzeria", "Curry", "Cantina", "Cafe", "Steakhouse", "Ramen", "Falafel", "Taqueria",
}

func restaurantName(rng *rand.Rand, idx int) string {
	adj := adjectives[rng.Intn(len(adjectives))]
	noun := nouns[rng.Intn(len(nouns))]
	return fmt.Sprintf("%s %s #%03d", adj, noun, idx)
}

func boundedRating(rng *rand.Rand) float64 {
	// Roughly clustered around ~8, bounded to (1.0, 10.0].
	raw := rng.NormFloat64()*0.75 + 8.0
	if raw <= 1.0 {
		return minScore
	}
	if raw > 10.0 {
		return maxScore
	}
	return raw
}

func GenerateMockUserRatings(seed int64, userCount, minRestaurants, maxRestaurants int) ([]UserRatings, error) {
	if userCount <= 0 {
		return nil, fmt.Errorf("userCount must be > 0")
	}
	if minRestaurants <= 0 || maxRestaurants <= 0 {
		return nil, fmt.Errorf("minRestaurants and maxRestaurants must be > 0")
	}
	if minRestaurants > maxRestaurants {
		return nil, fmt.Errorf("minRestaurants must be <= maxRestaurants")
	}

	rng := rand.New(rand.NewSource(seed))
	users := make([]UserRatings, 0, userCount)
	for u := 0; u < userCount; u++ {
		userID := fmt.Sprintf("user_%02d", u+1)
		n := rng.Intn(maxRestaurants-minRestaurants+1) + minRestaurants

		rr := make([]RestaurantRating, 0, n)
		for i := 1; i <= n; i++ {
			rr = append(rr, RestaurantRating{
				RestaurantName: restaurantName(rng, i),
				Rating:         boundedRating(rng),
			})
		}
		users = append(users, UserRatings{UserID: userID, Ratings: rr})
	}
	return users, nil
}

func DefaultMockUsers() ([]UserRatings, error) {
	return GenerateMockUserRatings(1337, 10, 50, 500)
}

func FlattenRatings(users []UserRatings) []map[string]any {
	out := make([]map[string]any, 0)
	for _, u := range users {
		for _, r := range u.Ratings {
			out = append(out, map[string]any{
				"user_id":         u.UserID,
				"restaurant_name": r.RestaurantName,
				"score":           r.Rating,
			})
		}
	}
	return out
}

