package mocks

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pavva91/gin-gorm-rest/config"
	"gorm.io/gorm"
)

type DbOrgMock struct {
	Mock   sqlmock.Sqlmock
	SqlDB  *sql.DB
	GormDB *gorm.DB
}

func (dbOrmMock DbOrgMock) ConnectToDB(cfg config.ServerConfig) {
}

func (dbOrmMock DbOrgMock) GetDB() *gorm.DB {
	return dbOrmMock.GormDB
}
