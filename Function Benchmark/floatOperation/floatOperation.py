import time
import sys
import math


def float_operations(n):
    for i in range(0, n):
        sin_i = math.sin(i)
        cos_i = math.cos(i)
        sqrt_i = math.sqrt(i)


def f(n):
    start = round(time.time(),6)
    n = int(n)
    result = float_operations(int(n))
    end = round(time.time(),6)

def handler(event, context):
    # n = event.get("num", 0)
    return f(10)