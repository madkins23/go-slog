# `flash` Handler

The `flash` package provides a feature-complete,[^1]
reasonably fast implementation of a `slog.Handler`.

After the `sloggy` handler turned out to be slower than desired it was cloned and renamed `flash`.
Then a series of edits were made to improve performance.

## Example

```go
logger := slog.New(flash.NewHandler(os.Stdout, nil))
logger.Info("hello", "count", 3)
```

## Performance Edits

A series of edits were made to `sloggy` to turn it into `flash`.
Some care has been taken to document these edits as they may be
representative of more general performance-enhancing principles in some cases.

### Remove `bytes.Buffer` Usage

Both `sloggy` and `flash` pre-format attributes added to loggers via `WithAttrs`.
This means that information has to be formatted to two destinations:
the prefix/suffix areas associated with handlers as well as
the eventual target for the log data.

For `sloggy` it made sense to write code that could format attributes to an `io.Writer`.
This worked well for the eventual target for the log data, which would always be an `io.Writer`.
For prefix/suffix data it made sense to use `bytes.Buffer` objects which
already implement the `io.Writer` interface.

Profiling indicated that the `bytes.Buffer` code was taking up a lot of CPU cycles,
as well as making a lot of requests for memory.
In addition, there were a lot of cases where Go formatting functions/methods
were configured to append formatted data to a pre-existing `byte` array.

For example, `strconv.AppendBool` takes a `byte` array,
appends the formatted string for a boolean, then returns a `byte` array that results,
whereas `strconv.FormatBool` calls `AppendBool` with an internal array and
returns that `[]byte` as a `string`.
Appending data to an array that is not filled to capacity is quite efficient.
Utilizing `bytes.Buffer` objects, while convenient, is not.

The `flash` edit for this issue changed from a common `io.Writer` interface to a `byte` array interface.
This works well for prefix/suffix data but requires a temporary `[]byte` object be
constructed and then written to the eventual `io.Writer` in `flash.Handler.Handle`.

[Benchmark tests](#benchmark-tests) have been added to compare the performance of the two approaches.
The command to run these tests is:
```
scripts/flash-bench Compose
```
Manual evaluation of this test over different time periods suggests:

* ~40% decrease in execution time
* ~24% decrease in bytes allocated
* ~69% decrease in the number of allocations

### Use Pools to Reduce Memory Allocations

The `sloggy` implementation allocated various blocks of memory and allowed them to be garbage collected.
Early profiling showed a memory allocation function taking a non-trivial amount of time,
and memory allocation numbers were high.

The `flash` implementation uses `sync.Pool` to recycle buffers of common sizes.
Currently, (`2024-03-09`) the following buffers are reused:

* log record output buffers and
* `composer` objects.

Initially there were two other pools:

* Basic fields arrays  
  These arrays were [removed later](#flatten-basic-field-array).
* Source records  
  These records are [no longer allocated on the heap](#use-local-variable-for-source-record).

To support the convenient use of `sync.Pool` there are two variants of
generic pools in file `flash/pools.go`.

[Benchmark tests](#benchmark-tests) have been added to compare the performance of the two approaches.
The command to run these tests is:
```
scripts/flash-bench Memory
```
Manual evaluation of this test over different time periods suggests:

* ~88% decrease in execution time
* ~97% decrease in bytes allocated
* no change in the number of allocations (?!?)

### Flatten Basic Field Array

In `sloggy` it made sense to create a small array of `slog.Attr`,
fill it with the "basic" fields (i.e. `time`, `level`, `msg`, and `source`),
and then compose all of those fields with `composer.addAttributes`.
The alternative, sending them in one at a time using `composer.addAttribute`,
added four more error result tests and seemed cluttered compared to
sending an array of attributes to `composer.addAttributes` all at once
and letting that method handle the individual error results.

In `flash` initially the small arrays were [moved to a `sync.Pool`](#use-pools-to-reduce-memory-allocations).
Later this was flattened out to four `composer.addAttribute` calls
along with the requisite additional error result handling.
This removed the `sync.Pool` for the small arrays, which were no longer required.

[Benchmark tests](#benchmark-tests) have been added to compare the performance of the two approaches.
The command to run these tests is:
```
scripts/flash-bench Basic
```
Manual evaluation of this test over different time periods suggests:

* ~36% decrease in execution time
* ~26% decrease in bytes allocated
* ~29% decrease in the number of allocations

In addition, the number of allocations for a simple log record,
which had been stubbornly `2` for `flash`, dropped to `1`.[^2]

### Use Generalized `Stringer` Interface

Several methods from `composer` that were redundant were removed because
the special cases all implemented `Stringer` which was already covered.
This removed some specialized methods on `composer`,
but probably didn't affect the performance very much.

No benchmark tests were done to compare the performance
as no measurable performance improvement was expected.

### Use Local Variable for Source Record

The `source` record represents the place in the code where the log record was generated.
The data is logged in a group, so placing it in a record makes sense.
This data is only gathered and logged if the `slog.HandlerOptions.AddSource` flag is `true`.

In `sloggy` the `source` record was allocated and returned.
In `flash` initially the small arrays were [moved to a `sync.Pool`](#use-pools-to-reduce-memory-allocations).
This meant allocating the records on the heap.
Later the record was allocated on the stack in `flash.Handler` and then
`flash.composer.loadSource` was called to populate the record.
Allocating on the stack is generally faster than on the heap.

[Benchmark tests](#benchmark-tests) have been added to compare the performance of the two approaches.
The command to run these tests is:
```
scripts/flash-bench Source
```
Manual evaluation of this test over different time periods suggests:

* ~10% decrease in execution time
* no change in bytes allocated
* no change in the number of allocations

### Call Before Visiting

A very small change was to check before invoking `slog.Value.Resolve`:
```go
if attr.Value.Kind() == slog.KindLogValuer {
    attr.Value.Resolve()
}
```

The theory was that a quick check would rule out a method call over a lot of attributes.
An assumption was made that such a quick check was possible but that is debatable.
The call to `slog.Value.Kind` is going to do essentially
the same check that `slog.Value.Resolve` will do right away,
and immediately return the same value if it not a `LogValuer`.
So this edit is basically only avoiding the method call and return overhead.

[Benchmark tests](#benchmark-tests) have been added to compare the performance of the two approaches.
The command to run these tests is:
```
scripts/flash-bench Resolve
```
Manual evaluation of this test over different time periods suggests:

* ~72% decrease in execution time
* no change in bytes allocated
* no change in the number of allocations

Caveats:

* Due to the tiny amount of code involved this is likely `72%` of almost nothing.
* This applies to cases where the attribute is not a `LogValuer`,
  for any `LogValuer` attribute it actually _adds_ time for the additional conditional.
  In all likelihood `LogValuer` attributes will be relatively rare and even when used
  only apply to a small percentage of all attributes, reducing this overhead.

### Benchmark Tests

The benchmark tests mentioned above have nothing to do with the `bench` package.
They were created specifically to verify performance edits.

**Note** that each test is focused on a specific part of the handler code.
That part of the code may not be executed very much,
so a great result for an individual performance edit may not mean much overall.

Test names are of the form `Benchmark<group><specific>` where `<group>`
specifies the related tests (different approaches to be compared) and
`<specific>` represents the approach to be tested.
For example: `BenchmarkSourceLoad` and `BenchmarkSourceNewReuse`
test [using a record for the `source` field data](#use-local-variable-for-source-record).

To run a group of tests:
```
go test -bench=Benchmark<group> -benchtime=<duration> flash/*.go
```
(where the `<duration>` any string that
[`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration) can handle).

If you can run a `bash` shell script you can use:
```
scripts/flash-bench <group>
```
to run tests with 1 second, 5 second, 15 second, and 60 second durations.

[^1]: Feature complete in this case being according to the verification test suite
defined in `go-slog/verify` and the tests and warnings defined therein.

[^2]: If only I knew where that last `1` came from.  :frowning_face: