from __future__ import annotations

import copy
import enum
import math
from dataclasses import dataclass
from typing import Any, Callable, Iterable

EPSILON = 1e-9
MIN_SCORE = 1.0 + EPSILON
MAX_SCORE = 10.0


class Strategy(str, enum.Enum):
    QUANTILE_MAP = "quantile_map"
    ZSCORE_SIGMOID = "zscore_sigmoid"
    PIECEWISE_BUCKET = "piecewise_bucket"


@dataclass
class _ScoredItem:
    original_index: int
    score: float


def _clamp(score: float) -> float:
    if score <= 1.0:
        return MIN_SCORE
    if score > MAX_SCORE:
        return MAX_SCORE
    return score


def _resolve_strategy(strategy: str | Strategy) -> Strategy:
    if isinstance(strategy, Strategy):
        return strategy
    try:
        return Strategy(strategy)
    except ValueError as exc:
        raise ValueError(f"unknown strategy: {strategy}") from exc


def _validate_input(scores: Iterable[float]) -> list[float]:
    validated: list[float] = []
    for score in scores:
        if math.isnan(score) or score <= 1.0 or score > 10.0:
            raise ValueError(f"score {score} outside supported range (1.0, 10.0]")
        validated.append(float(score))
    return validated


def _quantile_map(values: list[float]) -> list[float]:
    n = len(values)
    if n == 0:
        return []
    if n == 1:
        return [values[0]]

    ranked = sorted(
        [_ScoredItem(i, s) for i, s in enumerate(values)],
        key=lambda item: (item.score, item.original_index),
    )
    out = [0.0] * n
    for rank, item in enumerate(ranked):
        percentile = rank / (n - 1)
        mapped = MIN_SCORE + percentile * (MAX_SCORE - MIN_SCORE)
        out[item.original_index] = mapped
    return out


def _zscore_sigmoid(values: list[float]) -> list[float]:
    n = len(values)
    if n <= 1:
        return values[:]

    mean = sum(values) / n
    variance = sum((v - mean) ** 2 for v in values) / n
    std = math.sqrt(variance)
    if std < EPSILON:
        return _quantile_map(values)

    out: list[float] = []
    for value in values:
        z = (value - mean) / std
        logistic = 1.0 / (1.0 + math.exp(-z))
        mapped = MIN_SCORE + logistic * (MAX_SCORE - MIN_SCORE)
        out.append(_clamp(mapped))
    return out


def _piecewise_bucket(values: list[float], buckets: int = 4) -> list[float]:
    n = len(values)
    if n <= 1:
        return values[:]

    bucket_count = max(2, buckets)
    width = (MAX_SCORE - MIN_SCORE) / bucket_count
    bucketed: list[list[_ScoredItem]] = [[] for _ in range(bucket_count)]

    for idx, score in enumerate(values):
        raw = int((score - MIN_SCORE) / width) if width > 0 else 0
        bucket_idx = max(0, min(bucket_count - 1, raw))
        bucketed[bucket_idx].append(_ScoredItem(idx, score))

    out = [0.0] * n
    write_start = MIN_SCORE
    total = float(n)

    for items in bucketed:
        if not items:
            continue
        fraction = len(items) / total
        span = fraction * (MAX_SCORE - MIN_SCORE)
        write_end = min(MAX_SCORE, write_start + span)
        sorted_items = sorted(items, key=lambda item: (item.score, item.original_index))

        if len(sorted_items) == 1:
            mapped = _clamp((write_start + write_end) / 2.0)
            out[sorted_items[0].original_index] = mapped
        else:
            for pos, item in enumerate(sorted_items):
                local_p = pos / (len(sorted_items) - 1)
                mapped = write_start + local_p * (write_end - write_start)
                out[item.original_index] = _clamp(mapped)
        write_start = write_end

    return [_clamp(v) for v in out]


def redistribute(
    items: list[Any],
    score_getter: Callable[[Any], float],
    score_setter: Callable[[Any, float], None] | None = None,
    strategy: str | Strategy = Strategy.QUANTILE_MAP,
    options: dict[str, Any] | None = None,
) -> list[Any]:
    if not isinstance(items, list):
        raise TypeError("items must be a list")

    opts = options or {}
    selected = _resolve_strategy(strategy)

    scores = _validate_input(score_getter(item) for item in items)

    if selected is Strategy.QUANTILE_MAP:
        redistributed = _quantile_map(scores)
    elif selected is Strategy.ZSCORE_SIGMOID:
        redistributed = _zscore_sigmoid(scores)
    else:
        buckets = int(opts.get("bucket_count", 4))
        redistributed = _piecewise_bucket(scores, buckets=buckets)

    if score_setter is not None:
        for item, new_score in zip(items, redistributed):
            score_setter(item, _clamp(new_score))
        return items

    cloned = copy.deepcopy(items)
    for item, new_score in zip(cloned, redistributed):
        if isinstance(item, dict):
            item["score"] = _clamp(new_score)
        else:
            setattr(item, "score", _clamp(new_score))
    return cloned
