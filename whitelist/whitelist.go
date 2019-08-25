package whitelist

import (
	"encoding/json"
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
	"net/http"
)

func MyID(qq int64, name string, ret func(msg string)) {
	name, id, err := getUUID(name)
	if err != nil {
		cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("向Mojang查询玩家UUID失败: %v", err))
		return
	}

	owner, oldName, err := data.SetWhitelist(qq, name, id)
	if err != nil {
		cqp.AddLog(cqp.Error, "MyID", fmt.Sprintf("数据库操作失败: %v", err))
		return
	}

	// 若owner是当前处理的用户则说明绑定成功，否则就是失败
	if owner != qq {
		ret(fmt.Sprintf("账号%q当前被[CQ:at,qq=%d]占有", name, owner))
	} else {
		// 删除旧的白名单
		if oldName != nil {
			err := data.RemoveWhitelist(*oldName)
			if err != nil {
				ret(fmt.Sprintf("消除白名单%s失败: %v", *oldName, err))
				return
			}
		}

		// 添加白名单
		err := data.AddWhitelist(name)
		if err != nil {
			ret(fmt.Sprintf("添加白名单%s失败: %v", *oldName, err))
			return
		}

		ret(fmt.Sprintf("已为您添加白名单: %s", name))
	}
}

// getUUID 查询玩家的UUID
func getUUID(name string) (string, uuid.UUID, error) {
	var id uuid.UUID

	// 构造请求
	request, err := http.NewRequest("GET", "https://api.mojang.com/users/profiles/minecraft/"+name, nil)
	if err != nil {
		return name, id, err
	}

	// Golang默认的User-agent被屏蔽了
	request.Header.Set("User-agent", "SiS")

	// 发送Get请求
	resp, err := new(http.Client).Do(request)
	if err != nil {
		return name, id, err
	}
	defer resp.Body.Close()

	// 解析json返回值
	err = json.NewDecoder(resp.Body).Decode(&struct {
		Name *string
		ID   *uuid.UUID
	}{&name, &id})
	if err != nil {
		return name, id, err
	}

	// 检查返回码
	if resp.StatusCode != 200 {
		err = fmt.Errorf("服务器状态码非200: %v", resp.Status)
	}

	return name, id, err
}
