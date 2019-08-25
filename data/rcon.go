package data

import (
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/Tnze/go-mc/net"
)

var rcon net.RCONClientConn

var reopenRCON func() error

func openRCON(address, password string) (err error) {
	reopenRCON = func() error {
		rcon, err = net.DialRCON(address, password)
		return err
	}
	return reopenRCON()
}

// rconCommand 执行RCON命令，断线时尝试重连一次
func rconCommand(cmd string) (string, error) {
ReTry:
	err := rcon.Cmd(cmd)
	if err != nil {
		cqp.AddLog(cqp.Error, "rcon", "rcon添加白名单失败: "+err.Error())
		// 断线重连
		err = reopenRCON()
		if err != nil {
			return "", err
		}
		goto ReTry
	}

	resp, err := rcon.Resp()
	if err != nil {
		cqp.AddLog(cqp.Error, "rcon", "读rcon返回值失败: "+err.Error())
		// 不重连
		return "", err
	}

	cqp.AddLog(cqp.Info, "rcon", "RCON: "+resp)
	return resp, nil
}

func AddWhitelist(name string) error {
	_, err := rconCommand("whitelist add " + name)
	return err
}

func RemoveWhitelist(name string) error {
	_, err := rconCommand("whitelist add " + name)
	return err
}
