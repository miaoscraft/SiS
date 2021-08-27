// Package ping 提供内置指令"ping"的实现
package ping

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BaiMeow/SimpleBot/message"
	"net"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
)

func Ping(msg message.Msg, ret func(msg string)) bool {
	var args []string
	for _, m := range msg {
		if m.GetType() != "text" {
			continue
		}
		args = append(args, m.(message.Text).Text)
	}
	var (
		resp  []byte
		delay time.Duration
		err   error
	)

	addresses := getAddr(args)
	statuses := make([]status, len(addresses))
	for i, addr := range addresses {
		if d := data.Config.Ping.Timeout.Duration; d > 0 {
			//启用Timeout
			resp, delay, err = bot.PingAndListTimeout(addr, d)
		} else {
			//禁用Timeout
			resp, delay, err = bot.PingAndList(addr)
		}
		if err != nil {
			statuses[i].Error = err
			continue
		}

		err = json.Unmarshal(resp, &statuses[i])
		if err != nil {
			statuses[i].Error = err
			continue
		}

		// 延迟用手动填进去
		statuses[i].Delay = delay
		statuses[i].Address = addr
	}
	ret(render(statuses))
	return true
}

// 从[]string获取服务器地址和端口
// 支持的格式有:
// 	[ "ping" "play.miaoscraft.cn" ]
// 	[ "ping" "play.miaoscraft.cn:25565" ]
// 	[ "ping" "play.miaoscraft.cn" "25565" ]
func getAddr(args []string) (addrs []string) {
	args = args[1:] //去除第一个元素"ping"
	switch len(args) {
	default: // len >= 2
		return []string{net.JoinHostPort(args[0], args[1])}
	case 0: // 默认值
		args = append(args, data.Config.Ping.DefaultServer)
		fallthrough
	case 1:
		var addrErr *net.AddrError
		const missingPort = "missing port in address"
		addr := args[0]
		if _, _, err := net.SplitHostPort(addr); errors.As(err, &addrErr) && addrErr.Err == missingPort {
			_, addrsSRV, err := net.LookupSRV("minecraft", "tcp", addr)
			if err == nil && len(addrsSRV) > 0 {
				for _, addrSRV := range addrsSRV {
					addrs = append(addrs, net.JoinHostPort(addrSRV.Target, strconv.Itoa(int(addrSRV.Port))))
				}
				return
			}
			return []string{net.JoinHostPort(addr, "25565")}
		} else {
			return []string{addr}
		}
	}
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

	Address string        `json:"-"`
	Delay   time.Duration `json:"-"`
	Error   error         `json:"-"`
}

var tmp = template.Must(template.
	New("PingRet").
	Parse(`喵哈喽～{{ $list := .}}
{{ with index . 0 }}服务器版本: [{{ .Version.Protocol }}] {{ .Version.Name }}
每日消息: {{ .Description.ClearString }}
{{ range $index, $elem := $list }}延迟: {{if .Error}}请求失败：{{ .Error }}{{ else }}{{ .Delay }}{{ end }}
{{ end }}在线人数: {{ .Players.Online -}}/{{- .Players.Max }}
玩家列表:
{{ range .Players.Sample }}- [{{ .Name }}]
{{ end }}{{ end }}にゃ～`))

func render(statuses []status) string {
	var sb strings.Builder
	err := tmp.Execute(&sb, statuses)
	if err != nil {
		return fmt.Sprintf("似乎在渲染文字模版时出现了棘手的问题: %v", err)
	}
	cleanStr, _ := chat.TransCtrlSeq(sb.String(), false)
	return cleanStr
}
