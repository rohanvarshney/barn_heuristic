import pytest

from ranknorm.mock_data import flatten_ratings, generate_mock_user_ratings


def test_mock_data_is_deterministic():
    a = generate_mock_user_ratings(seed=123)
    b = generate_mock_user_ratings(seed=123)
    assert a == b


def test_mock_data_shape_and_bounds():
    users = generate_mock_user_ratings(
        seed=999, user_count=10, min_restaurants=50, max_restaurants=500
    )
    assert len(users) == 10
    for u in users:
        assert 50 <= len(u.ratings) <= 500
        for r in u.ratings:
            assert isinstance(r.restaurant_name, str) and r.restaurant_name
            assert r.rating > 1.0
            assert r.rating <= 10.0


def test_flatten_ratings_has_required_fields():
    users = generate_mock_user_ratings(seed=42)
    flat = flatten_ratings(users)
    assert flat
    sample = flat[0]
    assert {"user_id", "restaurant_name", "score"} <= set(sample.keys())


def test_generate_mock_user_ratings_invalid_inputs():
    with pytest.raises(ValueError, match="user_count must be > 0"):
        generate_mock_user_ratings(user_count=0)
    with pytest.raises(ValueError, match="user_count must be > 0"):
        generate_mock_user_ratings(user_count=-1)

    with pytest.raises(
        ValueError, match="min_restaurants and max_restaurants must be > 0"
    ):
        generate_mock_user_ratings(min_restaurants=0)
    with pytest.raises(
        ValueError, match="min_restaurants and max_restaurants must be > 0"
    ):
        generate_mock_user_ratings(max_restaurants=0)

    with pytest.raises(ValueError, match="min_restaurants must be <= max_restaurants"):
        generate_mock_user_ratings(min_restaurants=100, max_restaurants=50)
