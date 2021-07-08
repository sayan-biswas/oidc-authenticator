package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RootRouter(router *gin.Engine) {
	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"Application": "OIDC Authenticator",
			"Version":     "1.0.0",
			"Repository":  "github.com/sayan-biswas/oidc",
			"E-Mail":      "sayan-biswas@live.com",
		})
	})
}
