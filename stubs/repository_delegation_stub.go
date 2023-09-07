package stubs

import "github.com/pavva91/tezos-delegation-service/models"

type DelegationRepositoryStub struct {
	ListFn func() ([]models.Delegation, error)
	CreateFn func(*models.Delegation) (*models.Delegation, error)
}

func (stub DelegationRepositoryStub) List() ([]models.Delegation, error) {
	return stub.ListFn()
}

func (stub DelegationRepositoryStub) Create(delegation *models.Delegation) (*models.Delegation, error) {
	return stub.CreateFn(delegation)
}
