package data

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
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
	err = openDB(Config.Database.Driver, Config.Database.Source)
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
	ZeroBot struct {
		Host, Port, AccessToken string
	}
	// 游戏群
	GroupID int64
	// 管理群
	AdminID int64
	// 管理员
	Administrators []int64
	// 处理进群请求
	DealWithGroupRequest struct {
		Enable    bool
		CanReject bool
		CheckURL  string
		Token     string
	}

	// 数据库配置
	Database struct {
		Driver string
		Source string
	}

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
		Level     int64  // 所需权限
		Command   string // 指令本身
		Silent    bool   // 是否不回显
		AllowArgs bool   // 是否允许使用参数
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
		Logger.Warningf("配置文件中有未知数据: %q", uk)
	}

	// 替换文件路径Database中Source的文件路径
	Config.Database.Source, err = rendingDBSource(Config.Database.Source)
	if err != nil {
		return err
	}

	return nil
}

func rendingDBSource(raw string) (string, error) {
	var sb strings.Builder
	temp, err := template.
		New("DBSource").
		Funcs(template.FuncMap{
			"join": filepath.Join,
		}).
		Parse(raw)
	if err != nil {
		return "", fmt.Errorf("解析模版失败: %v", err)
	}

	err = temp.Execute(&sb, struct{ AppDir string }{AppDir})
	if err != nil {
		return "", fmt.Errorf("渲染模版失败: %v", err)
	}

	return sb.String(), nil
}
