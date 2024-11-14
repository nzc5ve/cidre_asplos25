from random import sample  
from os import path
import time                                                           
import sys
from datetime import datetime    

from jinja2 import Template

SCRIPT_DIR = path.abspath("/data")

def f(n):
    start = round(time.time(),6)
    #sleep_time = args.get("time","50")
    # start timing
    name = "username"
    size = 1000
    start = round(time.time(),6)
    cur_time = datetime.now()
    random_numbers = sample(range(0, 1000000), size)
    template = Template( open(path.join(SCRIPT_DIR, 'templates', 'template.html'), 'r').read())
    html = template.render(username = name, cur_time = cur_time, random_numbers = random_numbers)
    # end timing
    # dump stats
    end = round(time.time(),6)

def handler(event, context):
    # n = event.get("num", 0)
    return f(1)