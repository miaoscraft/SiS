package data

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var seizeUUID *sql.Stmt

func openDB(addr, user, pswd string) (err error) {
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/whitelist", user, pswd, addr))
	if err != nil {
		return err
	}

	if _, err = db.Exec(`
drop procedure if exists SeizeUUID;
create procedure SeizeUUID(in MyQQ bigint, in MyName text, in MyUUID binary(16))
begin
    insert ignore into whitelist_test
        (QQ, Name, UUID)
    values (MyQQ, MyName, MyUUID)
    on duplicate key update Name=MyName, UUID=MyUUID;

    select QQ
    from whitelist_test
    where UUID = MyUUID;
end;
`); err != nil {
		return err
	}
	return err
}

// SetWhitelist 尝试向数据库写入白名单数据，当ID未被占用时返回自己的name，当ID被占用则返回占用者的name
// 若原本该账号占有一个UUID，则会返回当时的UUID
func SetWhitelist(QQ int64, name string, ID uuid.UUID) (int64, string, error) {
	rows, err := db.Query("call SeizeUUID(?,?,?)", QQ, name, ID)
	if err != nil {
		return 0, "", err
	}

	var oldName string
	err = rows.Scan(&QQ, &oldName)
	return QQ, oldName, err
}
