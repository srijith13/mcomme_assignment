package db

import (
	"database/sql"
	"fmt"
	"monk-commerce/app/config"

	_ "github.com/lib/pq"
)

var connection = make(map[string]*sql.DB)

func CreateDbConPool() *sql.DB {
	dbinfo := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s sslmode=disable`, config.DbHost, config.DbUser, config.DbPass, "mct")
	if connection["mct"] == nil {
		connection["mct"], _ = sql.Open("postgres", dbinfo) //creating the db connection
	}
	return connection["mct"]
}
