import time
import random
from ranknorm.redistribute import _quantile_map, _piecewise_bucket, _zscore_sigmoid

random.seed(42)
values = [random.uniform(1.1, 9.9) for _ in range(100000)]

start = time.time()
for _ in range(10):
    _quantile_map(values)
end = time.time()
print(f"Quantile Map Baseline: {end - start:.4f}s")

start = time.time()
for _ in range(10):
    _piecewise_bucket(values)
end = time.time()
print(f"Piecewise Bucket Baseline: {end - start:.4f}s")

start = time.time()
for _ in range(10):
    _zscore_sigmoid(values)
end = time.time()
print(f"ZScore Sigmoid Baseline: {end - start:.4f}s")
