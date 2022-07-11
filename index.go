package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"t_unicorn/authPswdManager"
	"t_unicorn/dbHandler"
	"t_unicorn/jwtHandler"
	"t_unicorn/meth"
	"t_unicorn/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func RegistUser(w http.ResponseWriter, r *http.Request) {
	userPswd := r.FormValue("userpswd")
	userName := r.FormValue("username")
	db := dbHandler.SetupDB()

	SALT_SIZE := authPswdManager.GetSaltSize()
	new_salt := authPswdManager.GenerateRandomSaltHex(SALT_SIZE)
	var lastInsertID int
	saltedUserPswd := authPswdManager.HashPassword(userPswd, new_salt)
	query := dbHandler.GetRegistUserAuthQuery(r, saltedUserPswd)
	meth.PrintMessage("Regist UserAuth")
	err := db.QueryRow(query).Scan(&lastInsertID)
	meth.CheckErr(err)
	query = dbHandler.GetRegistUserAuthSaltQuery(lastInsertID, userName, new_salt)
	meth.PrintMessage("Regist UserAuthSalt")
	err = db.QueryRow(query).Scan(&lastInsertID)
	meth.CheckErr(err)
	query = dbHandler.GetRegistUserInfoQuery(r, lastInsertID)
	meth.PrintMessage("Regist UserInfo")
	err = db.QueryRow(query).Scan(&lastInsertID)
	meth.CheckErr(err)
	response := models.JsonResponse{Type: "success", Message: "%s registed."}
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("username")
	userPswd := r.FormValue("userpswd")
	db := dbHandler.SetupDB()
	meth.PrintMessage("Get User like login..")
	query := dbHandler.GetUserSaltQuery(r)
	rows, err := db.Query(query)
	defer rows.Close()
	var salt_hex string
	for rows.Next() {
		err := rows.Scan(&salt_hex)
		meth.CheckErr(err)
	}
	userPswdHash := authPswdManager.HashPassword(userPswd, salt_hex)
	meth.CheckErr(err)
	query = dbHandler.GetUserAuthQuery(userName, userPswdHash)
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
		return []byte(jwtHandler.GetJWTSignature()), nil
	})
	meth.CheckErr(err)
	c, _ := token.Claims.(jwt.MapClaims)
	fmt.Println(c["Email"])
	if _, ok := token.Claims.(jwt.Claims); ok && token.Valid {
		fmt.Printf("Token valid.")
	} else {
		fmt.Println(err)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := dbHandler.SetupDB()
	meth.PrintMessage("Getting Users..")
	query := dbHandler.GetUsersQuery()
	rows, err := db.Query(query)
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/RegistUser", RegistUser).Methods("POST")
	router.HandleFunc("/GetUsers", GetUsers).Methods("GET")
	router.HandleFunc("/JWTValidator", JWTValidator).Methods("GET")
	router.HandleFunc("/GetUser", GetUser).Methods("POST")
	fmt.Println("Server at 33443")
	log.Fatal(http.ListenAndServe("localhost:33443", router))
}
