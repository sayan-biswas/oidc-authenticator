package flags

import "github.com/spf13/pflag"

func init() {
	pflag.Bool("debug", false, "Debug Mode")
	pflag.String("logLevel", "info", "Logging Level - (debug, info, warn, error, fatal, panic)")
	pflag.String("logOutput", "console", "Log Output - (json, console)")
	pflag.String("config", "oidc.yaml", "Configuration File")
	pflag.Int("port", 4600, "Server Port")
	pflag.String("host", "localhost", "Server Hostname")
	pflag.Bool("tls", false, "Enable TLS")
	pflag.Parse()
}
