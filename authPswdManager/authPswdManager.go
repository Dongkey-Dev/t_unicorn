package authPswdManager

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"os"
	"strconv"
	"t_unicorn/meth"

	"github.com/joho/godotenv"
)

func generateRandomSaltHex(saltSize int) string {
	var salt = make([]byte, saltSize)
	_, err := rand.Read(salt[:])
	if err != nil {
		panic(err)
	}
	salt_hex := hex.EncodeToString(salt)
	return salt_hex
}

func getSaltSize() int {
	err := godotenv.Load()
	meth.CheckErr(err)
	SALT_SIZE, _ := strconv.Atoi(os.Getenv("SALT_SIZE"))
	return SALT_SIZE
}

func implementsBar(v interface{}) bool {
	type Barer interface {
		Bar() string
	}
	_, ok := v.(Barer)
	return ok
}

func HashPassword(password string, salt_hex string) string {
	salt, err := hex.DecodeString(salt_hex)
	if err != nil {
		panic(err)
	}
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()
	passwordBytes = append(passwordBytes, salt...)
	sha512Hasher.Write(passwordBytes)
	hashedPasswordBytes := sha512Hasher.Sum(nil)
	var base64EncodeedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)
	return base64EncodeedPasswordHash
}

func DoPasswordsMatch(hashedPassword, currPassword string, salt []byte) bool {
	salt_hex := hex.EncodeToString(salt)
	var currPasswordHash = HashPassword(currPassword, salt_hex)
	return hashedPassword == currPasswordHash
}

func GetSaltedUserPswd(userPswd string) (string, string) {
	SALT_SIZE := getSaltSize()
	new_salt := generateRandomSaltHex(SALT_SIZE)
	saltedUserPswd := HashPassword(userPswd, new_salt)
	return new_salt, saltedUserPswd
}
