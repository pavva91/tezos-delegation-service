package models

import (
	"time"

	"gorm.io/gorm"
)

type Delegation struct {
	gorm.Model `swaggerignore:"true"`
	// TODO: Timestamp must be ISO-8601 
	// TODO: time.RFC3339
	Timestamp  time.Time `json:"timestamp"`
	Amount     string    `json:"amount" binding:"required"`
	Delegator  string    `json:"delegator" binding:"required"`
	Block      int       `json:"block" binding:"required"`
}