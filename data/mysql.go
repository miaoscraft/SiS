package data

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func openDB(addr, user, pswd string) (err error) {
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/whitelist", user, pswd, addr))
	return err
}
