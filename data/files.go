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
		defer f.Close()

		_, err := io.Copy(f, strings.NewReader(content))
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

var defaultFiles = map[string]string{
	"conf.toml": `# SiS配置文件，在Wiki中有详细的配置说明：https://github.com/miaoscraft/SiS/wiki

GroupID = 123456789 # 游戏群群号
# AdminID = 123456789 # 管理群群号（可选）
# Administrators = [ 12345678910, 23456789101, 34567891011] # 设置管理员们（他们可以设置任何人的Level)

[Onebot]
Addr = "localhost:6700"
token = ""

[DealWithGroupRequest]# 处理入群请求。详细见：https://github.com/miaoscraft/SiS/wiki/%E8%87%AA%E5%8A%A8%E5%A4%84%E7%90%86%E5%85%A5%E7%BE%A4%E8%AF%B7%E6%B1%82
Enable = true # 启用
CanReject = false # 是否允许机器人拒绝请求
CheckURL = ""
Token = ""

[Database]
Driver = "sqlite3" # 数据库类型（仅支持mysql和sqlite3）
Source = "{{ join .AppDir \"data.db\"}}" # SQLite写法, 详细用法见https://github.com/mattn/go-sqlite3#dsn-examples
# Source = "用户:密码@tcp(地址:端口)/库名" # MySQL写法, 详细用法见https://github.com/go-sql-driver/mysql#dsn-data-source-name

[Ping] # Ping工具配置
DefaultServer = "play.miaoscraft.cn" # 默认目标服务器[:端口]，端口是可选的，默认为25565
Timeout = "60s" # 最长ping时间，为0时禁用。例如："300ms", "1.5h" 或 "2h45m"。可用的单位有 纳秒"ns", 微妙"us" (或 "µs"), 毫秒"ms", 秒"s", 分钟"m", 小时"h".

[RCON] # RCON配置
Address = "127.0.0.1:25575" #服务器地址:端口，必须写上端口
Password = "rcon_password" #服务器RCON密码，server.properties文件里的rcon.password

# 自定义命令配置
[Cmd.tps] # 命令名
Level = 0 # 执行该命令所需等级
Command = "tps" # 执行时实际发送的命令
# Silent = true # 禁用命令回显
# AllowArgs = true # 允许执行时附加参数（试验性功能，请自行考察安全性后慎重开启）

[Cmd."帮助"] # 中文命令需要引号，命令不可包含空格
Level = 0
Command = "help"
`,
}
