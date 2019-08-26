package syntax

import (
	"github.com/miaoscraft/SiS/whitelist"
	"regexp"
	"strings"
)

var (
	// 指令前缀，通常为cq码[CQ:at,qq=<机器人qq>]
	CmdPrefix string
)

var expMyID = regexp.MustCompile(`(?i)MyID\s*[=＝]\s*([[:word:]]{3,16})`)

// GroupMsg 处理从游戏群接收到的消息，若为合法命令则进行相应的处理。并发安全
// 返回值指示是否拦截本消息
func GroupMsg(from int64, msg string, ret func(msg string)) bool {
	// 识别@指令
	if strings.HasPrefix(msg, CmdPrefix) {

		return true
	}

	// 识别MyID指令
	if match := expMyID.FindStringSubmatch(msg); len(match) == 2 {
		whitelist.MyID(from, match[1], ret)
		return true
	}

	return false
}
