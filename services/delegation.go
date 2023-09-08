package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pavva91/tezos-delegation-service/dto"
	"github.com/pavva91/tezos-delegation-service/models"
	"github.com/pavva91/tezos-delegation-service/repositories"
	"github.com/rs/zerolog/log"
)

var (
	DelegationService DelegationServiceInterface = delegationServiceImpl{}
)

type DelegationServiceInterface interface {
	ListDelegations(year time.Time) ([]models.Delegation, error)
	PollDelegations(int, string, *sync.RWMutex, bool, chan<- error, <-chan struct{}) error
}

type delegationServiceImpl struct{}

func (service delegationServiceImpl) ListDelegations(year time.Time) ([]models.Delegation, error) {
	if year.IsZero() {
		return repositories.DelegationRepository.List()
	} else {
		return repositories.DelegationRepository.ListByYear(year)
	}
}

func SaveBulkDelegations(delegations []dto.DelegationResponseFromApi, rwmu *sync.RWMutex) ([]models.Delegation, error) {
	var savedDelegations []models.Delegation
	for _, r := range delegations {
		delegationModel := r.ToModel()
		// TODO: Check if is a replicate before adding to DB
		rwmu.Lock()
		createdDelegation, err := repositories.DelegationRepository.Create(delegationModel)
		if err != nil {
			log.Info().Err(err).Msg("Error Creating Delegation in DB")
			return nil, err
		}
		rwmu.Unlock()
		savedDelegations = append(savedDelegations, *delegationModel)
		log.Info().Msg("Delegation Created Correctly: " + strconv.Itoa(int(createdDelegation.ID)))
	}
	return savedDelegations, nil
}

func (service delegationServiceImpl) PollDelegations(periodInSeconds int, apiEndpoint string, rwmu *sync.RWMutex, quitOnError bool, errorCh chan<- error, interruptCh <-chan struct{}) error {

	oldTime := time.Now().UTC()

	for {
		select {
		case <-interruptCh:
			quitOnError = true
		default:
			// log.Info().Msg("Continue polling")
		}
		newTime := time.Now().UTC()

		client := http.Client{
			Timeout: 5 * time.Second,
		}
		// NOTE: Here I call only the date greater than previous call date (old timeNow) https://api.tzkt.io/v1/operations/delegations?timestamp.gt=2020-02-20T02:40:57Z
		response, err := client.Get(apiEndpoint + "/operations/delegations?timestamp.ge=" + oldTime.Format(time.RFC3339) + "&timestamp.lt=" + newTime.Format(time.RFC3339))
		if err != nil {
			log.Info().Err(err).Msg("Connectivity Error - No response from request")
			if quitOnError {
				errorCh <- err
				return err
			} else {
				time.Sleep(time.Duration(periodInSeconds) * time.Second)
				continue
			}
		}
		if response.StatusCode != http.StatusOK {
			err := errors.New("Get Response different than 200: " + strconv.Itoa(response.StatusCode))
			log.Info().Err(err).Msg("")
			if quitOnError {
				errorCh <- err
				return err
			} else {
				time.Sleep(time.Duration(periodInSeconds) * time.Second)
				continue
			}
		}

		defer response.Body.Close()
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Info().Err(err).Msg("Error reading response body")
			return err
		}

		var results []dto.DelegationResponseFromApi
		err = json.Unmarshal(responseBody, &results)
		if err != nil {
			log.Info().Err(err).Msg("Cannot unmarshal JSON")
			return err
		}

		savedDelegations, err := SaveBulkDelegations(results, rwmu)
		if err != nil {
			return err
		}
		log.Info().Msg(fmt.Sprintf("Saved Delegations: %d", len(savedDelegations)))

		oldTime = newTime
		time.Sleep(time.Duration(periodInSeconds) * time.Second)
	}
}
