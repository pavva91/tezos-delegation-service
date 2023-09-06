package dto

import (
	"strconv"
	"time"

	"github.com/pavva91/tezos-delegation-service/models"
)

// NOTE: Created with https://mholt.github.io/json-to-go/
type DelegationResponseFromApi struct {
	Type      string    `json:"type"`
	ID        int64     `json:"id"`
	Level     int       `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Counter   int       `json:"counter"`
	Sender    struct {
		Address string `json:"address"`
	} `json:"sender"`
	GasLimit     int `json:"gasLimit"`
	GasUsed      int `json:"gasUsed"`
	StorageLimit int `json:"storageLimit"`
	BakerFee     int `json:"bakerFee"`
	Amount       int `json:"amount"`
	PrevDelegate struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"prevDelegate,omitempty"`
	Status      string `json:"status"`
	NewDelegate struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"newDelegate,omitempty"`
	Errors []struct {
		Type string `json:"type"`
	} `json:"errors,omitempty"`
}

func (dto *DelegationResponseFromApi) ToModel() *models.Delegation {
	var delegation models.Delegation
	delegation.Timestamp = dto.Timestamp
	delegation.Amount = strconv.Itoa(dto.Amount)
	delegation.Delegator = dto.Sender.Address
	delegation.Block = dto.Block
	delegation.ID = uint(dto.ID)
	return &delegation
}
