package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "user"
	DB_PASSWORD = "password"
	DB_NAME     = "dbname"
	DB_HOST     = "127.0.0.1"
	DB_PORT     = "5432"
	DB_SSLMODE  = "disable"
)

func GetDbConnection() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME, DB_SSLMODE)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
