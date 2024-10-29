package db

import (
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
)

var dbCliMock *sql.DB
var sqlMock sqlmock.Sqlmock

func InitDbMock() error {
	var err error
	dbCliMock, sqlMock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return nil
}

func GetDbClientMock() *sql.DB {
	return dbCliMock
}
func GetSqlMock() sqlmock.Sqlmock {
	return sqlMock
}
