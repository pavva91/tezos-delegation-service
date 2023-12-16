package server

import (
	"github.com/pavva91/tezos-delegation-service/controllers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func mapUrls(apiVersion string) {
	api := router.Group(apiVersion)
	{
		api.GET("/xtz/delegations", controllers.Delegation.List)

		healthGroup := api.Group("health")
		{
			healthGroup.GET("", controllers.Health.Status)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
