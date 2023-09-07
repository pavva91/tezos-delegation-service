package stubs

import "github.com/pavva91/tezos-delegation-service/models"

type DelegationServiceStub struct {
	ListDelegationsFn func() ([]models.Delegation, error)
	PollDelegationsFn func(int) error
}

func (stub DelegationServiceStub) ListDelegations() ([]models.Delegation, error) {
	return stub.ListDelegationsFn()
}

func (stub DelegationServiceStub) PollDelegations(periodInSeconds int) error {
	return stub.PollDelegationsFn(periodInSeconds)
}
