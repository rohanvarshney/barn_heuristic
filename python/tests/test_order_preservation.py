import random

import pytest

from ranknorm.redistribute import Strategy, redistribute


def _make_items(seed: int, n: int = 250):
    rng = random.Random(seed)
    items = []
    for i in range(n):
        # Mix clustered + spread values, with occasional exact ties.
        if i % 11 == 0:
            score = 8.25
        elif i % 17 == 0:
            score = 6.75
        else:
            score = 1.000000001 + (9.0 * rng.random())
        if score <= 1.0:
            score = 1.000000001
        if score > 10.0:
            score = 10.0
        items.append({"restaurant_name": f"Restaurant #{i:03d}", "score": float(score)})
    return items


@pytest.mark.parametrize(
    "strategy,options",
    [
        (Strategy.QUANTILE_MAP, None),
        (Strategy.ZSCORE_SIGMOID, None),
        (Strategy.PIECEWISE_BUCKET, {"bucket_count": 6}),
    ],
)
def test_order_preserved_after_redistribution(strategy, options):
    items = _make_items(seed=2026, n=300)

    original_scores = [x["score"] for x in items]
    stable_sorted_indices = sorted(range(len(items)), key=lambda i: (original_scores[i], i))

    out = redistribute(items, lambda x: x["score"], strategy=strategy, options=options)
    new_scores = [x["score"] for x in out]

    # If original score is non-decreasing in this stable order,
    # then redistributed scores must also be non-decreasing.
    for prev_i, curr_i in zip(stable_sorted_indices, stable_sorted_indices[1:]):
        assert new_scores[prev_i] <= new_scores[curr_i]

