package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/oidc/internal/configurator"
	"github.com/sayan-biswas/oidc/internal/logger"
	"github.com/sayan-biswas/oidc/internal/middleware"
	"github.com/sayan-biswas/oidc/internal/router"
	"go.uber.org/zap"
)

var log = logger.Logger

type Server struct {
	httpServer   *http.Server
	configurator *configurator.Configurator
}

func New(configFile string) *Server {
	return &Server{
		httpServer:   new(http.Server),
		configurator: configurator.New(configFile),
	}
}

func (server *Server) Start() (serverError error) {
	defer func() {
		log.Sync()
		if panic := recover(); panic != nil {
			log.Error("Server error", zap.String("error", panic.(string)))
		}
	}()
	config := server.configurator.GetConfiguration()

	// Set Gin Mode
	if config.Server.Debug {
		gin.SetMode(gin.DebugMode)
		logger.SetLevel("debug")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize GIN
	ginEngine := gin.New()
	ginEngine.AppEngine = true

	// Middleware chain
	{
		ginEngine.Use(middleware.Logger(config))
		ginEngine.Use(middleware.Recovery(config))
		ginEngine.Use(middleware.CORS(config))
		ginEngine.Use(middleware.Session(config))
		ginEngine.Use(middleware.OIDC(config))
	}

	//Router chain
	{
		router.RootRouter(ginEngine)
		router.OIDCRouter(ginEngine)
		router.UserRouter(ginEngine)
		router.ServerRouter(ginEngine)
	}

	// Start server
	server.httpServer.Addr = config.Server.Host + ":" + strconv.Itoa(config.Server.Port)
	server.httpServer.Handler = ginEngine

	log.Info(
		"Server starting",
		zap.String("host",
			config.Server.Host),
		zap.Int("port", config.Server.Port),
		zap.Bool("tls", config.Server.TLS),
	)

	if config.Server.TLS {
		serverError = server.httpServer.ListenAndServeTLS(config.Server.Certificate, config.Server.PrivateKey)
	} else {
		serverError = server.httpServer.ListenAndServe()
	}

	if serverError != nil && serverError != http.ErrServerClosed {
		log.Error("Error starting server", zap.Error(serverError))
	}

	return
}

func (server *Server) Stop() (serverError error) {
	log.Info("Shutting down server")
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if serverError = server.httpServer.Shutdown(context); serverError != nil {
		log.Error("Server forced to shutdown", zap.Error(serverError))
	}
	log.Info("Server successfully shutdown")
	return
}
