package main

import (
	// docs "github.com/pavva91/gin-gorm-rest/docs"

	"os"
	"strconv"

	"github.com/pavva91/gin-gorm-rest/server"
	"github.com/rs/zerolog/log"
)

// import "github.com/pavva91/gin-gorm-rest/routes"

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
	// For activating debug
	// <F5> then:
	//	1. 2 (debug with argument)
	//	2. d
	isDebug := false
	if len(os.Args) == 2 {
		debugArg := os.Args[1]
		if debugArg == "d" || debugArg == "debug" {
			os.Setenv("SERVER_ENVIRONMENT", "dev")
			isDebug = true
		}
	}
	log.Info().Msg("Debug mode: " + strconv.FormatBool(isDebug))
	server.StartApplication()
}
