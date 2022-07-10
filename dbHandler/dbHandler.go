package dbHandler

import (
	"database/sql"
	"fmt"
	"os"
	"t_unicorn/meth"

	"github.com/joho/godotenv"
)

func SetupDB() *sql.DB {
	err := godotenv.Load(".env")
	meth.CheckErr(err)
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	dbinfo := fmt.Sprintf("user = %s password = %s dbname = %s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	meth.CheckErr(err)
	return db
}
