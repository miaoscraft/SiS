package data

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"path/filepath"
)

// 创建一些需要但不存在的文件
func initFiles() error {
	load := func(f *os.File, content []byte) error {
		defer f.Close()

		_, err := io.Copy(f, bytes.NewReader(content))
		if err != nil {
			return err
		}

		return nil
	}
	for fileName, fileContent := range defaultFiles {
		f, err := os.OpenFile(filepath.Join(AppDir, fileName), os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
		if os.IsExist(err) {
			continue
		} else if err != nil {
			return err
		}

		err = load(f, fileContent)
		if err != nil {
			return err
		}
	}
	return nil
}

//go:embed conf.toml
var confToml []byte

var defaultFiles = map[string][]byte{
	"conf.toml": confToml,
}
