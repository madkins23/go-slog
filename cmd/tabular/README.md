# `tabular`

The [`tabular`](../cmd/server/server.go) application consumes the output saved in a temporary directory by
`scripts/bench` and `scripts/verify` and prints various benchmark tables
and a list of warnings by handler for both benchmark and verification test runs.

## Running `tabular`

Usage:

	go run cmd/tabular/tabular.go [flags]

The flags are:

	-bench string
	    Load benchmark data from path (optional)
	-language value
	    One or more language tags to be tried, defaults to US English.
	-useWarnings
	    Show warning instead of known errors
	-verify string
	    Load verification data from path (optional)

Script:

The `scripts/tabulate` script will run `tabular`,
taking input from temporary files created by `scripts/verify` and `scripts/bench`.

Output:

```
Benchmark Attributes
  ╔══════════════════════╦═════════════╤═══════════════╤═════════════╤═════════════╤═════════════════╗
  ║ Handler              ║        Runs │         Ns/Op │   Allocs/Op │    Bytes/Op │          GB/Sec ║
  ╠══════════════════════╬═════════════╪═══════════════╪═════════════╪═════════════╪═════════════════╣
  ║ chanchal/zap         ║   1,228,219 │      1,041.00 │           5 │         418 │      470,809.52 ║
  ║ phsym/zerolog        ║   1,690,441 │        725.10 │           2 │         272 │      946,553.07 ║
  ║ samber/logrus        ║      50,173 │     25,886.00 │          90 │       8,519 │          810.17 ║
  ║ samber/zap           ║     215,511 │      6,005.00 │          46 │       6,649 │       14,318.54 ║
  ║ samber/zerolog       ║     224,841 │      5,123.00 │          54 │       4,837 │       17,820.35 ║
  ║ slog/json            ║     751,104 │      1,369.00 │           6 │         473 │      234,700.62 ║
  ╚══════════════════════╩═════════════╧═══════════════╧═════════════╧═════════════╧═════════════════╝

  ...[tables for other benchmark tests]...

Warnings for slog/json:
  Suggested
     2 [Duplicates] Duplicate field(s) found
         Verify:AttributeDuplicate: map[alpha:2 charlie:3]
         Verify:AttributeWithDuplicate: map[alpha:2 charlie:3]

  ...[warnings for other handlers]...
```
