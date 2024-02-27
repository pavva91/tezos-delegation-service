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

	"github.com/ilyakaznacheev/cleanenv"
)

func MustStartApplication() {
	useEnvVar := os.Getenv("USE_ENVVAR")
	log.Printf("Using envvar value, must be USE_ENVVAR=\"true\" to run with environment variable, otherwise will use config file by default: %s", useEnvVar)

	if useEnvVar == "true" {
		conns, err := strconv.Atoi(os.Getenv("DB_CONNECTIONS"))
		if err != nil {
			log.Panic().Msg(fmt.Sprintf("Incorrect DB connections, must be int: %s\nInterrupt execution", strconv.Itoa(conns)))
		}
		config.ServerConfigValues.Database.Connections = conns
		config.ServerConfigValues.Database.Name = os.Getenv("DB_NAME")
		config.ServerConfigValues.Database.Host = os.Getenv("DB_HOST")
		config.ServerConfigValues.Database.Password = os.Getenv("DB_PASSWORD")
		config.ServerConfigValues.Database.Port = os.Getenv("DB_PORT")
		config.ServerConfigValues.Database.Username = os.Getenv("DB_USERNAME")
		config.ServerConfigValues.Database.Timezone = os.Getenv("DB_TIMEZONE")
		config.ServerConfigValues.Server.Host = os.Getenv("SERVER_HOST")
		config.ServerConfigValues.Server.Port = os.Getenv("SERVER_PORT")
		config.ServerConfigValues.APIDelegations.Endpoint = os.Getenv("DELEGATION_API_ENDPOINT")
		delAPI := os.Getenv("DELEGATION_API_POLL_PERIOD")
		a, err := strconv.Atoi(delAPI)
		config.ServerConfigValues.APIDelegations.PollPeriodInSeconds = uint(a)
		delayAPI := os.Getenv("DELEGATION_API_DELAY_SECONDS")
		b, err := strconv.Atoi(delayAPI)
		config.ServerConfigValues.APIDelegations.DelayLocalTimestampInSeconds = uint(b)
	} else {

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
