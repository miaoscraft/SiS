package data

import (
	"github.com/google/uuid"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func TestOpenDatabase_sqlite3(t *testing.T) {
	err := openDB("sqlite3", "data.db")
	//err := openDB("mysql",
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("data.db")
	defer closeDB()

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
	//if err != nil {
	//	t.Fatal(err)
	//}

	err = closeDB()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetLevel(t *testing.T) {
	err := openDB("sqlite3", "data.db")
	//err := openDB("mysql",
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("data.db")
	defer closeDB()

	if level, err := GetLevel(3261340757); err != nil {
		t.Fatal(err)
	} else {
		t.Log(level)
	}

	if err := SetLevel(3261340757, 12); err != nil {
		t.Fatal(err)
	}

	if level, err := GetLevel(3261340757); err != nil {
		t.Fatal(err)
	} else {
		t.Log(level)
	}
}
