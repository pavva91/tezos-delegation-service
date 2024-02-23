package server

import (
	"github.com/pavva91/tezos-delegation-service/internal/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func mapUrls(apiVersion string) {
	api := router.Group(apiVersion)
	{
		api.GET("/xtz/delegations", handlers.Delegation.List)

		healthGroup := api.Group("health")
		{
			healthGroup.GET("", handlers.Health.Status)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
