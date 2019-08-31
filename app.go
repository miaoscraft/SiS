package main

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/miaoscraft/SiS/data"
	"github.com/miaoscraft/SiS/log"
	"github.com/miaoscraft/SiS/syntax"
	"github.com/miaoscraft/SiS/whitelist"
	"runtime/debug"
)

//go:generate cqcfg -c .
// cqp: 名称: SiS
// cqp: 版本: 1.0.0:0
// cqp: 作者: Tnze
// cqp: 简介: Minecraft服务器综合管理器
func main() { /*空*/ }

func init() {
	cqp.AppID = "cn.miaoscraft.sis"
	cqp.Enable = onStart
	cqp.Disable = onStop
	cqp.Exit = onStop

	cqp.GroupMsg = onGroupMsg
	cqp.GroupMemberDecrease = onGroupMemberDecrease

	whitelist.Logger = log.NewLogger("MyID")
	data.Logger = log.NewLogger("Data")
}

var Logger = log.NewLogger("Main")

// 插件生命周期开始
func onStart() int32 {
	defer panicConvert()

	// 连接数据源
	err := data.Init(cqp.GetAppDir())
	if err != nil {
		Logger.Errorf("初始化数据源失败: %v", err)
	}

	// 将登录账号载入命令解析器（用于识别@）
	syntax.CmdPrefix = fmt.Sprintf("[CQ:at,qq=%d]", cqp.GetLoginQQ())

	return 0
}

// 插件生命周期结束
func onStop() int32 {
	defer panicConvert()

	err := data.Close()
	if err != nil {
		Logger.Errorf("释放数据源失败: %v", err)
	}
	return 0
}

// 群消息事件
func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	defer panicConvert()

	if fromQQ == 80000000 { // 忽略匿名
		return Ignore
	}

	ret := func(resp string) {
		cqp.SendGroupMsg(fromGroup, resp)
	}

	switch fromGroup {
	case data.Config.AdminID:
		// 当前版本，管理群和游戏群收到的命令不做区分
		fallthrough
	case data.Config.GroupID:
		if syntax.GroupMsg(fromQQ, msg, ret) {
			return Intercept
		}
	}
	return Ignore
}

// 群成员减少事件
func onGroupMemberDecrease(subType, sendTime int32, fromGroup, fromQQ, beingOperateQQ int64) int32 {
	defer panicConvert()

	retValue := Ignore
	ret := func(resp string) {
		cqp.SendGroupMsg(fromGroup, resp)
		retValue = Intercept
	}
	// 尝试删白名单
	if fromGroup == data.Config.GroupID {
		whitelist.RemoveWhitelist(beingOperateQQ, ret)
	}
	return retValue
}

const (
	Ignore    int32 = 0 //忽略消息
	Intercept       = 1 //拦截消息
)

// 用于捕获所有panic，转换为酷Q的Fatal日志
func panicConvert() {
	if v := recover(); v != nil {
		// 在这里调用debug.Stack()获取调用栈
		Logger.Errorf("%v\n%s", v, debug.Stack())
	}
}
