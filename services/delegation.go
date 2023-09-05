package services

import (
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pavva91/gin-gorm-rest/models"
	"github.com/pavva91/gin-gorm-rest/repositories"
)

var (
	DelegationService DelegationServiceInterface = delegationServiceImpl{}
)

type DelegationServiceInterface interface {
	ListAllDelegations() ([]models.Delegation, error)
	PollDelegations() error
}

type delegationServiceImpl struct{}

func (service delegationServiceImpl) ListAllDelegations() ([]models.Delegation, error) {
	return repositories.DelegationRepository.ListDelegations()
}

func (service delegationServiceImpl) PollDelegations() error {

	time.Sleep(3 * time.Second)
	for i := 1; ; i++ {
		log.Info().Msg(strconv.Itoa(i))
		time.Sleep(1 * time.Second)
	}
}
