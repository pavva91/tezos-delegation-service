package mocks

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pavva91/tezos-delegation-service/config"
	"gorm.io/gorm"
)

type DBOrg struct {
	Mock   sqlmock.Sqlmock
	DB  *sql.DB
	GormDB *gorm.DB
}

func (dbm DBOrg) ConnectToDB(cfg config.ServerConfig) {
}

func (dbm DBOrg) GetDB() *gorm.DB {
	return dbm.GormDB
}
