# Using `log/slog` with `gin`

Package `gin` contains utilities for using `slog` with
[`gin-gonic/gin`](https://github.com/gin-gonic/gin).
In particular, this package provides `gin.Writer` which can be used to redirect Gin-internal logging:
```go
import (
    "github.com/gin-gonic/gin"
    ginslog "github.com/madkins23/go-slog/gin"
)

gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{})
gin.DefaultErrorWriter = ginslog.NewWriter(&ginslog.Options{Level: slog.LevelError})
```
Configure this before starting Gin and all the Gin-internal logging
should be redirected to the new `io.Writer` objects.
These objects will parse the Gin-internal logging formats and
use `slog` to do the actual logging, so the log output will all look the same.

The `gin.Writer` objects can further parse the "standard" Gin traffic lines containing:
```
200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"
```
To embed the traffic data at the top level of the log messages:
```go
gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{
	Traffic: ginslog.Traffic{Parse: true, Embed: true},
})
```
To aggregate the traffic data into a group named by `ginslog.DefaultTrafficGroup`:
```go
gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{
	Traffic: ginslog.Traffic{Parse: true},
})
```
To aggregate the traffic data into a group named `"bob"`:
```go
gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{
	Traffic: ginslog.Traffic{Parse: true, Group: "bob"},
})
```
Further options can be found in the code documentation of `go-slog/gin.Options`.

