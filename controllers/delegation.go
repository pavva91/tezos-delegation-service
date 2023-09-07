package controllers

import (
	"net/http"
	"time"

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
	var queryParameters ListDelegationsQueryParameters
	err := context.ShouldBind(&queryParameters)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Parse Query Parameters")
		errorMessage := errorhandling.SimpleErrorMessage{Message: err.Error()}
		context.JSON(http.StatusBadRequest, errorMessage)
		context.Abort()
		return
	}
	delegations, err := services.DelegationService.ListDelegations(queryParameters.Year)
	if err != nil {
		log.Err(err).Msg("Error listing delegations")
		errorMessage := errorhandling.SimpleErrorMessage{Message: "Error to list delegations"}
		context.JSON(http.StatusInternalServerError, errorMessage)
		context.Abort()
		return
	}
	// TODO: Add "data" into the response
	delegationResponses := new(dto.DelegationResponse).ToDtos(delegations)
	response := dto.DataDelegationSliceResponse{
		Data: delegationResponses,
	}
	context.JSON(http.StatusOK, &response)
	context.Abort()
	return
}

type ListDelegationsQueryParameters struct {
	Year time.Time `form:"year" time_format:"2006"`
}
