package data

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// 启动数据库
func openDB(driver, source string) (err error) {
	db, err = sql.Open(driver, source)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	err = initDB()
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %v", err)
	}

	return nil
}

// 关闭数据库
func closeDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// 初始化数据库
func initDB() error {
	// "QQ->UUID", "UUID->QQ", "QQ->Level",
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
		    QQ INTEGER PRIMARY KEY ,
		    UUID BLOB NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS auths(
		    QQ INTEGER PRIMARY KEY ,
		    Level INT DEFAULT 0
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// SetWhitelist 尝试向数据库写入白名单数据，当ID未被占用时返回自己的QQ，当ID被占用则返回占用者的QQ
// 若原本该账号占有一个UUID，则会返回当时的UUID
func SetWhitelist(QQ int64, ID uuid.UUID, onOldID func(oldID uuid.UUID) error, onSuccess func() error) (owner int64, err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		err = fmt.Errorf("数据库开始事务失败: %v", err)
		return
	}

	// 在函数结束时根据err判断是否应该Rollback或者Commit
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("%v，且无法回滚数据: %v", err, rollbackErr)
			}
		} else {
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("数据提交失败: %v", err)
			}
		}
	}()

	// 检查UUID是否被他人占用
	err = tx.QueryRow("SELECT QQ FROM users WHERE UUID=?", ID[:]).Scan(&owner)
	switch err {
	case sql.ErrNoRows:
		// 没有被占用
		owner = QQ
		err = nil
	case nil:
		// 被占用了
		if owner != QQ {
			return
		}
	default:
		// 查询出错
		err = fmt.Errorf("数据库查询是否有占用者失败: %v", err)
		return
	}

	// 查询是否有旧白名单
	var oldID uuid.UUID
	err = tx.QueryRow("SELECT UUID FROM users WHERE QQ=?", QQ).Scan(&oldID)

	switch err {
	case nil: // 有旧的UUID
		// 消除旧账号白名单
		if err = onOldID(oldID); err != nil {
			return
		}

		// 更新UUID
		if _, err = tx.Exec("UPDATE users SET UUID=? WHERE QQ=?", ID[:], QQ); err != nil {
			err = fmt.Errorf("数据库更新UUID失败: %v", err)
			return
		}
	case sql.ErrNoRows: // 没有旧UUID
		if _, err = tx.Exec("INSERT INTO users (QQ, UUID) VALUES (?,?)", QQ, ID[:]); err != nil {
			err = fmt.Errorf("数据库插入UUID失败: %v", err)
			return
		}
		err = nil

	default: //查询出错
		err = fmt.Errorf("查询旧UUID失败: %v", err)
		return
	}

	// 更新玩家UUID
	if err = onSuccess(); err != nil {
		return
	}
	return
}

// UnsetWhitelist 从数据库获取玩家绑定的ID，返回UUID并删除记录
func UnsetWhitelist(QQ int64, onHas func(ID uuid.UUID) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("数据库开始事务失败: %v", err)
	}

	// 在函数结束时根据err判断是否应该Rollback或者Commit
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("数据库操作失败: %v，且无法回滚数据: %v", err, rollbackErr)
			}
		} else {
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("数据提交失败: %v", err)
			}
		}
	}()

	rows, err := tx.Query("SELECT UUID FROM users WHERE QQ=?", QQ)
	if err != nil {
		return fmt.Errorf("数据库查询UUID失败: %v", err)
	}

	if rows.Next() {
		var oldID uuid.UUID
		if err := rows.Scan(&oldID); err != nil {
			return fmt.Errorf("数据库读取旧UUID失败: %v", err)
		}

		if err := onHas(oldID); err != nil {
			return err
		}

		if _, err := tx.Exec("DELETE FROM users WHERE QQ=?", QQ); err != nil {
			return fmt.Errorf("数据库删除UUID失败: %v", err)
		}
	}
	return nil
}

// GetWhitelistByQQ 从数据库读取玩家绑定的ID，若没有绑定ID则返回uuid.Nil
func GetWhitelistByQQ(QQ int64) (id uuid.UUID, err error) {
	err = db.QueryRow("SELECT UUID FROM users WHERE QQ=?", QQ).Scan(&id)
	if err == sql.ErrNoRows {
		return uuid.Nil, nil
	}
	if err != nil {
		return uuid.Nil, err
	}

	return
}

// GetWhitelistByUUID 从数据库读取绑定ID的玩家，若ID没有被绑定则则返回0
func GetWhitelistByUUID(ID uuid.UUID) (qq int64, err error) {
	err = db.QueryRow("SELECT QQ FROM users WHERE UUID=?", ID[:]).Scan(&qq)
	if err == sql.ErrNoRows {
		return qq, nil
	}

	return
}

// GetLevel 获取某人的权限等级
func GetLevel(QQ int64) (level int64, err error) {
	err = db.QueryRow("SELECT Level FROM auths WHERE QQ=?", QQ).Scan(&level)
	if err == sql.ErrNoRows {
		level = 0
		err = nil
	} else if err != nil {
		err = fmt.Errorf("查询Level失败: %v", err)
	}

	return
}

// SetLevel 设置某人的权限等级
func SetLevel(QQ, level int64) (err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		err = fmt.Errorf("数据库开始事务失败: %v", err)
		return
	}

	// 查询是否有记录
	var rows *sql.Rows
	rows, err = tx.Query("SELECT Level FROM auths WHERE QQ=?", QQ)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("数据库操作失败: %v，且无法回滚数据: %v", err, rollbackErr)
		}
		return fmt.Errorf("数据库查询等级失败: %v", err)
	}

	// 根据数据存在性判断采用INSERT还是UPDATE
	if rows.Next() {
		if err = rows.Close(); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return fmt.Errorf("关闭rows失败: %v，且无法回滚数据: %v", err, rollbackErr)
			}
			return fmt.Errorf("关闭rows失败: %v", err)
		}
		_, err = tx.Exec("UPDATE auths SET Level=? WHERE QQ=?", level, QQ)
	} else {
		_, err = tx.Exec("INSERT INTO auths (QQ, Level) VALUES (?,?)", QQ, level)
	}
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("数据库操作失败: %v，且无法回滚数据: %v", err, rollbackErr)
		}
		return fmt.Errorf("数据库操作失败: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("数据库提交数据失败: %v", err)
	}
	return
}
