package jwtHandler

import(
  "time"
  "t_unicorn/meth"
  "github.com/golang-jwt/jwt/v4"
)

func CreateJWT(Email string) (string, error){
  signingKey := []byte("t_unicorn") //need to change local env file
  Tok := jwt.New(jwt.SigningMethodHS256)
  claims := Tok.Claims.(jwt.MapClaims)
  claims["Email"] = Email
  claims["exp"] = time.Now().Add(time.Minute * 20).Unix()

  tk, err := Tok.SignedString(signingKey)
  meth.CheckErr(err)
  return tk, nil
}
