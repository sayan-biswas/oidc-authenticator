package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/configurator"
	"github.com/sayan-biswas/oidc/internal/logger"
	"go.uber.org/zap"
)

func Logger(config *configurator.Configuration) gin.HandlerFunc {
	var log = logger.Logger
	return func(context *gin.Context) {
		start := time.Now()
		context.Next()
		end := time.Now()

		if len(context.Errors) > 0 {
			for _, contextError := range context.Errors.Errors() {
				log.Error("Context error", zap.String("error", contextError))
			}
		} else {
			log.Debug(strconv.Itoa(context.Writer.Status()),
				zap.Int("status", context.Writer.Status()),
				zap.String("method", context.Request.Method),
				zap.String("path", context.Request.URL.Path),
				zap.String("query", context.Request.URL.RawQuery),
				zap.String("host", context.Request.Host),
				zap.String("referer", context.Request.Referer()),
				zap.String("ip", context.ClientIP()),
				zap.Duration("latency", end.Sub(start)),
			)
		}
	}
}
