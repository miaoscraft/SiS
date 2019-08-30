// Package customize 提供自定义指令的实现
package customize

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/miaoscraft/SiS/data"
)

// 检查命令是否匹配一个自定义命令，若是的话则丢到RCON执行
// args长度必须大于0
func Exec(cmd string, args []string, fromQQ int64, ret func(string)) bool {
	cmds, ok := data.Config.Cmd[args[0]]
	if !ok {
		return false
	}

	// 获取权限
	level, err := data.GetLevel(fromQQ)
	if err != nil {
		cqp.AddLog(cqp.Error, "Cmds", fmt.Sprintf("获取权限出错: %v", err))
		ret("当前没有办法验证权限呢")
		return false
	}
	// 权限确认
	if cmds.Level <= level {
		// 执行指令
		resp, err := data.RCONCmd(cmds.Command)
		if err != nil {
			cqp.AddLog(cqp.Error, "Cmds", fmt.Sprintf("执行命令出错: %v", err))
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
