package data

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 创建一些需要但不存在的文件
func initFiles() error {
	load := func(f *os.File, content string) error {
		_, err := io.Copy(f, strings.NewReader(content))
		if err != nil {
			_ = f.Close()
			return err
		}
		return f.Close()
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

var defaultFiles = map[string]string{
	"conf.toml": `# SiS配置文件

GroupID = 123456789 # 游戏群群号
# AdminID = 123456789 # 管理群群号（可选）

[Ping] # Ping工具配置
DefaultServer = "play.miaoscraft.cn" # 默认目标服务器[:端口]，端口是可选的，默认为25565
Timeout = "60s" # 最长ping时间，为0时禁用。例如："300ms", "1.5h" 或 "2h45m"。可用的单位有 纳秒"ns", 微妙"us" (或 "µs"), 毫秒"ms", 秒"s", 分钟"m", 小时"h".

[RCON] # RCON配置
Address = "127.0.0.1:25575" #服务器地址:端口，必须写上端口
Password = "your_password" #服务器RCON密码，server.properties文件里的rcon.password`,
}
