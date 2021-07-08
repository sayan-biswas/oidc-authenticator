package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/sayan-biswas/oidc/internal/flags"
	"github.com/sayan-biswas/oidc/internal/logger"
	"github.com/sayan-biswas/oidc/internal/server"
	"github.com/spf13/pflag"
)

//go:embed assets/logo.txt
var logo string
var log = logger.Logger

func main() {

	// Print logo
	fmt.Println(logo)

	// Get command line flags
	configFile := pflag.Lookup("config").Value.String()

	// Create an instance of server
	server := server.New(configFile)

	// Start server
	go server.Start()

	// Stop server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.Stop()

	defer log.Sync()

}
