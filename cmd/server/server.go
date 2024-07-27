package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time_tracker/internal/handler"
	"time_tracker/internal/service"
)

func StartServer(services *service.Service) {
	handlers := handler.NewHandler(services)
	app := fiber.New()

	handlers.Routes(app)

	host := viper.GetString("HTTP_HOST")
	port := viper.GetString("HTTP_PORT")

	logrus.Info("Starting HTTP server...")
	err := app.Listen(host + ":" + port)
	if err != nil {
		logrus.Fatal(err)
	}
}
