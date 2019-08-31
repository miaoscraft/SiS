// Package customize 提供自定义指令的实现
package customize

import (
	"fmt"
	"strconv"

	"github.com/miaoscraft/SiS/data"
)

// 检查命令是否匹配一个自定义命令，若是的话则丢到RCON执行
// args长度必须大于0
func Exec(args []string, fromQQ int64, ret func(string)) bool {
	cmds, ok := data.Config.Cmd[args[0]]
	if !ok {
		return false
	}

	// 获取权限
	level, err := data.GetLevel(fromQQ)
	if err != nil {
		Logger.Errorf("获取权限出错: %v", err)
		ret("当前没有办法验证权限呢")
		return false
	}
	// 权限确认
	if cmds.Level <= level {
		Logger.Infof("成员%d以等级%d执行指令%q", fromQQ, level, cmds.Command)
		// 执行指令
		resp, err := data.RCONCmd(cmds.Command)
		if err != nil {
			Logger.Errorf("执行命令出错: %v", err)
			ret("服务器被玩坏啦？！")
		}

		// 返回结果
		if !cmds.Silent {
			ret(resp)
		}
		return true

	} else {
		//权限不足
		ret("你不能够执行这个命令哦～")
		return false
	}
}

func Auth(args []string, fromQQ int64, ret func(string)) bool {
	// args: ["auth", "@Member" | "QQ-num", "level"]
	if len(args) != 3 {
		return false
	}
	var target, level int64
	var err error
	// 解析目标QQ
	if _, err = fmt.Sscanf(args[1], "[CQ:at,qq=%d]", &target); err != nil {
	} else if target, err = strconv.ParseInt(args[1], 10, 64); err != nil {
	} else {
		return false
	}
	// 解析权限等级
	if level, err = strconv.ParseInt(args[2], 10, 64); err != nil {
		return false
	}
	// 确认是否是超级管理员
	for _, v := range data.Config.Administrators {
		if v == fromQQ {
			// 该用户属于最高管理员
			Logger.Infof("将%d的权限设置为%d", target, level)
			err = data.SetLevel(target, level)
			if err != nil {
				Logger.Errorf("设置权限出错: %v", err)
				ret("这里出现了问题(つД`)ノ")
			} else {
				ret("成功了( ̀⌄ ́)")
			}
			return true
		}
	}
	return false
}

var Logger interface {
	Error(str string)
	Errorf(format string, args ...interface{})

	Waring(str string)
	Waringf(format string, args ...interface{})

	Info(str string)
	Infof(format string, args ...interface{})

	Debug(str string)
	Debugf(format string, args ...interface{})
}
