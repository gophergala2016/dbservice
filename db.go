package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func GetDbConnection() (*sql.DB, error) {
	config, err := ParseDbConfig(".")
	if err != nil {
		return nil, err
	}
	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Database, config.SslMode)
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
