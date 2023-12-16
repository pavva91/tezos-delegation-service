package models

import (
	"time"

	"gorm.io/gorm"
)

type Delegation struct {
	gorm.Model `swaggerignore:"true"`
	Timestamp  time.Time
	Amount     string
	Delegator  string
	Block      string
}
