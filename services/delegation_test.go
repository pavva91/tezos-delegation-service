package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/pavva91/tezos-delegation-service/dto"
	"github.com/pavva91/tezos-delegation-service/models"
	"github.com/pavva91/tezos-delegation-service/repositories"
	"github.com/pavva91/tezos-delegation-service/stubs"
	"github.com/stretchr/testify/assert"
)

func Test_ListDelegations_YearNonZero_Error(t *testing.T) {
	nonZeroValueDate := time.Now()
	errorMessage := "Unexpected Internal Error"
	unexpectedError := errors.New(errorMessage)
	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListByYearFn = func(time.Time) ([]models.Delegation, error) {
		return nil, unexpectedError
	}
	repositories.DelegationRepository = delegationRepositoryStub

	delegations, err := DelegationService.ListDelegations(nonZeroValueDate)

	assert.NotNil(t, err)
	assert.Nil(t, delegations)
	assert.Equal(t, errorMessage, err.Error())
}

func Test_ListDelegations_YearIsZero_Error(t *testing.T) {
	var zeroValueDate time.Time
	errorMessage := "Unexpected Internal Error"
	unexpectedError := errors.New(errorMessage)
	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListFn = func() ([]models.Delegation, error) {
		return nil, unexpectedError
	}
	repositories.DelegationRepository = delegationRepositoryStub

	delegations, err := DelegationService.ListDelegations(zeroValueDate)

	assert.NotNil(t, err)
	assert.Nil(t, delegations)
	assert.Equal(t, errorMessage, err.Error())
}

func Test_ListDelegations_YearNonZeroEmptyList_Empty(t *testing.T) {
	nonZeroValueDate := time.Now()
	emptyDelegationList := []models.Delegation{}

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListByYearFn = func(time.Time) ([]models.Delegation, error) {
		return emptyDelegationList, nil
	}
	repositories.DelegationRepository = delegationRepositoryStub

	delegations, err := DelegationService.ListDelegations(nonZeroValueDate)

	assert.Nil(t, err)
	assert.NotNil(t, delegations)
	assert.Equal(t, 0, len(delegations))
}

func Test_ListDelegations_YearIsZeroEmptyList_Empty(t *testing.T) {
	var zeroValueDate time.Time
	emptyDelegationList := []models.Delegation{}

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.ListFn = func() ([]models.Delegation, error) {
		return emptyDelegationList, nil
	}
	repositories.DelegationRepository = delegationRepositoryStub

	delegations, err := DelegationService.ListDelegations(zeroValueDate)

	assert.Nil(t, err)
	assert.NotNil(t, delegations)
	assert.Equal(t, 0, len(delegations))
}

func Test_PollDelegations_WrongApiEndpointScheme_Error(t *testing.T) {
	wrongApiEndpoint := "wrong"
	pollPeriodInSeconds := uint(1)
	expectedErrorContent1 := "Get \"" + wrongApiEndpoint + "/operations/delegations?timestamp.ge="
	expectedErrorContent2 := "unsupported protocol scheme"
	rwmutex := &sync.RWMutex{}
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})

	go DelegationService.PollDelegations(pollPeriodInSeconds, wrongApiEndpoint, rwmutex, stopOnError, errorCh, interruptCh)

	time.Sleep(5 * time.Second)
	interruptCh <- struct{}{}
	err := <-errorCh
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
	assert.Contains(t, err.Error(), expectedErrorContent2)
}

func Test_PollDelegations_WrongApiEndpointDomain_Error(t *testing.T) {
	wrongApiEndpoint := "http://wrong"
	pollPeriodInSeconds := uint(1)
	expectedErrorContent1 := "Get \"" + wrongApiEndpoint + "/operations/delegations?timestamp.ge="
	expectedErrorContent2 := "dial tcp: lookup wrong: no such host"
	rwmutex := &sync.RWMutex{}
	// TODO: stopOnError becomes a boolean
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})
	// TODO: create channel token to send signal

	go DelegationService.PollDelegations(pollPeriodInSeconds, wrongApiEndpoint, rwmutex, stopOnError, errorCh, interruptCh)

	time.Sleep(5 * time.Second)
	// TODO: create channel token to send signal
	interruptCh <- struct{}{}
	err := <-errorCh
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
	assert.Contains(t, err.Error(), expectedErrorContent2)
}

func Test_PollDelegations_Not200FromApiEndpoint_Error(t *testing.T) {
	pollPeriodInSeconds := uint(1)
	errorHttpStatus := http.StatusBadRequest
	expectedError := "Get Response different than 200: " + strconv.Itoa(errorHttpStatus)
	rwmutex := &sync.RWMutex{}
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})

	// Mock outbound http request https://medium.com/zus-health/mocking-outbound-http-requests-in-go-youre-probably-doing-it-wrong-60373a38d2aa
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/operations/delegations" {
			t.Errorf("Expected to request '/operations/delegations', got: %s", r.URL.Path)
		}
		w.WriteHeader(errorHttpStatus)
		w.Write([]byte{})
	}))
	defer server.Close()

	go DelegationService.PollDelegations(pollPeriodInSeconds, server.URL, rwmutex, stopOnError, errorCh, interruptCh)

	time.Sleep(3 * time.Second)
	interruptCh <- struct{}{}
	err := <-errorCh
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func Test_PollDelegations_ReturnedUnexpectedJSON_Error(t *testing.T) {
	pollPeriodInSeconds := uint(1)
	unexpectedJSON := []byte(`{"value":"fixed"}`)
	expectedError := "json: cannot unmarshal object into Go value of type []dto.DelegationResponseFromApi"
	rwmutex := &sync.RWMutex{}
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/operations/delegations" {
			t.Errorf("Expected to request '/operations/delegations', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(unexpectedJSON)
	}))
	defer server.Close()

	err := DelegationService.PollDelegations(pollPeriodInSeconds, server.URL, rwmutex, stopOnError, errorCh, interruptCh)
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func Test_PollDelegations_WorksThenApiGoDownAfter2Seconds_Error(t *testing.T) {
	pollPeriodInSeconds := uint(1)
	emptyListDelegations := []dto.DelegationResponseFromApi{}
	jsonResponse, err := json.Marshal(emptyListDelegations)
	if err != nil {
		return
	}
	// NOTE: error: "Get \"http://127.0.0.1:45969/operations/delegations?timestamp.ge=2023-09-07T08:49:59Z&timestamp.lt=2023-09-07T08:50:01Z\": dial tcp 127.0.0.1:45969: connect: connection refused"
	expectedErrorContent1 := "connect: connection refused"
	rwmutex := &sync.RWMutex{}
	stopOnError := false
	errorCh := make(chan error)
	interruptCh := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/operations/delegations" {
			t.Errorf("Expected to request '/operations/delegations', got: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))

	go func() {
		time.Sleep(2 * time.Second)
		server.Close()
	}()

	defer server.Close()

	go DelegationService.PollDelegations(pollPeriodInSeconds, server.URL, rwmutex, stopOnError, errorCh, interruptCh)

	time.Sleep(5 * time.Second)
	interruptCh <- struct{}{}
	err = <-errorCh
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), expectedErrorContent1)
}

func Test_SaveBulkDelegations_EmptySlice_NoError(t *testing.T) {
	var emptyList []dto.DelegationResponseFromApi
	rwmutex := &sync.RWMutex{}

	savedDelegations, err := SaveBulkDelegations(emptyList, rwmutex)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(savedDelegations))
}

func Test_SaveBulkDelegations_ErrorFromRepositoryCreate_Error(t *testing.T) {
	var emptyList []dto.DelegationResponseFromApi
	delegation1 := dto.DelegationResponseFromApi{}
	listOneElement := append(emptyList, delegation1)
	rwmutex := &sync.RWMutex{}

	errorMessage := "Unexpected Internal Error"
	unexpectedError := errors.New(errorMessage)

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) (*models.Delegation, error) {
		return nil, unexpectedError
	}
	repositories.DelegationRepository = delegationRepositoryStub

	savedDelegations, err := SaveBulkDelegations(listOneElement, rwmutex)
	fmt.Println(err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, errorMessage, err.Error())
	assert.Equal(t, 0, len(savedDelegations))
}

func Test_SaveBulkDelegations_OKList1Element_ReturnSavedDelegation(t *testing.T) {
	var emptyList []dto.DelegationResponseFromApi
	delegation := dto.DelegationResponseFromApi{}
	delegationModel1 := models.Delegation{}
	listOneElement := append(emptyList, delegation)
	rwmutex := &sync.RWMutex{}

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) (*models.Delegation, error) {
		return &delegationModel1, nil
	}
	repositories.DelegationRepository = delegationRepositoryStub

	savedDelegations, err := SaveBulkDelegations(listOneElement, rwmutex)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(savedDelegations))
}

func Test_SaveBulkDelegations_OKList2Element_ReturnSavedDelegation(t *testing.T) {
	var delegations []dto.DelegationResponseFromApi
	delegation := dto.DelegationResponseFromApi{}
	delegationModel1 := models.Delegation{}
	delegations = append(delegations, delegation)
	delegations = append(delegations, delegation)
	rwmutex := &sync.RWMutex{}

	delegationRepositoryStub := stubs.DelegationRepositoryStub{}
	delegationRepositoryStub.CreateFn = func(*models.Delegation) (*models.Delegation, error) {
		return &delegationModel1, nil
	}
	repositories.DelegationRepository = delegationRepositoryStub

	savedDelegations, err := SaveBulkDelegations(delegations, rwmutex)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(savedDelegations))
}
