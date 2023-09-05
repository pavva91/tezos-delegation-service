package repositories

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pavva91/gin-gorm-rest/db"
	"github.com/pavva91/gin-gorm-rest/mocks"
	"github.com/pavva91/gin-gorm-rest/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func Test_CreateUser_OK(t *testing.T) {
	// Mocks
	username := "alice"
	user := models.User{
		Username: username,
	}
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer sqlDb.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm stub connection", err)
	}
	dbMock := mocks.DbOrgMock{
		Mock:   mock,
		SqlDB:  sqlDb,
		GormDB: gormDB,
	}

	db.DbOrm = dbMock

	expectedQuery := `INSERT INTO "users" ("created_at","updated_at","deleted_at","name","username","email","password") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`

	// Stubs
	rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
	dbMock.Mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(AnyTime{}, AnyTime{}, nil, "", username, "", "").WillReturnRows(rows)

	// Call function to test
	userReturn, err := UserRepository.CreateUser(&user)

	// Check Values
	err = dbMock.Mock.ExpectationsWereMet()
	assert.Nil(t, err)
	assert.Equal(t, userReturn, &user)
}

func Test_CreateUser_Error(t *testing.T) {
	// Mocks
	username := "alice"
	user := models.User{
		Username: username,
	}
	sqlDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub databse connection", err)
	}
	defer sqlDb.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm stub connection", err)
	}
	dbMock := mocks.DbOrgMock{
		Mock:   mock,
		SqlDB:  sqlDb,
		GormDB: gormDB,
	}

	db.DbOrm = dbMock

	expectedQuery := `INSERT INTO "users" ("created_at","updated_at","deleted_at","name","username","email","password") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`

	// Stubs
	errorMessage := "unexpected error"
	unexpectedError := errors.New(errorMessage)
	dbMock.Mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(AnyTime{}, AnyTime{}, nil, "", username, "", "").WillReturnError(unexpectedError)

	// Call function to test
	userReturn, err := UserRepository.CreateUser(&user)

	// Check Values
	err = dbMock.Mock.ExpectationsWereMet()
	assert.Nil(t, err)
	assert.Nil(t, userReturn)
}
