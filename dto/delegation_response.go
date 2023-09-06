package dto

import (
	"time"

	"github.com/pavva91/tezos-delegation-service/models"
)

type DelegationResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Amount    string    `json:"amount"`
	Delegator string    `json:"delegator"`
	Block     string    `json:"block"`
}

func (dto *DelegationResponse) ToDto(delegation models.Delegation) {
	dto.Timestamp = delegation.Timestamp
	dto.Amount = delegation.Amount
	dto.Delegator = delegation.Delegator
	dto.Block = delegation.Block
}

func (dto DelegationResponse) ToDtos(delegationModels []models.Delegation) (delegationDtos []DelegationResponse) {
	delegationDtos = make([]DelegationResponse, len(delegationModels))
	for i, delegationModel := range delegationModels {
		delegationDtos[i].ToDto(delegationModel)
	}
	return delegationDtos
}
