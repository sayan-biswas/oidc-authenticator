package configurator

type Configuration struct {
	Server  Server
	Session Session
	Redis   Redis
	OIDC    OIDC
	CORS    CORS
}

type Server struct {
	Host        string
	Port        int
	Debug       bool
	TLS         bool
	Certificate string
	PrivateKey  string
}

type Session struct {
	Store    string
	Secret   string
	MaxAge   int
	HttpOnly bool
	Secure   bool
}

type Redis struct {
	Host        string
	Port        string
	Password    string
	Protocol    string
	Connections int
}

type OIDC struct {
	Provider     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type CORS struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}
