package stubs

import "github.com/pavva91/gin-gorm-rest/models"

type DelegationServiceStub struct {
	ListDelegationsFn func() ([]models.Delegation, error)
}

func (stub DelegationServiceStub) ListAllDelegations() ([]models.Delegation, error) {
	return stub.ListDelegationsFn()
}
