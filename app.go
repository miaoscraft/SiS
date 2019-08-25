package main

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/miaoscraft/SiS/data"
	"github.com/miaoscraft/SiS/syntax"
)

//go:generate cqcfg .
// cqp: 名称: SiS
// cqp: 版本: 1.0.0:2
// cqp: 作者: Tnze
// cqp: 简介: Minecraft服务器综合管理器
func main() { /*空*/ }

func init() {
	cqp.AppID = "cn.miaoscraft.sis"
	cqp.Enable = onEnable
	cqp.GroupMsg = onGroupMsg

}

func onEnable() int32 {
	// 连接数据源
	err := data.Init()
	if err != nil {
		cqp.AddLog(cqp.Error, "Init", fmt.Sprintf("初始化数据源失败: %v", err))
	}

	// 将登录账号载入命令解析器（用于识别@）
	syntax.CmdPrefix = fmt.Sprintf("[CQ:at,qq=%d]", cqp.GetLoginQQ())

	return 0
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	switch fromGroup {
	case data.Config.AdminID:
		// 当前版本，管理群和游戏群收到的命令不做区分
		fallthrough

	case data.Config.GroupID:
		syntax.GroupMsg(fromQQ, msg,
			func(resp string) { //callback
				cqp.SendGroupMsg(fromGroup, resp)
			})
	}
	return 0
}
