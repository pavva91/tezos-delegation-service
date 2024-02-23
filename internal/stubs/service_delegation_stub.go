package stubs

import (
	"time"

	"github.com/pavva91/tezos-delegation-service/internal/models"
)

type DelegationServiceStub struct {
	ListFn func(time.Time) ([]models.Delegation, error)
	PollFn func(uint, string, bool, chan<- error, <-chan struct{}) error
}

func (stub DelegationServiceStub) List(year time.Time) ([]models.Delegation, error) {
	return stub.ListFn(year)
}

func (stub DelegationServiceStub) Poll(periodInSeconds uint, apiEndpoint string, quitOnError bool, errorCh chan<- error, interruptCh <-chan struct{}) error {
	return stub.PollFn(periodInSeconds, apiEndpoint, quitOnError, errorCh, interruptCh)
}
