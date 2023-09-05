package stubs

import "github.com/pavva91/gin-gorm-rest/models"

type DelegationRepositoryStub struct {
	ListDelegationsFn func() ([]models.Delegation, error)
}

func (stub DelegationRepositoryStub) ListDelegations() ([]models.Delegation, error) {
	return stub.ListDelegationsFn()
}
