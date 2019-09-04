package whitelist

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
	"strconv"
)

func Info(args []string, fromQQ int64, ret func(string)) bool {
	// 找出当前想查询的人的QQ
	var target int64
	switch len(args) {
	case 1:
		target = fromQQ
	case 2:
		var qq int64
		if _, err := fmt.Sscanf(args[1], "[CQ:at,qq=%d]", &qq); err == nil {
			target = qq
		} else if qq, err = strconv.ParseInt(args[1], 10, 64); err == nil {
			target = qq
		}
	default:
		return false
	}

	// 查询本人的绑定
	ID, err := data.GetWhitelist(target)
	if err != nil {
		Logger.Errorf("读取玩家绑定的ID出错: %v", err)
		ret("数据库查询失败惹(つД`)ノ")
		return true
	}
	if ID == uuid.Nil {
		ret("这个还没有绑定白名单呢")
		return true
	}

	// 根据UUID找到名字
	name, err := getName(ID)
	if err != nil {
		ret("游戏名查询失败惹(つД`)ノ")
		return true
	}
	ret(name)

	return true
}
