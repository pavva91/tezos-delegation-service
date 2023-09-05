package server

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pavva91/gin-gorm-rest/config"
	"github.com/rs/zerolog/log"
)

var (
	router = gin.Default()
)

func NewRouter() *gin.Engine {

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS Configs based on SERVER_ENVIRONMENT variable
	switch env := config.ServerConfigValues.Server.Environment; env {
	case "dev":
		corsConfigDev := cors.DefaultConfig()
		corsConfigDev.AllowAllOrigins = true
		corsConfigDev.AllowHeaders = append(corsConfigDev.AllowHeaders, "Authorization")
		router.Use(cors.New(corsConfigDev))
	case "stage":
		log.Info().Msg("TODO: Stage environment Setup, for now allow all CORS")
		corsConfigStage := cors.DefaultConfig()
		corsConfigStage.AllowAllOrigins = true
		corsConfigStage.AllowHeaders = append(corsConfigStage.AllowHeaders, "Authorization")
		router.Use(cors.New(corsConfigStage))
	case "prod":
		corsConfigProd := cors.DefaultConfig()
		corsConfigProd.AllowOrigins = config.ServerConfigValues.Server.CorsAllowedClients
		router.Use(cors.New(corsConfigProd))
	default:
		log.Error().Msg(fmt.Sprintf("Incorrect Dev Environment: %s\nInterrupt execution", env))
		os.Exit(1)
	}

	return router
}

func MapUrls() {
	apiVersion := fmt.Sprintf("/%s/%s", config.ServerConfigValues.Server.ApiPath, config.ServerConfigValues.Server.ApiVersion)
	mapUrls(apiVersion)
}
