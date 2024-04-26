// Package gin contains utilities for using log/slog with gin-gonic/gin.
// In particular, this package provides gin.Writer which can be used to redirect Gin-internal logging:
//
//	import (
//	    "github.com/gin-gonic/gin"
//	    ginslog "github.com/madkins23/go-slog/gin"
//	)
//
//	gin.DefaultWriter = ginslog.NewWriter(&ginslog.Options{})
//	gin.DefaultErrorWriter = ginslog.NewWriter(&ginslog.Options{Level: slog.LevelError})
//
// Configure this before starting Gin and all the Gin-internal logging
// should be redirected to the new io.Writer objects.
// These objects will parse the Gin-internal logging formats and
// use log/slog to do the actual logging, so the log output will all look the same.
//
// The io.Writer objects provided by NewWriter can further parse the "standard" Gin traffic lines containing
// messages of the following format:
//
//	200 |  5.529751605s |             ::1 | GET      "/chart.svg?tag=With_Attrs_Attributes&item=MemBytes"
//
// Activate this additional parsing by setting Options.Parse to true.
// This will result the following embedded group:
//
//	gin.code=200 gin.elapsed=5.529751605s gin.client=::1 gin.method=GET gin.url=/chart.svg?tag=With_Attrs_Attributes&item=MemBytes
//
// Additionally setting Options.Embed to true will embed the group fields at the top level:
//
//	code=200 elapsed=5.529751605s client=::1 method=GET url=/chart.svg?tag=With_Attrs_Attributes&item=MemBytes
//
// Further options can be found in the code documentation for gin.Options.
package gin
