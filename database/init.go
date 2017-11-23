package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
)

var (
	Instance *sql.DB
	err error
)

func Init() error {
	Instance, err = sql.Open("mysql", "root:8343044@/api")								//Instance, err = sql.Open("mysql", "root:8343044@/api")
	if err != nil {
		return err
	}
	if err := Instance.Ping(); err != nil {
		return errors.New("Cannot establish any connections to the database.")
	}

	return nil
}