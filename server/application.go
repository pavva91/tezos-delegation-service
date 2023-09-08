package server

import (
	"fmt"
	"sync"

	"github.com/pavva91/tezos-delegation-service/config"
	"github.com/pavva91/tezos-delegation-service/db"
	"github.com/pavva91/tezos-delegation-service/docs"
	"github.com/pavva91/tezos-delegation-service/models"
	"github.com/pavva91/tezos-delegation-service/services"
	"github.com/rs/zerolog/log"

	"os"

	"github.com/ilyakaznacheev/cleanenv"
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

	rwMutex := &sync.RWMutex{}
	stopOnError := make(chan bool)
	errorCh := make(chan error)

	go services.DelegationService.PollDelegations(config.ServerConfigValues.ApiDelegations.PollPeriodInSeconds, config.ServerConfigValues.ApiDelegations.Endpoint, rwMutex, stopOnError, errorCh)

	// Create Router
	router := NewRouter()

	MapUrls()

	// Start Server
	router.Run(":" + config.ServerConfigValues.Server.Port)
}

func inititalizeDb() {

}
