package main

//go:generate swag init

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time_tracker/cmd/server"
	"time_tracker/configs"
	_ "time_tracker/docs"
	"time_tracker/internal/repository"
	"time_tracker/internal/service"
)

// @title			Time Tracker API
// @version		1.0
// @description	This is a simple API for time tracking.
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	rovsh.dev@gmail.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:80
// @BasePath		/
func main() {
	configs.LoadMain()
	dbConfig := configs.GetDBConfig()
	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		logrus.Fatalf("DB connection failed: %s", err.Error())
	}
	httpClient := fiber.AcquireClient()
	externalApiUrl := viper.GetString("PEOPLE_HOST")
	repos := repository.NewRepository(db, httpClient, externalApiUrl)
	services := service.NewService(repos)
	server.StartServer(services)
}
