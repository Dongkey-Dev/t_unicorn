package meth

import (
	"database/sql"
	"fmt"
)

func PrintMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func CheckErr(e error) {
	if e != nil && e != sql.ErrNoRows {
		panic(e)
	}
}
