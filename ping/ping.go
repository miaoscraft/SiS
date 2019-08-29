// Package ping 提供内置指令"ping"的实现
package ping

import (
	"encoding/json"
	"fmt"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func Ping(args []string, ret func(msg string)) bool {
	var (
		resp  []byte
		delay time.Duration
		err   error
	)
	if d := data.Config.Ping.Timeout.Duration; d > 0 {
		//启用Timeout
		addr, port := getAddr(args)
		resp, delay, err = bot.PingAndListTimeout(addr, port, d)
	} else {
		//禁用Timeout
		resp, delay, err = bot.PingAndList(getAddr(args))
	}
	if err != nil {
		ret(fmt.Sprintf("嘶...请求失败惹！: %v", err))
		return true
	}

	var s status
	err = json.Unmarshal(resp, &s)
	if err != nil {
		ret(fmt.Sprintf("嘶...解码失败惹！: %v", err))
		return true
	}
	// 延迟用手动填进去
	s.Delay = delay

	ret(s.String())

	return true
}

// 从[]string获取服务器地址和端口
// 支持的格式有:
// 	[ "ping" "play.miaoscraft.cn" ]
// 	[ "ping" "play.miaoscraft.cn:25565" ]
// 	[ "ping" "play.miaoscraft.cn" "25565" ]
func getAddr(args []string) (addr string, port int) {
	args = args[1:] //去除第一个元素"ping"
	// 默认值
	addr = data.Config.Ping.DefaultServer
	port = 25565

	// 在第二个参数内寻找端口
	if len(args) >= 2 {
		if p, err := strconv.Atoi(args[1]); err == nil {
			port = p
		}
	}

	// 如果有则加载第一个参数
	if len(args) >= 1 {
		addr = args[0]
	}

	// 在冒号后面寻找端口
	f := strings.Split(addr, ":")
	if len(f) >= 2 {
		if p, err := strconv.Atoi(f[1]); err == nil {
			port = p
		}
	}

	// 冒号前面是地址
	addr = f[0]

	return addr, port
}

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
	//favicon ignored

	Delay time.Duration `json:"-"`
}

var tmp = template.Must(template.
	New("PingRet").
	Parse(`喵哈喽～
服务器版本: [{{ .Version.Protocol }}] {{ .Version.Name }}
Motd: {{ .Description.ClearString }}
延迟: {{.Delay }}
在线人数: {{ .Players.Online -}}/{{- .Players.Max }}
玩家列表:
{{ range .Players.Sample }}- [{{ .Name }}]
{{ end }}にゃ～`))

func (s status) String() string {
	var sb strings.Builder
	err := tmp.Execute(&sb, s)
	if err != nil {
		return fmt.Sprintf("似乎在渲染文字模版时出现了棘手的问题: %v", err)
	}
	return sb.String()
}
