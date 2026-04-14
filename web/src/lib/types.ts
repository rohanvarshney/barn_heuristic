export type Strategy = 'quantile_map' | 'zscore_sigmoid' | 'piecewise_bucket';

export const STRATEGIES: Strategy[] = ['quantile_map', 'zscore_sigmoid', 'piecewise_bucket'];

export type RestaurantScore = {
  id: string;
  restaurantName: string;
  score: number;
};

export type UserRatings = {
  userId: string;
  ratings: RestaurantScore[];
};

