from __future__ import annotations

from dataclasses import dataclass
import random
from typing import Iterable

from .redistribute import MAX_SCORE, MIN_SCORE


@dataclass(frozen=True)
class RestaurantRating:
    restaurant_name: str
    rating: float


@dataclass(frozen=True)
class UserRatings:
    user_id: str
    ratings: tuple[RestaurantRating, ...]


_ADJECTIVES: tuple[str, ...] = (
    "Golden",
    "Rusty",
    "Hidden",
    "Spicy",
    "Cozy",
    "Lucky",
    "Midnight",
    "Sunny",
    "Royal",
    "Humble",
    "Crispy",
    "Saffron",
    "Electric",
    "Silver",
    "Jade",
    "Maple",
)

_NOUNS: tuple[str, ...] = (
    "Noodle",
    "Bistro",
    "Tavern",
    "Kitchen",
    "Diner",
    "Sushi",
    "Bakery",
    "Grill",
    "Pizzeria",
    "Curry",
    "Cantina",
    "Cafe",
    "Steakhouse",
    "Ramen",
    "Falafel",
    "Taqueria",
)


def _restaurant_name(rng: random.Random, idx: int) -> str:
    adj = _ADJECTIVES[rng.randrange(len(_ADJECTIVES))]
    noun = _NOUNS[rng.randrange(len(_NOUNS))]
    return f"{adj} {noun} #{idx:03d}"


def _bounded_rating(rng: random.Random) -> float:
    # Produce a realistic-ish, somewhat clustered distribution near 7–9,
    # while respecting the SDK contract (1.0, 10.0].
    raw = rng.gauss(mu=8.0, sigma=0.75)
    if raw <= 1.0:
        return MIN_SCORE
    if raw > 10.0:
        return MAX_SCORE
    return float(raw)


def generate_mock_user_ratings(
    *,
    seed: int = 1337,
    user_count: int = 10,
    min_restaurants: int = 50,
    max_restaurants: int = 500,
) -> list[UserRatings]:
    if user_count <= 0:
        raise ValueError("user_count must be > 0")
    if min_restaurants <= 0 or max_restaurants <= 0:
        raise ValueError("min_restaurants and max_restaurants must be > 0")
    if min_restaurants > max_restaurants:
        raise ValueError("min_restaurants must be <= max_restaurants")

    rng = random.Random(seed)
    users: list[UserRatings] = []

    for u in range(user_count):
        user_id = f"user_{u + 1:02d}"
        restaurant_count = rng.randint(min_restaurants, max_restaurants)

        rr: list[RestaurantRating] = []
        for i in range(restaurant_count):
            rr.append(
                RestaurantRating(
                    restaurant_name=_restaurant_name(rng, i + 1),
                    rating=_bounded_rating(rng),
                )
            )

        users.append(UserRatings(user_id=user_id, ratings=tuple(rr)))

    return users


def flatten_ratings(users: Iterable[UserRatings]) -> list[dict]:
    items: list[dict] = []
    for user in users:
        for r in user.ratings:
            items.append(
                {
                    "user_id": user.user_id,
                    "restaurant_name": r.restaurant_name,
                    "score": r.rating,
                }
            )
    return items
