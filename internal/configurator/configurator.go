package configurator

import (
	"path"
	"strings"

	"github.com/sayan-biswas/oidc/internal/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log = logger.Logger

type Configurator struct {
	*viper.Viper
}

func New(configFile string) (configurator *Configurator) {
	configurator = &Configurator{
		viper.New(),
	}
	configurator.Configure(configFile)
	return
}

func (configurator *Configurator) GetConfiguration() (config *Configuration) {
	if configError := configurator.Unmarshal(&config); configError != nil {
		log.Error("Invalid config file", zap.Error(configError))
	}
	return
}

func (configurator *Configurator) Configure(configFile string) {

	log.Info("Using config file", zap.String("configFile", configFile))

	// Viper flags
	configurator.AllowEmptyEnv(false)
	configurator.SetConfigType("yaml")

	// Change this to take path from command line or root
	dir, file := path.Split(configFile)
	if len(dir) == 0 {
		configurator.AddConfigPath(".")
	} else {
		configurator.AddConfigPath(dir)
	}
	configurator.SetConfigName(strings.TrimSuffix(file, path.Ext(file)))

	// Load server configurations
	if configError := configurator.ReadInConfig(); configError != nil {
		log.Error("Error reading config file", zap.Error(configError))
	}

	// Bind Flags
	configurator.BindPFlag("Server.Debug", pflag.Lookup("debug"))
	configurator.BindPFlag("Server.Host", pflag.Lookup("host"))
	configurator.BindPFlag("Server.Port", pflag.Lookup("port"))
	configurator.BindPFlag("Server.TLS", pflag.Lookup("tls"))

	// Defaults
	configurator.SetDefault("Server.Host", "localhost")
	configurator.SetDefault("Server.Port", 4600)
	configurator.SetDefault("Server.Debug", false)
	configurator.SetDefault("Server.TLS", false)

	configurator.SetDefault("Session.Store", "cookie")
	configurator.SetDefault("Session.Secret", "6251655468576D5A")
	configurator.SetDefault("Session.MaxAge", 1800)
	configurator.SetDefault("Session.HttpOnly", true)
	configurator.SetDefault("Session.Secure", false)

	configurator.SetDefault("Redis.Host", "localhost")
	configurator.SetDefault("Redis.Port", 6379)
	configurator.SetDefault("Redis.Password", "")
	configurator.SetDefault("Redis.Protocol", "tcp")
	configurator.SetDefault("Redis.Connections", 10)

	configurator.SetDefault("OIDC.Provider", "")
	configurator.SetDefault("OIDC.ClientID", "")
	configurator.SetDefault("OIDC.ClientSecret", "")
	configurator.SetDefault("OIDC.RedirectURL", "")
	configurator.SetDefault("OIDC.Scopes", []string{"profile"})

	configurator.SetDefault("CORS.AllowOrigins", []string{"*"})
	configurator.SetDefault("CORS.AllowMethods", []string{"GET", "POST", "HEAD", "OPTIONS"})
	configurator.SetDefault("CORS.AllowHeaders", []string{"Origin", "Content-Length", "Content-Type"})
	configurator.SetDefault("CORS.ExposeHeaders", []string{"Content-Length", "Content-Type"})
	configurator.SetDefault("CORS.AllowCredentials", true)
	configurator.SetDefault("CORS.MaxAge", 86400)

	// Environment Variables
	configurator.BindEnv("Server.Host", "SERVER_HOST")
	configurator.BindEnv("Server.Port", "SERVER_PORT")
	configurator.BindEnv("Server.Debug", "SERVER_DEBUG")
	configurator.BindEnv("Server.TLS", "SERVER_TLS")
	configurator.BindEnv("Server.Certificate", "SERVER_CERTIFICATE")
	configurator.BindEnv("Server.PrivateKey", "SERVER_PRIVATE_KEY")

	configurator.BindEnv("Session.Store", "SESSION_STORE")
	configurator.BindEnv("Session.Secret", "SESSION_SECRET")
	configurator.BindEnv("Session.MaxAge", "SESSION_MAX_AGE")
	configurator.BindEnv("Session.HttpOnly", "SESSION_HTTP_ONLY")
	configurator.BindEnv("Session.Secure", "SESSION_SECURE")

	configurator.BindEnv("Redis.Host", "REDIS_HOST")
	configurator.BindEnv("Redis.Port", "REDIS_PORT")
	configurator.BindEnv("Redis.Password", "REDIS_PASSWORD")
	configurator.BindEnv("Redis.Protocol", "REDIS_PROTOCOL")
	configurator.BindEnv("Redis.Connections", "REDIS_CONNECTIONS")

	configurator.BindEnv("OIDC.Provider", "OIDC_PROVIDER")
	configurator.BindEnv("OIDC.ClientID", "OIDC_CLIENT_ID")
	configurator.BindEnv("OIDC.ClientSecret", "OIDC_CLIENT_SECRET")
	configurator.BindEnv("OIDC.RedirectURL", "OIDC_REDIRECT_URL")
	configurator.BindEnv("OIDC.Scopes", "OIDC_SCOPE")

	configurator.BindEnv("CORS.AllowOrigins", "CORS_ALLOW_ORIGINS")
	configurator.BindEnv("CORS.AllowMethods", "CORS_ALLOW_METHODS")
	configurator.BindEnv("CORS.AllowHeaders", "CORS_ALLOW_HEADERS")
	configurator.BindEnv("CORS.ExposeHeaders", "CORS_EXPOSE_HEADERS")
	configurator.BindEnv("CORS.AllowCredentials", "CORS_ALLOW_CREDENTIALS")
	configurator.BindEnv("CORS.MaxAge", "CORS_MAX_AGE")

}
