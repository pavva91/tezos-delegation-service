package server

import (
	"fmt"
	"strconv"

	"github.com/pavva91/tezos-delegation-service/config"
	"github.com/pavva91/tezos-delegation-service/docs"
	"github.com/pavva91/tezos-delegation-service/internal/db"
	"github.com/pavva91/tezos-delegation-service/internal/models"
	"github.com/pavva91/tezos-delegation-service/internal/services"
	"github.com/rs/zerolog/log"

	"os"
)

func MustStartApplication() {
	delAPI := os.Getenv("DELEGATION_API_POLL_PERIOD")
	a, err := strconv.Atoi(delAPI)
	config.ServerConfigValues.APIDelegations.PollPeriodInSeconds = uint(a)
	delayAPI := os.Getenv("DELEGATION_API_DELAY_SECONDS")
	b, err := strconv.Atoi(delayAPI)
	config.ServerConfigValues.APIDelegations.DelayLocalTimestampInSeconds = uint(b)

	env := os.Getenv("SERVER_ENVIRONMENT")

	log.Info().Msg(fmt.Sprintf("Running Environment: %s", env))

	// Set Swagger Info
	docs.SwaggerInfo.Title = "Tezos Delegation Service"
	docs.SwaggerInfo.Description = "Service that gathers new delegations made on the Tezos protocol"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.BasePath = fmt.Sprintf("/%s/%s", config.ServerConfigValues.Server.APIPath, config.ServerConfigValues.Server.APIVersion)
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.ServerConfigValues.Server.Host, config.ServerConfigValues.Server.Port)
	// docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s/%s/%s", config.ServerConfigValues.Server.Host, config.ServerConfigValues.Server.Port, config.ServerConfigValues.Server.APIPath, config.ServerConfigValues.Server.APIVersion)
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Connect to DB
	db.ORM.MustConnectToDB(config.ServerConfigValues)
	err = db.ORM.GetDB().AutoMigrate(&models.Delegation{})
	if err != nil {
		log.Err(err).Msg("Incorrect DB migration")
		os.Exit(1)
	}

	inititalizeDB()

	stopOnError := false
	errorCh := make(chan error)
	quitOnErrorSignalCh := make(chan struct{})

	go func() {
		err := services.Delegation.Poll(config.ServerConfigValues.APIDelegations.PollPeriodInSeconds, config.ServerConfigValues.APIDelegations.Endpoint, stopOnError, errorCh, quitOnErrorSignalCh)
		if err != nil {
			log.Err(err).Msg("")
		}
	}()

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
