package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pavva91/tezos-delegation-service/dto"
	"github.com/pavva91/tezos-delegation-service/errorhandling"
	"github.com/pavva91/tezos-delegation-service/services"
	"github.com/rs/zerolog/log"
)

var (
	DelegationController = eventController{}
)

type eventController struct{}

// ListDelegations godoc
//
//	@Summary		List Delegations
//	@Description	List all the aggregated new delegations
//	@Tags			delegations
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	dto.DelegationResponse
//	@Failure		500	{object}	errorhandling.SimpleErrorMessage
//	@Router			/xtz/delegations [get]
func (controller eventController) ListDelegations(context *gin.Context) {
	delegations, err := services.DelegationService.ListAllDelegations()
	if err != nil {
		log.Err(err).Msg("Error listing delegations")
		errorMessage := errorhandling.SimpleErrorMessage{Message: "Error to list delegations"}
		context.JSON(http.StatusInternalServerError, errorMessage)
		context.Abort()
		return
	}
	// TODO: Add "data" into the response
	delegationResponses := new(dto.DelegationResponse).ToDtos(delegations)
	context.JSON(http.StatusOK, &delegationResponses)
	context.Abort()
	return
}
