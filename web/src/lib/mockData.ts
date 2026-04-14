import type { RestaurantScore, UserRatings } from './types';
import { MAX_SCORE, MIN_SCORE } from './normalize';

class SeededRandom {
  private state: number;

  constructor(seed: number) {
    this.state = seed >>> 0;
  }

  next(): number {
    this.state = (1664525 * this.state + 1013904223) >>> 0;
    return this.state / 0x100000000;
  }

  int(min: number, max: number): number {
    return Math.floor(this.next() * (max - min + 1)) + min;
  }

  normal(mean = 0, std = 1): number {
    const u1 = Math.max(this.next(), Number.MIN_VALUE);
    const u2 = this.next();
    const z0 = Math.sqrt(-2 * Math.log(u1)) * Math.cos(2 * Math.PI * u2);
    return mean + z0 * std;
  }
}

const ADJECTIVES = [
  'Golden',
  'Rusty',
  'Hidden',
  'Spicy',
  'Cozy',
  'Lucky',
  'Midnight',
  'Sunny',
  'Royal',
  'Humble',
  'Crispy',
  'Saffron',
];
const NOUNS = [
  'Bistro',
  'Kitchen',
  'Tavern',
  'Grill',
  'Ramen',
  'Pizzeria',
  'CurryHouse',
  'Cafe',
  'Diner',
  'Taqueria',
  'Bakery',
  'SushiBar',
];

const restaurantName = (rng: SeededRandom, idx: number): string => {
  const adj = ADJECTIVES[rng.int(0, ADJECTIVES.length - 1)];
  const noun = NOUNS[rng.int(0, NOUNS.length - 1)];
  return `${adj} ${noun} #${String(idx).padStart(3, '0')}`;
};

const boundedRating = (rng: SeededRandom): number => {
  const raw = rng.normal(8.0, 0.75);
  if (raw <= 1.0) return MIN_SCORE;
  if (raw > 10.0) return MAX_SCORE;
  return raw;
};

export const generateMockUsers = (
  seed = 1337,
  userCount = 10,
  minRestaurants = 50,
  maxRestaurants = 500,
): UserRatings[] => {
  const rng = new SeededRandom(seed);
  const users: UserRatings[] = [];

  for (let u = 0; u < userCount; u += 1) {
    const restaurantCount = rng.int(minRestaurants, maxRestaurants);
    const ratings: RestaurantScore[] = [];
    for (let i = 0; i < restaurantCount; i += 1) {
      ratings.push({
        id: `u${u + 1}-r${i + 1}`,
        restaurantName: restaurantName(rng, i + 1),
        score: boundedRating(rng),
      });
    }
    users.push({
      userId: `user_${String(u + 1).padStart(2, '0')}`,
      ratings,
    });
  }

  return users;
};

export const flattenUsers = (users: UserRatings[]): RestaurantScore[] => {
  return users.flatMap((u) => u.ratings.map((r) => ({ ...r, id: `${u.userId}-${r.id}` })));
};

