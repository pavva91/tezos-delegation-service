package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pavva91/tezos-delegation-service/config"
	"github.com/pavva91/tezos-delegation-service/dto"
	"github.com/pavva91/tezos-delegation-service/models"
	"github.com/pavva91/tezos-delegation-service/repositories"
	"github.com/rs/zerolog/log"
)

var (
	DelegationService DelegationServiceInterface = delegationServiceImpl{}
)

type DelegationServiceInterface interface {
	ListAllDelegations() ([]models.Delegation, error)
	PollDelegations(int) error
}

type delegationServiceImpl struct{}

func (service delegationServiceImpl) ListAllDelegations() ([]models.Delegation, error) {
	return repositories.DelegationRepository.List()
}

func (service delegationServiceImpl) PollDelegations(freqInSeconds int) error {

	oldTime := time.Now().UTC()
	time.Sleep(2 * time.Second)

	// NOTE: This will be a sync task, the polling will not start before this task is done
	// NOTE: Use RWMutex
	// time.Sleep(3 * time.Second)
	for i := 1; ; i++ {
		newTime := time.Now().UTC()
		log.Info().Msg(newTime.Format(time.RFC3339))
		log.Info().Msg(strconv.Itoa(i))
		// TODO: Here I call only the date greater than previous call date (old timeNow) https://api.tzkt.io/v1/operations/delegations?timestamp.gt=2020-02-20T02:40:57Z
		response, err := http.Get(config.ServerConfigValues.ApiDelegations.Endpoint + "/operations/delegations?timestamp.gt=" + oldTime.Format(time.RFC3339) + "&timestamp.lt=" + newTime.Format(time.RFC3339))
		if err != nil {
			log.Error().Err(err).Msg("No response from request")
			return err
		}
		// TODO: Parse raw response
		defer response.Body.Close()
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Info().Err(err).Msg("Error reading response body")
			return err
		}

		var result []dto.DelegationResponseFromApi
		err = json.Unmarshal(responseBody, &result)
		if err != nil {
			log.Info().Err(err).Msg("Cannot unmarshal JSON")
			return err
		}

		fmt.Println(len(result))
		if len(result) != 0 {
			// TODO: Add records to DB
			// TODO: From DTO to Model
			// TODO: Check if is a replicate before adding to DB
			fmt.Println(result[0].Hash)
		}

		fmt.Println(string(responseBody))
		oldTime = newTime
		time.Sleep(1 * time.Second)
		time.Sleep(time.Duration(freqInSeconds) * time.Second)
	}
}
