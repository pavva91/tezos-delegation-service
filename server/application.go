package server

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pavva91/gin-gorm-rest/config"
	"github.com/pavva91/gin-gorm-rest/db"
	"github.com/pavva91/gin-gorm-rest/docs"
	"github.com/pavva91/gin-gorm-rest/models"
	"github.com/pavva91/gin-gorm-rest/services"
)

func StartApplication() {
	env := os.Getenv("SERVER_ENVIRONMENT")

	log.Info().Msg(fmt.Sprintf("Running Environment: %s", env))

	switch env {
	case "dev":
		err := cleanenv.ReadConfig("./config/dev-config.yml", &config.ServerConfigValues)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	case "stage":
		err := cleanenv.ReadConfig("./config/stage-config.yml", &config.ServerConfigValues)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	case "prod":
		err := cleanenv.ReadConfig("./config/prod-config.yml", &config.ServerConfigValues)
		if err != nil {
			log.Err(err).Msg(err.Error())
		}
	default:
		log.Error().Msg(fmt.Sprintf("Incorrect Dev Environment: %s\nInterrupt execution", env))
		os.Exit(1)
	}

	// Set Swagger Info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "Insert here REST API Description"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.BasePath = fmt.Sprintf("/%s/%s", config.ServerConfigValues.Server.ApiPath, config.ServerConfigValues.Server.ApiVersion)
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.ServerConfigValues.Server.Host, config.ServerConfigValues.Server.Port)
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Connect to DB
	db.DbOrm.ConnectToDB(config.ServerConfigValues)
	db.DbOrm.GetDB().AutoMigrate(&models.Delegation{})

	inititalizeDb()

	go services.DelegationService.PollDelegations()

	// Create Router
	router := NewRouter()

	MapUrls()

	// Start Server
	router.Run(":" + config.ServerConfigValues.Server.Port)

}

func inititalizeDb() {
	
}
