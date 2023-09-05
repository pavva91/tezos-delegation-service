package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthController struct{}

var (
	HealthController = healthController{}
)

// HealthController godoc
//
//	@Summary		Check Status
//	@Description	Check the status of the REST API
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	message
//	@Router			/health [get]
func (controller healthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "Working!")
}
