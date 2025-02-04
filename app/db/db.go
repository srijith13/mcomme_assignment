package db

import (
	"database/sql"
	"fmt"
	"monk-commerce/app/config"

	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

func InitDb() (*sql.DB, error) {
	dbinfo := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s sslmode=disable`, config.DbHost, config.DbUser, config.DbPass, "mct")
	db, err := sql.Open("postgres", dbinfo) //creating the db connection
	if err != nil {
		log.Errorf("DB Connection Failed %v", err)
		return nil, err
	}
	return db, nil
}
