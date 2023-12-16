package server

import (
	"fmt"

	"github.com/pavva91/tezos-delegation-service/config"
	"github.com/pavva91/tezos-delegation-service/db"
	"github.com/pavva91/tezos-delegation-service/docs"
	"github.com/pavva91/tezos-delegation-service/models"
	"github.com/pavva91/tezos-delegation-service/services"
	"github.com/rs/zerolog/log"

	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustStartApplication() {
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
	docs.SwaggerInfo.Title = "Tezos Delegation Service"
	docs.SwaggerInfo.Description = "Service that gathers new delegations made on the Tezos protocol"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.BasePath = fmt.Sprintf("/%s/%s", config.ServerConfigValues.Server.APIPath, config.ServerConfigValues.Server.APIVersion)
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.ServerConfigValues.Server.Host, config.ServerConfigValues.Server.Port)
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Connect to DB
	db.ORM.MustConnectToDB(config.ServerConfigValues)
	err := db.ORM.GetDB().AutoMigrate(&models.Delegation{})
	if err != nil {
		log.Err(err).Msg("Incorrect DB migration")
		os.Exit(1)
	}

	inititalizeDB()

	stopOnError := false
	errorCh := make(chan error)
	quitOnErrorSignalCh := make(chan struct{})

	go services.Delegation.Poll(config.ServerConfigValues.APIDelegations.PollPeriodInSeconds, config.ServerConfigValues.APIDelegations.Endpoint, stopOnError, errorCh, quitOnErrorSignalCh)

	// Create Router
	router := NewRouter()

	MapUrls()

	// Start Server
	err = router.Run(":" + config.ServerConfigValues.Server.Port)
	if err != nil {
		log.Err(err).Msg("Error starting router")
		os.Exit(1)
	}
}

func inititalizeDB() {

}
