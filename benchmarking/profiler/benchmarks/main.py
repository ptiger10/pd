import pandas as pd
import datetime
import json
import sys


def main():
    # Start tests
    results = {
        "100k": {
            "sum": sumTest(),
            "mean": meanTest(),
            "median": medianTest(),
            "min": minTest(),
            "max": maxTest(),
            "std": stdTest(),
            "readCSVSum": readCSVSumTest(),
            }
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
            # print("{}: ".format(fn.__name__).ljust(15), end='')
            # print(
            # "{}".format(round(duration*1000000000), 0).rjust(10), " ns/op")
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


df = pd.read_csv('RandomNumbers.csv')


@timer(1000)
def sumTest():
    s = df.sum()
    assert round(s.iloc[0], 2) == 50408.63


@timer(1000)
def meanTest():
    s = df.mean()
    assert round(s.iloc[0], 2) == 0.5


@timer(1000)
def minTest():
    s = df.min()
    assert round(s.iloc[0], 2) == 0.0


@timer(1000)
def maxTest():
    s = df.max()
    assert round(s.iloc[0], 2) == 1.0


@timer(1000)
def stdTest():
    s = df.std()
    assert round(s.iloc[0], 2) == 0.29


@timer(100)
def medianTest():
    s = df.median()
    assert round(s.iloc[0], 2) == 0.5


@timer(50)
def readCSVSumTest():
    df = pd.read_csv('RandomNumbers.csv')
    s = df.sum()
    assert round(s.iloc[0], 2) == 50408.63


if __name__ == "__main__":
    main()
