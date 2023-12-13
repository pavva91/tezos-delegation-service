package db

import (
	"fmt"

	"github.com/pavva91/tezos-delegation-service/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DbOrm    OrmInterface = dbOrmImpl{}
	database *gorm.DB
)

type OrmInterface interface {
	ConnectToDB(cfg config.ServerConfig)
	GetDB() *gorm.DB
}

type dbOrmImpl struct{}

func (dbOrm dbOrmImpl) ConnectToDB(cfg config.ServerConfig) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Error().Msg(err.Error())
		panic(fmt.Errorf("error connecting db: %w", err))
	}

	db.Logger.LogMode(logger.Info)
	database = db
}

func (dbOrm dbOrmImpl) GetDB() *gorm.DB {
	return database
}
