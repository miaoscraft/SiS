package ping

import (
	"encoding/json"
	"fmt"
	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"strings"
	"text/template"
	"time"
)

func Ping(addr string, port int, ret func(msg string)) {
	resp, delay, err := bot.PingAndList(addr, port)
	if err != nil {
		ret(fmt.Sprintf("请求失败: %v", err))
		return
	}

	var s status
	err = json.Unmarshal(resp, &s)
	if err != nil {
		ret(fmt.Sprintf("解码失败: %v", err))
		return
	}
	s.Delay = delay // 延迟用手动填进去

	ret(String(s))

	return
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

	Delay time.Duration
}

var tmp = template.Must(template.
	New("PingRet").
	Parse(`服务器: [{{ .Version.Protocol }}] {{ .Version.Name }}
每日消息: {{ .Description }}
在线人数: {{ .Players.Online -}}/{{- .Players.Max }}
玩家列表:
{{ range .Players.Sample }}- [{{ .Name }}]
{{ end }}延迟: {{.Delay }}`))

func String(s status) string {
	var sb strings.Builder
	err := tmp.Execute(&sb, s)
	if err != nil {
		return fmt.Sprintf("文字模版渲染失败: %v", err)
	}
	return sb.String()
}
