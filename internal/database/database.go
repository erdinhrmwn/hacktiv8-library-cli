package database

import (
	"database/sql"

	"github.com/erdinhrmwn/hacktiv8-library-cli/config"
	_ "github.com/go-sql-driver/mysql"
)

func InitializeDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DatabaseDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}
