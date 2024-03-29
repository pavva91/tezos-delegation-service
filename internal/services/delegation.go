package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pavva91/tezos-delegation-service/config"
	"github.com/pavva91/tezos-delegation-service/internal/dto"
	"github.com/pavva91/tezos-delegation-service/internal/models"
	"github.com/pavva91/tezos-delegation-service/internal/repositories"
	"github.com/rs/zerolog/log"
)

var (
	Delegation DelegationServicer = delegation{}
)

type DelegationServicer interface {
	List(year time.Time) ([]models.Delegation, error)
	// PollDelegations Fuction that runs asynchronously for polling delegations
	Poll(periodInSeconds uint, apiEndpoint string, quitOnError bool, errorOutCh chan<- error, quitOnErrorTrueSignalInCh <-chan struct{}) error
}

type delegation struct{}

func (s delegation) List(year time.Time) ([]models.Delegation, error) {
	if year.IsZero() {
		return repositories.Delegation.List()
	}
	return repositories.Delegation.ListByYear(year)
}

func SaveBulkDelegations(delegations []dto.DelegationResponseFromAPI) ([]models.Delegation, error) {
	var savedDelegations []models.Delegation

	for _, d := range delegations {
		// TODO: Check if is a replicate before adding to DB
		delegationModel := d.ToModel()
		// NOTE: I Use gorm that is Thread-Safe, so a RWMutex is not needed on my side,
		// I just add it for showing what I would have done if I had to handle myself race conditions
		// rwmu.Lock()
		// NOTE: I could use defer rwmu.Unlock()
		// In this case I prefer to make the 2 explicit calls
		err := repositories.Delegation.Create(delegationModel)
		if err != nil {
			log.Info().Err(err).Msg("error creating delegation in db")
			// rwmu.Unlock()
			return nil, fmt.Errorf("repository error: %w", err)
		}
		// rwmu.Unlock()
		savedDelegations = append(savedDelegations, *delegationModel)
		log.Info().Msg("Delegation Created Correctly: " + strconv.Itoa(int(delegationModel.ID)))
	}
	return savedDelegations, nil
}

func (s delegation) Poll(periodInSeconds uint, apiEndpoint string, quitOnError bool, errorOutCh chan<- error, quitOnErrorTrueSignalInCh <-chan struct{}) error {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// NOTE: Using Now() can be a problem if timestamp are not in sync within servers. To be on the safe side, if there's no strict performance boundary is to put now a couple of minutes before:
	// oldTime := time.Now().UTC()
	oldTime := time.Now().UTC().Add(-time.Second * time.Duration(config.ServerConfigValues.APIDelegations.DelayLocalTimestampInSeconds))
	log.Info().Msg("Starting time polling: " + oldTime.Format(time.RFC3339Nano))

	time.Sleep(time.Duration(periodInSeconds) * time.Second)

	for {
		select {
		case <-quitOnErrorTrueSignalInCh:
			quitOnError = true
		default:
			// log.Info().Msg("Continue polling")
		}
		// NOTE: Using Now() can be a problem if timestamp are not in sync within servers. To be on the safe side, if there's no strict performance boundary is to put now a couple of minutes before
		// newTime := time.Now().UTC()
		newTime := time.Now().UTC().Add(-time.Second * time.Duration(config.ServerConfigValues.APIDelegations.DelayLocalTimestampInSeconds))

		// NOTE: Here I call only the date greater than previous call date (old timeNow) https://api.tzkt.io/v1/operations/delegations?timestamp.gt=2020-02-20T02:40:57Z
		response, err := client.Get(apiEndpoint + "/operations/delegations?timestamp.ge=" + oldTime.Format(time.RFC3339) + "&timestamp.lt=" + newTime.Format(time.RFC3339))
		if err != nil {
			log.Error().Err(err).Msg("connectivity error - no response from request")

			if quitOnError {
				errorOutCh <- err

				return err
			}
			time.Sleep(time.Duration(periodInSeconds) * time.Second)

			continue
		}

		if response.StatusCode != http.StatusOK {
			err := fmt.Errorf("get response different than 200: %w ", err)
			log.Error().Err(err).Msg("")

			if quitOnError {
				errorOutCh <- err

				return err
			}

			time.Sleep(time.Duration(periodInSeconds) * time.Second)

			continue
		}

		defer response.Body.Close()
		responseBody, err := io.ReadAll(response.Body)

		if err != nil {
			log.Error().Err(err).Msg("error reading response body")

			return err
		}

		var results []dto.DelegationResponseFromAPI
		err = json.Unmarshal(responseBody, &results)
		if err != nil {
			log.Error().Err(err).Msg("cannot unmarshal json")
			return err
		}

		savedDelegations, err := SaveBulkDelegations(results)
		if err != nil {
			return err
		}
		log.Info().Msg(fmt.Sprintf("saved delegations: %d", len(savedDelegations)))

		oldTime = newTime
		time.Sleep(time.Duration(periodInSeconds) * time.Second)
	}
}
