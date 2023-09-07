package repositories

import (
	"github.com/pavva91/tezos-delegation-service/db"
	"github.com/pavva91/tezos-delegation-service/models"
)

var (
	DelegationRepository DelegationRepositoryInterface = delegationRepositoryImpl{}
)

type DelegationRepositoryInterface interface {
	List() ([]models.Delegation, error)
	Create(delegation *models.Delegation) (*models.Delegation, error)
}

type delegationRepositoryImpl struct{}

func (repository delegationRepositoryImpl) List() ([]models.Delegation, error) {
	delegations := []models.Delegation{}
	err := db.DbOrm.GetDB().Order("timestamp DESC").Find(&delegations).Error
	if err != nil {
		return nil, err
	}
	return delegations, nil
}

func (repository delegationRepositoryImpl) Create(delegation *models.Delegation) (*models.Delegation, error) {
	err := db.DbOrm.GetDB().Create(&delegation).Error
	if err != nil {
		return nil, err
	}
	return delegation, nil
}

