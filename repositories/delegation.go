package repositories

import (
	"github.com/pavva91/tezos-delegation-service/db"
	"github.com/pavva91/tezos-delegation-service/models"
)

var (
	DelegationRepository DelegationRepositoryInterface = delegationRepositoryImpl{}
)

type DelegationRepositoryInterface interface {
	ListDelegations() ([]models.Delegation, error)
}

type delegationRepositoryImpl struct{}

func (repository delegationRepositoryImpl) ListDelegations() ([]models.Delegation, error) {
	delegations := []models.Delegation{}
	err := db.DbOrm.GetDB().Find(&delegations).Error
	if err != nil {
		return nil, err
	}
	return delegations, nil
}
