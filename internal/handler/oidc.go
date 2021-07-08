package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	oidc "github.com/sayan-biswas/oidc/pkg/gin-oidc"
)

func Login(context *gin.Context) {
	OIDC := oidc.Get(context)
	if _, isRedirect := context.GetQuery("code"); isRedirect {
		OIDC.RedirectHandler()
	} else {
		OIDC.LoginHandler()
	}
}

func Logout(context *gin.Context) {
	OIDC := oidc.Get(context)
	OIDC.LogoutHandler()
}

func Check(context *gin.Context) {
	OIDC := oidc.Get(context)
	if !OIDC.IsAuthenticated() {
		context.Status(http.StatusUnauthorized)
		return
	}
	session := sessions.Default(context)
	context.Header("User", session.Get("user").(string))
	context.JSON(http.StatusOK, gin.H{
		"userId": session.Get("user").(string),
	})
}

func GetTokens(context *gin.Context) {
	OIDC := oidc.Get(context)
	if !OIDC.IsAuthenticated() {
		context.Next()
		return
	}
	session := sessions.Default(context)
	context.JSON(http.StatusOK, gin.H{
		"AccessToken":  session.Get("access_token"),
		"RefreshToken": session.Get("refresh_token"),
		"IDToken":      session.Get("id_token"),
	})
}

func GetUserData(context *gin.Context) {
	OIDC := oidc.Get(context)
	if !OIDC.IsAuthenticated() {
		context.Next()
		return
	}
	session := sessions.Default(context)
	context.JSON(http.StatusOK, gin.H{
		"User": session.Get("user").(string),
	})
}
