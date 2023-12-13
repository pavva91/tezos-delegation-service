package server

import (
	"github.com/pavva91/tezos-delegation-service/controllers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func mapUrls(apiVersion string) {
	api := router.Group(apiVersion)
	{
		api.GET("/xtz/delegations", controllers.DelegationController.ListDelegations)

		healthGroup := api.Group("health")
		{
			healthGroup.GET("", controllers.HealthController.Status)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
