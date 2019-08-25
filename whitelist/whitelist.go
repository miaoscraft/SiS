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
		cqp.AddLog(cqp.Error, "myid", fmt.Sprintf("向Mojang查询玩家UUID失败: %v", err))
	}

	owner, oldName, err := data.SetWhitelist(qq, name, id)
	if err != nil {
		cqp.AddLog(cqp.Error, "myid", fmt.Sprintf("数据库操作失败: %v", err))
	}

	if owner != qq {
		ret(fmt.Sprintf("账号%q当前被[CQ:at,qq=%d]占有", name, owner))
	} else {
		// 删除旧的白名单
		if oldName != "" {
			err := data.RemoveWhitelist(oldName)
			if err != nil {
				ret(fmt.Sprintf("消除白名单%s失败: %v", oldName, err))
				return
			}
		}

		// 添加白名单
		err := data.AddWhitelist(name)
		if err != nil {
			ret(fmt.Sprintf("添加白名单%s失败: %v", oldName, err))
			return
		}

		ret(fmt.Sprintf("已为您添加白名单: %s", name))
	}
}

// getUUID 查询玩家的UUID
func getUUID(name string) (string, uuid.UUID, error) {
	var id uuid.UUID

	resp, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + name)
	if err != nil {
		return name, id, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&struct {
		Name *string
		ID   *uuid.UUID
	}{&name, &id})
	if err != nil {
		return name, id, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("服务器状态码非200: %v", resp.Status)
	}

	return name, id, err
}
