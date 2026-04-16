from ranknorm.redistribute import Strategy, redistribute


def _scores(items):
    return [item["score"] for item in items]


def test_quantile_map_default_evenly_spreads_clustered_scores():
    items = [{"id": i, "score": 8.0 + (0.01 * (i % 3))} for i in range(12)]
    out = redistribute(items, score_getter=lambda x: x["score"])
    scores = _scores(out)
    assert min(scores) > 1.0
    assert max(scores) <= 10.0
    assert len(set(round(v, 6) for v in scores)) == len(scores)


def test_strategy_switch_all_three():
    items = [{"id": i, "score": 7.5 + i * 0.1} for i in range(8)]
    q = redistribute(items, lambda x: x["score"], strategy=Strategy.QUANTILE_MAP)
    z = redistribute(items, lambda x: x["score"], strategy=Strategy.ZSCORE_SIGMOID)
    p = redistribute(items, lambda x: x["score"], strategy=Strategy.PIECEWISE_BUCKET)
    assert _scores(q) != _scores(z)
    assert _scores(z) != _scores(p)


def test_stable_tie_ordering_under_quantile():
    items = [
        {"id": "a", "score": 5.5},
        {"id": "b", "score": 5.5},
        {"id": "c", "score": 9.0},
    ]
    out = redistribute(items, lambda x: x["score"], strategy=Strategy.QUANTILE_MAP)
    score_by_id = {x["id"]: x["score"] for x in out}
    assert score_by_id["a"] < score_by_id["b"] < score_by_id["c"]


def test_raises_on_out_of_range_score():
    items = [{"score": 1.0}]
    try:
        redistribute(items, lambda x: x["score"])
        assert False, "expected ValueError"
    except ValueError:
        assert True


def test_raises_on_nan_score():
    items = [{"score": float("nan")}]
    try:
        redistribute(items, lambda x: x["score"])
        assert False, "expected ValueError"
    except ValueError:
        assert True


def test_empty_and_single_are_supported():
    assert redistribute([], lambda x: x["score"]) == []
    single = [{"score": 7.2}]
    out = redistribute(single, lambda x: x["score"])
    assert out[0]["score"] == 7.2
