package data

import (
	"encoding/binary"
	bolt "github.com/etcd-io/bbolt"
	"github.com/google/uuid"
)

var db *bolt.DB

func openDB(name string) (err error) {
	db, err = bolt.Open(name, 0666, nil)
	if err != nil {
		return err
	}

	err = initDB()
	if err != nil {
		return err
	}

	return nil
}

// 初始化数据库
func initDB() error {
	return db.Update(func(tx *bolt.Tx) error {
		for _, b := range []string{
			"QQ->UUID", "UUID->QQ",
		} {
			_, err := tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// SetWhitelist 尝试向数据库写入白名单数据，当ID未被占用时返回自己的QQ，当ID被占用则返回占用者的QQ
// 若原本该账号占有一个UUID，则会返回当时的UUID
func SetWhitelist(QQ int64, ID uuid.UUID, onOldID func(oldID uuid.UUID) error, onSuccess func() error) (owner int64, err error) {
	bytesQQ := int64Bits(QQ)

	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		qu := tx.Bucket([]byte("QQ->UUID"))
		uq := tx.Bucket([]byte("UUID->QQ"))

		// 读UUID主人
		bytesOwner := uq.Get(ID[:])
		if len(bytesOwner) == 8 {
			owner = int64(binary.BigEndian.Uint64(bytesOwner))
			return nil //请求读UUID已经被占据了，不写入新的数据
		}
		owner = QQ

		// 读旧ID
		bytesID := qu.Get(bytesQQ)
		if len(bytesID) == 16 {
			var oldID uuid.UUID
			copy(oldID[:], bytesID)
			// 读到了就可以删了,qu不用删是因为马上就要写入新的数据
			err := uq.Delete(bytesID)
			if err != nil {
				return err
			}

			err = onOldID(oldID) // 此时删除旧白名单，若失败则回滚
			if err != nil {
				return err
			}
		}

		// 记录新的UUID和QQ
		err = qu.Put(bytesQQ, ID[:])
		if err != nil {
			return err
		}
		err = uq.Put(ID[:], bytesQQ)
		if err != nil {
			return err
		}

		err = onSuccess() // 添加白名单，若出错则回滚
		if err != nil {
			return err
		}

		return nil
	})

	return
}

// UnsetWhitelist 从数据库获取玩家绑定的ID，返回UUID并删除记录
func UnsetWhitelist(QQ int64, onHas func(ID uuid.UUID) error) error {
	var UUID uuid.UUID
	return db.Update(func(tx *bolt.Tx) error {
		qu := tx.Bucket([]byte("QQ->UUID"))
		copy(UUID[:], qu.Get(int64Bits(QQ)))

		if UUID != uuid.Nil { //删除白名单，若失败则回滚
			err := onHas(UUID)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func int64Bits(n int64) (b []byte) {
	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return
}
