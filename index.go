package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"t_unicorn/authPswdManager"
	"t_unicorn/jwtHandler"
	"t_unicorn/meth"
	"t_unicorn/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "dongho"
	DB_PASSWORD = "qwer4321"
	DB_NAME     = "test_with_go"
	SALT_SIZE   = 16
)

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
	meth.PrintMessage("Regist Users..")
	err_1 := db.QueryRow(query_1).Scan(&lastInsertID)
	meth.CheckErr(err_1)
	query_2 := fmt.Sprintf(`
	INSERT INTO t_unicorn.user_auth_salt(
		user_id, username, salt
	) VALUES(
		'%d', '%s', '%s'
	)
	`, lastInsertID, userName, new_salt)
	err_2 := db.QueryRow(query_2).Scan(&lastInsertID)
	meth.CheckErr(err_2)
	response := models.JsonResponse{Type: "success", Message: "%s registed."}
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	userPswd := r.FormValue("userpswd")
	db := setupDB()
	meth.PrintMessage("Get User like login..")
	get_salt_query := fmt.Sprintf(`
		select salt from t_unicorn.user_auth_salt uas
		where uas.username = '%s'
	`, userName)
	rows, err := db.Query(get_salt_query)
	defer rows.Close()
	var salt_hex string
	for rows.Next() {
		err := rows.Scan(&salt_hex)
		meth.CheckErr(err)
	}
	userPswdHash := authPswdManager.HashPassword(userPswd, salt_hex)
	meth.CheckErr(err)
	query := fmt.Sprintf(`
		select user_id, username, email, created_on from t_unicorn.user_auth ua
		where ua.username = '%s' and 
		ua.password = '%s'
	`, userName, userPswdHash)
	rows, err = db.Query(query)
	meth.CheckErr(err)
	var users []models.User
	for rows.Next() {
		var user_id int
		var username string
		var email string
		var created_on string

		err = rows.Scan(&user_id, &username, &email, &created_on)
		meth.CheckErr(err)
		users = append(users, models.User{UserID: user_id, UserName: username, UserEmail: email, UserCreatedOn: created_on})
	}
	accessToken, err := jwtHandler.CreateJWT(users[0].UserEmail)
	meth.CheckErr(err)
	cookie := new(http.Cookie)
	cookie.Name = "access-token"
	cookie.Value = accessToken
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(time.Hour * 24)
	http.SetCookie(w, cookie)
	var response = models.JsonResponse{Type: "success__", Data: users}
	json.NewEncoder(w).Encode(response)
}

func JWTValidator(w http.ResponseWriter, r *http.Request) {
	T, err := r.Cookie("access-token")
	str_t := strings.Replace(T.String(), "access-token=", "", -1)
	meth.PrintMessage("TOKEN : " + str_t)

	token, err := jwt.Parse(str_t, func(token *jwt.Token) (interface{}, error) {
		return []byte("t_unicorn"), nil
	})
	meth.CheckErr(err)
	if _, ok := token.Claims.(jwt.Claims); ok && token.Valid {
		fmt.Printf("Token valid.")
	} else {
		fmt.Println(err)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	meth.PrintMessage("Getting Users..")
	rows, err := db.Query(`SELECT user_id, username, email, created_on FROM t_unicorn.user_auth;`)
	meth.CheckErr(err)
	fmt.Println(rows.Next())
	var users []models.User
	for rows.Next() {
		var user_id int
		var username string
		var email string
		var created_on string

		err = rows.Scan(&user_id, &username, &email, &created_on)
		meth.CheckErr(err)
		users = append(users, models.User{UserID: user_id, UserName: username, UserEmail: email, UserCreatedOn: created_on})
	}
	var response = models.JsonResponse{Type: "success", Data: users}
	json.NewEncoder(w).Encode(response)
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user = %s password = %s dbname = %s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	meth.CheckErr(err)
	return db
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/RegistUser", RegistUser).Methods("POST")
	router.HandleFunc("/GetUsers", GetUsers).Methods("GET")
	router.HandleFunc("/JWTValidator", JWTValidator).Methods("GET")
	router.HandleFunc("/GetUser", GetUser).Methods("POST")
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
