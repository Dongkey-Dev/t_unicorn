package main

import (
	"t_unicorn/authPswdManager"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "dongho"
	DB_PASSWORD = "qwer4321"
	DB_NAME     = "test_with_go"
	SALT_SIZE   = 16
)

type User struct {
	UserID        int    `json:"userid"`
	UserName      string `json:"username"`
	UserEmail     string `json:"useremail"`
	UserCreatedOn string `json:"usercreatedon"`
}

type UserSequenceID struct {
	nextval int `json:"userid"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []User `json:"data"`
	Message string `json:"message"`
}

type JsonResponseSequenceID struct {
	Type    string           `json:"type"`
	Data    []UserSequenceID `json:"data"`
	Message string           `json:"message"`
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func checkErr(err error) {
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
}

func RegistUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	userPswd := r.FormValue("userpswd")
	userEmail := r.FormValue("email")
	db := setupDB()
	new_salt := authPswdManager.GenerateRandomSaltHex(SALT_SIZE)
	var lastInsertID int
	saltedUserPswd := authPswdManager.HashPassword(userPswd, new_salt)
	query_1 := fmt.Sprintf(`
	INSERT INTO t_unicorn.user_auth(
		username, password, email
	) VALUES(
		'%s', '%s', '%s'
	) returning user_id
	`, userName, saltedUserPswd, userEmail)
	printMessage("Regist Users..")
	err_1 := db.QueryRow(query_1).Scan(&lastInsertID)
	checkErr(err_1)
	query_2 := fmt.Sprintf(`
	INSERT INTO t_unicorn.user_auth_salt(
		user_id, username, salt
	) VALUES(
		'%d', '%s', '%s'
	)
	`, lastInsertID, userName, new_salt)
	err_2 := db.QueryRow(query_2).Scan(&lastInsertID)
	checkErr(err_2)
	response := JsonResponse{Type: "success", Message: "%s registed."}
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	userPswd := r.FormValue("userpswd")

	db := setupDB()
	printMessage("Get User like login..")
	get_salt_query := fmt.Sprintf(`
		select salt from t_unicorn.user_auth_salt uas
		where uas.username = '%s'
	`, userName)
	rows, err := db.Query(get_salt_query)
	defer rows.Close()
	var salt_hex string
	for rows.Next() {
		err := rows.Scan(&salt_hex)
		checkErr(err)
	}
	userPswdHash := authPswdManager.HashPassword(userPswd, salt_hex)
	checkErr(err)
	query := fmt.Sprintf(`
		select user_id, username, email, created_on from t_unicorn.user_auth ua
		where ua.username = '%s' and 
		ua.password = '%s'
	`, userName, userPswdHash)
	rows, err = db.Query(query)
	checkErr(err)
	var users []User
	for rows.Next() {
		var user_id int
		var username string
		var email string
		var created_on string

		err = rows.Scan(&user_id, &username, &email, &created_on)
		checkErr(err)
		users = append(users, User{UserID: user_id, UserName: username, UserEmail: email, UserCreatedOn: created_on})
	}
	var response = JsonResponse{Type: "success__", Data: users}
	json.NewEncoder(w).Encode(response)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting Users..")
	rows, err := db.Query(`SELECT user_id, username, email, created_on FROM t_unicorn.user_auth;`)
	checkErr(err)
	fmt.Println(rows.Next())
	var users []User
	for rows.Next() {
		var user_id int
		var username string
		var email string
		var created_on string

		err = rows.Scan(&user_id, &username, &email, &created_on)

		// check errors
		checkErr(err)
		users = append(users, User{UserID: user_id, UserName: username, UserEmail: email, UserCreatedOn: created_on})
	}
	var response = JsonResponse{Type: "success", Data: users}
	json.NewEncoder(w).Encode(response)
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user = %s password = %s dbname = %s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	return db
}

func main() {
	router := mux.NewRouter()
	// router.HandleFunc("/Login/", UserLogin).Methods("GET")
	router.HandleFunc("/RegistUser", RegistUser).Methods("POST")
	router.HandleFunc("/GetUsers", GetUsers).Methods("GET")
	router.HandleFunc("/GetUser", GetUser).Methods("POST")
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

