package repositories

import (
	"time"

	"github.com/pavva91/tezos-delegation-service/db"
	"github.com/pavva91/tezos-delegation-service/models"
)

var (
	DelegationRepository Delegation = delegation{}
)

type Delegation interface {
	List() ([]models.Delegation, error)
	ListByYear(year time.Time) ([]models.Delegation, error)
	Create(d *models.Delegation) error
}

type delegation struct{}

func (r delegation) List() ([]models.Delegation, error) {
	delegations := []models.Delegation{}
	err := db.DbOrm.GetDB().Order("timestamp DESC").Find(&delegations).Error
	if err != nil {
		return nil, err
	}
	return delegations, nil
}

func (r delegation) ListByYear(year time.Time) ([]models.Delegation, error) {
	delegations := []models.Delegation{}
	err := db.DbOrm.GetDB().Where("EXTRACT(YEAR FROM timestamp) = ?", year.Year()).Order("timestamp DESC").Find(&delegations).Error
	if err != nil {
		return nil, err
	}
	return delegations, nil
}

func (r delegation) Create(d *models.Delegation) error {
	return db.DbOrm.GetDB().Create(&d).Error
}
