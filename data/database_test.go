package data

// 两个不是很完善的测试

import (
	"github.com/google/uuid"
	"testing"
)

func TestOpenDatabase_sqlite3(t *testing.T) {
	err := openDB("sqlite3", "data.db")
	if err != nil {
		t.Fatal(err)
	}
	//defer os.Remove("data.db")
	defer func() {
		err := closeDB()
		if err != nil {
			t.Fatal(err)
		}
	}()

	owner, err := SetWhitelist(3261340757, uuid.MustParse("58f6356eb30c48118bfcd72a9ee99e74"),
		func(id uuid.UUID) error {
			t.Log("old id:", id)
			return nil
		},
		func() error {
			t.Log("success")
			return nil
		})
	t.Log("owner:", owner)

	//err = UnsetWhitelist(3261340757, func(id uuid.UUID) error {
	//	t.Log(id)
	//	return nil
	//})
	if err != nil {
		t.Fatal(err)
	}

	err = closeDB()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLevel(t *testing.T) {
	err := openDB("sqlite3", "data.db")
	if err != nil {
		t.Fatal(err)
	}
	//defer os.Remove("data.db")
	defer func() {
		err := closeDB()
		if err != nil {
			t.Fatal(err)
		}
	}()

	if level, err := GetLevel(3261340757); err != nil {
		t.Fatal(err)
	} else {
		t.Log(level)
	}

	if err := SetLevel(3261340757, 11); err != nil {
		t.Fatal(err)
	}

	if level, err := GetLevel(3261340757); err != nil {
		t.Fatal(err)
	} else {
		t.Log(level)
	}
}
