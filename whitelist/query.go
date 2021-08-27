package whitelist

import (
	"fmt"
	"github.com/BaiMeow/SimpleBot/message"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
	"regexp"
	"strconv"
)

var expQQ = regexp.MustCompile(`^[0-9]{6,11}$`)          // 匹配一个QQ
var expName = regexp.MustCompile(`^[0-9A-Za-z_]{3,16}$`) // 匹配一个玩家名

func Info(args message.Msg, fromQQ int64, ret func(string)) bool {
	// 找出当前想查询的人的QQ
	switch len(args) {
	case 1:
		qqInfo(fromQQ, ret)
		return true
	case 2:
		switch args[1].GetType() {
		case "text":
			txt := args[1].(message.Text).Text
			if expQQ.MatchString(txt) {
				qq, err := strconv.ParseInt(txt, 10, 64)
				if err != nil {
					Logger.Waringf("%v", err)
					return false
				}
				qqInfo(qq, ret)
			} else if expName.MatchString(txt) {
				nameInfo(txt, ret)
			}
		case "at":
			id := args[1].(message.At).ID
			if id == "all" {
				return false
			}
			qq, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				Logger.Waringf("%v", err)
				return false
			}
			qqInfo(qq, ret)
			return true
		}
	}
	return false
}

func qqInfo(targetQQ int64, ret func(string)) {
	// 查询本人的绑定
	ID, err := data.GetWhitelistByQQ(targetQQ)
	if err != nil {
		Logger.Errorf("读取玩家绑定的ID出错: %v", err)
		ret("数据库查询失败惹(つД`)ノ")
		return
	}
	if ID == uuid.Nil {
		ret("这个还没有绑定白名单呢")
		return
	}

	// 根据UUID找到名字
	name, err := getName(ID)
	if err != nil {
		ret("游戏名查询失败惹(つД`)ノ")
		return
	}
	ret(name)
}

func nameInfo(targetName string, ret func(string)) {
	name, id, err := GetUUID(targetName)
	if err != nil {
		Logger.Errorf("查询UUID失败: %v", err)
		ret("查无此人")
		return
	}

	qq, err := data.GetWhitelistByUUID(id)
	if err != nil {
		Logger.Errorf("数据库查询QQ失败: %v", err)
		ret("数据库出问题了(つД`)ノ")
		return
	}

	if qq == 0 {
		ret(fmt.Sprintf("没人绑定%s哟~", name))
	} else {
		ret(fmt.Sprintf("啊呐占用%s的是%d哟", name, qq))
	}
}
