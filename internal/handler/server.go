package handler

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/logger"
)

func GetLogLevel(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"LogLevel": logger.Level.String(),
	})
}

func SetLogLevel(context *gin.Context) {
	setError := logger.SetLevel(context.Query("logLevel"))
	if setError != nil {
		context.AbortWithError(http.StatusBadRequest, setError)
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"LogLevel": logger.Level.String(),
	})
}

func GetStats(context *gin.Context) {
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	context.JSON(http.StatusOK, gin.H{
		"Allocation":       fmt.Sprintf("%v MiB", bToMb(memStat.Alloc)),
		"Total Allocation": fmt.Sprintf("%v MiB", bToMb(memStat.TotalAlloc)),
		"System":           fmt.Sprintf("%v MiB", bToMb(memStat.Sys)),
		"NumGC":            fmt.Sprintf("%v", memStat.NumGC),
	})
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
