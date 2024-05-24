Handler benchmarks are scored based on several metrics on each of various tests.
Metrics are worth different amounts depending on what they are.
The weights applied during this process are shown on the right.

Each combination of handler and test results in a single line of test output:

```
BenchmarkMadkinsFlash/BenchmarkSimple-8  3547497  327.1 ns/op  284.33 MB/s  24 B/op  1 allocs/op
```

From this line we get:
* the handler name (`BenchmarkMadkinsFlash`),
* the test name (`BenchmarkSimple`),
* the number of test runs (`3547497`),
* nanoseconds per operation (`327.1 ns/op`),
* memory bytes allocated per operation (`24 B/op`),
* separate memory allocations per operation (`1 allocs/op`), and
* estimated logging throughput per second (`284.33 MB/s`).

For each handler/test combination (single line or test results)
we use one or more of the following three data items:

* nanoseconds per operation,
* memory bytes allocated per operation, and
* separate memory allocations per operation.

These three items are combined over two steps.
First the test value ranges are acquired:

```
for each test
    for each handler
        for each of the three results described above
            track the highest and lowest value for the test over all handlers
```

Then the test scores are calculated:

```
for each handler
    for each test
        scorePerTest starts at zero
        for each of the three results described above
            convert the value to a fraction of
                the range of values for the test from the previous step
            scorePerTest = scorePerTest + weight(result) * 100.0 * the fraction
        scorePerTest /= sum of weight(result)
    scorePerHandler = average of scorePerTest for handler
```

Where the `weight(result)` comes from the predefined table shown above and to the right.
There is currently no weighting by test, all tests are considered equal.
