package jwtHandler

import (
	"os"
	"t_unicorn/meth"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func GetJWTClaims() jwt.Claims {
	Tok := jwt.New(jwt.SigningMethodHS256)
	claims := Tok.Claims.(jwt.MapClaims)
	return claims
}

func CreateJWT(Email string) (string, error) {
	err := godotenv.Load()
	meth.CheckErr(err)
	signingKey := []byte(os.Getenv("CLAIMS_WORD"))
	Tok := jwt.New(jwt.SigningMethodHS256)
	claims := Tok.Claims.(jwt.MapClaims)
	claims["Email"] = Email
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()

	tk, err := Tok.SignedString(signingKey)
	meth.CheckErr(err)
	return tk, nil
}

func GetJWTSignature() string {
	err := godotenv.Load()
	meth.CheckErr(err)
	return os.Getenv("CLAIMS_WORD")
}
