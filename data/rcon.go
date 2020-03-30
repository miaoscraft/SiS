package data

import (
	"errors"
	"github.com/Tnze/go-mc/chat"
	mcnet "github.com/Tnze/go-mc/net"
	"strings"
	"time"
)

var rconDialer func() (mcnet.RCONClientConn, error)

func openRCON(address, password string) error {
	rconDialer = func() (mcnet.RCONClientConn, error) {
		return mcnet.DialRCON(address, password)
	}
	return nil
}

// RCONCmd 执行RCON命令，每次都创建一个连接。
// 当ret不为nil时，通过回调方式返回RCON执行结果，ret可能被调用多次。
func RCONCmd(cmd string, ret func(string)) error {
	var r *mcnet.RCONConn

	if rconDialer == nil {
		return errors.New("RCON未设置")
	}
	conn, err := rconDialer()
	if err != nil {
		return err
	}
	r = conn.(*mcnet.RCONConn)

	err = r.Cmd(cmd)
	if err != nil {
		return err
	}

	go func() {
		defer r.Close()
		if ret == nil {
			return
		}
		tip := time.AfterFunc(time.Second, func() {
			ret("正在努力发送指令噢，请稍后~")
		})
		for {
			_ = r.SetDeadline(time.Now().Add(time.Second * 10))
			resp, err := r.Resp()
			if err != nil {
				Logger.Debugf("停止转发RCON返回值: %v", err)
				return
			}
			tip.Stop() // 不再发送提示

			Logger.Debugf("RCON返回: %q", resp)
			// 过滤掉末尾换行符、空格和零字符，过滤§格式字符串
			resp = chat.Message{Text: strings.TrimRight(resp, " \000\n")}.ClearString()
			ret(resp)
		}
	}()

	return nil
}

// AddWhitelist 从游戏服务器添加白名单
func AddWhitelist(name string) error {
	return RCONCmd("whitelist add "+name, nil)
}

// RemoveWhitelist 从游戏服务器删除白名单
func RemoveWhitelist(name string) error {
	return RCONCmd("whitelist remove "+name, nil)
}
