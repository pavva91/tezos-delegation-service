package stubs

import (
	"time"

	"github.com/pavva91/tezos-delegation-service/models"
)

type DelegationRepositoryStub struct {
	ListFn       func() ([]models.Delegation, error)
	ListByYearFn func(time.Time) ([]models.Delegation, error)
	CreateFn     func(*models.Delegation) (*models.Delegation, error)
}

func (stub DelegationRepositoryStub) List() ([]models.Delegation, error) {
	return stub.ListFn()
}

func (stub DelegationRepositoryStub) ListByYear(year time.Time) ([]models.Delegation, error) {
	return stub.ListByYearFn(year)
}

func (stub DelegationRepositoryStub) Create(delegation *models.Delegation) (*models.Delegation, error) {
	return stub.CreateFn(delegation)
}
