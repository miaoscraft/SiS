package data

import "github.com/Tnze/go-mc/net"

var rcon net.RCONClientConn

func openRCON(address, password string) (err error) {
	rcon, err = net.DialRCON(address, password)
	return
}
