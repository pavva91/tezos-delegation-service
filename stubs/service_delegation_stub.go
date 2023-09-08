package stubs

import (
	"sync"
	"time"

	"github.com/pavva91/tezos-delegation-service/models"
)

type DelegationServiceStub struct {
	ListDelegationsFn func(time.Time) ([]models.Delegation, error)
	PollDelegationsFn func(int, string, *sync.RWMutex, bool, chan<- error, <-chan struct{}) error
}

func (stub DelegationServiceStub) ListDelegations(year time.Time) ([]models.Delegation, error) {
	return stub.ListDelegationsFn(year)
}

func (stub DelegationServiceStub) PollDelegations(periodInSeconds int, apiEndpoint string, rwmu *sync.RWMutex, quitOnError bool, errorCh chan<- error, interruptCh <-chan struct{}) error {
	return stub.PollDelegationsFn(periodInSeconds, apiEndpoint, rwmu, quitOnError, errorCh, interruptCh)
}
