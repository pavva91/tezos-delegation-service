package main

import (
	// docs "github.com/pavva91/tezos-delegation-service/docs"

	"fmt"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pavva91/tezos-delegation-service/config"
	"github.com/pavva91/tezos-delegation-service/internal/server"
	"github.com/rs/zerolog/log"
)

// import "github.com/pavva91/tezos-delegation-service/routes"

//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// NOTE: For activating debug
	// <F5> then:
	//	1. 2 (debug with argument)
	//	2. d
	// NOTE: For deactivating debug (terminate session)
	// <F5> then select:
	// - Terminate Session (usually 2 if thread is stopped, 1 if thread is not stopped)
	isDebug := false
	if len(os.Args) == 2 {
		debugArg := os.Args[1]
		if debugArg == "d" || debugArg == "debug" {
			os.Setenv("SERVER_ENVIRONMENT", "dev")
			err := cleanenv.ReadConfig("./config/debug-config.yml", &config.ServerConfigValues)
			if err != nil {
				log.Error().Msg(err.Error())
			}
			isDebug = true
		}
	}
	log.Info().Msg("Debug mode: " + strconv.FormatBool(isDebug))

	if !isDebug {
		setupEnvVars()
	}

	// log.Info().Msg("GOMAXPROCS: " + strconv.Itoa(runtime.GOMAXPROCS(0)))
	server.MustStartApplication()
}

func setupEnvVars() {
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
	config.ServerConfigValues.Server.APIPath = os.Getenv("API_PATH")
	config.ServerConfigValues.Server.APIVersion = os.Getenv("API_VERSION")
}
