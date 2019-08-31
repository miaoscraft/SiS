package data

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"path/filepath"
	"time"
)

// AppDir 当前插件数据目录
var AppDir string

// Init 初始化插件的数据源，包括读取配置文件、建立数据库连接
func Init(dir string) error {
	AppDir = dir
	// 初始化默认文件
	err := initFiles()
	if err != nil {
		return fmt.Errorf("创建文件时出错: %v", err)
	}
	// 读取配置文件
	err = readConfig()
	if err != nil {
		return fmt.Errorf("读配置文件出错: %v", err)
	}

	// 连接数据库
	err = openDB(filepath.Join(AppDir, "data.db"))
	if err != nil {
		return fmt.Errorf("打开数据库出错: %v", err)
	}

	// 连接MC服务器
	err = openRCON(Config.RCON.Address, Config.RCON.Password)
	if err != nil {
		return fmt.Errorf("连接RCON出错: %v", err)
	}

	return nil
}

// Close 关闭所有打开的资源
func Close() error {
	err := closeDB()
	if err != nil {
		return err
	}
	return nil
}

var Config struct {
	// 游戏群
	GroupID int64
	// 管理群
	AdminID int64
	// MC服务器远程控制台
	RCON struct {
		Address  string
		Password string
	}
	// Ping工具配置
	Ping struct {
		DefaultServer string
		Timeout       duration
	}

	// 自定义命令
	Cmd map[string]struct {
		Level   int64
		Command string
		Silent  bool
	}
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func readConfig() error {
	md, err := toml.DecodeFile(filepath.Join(AppDir, "conf.toml"), &Config)
	if err != nil {
		return err
	}

	// 检查配置文件是否有多余数据，抛警告⚠️
	if uk := md.Undecoded(); len(uk) > 0 {
		Logger.Waringf("配置文件中有未知数据: %q", uk)
	}

	return nil
}
