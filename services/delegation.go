package services

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pavva91/gin-gorm-rest/config"
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

	oldTime := time.Now().UTC()
	time.Sleep(3 * time.Second)
	response, err := http.Get(config.ServerConfigValues.ApiDelegations.Endpoint + "/operations/delegations")
	if err != nil {
		return err
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	// TODO: First I populate the db for the first time with the simple call: https://api.tzkt.io/v1/operations/delegations
	fmt.Println(string(responseData))

	// NOTE: This will be a sync task, the polling will not start before this task is done
	// NOTE: Use RWMutex
	// time.Sleep(3 * time.Second)
	for i := 1; ; i++ {
		newTime := time.Now().UTC()
		log.Info().Msg(newTime.Format(time.RFC3339))
		log.Info().Msg(strconv.Itoa(i))
		// TODO: Here I call only the date greater than previous call date (old timeNow) https://api.tzkt.io/v1/operations/delegations?timestamp.gt=2020-02-20T02:40:57Z
		response, err := http.Get(config.ServerConfigValues.ApiDelegations.Endpoint + "/operations/delegations?timestamp.gt="+oldTime.Format(time.RFC3339))
		if err != nil {
			return err
		}
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(responseData))
		oldTime = newTime
		time.Sleep(1 * time.Second)
	}
}
