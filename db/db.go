package db

import (
	"fmt"

	"github.com/pavva91/gin-gorm-rest/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DbOrm DbOrmInterface = dbOrmImpl{}
)

type DbOrmInterface interface {
	ConnectToDB(cfg config.ServerConfig)
	GetDB() *gorm.DB
}

type dbOrmImpl struct{}

var database *gorm.DB

func (dbOrm dbOrmImpl) ConnectToDB(cfg config.ServerConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	db.Logger.LogMode(logger.Info)
	// TODO: Fix circular dependency if I put AutoMigrate here
	// FIX: move to mani.go
	// db.AutoMigrate(&models.User{}, &models.Event{})
	database = db
}

// func Migrate() {
// 	database.AutoMigrate(&models.User{}, &models.Event{})
// }

func (dbOrm dbOrmImpl) GetDB() *gorm.DB {
	return database
}
