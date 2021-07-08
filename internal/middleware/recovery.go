package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/configurator"
	"github.com/sayan-biswas/oidc/internal/logger"
	"go.uber.org/zap"
)

func Recovery(config *configurator.Configuration) gin.HandlerFunc {
	var log = logger.Logger
	return func(context *gin.Context) {
		defer func() {
			if panic := recover(); panic != nil {
				var brokenPipe bool
				if opError, ok := panic.(*net.OpError); ok {
					if syscallError, ok := opError.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(syscallError.Error()), "broken pipe") || strings.Contains(strings.ToLower(syscallError.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(context.Request, false)
				if brokenPipe {
					log.Error(context.Request.URL.Path,
						zap.Any("error", panic),
						zap.String("request", string(httpRequest)),
					)
					context.Error(panic.(error))
					context.Abort()
					return
				}

				if config.Server.Debug {
					log.Error("Panic",
						zap.Time("time", time.Now()),
						zap.Any("error", panic),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					log.Error("Panic",
						zap.Time("time", time.Now()),
						zap.Any("error", panic),
						zap.String("request", string(httpRequest)),
					)
				}
				context.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		context.Next()
	}
}
