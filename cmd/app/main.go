package main

import (
	"fmt"
	"log"

	"github.com/restlesswhy/btc-service/config"
	"github.com/restlesswhy/btc-service/internal/server"
	"github.com/restlesswhy/btc-service/pkg/logger"
)

// @title BTC-sercive swagger API
// @version 2.0
// @description BTC service with Hyperledger Fabric implementation.

// @contact.name German Generalov
// @contact.url http://github.com/restlesswhy

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4000
// @BasePath /api/v1/
// @schemes http
func main() {
	log.Println("Starting microservice")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.Named(fmt.Sprintf(`(%s)`, cfg.ServiceName))
	appLogger.Infof("CFG: %+v", cfg)

	if err := server.New(appLogger, cfg, nil).Run(); err != nil {
		appLogger.Fatal(err)
	}
}
