
import time
from datetime import datetime
import re
import sys
def f(n):

    start = round(time.time(),6)

    data = open("datapath","r").read()
    cleanup_re = re.compile('[a-z]+')
    result = cleanup_re.findall(data)
    print(result)
    end = round(time.time(),6)
    #print("{} {}, {}, {}, {}, {}".format(id_n, n, pred, start, end, end-start))
    #print(result)
def handler(event, context):
    # n = event.get("num", 0)
    return f(1)