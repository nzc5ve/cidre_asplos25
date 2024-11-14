import igraph
import time
import sys

def f(n):

    start = round(time.time(),6)
    n = int(n)

    size = int(n)
    graph = igraph.Graph.Barabasi(size, size)
    result = graph.bfs(0)
    end = round(time.time(),6)

def handler(event, context):
    # n = event.get("num", 0)
    return f(1000)