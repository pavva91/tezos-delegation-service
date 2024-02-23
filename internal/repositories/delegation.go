package repositories

import (
	"time"

	"github.com/pavva91/tezos-delegation-service/internal/db"
	"github.com/pavva91/tezos-delegation-service/internal/models"
)

var (
	Delegation DelegationRepositer = delegation{}
)

type DelegationRepositer interface {
	List() ([]models.Delegation, error)
	ListByYear(year time.Time) ([]models.Delegation, error)
	Create(d *models.Delegation) error
}

type delegation struct{}

func (r delegation) List() ([]models.Delegation, error) {
	delegations := []models.Delegation{}
	err := db.ORM.GetDB().Order("timestamp DESC").Find(&delegations).Error
	if err != nil {
		return nil, err
	}
	return delegations, nil
}

func (r delegation) ListByYear(year time.Time) ([]models.Delegation, error) {
	delegations := []models.Delegation{}
	err := db.ORM.GetDB().Where("EXTRACT(YEAR FROM timestamp) = ?", year.Year()).Order("timestamp DESC").Find(&delegations).Error
	if err != nil {
		return nil, err
	}
	return delegations, nil
}

func (r delegation) Create(d *models.Delegation) error {
	return db.ORM.GetDB().Create(&d).Error
}
