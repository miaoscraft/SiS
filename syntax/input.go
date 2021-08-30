// Package syntax 实现SiS机器人支持的语法的解析和执行
package syntax

import (
	"github.com/BaiMeow/SimpleBot/message"
	"github.com/miaoscraft/SiS/customize"
	"github.com/miaoscraft/SiS/ping"
	"github.com/miaoscraft/SiS/whitelist"
	"regexp"
)

//ID 登陆qq号，为识别@命令作准备
var ID int64

var expMyID = regexp.MustCompile(`^\s*(?i)MyID\s*[=＝]\s*([0-9A-Za-z_]{3,16})\s*$`)

// GroupMsg 处理从游戏群接收到的消息，若为合法命令则进行相应的处理。并发安全
// 返回值指示是否拦截本消息
func GroupMsg(from int64, msg message.Msg, ret func(msg string)) bool {
	// 识别MyID指令
	if len(msg) == 1 && msg[0].GetType() == "text" {
		if match := expMyID.FindStringSubmatch(msg[0].(message.Text).Text); len(match) == 2 {
			whitelist.MyID(from, match[1], ret)
			return true
		}
	}
	//@命令至少有两个消息段
	if len(msg) < 2 {
		return false
	}
	at, ok := msg[0].(message.At)
	if !ok || !at.IsAt(ID) {
		return false
	}
	args := msg[1:].Fields()
	if args[0].GetType() != "text" {
		return false
	}

	switch args[0].(message.Text).Text {
	case "ping": // ping指令
		return ping.Ping(args, ret)

	case "auth": // auth指令
		return customize.Auth(args, from, ret)

	case "info": // 白名单查询指令
		return whitelist.Info(args, from, ret)

	default: // 自定义指令
		return customize.Exec(args, from, ret)
	}
}
