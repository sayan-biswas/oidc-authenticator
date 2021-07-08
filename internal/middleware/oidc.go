package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/configurator"
	"github.com/sayan-biswas/oidc/internal/logger"
	oidc "github.com/sayan-biswas/oidc/pkg/gin-oidc"
	"go.uber.org/zap"
)

func OIDC(config *configurator.Configuration) gin.HandlerFunc {

	var log = logger.Logger

	OIDC, OIDCError := oidc.New(&oidc.Configuration{
		Provider:     config.OIDC.Provider,
		ClientID:     config.OIDC.ClientID,
		ClientSecret: config.OIDC.ClientSecret,
		RedirectURL:  config.OIDC.RedirectURL,
		Scopes:       config.OIDC.Scopes,
	})
	if OIDCError != nil {
		log.Error("Error connecting to OIDC server", zap.Error(OIDCError))
	}
	return func(context *gin.Context) {
		OIDC.SetContext(context)
	}
}

func CheckOIDC() gin.HandlerFunc {
	return func(context *gin.Context) {
		OIDC := oidc.Get(context)
		if !OIDC.IsAuthenticated() {
			context.Redirect(http.StatusSeeOther, "/login")
		}
	}
}
