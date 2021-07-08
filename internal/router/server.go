package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/handler"
)

func ServerRouter(router *gin.Engine) {
	user, ok := os.LookupEnv("SERVER_USER")
	if !ok {
		user = "admin"
	}
	routerGroup := router.Group("/server")
	routerGroup.Use(gin.BasicAuth(gin.Accounts{
		user: os.Getenv("SERVER_PASSWORD"),
	}))
	{
		routerGroup.GET("/status", handler.GetStats)
		routerGroup.GET("/configuration", handler.GetLogLevel)
		routerGroup.PATCH("/configuration", handler.SetLogLevel)
	}
}
