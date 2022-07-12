package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"t_unicorn/authPswdManager"
	dbManager "t_unicorn/dbManager"
	jwtManager "t_unicorn/jwtManager"
	. "t_unicorn/meth"
	"t_unicorn/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func RegistUser(w http.ResponseWriter, r *http.Request) {
	userPswd := r.FormValue("userpswd")
	userName := r.FormValue("username")
	db := dbManager.SetupDB()
	var lastInsertID int

	new_salt, saltedUserPswd := authPswdManager.GetSaltedUserPswd(userPswd)

	query := dbManager.GetRegistUserAuthQuery(r, saltedUserPswd)
	PrintMessage("Regist UserAuth")
	err := db.QueryRow(query).Scan(&lastInsertID)
	CheckErr(err)
	query = dbManager.GetRegistUserAuthSaltQuery(lastInsertID, userName, new_salt)
	PrintMessage("Regist UserAuthSalt")
	err = db.QueryRow(query).Scan(&lastInsertID)
	CheckErr(err)
	query = dbManager.GetRegistUserInfoQuery(r, lastInsertID)
	PrintMessage("Regist UserInfo")
	err = db.QueryRow(query).Scan(&lastInsertID)
	CheckErr(err)
	response := models.JsonResponse{Type: "success", Message: "%s registed."}
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	userPswd := r.FormValue("userpswd")
	db := dbManager.SetupDB()
	PrintMessage("Get User like login..")
	query := dbManager.GetUserSaltQuery(r)
	rows, err := db.Query(query)
	defer rows.Close()
	var salt_hex string
	for rows.Next() {
		err := rows.Scan(&salt_hex)
		CheckErr(err)
	}
	userPswdHash := authPswdManager.HashPassword(userPswd, salt_hex)
	CheckErr(err)
	query = dbManager.GetUserAuthQuery(userName, userPswdHash)
	rows, err = db.Query(query)
	CheckErr(err)
	var users []models.User
	for rows.Next() {
		var user_id int
		var username string
		var email string
		var created_on string

		err = rows.Scan(&user_id, &username, &email, &created_on)
		CheckErr(err)
		users = append(users, models.User{UserID: user_id, UserName: username, UserEmail: email, UserCreatedOn: created_on})
	}
	accessToken, err := jwtManager.CreateJWT(users[0].UserEmail)
	CheckErr(err)
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
	PrintMessage("TOKEN : " + str_t)
	token, err := jwt.Parse(str_t, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtManager.GetJWTSignature()), nil
	})
	CheckErr(err)
	c, _ := token.Claims.(jwt.MapClaims)
	fmt.Println(c["Email"])
	if _, ok := token.Claims.(jwt.Claims); ok && token.Valid {
		fmt.Printf("Token valid.")
	} else {
		fmt.Println(err)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := dbManager.SetupDB()
	PrintMessage("Getting Users..")
	query := dbManager.GetUsersQuery()
	rows, err := db.Query(query)
	CheckErr(err)
	fmt.Println(rows.Next())
	var users []models.User
	for rows.Next() {
		var user_id int
		var username string
		var email string
		var created_on string

		err = rows.Scan(&user_id, &username, &email, &created_on)
		CheckErr(err)
		users = append(users, models.User{UserID: user_id, UserName: username, UserEmail: email, UserCreatedOn: created_on})
	}
	var response = models.JsonResponse{Type: "success", Data: users}
	json.NewEncoder(w).Encode(response)
}

func MockRegistUsers(w http.ResponseWriter, r *http.Request) {
	db := dbManager.SetupDB()

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/RegistUser", RegistUser).Methods("POST")
	router.HandleFunc("/GetUsers", GetUsers).Methods("GET")
	router.HandleFunc("/JWTValidator", JWTValidator).Methods("GET")
	router.HandleFunc("/GetUser", GetUser).Methods("POST")
	router.HandleFunc("/MockRegistUsers", MockRegistUsers).Methods("GET")
	fmt.Println("Server at 33443")
	log.Fatal(http.ListenAndServe("localhost:33443", router))
}
