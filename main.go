package main

import (
	// docs "github.com/pavva91/tezos-delegation-service/docs"

	"os"
	"strconv"

	"github.com/pavva91/tezos-delegation-service/server"
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
			isDebug = true
		}
	}
	log.Info().Msg("Debug mode: " + strconv.FormatBool(isDebug))
	// log.Info().Msg("GOMAXPROCS: " + strconv.Itoa(runtime.GOMAXPROCS(0)))
	server.MustStartApplication()
}
