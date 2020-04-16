package data

import (
	"github.com/Tnze/go-mc/chat"
	mcnet "github.com/Tnze/go-mc/net"
	"github.com/fatih/pool"
	"math/rand"
	"net"
	"strings"
	"time"
)

var rcon pool.Pool

func openRCON(address, password string) (err error) {
	rcon, err = pool.NewChannelPool(1, 5, func() (net.Conn, error) {
		rcon, err := mcnet.DialRCON(address, password)
		return rcon.(*mcnet.RCONConn).Conn, err
	})
	return
}

// RCONCmd 执行RCON命令，断线时尝试重连一次
func RCONCmd(cmd string, ret func(string)) error {
	var r *mcnet.RCONConn
	for {
		conn, err := rcon.Get()
		if err != nil {
			return err
		}
		r = &mcnet.RCONConn{Conn: conn, ReqID: rand.Int31()}

		err = r.Cmd(cmd)
		if err == nil {
			break
		}
		Logger.Errorf("rcon执行失败: %v", err)
		if pc, ok := r.Conn.(*pool.PoolConn); ok {
			// close the underlying connection
			pc.MarkUnusable()
			pc.Close()
		}
	}
	go func() {
		defer r.Close()
		if ret == nil {
			return
		}
		for ret != nil {
			_ = r.SetWriteDeadline(time.Now().Add(time.Second * 10))
			resp, err := r.Resp()
			if err != nil {
				Logger.Infof("停止转发rcon返回值: %v", err)
				return
			}

			Logger.Infof("RCON返回: %q", resp)
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
