package data

import "testing"

// 又是一个很烂的单元测试
func TestReadConfig(t *testing.T) {
	if err := initFiles(); err != nil {
		t.Fatal(err)
	}

	if err := readConfig(); err != nil {
		t.Fatal(err)
	}
	t.Log(Config)
}
