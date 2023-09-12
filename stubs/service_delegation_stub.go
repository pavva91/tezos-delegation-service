package stubs

import (
	"time"

	"github.com/pavva91/tezos-delegation-service/models"
)

type DelegationServiceStub struct {
	ListDelegationsFn func(time.Time) ([]models.Delegation, error)
	PollDelegationsFn func(uint, string, bool, chan<- error, <-chan struct{}) error
}

func (stub DelegationServiceStub) ListDelegations(year time.Time) ([]models.Delegation, error) {
	return stub.ListDelegationsFn(year)
}

func (stub DelegationServiceStub) PollDelegations(periodInSeconds uint, apiEndpoint string, quitOnError bool, errorCh chan<- error, interruptCh <-chan struct{}) error {
	return stub.PollDelegationsFn(periodInSeconds, apiEndpoint, quitOnError, errorCh, interruptCh)
}
