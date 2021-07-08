package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/handler"
	"github.com/sayan-biswas/oidc/internal/middleware"
)

func UserRouter(router *gin.Engine) {
	routerGroup := router.Group("/user")
	routerGroup.Use(middleware.CheckOIDC())
	{
		routerGroup.GET("/data", handler.GetUserData)
		routerGroup.GET("/tokens", handler.GetTokens)
	}
}
