package repositories

import (
	"time"

	"github.com/pavva91/tezos-delegation-service/db"
	"github.com/pavva91/tezos-delegation-service/models"
)

var (
	DelegationRepository DelegationRepositoryInterface = delegationRepositoryImpl{}
)

type DelegationRepositoryInterface interface {
	List() ([]models.Delegation, error)
	ListByYear(year time.Time) ([]models.Delegation, error)
	Create(delegation *models.Delegation) error
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

func (repository delegationRepositoryImpl) ListByYear(year time.Time) ([]models.Delegation, error) {
	delegations := []models.Delegation{}
	err := db.DbOrm.GetDB().Where("EXTRACT(YEAR FROM timestamp) = ?", year.Year()).Order("timestamp DESC").Find(&delegations).Error
	if err != nil {
		return nil, err
	}
	return delegations, nil
}

func (repository delegationRepositoryImpl) Create(delegation *models.Delegation) error {
	return db.DbOrm.GetDB().Create(&delegation).Error
}
