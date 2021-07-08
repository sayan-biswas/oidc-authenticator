package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/configurator"
	"github.com/sayan-biswas/oidc/internal/logger"
	"go.uber.org/zap"
)

func Session(config *configurator.Configuration) (handler gin.HandlerFunc) {

	var log = logger.Logger

	options := sessions.Options{
		Path:     "/",
		MaxAge:   config.Session.MaxAge,
		HttpOnly: config.Session.HttpOnly,
		Secure:   config.Session.Secure,
	}

	if config.Session.Store == "redis" {
		store, storeError := redis.NewStore(
			config.Redis.Connections,
			config.Redis.Protocol,
			config.Redis.Host+":"+config.Redis.Port,
			config.Redis.Password,
			[]byte(config.Session.Secret),
		)
		if storeError != nil {
			log.Error("Error connecting to Redis", zap.Error(storeError))
		}
		store.Options(options)
		handler = sessions.Sessions("oidc-session", store)
	}

	if config.Session.Store == "cookie" {
		store := cookie.NewStore([]byte(config.Session.Secret))
		store.Options(options)
		handler = sessions.Sessions("oidc-session", store)
	}

	return
}
