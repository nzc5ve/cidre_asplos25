import igraph
import time
import sys

def f(n):

    start = round(time.time(),6)

    size = int(n)
    graph = igraph.Graph.Barabasi(size, size)

    result = graph.pagerank()
    end = round(time.time(),6)

def handler(event, context):
    # n = event.get("num", 0)
    return f(1000)