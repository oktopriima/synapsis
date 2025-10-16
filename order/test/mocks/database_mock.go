package mocks

import (
	"synapsis/order/database/connection"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Instance() (connection.DBInstance, sqlmock.Sqlmock) {
	gormDb, mock := setupMock()
	ins := new(connection.Instance)
	ins.GormDB = gormDb

	return ins, mock
}

func setupMock() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, mock
	}

	// Use postgres driver instead of mysql
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
	})
	if err != nil {
		panic("failed to connect gorm DB: " + err.Error())
	}

	return gormDB, mock
}
