package main

import (
	"fmt"

	"capstone-qris/internal/config"
)

// @title           QRIS Payment API
// @version         1.0
// @description     High-performance QRIS payment microservice with Redis caching and optimistic locking.
// @host      localhost:3000
// @BasePath  /
// @securityDefinitions.apikey  X-Client-Id
// @in                          header
// @name                        X-Client-Id
// @securityDefinitions.apikey  X-Client-Key
// @in                          header
// @name                        X-Client-Key
func main() {
	appCfg := config.LoadConfig()
	log := config.NewLogger(appCfg)
	db := config.NewDatabase(appCfg, log)
	rdb := config.NewRedis(appCfg, log)
	val := config.NewValidator()
	app := config.NewFiber(appCfg)

	config.Bootstrap(&config.BootstrapConfig{
		DB:        db,
		Redis:     rdb,
		App:       app,
		Log:       log,
		Validator: val,
		Config:    appCfg,
	})

	port := appCfg.Web.Port
	log.Infof("Starting QRIS Payment API on port %d", port)

	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.WithError(err).Fatal("server exited with error")
	}
}
