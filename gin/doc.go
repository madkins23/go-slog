// Package gin contains utilities for using log/slog with gin-gonic/gin.
// In particular, provides gin.Writer which can be used to redirect Gin-internal logging:
//
//	gin.DefaultWriter = gin.NewWriter(slog.LevelInfo, NoTraffic)
//	gin.DefaultErrorWriter = gin.NewWriter(slog.LevelError, NoTraffic)
package gin
