package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/handler"
)

func OIDCRouter(router *gin.Engine) {
	routerGroup := router.Group("/")
	{
		routerGroup.GET("/login", handler.Login)
		routerGroup.GET("/logout", handler.Logout)
		routerGroup.GET("/check", handler.Check)
	}
}
