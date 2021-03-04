package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/customize"
	"github.com/miaoscraft/SiS/data"
	"github.com/miaoscraft/SiS/log"
	"github.com/miaoscraft/SiS/syntax"
	"github.com/miaoscraft/SiS/whitelist"
	"net/http"
	"net/url"
	"regexp"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

// cqp: 名称: SiS
// cqp: 版本: 1.4.1:1
// cqp: 作者: Tnze
// cqp: 简介: Minecraft服务器综合管理器
func main() {
	onStart()
	zero.OnMessage(zero.OnlyGroup).
		Handle(func(ctx *zero.Ctx) {
			ret := func(resp string) {
				ctx.Send(resp)
			}
			// 将登录账号载入命令解析器（用于识别@）
			syntax.CmdPrefix = fmt.Sprintf("[CQ:at,qq=%d]", ctx.Event.SelfID)
			onGroupMsg(
				ctx.Event.GroupID,
				ctx.Event.Sender.ID,
				ctx.Event.Message.CQString(),
				ret,
			)
		})
	zero.OnNotice(func(ctx *zero.Ctx) bool {
		return ctx.Event.PostType == "notice" && ctx.Event.NoticeType == "group_decrease"
	}).Handle(func(ctx *zero.Ctx) {
		ret := func(resp string) {
			ctx.Send(resp)
		}
		onGroupMemberDecrease(ctx.Event.GroupID, ctx.Event.OperatorID, ctx.Event.UserID, ret)
	})
	zero.OnRequest().Handle(func(ctx *zero.Ctx) {
		ret := func(resp string) {
			ctx.Send(resp)
		}
		onGroupRequest(ctx, ctx.Event.SubType, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.Comment, ctx.Event.Flag, ret)
	})

	zero.Run(zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []string{"123456"},
		Driver: []zero.Driver{
			driver.NewWebSocketClient(
				data.Config.ZeroBot.Host,
				data.Config.ZeroBot.Port,
				data.Config.ZeroBot.AccessToken,
			),
		},
	})
	select {}
	onStop()
}

func init() {
	//cqp.AppID = "cn.miaoscraft.sis"

	customize.Logger = log.NewLogger("Cstm")
	whitelist.Logger = log.NewLogger("MyID")
	data.Logger = log.NewLogger("Data")
}

var Logger = log.NewLogger("Main")

// 插件生命周期开始
func onStart() int32 {
	// 连接数据源
	err := data.Init(".")
	if err != nil {
		Logger.Errorf("初始化数据源失败: %v", err)
	}
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
func onGroupMsg(fromGroup, fromQQ int64, msg string, ret func(resp string)) int32 {
	if fromQQ == 80000000 { // 忽略匿名
		return Ignore
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
func onGroupMemberDecrease(fromGroup, fromQQ, beingOperateQQ int64, ret func(resp string)) int32 {
	retValue := Ignore
	ret1 := func(resp string) {
		ret(resp)
		retValue = Intercept
	}
	// 尝试删白名单
	if fromGroup == data.Config.GroupID {
		whitelist.RemoveWhitelist(beingOperateQQ, ret1)
	}
	return retValue
}

// 入群请求事件
func onGroupRequest(ctx *zero.Ctx, subType string, fromGroup, fromQQ int64, msg, respFlag string, ret func(resp string)) int32 {
	if !data.Config.DealWithGroupRequest.Enable || // 功能未启用
		fromGroup != data.Config.GroupID { // 不是要管理的群 {
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
				ctx.SetGroupAddRequest(respFlag, subType, Deny, "")
				return Intercept
			}
			return Ignore
		}
		Logger.Infof("允许%d作为%s入群", fromQQ, name)
		ctx.SetGroupAddRequest(respFlag, subType, Allow, "")

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

	Allow = true  // 允许进群
	Deny  = false // 拒绝进群
)
