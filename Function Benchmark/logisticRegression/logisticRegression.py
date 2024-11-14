import sys
import numpy as np
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
import pandas as pd
import time
import re

cleanup_re = re.compile('[^a-z]+')


def cleanup(sentence):
    sentence = sentence.lower()
    sentence = cleanup_re.sub(' ', sentence).strip()
    return sentence


def f(n):
    start = round(time.time(),6)
    #sleep_time = args.get("time","50")
    df = pd.read_csv("datapath",format(n))

    df['train'] = df['Text'].apply(cleanup)

    tfidf_vector = TfidfVectorizer(min_df=100).fit(df['train'])

    train = tfidf_vector.transform(df['train'])

    model = LogisticRegression()
    model.fit(train, df['Score'])

    #model_file_path = "/dataout/reviews{}".format(n)
    #joblib.dump(model, model_file_path)
    end = round(time.time(),6)
    #print("{} {}, {}, {}, {}, {}".format(id_n, n, pred, start, end, end-start))
    #print(result)
def handler(event, context):
    # n = event.get("num", 0)
    return f(10)