import pandas as pd
import datetime
import json
import sys
import os


def main():
    # Start tests
    results = {
        "100k": {
            "sum": sumTest(),
            # "sumx10": sumTest100k10x(),
            # "readCSVSum10x": readCSVSumTest10x(),
            "mean": meanTest(),
            "min": minTest(),
            "max": maxTest(),
            "std": stdTest(),
            # "readCSVSum": readCSVSumTest(),
            },
        "500k": {
            "sum2": sumTest500(),
        #     "mean2": meanTest500(),
        },
        # "5m": {
        #     "sum": sumTest5m(),
        # }
    }
    json.dump(results, sys.stdout)


# timer computes the average duration across n tests
# returns the duration as string and nanoseconds
def timer(n):
    def decorator(fn):
        def wrapper(*args, **kwargs):
            times = []
            for i in range(n):
                start = datetime.datetime.now()
                fn(*args, **kwargs)
                end = datetime.datetime.now()
                duration = (end-start).total_seconds()
                times.append(duration)
            duration = sum(times)/len(times)
            ns = 1000000000
            mcs = 1000000
            ms = 1000
            if duration * mcs < 1:
                speed = "{:.1f}ns".format(duration*ns)
            if duration * ms < 1:
                speed = "{:.1f}μs".format(duration*mcs)
            elif duration < 1:
                speed = "{:.1f}ms".format(duration*ms)
            else:
                speed = "{:.1f}s".format(duration)
            return speed, int(duration*ns)
        return wrapper
    return decorator


def get_filepath(s):
    basename = files[s]
    thisFile = sys.argv[0]
    path = os.path.join(os.path.dirname(thisFile), basename)
    return path


files = {
    '100k': '../dataRandom100k1Col.csv',
    '100k10x': '../dataRandom100k10Col.csv',
    '500k': '../dataRandom500k2Col.csv',
    '5m': '../dataRandom5m1Col.csv',
}
df100 = pd.read_csv(get_filepath('100k'))
df100k10x = pd.read_csv(get_filepath('100k10x'))
df500 = pd.read_csv(get_filepath('500k'))
# df5m = pd.read_csv(get_filepath('5m'))


@timer(1000)
def sumTest():
    s = df100.sum()
    assert round(s.iloc[0], 2) == 50408.63


@timer(100)
def sumTest100k10x():
    s = df100k10x.sum()
    assert round(s.iloc[0], 2) == 50408.63


@timer(100)
def sumTest500():
    s = df500.sum()
    assert round(s.iloc[0], 2) == 130598.19


# @timer(20)
# def sumTest5m():
#     s = df5m.sum()
#     assert round(s.iloc[0], 2) == 2520431.67


@timer(1000)
def meanTest():
    s = df100.mean()
    assert round(s.iloc[0], 2) == 0.5


@timer(100)
def meanTest500():
    s = df500.mean()
    assert round(s.iloc[0], 2) == 0.26


@timer(1000)
def minTest():
    s = df100.min()
    assert round(s.iloc[0], 2) == 0.0


@timer(1000)
def maxTest():
    s = df100.max()
    assert round(s.iloc[0], 2) == 1.0


@timer(1000)
def stdTest():
    s = df100.std()
    assert round(s.iloc[0], 2) == 0.29


@timer(100)
def medianTest():
    s = df100.median()
    assert round(s.iloc[0], 2) == 0.5


@timer(50)
def readCSVSumTest():
    df = pd.read_csv(get_filepath('100k'))
    s = df.sum()
    assert round(s.iloc[0], 2) == 50408.63


@timer(20)
def readCSVSumTest10x():
    df = pd.read_csv(get_filepath('100k10x'))
    s = df.sum()
    assert round(s.iloc[0], 2) == 50408.63


if __name__ == "__main__":
    main()
