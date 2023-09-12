package models

import (
	"time"

	"gorm.io/gorm"
)

type Delegation struct {
	gorm.Model `swaggerignore:"true"`
	Timestamp  time.Time `json:"timestamp"`
	Amount     string    `json:"amount" binding:"required"`
	Delegator  string    `json:"delegator" binding:"required"`
	Block      string    `json:"block" binding:"required"`
}
