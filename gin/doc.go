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
package gin
