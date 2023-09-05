package dto

import "time"

type DelegationResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Amount    string    `json:"amount"`
	Delegator string    `json:"delegator"`
	Block     string    `json:"block"`
}
