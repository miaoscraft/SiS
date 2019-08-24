package main

import "github.com/Tnze/CoolQ-Golang-SDK/cqp"

//go:generate cqcfg .
// cqp: 名称: GoDemo
// cqp: 版本: 1.0.0:2
// cqp: 作者: Tnze
// cqp: 简介: 一个超棒的Go语言插件Demo，它会回复你的私聊消息~
func main() { /*此处应当留空*/ }

func init() {
	cqp.AppID = "me.cqp.tnze.demo" // TODO: 修改为这个插件的ID
	cqp.PrivateMsg = onPrivateMsg
}

func onPrivateMsg(subType, msgID int32, fromQQ int64, msg string, font int32) int32 {
	cqp.SendPrivateMsg(fromQQ, msg) //复读机
	return 0
}
