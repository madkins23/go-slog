# `server`

The [`server`](../cmd/server/server.go) application consumes the output of
`go test -bench` and serves several web pages that provide the output.

The root page shows links to various test data pages and the warnings:
![The root page shows links to various test data pages and the warnings.](../cmd/server/images/root.png)

Test pages show the same tables as `tabular` plus charts comparing the results:
![Test pages show the same tables as `tabular` plus charts comparing the results.](../cmd/server/images/test.png)

[Recent benchmark data](https://madkins23.github.io/go-slog/index.html).
