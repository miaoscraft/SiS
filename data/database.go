package data

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

type Invited struct {
	gorm.Model
	Name string `form:"Name" gorm:"not null;unique"`
	By   string
}

type User struct {
	QQ   int64  `gorm:"column:QQ;primary_key"`
	UUID []byte `gorm:"column:UUID;not null;unique;type:binary(16)"`
}

type Auth struct {
	QQ    int64 `gorm:"column:QQ;primary_key"`
	Level int64 `gorm:"column:Level;default:0"`
}

// 启动数据库
func openDB(driver, source string) (err error) {
	db, err = gorm.Open(driver, source)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	if errs := initDB(); errs != nil {
		return fmt.Errorf("初始化数据库失败: %v", errs)
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
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}
	if !db.HasTable(&Auth{}) {
		db.CreateTable(&Auth{})
	}
	return db.Error
}

// SetWhitelist 尝试向数据库写入白名单数据，当ID未被占用时返回自己的QQ，当ID被占用则返回占用者的QQ
// 若原本该账号占有一个UUID，则会返回当时的UUID
func SetWhitelist(QQ int64, ID uuid.UUID, onOldID func(oldID uuid.UUID) error) (owner int64, err error) {
	tx := db.Begin()
	if err = tx.Error; err != nil {
		return
	}

	// 在函数结束时根据err判断是否应该Rollback或者Commit
	defer func() {
		if err == nil {
			if err1 := tx.Commit().Error; err1 != nil {
				err = fmt.Errorf("数据库提交数据失败: %v", err1)
			}
		} else {
			if err1 := tx.Rollback().Error; err1 != nil {
				err = fmt.Errorf("数据库操作失败: %v，且无法回滚数据: %v", err, err1)
			}
		}
	}()

	//事先排除id被占用的情况
	qq, err1 := GetWhitelistByUUID(ID)
	if err1 != nil {
		return 0, err
	}
	//id有人在使用，直接返回
	if qq != 0 {
		return qq, nil
	}

	var user User
	// 判断qq号是否存在记录,没有就创建有就更新
	if gorm.IsRecordNotFoundError(tx.First(&user, QQ).Error) {
		return 0, tx.Create(&User{QQ: QQ, UUID: ID[:]}).Error
	}
	//从mc服务器删除老旧uuid
	var uu uuid.UUID
	if err := uu.Scan(user.UUID); err != nil {
		return 0, fmt.Errorf("数据库错误，无法从游戏服务器删除旧的uuid:%v", err)
	}
	if err1 := onOldID(uu); err1 != nil {
		return 0, fmt.Errorf("从mc服务器删除白名单失败:%v", err1)
	}
	user.UUID = ID[:]
	if err = tx.Save(&user).Error; err != nil {
		return 0, fmt.Errorf("将白名单保存到游戏服务器后更新数据库错误：%v", err)
	}
	return
}

// UnsetWhitelist 从数据库获取玩家绑定的ID，返回UUID并删除记录
func UnsetWhitelist(QQ int64, onHas func(ID uuid.UUID) error) (err error) {
	tx := db.Begin()
	if err := db.Error; err != nil {
		return fmt.Errorf("数据库开始事务失败: %v", err)
	}

	// 在函数结束时根据err判断是否应该Rollback或者Commit
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
				err = fmt.Errorf("数据库操作失败: %v，且无法回滚数据: %v", err, rollbackErr)
			}
		} else {
			if err = tx.Commit().Error; err != nil {
				err = fmt.Errorf("数据提交失败: %v", err)
			}
		}
	}()

	var user User
	if err := tx.Where("QQ = ?", QQ).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil // 没有数据
	} else if err != nil {
		return fmt.Errorf("数据库查询UUID失败: %v", err)
	}
	var uu uuid.UUID
	if err := uu.Scan(user.UUID); err != nil {
		return err
	}
	if err := onHas(uu); err != nil {
		return err
	}
	//如果主键为空gorm会删掉所有记录，非常危险需要提前检查一下
	if user.QQ == 0 {
		return fmt.Errorf("没有完全查询到")
	}
	if err := tx.Delete(&user).Error; err != nil {
		return fmt.Errorf("数据库删除UUID失败: %v", err)
	}

	return nil
}

// GetWhitelistByQQ 从数据库读取玩家绑定的ID，若没有绑定ID则返回uuid.Nil
func GetWhitelistByQQ(QQ int64) (uuid.UUID, error) {
	var user User
	err := db.Where("QQ=?", QQ).First(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return uuid.Nil, nil
	}
	if err != nil {
		return uuid.Nil, err
	}
	var uu uuid.UUID
	if err := uu.Scan(user.UUID); err != nil {
		return uuid.Nil, err
	}
	return uu, nil
}

// GetWhitelistByUUID 从数据库读取绑定ID的玩家，若ID没有被绑定则则返回0
func GetWhitelistByUUID(ID uuid.UUID) (int64, error) {
	var user User
	err := db.Where("UUID=?", ID[:]).First(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return user.QQ, nil
}

// GetLevel 获取某人的权限等级
func GetLevel(QQ int64) (level int64, err error) {
	var auth Auth
	err = db.Where("QQ=?", QQ).First(&auth).Error
	if err == gorm.ErrRecordNotFound {
		level = 0
		err = nil
	} else if err != nil {
		err = fmt.Errorf("查询Level失败: %v", err)
	}
	level = auth.Level
	return
}

// SetLevel 设置某人的权限等级
func SetLevel(QQ, level int64) (err error) {
	tx := db.Begin()
	if err != nil {
		err = fmt.Errorf("数据库开始事务失败: %v", err)
		return
	}
	defer func() {
		if err == nil {
			if err1 := tx.Commit().Error; err1 != nil {
				err = fmt.Errorf("数据库提交数据失败: %v", err)
			}
		} else {
			if err1 := tx.Rollback().Error; err1 != nil {
				err = fmt.Errorf("数据库操作失败: %v，且无法回滚数据: %v", err, err1)
			}
		}
	}()
	var auth Auth
	auth.QQ = QQ

	// 查询是否有记录
	err1 := tx.Where("QQ=?", QQ).First(&auth).Error

	auth.Level = level
	//没有就创建
	if gorm.IsRecordNotFoundError(err1) {
		return tx.Create(&auth).Error
	}
	if err1 != nil {
		return err1
	}
	//存在就更新
	return tx.Save(&auth).Error

}
