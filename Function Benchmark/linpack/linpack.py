import time
import sys
from numpy import matrix, linalg, random


def linpack(n):
    # LINPACK benchmarks
    ops = (2.0 * n) * n * n / 3.0 + (2.0 * n) * n

    # Create AxA array of random numbers -0.5 to 0.5
    A = random.random_sample((n, n)) - 0.5
    B = A.sum(axis=1)

    # Convert to matrices
    A = matrix(A)
    B = matrix(B.reshape((n, 1)))

    # Ax = B
    s = time.time()
    x = linalg.solve(A, B)
    latency = time.time() - s

    mflops = (ops * 1e-6 / latency)

    result = {
        'mflops': mflops,
        'latency': latency
    }

    return result


def f(n):
    start = round(time.time(),6)

    result = linpack(int(n))
    end = round(time.time(),6)


def handler(event, context):
    # n = event.get("num", 0)
    return f(5)