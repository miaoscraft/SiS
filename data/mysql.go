package data

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func openDB(addr, user, pswd, schema string) (err error) {
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pswd, addr, schema))
	if err != nil {
		return err
	}

	if err := initDB(schema); err != nil {
		return err
	}

	return nil
}

// 初始化数据库，自动检查表是否完整，重载储存过程
func initDB(schema string) error {
	// 建库
	if _, err := db.Exec(`create database if not exists ? character set UTF8;`, schema); err != nil {
		return err
	}

	// 建表
	if _, err := db.Exec(`
create table table_name
(
	QQ bigint null,
	Name text null,
	UUID binary(16) null,
	constraint table_name_pk
		primary key (QQ)
);
`); err != nil {
		return err
	}

	if _, err := db.Exec(`create unique index table_name_UUID_uindex on players (UUID);`); err != nil {
		return err
	}

	// 加载储存过程
	if _, err := db.Exec(`drop procedure if exists SeizeUUID;`); err != nil {
		return err
	}
	if _, err := db.Exec(`
create procedure SeizeUUID(in MyQQ bigint, in MyName text, in MyUUID binary(16))
begin
    declare oldName text;
    select Name into oldName from players where QQ = MyQQ;

    insert ignore into players
        (QQ, Name, UUID)
    values (MyQQ, MyName, MyUUID)
    on duplicate key update Name=MyName, UUID=MyUUID;

    select QQ, oldName
    from players
    where UUID = MyUUID;
end;
`); err != nil {
		return err
	}
	return nil
}

// SetWhitelist 尝试向数据库写入白名单数据，当ID未被占用时返回自己的QQ，当ID被占用则返回占用者的QQ
// 若原本该账号占有一个UUID，则会返回当时的UUID
func SetWhitelist(QQ int64, name string, ID uuid.UUID) (int64, *string, error) {
	rows, err := db.Query("call SeizeUUID(?,?,?);", QQ, name, ID[:])
	if err != nil {
		return 0, nil, err
	}
	if rows.Next() {
		var owner int64
		var oldName *string
		err = rows.Scan(&owner, &oldName)
		return owner, oldName, err
	}

	return 0, nil, errors.New("数据库没有返回数据")
}

// UnsetWhitelist 从数据库获取玩家绑定的ID，返回UUID并删除记录
func UnsetWhitelist(QQ int64) (uuid.UUID, bool, error) {
	rows, err := db.Query("select UUID from players where QQ=?", QQ)
	if err != nil {
		return uuid.Nil, false, err
	}
	// 先读数据
	var UUID uuid.UUID
	if rows.Next() {
		err = rows.Scan(&UUID)
		if err != nil {
			return uuid.Nil, false, err
		}
	}

	// 然后删除
	_, err = db.Exec("delete from players where QQ=?", QQ)
	if err != nil {
		return uuid.Nil, false, err
	}

	//返回的是读出来的数据
	return UUID, true, nil
}
