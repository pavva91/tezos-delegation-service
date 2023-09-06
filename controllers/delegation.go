package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pavva91/tezos-delegation-service/dto"
	"github.com/pavva91/tezos-delegation-service/services"
)

var (
	DelegationController = eventController{}
)

type eventController struct{}

// ListDelegations godoc
//
//	@Summary		List Delegations
//	@Description	List all the delegations
//	@Tags			delegations
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	dto.Delegation
//	@Router			/delegations [get]
//	@Schemes
func (controller eventController) ListDelegations(context *gin.Context) {
	delegations, err := services.DelegationService.ListAllDelegations()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error to list delegations", "error": err})
		context.Abort()
		return
	}
	// TODO: Add "data" into the response
	delegationResponses := new(dto.DelegationResponse).ToDtos(delegations)
	context.JSON(http.StatusOK, &delegationResponses)
	context.Abort()
	return
}
