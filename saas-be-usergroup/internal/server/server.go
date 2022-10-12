package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/gofiber/fiber/v2"
	baseApp "github.com/saas-be-usergroup/internal/app"
	"github.com/saas-be-usergroup/pkg/logger"
	"github.com/saas-be-usergroup/pkg/postgres"
	"github.com/saas-be-usergroup/pkg/viper"
	viperPkg "github.com/spf13/viper"
)

func Run() {

	//load config
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "../..")
	config := &viper.EnvConfig{
		FileName: "config",
		FileType: "yaml",
		Path:     basepath,
	}
	if err := config.ReadConfig(); err != nil {
		log.Fatal(err)
	}

	// //load connection config
	pg, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	// //load connection redis
	// _, err = redis.Connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	zap, err := logger.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	//load fiber
	app := fiber.New(fiber.Config{
		IdleTimeout: 5,
	})
	rh := &baseApp.Handlers{
		Postgres: pg,
		R:        app,
		Logger:   zap,
	}
	rh.SetupRouter()

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":" + viperPkg.GetString("server.port")); err != nil {
			log.Panicf("failed listen into port %v", err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	var _ = <-c // This blocks the main thread until an interrupt is received
	log.Println("gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	pg.Close()
	fmt.Println("services was successful shutdown.")
}
