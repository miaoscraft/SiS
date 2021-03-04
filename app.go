package main

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/customize"
	"github.com/miaoscraft/SiS/data"
	"github.com/miaoscraft/SiS/log"
	"github.com/miaoscraft/SiS/syntax"
	"github.com/miaoscraft/SiS/whitelist"
	"net/http"
	"net/url"
	"regexp"
)

//go:generate cqcfg -c .
// cqp: 名称: SiS
// cqp: 版本: 1.4.1:1
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
	cqp.GroupRequest = onGroupRequest

	customize.Logger = log.NewLogger("Cstm")
	whitelist.Logger = log.NewLogger("MyID")
	data.Logger = log.NewLogger("Data")
}

var Logger = log.NewLogger("Main")

// 插件生命周期开始
func onStart() int32 {
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
	err := data.Close()
	if err != nil {
		Logger.Errorf("释放数据源失败: %v", err)
	}
	return 0
}

// 群消息事件
func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
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

// 入群请求事件
func onGroupRequest(subType, sendTime int32, fromGroup, fromQQ int64, msg, respFlag string) int32 {
	if !data.Config.DealWithGroupRequest.Enable || // 功能未启用
		fromGroup != data.Config.GroupID || // 不是要管理的群
		subType != 1 { // 不是他人要申请入群
		return Ignore
	}
	Logger := log.NewLogger("DwGR")
	for _, name := range regexp.MustCompile(`[0-9A-Za-z_]{3,16}`).FindAllString(msg, 3) {
		name, id, err := whitelist.GetUUID(name)
		if err != nil {
			Logger.Infof("处理%d的入群请求，检查游戏名失败: %v", fromQQ, err)
			continue
		}

		if ok, err := checkRequest(name, id); err != nil {
			Logger.Errorf("服务器检查出错: %v", err)
			return Ignore
		} else if !ok {
			Logger.Infof("服务器拒绝%d作为%s入群: %v", fromQQ, name, err)
			if data.Config.DealWithGroupRequest.CanReject {
				cqp.SetGroupAddRequest(respFlag, subType, Deny, "")
				return Intercept
			}
			return Ignore
		}
		Logger.Infof("允许%d作为%s入群", fromQQ, name)
		cqp.SetGroupAddRequest(respFlag, subType, Allow, "")

		ret := func(resp string) { cqp.SendGroupMsg(fromGroup, resp) }
		whitelist.MyID(fromQQ, name, ret)
		return Intercept
	}
	return Ignore
}

func checkRequest(name string, id uuid.UUID) (bool, error) {
	if data.Config.DealWithGroupRequest.CheckURL == "" {
		return true, nil
	}
	resp, err := http.PostForm(
		data.Config.DealWithGroupRequest.CheckURL,
		url.Values{
			"Token": []string{data.Config.DealWithGroupRequest.Token},
			"Name":  []string{name},
			"UUID":  []string{id.String()},
		},
	)
	if err != nil {
		return false, fmt.Errorf("请求出错: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	return false, nil
}

const (
	Ignore    int32 = 0 //忽略消息
	Intercept       = 1 //拦截消息

	Allow = 1 // 允许进群
	Deny  = 2 // 拒绝进群
)
