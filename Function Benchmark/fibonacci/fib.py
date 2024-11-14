import time
import os
import sys

"""
def timer(sleep_time):
    # sleep
    time.sleep(int(sleep_time)/1000)
    os._exit(1)
"""
def fib(n):
    if n<0:
        print("Incorrect input")
    # First Fibonacci number is 0
    elif n==1:
        return 0
    # Second Fibonacci number is 1
    elif n==2:
        return 1
    else:
        return fib(n-1)+fib(n-2)

def f(n):
    """
    f = open("/proc/{pid}/stat".format(pid=os.getpid()), 'r')
    cpu = f.read().split()[-14]
    f.close()
    """
    #cpu = open("/proc/{pid}/stat".format(pid=os.getpid()), 'r').read().split()[-14]
    start = round(time.time(),6)
    #sleep_time = args.get("time","50")
    n = int(n)
    #thread = threading.Thread(target=timer,args=(sleep_time,))
    #thread.start()
    result = fib(int(n))
    #thread.join()
    end = round(time.time(),6)
    return result, end - start

def handler(event, context):
    # n = event.get("num", 0)
    return f(10)