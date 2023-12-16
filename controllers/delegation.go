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
	Delegation = delegation{}
)

type delegation struct{}

// List godoc
//
//	@Summary		List Delegations
//	@Description	List all the aggregated new delegations
//	@Tags			delegations
//	@Accept			json
//	@Produce		json
//	@Param			year	query		string	false	"Filter results by year"
//	@Success		200		{object}	dto.DataDelegationSliceResponse
//	@Failure		400		{object}	errorhandling.SimpleErrorMessage
//	@Failure		500		{object}	errorhandling.SimpleErrorMessage
//	@Router			/xtz/delegations [get]
func (controller delegation) List(c *gin.Context) {
	var queryParameters dto.ListDelegationsQueryParameters
	err := c.ShouldBind(&queryParameters)
	if err != nil {
		log.Error().Err(err).Msg("Unable to Parse Query Parameters")
		errorMessage := errorhandling.SimpleErrorMessage{Message: err.Error()}
		c.JSON(http.StatusBadRequest, errorMessage)
		c.Abort()

		return
	}
	
	delegations, err := services.Delegation.List(queryParameters.Year)

	if err != nil {
		log.Err(err).Msg("Error listing delegations")
		errorMessage := errorhandling.SimpleErrorMessage{Message: "Error to list delegations"}
		c.JSON(http.StatusInternalServerError, errorMessage)
		c.Abort()

		return
	}
	delegationResponses := new(dto.DelegationResponse).ToDtos(delegations)
	response := dto.DataDelegationSliceResponse{
		Data: delegationResponses,
	}
	c.JSON(http.StatusOK, &response)
	c.Abort()
}
