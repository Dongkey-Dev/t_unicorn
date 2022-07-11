package dbHandler

import (
	"database/sql"
	"fmt"
	"net/http"
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
	dbinfo := fmt.Sprintf(
		"user = %s password = %s dbname = %s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	meth.CheckErr(err)
	return db
}

func GetRegistUserAuthQuery(r *http.Request, saltedUserPswd string) string {
	rfv := r.FormValue
	query := fmt.Sprintf(`
	INSERT INTO t_unicorn.user_auth(
		username, password, email
	) VALUES(
		'%s', '%s', '%s'
	) returning user_id
	`, rfv("username"), saltedUserPswd, rfv("email"))
	return query
}

func GetRegistUserAuthSaltQuery(lastInsertID int, userName string, new_salt string) string {
	query := fmt.Sprintf(`
	INSERT INTO t_unicorn.user_auth_salt(
		user_id, username, salt
	) VALUES(
		'%d', '%s', '%s'
	)
	`, lastInsertID, userName, new_salt)
	return query
}

func GetRegistUserInfoQuery(r *http.Request, lastInsertID int) string {
	rfv := r.FormValue
	query := fmt.Sprintf(`
	INSERT INTO t_unicorn.user_info(
		user_id, name, dob, gender, phone
	) VALUES(
		'%d', '%s', '%s', '%s', '%s'
	)
	`, lastInsertID, rfv("name"), rfv("dob"), rfv("gender"), rfv("phone"))
	return query
}

func GetUserSaltQuery(r *http.Request) string {
	rfv := r.FormValue
	query := fmt.Sprintf(`
		select salt from t_unicorn.user_auth_salt uas
		where uas.username = '%s'
	`, rfv("username"))
	return query
}

func GetUserAuthQuery(userName string, userPswdHash string) string {
	query := fmt.Sprintf(`
		select user_id, username, email, created_on from t_unicorn.user_auth ua
		where ua.username = '%s' and 
		ua.password = '%s'
	`, userName, userPswdHash)
	return query
}

func GetUsersQuery() string {
	query := fmt.Sprintf(`SELECT user_id, username, email, created_on FROM t_unicorn.user_auth;`)
	return query
}
