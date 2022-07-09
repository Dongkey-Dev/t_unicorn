package meth

import(
  "database/sql"
)

func CheckErr(e error){
  if e!=nil && e!=sql.ErrNoRows{
    panic(e)
  }
}
