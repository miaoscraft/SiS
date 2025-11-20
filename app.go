package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/BaiMeow/SimpleBot/bot"
	"github.com/BaiMeow/SimpleBot/driver"
	"github.com/BaiMeow/SimpleBot/message"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/customize"
	"github.com/miaoscraft/SiS/data"
	"github.com/miaoscraft/SiS/log"
	"github.com/miaoscraft/SiS/syntax"
	"github.com/miaoscraft/SiS/whitelist"
)

var b *bot.Bot

func main() {
	customize.Logger = log.NewLogger("Cstm")
	whitelist.Logger = log.NewLogger("MyID")
	data.Logger = log.NewLogger("Data")
	// 连接数据源
	err := data.Init()
	if err != nil {
		Logger.Errorf("初始化数据源失败: %v", err)
	}
	//使用正向ws驱动
	b = bot.New(driver.NewWsDriver(data.Config.OneBot.Addr, data.Config.OneBot.Token))
	//注册事件
	b.Attach(&bot.GroupMsgHandler{
		Priority: 1,
		F:        onGroupMsg,
	})
	b.Attach(&bot.GroupAddHandler{
		Priority: 1,
		F:        onGroupRequest,
	})
	b.Attach(&bot.GroupDecreaseHandler{
		Priority: 1,
		F:        onGroupMemberDecrease,
	})
	err = b.Run()
	if err != nil {
		Logger.Error(err.Error())
	}
	// 将登录账号载入命令解析器（用于识别@）
	syntax.ID = b.GetID()
	defer onStop()
	select {}
}

var Logger = log.NewLogger("Main")

// 插件生命周期结束
func onStop() {
	err := data.Close()
	if err != nil {
		Logger.Errorf("释放数据源失败: %v", err)
	}
}

// 群消息事件
func onGroupMsg(MsgID int32, fromGroup, fromQQ int64, Msg message.Msg) bool {
	ret := func(resp string) {
		b.SendGroupMsg(fromGroup, message.CQstrToArrayMessage(resp).ToMsgStruct())
	}

	switch fromGroup {
	case data.Config.AdminID:
		// 当前版本，管理群和游戏群收到的命令不做区分
		fallthrough
	case data.Config.GroupID:
		if syntax.GroupMsg(fromQQ, Msg, ret) {
			return Intercept
		}
	}
	return Ignore
}

// 群成员减少事件
func onGroupMemberDecrease(fromGroup, fromQQ, beingOperateQQ int64) bool {
	retValue := Ignore
	ret := func(resp string) {
		b.SendGroupMsg(fromGroup, message.New().Text(resp))
		retValue = Intercept
	}
	// 尝试删白名单
	if fromGroup == data.Config.GroupID {
		whitelist.RemoveWhitelist(beingOperateQQ, ret)
	}
	return retValue
}

// 入群请求事件
func onGroupRequest(request *bot.GroupRequest) bool {
	if !data.Config.DealWithGroupRequest.Enable || // 功能未启用
		request.GroupID != data.Config.GroupID { // 不是要管理的群
		return Ignore
	}
	Logger := log.NewLogger("DwGR")
	for _, name := range regexp.MustCompile(`[0-9A-Za-z_]{3,16}`).FindAllString(request.Comment, 3) {
		name, id, err := whitelist.GetUUID(name)
		if err != nil {
			Logger.Infof("处理%d的入群请求，检查游戏名失败: %v", request.UserID, err)
			continue
		}
		if ok, err := checkRequest(name, id); err != nil {
			Logger.Errorf("服务器检查出错: %v", err)
			return Ignore
		} else if !ok {
			Logger.Infof("服务器拒绝%d作为%s入群: %v", request.UserID, name, err)
			if data.Config.DealWithGroupRequest.CanReject {
				request.Reject("")
				return Intercept
			}
			return Ignore
		}
		Logger.Infof("允许%d作为%s入群", request.UserID, name)
		request.Agree()

		ret := func(resp string) { b.SendGroupMsg(request.GroupID, message.New().Text(resp)) }
		whitelist.MyID(request.UserID, name, ret)
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
	Ignore    = false //忽略消息
	Intercept = true  //拦截消息
)
